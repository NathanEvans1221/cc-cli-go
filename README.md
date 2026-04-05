# CC-CLI-Go - Claude Code CLI Implementation / CC-CLI-Go - Claude Code CLI 實作

> A Go implementation of Claude Code CLI for learning purposes. / 以 Go 實作的 Claude Code CLI 學習專案。
>
> This project implements a minimal but functional CLI tool that interacts with Anthropic's Claude API. / 本專案實作一個精簡但功能完整的 CLI 工具，與 Anthropic Claude API互動。

**Project Location / 專案位置**:

```
/Users/user-name/github/cc-cli-go
```

---

## Features / 功能

### Core Features / 核心功能

- ✅ **CLI Entry & Argument Parsing / CLI 入口與參數解析** (using Cobra / 使用 Cobra)
- ✅ **Anthropic API Streaming Client / Anthropic API 串流客戶端** (SSE)
- ✅ **Query Loop / 查詢迴圈** (Channel-based streaming / Channel 串流)
- ✅ **Tool System / 工具系統** (Interface-based, extensible / Interface 架構，可擴展)
- ✅ **Bash Tool / Bash 工具** (Execute shell commands / 執行 shell 命令)
- ✅ **Read Tool / Read 工具** (Read file contents / 讀取檔案內容)
- ✅ **Edit Tool / Edit 工具** (Edit files by exact string replacement / 字串替換編輯檔案)
- ✅ **Basic TUI / 基本 TUI** (Bubble Tea framework / Bubble Tea 框架)
- ✅ **Message Rendering / 訊息渲染**
- ✅ **Streaming Event Handling / 串流事件處理**
- ✅ **Concurrent Tool Execution / 工具並行執行**

### Technical Highlights / 技術亮點

| Feature / 功能                              | Description / 描述                                                                           |
| ------------------------------------------- | -------------------------------------------------------------------------------------------- |
| **Idiomatic Go / Go 慣用寫法**              | Follows Go best practices and conventions /遵循 Go 最佳實踐與慣例                            |
| **Channel-based Streaming / Channel 串流**  | Uses Go channels instead of generators / 使用 Go channel 而非 generator                      |
| **Interface-based Design / Interface 設計** | Tool system is easily extensible / 工具系統易於擴展                                          |
| **Concurrent Execution /並行執行**          | Parallel tool execution with goroutines + WaitGroup / 使用 goroutine + WaitGroup並行執行工具 |
| **Single Binary /單一二進位檔**             | No runtime dependencies /無執行時依賴                                                        |

---

## Directory Structure / 目錄結構

```
cc-cli-go/
├── cmd/
│   └── cc-cli-go/
│       └── main.go              # Application entry point /應用程式入口
├── internal/
│   ├── api/                     # Anthropic API client / API 客戶端
│   │   ├── client.go            # HTTP client setup / HTTP 客戶端設定
│   │   ├── config.go            # API configuration / API 設定
│   │   ├── events.go            # Event parsing / 事件解析
│   │   ├── request.go           # Request builder / 請求建構
│   │   └── streaming.go         # SSE streaming / SSE 串流
│   ├── cli/                     # CLI framework / CLI 框架
│   │   ├── root.go              # Root command /根命令
│   │   └── run.go               # Interactive mode command /互動模式命令
│   ├── query/                   # Query engine / 查詢引擎
│   │   ├── engine.go            # Core query loop /核心查詢迴圈
│   │   └── types.go             # Query types / 查詢類型
│   ├── tools/                   # Tool system / 工具系統
│   │   ├── tool.go              # Tool interface / 工具介面
│   │   ├── registry.go          # Tool registry / 工具註冊
│   │   ├── bash/
│   │   │   └── bash.go          # Bash tool / Bash 工具
│   │   ├── edit/
│   │   │   └── edit.go          # Edit tool / Edit 工具
│   │   └── read/
│   │       └── read.go          # Read tool / Read 工具
│   ├── tui/                     # Terminal UI /終端 UI
│   │   ├── app.go               # Bubble Tea app / Bubble Tea應用
│   │   └── model.go             # TUI model / TUI 模型
│   └── types/                   # Shared types / 共用類型
│       ├── content.go           # Content block types /內容區塊類型
│       ├── message.go           # Message types /訊息類型
│       ├── usage.go             # Token usage / Token 使用量
│       └── uuid.go              # UUID generation / UUID 生成
├── bin/                         # Compiled binaries /編譯後的二進位檔
│   └── cc-cli-go
├── go.mod                       # Go module definition / Go模組定義
├── go.sum                       # Dependency checksums / 依賴校驗
├── TODO.md                      # Task tracking / 任務追蹤
└── .gitignore
```

