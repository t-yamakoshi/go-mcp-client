# Go MCP クライアント

クリーンアーキテクチャの原則に従って構築された Model Context Protocol (MCP) クライアントの Go 実装です。

## 概要

このプロジェクトは、AI アシスタントが外部データソースやツールに標準化されたインターフェースを通じて接続することを可能にする Model Context Protocol のクライアントを実装しています。アプリケーションは保守性、テスタビリティ、関心の分離を確保するためにクリーンアーキテクチャの原則に従って構造化されています。

## アーキテクチャ

このプロジェクトは以下の層でクリーンアーキテクチャの原則に従っています：

### 1. ドメイン層 (`pkg/domain/`)
- **エンティティ** (`entity/`): コアビジネスオブジェクト（Message、Tool、Connection など）
- **リポジトリ** (`repository/`): データアクセスの抽象インターフェース
- **サービス** (`service/`): ビジネスロジックの抽象インターフェース
- **設定** (`config/`): アプリケーション設定の構造体
- **レスポンス** (`response/`): プロトコルレスポンスの構造体

### 2. ユースケース層 (`pkg/usecase/`)
- **ビジネスロジック**: アプリケーション固有のビジネスルール
- **オーケストレーション**: 異なるドメインサービス間の調整
- **入力/出力ポート**: アプリケーションが外部システムとどのように相互作用するかを定義

### 3. インターフェース層 (`pkg/interfaces/`)
- **CLI ハンドラー**: コマンドラインインターフェースの実装
- **HTTP ハンドラー**: HTTP インターフェースの実装（将来の拡張用）
- **コントローラー**: ユーザー入力の処理と出力のフォーマット

### 4. インフラストラクチャ層 (`pkg/infrastructure/`)
- **MCP リポジトリ**: MCP プロトコルの WebSocket 実装
- **設定リポジトリ**: ファイルベースの設定ストレージ
- **外部サービス**: データベース、外部 API など

## 機能

- WebSocket ベースの MCP サーバーとの通信
- MCP プロトコルバージョン 2024-11-05 のサポート
- 設定可能なクライアント設定
- グレースフルシャットダウン処理
- 拡張可能なメッセージハンドラーシステム
- クリーンアーキテクチャ設計
- 依存性注入（Wire を使用）
- 関心の分離

## インストール

```bash
git clone <repository-url>
cd go-mcp-client
go mod tidy
```

## 使用方法

### 基本的な使用方法

```bash
# デフォルト設定で実行
go run cmd/mcpclient/main.go

# カスタムサーバーURLで実行
go run cmd/mcpclient/main.go -server ws://localhost:3000

# カスタム設定ファイルで実行
go run cmd/mcpclient/main.go -config my-config.json
```

### 設定

クライアントは JSON 設定ファイルを使用して設定できます。設定ファイルが提供されない場合、デフォルト設定が作成されます。

設定例 (`config.json`):

```json
{
  "server_url": "ws://localhost:3000",
  "client_info": {
    "name": "go-mcp-client",
    "version": "1.0.0"
  },
  "log_level": "info"
}
```

### コマンドライン引数

- `-config`: 設定ファイルのパス（デフォルト: `config.json`）
- `-server`: MCP サーバーURL（設定ファイルを上書き）

## プロジェクト構造

```
go-mcp-client/
├── cmd/
│   └── mcpclient/                    # メインエントリポイント（DI設定付き）
│       ├── main.go
│       ├── wire.go                   # Wire設定
│       └── wire_gen.go               # 自動生成されたDIコード
├── pkg/
│   ├── domain/                       # ドメイン層（最も内側）
│   │   ├── entity/                   # エンティティ
│   │   │   ├── message.go            # メッセージ関連
│   │   │   ├── tool.go               # ツール関連
│   │   │   ├── connection.go         # 接続関連
│   │   │   └── info.go               # クライアント・サーバー情報
│   │   ├── repository/               # リポジトリインターフェース
│   │   │   ├── mcp_repository.go     # MCPリポジトリ
│   │   │   └── config_repository.go  # 設定リポジトリ
│   │   ├── service/                  # サービスインターフェース
│   │   │   ├── mcp_service.go        # MCPサービス
│   │   │   └── config_service.go     # 設定サービス
│   │   ├── config/                   # 設定構造体
│   │   │   └── config.go
│   │   └── response/                 # レスポンス構造体
│   │       └── response.go
│   ├── usecase/                      # ユースケース層
│   │   ├── mcp_usecase.go            # MCPビジネスロジック
│   │   └── config_usecase.go         # 設定管理ビジネスロジック
│   ├── interfaces/                   # インターフェース層
│   │   └── handler.go                # CLI/HTTPハンドラー
│   └── infrastructure/               # インフラストラクチャ層（最も外側）
│       ├── mcp_repository.go         # WebSocket実装
│       └── config_repository.go      # ファイルベース設定ストレージ
├── test-server/                      # テスト用MCPサーバー
│   └── main.go
├── go.mod
└── README.md
```

## 依存関係の流れ

依存関係の流れはクリーンアーキテクチャの原則に従います：

```
Interfaces → Use Cases → Domain ← Infrastructure
     ↓           ↓         ↑           ↑
   (Input)   (Business)  (Core)    (External)
```

- **インターフェース**は**ユースケース**に依存
- **ユースケース**は**ドメイン**インターフェースに依存
- **インフラストラクチャ**は**ドメイン**インターフェースを実装
- **ドメイン**は他の層に依存しない

## 開発

### ビルド

```bash
go build -o mcp-client ./cmd/mcpclient
```

### テストの実行

```bash
go test ./...
```

### 新しい機能の追加

1. **ドメイン層**: エンティティとインターフェースを定義
2. **ユースケース層**: ビジネスロジックを実装
3. **インターフェース層**: ユーザーインターフェースハンドラーを追加
4. **インフラストラクチャ層**: 外部統合を実装

## MCP プロトコルサポート

このクライアントは以下の MCP プロトコル機能をサポートしています：

- 接続確立と初期化
- カスタムハンドラーによるメッセージ処理
- ツール一覧とツール呼び出し（フレームワーク準備完了）
- Ping/pong ハートビート機構

## 実際の使用例

### 1. テストサーバーとの接続

```bash
# テストサーバーを起動
go run test-server/main.go

# 別のターミナルでクライアントを接続
./mcp-client -server ws://localhost:3000
```

### 2. Claude Desktop の MCP サーバーとの接続

```bash
# Claude Desktop が起動している場合
./mcp-client -server ws://localhost:3000
```

### 3. カスタム MCP サーバーとの接続

```bash
./mcp-client -server ws://your-mcp-server.com:3000
```

## クリーンアーキテクチャの利点

- **テスタビリティ**: 各層を独立してテスト可能
- **保守性**: 関心の分離による明確な責任分担
- **柔軟性**: 実装の簡単な入れ替えが可能
- **スケーラビリティ**: コンポーネント間の明確な境界
- **独立性**: ビジネスロジックが外部フレームワークに依存しない

## 今後の拡張予定

1. **テストの追加**: 各層のユニットテスト
2. **ログシステム**: 構造化ログの実装
3. **エラーハンドリング**: より詳細なエラー処理
4. **HTTP API**: RESTful API の追加
5. **メトリクス**: パフォーマンス監視
6. **設定検証**: より厳密な設定バリデーション

## 貢献

1. リポジトリをフォーク
2. 機能ブランチを作成
3. クリーンアーキテクチャの原則に従って変更を加える
4. 新機能のテストを追加
5. プルリクエストを送信

## ライセンス

[ライセンスを追加してください] 
