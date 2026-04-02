# Core Features - CC-CLI-Go

> **目的**: 定義 Go 版本需要實現的核心功能清單

---

## 決策摘要

| 項目 | 決策 |
|------|------|
| MCP | ❌ 不保留 |
| TUI Framework | Bubble Tea |
| Provider 支援 | 只支援 Anthropic API |
| 目標 | 學習導向的精簡版 CLI |

---

## 功能分類

### ✅ Phase 1: 核心功能 (必須實現)

#### 1.1 CLI Entry & Argument Parsing

**來源**: `src/entrypoints/cli.tsx`, `src/main.tsx`

**功能**:
- [x] 基本命令列參數解析
- [x] `-p` / `--print` 非互動模式
- [x] `--version` 版本顯示
- [x] `--help` 幫助信息
- [x] `-c` / `--continue` 繼續上次會話
- [x] `--resume <session-id>` 恢復指定會話
- [x] `--model <model>` 指定模型
- [x] `[prompt]` 直接輸入提示

**Go 技術**: `cobra` + `pflag`

**預估工時**: 2-3 小時

---

#### 1.2 Anthropic API Streaming Client

**來源**: `src/services/api/claude.ts`

**功能**:
- [x] 建構 API request (system prompt, messages, tools)
- [x] SSE streaming 連接
- [x] Event parsing:
  - [x] `message_start`
  - [x] `content_block_start` (text, thinking, tool_use)
  - [x] `content_block_delta` (text, input_json_delta)
  - [x] `content_block_stop`
  - [x] `message_delta` (stop_reason, usage)
  - [x] `message_stop`
  - [x] `ping`
  - [x] `error`
- [x] Error handling (rate limit, overload, network)
- [x] Token usage tracking

**Go 技術**: `net/http`, `bufio.Scanner`, `encoding/json`

**API 端點**: `https://api.anthropic.com/v1/messages`

**Request Headers**:
```
x-api-key: <API_KEY>
anthropic-version: 2023-06-01
anthropic-beta: prompt-caching-2024-07-31
anthropic-beta: token-efficient-tools-2025-02-19
content-type: application/json
```

**預估工時**: 4-6 小時

---

#### 1.3 Query Loop

**來源**: `src/query.ts`

**功能**:
- [x] Message 管理 (append, compact boundary)
- [x] System prompt 注入
- [x] Streaming response 處理
- [x] Tool use 檢測與觸發
- [x] Tool result 收集
- [x] 多輪對話循環
- [x] Interrupt handling (Ctrl+C)

**核心邏輯**:
```
while (true) {
    1. Build request (messages + system prompt + tools)
    2. Call API streaming
    3. Yield events to UI
    4. Collect tool_uses
    5. If tool_uses:
        - Execute tools
        - Append tool_results
        - Continue loop
    6. Else:
        - Return terminal
}
```

**預估工時**: 4-5 小時

---

#### 1.4 Tool System - Core Tools

**來源**: `src/tools/`, `src/Tool.ts`

**必要工具**:

| Tool | 功能 | 優先級 | 預估工時 |
|------|------|--------|---------|
| **Bash** | 執行 shell 命令 | P0 | 2h |
| **Read** | 讀取檔案內容 | P0 | 1h |
| **Edit** | 編輯檔案 (字串替換) | P0 | 2h |
| **Write** | 寫入新檔案 | P0 | 1h |
| **Glob** | 檔案模式匹配 | P0 | 1.5h |
| **Grep** | 搜尋檔案內容 | P0 | 2h |

**Tool Interface (Go)**:
```go
type Tool interface {
    Name() string
    Description() string
    InputSchema() map[string]interface{}  // JSON Schema
    Execute(ctx context.Context, input map[string]interface{}, tc *ToolContext) (*ToolResult, error)
    IsReadOnly(input map[string]interface{}) bool
    IsConcurrencySafe(input map[string]interface{}) bool
    CheckPermissions(input map[string]interface{}, ctx *ToolContext) PermissionResult
    UserFacingName(input map[string]interface{}) string
}
```

**預估工時**: 8-10 小時

---

#### 1.5 Permission System

**來源**: `src/utils/permissions/`, `src/hooks/useCanUseTool.tsx`

**功能**:
- [x] Permission modes: `default`, `accept`, `plan`, `auto`
- [x] Permission rules (allow, deny, ask)
- [x] Rule matching (tool name + input pattern)
- [x] Interactive permission dialog (TUI)
- [x] Rule persistence (settings.json)
- [x] Dangerous command detection

**Permission Rule 格式**:
```json
{
  "toolName": "Bash",
  "inputPattern": "git *",
  "behavior": "allow"
}
```

