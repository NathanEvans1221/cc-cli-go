# TODO - cc-cli-go 開發任務清單 / Development Task List

> **狀態更新日期 / Status Update Date**: 2026-04-01
> **預估工時 / Estimated Hours**: ~25-35 小時 / hours

---

## ✅ 已完成功能 / Completed Features

**Phase 1: 核心功能 / Core Features**

- [x] CLI Entry & Argument Parsing / CLI 入口與參數解析 (using Cobra)
- [x] Anthropic API Streaming Client / Anthropic API 串流客戶端 (SSE)
- [x] Query Loop / 查詢迴圈 (Channel-based streaming / Channel 串流)
- [x] Tool System / 工具系統 (Interface-based, extensible / Interface 架構，可擴展)
- [x] Bash Tool / Bash 工具 (Execute shell commands / 執行 shell 命令)
- [x] Read Tool / Read 工具 (Read file contents / 讀取檔案內容)
- [x] Edit Tool / Edit 工具 (Edit files by exact string replacement / 字串替換編輯檔案)

**Phase 2: TUI (基本) / Basic TUI**

- [x] Basic TUI / 基本 TUI (Bubble Tea framework / Bubble Tea 框架)
- [x] Message Rendering / 訊息渲染
- [x] Streaming Event Handling / 串流事件處理
- [x] Concurrent Tool Execution / 工具並行執行
- [x] Keyboard Shortcuts / 鍵盤快捷鍵 (Ctrl+C, Ctrl+D, Escape)

---

## ❌ 待實作功能 / Pending Features

### 🔴 Phase 1: 核心工具 / Core Tools (P0 - 必要 / Required)

#### Tools 工具

- [ ] **Write Tool / Write 工具**
  - 功能: Create new files / 建立新檔案
  - 預估工時 / Est. Hours: 1h
  - 位置: `internal/tools/write/write.go`

- [ ] **Glob Tool / Glob 工具**
  - 功能: File pattern matching / 檔案模式匹配
  - 預估工時 / Est. Hours: 1.5h
  - 位置: `internal/tools/glob/glob.go`
  - 技術: `path/filepath.Glob`

- [ ] **Grep Tool / Grep 工具**
  - 功能: Search file contents / 搜尋檔案內容
  - 預估工時 / Est. Hours: 2h
  - 位置: `internal/tools/grep/grep.go`
  - 技術: `regexp`, `bufio.Scanner`

#### Permission System 權限系統

- [ ] **Permission System / 權限系統**
  - 功能: Tool approval dialogs / 工具核准對話框
  - 預估工時 / Est. Hours: 4-5h
  - 位置: `internal/permission/`
  - 子任務 / Subtasks:
    - [ ] Permission modes (default, accept, plan, auto) / 權限模式
    - [ ] Permission rules (allow, deny, ask) / 權限規則
    - [ ] Rule matching (tool name + input pattern) / 規則匹配
    - [ ] Dangerous command detection /危險命令檢測

#### Context Building 環境建構

- [ ] **Context Building / 環境建構**
  - 功能: Build context for API requests / 環境資訊建構
  - 預估工時 / Est. Hours: 2-3h
  - 位置: `internal/context/`
  - 子任務 / Subtasks:
    - [ ] Git status (`git status --porcelain`) / Git 狀態
    - [ ] Git branch (`git branch --show-current`) / Git 分支
    - [ ] Current working directory / 當前工作目錄
    - [ ] Date/time / 日期時間

- [ ] **CLAUDE.md Discovery / CLAUDE.md 發現**
  - 功能: Automatic context loading / 自動載入環境
  - 預估工時 / Est. Hours: 1h
  - 位置: `internal/context/claudemd.go`
  - 子任務 / Subtasks:
    - [ ] Search from current directory upward / 從當前目錄向上搜尋
    - [ ] Merge all found CLAUDE.md files / 合併所有找到的檔案

#### Session Storage 會話儲存

- [ ] **Session Storage / 會話儲存**
  - 功能: Conversation persistence / 對話持久化
  - 預估工時 / Est. Hours: 3-4h
  - 位置: `internal/session/`
  - 子任務 / Subtasks:
    - [ ] Session ID generation / Session ID 生成
    - [ ] Transcript storage (JSONL format) / Transcript 儲存 (JSONL 格式)
    - [ ] Session metadata / Session 元資料
    - [ ] Session cleanup (old sessions) / Session 清理

- [ ] **Resume Functionality / 恢復功能**
  - 功能: Continue previous sessions / 繼續之前的會話
  - 預估工時 / Est. Hours: 1h
  - CLI flags:
    - `-c` / `--continue`: Continue last session / 繼續上次會話
    - `--resume <session-id>`: Resume specific session / 恢復指定會話

---

### 🟡 Phase 2: TUI 增強 / TUI Enhancements (P0 - 必要 / Required)

#### Permission Dialog 權限對話框

- [ ] **Permission Dialog / 權限對話框**
  - 功能: Interactive permission UI /互動式權限 UI
  - 預估工時 / Est. Hours: 3-4h
  - 位置: `internal/tui/permission.go`
  - 子任務 / Subtasks:
    - [ ] Tool use request display / 工具使用請求顯示
    - [ ] Allow / Deny / Always Allow / Always Deny buttons / 按鈕
    - [ ] Input preview (file paths, commands) /輸入預覽
    - [ ] Keyboard navigation / 鍵盤導航

