# 使用者指南

## 安裝

在專案根目錄執行 `go build -o bin/cc-cli-go.exe ./cmd/cc-cli-go`。執行前設定環境變數 `ANTHROPIC_API_KEY`。

## 設定

全域檔案是使用者目錄下的 `.claude/settings.json`。可設定 `permission.mode`（`default`、`accept`、`plan`、`auto`）、`api.model`、`ui.language` 與 `ui.theme`。`ui.language` 設為 `zh-TW` 時，使用者標籤是「你:」，進行中的查詢顯示「思考中...」。

## 一般執行

`cc-cli-go --version` 印出版本。`cc-cli-go run` 進入互動模式。`cc-cli-go run --continue` 接上一輪會話，`cc-cli-go run --resume <id>` 指定會話。找不到會話時會警告並開新會話。

## 疑難排解

沒有 API key 時程式直接結束並說明要設定 `ANTHROPIC_API_KEY`。權限被拒的工具不會執行。查詢錯誤若帶有建議，畫面會顯示該建議。
