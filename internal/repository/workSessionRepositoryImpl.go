package repository

import (
	"database/sql"
	"play-wails/internal/model"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type workSessionRepositoryImpl struct {
	db *sqlx.DB
}

// UUIDはTEXTのためstringで受ける
type workSessionRow struct {
	ID        string     `db:"id"`
	RunID     string     `db:"run_id"`
	TaskID    string     `db:"task_id"`
	StartTime time.Time  `db:"start_time"`
	EndTime   *time.Time `db:"end_time"`
}

/*
 * レコードをモデルに変換
 *
 * @param row レコード
 * @return モデル
 */
func rowToWorkSession(row *workSessionRow) *model.WorkSession {

	s := &model.WorkSession{
		StartTime: row.StartTime,
		EndTime:   row.EndTime,
	}
	s.ID, _ = uuid.Parse(row.ID)
	s.RunID, _ = uuid.Parse(row.RunID)
	s.TaskID, _ = uuid.Parse(row.TaskID)

	return s
}

/*
 * 実装クラスのインスタンス生成
 *
 * @param db データベース
 * @return インスタンス
 */
func NewWorkSessionRepositoryImpl(db *sql.DB) WorkSessionRepository {
	return &workSessionRepositoryImpl{db: sqlx.NewDb(db, "libsql")}
}

/*
 * レコード作成
 *
 * @param session レコード
 * @return エラー
 */
func (r *workSessionRepositoryImpl) Create(session *model.WorkSession) error {
	query := `INSERT INTO work_sessions (
		id
		, run_id
		, task_id
		, start_time
		, end_time
	) VALUES (
		:id
		, :run_id
		, :task_id
		, :start_time
		, :end_time
	)`

	// インサート処理実行
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":         session.ID.String(),
		"run_id":     session.RunID.String(),
		"task_id":    session.TaskID.String(),
		"start_time": session.StartTime,
		"end_time":   session.EndTime,
	})

	return err
}

/*
 * レコードを取得
 *
 * @param id レコードID
 * @return レコード, エラー
 */
func (r *workSessionRepositoryImpl) FindByID(id uuid.UUID) (*model.WorkSession, error) {
	var row workSessionRow
	err := r.db.Get(&row,
		`SELECT 
			id
			, run_id
			, task_id
			, start_time
			, end_time 
		FROM work_sessions 
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

	// レコードをモデルに変換
	return rowToWorkSession(&row), nil
}

/*
 * レコードを更新
 *
 * @param session レコード
 * @return エラー
 */
func (r *workSessionRepositoryImpl) Update(session *model.WorkSession) error {
	query := `UPDATE work_sessions SET 
		run_id = :run_id
		, task_id = :task_id
		, start_time = :start_time
		, end_time = :end_time
	WHERE id = :id`

	// 更新処理実行
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":         session.ID.String(),
		"run_id":     session.RunID.String(),
		"task_id":    session.TaskID.String(),
		"start_time": session.StartTime,
		"end_time":   session.EndTime,
	})
	return err
}

/*
 * レコードを取得
 *
 * @param runID レコードID
 * @return レコード, エラー
 */
func (r *workSessionRepositoryImpl) ListByRunID(runID uuid.UUID) ([]*model.WorkSession, error) {

	var rows []workSessionRow
	err := r.db.Select(&rows,
		`SELECT 
			id
			, run_id
			, task_id
			, start_time
			, end_time 
		FROM work_sessions 
		WHERE run_id = ? 
		ORDER BY start_time`,
		runID.String(),
	)

	// エラーチェック
	if err != nil {
		return nil, err
	}

	// ワークセッションを全てリストに追加
	list := make([]*model.WorkSession, 0, len(rows))
	for i := range rows {
		list = append(list, rowToWorkSession(&rows[i]))
	}

	return list, nil
}

/*
 * レコードを削除
 *
 * @param id レコードID
 * @return エラー
 */
func (r *workSessionRepositoryImpl) Delete(id uuid.UUID) error {
	_, err := r.db.Exec(
		`DELETE FROM work_sessions WHERE id = ?`,
		id.String(),
	)

	return err
}
