// Builds web/dist; files under a/ are named by content.
import { execFileSync } from 'node:child_process'
import { createHash } from 'node:crypto'
import { copyFileSync, mkdirSync, readFileSync, renameSync, rmSync, statSync, writeFileSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const web = dirname(fileURLToPath(import.meta.url))
const root = dirname(web)
const dist = join(web, 'dist')

rmSync(dist, { recursive: true, force: true })
mkdirSync(join(dist, 'a'), { recursive: true })
execFileSync('go', ['build', '-trimpath', '-ldflags=-s -w', '-tags', 'mygo_noinspector', '-o',
  join(dist, 'app.wasm'), './cmd/gallery'], { cwd: root, stdio: 'inherit', env: { ...process.env, GOOS: 'js', GOARCH: 'wasm' } })
const goroot = execFileSync('go', ['env', 'GOROOT'], { encoding: 'utf8' }).trim()
copyFileSync(join(goroot, 'lib', 'wasm', 'wasm_exec.js'), join(dist, 'wasm_exec.js'))

// name moves a file into a/ by its content hash.
function name (base, file, ext) {
  const hash = createHash('sha256').update(readFileSync(join(dist, file))).digest('hex').slice(0, 12)
  const to = `a/${base}.${hash}.${ext}`
  renameSync(join(dist, file), join(dist, to))
  return to
}

const size = statSync(join(dist, 'app.wasm')).size
const wasm = name('app', 'app.wasm', 'wasm')
const exec = name('wasm_exec', 'wasm_exec.js', 'js')
writeFileSync(join(dist, 'main.js'),
  readFileSync(join(web, 'main.js'), 'utf8').replaceAll('__WASM__', wasm).replaceAll('__SIZE__', String(size)))
const main = name('main', 'main.js', 'js')
writeFileSync(join(dist, 'index.html'),
  readFileSync(join(web, 'index.html'), 'utf8').replace('__EXEC__', exec).replace('__MAIN__', main))
copyFileSync(join(web, '_headers'), join(dist, '_headers'))
console.log(`web/dist built: ${main} ${wasm} (${size} bytes)`)
