# pf-developer-portal

学習用の OpenAPI ポータルです。手置きの YAML をカタログとリファレンスに描き、example からモック応答します。他プロダクトへのライブプロキシや Git clone はありません。**本番 API ポータルの置き換えではありません。**

```powershell
go test ./...
go run ./cmd/portal
```

- カタログ: http://localhost:8111
- 例: http://localhost:8111/docs/payments

モックはスキーマ検証あり、永続化なしです。必須フィールドや `Idempotency-Key` が無いと 400、未知パスは 404 です。

YAML に載っている仕様:

| slug | 内容 |
| --- | --- |
| `payments` | 架空のチャージ API |
| `commerce-catalog` | コマース商品のサブセット |
| `content-blog` | ブログ公開記事のサブセット |

テストで AWS キー風などの文字列が YAML に混ざらないようにしています。

Compose は `deploy/` です。ポータルは http://localhost:8111 です。

## OpenAPI の破壊的変更

```powershell
go run ./cmd/oasdiff-gate testdata\openapi\base.yaml testdata\openapi\compatible.yaml
go run ./cmd/oasdiff-gate testdata\openapi\base.yaml testdata\openapi\breaking.yaml
```

後者は終了コード 1 です。`.github/workflows/openapi-breaking.yml` が同じフィクスチャを CI で回します。他リポジトリへは `examples/oasdiff-action.yml` をコピーできます。

設計の詳細は [portfolio-plan](https://github.com/maeplego/portfolio-plan) の `portfolio-plan/developer-platform/docs/` です。
