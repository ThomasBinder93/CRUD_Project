import { test, expect } from '@playwright/test';

test.describe('CRUD Application', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should load the page', async ({ page }) => {
    await expect(page).toHaveTitle('CRUD Demo');
    await expect(page.locator('h1')).toHaveText('CRUD Demo');
    await expect(page.locator('#nameInput')).toBeVisible();
    await expect(page.locator('button').filter({ hasText: 'Erstellen' })).toBeVisible();
  });

  test('should create a new item', async ({ page }) => {
    const itemName = 'Test Item ' + Date.now();

    // Fill the input
    await page.locator('#nameInput').fill(itemName);

    // Click create button
    await page.locator('button').filter({ hasText: 'Erstellen' }).click();

    // Check success message
    await expect(page.locator('#message')).toHaveText('Item erfolgreich erstellt');
    await expect(page.locator('#message')).toHaveCSS('color', 'rgb(0, 128, 0)'); // green

    // Check item appears in list
    await expect(page.locator(`#list input[value="${itemName}"]`)).toHaveCount(1);
  });

  test('should show error for empty item creation', async ({ page }) => {
    // Leave input empty
    await page.locator('#nameInput').clear();

    // Click create button
    await page.locator('button').filter({ hasText: 'Erstellen' }).click();

    // Check error message
    await expect(page.locator('#message')).toHaveText('Fehler: item name cannot be empty');
    await expect(page.locator('#message')).toHaveCSS('color', 'rgb(255, 0, 0)'); // red

    // Check no item was added (assuming list starts empty or check count)
    const listItems = page.locator('#list tr');
    const initialCount = await listItems.count();
    // After error, count should remain the same
    await expect(listItems).toHaveCount(initialCount);
  });

  test('should update an item', async ({ page }) => {
    const originalName = 'Original Item ' + Date.now();
    const updatedName = 'Updated Item ' + Date.now();

    // Create an item first
    await page.locator('#nameInput').fill(originalName);
    await page.locator('button').filter({ hasText: 'Erstellen' }).click();
    await expect(page.locator('#message')).toHaveText('Item erfolgreich erstellt');

    // Find the item and update it
    const itemRow = page.locator('#list tr').filter({ hasText: originalName });
    await expect(itemRow).toBeVisible();

    // Update the input field for that item
    await itemRow.locator('input').fill(updatedName);

    // Click update button
    await itemRow.locator('button').filter({ hasText: 'Ändern' }).click();

    // Check success message
    await expect(page.locator('#message')).toHaveText('Item erfolgreich aktualisiert');

    // Check item was updated
    await expect(page.locator('#list')).toContainText(updatedName);
    await expect(page.locator('#list')).not.toContainText(originalName);
  });

  test('should delete an item', async ({ page }) => {
    const itemName = 'Item to Delete ' + Date.now();

    // Create an item first
    await page.locator('#nameInput').fill(itemName);
    await page.locator('button').filter({ hasText: 'Erstellen' }).click();
    await expect(page.locator('#message')).toHaveText('Item erfolgreich erstellt');

    // Find the item and delete it
    const itemRow = page.locator('#list tr').filter({ hasText: itemName });
    await expect(itemRow).toBeVisible();

    // Click delete button
    await itemRow.locator('button').filter({ hasText: 'Löschen' }).click();

    // Check success message
    await expect(page.locator('#message')).toHaveText('Item erfolgreich gelöscht');

    // Check item was removed
    await expect(page.locator('#list')).not.toContainText(itemName);
  });

  test('should handle update of empty name', async ({ page }) => {
    const originalName = 'Item to Update ' + Date.now();

    // Create an item first
    await page.locator('#nameInput').fill(originalName);
    await page.locator('button').filter({ hasText: 'Erstellen' }).click();
    await expect(page.locator('#message')).toHaveText('Item erfolgreich erstellt');

    // Find the item and try to update to empty
    const itemRow = page.locator('#list tr').filter({ hasText: originalName });
    await itemRow.locator('input').clear();

    // Click update button
    await itemRow.locator('button').filter({ hasText: 'Ändern' }).click();

    // Check error message
    await expect(page.locator('#message')).toHaveText('Fehler: item name cannot be empty');

    // Check item still has original name
    await expect(page.locator('#list')).toContainText(originalName);
  });
});