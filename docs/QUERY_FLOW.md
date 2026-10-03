# 查詢與工具流程

`cmd/cc-cli-go` 進入 `internal/cli`。`prepareInteractive` 讀取設定、建立 `api.Client`、註冊工具、建立 `query.Engine`，然後才交給 TUI。測試可以替換 `runTUI`，避免啟動會阻塞的程式。

## 一輪查詢

1. `Engine.Query` 把目前訊息送到 `Client.Stream`。
2. 沒有 `tool_use` 時，這一輪以 `completed` 結束。
3. 有 `tool_use` 時，`executeTools` 依權限決定是否呼叫 `Tool.Execute`。`allow` 會執行；`deny` 與尚未核准的 `ask` 只回錯誤結果。
4. 工具結果組成下一則 user 訊息，進入下一輪，直到沒有工具或達到 `MaxTurns`。
5. `ToolTimeout` 大於 0 時，逾時的工具回傳錯誤，不會回報成功。

## 擴充點

新工具實作 `tools.Tool` 並在 `prepareInteractive` 註冊。模型傳輸實作 `Stream`，測試用 `NewEngineWithStreamer`。這份說明不要求 Makefile 或 GitHub Actions。
