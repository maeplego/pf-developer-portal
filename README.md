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

### `PORTAL_*_URL`（Compose vs Kubernetes）

ホームの「他ツール」リンクは次の環境変数です。

| 変数 | Compose 例 | Kubernetes |
| --- | --- | --- |
| `PORTAL_CI_DASH_URL` | `http://localhost:8115` | 空のまま（ci-dash はクラスタ非搭載） |
| `PORTAL_REVIEW_URL` | `http://localhost:8118` | 空のまま |
| `PORTAL_SCANNER_URL` | 任意（ホストの docs URL など） | 空のまま |

K8s に載るのは **portal のみ**（overlay B/D/E）。ci-dash / review / scanner / CLI はホストまたは Compose 専用です。

## OpenAPI の破壊的変更

```powershell
go run ./cmd/oasdiff-gate testdata\openapi\base.yaml testdata\openapi\compatible.yaml
go run ./cmd/oasdiff-gate testdata\openapi\base.yaml testdata\openapi\breaking.yaml
```

後者は終了コード 1 です。`.github/workflows/openapi-breaking.yml` が同じフィクスチャを CI で回します。他リポジトリへは `examples/oasdiff-action.yml` をコピーできます。

設計の詳細は [portfolio-plan](https://github.com/maeplego/portfolio-plan) の `portfolio-plan/developer-platform/docs/` です。

## ライセンスと利用条件

本リポジトリは **デモ・学習・社内評価用** です。現状品質に **保証はありません**。

- 許可: クローン、ローカル実行、学習、非本番の評価
- 別契約が必要: 本番運用、有償サービスへの組込み、再販・托管の提供

詳細は [LICENSE](./LICENSE) と [licensing.md](https://github.com/maeplego/portfolio-plan/blob/master/portfolio-plan/licensing.md) を参照してください。

