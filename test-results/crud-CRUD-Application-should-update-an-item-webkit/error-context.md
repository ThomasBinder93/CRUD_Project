# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: crud.spec.js >> CRUD Application >> should update an item
- Location: tests\crud.spec.js:50:7

# Error details

```
Error: expect(locator).toBeVisible() failed

Locator: locator('#list tr').filter({ hasText: 'Original Item 1779030863085' })
Expected: visible
Timeout: 5000ms
Error: element(s) not found

Call log:
  - Expect "toBeVisible" with timeout 5000ms
  - waiting for locator('#list tr').filter({ hasText: 'Original Item 1779030863085' })

```

```yaml
- heading "CRUD Demo" [level=1]
- textbox "Name"
- button "Erstellen"
- table:
  - row "test Ändern Löschen":
    - cell "test":
      - textbox: test
    - cell "Ändern":
      - button "Ändern"
    - cell "Löschen":
      - button "Löschen"
  - row "Test Item Ändern Löschen":
    - cell "Test Item":
      - textbox: Test Item
    - cell "Ändern":
      - button "Ändern"
    - cell "Löschen":
      - button "Löschen"
  - row "Original Item 1779030401741 Ändern Löschen":
    - cell "Original Item 1779030401741":
      - textbox: Original Item 1779030401741
    - cell "Ändern":
      - button "Ändern"
    - cell "Löschen":
      - button "Löschen"
  - row "Item to Update 1779030401737 Ändern Löschen":
    - cell "Item to Update 1779030401737":
      - textbox: Item to Update 1779030401737
    - cell "Ändern":
      - button "Ändern"
    - cell "Löschen":
      - button "Löschen"
  - row "Item to Delete 1779030401737 Ändern Löschen":
    - cell "Item to Delete 1779030401737":
      - textbox: Item to Delete 1779030401737
    - cell "Ändern":
      - button "Ändern"
    - cell "Löschen":
      - button "Löschen"
  - row "Test Item 1779030401737 Ändern Löschen":
    - cell "Test Item 1779030401737":
      - textbox: Test Item 1779030401737
    - cell "Ändern":
      - button "Ändern"
    - cell "Löschen":
      - button "Löschen"
  - row "Test Item 1779030403660 Ändern Löschen":
    - cell "Test Item 1779030403660":
      - textbox: Test Item 1779030403660
    - cell "Ändern":
      - button "Ändern"
    - cell "Löschen":
      - button "Löschen"
  - row "Test Item 1779030404055 Ändern Löschen":
    - cell "Test Item 1779030404055":
      - textbox: Test Item 1779030404055
    - cell "Ändern":
      - button "Ändern"
    - cell "Löschen":
      - button "Löschen"
  - row "Item to Delete 1779030404310 Ändern Löschen":
    - cell "Item to Delete 1779030404310":
      - textbox: Item to Delete 1779030404310
    - cell "Ändern":
      - button "Ändern"
    - cell "Löschen":
      - button "Löschen"
  - row "Original Item 1779030405054 Ändern Löschen":
    - cell "Original Item 1779030405054":
      - textbox: Original Item 1779030405054
    - cell "Ändern":
      - button "Ändern"
    - cell "Löschen":
      - button "Löschen"
  - row "Item to Update 1779030406578 Ändern Löschen":
    - cell "Item to Update 1779030406578":
      - textbox: Item to Update 1779030406578
    - cell "Ändern":
      - button "Ändern"
    - cell "Löschen":
      - button "Löschen"
  - row "Original Item 1779030408037 Ändern Löschen":
    - cell "Original Item 1779030408037":
      - textbox: Original Item 1779030408037
    - cell "Ändern":
      - button "Ändern"
    - cell "Löschen":
      - button "Löschen"
  - row "Item to Update 1779030408152 Ändern Löschen":
    - cell "Item to Update 1779030408152":
      - textbox: Item to Update 1779030408152
    - cell "Ändern":
      - button "Ändern"
    - cell "Löschen":
      - button "Löschen"
  - row "Item to Delete 1779030408151 Ändern Löschen":
    - cell "Item to Delete 1779030408151":
      - textbox: Item to Delete 1779030408151
    - cell "Ändern":
      - button "Ändern"
    - cell "Löschen":
      - button "Löschen"
  - row "Original Item 1779030862400 Ändern Löschen":
    - cell "Original Item 1779030862400":
      - textbox: Original Item 1779030862400
    - cell "Ändern":
      - button "Ändern"
    - cell "Löschen":
      - button "Löschen"
  - row "Original Item 1779030863085 Ändern Löschen":
    - cell "Original Item 1779030863085":
      - textbox: Original Item 1779030863085
    - cell "Ändern":
      - button "Ändern"
    - cell "Löschen":
      - button "Löschen"
```

# Test source

