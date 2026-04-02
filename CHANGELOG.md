# CHANGELOG

本專案所有重要變更都將記錄在此文件中。

格式基於 [Keep a Changelog](https://keepachangelog.com/zh-TW/1.0.0/)，
並且本專案遵守 [Semantic Versioning](https://semver.org/lang/zh-TW/) 版本號規則。

## [Unreleased]

### 新增功能 / Added

#### 錯誤處理 / Error Handling

- **Unified Error Types**: 統一錯誤類型系統，支援 API/Tool/Permission/Config/Session/Internal 錯誤。
- **Error Wrapping**: 錯誤包裝器，支援錯誤鏈與上下文資訊。
- **User-friendly Messages**: 使用者友善錯誤訊息，包含建議解決方案。
- **API Errors**: API 錯誤處理（連線/認證/速率限制/逾時）。
- **Tool Errors**: 工具錯誤處理（權限/輸入驗證/執行/檔案操作）。

#### 環境壓縮 / Context Compaction

- **Context Compaction System**: 自動環境壓縮系統，當對話超過閾值時自動摘要舊訊息。
- **Summary Generation**: 智能摘要生成，提取主題、工具使用統計。
- **Manual Compaction**: 支援手動壓縮功能。
- **Token Estimation**: Token 估算與閾值檢測。

#### 程式碼品質 / Code Quality

- **Package Documentation**: 核心模組完整 Package 註解。
- **Type Documentation**: 所有公開類型詳細註解。
- **Function Documentation**: 重要函數功能說明。
- **Code Readability**: 程式碼可讀性改善。

#### 設定管理 / Configuration Management

- **Configuration System**: 完整設定管理系統，支援全域與專案設定檔。
- **Settings Structure**: 支援權限模式、權限規則、工具設定、API 設定。
- **Global Settings**: `~/.claude/settings.json` 全域設定檔。
- **Project Settings**: `.claude/settings.json` 專案設定檔，覆蓋全域設定。
- **Settings Validation**: 完整設定驗證功能，確保設定格式正確。
- **Settings Merging**: 專案設定優先於全域設定，自動合併機制。
- **Config Tests**: 20 個測試案例，覆蓋率 74.5%。

#### 測試框架 / Testing Framework

- **TESTING.md**: 完整測試策略文件，包含測試原則、工具、範例與最佳實踐。
- **testutil Package**: 測試輔助工具 package，提供常用斷言與檔案系統測試工具。
- **工具測試**:
  - Write Tool Tests: 完整單元測試（覆蓋率 81.8%，9 個測試案例）
  - Glob Tool Tests: 完整單元測試（覆蓋率 84.1%，10 個測試案例）
  - Grep Tool Tests: 完整單元測試（覆蓋率 80.0%，10 個測試案例）
  - Read Tool Tests: 完整單元測試（覆蓋率 80.0%，10 個測試案例）
  - Edit Tool Tests: 完整單元測試（覆蓋率 89.3%，11 個測試案例）
  - Bash Tool Tests: 完整單元測試（覆蓋率 93.9%，12 個測試案例）
- **系統測試**:
  - Permission System Tests: 完整單元測試（覆蓋率 90.9%，15 個測試案例）
  - Session Storage Tests: 完整單元測試（覆蓋率 88.0%，10 個測試案例）
  - Context Building Tests: 完整單元測試（覆蓋率 95.7%，14 個測試案例）
  - API Client Tests: 完整單元測試（覆蓋率 44.2%，11 個測試案例）
- **整合測試**:
  - Integration Tests: 6 個整合測試案例（測試工具執行與權限系統整合）
  - Mock Client: API Client Mock 用於測試
- **測試統計**: 總計 127 個測試案例，全部通過 ✅，平均覆蓋率 81.9%

#### 工具系統 / Tools

- **Write Tool**: 建立新檔案工具，僅允許建立新檔案（若檔案已存在則拒絕），支援自動建立父目錄。
- **Glob Tool**: 檔案模式匹配工具，支援 glob 模式搜尋（如 `**/*.go`），返回匹配的檔案列表。
- **Grep Tool**: 檔案內容搜尋工具，支援正則表達式搜尋，支援檔案類型過濾（如 `*.go`）。

#### 環境建構 / Context Building

- **Context Building**: 自動收集環境資訊（工作目錄、Git 狀態、Git 分支、日期時間），並加入 system prompt。
- **CLAUDE.md Discovery**: 自動從當前目錄向上搜尋 CLAUDE.md 和 GEMINI.md 檔案，並合併至 system prompt。

#### TUI 增強 / TUI Enhancement

- **Input Handling**: 多行輸入支援，歷史導航（up/down 鍵），貼上處理，自動調整高度。
- **Textarea Integration**: 使用 textarea.Model 取代 textinput.Model，支援多行編輯。
- **History Management**: 命令歷史管理（最多 1000 條），支援歷史導航。
- **Smart Navigation**: 智能導航（單行模式時啟用歷史導航，多行模式時保留游標移動）。

- **Session Storage**: 會話持久化功能，支援 JSONL 格式儲存對話記錄與元資料。
- **Resume Functionality**: 恢復會話功能，支援 `-c` / `--continue` 繼續上次會話，`--resume <session-id>` 恢復指定會話。
- **Session Cleanup**: 自動清理舊會話功能。

#### 權限系統 / Permission System

- **Permission System**: 工具核准系統，支援多種權限模式（default、accept、plan、auto）與權限規則（allow、deny、ask）。
- **Permission Dialog**: 互動式權限對話框，支援 Allow / Deny / Always Allow / Always Deny 選項。
- **Dangerous Command Detection**: 危險命令檢測功能，自動識別並警告危險操作（如 `rm -rf`、`DROP TABLE`、`git push --force`）。

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