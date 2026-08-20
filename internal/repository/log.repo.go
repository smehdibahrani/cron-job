package repository

import (
	"cron_job/internal/config"
	"cron_job/internal/entity"
	"fmt"
	"gorm.io/gorm"
)

type LogRepository struct {
	db *gorm.DB
}

func NewLogRepository() *LogRepository {
	return &LogRepository{db: config.DB}
}

func (r *LogRepository) Create(log *entity.Log) (entity.Log, error) {
	err := r.db.Create(&log).Error
	return *log, err
}

func (r *LogRepository) GetAllByUserId(userId uint) ([]entity.Log, error) {
	var logs []entity.Log
	err := r.db.Where("user_id = ?", userId).Find(&logs).Error
	return logs, err
}

func (r *LogRepository) GetAllByJobIdAndUserId(jobId uint, userId uint) ([]entity.Log, error) {
	var logs []entity.Log
	err := r.db.Preload("RequestHttp").
		Where("job_id = ? and user_id = ?", jobId, userId).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

func (r *LogRepository) Clean() {
	res := r.db.Raw(`WITH user_log_limits AS (
    SELECT
        up.user_id,
        COALESCE(p.save_logs, 50) AS max_logs
    FROM user_plans up
    INNER JOIN plans p ON up.plan_id = p.id
    WHERE up.is_active = true
),
ranked_logs AS (
    SELECT
        l.id,
        l.user_id,
        ROW_NUMBER() OVER (PARTITION BY l.user_id ORDER BY l.created_at DESC) AS rn
    FROM logs l
    inner JOIN user_log_limits as ull ON l.user_id = ull.user_id
),
logs_to_delete AS (
    SELECT rl.id
    FROM ranked_logs rl
    INNER JOIN user_log_limits ull ON rl.user_id = ull.user_id
    WHERE rl.rn > ull.max_logs
)

-- Debugging: Check what would be deleted
-- SELECT * FROM logs_to_delete;

-- Perform deletion
DELETE FROM logs
WHERE id IN (SELECT id FROM logs_to_delete);`).Error
	fmt.Println("log clean result:", res)
}
