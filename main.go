package main

import (
	"fmt"
	"syscall"
	"time"
	"math/rand"
)

var (
	user32                = syscall.NewLazyDLL("user32.dll")
	kernel32              = syscall.NewLazyDLL("kernel32.dll")
	procGetAsyncKeyState  = user32.NewProc("GetAsyncKeyState")
	procKeybdEvent        = user32.NewProc("keybd_event")
	procSetConsoleCtrlHandler = kernel32.NewProc("SetConsoleCtrlHandler")

	running = false
)
var (
	procSetConsoleCP  = kernel32.NewProc("SetConsoleCP")
	procSetConsoleOutputCP = kernel32.NewProc("SetConsoleOutputCP")
)

func setConsoleUTF8() {
	// Устанавливаем кодировку ввода и вывода в UTF-8 (кодировка 65001)
	procSetConsoleCP.Call(65001)
	procSetConsoleOutputCP.Call(65001)
}

// Константы клавиш
const (
	VK_F       = 0x46 // Клавиша "F"
	VK_MENU    = 0x12 // Alt
	VK_1       = 0x31
	VK_2       = 0x32
	KEYEVENTF_KEYUP = 0x02
)

func main() {
	setConsoleUTF8()
	fmt.Println("👊 Жми Alt+1, чтобы начать кликать F, и Alt+2, чтобы остановить! 🛠️")

	// Установка обработчика для выхода через Ctrl+C
	setConsoleCtrlHandler()

	// Основной цикл для проверки клавиш
	go func() {
		for {
			// Проверка Alt+1
			if isKeyPressed(VK_MENU) && isKeyPressed(VK_1) {
				fmt.Println("🔥 Начинаю жать F")
				running = true
			}

			// Проверка Alt+2
			if isKeyPressed(VK_MENU) && isKeyPressed(VK_2) {
				fmt.Println("🛑 Останавливаю!")
				running = false
			}

			time.Sleep(50 * time.Millisecond) // Небольшая задержка
		}
	}()

	// Горутина для кликов F
	go func() {
		for {
			if running {
				pressKey(VK_F)
				// Рандомная задержка от 42 до 78 мс
				delay := time.Duration(rand.Intn(37)+42) * time.Millisecond
				fmt.Printf("🎮 Кликнул F, задержка %d мс\n", delay.Milliseconds())
				time.Sleep(delay)
			}
		}
	}()

	// Блокируем основной поток, чтобы программа не завершилась
	select {}
}

// Проверка нажатия клавиши
func isKeyPressed(vkCode int) bool {
	ret, _, _ := procGetAsyncKeyState.Call(uintptr(vkCode))
	return ret&0x8000 != 0
}

// Нажатие клавиши
func pressKey(vkCode int) {
	procKeybdEvent.Call(uintptr(vkCode), 0, 0, 0)
	time.Sleep(10 * time.Millisecond)
	procKeybdEvent.Call(uintptr(vkCode), 0, KEYEVENTF_KEYUP, 0)
}

// Обработчик для выхода через Ctrl+C
func setConsoleCtrlHandler() {
	handler := syscall.NewCallback(func(ctrlType uint) int {
		if ctrlType == 0 {
			fmt.Println("🚪 Выход по Ctrl+C")
			return 1
		}
		return 0
	})
	procSetConsoleCtrlHandler.Call(handler, 1)
}