**預估工時**: 4-5 小時

---

#### 1.6 Context Building

**來源**: `src/context.ts`, `src/utils/claudemd.ts`

**功能**:
- [x] Git status (`git status --porcelain`)
- [x] Git branch (`git branch --show-current`)
- [x] Current working directory
- [x] Date/time
- [x] CLAUDE.md 發現與讀取
  - [x] 從當前目錄往上搜尋
  - [x] 合併所有找到的 CLAUDE.md
- [x] Memory files (可選)

**預估工時**: 2-3 小時

---

#### 1.7 Session Storage

**來源**: `src/utils/sessionStorage.ts`, `src/utils/conversationRecovery.ts`

**功能**:
- [x] Session ID 生成
- [x] Transcript 儲存 (JSONL 格式)
- [x] Session metadata
- [x] Session 列表查詢
- [x] Session 恢復 (resume)
- [x] Session 清理 (舊會話)

**檔案結構**:
```
~/.claude/
├── sessions/
│   ├── <uuid>.jsonl
│   └── <uuid>.metadata.json
├── projects/
│   └── <cwd-hash>/
│       └── sessions/
└── settings.json
```

**預估工時**: 3-4 小時

---

### ✅ Phase 2: TUI (必須實現)

#### 2.1 REPL Screen (Bubble Tea)

**來源**: `src/screens/REPL.tsx`

**功能**:
- [x] 主界面佈局
- [x] 用戶輸入區
- [x] 訊息列表顯示
- [x] Loading spinner
- [x] Status bar (model, tokens, etc.)

**Bubble Tea 模型**:
```go
type Model struct {
    messages   []Message
    input      textinput.Model
    viewport   viewport.Model
    spinner    spinner.Model
    ready      bool
}

func (m Model) Init() tea.Cmd
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m Model) View() string
```

**預估工時**: 6-8 小時

---

#### 2.2 Message Rendering

**來源**: `src/components/Messages.tsx`, `src/components/MessageRow.tsx`

**功能**:
- [x] User message rendering
- [x] Assistant message rendering (text, thinking)
- [x] Tool use message rendering (collapsed/expanded)
- [x] Tool result message rendering
- [x] System message rendering
- [x] Code syntax highlighting (可選)
- [x] Markdown rendering (可選)

**預估工時**: 4-5 小時

---

#### 2.3 Permission Dialog

**來源**: `src/components/permissions/`

**功能**:
- [x] Tool use 請求顯示
- [x] Allow / Deny / Always Allow / Always Deny 按鈕
- [x] Input preview (檔案路徑, 命令等)
- [x] 鍵盤導航

**預估工時**: 3-4 小時

---

#### 2.4 Input Handling

**來源**: `src/components/PromptInput/`

**功能**:
- [x] Multi-line input
- [x] History navigation (上/下鍵)
- [x] Auto-completion (可選)
- [x] Paste handling
- [x] Keyboard shortcuts:
  - [x] Enter: Submit
  - [x] Ctrl+C: Interrupt
  - [x] Ctrl+D: Exit
  - [x] Escape: Cancel current action

**預估工時**: 3-4 小時

---

### ⚪ Phase 3: 進階功能 (可選)

#### 3.1 Context Compaction

**來源**: `src/services/compact/`

**功能**:
- [x] Auto-compact when context exceeds threshold
- [x] Manual `/compact` command
- [x] Summary generation

**預估工時**: 4-6 小時

**優先級**: 中 (對長對話有幫助)

---

#### 3.2 Additional Tools

**可選工具**:

| Tool | 功能 | 優先級 |
|------|------|--------|
| Agent | Subagent spawning | 低 |
| TodoWrite | Task tracking | 低 |
| WebFetch | Fetch URLs | 低 |
| WebSearch | Web search | 低 |
| NotebookEdit | Jupyter notebook editing | 非常低 |

**預估工時**: 依工具而定

---

#### 3.3 Configuration Management

**來源**: `src/utils/config.ts`, `src/utils/settings/`

**功能**:
- [x] Global settings (`~/.claude/settings.json`)
- [x] Project settings (`.claude/settings.json`)
- [x] Settings schema validation
- [x] Settings hot-reload

**預估工時**: 2-3 小時

**優先級**: 中

---

## ❌ 明確移除的功能

### 移除原因

