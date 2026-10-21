// Ships files first, then the page, so none goes missing.
import { execSync } from 'node:child_process'
import { cpSync, mkdirSync, mkdtempSync, readFileSync, rmSync, writeFileSync } from 'node:fs'
import { tmpdir } from 'node:os'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const web = dirname(fileURLToPath(import.meta.url))
const site = 'https://crychicui.zhangyanghaha0407.workers.dev'
// --first deploys a Worker that has no live page yet.
const first = process.argv.includes('--first')

function fail (msg) {
  console.error('deploy: ' + msg)
  process.exit(1)
}

const run = (cmd, env = {}) => execSync(cmd, { cwd: web, stdio: 'inherit', env: { ...process.env, ...env } })
const root = dirname(web)
const out = (cmd, cwd = web) => execSync(cmd, { cwd, encoding: 'utf8' })

// patient retries check every 2 s for a minute.
async function patient (what, check) {
  for (let i = 0; i < 30; i++) {
    if (await check()) return
    await new Promise(resolve => setTimeout(resolve, 2000))
  }
  fail('gave up waiting for ' + what)
}

async function get (url) {
  const res = await fetch(url, { cache: 'no-store' })
  if (!res.ok) fail(`${url}: HTTP ${res.status}`)
  return Buffer.from(await res.arrayBuffer())
}

// upload makes a version from the assets in dir.
function upload (dir, message) {
  const text = out(`npx wrangler versions upload --assets "${dir}" --message "${message}"`)
  console.error(text)
  const id = text.match(/^Worker Version ID: (\S+)/m)?.[1]
  const preview = text.match(/^Version Preview URL: (\S+)/m)?.[1]
  if (!id || !preview) fail('wrangler printed no version id or preview URL')
  return { id, preview }
}

const promote = (id, message) => run(`npx wrangler versions deploy "${id}@100%" --yes --message "${message}"`)
const live = async main => (await fetch(site + '/', { cache: 'no-store' }).then(r => r.text())).includes(main)

if (out('git status --porcelain').trim()) fail('commit first; each version names its commit')
run('node build.mjs')
const rev = out('git rev-parse --short HEAD').trim()
const stats = out('go run ./cmd/gallery -stats', root).trim().split(/\s+/).map(Number)
if (stats.length !== 2 || !stats.every(Number.isFinite) || stats.some(v => v <= 0)) fail('gallery -stats printed no counts')
const [components, categories] = stats
const testEnv = base => ({
  BASE_URL: base,
  EXPECTED_COMPONENTS: String(components),
  EXPECTED_CATEGORIES: String(categories),
})
const mainOf = html => html.match(/src="(a\/main\.[0-9a-f]+\.js)"/)?.[1]
const main = mainOf(readFileSync(join(web, 'dist', 'index.html'), 'utf8'))
if (!main) fail('dist/index.html names no main script')

if (first) {
  run(`npx wrangler deploy --message "${rev}"`)
} else {
  // The live build's files ride along, for pages already open.
  const old = (await get(site + '/')).toString()
  const oldMain = mainOf(old)
  if (!oldMain) fail('the live page names no main script')
  for (const f of old.matchAll(/src="(a\/[^"]+)"/g)) writeFileSync(join(web, 'dist', f[1]), await get(`${site}/${f[1]}`))
  const wasm = readFileSync(join(web, 'dist', oldMain), 'utf8').match(/fetch\('(a\/app\.[0-9a-f]+\.wasm)'\)/)?.[1]
  if (!wasm) fail('the live main script names no wasm')
  writeFileSync(join(web, 'dist', wasm), await get(`${site}/${wasm}`))
  const stage = mkdtempSync(join(tmpdir(), 'crychic-'))
  mkdirSync(join(stage, 'a'))
  writeFileSync(join(stage, 'index.html'), old)
  cpSync(join(web, 'dist', 'a'), join(stage, 'a'), { recursive: true })
  cpSync(join(web, 'dist', '_headers'), join(stage, '_headers'))
  // Step one: the live page beside the new files.
  const files = upload(stage, `${rev} files`)
  promote(files.id, `${rev} files`)
  await patient('the new files to go live', async () => (await fetch(`${site}/${main}`, { cache: 'no-store' })).ok)
  rmSync(stage, { recursive: true, force: true })
  // Step two: the new page, tested on its preview first.
  const page = upload(join(web, 'dist'), rev)
  run('npx playwright test', testEnv(page.preview))
  promote(page.id, rev)
}
await patient('the new page to go live', () => live(main))
run('npx playwright test', testEnv(site))
console.log(`deploy: ${rev} live at ${site}`)