---

## Documentation Index / 文件索引

### Root Documents / 根目錄文件

- [README.md](README.md) - Project overview, setup, usage, and development notes / 專案總覽、安裝、使用與開發說明
- [CHANGELOG.md](CHANGELOG.md) - Version history and notable changes / 版本歷史與重要變更
- [TODO.md](TODO.md) - Task tracking, completion status, and roadmap / 任務追蹤、完成狀態與路線圖
- [TESTING.md](TESTING.md) - Testing strategy, commands, and best practices / 測試策略、指令與最佳實踐
- [GIT_REMOTE_SETUP.md](GIT_REMOTE_SETUP.md) - Git remote setup guide / Git 遠端設定指南
- [.agent_task_state.md](.agent_task_state.md) - Agent task state snapshot / Agent 任務狀態快照

### docs/ Documents / docs 目錄文件

- [docs/IMPLEMENTATION_PLAN.md](docs/IMPLEMENTATION_PLAN.md) - Implementation phases and task breakdown / 實作階段與任務拆解
- [docs/GO_MODULE_STRUCTURE.md](docs/GO_MODULE_STRUCTURE.md) - Go module and package structure / Go 模組與套件結構
- [docs/ARCHITECTURE_ANALYSIS.md](docs/ARCHITECTURE_ANALYSIS.md) - Source architecture analysis (TS/Bun reference) / 原始架構分析（TS/Bun 參考）
- [docs/GO_ARCHITECTURE.md](docs/GO_ARCHITECTURE.md) - Go architecture design and data flow / Go 架構設計與資料流
- [docs/CORE_FEATURES.md](docs/CORE_FEATURES.md) - Core feature scope and priorities / 核心功能範圍與優先級
- [docs/TECH_DECISIONS.md](docs/TECH_DECISIONS.md) - ADR-style technical decisions / ADR 形式技術決策
- [docs/OPEN_QUESTIONS.md](docs/OPEN_QUESTIONS.md) - Discussion topics, status tracking, and unfinished items / 討論題目、狀態追蹤與未完成事項

### Documentation Maintenance Rules / 文件維護規則

1. 若新增、重新命名或刪除任何 `.md` 檔案，必須同步更新本節「Documentation Index / 文件索引」。
2. `README.md` 應保持為文件入口，所有重要說明文件需可從本檔案直接點擊進入。
3. PR 若包含文件異動，需一併檢查連結有效性與簡介是否仍正確。
4. 若內容狀態已變更（例如完成度、測試數字、版本），需同步更新對應文件（至少 `README.md`、`TODO.md`、`CHANGELOG.md`）。

---

## Prerequisites / 先決條件

- **Go 1.21 or higher / Go 1.21 或更高版本**
- **Anthropic API Key / Anthropic API 金鑰**

---

## Installation / 安裝

### Clone the Repository / 克隆儲存庫

```bash
cd /Users/user-name/github
git clone <repository-url> cc-cli-go
cd cc-cli-go
```

### Windows 11 Example / Windows 11 範例

```powershell
cd D:\github\chiisen
git clone <repository-url> cc-cli-go
cd cc-cli-go
```

### Install Dependencies / 安裝依賴

```bash
go mod download
```

### Build / 建構

```bash
# Build binary / 建構二進位檔
go build -o bin/cc-cli-go ./cmd/cc-cli-go

# Windows 11: build with .exe extension / Windows 11: 使用 .exe 副檔名建構
go build -o bin/cc-cli-go.exe ./cmd/cc-cli-go

# Or use go install / 或使用 go install
go install ./cmd/cc-cli-go
```

---

## Usage / 使用方式

### Set API Key / 設定 API 金鑰

```bash
export ANTHROPIC_API_KEY="your-api-key-here"
```

### Windows 11 PowerShell / Windows 11 PowerShell

```powershell
$env:ANTHROPIC_API_KEY = "your-api-key-here"
```

### Windows 11 CMD / Windows 11 CMD

```cmd
setx ANTHROPIC_API_KEY "your-api-key-here"
```

