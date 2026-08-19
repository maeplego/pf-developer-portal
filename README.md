# pf-developer-portal

P11 OpenAPI developer portal (idea 29). **Hand-placed YAML** is rendered as a catalog + reference, and a **mock** answers from examples. No GitHub clone, no exploit/PoC, no live proxy to other products.

Learning / portfolio sample. Formal docs: `project/portfolio-plan/developer-platform/docs/`. CI dashboard and code review live in other P11 repos and are not this slice.

## Demo

```powershell
go test ./...
go run ./cmd/portal
# http://localhost:8111
# http://localhost:8111/docs/payments
# curl.exe -s http://localhost:8111/health
```

Mock (example-first, request schema validated, no persistence):

```powershell
curl.exe -s -D - -X POST http://localhost:8111/mock/payments/v1/charges -H "Content-Type: application/json" -H "Idempotency-Key: demo-charge-001" -d "{\"amountMinor\":1299,\"currency\":\"JPY\",\"source\":\"tok_demo_visa\"}"
```

Missing `amountMinor` or `Idempotency-Key` returns **400**. Unknown paths return **404**.

## Specs in this repo

| slug | file | why |
| --- | --- | --- |
| `payments` | `specs/payments-v1.yaml` | idea 29 fictional charges API |
| `commerce-catalog` | `specs/commerce-catalog-v1.yaml` | P06 catalog subset, hand-placed |
| `content-blog` | `specs/content-blog-v1.yaml` | P08 public posts subset, hand-placed |

Examples are linted in tests so AWS-key / GitHub PAT / PEM shapes cannot ship in YAML.

## Compose

```powershell
copy deploy\.env.example deploy\.env
docker compose -f deploy\compose.yaml --env-file deploy\.env up --build
```

Portal: http://localhost:8111

Overlay B（Docker Desktop Kubernetes。他 overlay と同時に載せない）:

```powershell
cd ..\pf-cloud-k8s
.\scripts\cluster-smoke-b-collab.ps1
```

http://portal.localhost

## HTTP

| Method | Path | Role |
| --- | --- | --- |
| GET | `/health` `/ready` | liveness / ready |
| GET | `/` | catalog HTML |
| GET | `/docs/{slug}` | reference + Try it out |
| GET | `/api/catalog` | JSON list |
| GET | `/api/specs/{slug}` | raw YAML |
| * | `/mock/{slug}/...` | example mock |

## Not in this slice

oasdiff Action, CI pipeline dashboard, PR review UI, spec upload admin, Git sync, API keys.
