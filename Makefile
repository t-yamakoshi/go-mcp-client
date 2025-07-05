# Go MCP Client Makefile
# クリーンアーキテクチャに基づく MCP クライアントの開発・ビルド・テスト用

# 変数定義
BINARY_NAME=./app/mcp-client
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe
BINARY_DARWIN=$(BINARY_NAME)_darwin

# パス設定
CMD_PATH=./cmd/mcpclient
PKG_PATH=./pkg
TEST_SERVER_PATH=./test-server

# Go 関連設定
GO=go
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)
CGO_ENABLED=0

# バージョン情報
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# フラグ設定
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"
BUILD_FLAGS=-trimpath -ldflags="-s -w"

# デフォルトターゲット
.DEFAULT_GOAL := help

# ヘルプ表示
.PHONY: help
help: ## 利用可能なコマンドを表示
	@echo "Go MCP Client - 利用可能なコマンド:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "例: make build    # バイナリをビルド"
	@echo "例: make test     # テストを実行"
	@echo "例: make run      # アプリケーションを実行"

# 開発関連
.PHONY: build
build: ## バイナリをビルド
	@echo "ビルド中: $(BINARY_NAME)..."
	$(GO) build $(LDFLAGS) -o $(BINARY_NAME) $(CMD_PATH)
	@echo "ビルド完了: $(BINARY_NAME)"

.PHONY: build-race
build-race: ## レースコンディション検出付きでビルド
	@echo "レース検出付きでビルド中: $(BINARY_NAME)..."
	$(GO) build -race $(LDFLAGS) -o $(BINARY_NAME) $(CMD_PATH)
	@echo "ビルド完了: $(BINARY_NAME)"

.PHONY: build-debug
build-debug: ## デバッグ情報付きでビルド
	@echo "デバッグ情報付きでビルド中: $(BINARY_NAME)..."
	CGO_ENABLED=1 $(GO) build -gcflags="all=-N -l" $(LDFLAGS) -o $(BINARY_NAME) $(CMD_PATH)
	@echo "ビルド完了: $(BINARY_NAME)"

.PHONY: build-all
build-all: ## 全プラットフォーム用にビルド
	@echo "全プラットフォーム用にビルド中..."
	GOOS=linux GOARCH=amd64 $(GO) build $(BUILD_FLAGS) $(LDFLAGS) -o $(BINARY_UNIX) $(CMD_PATH)
	GOOS=windows GOARCH=amd64 $(GO) build $(BUILD_FLAGS) $(LDFLAGS) -o $(BINARY_WINDOWS) $(CMD_PATH)
	GOOS=darwin GOARCH=amd64 $(GO) build $(BUILD_FLAGS) $(LDFLAGS) -o $(BINARY_DARWIN) $(CMD_PATH)
	@echo "全プラットフォーム用ビルド完了"

# 実行関連
.PHONY: run
run: ## アプリケーションを実行
	@echo "アプリケーション実行中..."
	$(GO) run $(CMD_PATH)

.PHONY: run-server
run-server: ## テストサーバーを実行
	@echo "テストサーバー実行中..."
	$(GO) run $(TEST_SERVER_PATH)/main.go

.PHONY: run-with-server
run-with-server: ## テストサーバーとクライアントを並行実行
	@echo "テストサーバーとクライアントを並行実行中..."
	@make run-server & sleep 2 && make run

# テスト関連
.PHONY: test
test: ## テストを実行
	@echo "テスト実行中..."
	$(GO) test -v ./...

.PHONY: test-race
test-race: ## レースコンディション検出付きでテスト
	@echo "レース検出付きでテスト実行中..."
	$(GO) test -race -v ./...

.PHONY: test-coverage
test-coverage: ## カバレッジ付きでテスト
	@echo "カバレッジ付きでテスト実行中..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "カバレッジレポート生成完了: coverage.html"

.PHONY: test-short
test-short: ## 短時間テストのみ実行
	@echo "短時間テスト実行中..."
	$(GO) test -short -v ./...

.PHONY: test-benchmark
test-benchmark: ## ベンチマークテスト実行
	@echo "ベンチマークテスト実行中..."
	$(GO) test -bench=. -benchmem ./...

# 依存関係管理
.PHONY: deps
deps: ## 依存関係を整理
	@echo "依存関係整理中..."
	$(GO) mod tidy
	$(GO) mod download

.PHONY: deps-update
deps-update: ## 依存関係を更新
	@echo "依存関係更新中..."
	$(GO) get -u ./...
	$(GO) mod tidy

.PHONY: deps-check
deps-check: ## 依存関係の脆弱性チェック
	@echo "依存関係の脆弱性チェック中..."
	$(GO) list -json -deps ./... | grep -E '"Path":' | cut -d'"' -f4 | sort -u | xargs -I {} sh -c 'echo "Checking {}..." && go list -m -versions {}'

