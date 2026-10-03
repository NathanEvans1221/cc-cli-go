# API 使用指南

這個套件是 `internal/api`。客戶端會重用 `http.Transport` 連線，並對連線失敗、HTTP 429 與 503 最多重試 3 次。

## 客戶端

`NewClient(apiKey)` 會設定 Anthropic 的 base URL 與標頭 `x-api-key`、`anthropic-version`。`Stream` 對 `POST /messages` 讀取 SSE，略過沒有 `data: ` 前綴的行。

## 請求與回應

`NewRequest(model, maxTokens)` 預設 `stream: true`。`SetSystem`、`AddMessage` 與 `AddTool` 組出訊息。串流事件用 `ParseMessageStart`、`ParseContentBlock` 與 `ParseDelta` 轉成 `internal/types`。

## 錯誤處理

非 200 狀態會從 `Stream` 回傳錯誤。可重試的 429 與 503 會在 `doRequest` 裡用完次數後回傳最後一次錯誤。不要在測試裡呼叫正式的 Anthropic 端點；把 `baseURL` 指到注入的 HTTP 伺服器。
