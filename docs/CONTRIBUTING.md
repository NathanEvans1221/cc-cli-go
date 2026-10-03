# 參與開發

## 開發環境

需要 Go 1.25。在專案根目錄執行 `go test ./internal/...`。本倉庫沒有要求 Makefile 或 GitHub Actions。

## 風格

變更保持在相關套件內。測試要呼叫實際函式；網路與模型傳輸用注入的假伺服器或 `NewEngineWithStreamer`。提交訊息使用 `<type>(<scope>): <繁體中文主旨>`。

## 送出變更

以分支開發，確認 `go test` 與 `git diff --check`。說明行為變更，並在 `CHANGELOG.md` 的 Unreleased 記下使用者看得到的修正或文件。