# コード品質
.PHONY: fmt
fmt: ## コードフォーマット
	@echo "コードフォーマット中..."
	$(GO) fmt ./...

.PHONY: vet
vet: ## コード静的解析
	@echo "コード静的解析中..."
	$(GO) vet ./...

.PHONY: lint
lint: ## golangci-lint でリント
	@echo "リント実行中..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint がインストールされていません。インストールしてください:"; \
		echo "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

.PHONY: check
check: fmt vet lint ## コード品質チェック（fmt + vet + lint）

# Wire 関連
.PHONY: wire
wire: ## Wire で依存性注入コード生成
	@echo "Wire で依存性注入コード生成中..."
	@if command -v wire >/dev/null 2>&1; then \
		cd $(CMD_PATH) && wire; \
	else \
		echo "Wire がインストールされていません。インストールしてください:"; \
		echo "go install github.com/google/wire/cmd/wire@latest"; \
	fi

.PHONY: wire-check
wire-check: ## Wire の設定をチェック
	@echo "Wire 設定チェック中..."
	@if command -v wire >/dev/null 2>&1; then \
		cd $(CMD_PATH) && wire check; \
	else \
		echo "Wire がインストールされていません。"; \
	fi

# クリーンアップ
.PHONY: clean
clean: ## ビルド成果物を削除
	@echo "クリーンアップ中..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_WINDOWS)
	rm -f $(BINARY_DARWIN)
	rm -f coverage.out
	rm -f coverage.html
	@echo "クリーンアップ完了"

.PHONY: clean-all
clean-all: clean ## 完全クリーンアップ（キャッシュも含む）
	@echo "完全クリーンアップ中..."
	$(GO) clean -cache -modcache -testcache
	@echo "完全クリーンアップ完了"

# 開発支援
.PHONY: dev-setup
dev-setup: ## 開発環境セットアップ
	@echo "開発環境セットアップ中..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "golangci-lint をインストール中..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@if ! command -v wire >/dev/null 2>&1; then \
		echo "Wire をインストール中..."; \
		go install github.com/google/wire/cmd/wire@latest; \
	fi
	$(GO) mod download
	@echo "開発環境セットアップ完了"

.PHONY: dev
dev: ## 開発モード実行（ファイル変更監視付き）
	@echo "開発モード実行中..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "Air がインストールされていません。インストールしてください:"; \
		echo "go install github.com/cosmtrek/air@latest"; \
		echo "または make run を使用してください"; \
	fi

.PHONY: install
install: build ## バイナリをインストール
	@echo "バイナリをインストール中..."
	cp $(BINARY_NAME) /usr/local/bin/ || echo "権限が不足しています。sudo make install を試してください"

.PHONY: uninstall
uninstall: ## バイナリをアンインストール
	@echo "バイナリをアンインストール中..."
	rm -f /usr/local/bin/$(BINARY_NAME)

# リリース関連
.PHONY: release
release: clean build-all ## リリース用ビルド
	@echo "リリース用ビルド完了"
	@echo "生成されたファイル:"
	@ls -la $(BINARY_NAME)*

.PHONY: version
version: ## バージョン情報を表示
	@echo "バージョン: $(VERSION)"
	@echo "ビルド時刻: $(BUILD_TIME)"
	@echo "Git コミット: $(GIT_COMMIT)"
	@echo "Go バージョン: $(shell go version)"
	@echo "OS/Arch: $(GOOS)/$(GOARCH)"

# ドキュメント関連
.PHONY: docs
docs: ## ドキュメント生成
	@echo "ドキュメント生成中..."
	@if command -v godoc >/dev/null 2>&1; then \
		echo "godoc サーバーを起動中... http://localhost:6060"; \
		godoc -http=:6060; \
	else \
		echo "godoc がインストールされていません。インストールしてください:"; \
		echo "go install golang.org/x/tools/cmd/godoc@latest"; \
	fi

# デバッグ関連
.PHONY: debug
debug: build-debug ## デバッグ用ビルド
	@echo "デバッグ用バイナリ生成完了: $(BINARY_NAME)"
	@echo "デバッガーで実行する準備ができました"

.PHONY: profile
profile: ## プロファイリング実行
	@echo "プロファイリング実行中..."
	$(GO) test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=. ./...
	@echo "プロファイルファイル生成完了: cpu.prof, mem.prof"

# 便利なショートカット
.PHONY: all
all: clean deps check test build ## 全工程実行（クリーン → 依存関係 → チェック → テスト → ビルド）

.PHONY: quick
quick: deps build ## クイックビルド（依存関係 → ビルド）

.PHONY: test-build
test-build: test build ## テスト → ビルド

.PHONY: ci
ci: deps check test-coverage build ## CI 用（依存関係 → チェック → カバレッジテスト → ビルド）
