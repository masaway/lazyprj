# README ビジュアル素材 計画

## 目的

README にスクリーンショット・GIF を追加し、初見ユーザーが「何ができるか」「どう使うか」を
テキストを読まずに把握できるようにする。

---

## 素材一覧

| # | 種別 | ファイル名（案） | 挿入箇所 |
|---|------|----------------|---------|
| 1 | GIF | `demo-overview.gif` | README 冒頭（`# muxflow` 直下） |
| 2 | GIF | `demo-first-run.gif` | 「初回起動」セクション |
| 3 | GIF | `demo-editor.gif` | 「起動コマンドを登録する」セクション |
| 4 | PNG | `screenshot-main.png` | 「2回目以降の起動」セクション、または概要説明 |
| 5 | PNG | `screenshot-scan.png` | 「プロジェクトをスキャン」手順 |

---

## 各素材の詳細

### 1. `demo-overview.gif` — ヒーローデモ（最重要）

**目的**: README を開いた人が最初に見るデモ。「これが muxflow だ」を一発で伝える。

**収録する操作フロー**:

1. `muxflow` 起動 → メイン画面表示（複数プロジェクトが登録済みの状態）
2. `j` / `k` でプロジェクトを移動（右パネルの Detail が更新されることを見せる）
3. `Enter` でセッション起動 → tmux にアタッチ
4. （アタッチ後 `tmux detach` などで戻り、再度起動 → `Enter` で即アタッチになることを見せる）

**収録時の準備**:
- 5〜8 個のプロジェクトを事前に登録しておく
- 少なくとも 1 プロジェクトに `docker compose up` など複数ウィンドウ構成を登録しておく
- ターミナルウィンドウは 120x30 程度に設定

**推奨収録時間**: 20〜30 秒

---

### 2. `demo-first-run.gif` — 初回セットアップ

**目的**: 「インストール直後に何をすればいいか」の流れを伝える。

**収録する操作フロー**:

1. 設定ファイルを削除してまっさらな状態にする（収録前の準備）
2. `muxflow` 起動 → Setup 画面が自動で開く
3. スキャンディレクトリ（例: `~/work`）を入力 → `Enter` で保存
4. スキャン画面が自動で開く → `a` で全選択 → `Enter` で登録
5. メイン画面に戻る → プロジェクト一覧が表示された状態で終わる

**推奨収録時間**: 25〜40 秒

---

### 3. `demo-editor.gif` — 起動コマンド登録

**目的**: muxflow の一番の使いどころ（コマンド登録）を伝える。

**収録する操作フロー**:

1. メイン画面でプロジェクトにカーソルを合わせる
2. `e` でエディタ画面を開く（オーバーレイ表示）
3. `a` でウィンドウを追加 → ウィンドウ名を入力 → `Enter`
4. `Tab` で Panes パネルに移動 → `a` でペイン追加
5. ペイン編集フォームで Dir / Command / Execute を入力 → `Enter` で保存
6. `w` で全体保存 → `Esc` でメイン画面に戻る

**推奨収録時間**: 30〜45 秒

---

### 4. `screenshot-main.png` — メイン画面スクリーンショット

**目的**: メイン画面のレイアウト（左: プロジェクト一覧、右: 詳細）を静的に伝える。

**撮影条件**:
- 複数プロジェクトが登録されている
- 起動中セッション（●）と停止中セッション（○）が両方見える
- 自動起動フラグ（★）が付いているプロジェクトが含まれている
- Detail パネルに複数ウィンドウ・ペインの設定が表示されている

---

### 5. `screenshot-scan.png` — スキャン画面スクリーンショット

**目的**: スキャン画面の UI を見せる。テキスト説明だけでは伝わりにくい「スキップ済みセクション」も含める。

**撮影条件**:
- 複数のリポジトリがスキャン結果として表示されている
- いくつかが `[x]` で選択済みの状態
- 「スキップ済み」セクションも表示されている（展開状態）

---

## 収録ツール（推奨）

| ツール | 用途 | 備考 |
|--------|------|------|
| [vhs](https://github.com/charmbracelet/vhs) | GIF 自動生成 | `.tape` スクリプトで再現可能。Bubble Tea 製ツールとの相性が良い |
| [asciinema](https://asciinema.org/) | ターミナル録画 | `agg` で GIF に変換できる |
| [terminalizer](https://github.com/faressoft/terminalizer) | GIF 録画 | 設定ファイルでフレームレート等を調整可 |
| `scrot` / `gnome-screenshot` | PNG スクリーンショット | Linux 向け |
| macOS `Cmd+Shift+4` | PNG スクリーンショット | macOS 向け |

**推奨: vhs**

muxflow 自体が Charmbracelet 製ライブラリを使っており、vhs との親和性が高い。
`.tape` スクリプトをリポジトリに含めておくと、将来の UI 変更時に再生成しやすい。

---

## 素材の保存場所

```
muxflow/
└── docs/
    └── assets/
        ├── demo-overview.gif
        ├── demo-first-run.gif
        ├── demo-editor.gif
        ├── screenshot-main.png
        └── screenshot-scan.png
```

README からは相対パスで参照:
```markdown
![muxflow demo](docs/assets/demo-overview.gif)
```

---

## README への挿入イメージ

```markdown
# muxflow

[![License: MIT](...)](#)

![muxflow demo](docs/assets/demo-overview.gif)

tmuxセッションをTUIで管理するツールです。...

---

## 初回起動

![初回セットアップ](docs/assets/demo-first-run.gif)

### 1. スキャンディレクトリを設定する
...

---

## 起動コマンドを登録する

![コマンド登録デモ](docs/assets/demo-editor.gif)

...
```

---

## 優先順位

1. **`demo-overview.gif`** — インパクト最大。最初に作る
2. **`screenshot-main.png`** — 静的なので作りやすい。overview GIF と合わせて投入
3. **`demo-editor.gif`** — muxflow の差別化ポイントなので重要
4. **`demo-first-run.gif`** — セットアップの不安を取り除く
5. **`screenshot-scan.png`** — あると丁寧だが優先度は低め
