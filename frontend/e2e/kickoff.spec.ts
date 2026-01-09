import { test, expect } from '@playwright/test'

test.describe('Kickoff Wizard', () => {
  test('completes full wizard flow and generates prompt', async ({ page }) => {
    await page.goto('/kickoff')

    // Step 1: Project Identity
    await expect(page.getByRole('heading', { name: 'Project Identity' })).toBeVisible()
    await page.getByPlaceholder('e.g., A task management app').fill('A todo app for developers')
    await page.getByRole('button', { name: 'Next' }).click()

    // Step 2: Success Criteria
    await expect(page.getByRole('heading', { name: 'Success Criteria' })).toBeVisible()
    await page.getByPlaceholder('e.g., Users can create').fill('Users can add and complete todos')
    await page.getByRole('button', { name: 'Next' }).click()

    // Step 3: Users & Roles
    await expect(page.getByRole('heading', { name: 'Users & Roles' })).toBeVisible()
    await page.getByPlaceholder('e.g., Anonymous visitors').fill('Authenticated users only')
    await page.getByRole('button', { name: 'Next' }).click()

    // Step 4: Data Sensitivity
    await expect(page.getByRole('heading', { name: 'Data Sensitivity' })).toBeVisible()
    await page.getByPlaceholder('e.g., User emails').fill('User emails, todo content')
    await page.getByPlaceholder('e.g., 2 years').fill('1 year')
    await page.getByRole('button', { name: 'Next' }).click()

    // Step 5: Auth Model
    await expect(page.getByRole('heading', { name: 'Auth Model' })).toBeVisible()
    await page.getByRole('combobox').selectOption('basic')
    await page.getByRole('button', { name: 'Next' }).click()

    // Step 6: Concurrency
    await expect(page.getByRole('heading', { name: 'Concurrency' })).toBeVisible()
    await page.getByPlaceholder('e.g., Multiple users').fill('Single user per account')
    await page.getByRole('button', { name: 'Next' }).click()

    // Step 7: Risks & Tradeoffs
    await expect(page.getByRole('heading', { name: 'Risks & Tradeoffs' })).toBeVisible()
    await page.locator('#top-risks').fill('Data loss')
    await page.locator('#mitigations').fill('Daily backups')
    await page.getByRole('button', { name: 'Next' }).click()

    // Step 8: Boundaries
    await expect(page.getByRole('heading', { name: 'Boundaries' })).toBeVisible()
    await page.getByPlaceholder('e.g., Task titles are public').fill('All data is private')
    await page.getByRole('button', { name: 'Next' }).click()

    // Step 9: Non-Goals
    await expect(page.getByRole('heading', { name: 'Non-Goals' })).toBeVisible()
    await page.getByPlaceholder('e.g., Mobile app').fill('Mobile app, team features')
    await page.getByRole('button', { name: 'Next' }).click()

    // Step 10: Constraints
    await expect(page.getByRole('heading', { name: 'Constraints' })).toBeVisible()
    await page.getByPlaceholder('e.g., Must ship in 2 weeks').fill('Ship in 1 week')

    // Generate prompt
    await page.getByRole('button', { name: 'Generate' }).click()

    // Verify preview appears
    await expect(page.getByRole('heading', { name: 'Generated Kickoff Prompt' })).toBeVisible()
    await expect(page.getByText('kickoff-prompt.md')).toBeVisible()

    // Verify copy button exists
    await expect(page.getByRole('button', { name: 'Copy' })).toBeVisible()

    // Verify download button exists
    await expect(page.getByRole('button', { name: 'Download' })).toBeVisible()

    // Verify prompt content contains input
    await expect(page.locator('pre')).toContainText('todo app')
  })
})
