package controller

import (
	"play-wails/internal/model"
	"play-wails/internal/service"

	"github.com/google/uuid"
)

/*
 * WorkSessionController は作業時間計測のオーケストレーションを行う
 * フロントからのリクエストを受け、Service を呼び出して結果を返す
 */
type WorkSessionController struct {
	workSessionService *service.WorkSessionService
}

/*
 * 実装クラスのインスタンス生成
 *
 * @param workSessionService 作業セッションサービス
 * @return インスタンス
 */
func NewWorkSessionController(workSessionService *service.WorkSessionService) *WorkSessionController {
	return &WorkSessionController{workSessionService: workSessionService}
}

/*
 * 計測を開始し、新規作業セッションを返す
 *
 * @param taskID タスクID（UUID文字列）
 * @return 作業セッション, エラー
 */
func (c *WorkSessionController) Start(taskID string) (*model.WorkSession, error) {
	// タスクIDをUUIDに変換
	id, err := uuid.Parse(taskID)
	if err != nil {
		return nil, err
	}

	// 作業セッションを開始
	return c.workSessionService.Start(id)
}

/*
 * 指定した作業セッションを停止する
 *
 * @param sessionID 作業セッションID（UUID文字列）
 * @return エラー
 */
func (c *WorkSessionController) Stop(sessionID string) error {
	// 作業セッションIDをUUIDに変換
	id, err := uuid.Parse(sessionID)
	if err != nil {
		return err
	}
	return c.workSessionService.Stop(id)
}

/*
 * 同一計測として作業を再開し、新規作業セッションを返す
 *
 * @param taskID タスクID（UUID文字列）
 * @param runID 計測実行のグループID（UUID文字列）
 * @return 作業セッション, エラー
 */
// taskID: タスクID（UUID文字列）, runID: 計測実行のグループID（UUID文字列）
func (c *WorkSessionController) Resume(taskID string, runID string) (*model.WorkSession, error) {
	// タスクIDをUUIDに変換
	tid, err := uuid.Parse(taskID)
	if err != nil {
		return nil, err
	}

	// 計測実行のグループIDをUUIDに変換
	rid, err := uuid.Parse(runID)
	if err != nil {
		return nil, err
	}

	// 作業セッションを再開
	return c.workSessionService.Resume(tid, rid)
}

/*
 * 指定した作業セッションを取得する
 *
 * @param sessionID 作業セッションID（UUID文字列）
 * @return 作業セッション, エラー
 */
func (c *WorkSessionController) Current(sessionID string) (*model.WorkSession, error) {

	// 作業セッションIDをUUIDに変換
	id, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, err
	}

	// 作業セッションを取得
	return c.workSessionService.Current(id)
}

/*
 * 計測を完了し、累計時間でTimeRecordを1件作成
 *
 * @param runID 計測実行のグループID（UUID文字列）
 * @return 作成したTimeRecord, エラー
 */
func (c *WorkSessionController) Complete(runID string) (*model.TimeRecord, error) {

	// 計測実行のグループIDをUUIDに変換
	id, err := uuid.Parse(runID)
	if err != nil {
		return nil, err
	}

	// 計測を完了し、累計時間でTimeRecordを1件作成
	return c.workSessionService.Complete(id)
}