| 功能 | 移除原因 |
|------|---------|
| **MCP (Model Context Protocol)** | 複雜度高，學習導向可省略 |
| **Computer Use** | 需要 native 模組，複雜 |
| **Analytics / GrowthBook / Sentry** | 隱私考量，學習不需要 |
| **Voice Mode** | 需要 audio capture native |
| **Magic Docs** | 已移除功能 |
| **LSP Server** | 複雜，非核心 |
| **Plugins / Marketplace** | 複雜，非核心 |
| **Bridge Mode / Remote Control** | 企業功能，非核心 |
| **Background Sessions / Daemon** | 複雜，非核心 |
| **Templates** | 企業功能，非核心 |
| **Coordinator / Kairos** | 實驗功能 |
| **SSH Remote** | 複雜，非核心 |
| **Direct Connect** | 複雜，非核心 |
| **Multi-provider (Bedrock/Vertex/Azure)** | 只支援 Anthropic API |

---

## 功能實現優先順序

### Sprint 1: 基礎架構 (Day 1-2)

```
1. CLI entry + cobra setup
2. Config file structure
3. API client (basic streaming)
4. Simple query loop (no tools)
```

### Sprint 2: 核心工具 (Day 3-4)

```
5. Tool system framework
6. Bash tool
7. Read tool
8. Edit tool
9. Write tool
```

### Sprint 3: 搜尋工具 (Day 5)

```
10. Glob tool
11. Grep tool
```

### Sprint 4: 權限系統 (Day 6)

```
12. Permission system
13. Permission dialog (TUI)
```

### Sprint 5: Context & Session (Day 7)

```
14. Context building
15. CLAUDE.md discovery
16. Session storage
17. Resume functionality
```

### Sprint 6: TUI Polish (Day 8-9)

```
18. REPL screen
19. Message rendering
20. Input handling
21. Keyboard shortcuts
```

### Sprint 7: 測試 & 優化 (Day 10)

```
22. Integration testing
23. Error handling
24. Performance optimization
25. Documentation
```

---

## 功能清單總覽

### 必要功能 (P0) - 25 項

| # | 功能 | 預估工時 |
|---|------|---------|
| 1 | CLI argument parsing | 2-3h |
| 2 | API streaming client | 4-6h |
| 3 | Query loop | 4-5h |
| 4 | Bash tool | 2h |
| 5 | Read tool | 1h |
| 6 | Edit tool | 2h |
| 7 | Write tool | 1h |
| 8 | Glob tool | 1.5h |
| 9 | Grep tool | 2h |
| 10 | Permission system | 4-5h |
| 11 | Context building | 2-3h |
| 12 | CLAUDE.md discovery | 1h |
| 13 | Session storage | 3-4h |
| 14 | Resume functionality | 1h |
| 15 | REPL screen (Bubble Tea) | 6-8h |
| 16 | Message rendering | 4-5h |
| 17 | Permission dialog | 3-4h |
| 18 | Input handling | 3-4h |
| 19 | Interrupt handling | 1h |
| 20 | Error display | 1h |
| 21 | Token tracking | 1h |
| 22 | Model selection | 0.5h |
| 23 | Help/version display | 0.5h |
| 24 | Settings management | 2-3h |
| 25 | Logging | 1h |

**總計**: ~50-65 小時

---

### 可選功能 (P1/P2) - 8 項

| # | 功能 | 優先級 | 預估工時 |
|---|------|--------|---------|
| 1 | Context compaction | P1 | 4-6h |
| 2 | Code syntax highlighting | P1 | 2-3h |
| 3 | Markdown rendering | P1 | 2-3h |
| 4 | Agent tool (subagent) | P2 | 4-6h |
| 5 | TodoWrite tool | P2 | 1-2h |
| 6 | WebFetch tool | P2 | 1-2h |
| 7 | Auto-completion | P2 | 2-3h |
| 8 | Settings hot-reload | P2 | 1h |

**總計**: ~17-26 小時

---

## 驗收標準

### Phase 1 完成標準

- [x] 可以啟動 CLI 並顯示幫助
- [x] 可以輸入提示並獲得回應
- [x] 可以執行 Bash 命令 (經過權限確認)
- [x] 可以讀取、編輯、寫入檔案
- [x] 可以搜尋檔案 (Glob, Grep)
- [x] 可以恢復之前的會話
- [x] 可以看到 CLAUDE.md 的內容被注入

### Phase 2 完成標準

- [x] TUI 顯示正常
- [x] 可以滾動查看歷史訊息
- [x] 可以用鍵盤操作 (導航、提交、中斷)
- [x] 權限對話框可以正常運作

### 最終驗收

- [x] 所有核心工具可以正常運作
- [x] 長對話不會出現記憶體問題
- [x] 錯誤處理完善
- [x] 文檔完整

---

*Document generated for Go porting reference*