package service

import (
	"play-wails/internal/model"
	"play-wails/internal/repository"

	"github.com/google/uuid"
)

type TimeRecordService struct {
	trepo repository.TimeRecordRepository
}

/*
 * 実装クラスのインスタンス生成
 *
 * @param trepo 時間計測レコードリポジトリ
 * @return インスタンス
 */
func NewTimeRecordService(trepo repository.TimeRecordRepository) *TimeRecordService {
	return &TimeRecordService{trepo: trepo}
}

/*
 * 計測結果一覧を取得する
 * 論理削除済みは除外
 *
 * @return 計測結果一覧, エラー
 */
func (s *TimeRecordService) List() ([]*model.TimeRecord, error) {
	return s.trepo.List(true)
}

/*
 * 指定IDの計測結果を取得する
 *
 * @param id 計測結果ID
 * @return 計測結果, エラー
 */
func (s *TimeRecordService) Get(id uuid.UUID) (*model.TimeRecord, error) {
	return s.trepo.FindByID(id)
}

/*
 * 計測結果の時間（開始・終了・作業時間）を更新する
 *
 * @param record 計測結果
 * @return エラー
 */
func (s *TimeRecordService) Update(record *model.TimeRecord) error {
	return s.trepo.Update(record)
}

/*
 * 計測結果を論理削除する
 *
 * @param id 計測結果ID
 * @return エラー
 */
func (s *TimeRecordService) Delete(id uuid.UUID) error {
	return s.trepo.Delete(id)
}
