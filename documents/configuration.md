# 設定ファイル

設定ファイルは初回起動時に自動作成されます。

| パス | 説明 |
|------|------|
| `~/.config/muxflow/config.json` | 設定ファイル |

`XDG_CONFIG_HOME` が設定されている場合は `$XDG_CONFIG_HOME/muxflow/config.json` を使用します。

---

## 設定の構造

```json
{
  "projects": [
    {
      "name": "myapp",
      "path": "~/work/myapp",
      "auto_start": true,
      "description": "メインの開発プロジェクト",
      "windows": [
        {
          "name": "dev",
          "layout": "even-horizontal",
          "panes": [
            { "dir": ".", "command": "", "execute": false },
            { "dir": ".", "command": "npm run dev", "execute": true }
          ]
        }
      ]
    }
  ],
  "skipped_paths": [
    "/home/user/work/old-project"
  ],
  "settings": {
    "scan_directory": "/home/user/work"
  }
}
```

### プロジェクトフィールド

| フィールド | 型 | 説明 |
|-----------|-----|------|
| `name` | string | tmuxセッション名 |
| `path` | string | プロジェクトのルートディレクトリ |
| `auto_start` | bool | 起動時に自動起動するか（省略時: false） |
| `description` | string | 説明文（省略可） |
| `windows` | array | ウィンドウ・ペイン構成 |

### スキップ済みパス

`skipped_paths` にはスキャン画面で `x` キーによりスキップしたプロジェクトのフルパスが保存されます。  
一覧に表示したい場合は該当パスを削除するか、スキャン画面のスキップ済みセクションから `x` で解除できます。

---

## tmuxレイアウト

エディタ画面でウィンドウのレイアウトを選択できます（プレビュー付き）。

| レイアウト | 説明 |
|-----------|------|
| `even-horizontal` | ペインを左右に均等分割 |
| `even-vertical` | ペインを上下に均等分割 |
| `main-horizontal` | 上に広いメイン、下に小ペイン |
| `main-vertical` | 左に広いメイン、右に小ペイン |
| `tiled` | グリッド状に配置 |
