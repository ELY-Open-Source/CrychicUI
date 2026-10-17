# AGENTS.md

CrychicUI: 1500 components for MyGo native UI, in MyGo's own look, with motion and sound cues. Go module
`git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI`.
Progress truth: `TASKS.md`. Requirements: `PRD.md` and `docs/prd/catalog-*.md` (both local, not committed).

## Stack

| Item | Value |
| --- | --- |
| Go | 1.27.1 |
| MyGo | `github.com/egoist/mygo v0.5.5`, pinned; CLI pinned as `go tool mygo` |
| Sound | `github.com/egoist/mygo/audio`, cues synthesized in `sound` |
| Icons | Lucide (ISC), `icon/svg` |
| Web | `cmd/gallery` built for `js/wasm`; MyGo's browser backend draws it with WebGL 2 |
| Hosting | Cloudflare Workers static assets from `web/dist`; code on Gitea `git.zacharyzhang.com` |

Read MyGo source before using an API: `go list -m -f '{{.Dir}}' github.com/egoist/mygo` (`ui/*.go`, `ui/doc.go`,
`docs/ui/*.md`, `audio/*.go`). Never guess.

## Layout

```text
sound/              synthesized cues, Play, mute, Recent
icon/               Lucide SVGs, parsed on first use
kit/                helpers shared by three or more components
<category>/         one component per file, package = folder name
gallery/            registry, demos (demos_<category>*.go), app view, tests
cmd/gallery/        one main for desktop and browser; hook_js.go for tests
web/                shell, build.mjs, deploy.mjs, wrangler.jsonc, _headers, Playwright tests
internal/archtest/  import and file rules
```

## Tiers

Imports point down only. Same-tier imports are forbidden. `internal/archtest` enforces it.

0 sound, icon · 1 kit · 2 layout, typography, button, display · 3 input, choice, feedback, motion, fx ·
4 form, datetime, color, navigation, overlay, data, chart, media ·
5 chat, agent, code, files, shop, finance, dashboard, schedule, project, mail, social, account, devtools, game, learn,
health, home, marketing, content, desktop, geo · 6 gallery · 7 cmd/gallery

## Look

- Colors from `c.Theme()` only, tints by `Mix` and `Alpha`. No hex in components.
- Sizes from `t.Space(n)`, `t.Rem(n)`, `t.Radius`. Pills use radius 999. Hairlines are 1 DIP of `Border`.
- One soft shadow at most on cards; larger only on popups. No glow.
- Lucide icons in the text color, 16 DIP in controls.
- Focus ring and disabled dimming are MyGo's. Never redraw them.

## Components

- `func Name(c *ui.Context, required...) ui.Element`; optional knobs in a final `NameOptions` (never variadic). Zero
  options give the default look. Containers take content as a final `func()`. Several actions: a `NameResult` of
  named `ui.Element` fields.
- Value-bound components take `*T` and call `SetChanged` on the returned element when input changes it.
- Stateful components get their key from the caller (`c.Key`); private state lives in `ui.Local`.
- Every component moves: pose, color animation, transition, content transition, keyframes, loop or animated drawing.
  Reduce Motion is MyGo's job.
- Every interactive component plays one cue on its action through `kit` or `sound.Play`, from a handled action:
  views may build several times for one event.
- Copy shown by a component comes from the caller; built-in defaults are English and overridable.
- Fail fast: panic on programmer errors. Never swallow an error.
- Composites log state changes with `kit.Log`. Never log typed text or secrets: log counts.
- No element or parts stored across builds; no goroutine touches elements. `go tool mygo vet ./...` passes.
- Files under 500 lines (lock files exempt). Comments one line, ten words at most, only when needed. No dead
  code, no versioned APIs.
- A category file holds one exported component, named in snake case (`KPIGrid` in `kpi_grid.go`), or private
  helpers only.

## Gallery entry

Each component adds one `Entry` in `gallery/demos_<category>*.go`: ID, Name, Category, Summary, Demo. Each file's list
joins `entries` in `gallery/demos.go`. The demo is a top-level function and reports every event through `log`. The
code panel shows the demo function's own source, cut from the embedded `demos*.go`.

## Tests

- `go vet ./... && go tool mygo vet ./... && go test ./...` and `GOOS=js GOARCH=wasm go vet ./...` before done.
- Behavior tests drive `ui.NewTester` and assert values; interactive tests assert the cue in `sound.Recent()`.
  Test binaries start muted, so tests make no sound.
  Advance the tester's clock; never sleep.
- Web: `npm install && npm test` in `web/` builds `web/dist` and runs Playwright on a local server;
  `BASE_URL=<url> npx playwright test` checks a deployed site. Tests read `window.__crychic`.

## Deploy

- `npm run deploy` in `web/` (clean tree only): live page beside new files first, then the new page tested on its
  preview URL, promoted, and tested live at https://crychicui.zhangyanghaha0407.workers.dev. The Worker's first
  deploy is `npm run deploy -- --first`.
- Rollback: `npx wrangler rollback <version-id>` in `web/`.

## Review

Each TASKS.md item is reviewed in a Herdr pane split to the right by
`codex --yolo -m gpt-6.1-sol -c model_reasoning_effort=xhigh -c service_tier=fast`. At most three rounds. The reviewer
reports defects against PRD and this file, ends with `VERDICT: PASS` or `VERDICT: FAIL`, and does not reject work for
missing enhancements.

## Secrets

The Gitea token lives only in the environment (`GITEA_TOKEN`). Never write it to a file, a remote URL, a commit or a
log. Pushes pass it as an HTTP header.
