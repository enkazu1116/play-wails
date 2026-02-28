package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

/*
 * TursoDBの構造体
 */
type TursoDB struct {
	db *sql.DB
}

/*
 * TursoDBのインスタンスを作成
 *
 * @return *sql.DB, error
 */
func NewTursoDB() (*TursoDB, error) {

	// .envファイルから環境変数を取り込む
	url, token, err := LoadEnv()
	if err != nil {
		log.Fatal("【ERROR】.envファイルの読み込みに失敗しました")
		return nil, err
	}

	// 接続URLの作成
	dsn := CreateTursoURL(url, token)

	// DB接続
	db, err := sql.Open("libsql", dsn)
	if err != nil {
		log.Fatal("【ERROR】DB接続に失敗しました")
		return nil, err
	}

	// 最大接続数・最大空き接続数を設定
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	return &TursoDB{db: db}, nil
}

/*
 * TursoDBのインスタンスをクローズ
 */
func (t *TursoDB) Close() error {
	return t.db.Close()
}

/*
 * 接続確認と簡単なクエリ実行
 * - PingContext で接続テスト
 * - SELECT 1 でクエリ実行テスト
 */
func (t *TursoDB) HealthCheck(ctx context.Context) error {
	// 接続確認
	if err := t.db.PingContext(ctx); err != nil {
		return err
	}

	// クエリ実行確認
	var n int
	if err := t.db.QueryRowContext(ctx, "SELECT 1").Scan(&n); err != nil {
		return err
	}

	return nil
}

/*
 * .envファイルからDB接続に必要な環境変数を取り込む
 */
func LoadEnv() (string, string, error) {
	// .envファイルから環境変数を読み込む
	err := godotenv.Load()
	if err != nil {
		log.Fatal("【ERROR】.envファイルの読み込みに失敗しました")
	}

	// DBの接続情報を定義
	// dbName := "todo-app"
	url := os.Getenv("TURSO_DATABASE_URL")
	token := os.Getenv("TURSO_AUTH_TOKEN")

	return url, token, nil
}

/*
 * 接続URLの作成
 *
 * @param url データベースのURL
 * @param token データベースのトークン
 *
 * @return 接続URL
 */
func CreateTursoURL(url string, token string) string {

	// Strings.Builderを使用して、URLを組み立てる
	var sb strings.Builder

	sb.WriteString("libsql://")
	sb.WriteString(url)
	sb.WriteString(".turso.io?authToken=")
	sb.WriteString(token)

	return sb.String()
}
