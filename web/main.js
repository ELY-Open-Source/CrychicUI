// Streams the wasm, runs Go, shows crashes.
const stderr = []
const write = globalThis.fs.writeSync.bind(globalThis.fs)
globalThis.fs.writeSync = (fd, buf) => {
  if (fd === 2) stderr.push(new TextDecoder().decode(buf))
  return write(fd, buf)
}

function crash (msg) {
  const why = [msg, stderr.join('').slice(-4000)].filter(Boolean).join('\n\n')
  console.error('crychic: crash', why)
  document.getElementById('why').textContent = why
  document.getElementById('over').hidden = false
}
document.getElementById('reload').onclick = () => location.reload()

async function load () {
  const t0 = performance.now()
  const res = await fetch('__WASM__')
  if (!res.ok) throw new Error(`${res.url}: HTTP ${res.status}`)
  const total = Number('__SIZE__')
  const reader = res.body.getReader()
  const parts = []
  let got = 0
  for (;;) {
    const { done, value } = await reader.read()
    if (done) break
    parts.push(value)
    got += value.length
    const pct = Math.min(100, Math.floor(got / total * 100))
    document.getElementById('fill').style.width = pct + '%'
    document.getElementById('loading').setAttribute('aria-valuenow', pct)
  }
  const bytes = new Uint8Array(got)
  let off = 0
  for (const p of parts) {
    bytes.set(p, off)
    off += p.length
  }
  const go = new Go()
  const { instance } = await WebAssembly.instantiate(bytes, go.importObject)
  go.run(instance).then(() => crash('The Go program exited.'), (err) => crash(String(err && err.stack || err)))
  await new Promise((resolve, reject) => {
    const until = performance.now() + 30000
    const poll = () => {
      if (window.__crychic?.ready) return resolve()
      if (performance.now() > until) return reject(new Error('No frame within 30 s.'))
      requestAnimationFrame(poll)
    }
    poll()
  })
  const first = Math.round(performance.now() - t0)
  document.getElementById('loading').remove()
  window.__crychicBoot = { wasm: '__WASM__', bytes: got, first }
  console.info(`crychic: wasm ${(got / 1e6).toFixed(1)} MB, first frame ${first} ms`)
}

load().catch((err) => crash(String(err && err.stack || err)))
