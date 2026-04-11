import { expect, test } from '@playwright/test'

/**
 * Covers UI flows: empty state, node CRUD, edge validation + CRUD,
 * parameter edits, evaluate success/error, sample load + evaluate.
 * Pauses keep the recorded video readable (not ~1s long).
 */
test('full feature walkthrough with readable pacing', async ({ page }) => {
  test.setTimeout(180_000)

  await page.goto('/')
  await page.waitForTimeout(800)

  const nodePanel = page
    .locator('section.panel')
    .filter({ has: page.getByRole('heading', { name: 'Nodes' }) })
  const edgePanel = page
    .locator('section.panel')
    .filter({ has: page.getByRole('heading', { name: 'Edges' }) })
  const paramPanel = page
    .locator('section.panel')
    .filter({ has: page.getByRole('heading', { name: 'Evaluation Parameters' }) })

  // --- Initial empty state ---
  await expect(
    page.getByText('No nodes yet. Add claims or evidence.'),
  ).toBeVisible()
  await expect(
    page.getByRole('button', { name: 'Evaluate argument graph' }),
  ).toBeDisabled()
  await expect(
    page.getByText('No result yet. Build a graph and run evaluation.'),
  ).toBeVisible()
  await page.waitForTimeout(600)

  // --- Nodes: add two, edit labels & kind ---
  await nodePanel.getByRole('button', { name: 'Add Node' }).click()
  await page.waitForTimeout(500)
  await nodePanel.getByRole('button', { name: 'Add Node' }).click()
  await page.waitForTimeout(500)

  const nodeItems = nodePanel.locator('.list-item')
  await expect(nodeItems).toHaveCount(2)

  await nodeItems.nth(0).getByLabel('Label').fill('City center traffic restriction')
  await nodeItems.nth(1).getByLabel('Label').fill('Health impact evidence')
  await nodeItems.nth(1).getByLabel('Kind').selectOption('evidence')
  await page.waitForTimeout(500)

  // --- Nodes: add temporary node then remove (delete path) ---
  await nodePanel.getByRole('button', { name: 'Add Node' }).click()
  await page.waitForTimeout(400)
  await expect(nodeItems).toHaveCount(3)
  await nodeItems
    .nth(2)
    .getByLabel('Label')
    .fill('Temporary node (will be removed)')
  await page.waitForTimeout(400)
  await nodeItems.nth(2).getByRole('button', { name: 'Remove' }).click()
  await expect(nodeItems).toHaveCount(2)
  await page.waitForTimeout(600)

  // --- Edges: validation (missing endpoints) ---
  await edgePanel.getByRole('button', { name: 'Add edge' }).click()
  await expect(
    edgePanel.getByText('Please select both from and to nodes.'),
  ).toBeVisible()
  await page.waitForTimeout(700)

  // --- Edges: validation (self-loop) ---
  await edgePanel.getByLabel('from').selectOption('node_1')
  await edgePanel.getByLabel('to').selectOption('node_1')
  await edgePanel.getByRole('button', { name: 'Add edge' }).click()
  await expect(edgePanel.getByText('Self-loops are not allowed.')).toBeVisible()
  await page.waitForTimeout(700)

  // --- Edges: validation (weight out of range) ---
  await edgePanel.getByLabel('from').selectOption('node_2')
  await edgePanel.getByLabel('to').selectOption('node_1')
  await edgePanel.getByLabel('weight').fill('2')
  await edgePanel.getByRole('button', { name: 'Add edge' }).click()
  await expect(
    edgePanel.getByText('Weight must be a number between 0 and 1.'),
  ).toBeVisible()
  await edgePanel.getByLabel('weight').fill('0.8')
  await page.waitForTimeout(500)

  // --- Edges: valid add ---
  await edgePanel.getByLabel('from').selectOption('node_2')
  await edgePanel.getByLabel('to').selectOption('node_1')
  await edgePanel.getByLabel('weight').fill('0.8')
  await edgePanel.getByRole('button', { name: 'Add edge' }).click()
  await page.waitForTimeout(800)

  // --- Edges: delete then re-add ---
  await edgePanel.getByRole('button', { name: 'delete' }).first().click()
  await page.waitForTimeout(600)
  await edgePanel.getByLabel('from').selectOption('node_2')
  await edgePanel.getByLabel('to').selectOption('node_1')
  await edgePanel.getByLabel('relation').selectOption('support')
  await edgePanel.getByLabel('weight').fill('0.75')
  await edgePanel.getByRole('button', { name: 'Add edge' }).click()
  await page.waitForTimeout(800)

  // --- Parameters: all four fields ---
  await paramPanel.getByLabel(/damping/i).fill('0.72')
  await paramPanel.getByLabel(/baseline/i).fill('0.45')
  await paramPanel.getByLabel(/epsilon/i).fill('0.0005')
  await paramPanel.getByLabel(/maxIterations/i).fill('180')
  await page.waitForTimeout(1000)

  // --- Evaluate: success ---
  await page.getByRole('button', { name: 'Evaluate argument graph' }).click()
  await expect(page.getByText('Claim Ranking')).toBeVisible()
  await expect(page.getByText('Converged')).toBeVisible()
  await expect(page.getByRole('heading', { name: 'Warnings' })).toBeVisible()
  await expect(page.getByRole('heading', { name: 'Errors' })).toBeVisible()
  await page.waitForTimeout(1500)

  // --- Evaluate: validation error (empty label) ---
  await nodeItems.nth(0).getByLabel('Label').fill('')
  await page.getByRole('button', { name: 'Evaluate argument graph' }).click()
  await expect(page.getByRole('alert')).toBeVisible()
  await expect(page.getByRole('alert')).toContainText('graph validation failed')
  await page.waitForTimeout(1200)

  // --- Restore & re-evaluate ---
  await nodeItems
    .nth(0)
    .getByLabel('Label')
    .fill('City center traffic restriction')
  await page.getByRole('button', { name: 'Evaluate argument graph' }).click()
  await expect(page.getByText('Claim Ranking')).toBeVisible()
  await page.waitForTimeout(1200)

  // --- Sample loader + evaluate ---
  await page.getByRole('button', { name: 'Load sample debate graph' }).click()
  await page.waitForTimeout(800)
  await page.getByRole('button', { name: 'Evaluate argument graph' }).click()
  await expect(page.getByText('Claim Ranking')).toBeVisible()
  await page.waitForTimeout(2000)
})
