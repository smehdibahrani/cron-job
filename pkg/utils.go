package pkg

import (
	"cron_job/internal/entity"
	"cron_job/pkg/exception"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/sethvargo/go-password/password"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetUserIdFromReq(ctx *gin.Context) uint {
	value, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"exception": "user id not found"})
		return 0
	}
	userId := value.(uint)
	return userId
}

func GetUserFromReq(ctx *gin.Context) entity.User {
	value, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"exception": "user not found"})
		return entity.User{}
	}
	user, ok := value.(entity.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"exception": "user type assertion failed"})
		return entity.User{}
	}
	return user
}

func GetIntParam(ctx *gin.Context, key string) uint {
	id := ctx.Param(key)
	val, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"exception": "parameter " + key + " not found"})
	}
	return uint(val)
}

func HandleError(c *gin.Context, appErr *exception.AppError) {
	if appErr != nil {
		c.JSON(appErr.Code, gin.H{
			"message":    appErr.Message,
			"detail":     appErr.Detail,
			"statusCode": appErr.Code,
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"message":    "internal server exception",
		"statusCode": http.StatusInternalServerError,
	})
	log.Printf("Unexpected exception : %v", appErr)
	return
}

func GenerateRandomSixDigitNumber() string {
	rand.Seed(time.Now().UnixNano())     // Seed the random number generator
	number := rand.Intn(900000) + 100000 // Generate a number between 100000 and 999999
	return strconv.Itoa(number)
}

func GeneratePassword(length int) (string, error) {
	return password.Generate(length, length-7, 0, false, false)
}

func IsValidCron(cronStr string) bool {
	_, err5 := cron.ParseStandard(cronStr)
	if err5 == nil {
		return true
	}
	return false
}

func ParseIntervalMinutes(cronStr string) uint {
	// Parse the cron expression
	schedule, err := cron.ParseStandard(cronStr)
	if err != nil {
		return 0 // Return error if parsing fails
	}

	// Get the next two execution times to calculate the interval
	now := time.Now()
	nextExecutionTime := schedule.Next(now)
	secondNextExecutionTime := schedule.Next(nextExecutionTime)

	// Calculate the difference in minutes between the two execution times
	durationUntilNext := secondNextExecutionTime.Sub(nextExecutionTime)
	intervalMinutes := uint(durationUntilNext.Minutes())
	fmt.Println("intervalMinutes", intervalMinutes)
	return intervalMinutes
}

func ParseExecutionPerDay(cronStr string) uint {
	// Parse the cron string
	_, err := cron.ParseStandard(cronStr)
	if err != nil {
		return 0 // Return error if parsing fails
	}

	// Define how many executions will occur
	executionCount := 0

	// Get the components of the cron expressions
	fields := strings.Split(cronStr, " ")
	if len(fields) < 5 {
		return 0
	}

	// Extract minute and hour components
	minuteField := fields[0]
	hourField := fields[1]

	// Calculate executions based on minutes and hours
	minuteCount := countOccurrences(minuteField, 60) // Can execute max 60 times per hour
	hourCount := countOccurrences(hourField, 24)     // Can execute max 24 times per day

	// Total executions per day is the hour count multiplied by the number of minutes
	executionCount = minuteCount * hourCount

	return uint(executionCount)
}

func countOccurrences(field string, max int) int {
	if field == "*" {
		return max // Executes every unit of time (i.e., every minute/hour)
	}

	// Split the field into parts (e.g., "*/5", "0-10", "1,5", etc.)
	parts := strings.Split(field, ",")
	count := 0
	for _, part := range parts {
		if strings.Contains(part, "/") {
			// Handle intervals (e.g., "*/5" means every 5 units)
			interval := strings.Split(part, "/")[1]
			intVal, _ := strconv.Atoi(interval)
			count += max / intVal // Executions for intervals
		} else if strings.Contains(part, "-") {
			// Handle ranges (e.g., "1-5")
			rangeParts := strings.Split(part, "-")
			start, _ := strconv.Atoi(rangeParts[0])
			end, _ := strconv.Atoi(rangeParts[1])
			count += (end - start + 1) // Range execution count
		} else {
			// It's a specific minute or hour
			count++
		}
	}
	return count
}

func CalculateEpd(cronExpr string) uint {
	epd := ParseExecutionPerDay(cronExpr)
	fmt.Println("epd", epd)
	return epd
}

func MapToDto(c *gin.Context, dto interface{}) *exception.AppError {
	if err := c.ShouldBindJSON(&dto); err != nil {
		return exception.NewBadRequest(err.Error(), "")
	}
	return nil
}
