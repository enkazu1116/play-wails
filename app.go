package main

import (
	"context"
	"fmt"
	"play-wails/infarstructure/db"
)

/*
 * アプリの構造体
 * コンテキストとDBのインスタンスを保持
 * 起動はturso.goが責務を持つ
 */
type App struct {
	ctx context.Context
	db  *db.TursoDB
}

/*
 * アプリのインスタンスを作成
 */
func NewApp(db *db.TursoDB) *App {
	return &App{
		db: db,
	}
}

/*
 * アプリの起動
 */
func (a *App) startup(ctx context.Context) {
	// コンテキストを保存
	a.ctx = ctx
}

func (a *App) Greet(name string) string {
	// Turso への接続＆クエリ実行テスト
	if err := a.db.HealthCheck(a.ctx); err != nil {
		return fmt.Sprintf("Hello %s, but Turso error: %v", name, err)
	}

	return fmt.Sprintf("Hello %s, Turso connection & query OK!", name)
}
