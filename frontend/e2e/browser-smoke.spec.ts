import { expect, test } from '@playwright/test'

test('sample graph evaluates in UI', async ({ page }) => {
  await page.goto('/')
  await page.getByRole('button', { name: 'Load sample debate graph' }).click()
  await page.getByRole('button', { name: 'Evaluate argument graph' }).click()
  await expect(page.getByText('Claim Ranking')).toBeVisible()
  await expect(page.getByText('Converged')).toBeVisible()
})
