package model

import (
	"time"

	"github.com/google/uuid"
)

/*
 * 時間計測
 * 完了時に1件作成。同一のRunIDをWorkSessionグループの累計時間を保持
 */
type TimeRecord struct {
	ID         uuid.UUID     `json:"id"`
	RunID      uuid.UUID     `json:"run_id"`
	TaskID     uuid.UUID     `json:"task_id"`
	DeleteFlag bool          `json:"delete_flag"`
	StartTime  time.Time     `json:"start_time"`
	EndTime    time.Time     `json:"end_time"`
	Duration   time.Duration `json:"duration"`
}
