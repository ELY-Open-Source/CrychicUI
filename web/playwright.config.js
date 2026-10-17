import { defineConfig } from '@playwright/test'

// BASE_URL runs the tests against a deployed site.
const url = process.env.BASE_URL

export default defineConfig({
  testDir: 'tests',
  timeout: 90000,
  use: { baseURL: url || 'http://localhost:4173', viewport: { width: 1280, height: 820 } },
  webServer: url
    ? undefined
    : { command: 'python3 -m http.server 4173 --directory dist', url: 'http://localhost:4173', reuseExistingServer: true }
})