### Run Interactive Mode / 執行互動模式

```bash
# Using compiled binary / 使用編譯後的二進位檔
./bin/cc-cli-go run

# Or using go run / 或使用 go run
go run ./cmd/cc-cli-go run
```

### Windows 11 Run / Windows 11 執行

```powershell
# If you built with .exe / 如果已建構為 .exe
.\bin\cc-cli-go.exe run

# If you built without .exe / 如果沒有 .exe 副檔名
.\bin\cc-cli-go run
```

### Check Version / 檢查版本

```bash
./bin/cc-cli-go --version
# Output /輸出: cc-cli-go version 0.1.0
```

### Windows 11 Check Version / Windows 11 檢查版本

```powershell
# If you built with .exe / 如果已建構為 .exe
.\bin\cc-cli-go.exe --version

# If you built without .exe / 如果沒有 .exe 副檔名
.\bin\cc-cli-go --version
```

### Available Commands / 可用命令

```bash
./bin/cc-cli-go --help
```

---

## Testing / 測試

### Build Test / 建構測試

```bash
# Build all packages / 建構所有套件
go build ./...

# Build specific package / 建構特定套件
go build ./internal/api
go build ./internal/tools
go build ./internal/query
```

### Run Tests / 執行測試

```bash
# Run all tests (when available) / 執行所有測試（當可用時）
go test ./...

# Run tests with coverage / 執行測試並產生覆蓋率報告
go test -cover ./...

# Run specific package tests / 執行特定套件測試
go test ./internal/tools/bash
```

### Manual Testing / 手動測試

1. **Test Version Command / 測試版本命令**

   ```bash
   ./bin/cc-cli-go --version
   ```

2. **Test Interactive Mode / 測試互動模式**

   ```bash
   ANTHROPIC_API_KEY=your-key ./bin/cc-cli-go run
   ```

   - Type a message and press Enter /輸入訊息並按 Enter
   - Watch the streaming response / 觀看串流回應
   - Press Ctrl+C to exit / 按 Ctrl+C 退出

3. **Test Tools (when API is available) / 測試工具（當 API 可用時）**
   - The assistant can use Bash, Read, and Edit tools /助手可使用 Bash、Read、Edit 工具
   - Tools are executed concurrently when safe / 工具在安全時並行執行
   - Results are displayed in the TUI /結果顯示於 TUI

---

## Development /開發

### Project Statistics /專案統計

| Item / 項目               | Value /數值                                           |
| ------------------------- | ----------------------------------------------------- |
| **Total Files /總檔案數** | 22 Go source files / Go 源碼檔案                      |
| **Total Lines /總行數**   | ~1,500+ LOC                                           |
| **Packages /套件數**      | 8 internal packages /內部套件                         |
| **Dependencies /依賴**    | cobra (CLI), bubbletea (TUI), bubbles, lipgloss, uuid |

### Architecture /架構

```
User Input → TUI (Bubble Tea) → Query Engine → API Client → Claude API
使用者輸入 → TUI (Bubble Tea) → 查詢引擎 → API 客戶端 → Claude API
                 ↓                      ↓
             Messages            Tool Execution
             訊息列表            工具執行
                                        ↓
                               Tool Registry
                               工具註冊
                                        ↓
                             Bash / Read / Edit
                             Bash / Read / Edit 工具
```

### Key Design Decisions /關鍵設計決策

1. **Channel-based Streaming / Channel 串流**: Uses Go channels for event streaming instead of TypeScript generators / 使用 Go channel 進行事件串流而非 TypeScript generator

2. **Interface-based Tools / Interface 工具**: Tool interface allows easy addition of new tools / 工具介面允許輕鬆新增工具

3. **Concurrent Execution /並行執行**: Tools are executed in parallel when concurrency-safe / 工具在並行安全時並行執行

4. **Explicit Error Handling /明確錯誤處理**: Go's explicit error handling pattern / Go 的明確錯誤處理模式

---

## Future Enhancements /未來增強功能

Potential features for future development / 未來開發的潛在功能：

### Core Tools (P0 - Required) / 核心工具 (P0 - 必要) ✅

- [x] **Write Tool / Write 工具** - Create new files / 建立新檔案
- [x] **Glob Tool / Glob 工具** - File pattern matching / 檔案模式匹配
- [x] **Grep Tool / Grep 工具** - Search file contents / 搜尋檔案內容
- [x] **Permission System / 權限系統** - Tool approval dialogs / 工具核准對話框

