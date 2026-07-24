import subprocess
import time
import unittest
import urllib.request
from pathlib import Path

from selenium import webdriver
from selenium.webdriver import ActionChains
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.select import Select
from selenium.webdriver.support.ui import WebDriverWait

BASE_URL = "http://127.0.0.1:8080/"
PROJECT_ROOT = Path(__file__).resolve().parents[2]


class CRUDSeleniumTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.server = cls._start_server()
        cls.driver = cls._create_driver()
        cls.wait = WebDriverWait(cls.driver, 10)
        cls.driver.get(BASE_URL)

    @classmethod
    def tearDownClass(cls):
        cls.driver.quit()
        if cls.server is not None:
            cls.server.terminate()
            try:
                cls.server.wait(timeout=10)
            except subprocess.TimeoutExpired:
                cls.server.kill()

    @classmethod
    def _start_server(cls):
        try:
            with urllib.request.urlopen(BASE_URL, timeout=2):
                return None
        except Exception:
            pass

        process = subprocess.Popen(
            ["go", "run", "main.go"],
            cwd=str(PROJECT_ROOT),
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
        )

        deadline = time.time() + 45
        while time.time() < deadline:
            try:
                with urllib.request.urlopen(BASE_URL, timeout=1):
                    return process
            except Exception:
                if process.poll() is not None:
                    output = process.stdout.read() if process.stdout else ""
                    raise RuntimeError(f"The Go application exited before Selenium tests could start. Output: {output}")
                time.sleep(0.5)

        process.terminate()
        raise RuntimeError("Timed out while starting the Go application for Selenium tests.")

    @classmethod
    def _create_driver(cls):
        options = Options()
        options.add_argument("--headless=new")
        options.add_argument("--no-sandbox")
        options.add_argument("--disable-dev-shm-usage")
        options.add_argument("--window-size=1920,1080")
        return webdriver.Chrome(options=options)

    def setUp(self):
        self.driver.delete_all_cookies()
        self.driver.get(BASE_URL)
        self.wait.until(EC.visibility_of_element_located((By.ID, "nameInput")))

    def _create_item(self, name, description="", status="active"):
        self.driver.find_element(By.ID, "nameInput").clear()
        self.driver.find_element(By.ID, "nameInput").send_keys(name)
        self.driver.find_element(By.ID, "descriptionInput").clear()
        self.driver.find_element(By.ID, "descriptionInput").send_keys(description)
        select = Select(self.driver.find_element(By.ID, "statusInput"))
        select.select_by_value(status)
        self.driver.find_element(By.XPATH, "//button[normalize-space()='Erstellen']").click()

    def _wait_for_message(self, expected_text, is_error=False):
        message = self.wait.until(EC.visibility_of_element_located((By.ID, "message")))
        self.wait.until(lambda _: message.text == expected_text)
        color = message.value_of_css_property("color")
        if is_error:
            self.assertIn("255, 0, 0", color)
        else:
            self.assertIn("0, 128, 0", color)
        return message

    def _find_item_row(self, item_name):
        return self.wait.until(
            lambda _: self.driver.find_element(
                By.XPATH,
                f"//table[@id='list']//tr[.//input[@value='{item_name}']]",
            )
        )

    def test_should_load_the_page(self):
        self.assertEqual(self.driver.title, "CRUD Demo")
        self.assertEqual(self.driver.find_element(By.TAG_NAME, "h1").text, "CRUD Demo")
        self.assertTrue(self.driver.find_element(By.ID, "nameInput").is_displayed())
        self.assertTrue(
            self.driver.find_element(By.XPATH, "//button[normalize-space()='Erstellen']").is_displayed()
        )

    def test_should_create_a_new_item(self):
        item_name = f"Selenium Item {int(time.time() * 1000)}"
        self._create_item(item_name)
        self._wait_for_message("Item erfolgreich erstellt")
        self.assertTrue(self._find_item_row(item_name).is_displayed())

    def test_should_show_error_for_empty_item_creation(self):
        before_count = len(self.driver.find_elements(By.CSS_SELECTOR, "#list tr"))
        self.driver.find_element(By.ID, "nameInput").clear()
        self.driver.find_element(By.XPATH, "//button[normalize-space()='Erstellen']").click()
        self._wait_for_message("Fehler: item name cannot be empty", is_error=True)
        after_count = len(self.driver.find_elements(By.CSS_SELECTOR, "#list tr"))
        self.assertEqual(before_count, after_count)

    def test_should_update_an_item(self):
        original_name = f"Original {int(time.time() * 1000)}"
        updated_name = f"Updated {int(time.time() * 1000 + 1)}"

        self._create_item(original_name)
        self._wait_for_message("Item erfolgreich erstellt")
        row = self._find_item_row(original_name)
        row.find_element(By.XPATH, ".//input[contains(@id,'name-')]" ).clear()
        row.find_element(By.XPATH, ".//input[contains(@id,'name-')]" ).send_keys(updated_name)
        row.find_element(By.XPATH, ".//button[normalize-space()='Ändern']").click()
        self._wait_for_message("Item erfolgreich aktualisiert")
        self.assertTrue(self._find_item_row(updated_name).is_displayed())
        self.assertEqual(len(self.driver.find_elements(By.XPATH, f"//input[@value='{original_name}']")), 0)

    def test_should_delete_an_item(self):
        item_name = f"Delete Me {int(time.time() * 1000)}"
        self._create_item(item_name)
        self._wait_for_message("Item erfolgreich erstellt")
        row = self._find_item_row(item_name)
        row.find_element(By.XPATH, ".//button[normalize-space()='Löschen']").click()
        self._wait_for_message("Item erfolgreich gelöscht")
        self.assertEqual(len(self.driver.find_elements(By.XPATH, f"//input[@value='{item_name}']")), 0)

    def test_should_handle_update_of_empty_name(self):
        item_name = f"Empty Update {int(time.time() * 1000)}"
        self._create_item(item_name)
        self._wait_for_message("Item erfolgreich erstellt")
        row = self._find_item_row(item_name)
        row.find_element(By.XPATH, ".//input[contains(@id,'name-')]" ).clear()
        row.find_element(By.XPATH, ".//button[normalize-space()='Ändern']").click()
        self._wait_for_message("Fehler: item name cannot be empty", is_error=True)
        self.assertTrue(self._find_item_row(item_name).is_displayed())

    def test_persistence_after_reload(self):
        item_name = f"Persist {int(time.time() * 1000)}"
        self._create_item(item_name)
        self._wait_for_message("Item erfolgreich erstellt")
        self.driver.refresh()
        self.wait.until(EC.visibility_of_element_located((By.ID, "nameInput")))
        self.assertTrue(self._find_item_row(item_name).is_displayed())

    def test_multiple_rapid_creates(self):
        names = [f"Bulk {int(time.time() * 1000 + index)}" for index in range(5)]
        for name in names:
            self._create_item(name)
            self._wait_for_message("Item erfolgreich erstellt")

        for name in names:
            self.assertTrue(self._find_item_row(name).is_displayed())

    def test_accessibility_smoke_keyboard_focus_order(self):
        self.driver.find_element(By.ID, "nameInput").click()
        ActionChains(self.driver).send_keys(Keys.TAB).perform()
        self.assertTrue(self.driver.switch_to.active_element == self.driver.find_element(By.ID, "descriptionInput"))
        ActionChains(self.driver).send_keys(Keys.TAB).perform()
        self.assertTrue(self.driver.switch_to.active_element == self.driver.find_element(By.ID, "statusInput"))
        ActionChains(self.driver).send_keys(Keys.TAB).perform()
        self.assertTrue(self.driver.switch_to.active_element == self.driver.find_element(By.XPATH, "//button[normalize-space()='Erstellen']"))


if __name__ == "__main__":
    unittest.main(verbosity=2)
