import { brotliCompressSync } from 'node:zlib'
import { test, expect } from '@playwright/test'

// pageErrors keeps each page's console and page errors.
const pageErrors = new WeakMap()

test.beforeEach(async ({ page }) => {
  const seen = []
  pageErrors.set(page, seen)
  page.on('console', m => { if (m.type() === 'error') seen.push(m.text()) })
  page.on('pageerror', e => seen.push(e.message))
})

test.afterEach(async ({ page }) => {
  expect(pageErrors.get(page)).toEqual([])
})

async function open (page, hash = '') {
  await page.goto('/' + hash)
  await page.waitForFunction(() => window.__crychicBoot)
}

const where = page => page.evaluate(() => window.__crychic.page())

test('boots, reports, and logs its budget', async ({ page, request }) => {
  await open(page)
  expect(await where(page)).toBe('/')
  const stats = await page.evaluate(() => window.__crychic.stats())
  expect(stats.categories).toBe(38)
  expect(Array.isArray(await page.evaluate(() => window.__crychic.recent()))).toBe(true)
  expect(await page.locator('#mygo canvas').count()).toBeGreaterThan(0)
  const boot = await page.evaluate(() => window.__crychicBoot)
  const wasm = await (await request.get('/' + boot.wasm)).body()
  const br = brotliCompressSync(wasm).length
  console.log(`budget: wasm ${(wasm.length / 1e6).toFixed(1)} MB, brotli ${(br / 1e6).toFixed(1)} MB, ` +
    `first frame ${boot.first} ms, ${stats.components} components`)
})

test('deep link opens a category', async ({ page }) => {
  await open(page, '#/c/layout')
  expect(await where(page)).toBe('/c/layout')
})

test('the hash and history drive the router', async ({ page }) => {
  await open(page)
  await page.evaluate(() => { location.hash = '#/c/chart' })
  await page.waitForFunction(() => window.__crychic.page() === '/c/chart')
  await page.goBack()
  await page.waitForFunction(() => window.__crychic.page() === '/')
})