### Context & Session (P0 - Required) / 環境與會話 (P0 - 必要) ✅

- [x] **Context Building / 環境建構** - Git status, branch, cwd / Git 狀態、分支、當前目錄
- [x] **CLAUDE.md Discovery / CLAUDE.md 發現** - Automatic context loading / 自動載入環境
- [x] **Session Storage / 會話儲存** - Conversation persistence / 對話持久化
- [x] **Resume Functionality / 恢復功能** - `-c`, `--resume` flags / `-c`, `--resume` 參數

### TUI Enhancement (P0 - Required) / TUI 增強 (P0 - 必要) ✅

- [x] **Permission Dialog / 權限對話框** - Allow/Deny buttons / 允許/拒絕按鈕
- [x] **Input Handling / 輸入處理** - Multi-line, history navigation / 多行輸入、歷史導航
- [x] **Keyboard Shortcuts / 鍵盤快捷鍵** - Ctrl+C, Ctrl+D, Escape

### Advanced Features (P1/P2 - Optional) / 進階功能 (P1/P2 - 可選) ✅

- [x] **Context Compaction / 環境壓縮** - Auto-compact when context exceeds threshold / 自動壓縮
- [x] **Configuration Management / 設定管理** - settings files / 設定檔管理
- [x] **Comprehensive Test Suite / 完整測試套件** - Unit and integration tests / 單元與整合測試
- [x] **Improved Error Handling / 增強錯誤處理** - Better error messages / 更好的錯誤訊息

> See `TODO.md` for detailed task tracking. / 詳細任務追蹤請見 `TODO.md`。

---

## Troubleshooting /故障排除

### Build Errors / 建構錯誤

```bash
# Clean build cache / 清理建構快取
go clean -cache

# Update dependencies / 更新依賴
go mod tidy
go mod download

# Rebuild / 重新建構
go build -o bin/cc-cli-go ./cmd/cc-cli-go
```

### API Key Issues / API 金鑰問題

```bash
# Verify API key is set / 驗證 API 金鑰已設定
echo $ANTHROPIC_API_KEY

# Set API key / 設定 API 金鑰
export ANTHROPIC_API_KEY="sk-ant-api03-..."
```

### Runtime Errors / 執行時錯誤

| Error / 錯誤                                             | Solution / 解決方案                                            |
| -------------------------------------------------------- | -------------------------------------------------------------- |
| **"ANTHROPIC_API_KEY environment variable is required"** | Set your API key / 設定您的 API 金鑰                           |
| **Connection errors /連接錯誤**                          | Check network connectivity / 檢查網路連接                      |
| **Rate limiting /速率限制**                              | Wait and retry with exponential backoff / 等待並以指数退避重試 |

---

## Contributing / 貢獻指南

1. Fork the repository / Fork 儲存庫
2. Create a feature branch / 建立功能分支
3. Make your changes / 進行變更
4. Run tests / 執行測試
5. Submit a pull request / 提交 Pull Request

---

## License / 授權條款

This project is for learning purposes. See LICENSE file for details. / 本專案僅供學習用途。詳見 LICENSE 檔案。

---

## Acknowledgments / 致謝

- Original Claude Code CLI by Anthropic / Anthropic 原始 Claude Code CLI
- Bubble Tea framework by Charm / Charm Bubble Tea 框架
- Cobra CLI framework by Spf13 / Spf13 Cobra CLI 框架

---

## Version History / 版本歷史

| Version / 版本 | Date / 日期 | Features / 功能                                                 |
| -------------- | ----------- | --------------------------------------------------------------- |
| **v0.1.0**     | 2026-04-01  | Initial release with core features / 初版發布，含核心功能       |
|                |             | • Basic CLI structure / 基本 CLI 結構                           |
|                |             | • API streaming client / API 串流客戶端                         |
|                |             | • Tool system (Bash, Read, Edit) / 工具系統（Bash、Read、Edit） |
|                |             | • Query engine / 查詢引擎                                       |
|                |             | • Basic TUI / 基本 TUI                                          |

---

> **Note / 注意**: This is a learning project and not intended for production use. / 本專案僅供學習用途，不建議用於生產環境。For production use, please refer to the official Claude Code CLI. / 生產用途請參考官方 Claude Code CLI。
