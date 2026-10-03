# 工具開發指南

新工具放在 `internal/tools/<name>`，實作 `tools.Tool`。

## 介面

必須提供 `Name`、`Description`、`InputSchema`、`Execute`、`IsEnabled`、`IsReadOnly`、`IsConcurrencySafe` 與 `UserFacingName`。`Execute` 回傳 `*tools.ToolResult`。失敗時設 `IsError`，不要把失敗內容當成成功。

## 實作步驟

1. 在 `New()` 回傳具體型別。需要網路的工具把 `*http.Client` 做成欄位，測試注入 `httptest`。
2. 在 `internal/cli/run.go` 的 `prepareInteractive` 呼叫 `toolReg.Register`。
3. 用 `context` 的取消或 `QueryParams.ToolTimeout` 結束長時間工作，並回傳錯誤結果。

## 慣例

`WebFetch` 與 `WebSearch` 不在測試裡打外部網路。`WebSearch` 沒有設定 `BaseURL` 時回傳錯誤。`Agent` 只呼叫注入的 `Runner`。