#### Input Handling輸入處理

- [ ] **Input Handling /輸入處理**
  - 功能: Enhanced input controls / 增強輸入控制
  - 預估工時 / Est. Hours: 3-4h
  - 位置: `internal/tui/input.go`
  - 子任務 / Subtasks:
    - [ ] Multi-line input / 多行輸入
    - [ ] History navigation (up/down keys) / 歷史導航
    - [ ] Paste handling / 貼上處理

- [x] **Keyboard Shortcuts / 鍵盤快捷鍵**
  - 功能: Essential keyboard controls / 基本鍵盤控制
  - 預估工時 / Est. Hours: 1h
  - 子任務 / Subtasks:
    - [x] Enter: Submit / 提交
    - [x] Ctrl+C: Interrupt / 中斷
    - [x] Ctrl+D: Exit / 退出
    - [x] Escape: Cancel / 取消

---

### 🟢 Phase 3: 進階功能 / Advanced Features (P1/P2 - 可選 / Optional)

#### Context Compaction 環境壓縮

- [ ] **Context Compaction / 環境壓縮**
  - 功能: Auto-compact when context exceeds threshold / 自動壓縮
  - 預估工時 / Est. Hours: 4-6h
  - 位置: `internal/services/compact/`
  - 子任務 / Subtasks:
    - [ ] Auto-compact trigger / 自動壓縮觸發
    - [ ] Manual `/compact` command / 手動命令
    - [ ] Summary generation / 摘要生成

#### Configuration Management 設定管理

- [ ] **Configuration Management / 設定管理**
  - 功能: Settings files management / 設定檔管理
  - 預估工時 / Est. Hours: 2-3h
  - 位置: `internal/config/`
  - 子任務 / Subtasks:
    - [ ] Global settings (`~/.claude/settings.json`) / 全域設定
    - [ ] Project settings (`.claude/settings.json`) /專案設定
    - [ ] Settings schema validation / 設定驗證

#### Quality Assurance 品質保證

- [ ] **Comprehensive Test Suite / 完整測試套件**
  - 功能: Unit and integration tests / 單元與整合測試
  - 預估工時 / Est. Hours: 6-8h
  - 子任務 / Subtasks:
    - [ ] Tool tests / 工具測試
    - [ ] API client tests / API 客戶端測試
    - [ ] Query engine tests / 查詢引擎測試

- [ ] **Improved Error Handling / 增強錯誤處理**
  - 功能: Better error messages and recovery /更好的錯誤訊息與恢復
  - 預估工時 / Est. Hours: 2-3h
  - 子任務 / Subtasks:
    - [ ] API error handling / API 錯誤處理
    - [ ] Tool error handling / 工具錯誤處理
    - [ ] User-friendly error messages / 使用者友善錯誤訊息

---

## 📊 進度統計 / Progress Statistics

| 項目 / Item | 數量 / Count |
|-------------|-------------|
| 已完成 / Completed | 12 |
| 待實作 (P0) / Pending (P0) | 8 主要任務 / main tasks |
| 待實作 (P1/P2) / Pending (P1/P2) | 5 主要任務 / main tasks |
| 完成率 / Completion Rate | ~60% |

---

## 🎯 下一步優先順序 / Next Priority Order

### Sprint 1: 核心工具 / Core Tools (建議先完成 / Recommended First)

1. Write Tool / Write 工具 (1h)
2. Glob Tool / Glob 工具 (1.5h)
3. Grep Tool / Grep 工具 (2h)

### Sprint 2: 環境建構 / Context Building

4. Context Building / 環境建構 (2-3h)
5. CLAUDE.md Discovery / CLAUDE.md 發現 (1h)

### Sprint 3: 權限系統 / Permission System

6. Permission System / 權限系統 (4-5h)
7. Permission Dialog / 權限對話框 (3-4h)

### Sprint 4: 會話管理 / Session Management

8. Session Storage / 會話儲存 (3-4h)
9. Resume Functionality / 恢復功能 (1h)

### Sprint 5: TUI 增強 / TUI Enhancement

10. Input Handling /輸入處理 (3-4h)
11. Keyboard Shortcuts / 鍵盤快捷鍵 (1h)

---

## 📝 實作筆記 / Implementation Notes

### Tool Interface Pattern / 工具介面模式

所有工具都遵循相同的介面 / All tools follow the same interface:

```go
type Tool interface {
    Name() string
    Description() string
    InputSchema() map[string]interface{}
    Execute(ctx context.Context, input map[string]interface{}, tc *ToolContext) (*ToolResult, error)
    IsReadOnly(input map[string]interface{}) bool
    IsConcurrencySafe(input map[string]interface{}) bool
}
```

### Session Storage Structure / Session 儲存結構

```
~/.claude/
├── sessions/
│   ├── <uuid>.jsonl              # Transcript / 對話記錄
│   └── <uuid>.metadata.json      # Metadata / 元資料
├── projects/
│   └── <cwd-hash>/
│       └── sessions/
└── settings.json                 # Global settings / 全域設定
```

---

*最後更新 / Last Updated: 2026-04-01*