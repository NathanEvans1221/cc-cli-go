# TESTING - 測試策略與實作指南 / Testing Strategy & Implementation Guide

> **文件版本 / Version**: 1.0  
> **建立日期 / Created**: 2026-04-02  
> **適用專案 / Project**: cc-cli-go

---

## 📋 目錄 / Table of Contents

1. [測試策略 / Testing Strategy](#測試策略--testing-strategy)
2. [測試工具與技術 / Testing Tools & Techniques](#測試工具與技術--testing-tools--techniques)
3. [測試目錄結構 / Testing Directory Structure](#測試目錄結構--testing-directory-structure)
4. [單元測試 TODO List](#單元測試-todo-list)
5. [測試範例 / Test Examples](#測試範例--test-examples)
6. [如何執行測試 / How to Run Tests](#如何執行測試--how-to-run-tests)
7. [測試最佳實踐 / Testing Best Practices](#測試最佳實踐--testing-best-practices)

---

## 🎯 測試策略 / Testing Strategy

### 測試原則 / Testing Principles

1. **測試金字塔 / Testing Pyramid**
   - 單元測試（70%）：測試個別函數與方法
   - 整合測試（20%）：測試多個模組協作
   - E2E 測試（10%）：測試完整使用者流程

2. **FIRST 原則**
   - **Fast**：測試必須快速執行
   - **Independent**：測試之間互不依賴
   - **Repeatable**：測試結果可重現
   - **Self-validating**：測試自動判斷成功/失敗
   - **Timely**：測試與程式碼同步撰寫

3. **AAA 模式**
   - **Arrange**：準備測試環境與資料
   - **Act**：執行被測試的程式碼
   - **Assert**：驗證結果是否符合預期

### 測試覆蓋率目標 / Coverage Goals

| 模組 / Module            | 目標覆蓋率 / Target Coverage | 優先級 / Priority |
| ------------------------ | ---------------------------- | ----------------- |
| Tools（Write/Glob/Grep） | ≥ 80%                        | P0                |
| Permission System        | ≥ 85%                        | P0                |
| Session Storage          | ≥ 75%                        | P0                |
| Context Building         | ≥ 70%                        | P1                |
| API Client               | ≥ 60% (使用 mock)            | P1                |
| TUI Components           | ≥ 50%                        | P2                |

---

## 🔧 測試工具與技術 / Testing Tools & Techniques

### Go 測試框架 / Go Testing Framework

- **testing package**: Go 內建測試框架
- **testify**: Assertion 輔助庫（可選）
- **gomock**: Mock 生成工具（可選）
- **testfixtures**: 測試資料管理（可選）

### 測試命令 / Test Commands

```bash
# 執行所有測試
go test ./...

# 執行特定 package 測試
go test ./internal/tools/write

# 執行測試並顯示詳細輸出
go test -v ./...

# 執行測試並生成覆蓋率報告
go test -cover ./...

# 生成詳細覆蓋率報告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 執行效能測試
go test -bench=. ./...

# 執行競態檢測
go test -race ./...
```

---

## 📁 測試目錄結構 / Testing Directory Structure

```
internal/
├── tools/
│   ├── write/
│   │   ├── write.go
│   │   └── write_test.go          # Write Tool 測試
│   ├── glob/
│   │   ├── glob.go
│   │   └── glob_test.go            # Glob Tool 測試
│   ├── grep/
│   │   ├── grep.go
│   │   └── grep_test.go            # Grep Tool 測試
│   ├── bash/
│   │   ├── bash.go
│   │   └── bash_test.go            # Bash Tool 測試
│   ├── read/
│   │   ├── read.go
│   │   └── read_test.go            # Read Tool 測試
│   └── edit/
│   │   ├── edit.go
│   │   └── edit_test.go            # Edit Tool 測試
├── permission/
│   ├── types.go
│   ├── types_test.go               # Permission Types 測試
│   ├── dangerous.go
│   └── dangerous_test.go           # Dangerous Command 測試
├── session/
│   ├── session.go
│   └── session_test.go             # Session Storage 測試
├── context/
│   ├── context.go
│   ├── context_test.go             # Context Building 測試
│   ├── claudemd.go
│   └── claudemd_test.go            # CLAUDE.md Discovery 測試
├── api/
│   ├── client.go
│   ├── client_test.go              # API Client 測試 (使用 mock)
│   └── mock_client.go              # Mock API Client
├── testutil/
│   ├── testutil.go                 # 測試輔助工具
│   ├── filesystem.go               # 測試檔案系統輔助
│   └── mock.go                     # Mock 生成工具
tests/
├── integration/
│   └── integration_test.go         # 整合測試
└── e2e/
    └── e2e_test.go                 # E2E 測試
```

---

## ✅ 單元測試 TODO List

### 🔴 P0 - 必要測試 / Required Tests

#### Tools 工具測試

- [x] **Write Tool Tests**
  - [x] Test successful file creation
  - [x] Test file already exists error
  - [x] Test parent directory creation
  - [x] Test permission denied error
  - [x] Test empty file path error

- [x] **Glob Tool Tests**
  - [x] Test pattern matching with wildcard
  - [x] Test pattern matching with double star
  - [x] Test empty pattern error
  - [x] Test directory not found error
  - [x] Test result sorting

- [x] **Grep Tool Tests**
  - [x] Test regex pattern matching
  - [x] Test file type filtering
  - [x] Test empty pattern error
  - [x] Test invalid regex error
  - [x] Test result formatting
  - [x] Test large file handling

- [x] **Read Tool Tests**
  - [x] Test file reading with line numbers
  - [x] Test offset and limit parameters
  - [x] Test file not found error
  - [x] Test empty file handling

- [x] **Edit Tool Tests**
  - [x] Test successful edit
  - [x] Test old_string not found error
  - [x] Test multiple matches error
  - [x] Test exact string matching

#### Permission System 權限系統測試

- [x] **Permission Checker Tests**
  - [x] Test default mode behavior
  - [x] Test accept mode behavior
  - [x] Test plan mode behavior
  - [x] Test auto mode behavior
  - [x] Test rule matching
  - [x] Test pattern matching

- [x] **Dangerous Command Detection Tests**
  - [x] Test rm -rf detection
  - [x] Test DROP TABLE detection
  - [x] Test git push --force detection
  - [x] Test safe command handling
  - [x] Test all dangerous patterns

#### Session Storage 會話儲存測試

- [x] **Session Tests**
  - [x] Test session creation
  - [x] Test session save/load
  - [x] Test message addition
  - [x] Test session cleanup
  - [x] Test resume functionality
  - [x] Test JSONL format validation

---

### 🟡 P1 - 重要測試 / Important Tests

#### Context Building 環境建構測試

- [x] **Context Tests**
  - [x] Test Git branch detection
  - [x] Test Git status detection
  - [x] Test working directory detection
  - [x] Test system prompt generation

- [x] **CLAUDE.md Discovery Tests**
  - [x] Test upward directory search
  - [x] Test CLAUDE.md loading
  - [x] Test GEMINI.md loading
  - [x] Test multiple file merging

#### API Client 測試

- [x] **API Client Tests**
  - [x] Test streaming request (using mock)
  - [x] Test error handling
  - [x] Test request building
  - [x] Test response parsing

---

### 🟢 P2 - 可選測試 / Optional Tests

#### TUI Components 測試

- [x] **TUI Model Tests**
  - [x] Test input handling
  - [x] Test message rendering
  - [x] Test keyboard shortcuts

- [x] **Permission Dialog Tests**
  - [x] Test dialog rendering
  - [x] Test button selection
  - [x] Test decision handling

#### Integration Tests 整合測試

- [x] **Integration Tests**
  - [x] Test full conversation flow
  - [x] Test tool execution with permission
  - [x] Test session persistence flow

---

## 💡 測試範例 / Test Examples

### Write Tool 測試範例

```go
// internal/tools/write/write_test.go
package write

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/user-name/cc-cli-go/internal/tools"
)

func TestWriteTool_Execute_Success(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	content := "Hello, World!"

	tool := New()
	input := map[string]interface{}{
		"file_path": filePath,
		"content":   content,
	}
	tc := &tools.ToolContext{WorkingDir: tmpDir}

	// Act
	result, err := tool.Execute(context.Background(), input, tc)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.IsError {
		t.Errorf("expected success, got error: %s", result.Content)
	}

	// Verify file was created
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}

	if string(data) != content {
		t.Errorf("expected content %s, got %s", content, string(data))
	}
}

func TestWriteTool_Execute_FileExists(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "existing.txt")

	// Create existing file
	if err := os.WriteFile(filePath, []byte("old content"), 0644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	tool := New()
	input := map[string]interface{}{
		"file_path": filePath,
		"content":   "new content",
	}
	tc := &tools.ToolContext{WorkingDir: tmpDir}

	// Act
	result, err := tool.Execute(context.Background(), input, tc)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.IsError {
		t.Error("expected error for existing file")
	}

	if !contains(result.Content, "already exists") {
		t.Errorf("expected 'already exists' error, got: %s", result.Content)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
```

### Permission System 測試範例

```go
// internal/permission/types_test.go
package permission

import (
	"testing"
)

func TestChecker_Check_DefaultMode(t *testing.T) {
	// Arrange
	checker := NewChecker(ModeDefault)

	// Act & Assert - Read tool should be allowed by default
	decision := checker.Check("Read", map[string]interface{}{
		"file_path": "/test.txt",
	})

	if decision.Behavior != BehaviorAllow {
		t.Errorf("expected Allow for Read tool, got %s", decision.Behavior)
	}

	// Act & Assert - Bash tool should ask by default
	decision = checker.Check("Bash", map[string]interface{}{
		"command": "ls",
	})

	if decision.Behavior != BehaviorAsk {
		t.Errorf("expected Ask for Bash tool, got %s", decision.Behavior)
	}
}

func TestChecker_Check_AcceptMode(t *testing.T) {
	// Arrange
	checker := NewChecker(ModeAccept)

	// Act
	decision := checker.Check("Bash", map[string]interface{}{
		"command": "rm -rf /",
	})

	// Assert - Accept mode should allow everything
	if decision.Behavior != BehaviorAllow {
		t.Errorf("expected Allow in accept mode, got %s", decision.Behavior)
	}
}

func TestChecker_Check_DangerousCommand(t *testing.T) {
	// Arrange
	checker := NewChecker(ModeDefault)

	// Act
	decision := checker.Check("Bash", map[string]interface{}{
		"command": "rm -rf /important/data",
	})

	// Assert
	if decision.Behavior != BehaviorAsk {
		t.Errorf("expected Ask for dangerous command, got %s", decision.Behavior)
	}

	if !contains(decision.Reason, "dangerous") {
		t.Errorf("expected dangerous reason, got: %s", decision.Reason)
	}
}
```

### Session Storage 測試範例

```go
// internal/session/session_test.go
package session

import (
	"testing"
)

func TestSession_NewSession(t *testing.T) {
	// Arrange
	projectID := "/test/project"

	// Act
	session := NewSession(projectID)

	// Assert
	if session.ID == "" {
		t.Error("expected session ID to be generated")
	}

	if session.ProjectID != projectID {
		t.Errorf("expected project ID %s, got %s", projectID, session.ProjectID)
	}

	if len(session.Messages) != 0 {
		t.Error("expected empty messages list")
	}
}

func TestSession_AddMessage(t *testing.T) {
	// Arrange
	session := NewSession("/test")
	msg := types.NewUserMessage("test message")

	// Act
	session.AddMessage(msg)

	// Assert
	if len(session.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(session.Messages))
	}

	if session.Messages[0] != msg {
		t.Error("message not added correctly")
	}
}

func TestSession_SaveAndLoad(t *testing.T) {
	// Arrange
	session := NewSession("/test")
	session.AddMessage(types.NewUserMessage("test"))

	// Act - Save
	if err := session.Save(); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	// Act - Load
	loaded, err := LoadSession(session.ID)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	// Assert
	if loaded.ID != session.ID {
		t.Errorf("expected ID %s, got %s", session.ID, loaded.ID)
	}

	if len(loaded.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(loaded.Messages))
	}
}
```

---

## 🚀 如何執行測試 / How to Run Tests

### 本地開發環境 / Local Development

```bash
# 1. 執行所有測試
go test ./...

# 2. 執行特定模組測試
go test ./internal/tools/write

# 3. 執行測試並顯示詳細輸出
go test -v ./internal/permission

# 4. 執行測試並生成覆蓋率報告
go test -cover ./...

# 5. 生成 HTML 覆蓋率報告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# 6. 執行效能測試
go test -bench=. ./internal/tools

# 7. 執行競態檢測
go test -race ./...
```

---

## 📚 測試最佳實踐 / Testing Best Practices

### 1. 使用 t.TempDir() 管理測試檔案

```go
func TestFileOperation(t *testing.T) {
	// t.TempDir() 會在測試結束後自動清理
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	// ...
}
```

### 2. 使用 Table-driven Tests

```go
func TestDangerousCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected bool
	}{
		{"rm -rf", "rm -rf /data", true},
		{"DROP TABLE", "DROP TABLE users;", true},
		{"safe command", "ls -la", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDangerousCommand("Bash", map[string]interface{}{
				"command": tt.command,
			})
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
```

### 3. 使用 Mock 避免外部依賴

```go
// 使用介面而非直接依賴
type APIClient interface {
	Stream(ctx context.Context, req *Request) (<-chan Event, error)
}

// 測試中使用 Mock
type MockAPIClient struct {
	events []Event
}

func (m *MockAPIClient) Stream(ctx context.Context, req *Request) (<-chan Event, error) {
	ch := make(chan Event)
	go func() {
		for _, e := range m.events {
			ch <- e
		}
		close(ch)
	}()
	return ch, nil
}
```

### 4. 測試錯誤情況

```go
func TestErrorHandling(t *testing.T) {
	tool := New()

	// 測試空路徑
	result, _ := tool.Execute(context.Background(), map[string]interface{}{
		"file_path": "",
	}, nil)

	if !result.IsError {
		t.Error("expected error for empty file path")
	}
}
```

### 5. 清理測試資源

```go
func TestWithCleanup(t *testing.T) {
	file, err := os.Create("test.txt")
	if err != nil {
		t.Fatal(err)
	}

	// 使用 t.Cleanup 確保清理
	t.Cleanup(func() {
		file.Close()
		os.Remove("test.txt")
	})

	// 測試程式碼...
}
```

---

## 📊 測試報告範例 / Test Report Examples

### 覆蓋率報告

```
github.com/user-name/cc-cli-go/internal/tools/write/write.go      80.5%
github.com/user-name/cc-cli-go/internal/permission/types.go        85.2%
github.com/user-name/cc-cli-go/internal/session/session.go         75.8%
total:                                                            78.3%
```

### 測試執行結果

```
=== RUN   TestWriteTool_Execute_Success
--- PASS: TestWriteTool_Execute_Success (0.02s)
=== RUN   TestWriteTool_Execute_FileExists
--- PASS: TestWriteTool_Execute_FileExists (0.01s)
PASS
ok      github.com/user-name/cc-cli-go/internal/tools/write    0.150s
```

---

## 🔗 相關資源 / Related Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Go Testing Best Practices](https://golang.org/doc/effective_go#testing)
- [testify Assertion Library](https://github.com/stretchr/testify)
- [gomock Mock Generator](https://github.com/golang/mock)

---

_最後更新 / Last Updated: 2026-04-02_
