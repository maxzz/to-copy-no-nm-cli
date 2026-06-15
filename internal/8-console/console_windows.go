//go:build windows

package console

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	enableVirtualTerminalProcessing = 0x0004
)

var (
	kernel32           = windows.NewLazySystemDLL("kernel32.dll")
	procGetStdHandle   = kernel32.NewProc("GetStdHandle")
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
	procGetch          = windows.NewLazySystemDLL("msvcrt.dll").NewProc("_getch")
)

const stdOutputHandle = ^uintptr(10) // -11 as uintptr
const stdInputHandle = ^uintptr(9)   // -10 as uintptr

func init() {
	enableColors()
}

func enableColors() {
	handle, _, _ := procGetStdHandle.Call(stdOutputHandle)
	if handle == 0 {
		return
	}

	var mode uint32
	r1, _, _ := procGetConsoleMode.Call(handle, uintptr(unsafe.Pointer(&mode)))
	if r1 == 0 {
		return
	}

	procSetConsoleMode.Call(handle, uintptr(mode|enableVirtualTerminalProcessing))
}

const (
	colorRed    = "\x1b[31m"
	colorGreen  = "\x1b[32m"
	colorYellow = "\x1b[33m"
	colorDim    = "\x1b[2m\x1b[90m"
	colorReset  = "\x1b[0m"
)

// PrintError writes art and err in red, then waits for a key press before exiting.
func PrintError(err error, art string) {
	fmt.Print(colorRed)
	fmt.Print(art)
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	fmt.Print(colorReset)
	fmt.Print("Press any key to close...")
	waitForKey()
	os.Exit(1)
}

func PrintSuccess(art string) {
	fmt.Print(colorGreen)
	fmt.Print(art)
	fmt.Print(colorReset)
	time.Sleep(1500 * time.Millisecond)
	os.Exit(0)
}

func waitForKey() {
	handle, _, _ := procGetStdHandle.Call(stdInputHandle)
	if handle == 0 {
		waitForKeyFallback()
		return
	}

	var mode uint32
	r1, _, _ := procGetConsoleMode.Call(handle, uintptr(unsafe.Pointer(&mode)))
	if r1 == 0 {
		waitForKeyFallback()
		return
	}

	procGetch.Call()
}

func waitForKeyFallback() {
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadByte()
}
