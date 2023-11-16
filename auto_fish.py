import time
import pyautogui


def mouse():
    while True:
        print(pyautogui.position())
        time.sleep(1)
# mouse()
def auto():
    while True:
        pyautogui.moveTo(1655, 246)
        pyautogui.moveTo(1951, 1355)
        pyautogui.click()
        pyautogui.write("/debug2 10376 824636140624 3")
        pyautogui.hotkey('enter')
        pyautogui.moveTo(1923, 1283)
        time.sleep(2)
        pyautogui.rightClick()
        pyautogui.moveTo(2008, 1111)
        pyautogui.click()
        pyautogui.moveTo(1914, 566)
        pyautogui.click()
        pyautogui.hotkey('enter')
        time.sleep(10)
        pyautogui.moveTo(1655, 246)
        pyautogui.click()
        time.sleep(8.2)

auto()