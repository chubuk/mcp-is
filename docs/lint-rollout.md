# Rollout: Go linters (golangci-lint) for mcp-is

You are implementing this plan in the mcp-is repository.
Do not invent scope. Follow the steps in order. Stop and report if a step is blocked.

## Goal

Add a high-signal lint gate to mcp-is:
- `make lint` as the single entry point
- the same check in CI (job `lint`)
- narrow enable-list (no "all linters")

## Non-goals

- Do not enable every golangci-lint linter.
- Do not enable `unused` (catalog of 282 MCP tools is intentionally "dead" in places).
- Do not add `revive` in this pass.
- Do not change business logic except for genuine P0 findings (see Step 4).
- Do not run `deploy` / `docker push` / anything destructive.

## Context (verified)

| Item | Value |
|---|---|
| Project root | `/home/chubuk/projects/mcp-is` |
| Go module | `github.com/chubuk/mcp-is` |
| go.mod | `go 1.25.5` |
| Size | ~777 `*.go`, ~212k LOC, ~314 `*_test.go` |
| CI | `.github/workflows/ci.yml` — job `lint` currently runs only `go vet ./...` |
| Makefile | has `build run test docker-up docker-down logs clean deb` — **no** `lint` |
| golangci-lint | **not installed**, **no** `.golangci.yml` |
| CONTRIBUTING.md | mentions `golangci-lint run` under "if installed" — make it required |
| Local Go | `go` + `gofmt` present; `golangci-lint` / `goimports` may be missing |

## Enable-list (final for this pass)

| Linter | Why |
|---|---|
| `govet` | already in CI, keep |
| `staticcheck` | real bugs, deprecated, dead stores |
| `errcheck` | ignored `err` — critical for auth/storage/webhooks |
| `ineffassign` | useless assignments |
| `gofmt` | formatting (via `formatters` in golangci-lint v2) |

**Toolchain:** golangci-lint **v2.15.0** (v1.x incompatible with `go 1.25.5`).

Leave disabled: `unused`, `revive`, style-only ST10xx packages (see config).

---

## Step 0 — Preflight (read-only)

1. Confirm `go version` and `go test ./...` currently pass (or note failures).
2. `ls .golangci.yml` — if it already exists, STOP and report (do not overwrite).
3. Read `Makefile`, `.github/workflows/ci.yml`, `CONTRIBUTING.md` Go section.

## Step 1 — Install golangci-lint (pinned)

Install a **pinned** stable version (do not use `latest`).

**IMPORTANT (2026-10-22):** golangci-lint **v1.x does not work** with `go.mod` `go 1.25.5`
(v1.64.8 is built with Go 1.24). Use **golangci-lint v2** (stable `v2.15.0` at time of writing).
Config schema is **v2** (see Step 2) — not the v1 format.

Examples (adjust version):

```bash
# binary install to ~/go/bin (user-local, no sudo)
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
  | sh -s -- -b "$(go env GOPATH)/bin" v2.15.0

export PATH="$(go env GOPATH)/bin:$PATH"
golangci-lint version
```

Record the exact version you installed in `docs/lint-rollout.md` (fill in below).

- Installed version: `v2.15.0`

## Step 2 — Add `.golangci.yml` (create at repo root) — **v2 schema**

```yaml
version: "2"

run:
  timeout: 5m
  tests: true

linters:
  default: none
  enable:
    - govet
    - staticcheck
    - errcheck
    - ineffassign
  settings:
    errcheck:
      check-blank: false
    staticcheck:
      checks:
        - all
        - "-ST1000"
        - "-ST1003"
        - "-ST1016"
        - "-ST1020"
        - "-ST1021"
        - "-ST1022"

formatters:
  enable:
    - gofmt

issues:
  max-issues-per-linter: 0
  max-same-issues: 0
```

v1 → v2 mapping used here:
- `linters.disable-all: true` → `linters.default: none`
- `linters-settings` → `linters.settings`
- `gofmt` linter → `formatters.enable: [gofmt]`
- `exclude-use-default: false` — dropped in v2 (defaults already narrow with `default: none`)

