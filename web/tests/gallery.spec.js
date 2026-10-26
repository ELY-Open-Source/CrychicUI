import { test, expect } from '@playwright/test'

// Each category's demo and a cue point the tester found.
const clicks = [
  ['/c/layout/sidebar-layout', 268, 44], ['/c/typography/clamp', 28, 92], ['/c/button/solid', 28, 28],
  ['/c/display/list-item', 28, 28], ['/c/input/password-field', 780, 60], ['/c/choice/check', 28, 28],
  ['/c/feedback/dismissible-alert', 476, 44], ['/c/motion/shared-move', 44, 108], ['/c/fx/confetti', 28, 92],
  ['/c/form/form-card', 252, 188], ['/c/datetime/date-range-picker', 28, 28], ['/c/color/well', 28, 28],
  ['/c/navigation/navbar', 124, 44], ['/c/overlay/modal', 28, 28], ['/c/data/expand-table', 44, 76],
  ['/c/chart/line', 44, 28], ['/c/media/image-view', 28, 284], ['/c/chat/bubble', 28, 92],
  ['/c/agent/prompt-box', 44, 92], ['/c/code/code-block', 460, 44], ['/c/files/file-list', 44, 76],
  ['/c/shop/product-card', 44, 364], ['/c/finance/balance', 364, 44], ['/c/dashboard/revenue', 188, 76],
  ['/c/schedule/availability', 92, 60], ['/c/project/board', 220, 172], ['/c/mail/inbox', 44, 76],
  ['/c/social/post', 44, 412], ['/c/account/sign-in', 316, 236], ['/c/devtools/log-view', 460, 44],
  ['/c/game/dice', 172, 156], ['/c/learn/flashcard', 28, 28], ['/c/health/water', 44, 108],
  ['/c/home/light', 364, 60], ['/c/marketing/cta', 236, 140], ['/c/content/doc-toc', 28, 76],
  ['/c/desktop/title-tabs', 348, 44], ['/c/geo/map-canvas', 476, 44],
]

const errors = []

test.beforeEach(async ({ page }) => {
  errors.length = 0
  page.on('console', m => { if (m.type() === 'error') errors.push(m.text()) })
  page.on('pageerror', e => errors.push(e.message))
})

test.afterEach(() => { expect(errors).toEqual([]) })

async function boot (page, hash) {
  await page.goto('/' + hash)
  await page.waitForFunction(() => window.__crychicBoot)
}

async function go (page, path) {
  await page.evaluate(p => { location.hash = '#' + p }, path)
  // shown is empty on the not-found page.
  await page.waitForFunction(p => window.__crychic.shown() === p, path)
}

test('a deep link opens a component and its demo', async ({ page }) => {
  await boot(page, '#/c/game/dice')
  await page.waitForFunction(() => window.__crychic.shown() === '/c/game/dice')
  await page.waitForFunction(() => window.__crychic.demo().w > 0)
})

test('every category route opens', async ({ page }) => {
  await boot(page, '')
  const keys = [...new Set(clicks.map(([p]) => p.split('/')[2]))]
  expect(keys.length).toBe((await page.evaluate(() => window.__crychic.stats())).categories)
  for (const k of keys) await go(page, '/c/' + k)
  await page.evaluate(() => { location.hash = '#/c/not-a-category' })
  await page.waitForFunction(() => window.__crychic.page() === '/c/not-a-category')
  expect(await page.evaluate(() => window.__crychic.shown())).toBe('')
})

test('a click on one demo per category plays a cue', async ({ page }) => {
  test.setTimeout(300000)
  await boot(page, '')
  for (const [path, dx, dy] of clicks) {
    await go(page, path)
    await page.waitForFunction(() => window.__crychic.demo().w > 0)
    // Arrivals finish before the click, as in the Go tester.
    await page.waitForTimeout(2000)
    const d = await page.evaluate(() => window.__crychic.demo())
    const before = JSON.stringify(await page.evaluate(() => window.__crychic.recent()))
    await page.mouse.click(d.x + dx + 6, d.y + dy + 6)
    await expect.soft.poll(async () => JSON.stringify(await page.evaluate(() => window.__crychic.recent())),
      { message: path, timeout: 3000 }).not.toBe(before)
  }
})
