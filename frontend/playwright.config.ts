import { defineConfig } from '@playwright/test'

export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: 0,
  workers: 1,
  reporter: 'line',
  timeout: 180_000,
  use: {
    baseURL: 'http://localhost:5173',
    trace: 'off',
    viewport: { width: 1280, height: 720 },
    video: { mode: 'on', size: { width: 1280, height: 720 } },
    launchOptions: {
      // Slows actions so demo videos are not ~1s long
      slowMo: 120,
    },
  },
})
