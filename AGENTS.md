# AGENTS.md

CrychicUI: 1500 components for MyGo native UI, in MyGo's own look, with motion and sound cues. Go module
`github.com/ELY-Open-Source/CrychicUI`.
Progress truth: `TASKS.md`. Requirements: `PRD.md` and `docs/prd/catalog-*.md` (both local, not committed).

## Stack

| Item | Value |
| --- | --- |
| Go | 1.27.1 |
| MyGo | `github.com/egoist/mygo v0.5.5`, pinned; CLI pinned as `go tool mygo` |
| Sound | `github.com/egoist/mygo/audio`, cues, notes and drums synthesized in `sound`; nothing else opens audio |
| Icons | Lucide (ISC), `icon/svg` |
| Web | `cmd/gallery` built for `js/wasm`, shipped in 20 MiB parts (Workers cap files at 25 MiB); MyGo's browser backend draws it with WebGL 2 |
| Hosting | Cloudflare Workers static assets from `web/dist`; code on GitHub `github.com/ELY-Open-Source/CrychicUI` |

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
  options give the default look. Containers take content as a final `func()`, after any options. Several actions:
  named `ui.Element` results, root first; `mygo vet` forbids elements in structs.
- Value-bound components take `*T` and call `SetChanged` on the returned element when input changes it.
- Forms that submit return `(root, send)`; `send` stays disabled until the value is valid.
- Stateful components get their key from the caller (`c.Key`); private state lives in `ui.Local`.
- Every component moves: pose, color animation, transition, content transition, keyframes, loop or animated drawing.
  Reduce Motion is MyGo's job.
- Every interactive component plays one cue on its action through `kit` or `sound.Play`, from a handled action:
  views may build several times for one event.
- Copy shown by a component comes from the caller; built-in defaults are English and overridable.
- Fail fast: panic on programmer errors. Never swallow an error.
- Composites log state changes with `kit.Log`. Never log typed text or secrets: log counts.
- No element or parts stored across builds; no goroutine touches elements. `go tool mygo vet ./...` passes.
- A component deciding from laid-out sizes (`LayoutBounds`, a `ScrollState`'s `Max`) calls `kit.Settle` with the
  values it used; paint rebuilds when layout disagrees. Arrivals on first build use `kit.Arrive`; `Enter` only
  runs inside a shown parent.
- `Clip` clips children only: `Draw` or `DrawOver` paint that can pass the box goes through `kit.Inside`.
  `Painter.Line` draws only horizontal and vertical lines; stroke a `ui.Path` for others.
- Files under 500 lines (lock files exempt). Comments one line, ten words at most, only when needed. No dead
  code, no versioned APIs.
- A category file holds one exported component, named in snake case (`KPIGrid` in `kpi_grid.go`), or private
  helpers only.
- Quality bar is the layout–chat categories (`chat.ConversationList`): real structure (icon, avatar via
  `display.Avatar`, status badge, meta line), hover and selected colors animated, staggered `kit.Arrive` on lists,
  empty state where a list can empty, every catalog clause implemented. A bordered box of plain text fails.
- MyGo can shrink single-line text in an auto-width row by a fraction and ellipsize it: give such text
  `Shrink(0)` (cap long copy with `MaxWidth` plus `Wrap`). Text that should truncate needs `Grow(1).MinWidth(0)`.
- Panics name the component and caller label: `fmt.Sprintf("pkg: Name %q needs ...", label)`.

## Gallery entry

Each component adds one `Entry` in `gallery/demos_<category>*.go`: ID, Name, Category, Summary, Demo. Each file's list
joins `entries` in `gallery/demos.go`. The demo is a top-level function and reports every event through `log`. The
code panel shows the demo function's own source, cut from the embedded `demos*.go`. Demos sit in a fixed width
(320–520) and use realistic sample data: full names, real titles, 3–6 rows for lists, photos from `photos()`.
Filler such as "A", "Item", "Notes" fails review.

## Tests

- `go vet ./... && go tool mygo vet ./... && go test ./...` and `GOOS=js GOARCH=wasm go vet ./...` before done.
- Behavior tests drive `ui.NewTester` and assert values; interactive tests assert the cue in `sound.Recent()`.
  Test binaries start muted, so tests make no sound.
  Advance the tester's clock; never sleep.
- Web: `npm install && npm test` in `web/` builds `web/dist` and runs Playwright on a local server;
  `BASE_URL=<url> npx playwright test` checks a deployed site. Tests read `window.__crychic`.

## Deploy

- Each reviewed item is committed, pushed, then deployed before the next item starts.
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

Tokens (GitHub, Cloudflare) live only in the environment. Never write one to a file, a remote URL, a commit or a
log.
