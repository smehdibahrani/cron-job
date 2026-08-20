CREATE TYPE NotificationType AS ENUM ('email', 'webhook'); -- Create the enum type in PostgreSQL
CREATE TYPE NotificationAction AS ENUM ('job_failing', 'after_each_job_execution'); -- Create the enum type in PostgreSQL
CREATE TYPE method AS ENUM ('get', 'post','put','delete','patch','head','options'); -- Create the enum type in PostgreSQL
