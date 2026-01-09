import { test, expect } from '@playwright/test'

test.describe('Hooks Generator', () => {
  test('generates hooks from preset selection', async ({ page }) => {
    await page.goto('/hooks')

    // Verify preset cards are visible
    await expect(page.getByText('Light')).toBeVisible()
    await expect(page.getByText('Default')).toBeVisible()

    // Select Basic preset
    await page.getByText('Basic').click()

    // Verify tech stack checkboxes
    await expect(page.getByLabel('Go')).toBeChecked()
    await expect(page.getByLabel('TypeScript')).toBeChecked()

    // Generate
    await page.getByRole('button', { name: 'Generate Hooks' }).click()

    // Verify hook preview appears
    await expect(page.getByRole('heading', { name: 'Generated Hooks' })).toBeVisible()

    // Verify output controls
    await expect(page.getByRole('button', { name: 'Copy' })).toBeVisible()
    await expect(page.getByRole('button', { name: 'Download' })).toBeVisible()
    await expect(page.getByRole('button', { name: 'Download All' })).toBeVisible()

    // Verify hook content is present
    await expect(page.locator('pre')).toBeVisible()
  })
})
