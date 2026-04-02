# CHANGELOG

本專案所有重要變更都將記錄在此文件中。

格式基於 [Keep a Changelog](https://keepachangelog.com/zh-TW/1.0.0/)，
並且本專案遵守 [Semantic Versioning](https://semver.org/lang/zh-TW/) 版本號規則。

## [Unreleased]

### 新增功能 / Added

#### 工具系統 / Tools

- **Write Tool**: 建立新檔案工具，僅允許建立新檔案（若檔案已存在則拒絕），支援自動建立父目錄。
- **Glob Tool**: 檔案模式匹配工具，支援 glob 模式搜尋（如 `**/*.go`），返回匹配的檔案列表。
- **Grep Tool**: 檔案內容搜尋工具，支援正則表達式搜尋，支援檔案類型過濾（如 `*.go`）。

## [0.1.0] - 2026-04-01

### 新增功能 / Added

#### 核心功能 / Core Features

- CLI 入口與參數解析（使用 Cobra）
- Anthropic API 串流客戶端（SSE）
- 查詢迴圈（Channel-based streaming）
- 工具系統（Interface-based, extensible）

#### 工具 / Tools

- Bash Tool: 執行 shell 命令
- Read Tool: 讀取檔案內容（支援行號顯示與分頁）
- Edit Tool: 字串替換編輯檔案

#### TUI / Terminal User Interface

- 基本 TUI（Bubble Tea framework）
- 訊息渲染
- 串流事件處理
- 工具並行執行
- 鍵盤快捷鍵（Ctrl+C, Ctrl+D, Escape）