```ts
  1   | import { test, expect } from '@playwright/test';
  2   | 
  3   | test.describe('CRUD Application', () => {
  4   |   test.beforeEach(async ({ page }) => {
  5   |     await page.goto('/');
  6   |   });
  7   | 
  8   |   test('should load the page', async ({ page }) => {
  9   |     await expect(page).toHaveTitle('CRUD Demo');
  10  |     await expect(page.locator('h1')).toHaveText('CRUD Demo');
  11  |     await expect(page.locator('#nameInput')).toBeVisible();
  12  |     await expect(page.locator('button').filter({ hasText: 'Erstellen' })).toBeVisible();
  13  |   });
  14  | 
  15  |   test('should create a new item', async ({ page }) => {
  16  |     const itemName = 'Test Item ' + Date.now();
  17  | 
  18  |     // Fill the input
  19  |     await page.locator('#nameInput').fill(itemName);
  20  | 
  21  |     // Click create button
  22  |     await page.locator('button').filter({ hasText: 'Erstellen' }).click();
  23  | 
  24  |     // Check success message
  25  |     await expect(page.locator('#message')).toHaveText('Item erfolgreich erstellt');
  26  |     await expect(page.locator('#message')).toHaveCSS('color', 'rgb(0, 128, 0)'); // green
  27  | 
  28  |     // Check item appears in list
  29  |     await expect(page.locator('#list')).toContainText(itemName);
  30  |   });
  31  | 
  32  |   test('should show error for empty item creation', async ({ page }) => {
  33  |     // Leave input empty
  34  |     await page.locator('#nameInput').clear();
  35  | 
  36  |     // Click create button
  37  |     await page.locator('button').filter({ hasText: 'Erstellen' }).click();
  38  | 
  39  |     // Check error message
  40  |     await expect(page.locator('#message')).toHaveText('Fehler: item name cannot be empty');
  41  |     await expect(page.locator('#message')).toHaveCSS('color', 'rgb(255, 0, 0)'); // red
  42  | 
  43  |     // Check no item was added (assuming list starts empty or check count)
  44  |     const listItems = page.locator('#list tr');
  45  |     const initialCount = await listItems.count();
  46  |     // After error, count should remain the same
  47  |     await expect(listItems).toHaveCount(initialCount);
  48  |   });
  49  | 
  50  |   test('should update an item', async ({ page }) => {
  51  |     const originalName = 'Original Item ' + Date.now();
  52  |     const updatedName = 'Updated Item ' + Date.now();
  53  | 
  54  |     // Create an item first
  55  |     await page.locator('#nameInput').fill(originalName);
  56  |     await page.locator('button').filter({ hasText: 'Erstellen' }).click();
  57  |     await expect(page.locator('#message')).toHaveText('Item erfolgreich erstellt');
  58  | 
  59  |     // Find the item and update it
  60  |     const itemRow = page.locator('#list tr').filter({ hasText: originalName });
> 61  |     await expect(itemRow).toBeVisible();
      |                           ^ Error: expect(locator).toBeVisible() failed
  62  | 
  63  |     // Update the input field for that item
  64  |     await itemRow.locator('input').fill(updatedName);
  65  | 
  66  |     // Click update button
  67  |     await itemRow.locator('button').filter({ hasText: 'Ändern' }).click();
  68  | 
  69  |     // Check success message
  70  |     await expect(page.locator('#message')).toHaveText('Item erfolgreich aktualisiert');
  71  | 
  72  |     // Check item was updated
  73  |     await expect(page.locator('#list')).toContainText(updatedName);
  74  |     await expect(page.locator('#list')).not.toContainText(originalName);
  75  |   });
  76  | 
  77  |   test('should delete an item', async ({ page }) => {
  78  |     const itemName = 'Item to Delete ' + Date.now();
  79  | 
  80  |     // Create an item first
  81  |     await page.locator('#nameInput').fill(itemName);
  82  |     await page.locator('button').filter({ hasText: 'Erstellen' }).click();
  83  |     await expect(page.locator('#message')).toHaveText('Item erfolgreich erstellt');
  84  | 
  85  |     // Find the item and delete it
  86  |     const itemRow = page.locator('#list tr').filter({ hasText: itemName });
  87  |     await expect(itemRow).toBeVisible();
  88  | 
  89  |     // Click delete button
  90  |     await itemRow.locator('button').filter({ hasText: 'Löschen' }).click();
  91  | 
  92  |     // Check success message
  93  |     await expect(page.locator('#message')).toHaveText('Item erfolgreich gelöscht');
  94  | 
  95  |     // Check item was removed
  96  |     await expect(page.locator('#list')).not.toContainText(itemName);
  97  |   });
  98  | 
  99  |   test('should handle update of empty name', async ({ page }) => {
  100 |     const originalName = 'Item to Update ' + Date.now();
  101 | 
  102 |     // Create an item first
  103 |     await page.locator('#nameInput').fill(originalName);
  104 |     await page.locator('button').filter({ hasText: 'Erstellen' }).click();
  105 |     await expect(page.locator('#message')).toHaveText('Item erfolgreich erstellt');
  106 | 
  107 |     // Find the item and try to update to empty
  108 |     const itemRow = page.locator('#list tr').filter({ hasText: originalName });
  109 |     await itemRow.locator('input').clear();
  110 | 
  111 |     // Click update button
  112 |     await itemRow.locator('button').filter({ hasText: 'Ändern' }).click();
  113 | 
  114 |     // Check error message
  115 |     await expect(page.locator('#message')).toHaveText('Fehler: item name cannot be empty');
  116 | 
  117 |     // Check item still has original name
  118 |     await expect(page.locator('#list')).toContainText(originalName);
  119 |   });
  120 | });
```