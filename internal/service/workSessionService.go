package service

import (
	"errors"
	"play-wails/internal/model"
	"play-wails/internal/repository"
	"time"

	"github.com/google/uuid"
)

type WorkSessionService struct {
	wrepo repository.WorkSessionRepository
	trepo repository.TimeRecordRepository
}

func NewWorkSessionService(wrepo repository.WorkSessionRepository, trepo repository.TimeRecordRepository) *WorkSessionService {
	return &WorkSessionService{wrepo: wrepo, trepo: trepo}
}

/*
 * 新規作業セッションを生成し、作業を開始する（新規 RunID を発行）
 *
 * @param taskID タスクID
 * @return 作業セッション（RunID を含む）, エラー
 */
func (s *WorkSessionService) Start(taskID uuid.UUID) (*model.WorkSession, error) {
	runID := uuid.New()
	session := &model.WorkSession{
		ID:        uuid.New(),
		RunID:     runID,
		TaskID:    taskID,
		StartTime: time.Now(),
		EndTime:   nil,
	}

	if err := s.wrepo.Create(session); err != nil {
		return nil, err
	}

	return session, nil
}

/*
 * 同一計測実行として作業を再開する（既存 RunID で新規 WorkSession を作成）
 *
 * @param taskID タスクID
 * @param runID 計測実行のグループID（開始時に発行した RunID）
 * @return 作業セッション, エラー
 */
func (s *WorkSessionService) Resume(taskID uuid.UUID, runID uuid.UUID) (*model.WorkSession, error) {
	session := &model.WorkSession{
		ID:        uuid.New(),
		RunID:     runID,
		TaskID:    taskID,
		StartTime: time.Now(),
		EndTime:   nil,
	}

	if err := s.wrepo.Create(session); err != nil {
		return nil, err
	}

	return session, nil
}

/*
 * 作業セッションを停止する
 *
 * @param sessionID 作業セッションID
 * @return エラー
 * @error エラー
 */
func (s *WorkSessionService) Stop(sessionID uuid.UUID) error {
	session, err := s.wrepo.FindByID(sessionID)
	if err != nil {
		return err
	}

	if err := session.Stop(time.Now()); err != nil {
		return err
	}

	return s.wrepo.Update(session)
}

/*
 * 作業セッションを取得する
 *
 * @param sessionID 作業セッションID
 * @return 作業セッション, エラー
 */
func (s *WorkSessionService) Current(sessionID uuid.UUID) (*model.WorkSession, error) {
	return s.wrepo.FindByID(sessionID)
}

/*
 * 計測を完了し、同一 RunID の全 WorkSession の累計で TimeRecord を1件作成する
 *
 * @param runID 計測実行のグループID
 * @return 作成した TimeRecord, エラー
 */
func (s *WorkSessionService) Complete(runID uuid.UUID) (*model.TimeRecord, error) {
	sessions, err := s.wrepo.ListByRunID(runID)
	if err != nil {
		return nil, err
	}
	if len(sessions) == 0 {
		return nil, errors.New("該当する作業セッションがありません")
	}

	var total time.Duration
	var firstStart, lastEnd time.Time
	var taskID uuid.UUID

	// 作業セッションを処理
	for i, sess := range sessions {
		// 最初の作業セッションの開始時刻とタスクIDを取得
		if i == 0 {
			firstStart = sess.StartTime
			taskID = sess.TaskID
		}

		// 作業セッションが停止されていない場合はエラー
		if sess.EndTime == nil {
			return nil, errors.New("未停止のセッションがあります。完了前に停止してください")
		}

		// 作業セッションの時間を累計
		total += sess.Duration()

		// 最後の作業セッションの終了時刻を更新
		if lastEnd.IsZero() || sess.EndTime.After(lastEnd) {
			lastEnd = *sess.EndTime
		}
	}

	// 時間計測モデルを作成
	record := &model.TimeRecord{
		ID:         uuid.New(),
		RunID:      runID,
		TaskID:     taskID,
		DeleteFlag: false,
		StartTime:  firstStart,
		EndTime:    lastEnd,
		Duration:   total,
	}

	// 時間計測レコードを作成
	if err := s.trepo.Create(record); err != nil {
		return nil, err
	}

	return record, nil
}
