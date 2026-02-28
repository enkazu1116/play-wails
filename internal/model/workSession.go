package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

/*
 * 作業セッション（1回の開始〜停止の履歴）
 * RunIDで同じ計測実行に属するセッションをグループ化する
 */
type WorkSession struct {
	ID        uuid.UUID  `json:"id"`
	RunID     uuid.UUID  `json:"run_id"`
	TaskID    uuid.UUID  `json:"task_id"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}

/*
 * 作業セッションが実行中か判定
 */
func (s WorkSession) IsRunning() bool {
	return s.EndTime == nil
}

/*
 * 作業セッションを停止する
 *
 * @param now 現在時刻
 */
func (s *WorkSession) Stop(now time.Time) error {
	if s.EndTime != nil {
		return errors.New("【ERROR】作業セッションは既に停止されています。")
	}
	if now.Before(s.StartTime) {
		return errors.New("【ERROR】作業セッションの開始時刻が過去の時刻です。")
	}

	s.EndTime = &now
	return nil
}

/*
 * セッションの作業時間（停止済みの場合のみ有効）
 *
 * @return 作業時間
 */
func (s WorkSession) Duration() time.Duration {
	if s.EndTime == nil {
		return 0
	}
	return s.EndTime.Sub(s.StartTime)
}
