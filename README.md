# Dialectic Lab

Dialectic Lab は、**議論を重み付きグラフとしてモデル化し、主張の受容可能性を数値評価する Web アプリケーション**です。

- **Frontend**: React + TypeScript + Vite
- **Backend**: Go (`net/http`)

単なる CRUD ではなく、支持 (`support`) / 反駁 (`attack`) 関係を反復計算し、収束情報付きで結果を返すのが特徴です。

---

## 1. アプリ概要

ユーザーは以下を行えます。

1. ノード（主張 / 根拠）を作成・編集・削除
2. エッジ（support / attack, weight 0〜1）を設定
3. 評価パラメータ（damping, epsilon, maxIterations, baseline）を調整
4. Go バックエンドで評価を実行
5. ランキング、収束情報、診断メッセージを確認

---

## 2. ディレクトリ構成

```txt
.
├─ backend
│  ├─ cmd/server/main.go
│  └─ internal
│     ├─ domain
│     ├─ eval
│     └─ http
└─ frontend
   ├─ src
   │  ├─ components
   │  ├─ api.ts
   │  └─ types.ts
   └─ vite.config.ts
```

---

## 3. ローカル実行手順

前提:
- Node.js 22+
- npm 10+
- Go 1.22+

### 3.1 Backend

```bash
cd backend
go test ./...
go run ./cmd/server/main.go
```

デフォルトで `http://localhost:8080` で待ち受けます。  
`PORT` 環境変数で変更可能です。

### 3.2 Frontend

```bash
cd frontend
npm install
npm run lint
npm run test
npm run dev
```

デフォルトで `http://localhost:5173` で起動します。  
`/api` は Vite の proxy で `http://localhost:8080` に転送されます。

---

## 4. API

### `GET /api/health`

Response:

```json
{ "status": "ok" }
```

### `POST /api/evaluate`

Request:

```json
{
  "nodes": [
    { "id": "claim_core", "label": "Main claim", "kind": "claim" },
    { "id": "ev_a", "label": "Evidence A", "kind": "evidence" }
  ],
  "edges": [
    { "from": "ev_a", "to": "claim_core", "relation": "support", "weight": 0.8 }
  ],
  "params": {
    "damping": 0.72,
    "epsilon": 0.0005,
    "maxIterations": 180,
    "baseline": 0.5
  }
}
```

Successful response:

```json
{
  "scores": [
    { "nodeId": "claim_core", "score": 0.7342 },
    { "nodeId": "ev_a", "score": 0.5 }
  ],
  "meta": {
    "converged": true,
    "iterations": 9,
    "maxDelta": 0.00024
  },
  "diagnostics": {
    "warnings": [],
    "errors": []
  }
}
```

Validation error response (HTTP 400):

```json
{
  "message": "graph validation failed",
  "diagnostics": {
    "warnings": [],
    "errors": ["edge[0] weight must be in [0,1]"]
  }
}
```

---

## 5. 評価アルゴリズム（概要）

1. 各ノードを baseline で初期化
2. 各反復で流入エッジを集約
   - support: 加算
   - attack: 減算
3. damping を使って前回スコアと合成
4. `[0,1]` にクランプ
5. `maxDelta < epsilon` で収束判定、または `maxIterations` 到達で停止

---

## 6. セキュリティ / 品質に関する方針

- Go バックエンド:
  - 標準ライブラリ中心（`net/http`, `encoding/json`）で依存を最小化
  - 入力バリデーションを実施（重み範囲、ノード参照整合、自己ループ禁止など）
  - JSON エラー時に 400 を返し、内部例外を 500 と分離
- Frontend:
  - 型安全（TypeScript）
  - API エラー表示を明示
  - Lint / Test / Build を CI 相当で通過

---

## 7. 今後の拡張アイデア

- グラフ可視化（force-directed layout）
- セッション保存（SQLite/PostgreSQL）
- 複数シナリオ比較（パラメータ感度分析）
- 主要ノードへの寄与分解（explainability）