Notes:
- `default: none` keeps the signal high.
- ST10xx disabled = no naming/style holy wars in this pass.
- After the first run, if there are **false positives**, add narrow `linters.exclusions.rules`
  (file/line/linter specific). Do not disable a whole linter for one FP.

## Step 3 — Makefile

Add targets (keep existing targets intact):

```make
.PHONY: lint lint-fix

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run ./... --fix
```

Also add `lint` to the `.PHONY` line if it already exists (do not drop other phony names).

## Step 4 — First run + triage (critical)

```bash
golangci-lint run ./... 2>&1 | tee /tmp/mcp-is-lint-first.txt
echo "exit=$?"
```

Triage findings into:

| Priority | Meaning | Action |
|---|---|---|
| **P0** | real bugs: ignored err in auth/storage/webhooks, nil deref, broken error wrap | **Fix code** in this pass |
| **P1** | staticcheck correctness | fix if small; else note and continue |
| **P2** | ineffassign / style-ish | leave or fix if trivial |

Rules:
- P0 fixes must stay behavior-preserving except the actual bug fix.
- Do **not** mass-reformat unrelated files.
- If P0 count is large (>30), fix the worst (auth/storage/webhooks first) and list the rest
  in the final report — do not silently skip.
- Record summary in final report: counts by linter and priority.

## Step 5 — CI (`.github/workflows/ci.yml`)

Extend job `lint` **after** `go vet ./...` (keep go vet):

```yaml
      - run: go vet ./...

      - name: golangci-lint
        uses: golangci/golangci-lint-action@v9
        with:
          version: v2.15.0   # MUST match the version installed in Step 1 (golangci-lint v2)
```

**v2 required** — do not pin `v1.64.8` (incompatible with `go 1.25.5`).
If `golangci-lint-action@v9` is unavailable, use the latest official action that supports
golangci-lint v2 — pin the linter version either way.

Keep `test`, `integration`, `build`, `docker`, `deploy` jobs unchanged.

## Step 6 — Docs

`CONTRIBUTING.md` — in the Go section / setup:
- Change "Run linter (if installed)" → **required** step `make lint`.
- Mention `make lint-fix` for auto-fixable issues (gofmt).
- Note that CI runs the same golangci-lint version/config.

Optional short `docs/lint.md`:
- how to install
- `make lint` / `make lint-fix`
- enable-list rationale (errcheck for auth/webhooks; unused left off)

## Step 7 — Verify

1. `make lint` → runs; exit 0 after P0 fixes (or documented remaining P2 with `--new-from-rev=HEAD~1` if you must — prefer clean).
2. `go vet ./...` → clean.
3. `go test ./...` → pass (same as preflight).
4. Negative smoke (optional, revert after): introduce a discarded `err` in a scratch file → `make lint` fails.
5. Do **not** push unless the user asks. Leave changes in working tree or commit only if instructed.

## Step 8 — Final report

Write a short report to the user (chat), including:

1. Files created/modified (list).
2. golangci-lint version used.
3. First-run finding counts (P0/P1/P2) and what was fixed.
4. Remaining findings (if any) and why deferred.
5. Whether `make lint`, `go vet`, `go test` are green.
6. Next optional steps (revive, pre-commit hook, baseline cleanup).

## Deliverable checklist

- [ ] `.golangci.yml` at repo root (config above)
- [ ] `make lint` and `make lint-fix` work
- [ ] CI job `lint` runs `go vet` + `golangci-lint` (same version)
- [ ] CONTRIBUTING.md updated
- [ ] P0 findings fixed or explicitly listed
- [ ] Report delivered

## Open questions for the user (ask only if blocked)

- Preferred golangci-lint version if install is ambiguous
- Whether to commit changes
- Whether to enable `revive` in a follow-up

## Command cheat-sheet

```bash
cd /home/chubuk/projects/mcp-is
make lint
make lint-fix
go vet ./...
go test ./...
golangci-lint run ./...
```
