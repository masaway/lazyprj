# Contributing to lazyprj

バグ報告・機能提案・プルリクエスト、どれも歓迎です。

---

## バグ報告・機能提案

[Issues](https://github.com/masaway/lazyprj/issues) からテンプレートを選んで投稿してください。

---

## ブランチ戦略

[GitHub Flow](https://docs.github.com/en/get-started/using-github/github-flow) を採用しています。

1. `main` から作業ブランチを切る
2. ブランチ上で実装・コミット
3. PRを作成して `main` へマージ
4. マージ後はブランチを削除

```bash
git switch -c feature/your-feature-name
# ... 実装 ...
git push origin feature/your-feature-name
# → GitHub上でPRを作成
```

ブランチ名は `feature/`, `fix/`, `docs/` などのプレフィックスをつけると分かりやすいです。

`main` への直pushはブランチ保護により禁止されています。

---

## プルリクエスト

### 開発環境のセットアップ

```bash
git clone https://github.com/masaway/lazyprj
cd lazyprj
go mod tidy
```

### ビルドと動作確認

```bash
go build -o lazyprj .
./lazyprj
```

### PRを送る前に

- `go build` が通ること
- 既存の動作を壊していないこと（手動で一通り操作して確認）
- コミットメッセージは変更内容が伝わるように書いてください

### アーキテクチャについて

[CLAUDE.md](CLAUDE.md) にコード構成の概要があります。実装の参考にしてください。
