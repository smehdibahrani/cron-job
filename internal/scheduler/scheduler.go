package scheduler

import (
	"bytes"
	"cron_job/internal/entity"
	"cron_job/internal/usecase"
	"cron_job/pkg/email"
	"cron_job/pkg/redis"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var jobUseCase *usecase.JobUseCase
var logUseCase *usecase.LogUseCase
var userUseCase *usecase.UserUseCase

func Init() {
	jobUseCase = usecase.NewJobUseCase()
	logUseCase = usecase.NewLogUseCase()
	userUseCase = usecase.NewUserUseCase()
}

type HeaderValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func AddToRedisQueue(job entity.Job) {
	fmt.Println("adding new job", job.ID)
	// convert cron to timestamp
	nextTime := getNextCronTime(job.Schedule)
	redis.AddJob(job, nextTime)
}

func RemoveFromRedisQueue(jobId uint) {
	fmt.Println("removing job", jobId)
	redis.RemoveJob(jobId)
}

func getNextCronTime(schedule string) time.Time {
	cronParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	scheduleEntry, _ := cronParser.Parse(schedule)
	now := time.Now()
	return scheduleEntry.Next(now)
}

func startHourlyLogCleanup() {
	// Create a ticker that fires every hour
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	// Run cleanup immediately on startup
	fmt.Println("Starting hourly log cleanup...")
	logUseCase.Clean()

	// Then run every hour
	for {
		select {
		case <-ticker.C:
			fmt.Println("Running hourly log cleanup...")
			logUseCase.Clean()
		}
	}
}

func StartScheduler() {
	// Start hourly log cleanup in a separate goroutine
	go startHourlyLogCleanup()

	// Main job execution loop
	for {
		jobIdsStr := redis.GetJobs()
		fmt.Println("jobIdsStr", jobIdsStr)
		for _, jobIdStr := range jobIdsStr {
			jobId, _ := strconv.Atoi(jobIdStr)
			go executeJob(uint(jobId))
		}
		time.Sleep(time.Minute)
	}
}

func executeJob(jobId uint) {
	fmt.Println("executing job")
	job, err := jobUseCase.GetById(jobId)
	if err != nil {
		fmt.Println("exception get job", err)
		return
	}
	redis.RemoveJob(job.ID)
	AddToRedisQueue(job)

	response, responseTime, err := callHttpApi(job.RequestHttp)
	if err != nil {
		fmt.Println("Error create Request job:", err)
		createLog(job, response, time.Duration(responseTime), err.Error())
		failureTimes := redis.SetJobFailureTimes(job.ID)
		handleNotification(job, true, failureTimes)
		jobUseCase.UpdateTotal(job.ID, job.TotalSuccess, job.TotalFail+1)
		return
	}

	createLog(job, response, time.Duration(responseTime), "")
	fmt.Println("Executed job:", job.Name, response.Status)
	handleNotification(job, false, 0)
	jobUseCase.UpdateTotal(job.ID, job.TotalSuccess+1, job.TotalFail)
}

func createLog(job entity.Job, rsp *http.Response, resTime time.Duration, errMsg string) {
	responseBody := ""
	statusCode := 0
	if rsp != nil {
		body, err := io.ReadAll(rsp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			return
		}
		responseBody = string(body)
		statusCode = rsp.StatusCode
		defer rsp.Body.Close()
	}
	logData := entity.Log{
		RequestHttpId: job.RequestHttp.ID,
		JobId:         job.ID,
		UserId:        job.UserId,
		Method:        string(job.RequestHttp.Method),
		Url:           job.RequestHttp.Url,
		Res:           responseBody,
		ResStatus:     uint(statusCode),
		ResTime:       uint(resTime.Seconds()),
		ErrorMessage:  errMsg,
	}
	_, err := logUseCase.Create(&logData)
	if err != nil {
		fmt.Println("Error create log:", err)
	}

}

func handleNotification(job entity.Job, isFailure bool, failureTimes uint) {
	user, _ := userUseCase.GetById(job.UserId)
	if user.Email == "" {
		return
	}
	if len(job.Notifications) > 0 {
		notif := job.Notifications[0]
		if notif.Type == "email" {
			if !isFailure && notif.Action == "after_each_job_execution" {
				go email.SendJobSuccessExecutedEmail(user.Email, job.ID, job.Name)
			} else if notif.Action == "job_failing" {
				if notif.Sensitivity != nil && *notif.Sensitivity <= failureTimes {
					go email.SendJobFailureEmail(user.Email, job.ID, job.Name)
					redis.RemoveJobFailureTimes(job.ID)
				}
			}
		}
	}
}

func callHttpApi(reqHttp entity.RequestHttp) (*http.Response, uint, error) {
	startTime := time.Now()
	req, err := http.NewRequest(strings.ToUpper(string(reqHttp.Method)), reqHttp.Url, bytes.NewBuffer(reqHttp.Body))
	if err != nil {
		fmt.Println("Error create Request job:", err)
		return nil, 0, err
	}

	var headers []HeaderValue
	err = json.Unmarshal(reqHttp.Headers, &headers)
	if err != nil {
		fmt.Println("Error Unmarshal headers:", err)
		return nil, 0, err
	}
	for _, header := range headers {
		req.Header.Set(header.Key, header.Value)
	}
	client := &http.Client{}
	client.Timeout = time.Duration(reqHttp.TimeOut) * time.Second
	rsp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error executing job:", err)
		return nil, uint(time.Since(startTime)), err
	}

	return rsp, uint(time.Since(startTime)), nil

}
