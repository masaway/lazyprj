# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## ビルドと実行

```bash
# ビルド
go build -o lazyprj .

# 実行（リポジトリルートから）
./lazyprj

# 依存関係の更新
go mod tidy
```

## アーキテクチャ概要

tmuxセッションをTUIで管理するGoアプリケーション。[Bubble Tea](https://github.com/charmbracelet/bubbletea) の Model-Update-View パターンで実装されている。

### 画面遷移

`internal/ui/app.go` の `App` struct が Bubble Tea のルートモデル。`currentScreen` フィールドで画面を切り替える：

- `screenMain` → `App`（メイン画面、プロジェクト一覧＋詳細）
- `screenEditor` → `EditorModel`（ウィンドウ・ペイン編集）
- `screenScan` → `ScanModel`（未登録プロジェクトのスキャン）

### アタッチの遅延実行

tmuxのアタッチは `tea.Quit` 後に実行する必要があるため、`App.pendingAttach` にセッション名を保存し、`main.go` で `p.Run()` の戻り値から `PendingAttach()` を呼び出して実行する。

### 設定ファイルの優先順位

`internal/config/config.go` の `GetConfigPath()` が `config/personal.json` → `config/default.json` の順で参照する。保存は常に `personal.json` に行われる。

### リポジトリルートの検出

`main.go` の `findRepoRoot()` が `config/` ディレクトリと `start-projects.sh` の両方が存在するディレクトリを上向きに探索してリポジトリルートとする。

### モジュール構成

| パッケージ | 役割 |
|-----------|------|
| `internal/config` | 設定のロード・保存・データ構造・プロジェクトスキャン |
| `internal/tmux` | tmuxコマンド操作（セッション作成・停止・アタッチ） |
| `internal/ui` | Bubble Tea の全画面モデル、スタイル定義 |
