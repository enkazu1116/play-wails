package controller

import (
	"play-wails/internal/model"
	"play-wails/internal/service"

	"github.com/google/uuid"
)

type TimeRecordController struct {
	timeRecordService *service.TimeRecordService
}

/*
 * 実装クラスのインスタンス生成
 *
 * @param timeRecordService 時間計測レコードサービス
 * @return インスタンス
 */
func NewTimeRecordController(timeRecordService *service.TimeRecordService) *TimeRecordController {
	return &TimeRecordController{timeRecordService: timeRecordService}
}

/*
 * 計測結果一覧を取得する
 * 論理削除済みは除外
 *
 * @return 計測結果一覧, エラー
 */
func (c *TimeRecordController) List() ([]*model.TimeRecord, error) {
	return c.timeRecordService.List()
}

/*
 * 指定IDの計測結果を取得する
 *
 * @param id 計測結果ID（UUID文字列）
 * @return 計測結果, エラー
 */
func (c *TimeRecordController) Get(id string) (*model.TimeRecord, error) {
	// 計測結果IDをUUIDに変換
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return c.timeRecordService.Get(uid)
}

/*
 * 計測結果の時間を変更する
 *
 * @param record 計測結果
 * @return エラー
 */
func (c *TimeRecordController) Update(record *model.TimeRecord) error {
	return c.timeRecordService.Update(record)
}

/*
 * 計測結果を論理削除する
 *
 * @param id 計測結果ID（UUID文字列）
 * @return エラー
 */
func (c *TimeRecordController) Delete(id string) error {

	// 計測結果IDをUUIDに変換
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return c.timeRecordService.Delete(uid)
}
