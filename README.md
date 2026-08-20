# Cron Job

HTTP cron scheduler. Create jobs that call an HTTP endpoint on a cron schedule, keep execution logs, and get notified when a job succeeds or fails.

## Problem

Many backends need a request to run on a timer: a health check, a webhook, a cache refresh, a report, a retry. That usually means a server with crontab, a queue worker, or a third-party ping service.

This project is that scheduler as an API. You register, define the URL, method, headers, body, and cron expression, then the service:

- stores the job in PostgreSQL
- queues the next run in Redis
- executes the HTTP request when due
- records status, response body, and duration
- can email you on failure (after N consecutive fails) or after every run

You do not need to run crontab yourself. Jobs are grouped, can be paused, and belong to the authenticated user.

## How it works

1. GORM migrates tables on startup. Custom PostgreSQL enums must already exist (see [Database init](#database-init-initsql)).
2. Active jobs are stored in a Redis sorted set (`jobs`), scored by the next Unix timestamp.
3. A loop wakes every minute, pulls due job IDs, executes each HTTP request, writes a log, and re-queues the next run.
4. Execution logs older than the user’s plan limit are cleaned every hour.

## Requirements

- Go 1.24+
- PostgreSQL
- Redis (`localhost:6379`)

## Setup

### 1. Create the database

```bash
createdb cronJob
```

Use the same name as `DB_NAME` in `.env`.

### 2. Run init SQL

GORM can create tables, but not these PostgreSQL enum types. Run `init.sql` **once, before the first start**:

```bash
psql -d cronJob -f init.sql
```

If you skip this step, AutoMigrate fails when it tries to use `notificationtype`, `notificationaction`, and `method`.

### 3. Environment

Create a `.env` in the project root:

```env
# app
ADDR=8081
DEBUG_MODE=true

# postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=cronJob
DB_USER=postgres
DB_PASS=postgres
```

| Variable     | Default   | Description                          |
|--------------|-----------|--------------------------------------|
| `ADDR`       | `8080`    | HTTP port                            |
| `DEBUG_MODE` | `false`   | Gin debug vs release mode            |
| `DB_HOST`    | `localhost` | PostgreSQL host                    |
| `DB_PORT`    | `5432`    | PostgreSQL port                      |
| `DB_NAME`    | `cronjob` | Database name                        |
| `DB_USER`    | `postgres`| Database user                        |
| `DB_PASS`    | `postgres`| Database password                    |

Redis is currently hardcoded to `localhost:6379`.

### 4. Run

```bash
go mod download
go run ./cmd/main.go
```

The API listens on `http://localhost:<ADDR>` (default `8081` if set in `.env`).

Check that Postgres and Redis are reachable:

```bash
curl http://localhost:8081/api/app/health
```

## Database init (`init.sql`)

`init.sql` only creates enums. Tables (`users`, `jobs`, `groups`, `notifications`, `request_https`, `logs`) are created by GORM AutoMigrate on startup.

```sql
CREATE TYPE NotificationType AS ENUM ('email', 'webhook');
CREATE TYPE NotificationAction AS ENUM ('job_failing', 'after_each_job_execution');
CREATE TYPE method AS ENUM ('get', 'post','put','delete','patch','head','options');
```

| Type                 | Values                                      | Used for                                      |
|----------------------|---------------------------------------------|-----------------------------------------------|
| `NotificationType`   | `email`, `webhook`                          | How to notify after a job run                 |
| `NotificationAction` | `job_failing`, `after_each_job_execution`   | When to notify                                |
| `method`             | `get`, `post`, `put`, `delete`, `patch`, `head`, `options` | HTTP method of the scheduled request |

`job_failing` uses `sensitivity`: notify after that many consecutive failures. `after_each_job_execution` notifies on every successful run. Email is implemented; webhook is reserved on the type.

Run this only on a fresh database. Re-running `CREATE TYPE` on an existing DB will error if the types already exist.

## Auth

Protected routes expect:

```
Authorization: Bearer <accessToken>
```

Access tokens last 24 hours. Refresh tokens last 30 days and are sent with:

```
refresh-token: Bearer <refreshToken>
```

## APIs

Base URL: `/api`

### App (public)

| Method | Path              | Description                                      |
|--------|-------------------|--------------------------------------------------|
| GET    | `/api/app/version`| App version (`0.0.1`)                            |
| GET    | `/api/app/health` | Postgres + Redis ping. `200` if both are ok      |

### Auth (public, except refresh)

#### `POST /api/auth/register`

Creates a user. Password min length 8.

```json
{
  "firstName": "Mehdi",
  "lastName": "Dev",
  "email": "mehdi@example.com",
  "password": "secret123"
}
```

```json
{ "message": "user register success" }
```

#### `POST /api/auth/login`

```json
{
  "email": "mehdi@example.com",
  "password": "secret123"
}
```

```json
{
  "firstName": "Mehdi",
  "lastName": "Dev",
  "email": "mehdi@example.com",
  "createdAt": "...",
  "accessToken": "...",
  "refreshToken": "..."
}
```

#### `POST /api/auth/refresh-token`

Header: `refresh-token: Bearer <refreshToken>`

Returns a new `accessToken` and `refreshToken` (same shape as login).

### Users (JWT)

| Method | Path                         | Description              |
|--------|------------------------------|--------------------------|
| GET    | `/api/users/`                | Current user profile     |
| PUT    | `/api/users/`                | Update first/last name   |
| POST   | `/api/users/change-password` | Change password          |

**Update**

```json
{
  "firstName": "Mehdi",
  "lastName": "Updated"
}
```

**Change password** (`currentPassword` min 6, `newPassword` min 8)

```json
{
  "currentPassword": "secret123",
  "newPassword": "newSecret1"
}
```

User responses return `firstName`, `lastName`, `email`, `createdAt`.

### Groups (JWT)

Jobs belong to a group. A default group (`name: default`, `defGrp: true`) is used when a job is created without `groupId`. The default group cannot be deleted; its jobs are moved to default if you delete another group.

| Method | Path                    | Description                         |
|--------|-------------------------|-------------------------------------|
| POST   | `/api/groups/`          | Create group                        |
| GET    | `/api/groups/`          | List current user’s groups          |
| GET    | `/api/groups/:id`       | Get one group                       |
| GET    | `/api/groups/:id/jobs`  | Jobs in that group                  |
| PUT    | `/api/groups/:id`       | Update group                        |
| DELETE | `/api/groups/:id`       | Delete group (not the default one)  |

**Create**

```json
{
  "name": "production",
  "tagName": "prod",
  "description": "prod health checks"
}
```

**Response**

```json
{
  "id": 1,
  "name": "production",
  "tagName": "prod",
  "description": "prod health checks",
  "defGrp": false,
  "jobCount": 0,
  "createdAt": "..."
}
```

### Jobs (JWT)

Cron expressions are standard 5-field cron (`minute hour day-of-month month day-of-week`), e.g. `*/5 * * * *` every 5 minutes.

If `jobData.groupId` is `0` or omitted, the user’s default group is used.

An active job is pushed to the Redis queue. Inactive jobs are removed from the queue.

| Method | Path                 | Description                    |
|--------|----------------------|--------------------------------|
| POST   | `/api/jobs`          | Create job                     |
| GET    | `/api/jobs`          | List current user’s jobs       |
| GET    | `/api/jobs/:id`      | Get one job                    |
| GET    | `/api/jobs/:id/logs` | Execution logs for that job    |
| PUT    | `/api/jobs/:id`      | Update job, HTTP request, notif|
| DELETE | `/api/jobs/:id`      | Soft-delete job and unqueue it |

**Create**

```json
{
  "jobData": {
    "name": "API health check",
    "description": "Ping production /health",
    "schedule": "*/5 * * * *",
    "isActive": true,
    "groupId": 1
  },
  "httpRequest": {
    "url": "https://api.example.com/health",
    "method": "get",
    "headers": [{ "key": "Accept", "value": "application/json" }],
    "body": {},
    "timeOut": 30
  },
  "notification": {
    "type": "email",
    "action": "job_failing",
    "sensitivity": 3
  }
}
```

`notification` is optional.

| Field | Notes |
|-------|--------|
| `httpRequest.method` | One of: `get`, `post`, `put`, `delete`, `patch`, `head`, `options` |
| `httpRequest.headers` | JSON array of `{ "key", "value" }` |
| `httpRequest.body` | JSON body for POST/PUT/PATCH |
| `httpRequest.timeOut` | Seconds |
| `notification.type` | `email` or `webhook` |
| `notification.action` | `job_failing` or `after_each_job_execution` |
| `notification.sensitivity` | Consecutive failures before email (`job_failing`) |

**Job response**

```json
{
  "id": 1,
  "name": "API health check",
  "groupId": 1,
  "groupName": "",
  "schedule": "*/5 * * * *",
  "executionPerDay": 288,
  "totalSuccess": 0,
  "totalFail": 0,
  "isActive": true,
  "requestHttp": { "url": "...", "method": "get", "timeOut": 30 },
  "notification": {
    "type": "email",
    "action": "job_failing",
    "sensitivity": 3
  },
  "createdAt": "...",
  "updatedAt": "..."
}
```

**Logs** (`GET /api/jobs/:id/logs`) — newest first:

```json
[
  {
    "id": 10,
    "jobId": 1,
    "userId": 1,
    "url": "https://api.example.com/health",
    "method": "get",
    "res": "{\"status\":\"ok\"}",
    "resStatus": 200,
    "resTime": 1,
    "errorMessage": "",
    "createdAt": "..."
  }
]
```

`resTime` is in seconds. `errorMessage` is set when the HTTP call fails (timeout, network error).

## Project layout

```
cmd/main.go              entrypoint: DB, Redis, scheduler, Gin
init.sql                 PostgreSQL enum types (run once)
internal/config          env + GORM Postgres
internal/entity          User, Job, Group, Notification, RequestHttp, Log
internal/handler         HTTP handlers / DTOs
internal/route           /api/app, /auth, /users, /jobs, /groups
internal/scheduler       Redis queue + HTTP execution + notifications
internal/usecase         business rules
internal/repository      GORM access
pkg/jwt, pkg/redis, pkg/email
```

## Example flow

```bash
# register + login
curl -X POST http://localhost:8081/api/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"firstName":"Mehdi","lastName":"Dev","email":"mehdi@example.com","password":"secret123"}'

TOKEN=$(curl -s -X POST http://localhost:8081/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"mehdi@example.com","password":"secret123"}' | jq -r .accessToken)

# create a job
curl -X POST http://localhost:8081/api/jobs \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{
    "jobData": {"name":"ping","schedule":"*/5 * * * *","isActive":true},
    "httpRequest": {"url":"https://httpbin.org/get","method":"get","timeOut":15}
  }'
```
