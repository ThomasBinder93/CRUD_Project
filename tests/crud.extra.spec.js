import { test, expect } from '@playwright/test';

test.describe('CRUD Extra Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('persistence after reload', async ({ page }) => {
    const name = 'Persist ' + Date.now();

    await page.locator('#nameInput').fill(name);
    await page.locator('button').filter({ hasText: 'Erstellen' }).click();
    await expect(page.locator('#message')).toHaveText('Item erfolgreich erstellt');

    await page.reload();
    await expect(page.locator(`#list input[value="${name}"]`)).toHaveCount(1);
  });

  test('multiple rapid creates', async ({ page }) => {
    const names = [];
    for (let i = 0; i < 5; i++) {
      const name = `Bulk-${Date.now()}-${i}`;
      names.push(name);
      await page.locator('#nameInput').fill(name);
      await page.locator('button').filter({ hasText: 'Erstellen' }).click();
      await expect(page.locator('#message')).toHaveText('Item erfolgreich erstellt');
    }

    for (const n of names) {
      await expect(page.locator(`#list input[value="${n}"]`)).toHaveCount(1);
    }
  });

  test('server error handling on create', async ({ page }) => {
    await page.route('**/api/items', route => {
      if (route.request().method() === 'POST') {
        route.fulfill({
          status: 500,
          contentType: 'application/json',
          body: JSON.stringify({ error: { details: 'internal server error' } })
        });
      } else {
        route.continue();
      }
    });

    await page.locator('#nameInput').fill('ErrTest ' + Date.now());
    await page.locator('button').filter({ hasText: 'Erstellen' }).click();

    await expect(page.locator('#message')).toHaveText('Fehler: internal server error');
  });

  test('show error on update/delete 404', async ({ page }) => {
    await page.route('**/api/items/**', route => {
      const method = route.request().method();
      if (method === 'PUT' || method === 'DELETE') {
        route.fulfill({
          status: 404,
          contentType: 'application/json',
          body: JSON.stringify({ error: { details: 'not found' } })
        });
      } else {
        route.continue();
      }
    });

    const itemName = 'ToBeNotFound ' + Date.now();
    await page.locator('#nameInput').fill(itemName);
    await page.locator('button').filter({ hasText: 'Erstellen' }).click();
    await expect(page.locator('#message')).toHaveText('Item erfolgreich erstellt');

    const itemRow = page.locator('#list tr').filter({
      has: page.locator(`input[value="${itemName}"]`),
    });
    await expect(itemRow).toBeVisible();

    await itemRow.locator('button').filter({ hasText: 'Ändern' }).click();
    await expect(page.locator('#message')).toHaveText('Fehler: not found');

    await itemRow.locator('button').filter({ hasText: 'Löschen' }).click();
    await expect(page.locator('#message')).toHaveText('Fehler: not found');
  });

  test('xss sanitization (no dialogs executed)', async ({ page }) => {
    let dialogSeen = false;
    page.on('dialog', async dialog => {
      dialogSeen = true;
      await dialog.dismiss();
    });

    const payload = '<script>window.__XSS=1;</script>';
    await page.locator('#nameInput').fill(payload);
    await page.locator('button').filter({ hasText: 'Erstellen' }).click();
    await expect(page.locator('#message')).toHaveText('Item erfolgreich erstellt');

    await expect(page.locator(`#list input[value="${payload}"]`)).toHaveCount(1);
    expect(dialogSeen).toBeFalsy();
  });

  test('accessibility smoke: keyboard focus order', async ({ page }) => {
    await page.keyboard.press('Tab'); // focus name input
    await expect(page.locator('#nameInput')).toBeFocused();
    await page.keyboard.press('Tab'); // description
    await expect(page.locator('#descriptionInput')).toBeFocused();
    await page.keyboard.press('Tab'); // status
    await expect(page.locator('#statusInput')).toBeFocused();
    await page.keyboard.press('Tab'); // create button
    await expect(page.locator('button').filter({ hasText: 'Erstellen' })).toBeFocused();
  });
});
