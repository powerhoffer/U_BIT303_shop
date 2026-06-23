const fs = require('fs');
const path = require('path');

const { chromium } = require('C:/Users/92569/.cache/codex-runtimes/codex-primary-runtime/dependencies/node/node_modules/playwright');

const root = path.resolve(__dirname, '..');
const outDir = path.join(root, 'output', 'progress-brief', 'screenshots');
fs.mkdirSync(outDir, { recursive: true });

const baseUrl = process.env.FRONTEND_URL || 'http://127.0.0.1:8080';

async function waitAndShot(page, route, filename, selector) {
  await page.goto(`${baseUrl}/#${route}`, { waitUntil: 'networkidle' });
  if (selector) {
    await page.waitForSelector(selector, { timeout: 15000 });
  }
  await page.waitForTimeout(500);
  await page.screenshot({ path: path.join(outDir, filename), fullPage: true });
}

async function clickIfVisible(page, text) {
  const locator = page.getByText(text, { exact: true }).first();
  if (await locator.count()) {
    await locator.click();
    await page.waitForTimeout(500);
    return true;
  }
  return false;
}

(async () => {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage({ viewport: { width: 1440, height: 1000 }, deviceScaleFactor: 1 });

  await page.goto(`${baseUrl}/#/login`, { waitUntil: 'networkidle' });
  await page.screenshot({ path: path.join(outDir, '01-login.png'), fullPage: true });

  await page.fill('input[name="username"]', 'root');
  await page.fill('input[name="password"]', '123456');
  const checkbox = page.locator('.remember input[type="checkbox"]').first();
  if (await checkbox.count()) {
    await checkbox.check({ force: true });
  }
  await page.getByRole('button', { name: /login/i }).click();
  await page.waitForURL(/#\/dashboard|#\//, { timeout: 20000 }).catch(() => {});
  await page.waitForLoadState('networkidle');

  await waitAndShot(page, '/dashboard', '02-dashboard.png', '.dashboard-container');
  await waitAndShot(page, '/employee/list', '03-employee-list.png', '.page-container');
  await clickIfVisible(page, 'New Employee');
  await page.screenshot({ path: path.join(outDir, '04-employee-form.png'), fullPage: true });
  await page.keyboard.press('Escape').catch(() => {});

  await waitAndShot(page, '/points/my', '05-my-credits.png', '.page-container');
  await waitAndShot(page, '/points/manage', '06-credit-operations.png', '.page-container');
  await waitAndShot(page, '/category/list', '07-category-list.png', '.page-container');
  await waitAndShot(page, '/goods/list', '08-goods-list.png', '.page-container');

  await clickIfVisible(page, 'New Goods');
  await page.screenshot({ path: path.join(outDir, '09-goods-form.png'), fullPage: true });
  await page.keyboard.press('Escape').catch(() => {});
  await page.waitForTimeout(300);

  const detailButton = page.getByText('Detail', { exact: true }).first();
  if (await detailButton.count()) {
    await detailButton.click();
    await page.waitForTimeout(500);
    await page.screenshot({ path: path.join(outDir, '10-goods-detail.png'), fullPage: true });
  } else {
    await page.screenshot({ path: path.join(outDir, '10-goods-detail-unavailable.png'), fullPage: true });
  }

  await browser.close();
  console.log(outDir);
})();
