import { test, expect } from '@playwright/test'

test.describe('Steering Generator', () => {
  test('generates steering files from config', async ({ page }) => {
    await page.goto('/steering')

    // Fill required fields
    await page.getByLabel('Project Name').fill('Test Project')
    await page.getByLabel('Project Description').fill('A test project')

    // Fill tech stack
    await page.getByLabel('Backend').fill('Go')
    await page.getByLabel('Frontend').fill('React')
    await page.getByLabel('Database').fill('PostgreSQL')

    // Enable conditional files
    await page.getByLabel('Include conditional steering').check()

    // Generate
    await page.getByRole('button', { name: 'Generate Steering Files' }).click()

    // Verify file preview appears
    await expect(page.getByRole('heading', { name: 'Generated Files' })).toBeVisible()

    // Verify tabs exist (at least product.md)
    await expect(page.getByRole('tab', { name: 'product.md' })).toBeVisible()

    // Verify output panel controls
    await expect(page.getByRole('button', { name: 'Copy' })).toBeVisible()
    await expect(page.getByRole('button', { name: 'Download' })).toBeVisible()
    await expect(page.getByRole('button', { name: 'Download All' })).toBeVisible()

    // Verify content contains project name
    await expect(page.locator('pre')).toContainText('Test Project')
  })
})
