package repository

import (
	"database/sql"
	"play-wails/internal/model"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type timeRecordRepositoryImpl struct {
	db *sqlx.DB
}

// UUIDはTEXT、durationはナノ秒整数
type timeRecordRow struct {
	ID         string    `db:"id"`
	RunID      string    `db:"run_id"`
	TaskID     string    `db:"task_id"`
	DeleteFlag int       `db:"delete_flag"`
	StartTime  time.Time `db:"start_time"`
	EndTime    time.Time `db:"end_time"`
	DurationNs int64     `db:"duration_ns"`
}

/*
 * レコードをモデルに変換
 *
 * @param row レコード
 * @return モデル
 */
func rowToTimeRecord(row *timeRecordRow) *model.TimeRecord {
	t := &model.TimeRecord{
		StartTime:  row.StartTime,
		EndTime:    row.EndTime,
		Duration:   time.Duration(row.DurationNs),
		DeleteFlag: row.DeleteFlag != 0,
	}
	t.ID, _ = uuid.Parse(row.ID)
	t.RunID, _ = uuid.Parse(row.RunID)
	t.TaskID, _ = uuid.Parse(row.TaskID)
	return t
}

/*
 * 実装クラスのインスタンス生成
 *
 * @param db データベース
 * @return インスタンス
 */
func NewTimeRecordRepositoryImpl(db *sql.DB) TimeRecordRepository {
	return &timeRecordRepositoryImpl{db: sqlx.NewDb(db, "libsql")}
}

/*
 * レコード作成
 *
 * @param record レコード
 * @return エラー
 */
func (r *timeRecordRepositoryImpl) Create(record *model.TimeRecord) error {

	// 削除フラグを取得
	deleteFlag := 0
	if record.DeleteFlag {
		deleteFlag = 1
	}

	// インサートクエリ作成
	query := `INSERT INTO time_records (
		id
		, run_id
		, task_id
		, delete_flag
		, start_time
		, end_time
		, duration_ns
	) VALUES (
		:id
		, :run_id
		, :task_id
		, :delete_flag
		, :start_time
		, :end_time
		, :duration_ns
	)`

	// インサート処理実行
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":          record.ID.String(),
		"run_id":      record.RunID.String(),
		"task_id":     record.TaskID.String(),
		"delete_flag": deleteFlag,
		"start_time":  record.StartTime,
		"end_time":    record.EndTime,
		"duration_ns": record.Duration.Nanoseconds(),
	})

	return err
}

/*
 * レコードを取得
 *
 * @param id レコードID
 * @return レコード, エラー
 */
func (r *timeRecordRepositoryImpl) FindByID(id uuid.UUID) (*model.TimeRecord, error) {
	var row timeRecordRow
	err := r.db.Get(&row,
		`SELECT 
			id
			, run_id
			, task_id
			, delete_flag
			, start_time
			, end_time
			, duration_ns 
		FROM time_records 
		WHERE id = ?`,
		id.String(),
	)

	// エラーチェック
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return rowToTimeRecord(&row), nil
}

/*
 * レコードを更新（開始・終了時刻・作業時間）
 *
 * @param record レコード
 * @return エラー
 */
func (r *timeRecordRepositoryImpl) Update(record *model.TimeRecord) error {
	query :=
		`UPDATE time_records 
			SET start_time = :start_time
			, end_time = :end_time
			, duration_ns = :duration_ns 
		WHERE id = :id`

	// 更新処理実行
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":          record.ID.String(),
		"start_time":  record.StartTime,
		"end_time":    record.EndTime,
		"duration_ns": record.Duration.Nanoseconds(),
	})

	return err
}

/*
 * レコード一覧を取得
 *
 * @param excludeDeleted true のとき論理削除済みを除外
 * @return レコード一覧, エラー
 */
func (r *timeRecordRepositoryImpl) List(excludeDeleted bool) ([]*model.TimeRecord, error) {
	query :=
		`SELECT 
			id
			, run_id
			, task_id
			, delete_flag
			, start_time
			, end_time
			, duration_ns 
		FROM time_records`

	// 論理削除済みを除外する場合
	if excludeDeleted {
		query += ` WHERE delete_flag = 0`
	}
	query += ` ORDER BY start_time DESC`

	// レコード一覧を取得
	var rows []timeRecordRow
	err := r.db.Select(&rows, query)
	if err != nil {
		return nil, err
	}

	// レコード一覧をモデルに変換
	list := make([]*model.TimeRecord, 0, len(rows))
	for i := range rows {
		list = append(list, rowToTimeRecord(&rows[i]))
	}

	return list, nil
}

/*
 * レコードを論理削除
 *
 * @param id レコードID
 * @return エラー
 */
func (r *timeRecordRepositoryImpl) Delete(id uuid.UUID) error {
	_, err := r.db.Exec(
		`UPDATE time_records SET delete_flag = 1 WHERE id = ?`,
		id.String(),
	)
	return err
}
