package redis

import (
	"context"
	"cron_job/internal/entity"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var Rdb *redis.Client
var ctx = context.Background()

func Init() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
}
func SetValueWithExpireTime(key string, value interface{}, expiration time.Duration) {
	Rdb.Set(ctx, key, value, expiration)
}

func GetValue(key string) string {
	result := Rdb.Get(ctx, key)
	return result.Val()
}

func AddJob(job entity.Job, nextTime time.Time) {
	res := Rdb.ZAdd(ctx, "jobs", &redis.Z{
		Score:  float64(nextTime.Unix()),
		Member: job.ID,
	})
	fmt.Println(res.String())
}

func GetJobs() []string {
	jobIds, _ := Rdb.ZRangeByScore(ctx, "jobs", &redis.ZRangeBy{
		Min: "0",
		Max: fmt.Sprintf("%d", time.Now().Unix()),
	}).Result()

	return jobIds
}

func RemoveJob(jobId uint) {
	Rdb.ZRem(ctx, "jobs", jobId)
}

func SetJobFailureTimes(jobId uint) uint {
	key := fmt.Sprintf("jobs-failures:%d", jobId)
	value := Rdb.Incr(ctx, key)
	return uint(value.Val())
}

func RemoveJobFailureTimes(jobId uint) {
	key := fmt.Sprintf("jobs-failures:%d", jobId)
	Rdb.Set(ctx, key, 0, 0)
}
