package main

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"
)

// Windows API constants
const (
	WM_USER      = 0x0400
	WM_TRAYICON  = WM_USER + 1
	WM_KEYDOWN   = 0x0100
	WM_KEYUP     = 0x0101
	WM_COMMAND   = 0x0111
	WM_DESTROY   = 0x0002
	WM_RBUTTONUP = 0x0205

	ID_TOGGLE  = 1001
	ID_EXIT    = 1002
	TOGGLE_KEY = 0x79 // F10

	WH_KEYBOARD_LL = 13

	INPUT_KEYBOARD    = 1
	KEYEVENTF_KEYUP   = 0x0002
	KEYEVENTF_UNICODE = 0x0004

	VK_BACK    = 0x08
	VK_TAB     = 0x09
	VK_RETURN  = 0x0D
	VK_SHIFT   = 0x10
	VK_CONTROL = 0x11
	VK_SPACE   = 0x20
	VK_F10     = 0x79

	NIF_MESSAGE = 0x00000001
	NIF_ICON    = 0x00000002
	NIF_TIP     = 0x00000004
	NIM_ADD     = 0x00000000
	NIM_MODIFY  = 0x00000001
	NIM_DELETE  = 0x00000002

	IDI_APPLICATION = 32512
	IDC_ARROW       = 32512
	MF_STRING       = 0x00000000
	MF_SEPARATOR    = 0x00000800
	TPM_RIGHTBUTTON = 0x0002
)

// Windows API structures
type POINT struct {
	X, Y int32
}

type MSG struct {
	Hwnd    syscall.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

type WNDCLASSW struct {
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     syscall.Handle
	HIcon         syscall.Handle
	HCursor       syscall.Handle
	HbrBackground syscall.Handle
	LpszMenuName  *uint16
	LpszClassName *uint16
}

type KBDLLHOOKSTRUCT struct {
	VkCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

type KEYBDINPUT struct {
	Wvk         uint16
	Wscan       uint16
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

type INPUT struct {
	Type uint32
	Ki   KEYBDINPUT
	_    [8]byte // padding for union
}

type NOTIFYICONDATAW struct {
	CbSize           uint32
	Hwnd             syscall.Handle
	UID              uint32
	UFlags           uint32
	UCallbackMessage uint32
	HIcon            syscall.Handle
	SzTip            [128]uint16
	DwState          uint32
	DwStateMask      uint32
	SzInfo           [256]uint16
	UVersion         uint32
	SzInfoTitle      [64]uint16
	DwInfoFlags      uint32
}

// Global variables
var (
	keyboardState = &KeyboardState{
		enabled:           false,
		inputBuffer:       "",
		lastBengaliOutput: "",
		mutex:             sync.Mutex{},
	}

	bengaliKeyboard  = NewBengaliKeyboard()
	mainWindowHandle atomic.Value

	// Windows API DLLs
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")

	// Windows API functions
	registerClassW      = user32.NewProc("RegisterClassW")
	createWindowExW     = user32.NewProc("CreateWindowExW")
	defWindowProcW      = user32.NewProc("DefWindowProcW")
	getMessageW         = user32.NewProc("GetMessageW")
	translateMessage    = user32.NewProc("TranslateMessage")
	dispatchMessageW    = user32.NewProc("DispatchMessageW")
	postQuitMessage     = user32.NewProc("PostQuitMessage")
	setWindowsHookExW   = user32.NewProc("SetWindowsHookExW")
	callNextHookEx      = user32.NewProc("CallNextHookEx")
	unhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	getAsyncKeyState    = user32.NewProc("GetAsyncKeyState")
	sendInput           = user32.NewProc("SendInput")
	loadIconW           = user32.NewProc("LoadIconW")
	loadCursorW         = user32.NewProc("LoadCursorW")
	createPopupMenu     = user32.NewProc("CreatePopupMenu")
	appendMenuW         = user32.NewProc("AppendMenuW")
	getCursorPos        = user32.NewProc("GetCursorPos")
	setForegroundWindow = user32.NewProc("SetForegroundWindow")
	trackPopupMenu      = user32.NewProc("TrackPopupMenu")
	destroyMenu         = user32.NewProc("DestroyMenu")
	getModuleHandleW    = kernel32.NewProc("GetModuleHandleW")
	shellNotifyIconW    = shell32.NewProc("Shell_NotifyIconW")
)

type KeyboardState struct {
	enabled           bool
	inputBuffer       string
	lastBengaliOutput string
	mutex             sync.Mutex
}

type BengaliKeyboard struct {
	keymap *KeyMap
}

func NewBengaliKeyboard() *BengaliKeyboard {
	return &BengaliKeyboard{
		keymap: NewKeyMap(),
	}
}

func (bk *BengaliKeyboard) ConvertText(input string) string {
	var result strings.Builder
	chars := []rune(input)
	i := 0

	for i < len(chars) {
		longestMatch := ""
		longestBengali := ""
		longestLen := 0
		isVowel := false

		// Try to find the longest matching pattern
		for pattern, bengaliChar := range bk.keymap.Patterns {
			patternRunes := []rune(pattern)
			if i+len(patternRunes) <= len(chars) {
				slice := string(chars[i : i+len(patternRunes)])
				if slice == pattern && len(patternRunes) > longestLen {
					longestMatch = pattern
					longestBengali = bengaliChar.Bengali
					longestLen = len(patternRunes)
					isVowel = bengaliChar.IsVowel
				}
			}
		}

		if longestLen > 0 {
			// Special handling for vowels
			if isVowel {
				resultStr := result.String()
				if len(resultStr) > 0 && endsWithConsonant(resultStr) {
					if longestMatch == "o" {
						// "o" after consonant is inherent vowel - add nothing
					} else if diacritic, exists := bk.keymap.VowelDiacritics[longestMatch]; exists {
						result.WriteString(diacritic)
					} else {
						result.WriteString(longestBengali)
					}
				} else {
					// Independent vowel
					result.WriteString(longestBengali)
				}
			} else {
				// Consonant or other character
				result.WriteString(longestBengali)
			}
			i += longestLen
		} else {
			result.WriteRune(chars[i])
			i++
		}
	}

	return result.String()
}

func endsWithConsonant(text string) bool {
	if len(text) == 0 {
		return false
	}
	runes := []rune(text)
	lastChar := runes[len(runes)-1]
	return isBengaliConsonant(lastChar)
}

func isBengaliConsonant(ch rune) bool {
	return (ch >= '\u0995' && ch <= '\u09B9') || // ক to হ
		ch == '\u09DC' || ch == '\u09DD' || // ড় and ঢ়
		ch == '\u09DF' || // য়
		ch == '\u09CE' // ৎ
}

func main() {
	fmt.Println("Bengali Keyboard starting...")

	hInstance, _, _ := getModuleHandleW.Call(0)

	className := stringToUTF16("BengaliKeyboardClass")
	wc := WNDCLASSW{
		LpfnWndProc:   syscall.NewCallback(windowProc),
		HInstance:     syscall.Handle(hInstance),
		HCursor:       syscall.Handle(loadCursor()),
		LpszClassName: &className[0],
	}

	registerClassW.Call(uintptr(unsafe.Pointer(&wc)))

	windowName := stringToUTF16("Bengali Keyboard")
	hwnd, _, _ := createWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(&className[0])),
		uintptr(unsafe.Pointer(&windowName[0])),
		0,
		0, 0, 0, 0,
		0, 0,
		hInstance,
		0,
	)

	mainWindowHandle.Store(syscall.Handle(hwnd))

	hook, _, _ := setWindowsHookExW.Call(
		WH_KEYBOARD_LL,
		syscall.NewCallback(keyboardHookProc),
		hInstance,
		0,
	)

	if hook == 0 {
		fmt.Println("Failed to install keyboard hook")
		return
	}

	fmt.Println("Creating tray icon...")
	createTrayIcon(syscall.Handle(hwnd))
	fmt.Println("Tray icon created. Application running...")

	var msg MSG
	for {
		ret, _, _ := getMessageW.Call(
			uintptr(unsafe.Pointer(&msg)),
			0, 0, 0,
		)
		if ret == 0 || ret == ^uintptr(0) { // 0 = WM_QUIT, -1 = error
			break
		}
		translateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		dispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}

	unhookWindowsHookEx.Call(hook)
}

func windowProc(hwnd syscall.Handle, msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case WM_TRAYICON:
		if uint32(lparam) == WM_RBUTTONUP {
			showContextMenu(hwnd)
		}
		return 0
	case WM_COMMAND:
		switch uint32(wparam) & 0xFFFF {
		case ID_TOGGLE:
			toggleKeyboard()
			updateTrayIcon(hwnd)
		case ID_EXIT:
			postQuitMessage.Call(0)
		}
		return 0
	case WM_DESTROY:
		removeTrayIcon(hwnd)
		postQuitMessage.Call(0)
		return 0
	}
	ret, _, _ := defWindowProcW.Call(uintptr(hwnd), uintptr(msg), wparam, lparam)
	return ret
}

func keyboardHookProc(code int32, wparam, lparam uintptr) uintptr {
	if code >= 0 {
		kbdStruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lparam))
		vkCode := kbdStruct.VkCode

		// Check for toggle key (F10)
		if wparam == WM_KEYDOWN && vkCode == TOGGLE_KEY {
			toggleKeyboard()
			if hwnd := mainWindowHandle.Load(); hwnd != nil {
				updateTrayIcon(hwnd.(syscall.Handle))
			}
			return 1
		}

		// Check for Ctrl key combinations
		ctrlPressed := isKeyPressed(VK_CONTROL)
		if ctrlPressed && wparam == WM_KEYDOWN {
			switch vkCode {
			case 0x43, 0x56, 0x41, 0x58, 0x5A, 0x59: // Ctrl+C,V,A,X,Z,Y
				ret, _, _ := callNextHookEx.Call(0, uintptr(code), wparam, lparam)
				return ret
			}
		}

		keyboardState.mutex.Lock()
		enabled := keyboardState.enabled
		keyboardState.mutex.Unlock()

		if enabled && wparam == WM_KEYDOWN {
			if ch := vkToChar(vkCode); ch != 0 {
				if processCharacter(ch) {
					return 1
				}
			}
		}
	}

	ret, _, _ := callNextHookEx.Call(0, uintptr(code), wparam, lparam)
	return ret
}

func processCharacter(ch rune) bool {
	keyboardState.mutex.Lock()
	defer keyboardState.mutex.Unlock()

	if ch == '\b' { // Backspace
		if len(keyboardState.inputBuffer) > 0 {
			runes := []rune(keyboardState.inputBuffer)
			keyboardState.inputBuffer = string(runes[:len(runes)-1])
		}
		return false
	} else if ch == ' ' || ch == '\n' || ch == '\t' {
		// Word boundary - process current word
		if len(keyboardState.inputBuffer) > 0 {
			word := keyboardState.inputBuffer
			bengaliWord := bengaliKeyboard.ConvertText(word)

			// Clear the buffer
			keyboardState.inputBuffer = ""
			keyboardState.lastBengaliOutput = ""

			// If we have a valid Bengali conversion and it's different from input
			if len(bengaliWord) > 0 && bengaliWord != word {
				keyboardState.mutex.Unlock()

				// Remove the English word
				for i := 0; i < len([]rune(word)); i++ {
					sendBackspace()
				}

				// Send the Bengali word
				sendUnicodeText(bengaliWord)

				// Send the space/newline/tab that triggered the conversion
				sendCharacter(ch)

				keyboardState.mutex.Lock()
				return true // Suppress the original space/newline/tab
			}
		}

		// Clear buffer even if no conversion happened
		keyboardState.inputBuffer = ""
		keyboardState.lastBengaliOutput = ""
		return false // Allow the space/newline/tab
	} else if isValidInputChar(ch) {
		// Add character to buffer but don't convert yet
		keyboardState.inputBuffer += string(ch)
		return false // Allow the character to be typed normally
	} else {
		// Non-matching character, clear buffer
		if len(keyboardState.inputBuffer) > 0 {
			keyboardState.inputBuffer = ""
			keyboardState.lastBengaliOutput = ""
		}
		return false // Allow the character
	}
}

func isValidInputChar(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') || ch == '.' || ch == ':' || ch == '$' || ch == '_'
}

func sendBackspace() {
	input := INPUT{
		Type: INPUT_KEYBOARD,
		Ki: KEYBDINPUT{
			Wvk: VK_BACK,
		},
	}
	sendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))

	input.Ki.DwFlags = KEYEVENTF_KEYUP
	sendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
}

func sendUnicodeText(text string) {
	for _, ch := range text {
		input := INPUT{
			Type: INPUT_KEYBOARD,
			Ki: KEYBDINPUT{
				Wscan:   uint16(ch),
				DwFlags: KEYEVENTF_UNICODE,
			},
		}
		sendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))

		input.Ki.DwFlags = KEYEVENTF_UNICODE | KEYEVENTF_KEYUP
		sendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
	}
}

func sendCharacter(ch rune) {
	input := INPUT{
		Type: INPUT_KEYBOARD,
	}

	switch ch {
	case ' ':
		input.Ki.Wvk = VK_SPACE
	case '\n':
		input.Ki.Wvk = VK_RETURN
	case '\t':
		input.Ki.Wvk = VK_TAB
	default:
		input.Ki.Wscan = uint16(ch)
		input.Ki.DwFlags = KEYEVENTF_UNICODE
	}

	sendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))

	input.Ki.DwFlags |= KEYEVENTF_KEYUP
	sendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
}

func createTrayIcon(hwnd syscall.Handle) {
	var nid NOTIFYICONDATAW
	nid.CbSize = uint32(unsafe.Sizeof(nid))
	nid.Hwnd = hwnd
	nid.UID = 1
	nid.UFlags = NIF_ICON | NIF_MESSAGE | NIF_TIP
	nid.UCallbackMessage = WM_TRAYICON

	keyboardState.mutex.Lock()
	enabled := keyboardState.enabled
	keyboardState.mutex.Unlock()

	icon, _, _ := loadIconW.Call(0, IDI_APPLICATION)
	nid.HIcon = syscall.Handle(icon)

	var tooltip string
	if enabled {
		tooltip = "Bengali Keyboard - Enabled (F10 to toggle)"
	} else {
		tooltip = "Bengali Keyboard - Disabled (F10 to toggle)"
	}

	tooltipUTF16 := stringToUTF16(tooltip)
	copy(nid.SzTip[:], tooltipUTF16[:min(len(tooltipUTF16), 127)])

	shellNotifyIconW.Call(NIM_ADD, uintptr(unsafe.Pointer(&nid)))
}

func updateTrayIcon(hwnd syscall.Handle) {
	var nid NOTIFYICONDATAW
	nid.CbSize = uint32(unsafe.Sizeof(nid))
	nid.Hwnd = hwnd
	nid.UID = 1
	nid.UFlags = NIF_ICON | NIF_TIP

	keyboardState.mutex.Lock()
	enabled := keyboardState.enabled
	keyboardState.mutex.Unlock()

	icon, _, _ := loadIconW.Call(0, IDI_APPLICATION)
	nid.HIcon = syscall.Handle(icon)

	var tooltip string
	if enabled {
		tooltip = "Bengali Keyboard - Enabled (F10 to toggle)"
	} else {
		tooltip = "Bengali Keyboard - Disabled (F10 to toggle)"
	}

	tooltipUTF16 := stringToUTF16(tooltip)
	copy(nid.SzTip[:], tooltipUTF16[:min(len(tooltipUTF16), 127)])

	shellNotifyIconW.Call(NIM_MODIFY, uintptr(unsafe.Pointer(&nid)))
}

func removeTrayIcon(hwnd syscall.Handle) {
	var nid NOTIFYICONDATAW
	nid.CbSize = uint32(unsafe.Sizeof(nid))
	nid.Hwnd = hwnd
	nid.UID = 1

	shellNotifyIconW.Call(NIM_DELETE, uintptr(unsafe.Pointer(&nid)))
}

func showContextMenu(hwnd syscall.Handle) {
	hmenu, _, _ := createPopupMenu.Call()

	keyboardState.mutex.Lock()
	enabled := keyboardState.enabled
	keyboardState.mutex.Unlock()

	var toggleText string
	if enabled {
		toggleText = "Disable Bengali Keyboard"
	} else {
		toggleText = "Enable Bengali Keyboard"
	}

	toggleTextUTF16 := stringToUTF16(toggleText)
	exitTextUTF16 := stringToUTF16("Exit")

	appendMenuW.Call(hmenu, MF_STRING, ID_TOGGLE, uintptr(unsafe.Pointer(&toggleTextUTF16[0])))
	appendMenuW.Call(hmenu, MF_SEPARATOR, 0, 0)
	appendMenuW.Call(hmenu, MF_STRING, ID_EXIT, uintptr(unsafe.Pointer(&exitTextUTF16[0])))

	var pt POINT
	getCursorPos.Call(uintptr(unsafe.Pointer(&pt)))

	setForegroundWindow.Call(uintptr(hwnd))
	trackPopupMenu.Call(
		hmenu,
		TPM_RIGHTBUTTON,
		uintptr(pt.X),
		uintptr(pt.Y),
		0,
		uintptr(hwnd),
		0,
	)

	destroyMenu.Call(hmenu)
}

func toggleKeyboard() {
	keyboardState.mutex.Lock()
	defer keyboardState.mutex.Unlock()
	keyboardState.enabled = !keyboardState.enabled
	keyboardState.inputBuffer = ""
	keyboardState.lastBengaliOutput = ""
}

func vkToChar(vkCode uint32) rune {
	shiftPressed := isKeyPressed(VK_SHIFT)

	switch {
	case vkCode >= 0x41 && vkCode <= 0x5A: // A-Z
		if shiftPressed {
			return rune(vkCode - 0x41 + 'A')
		}
		return rune(vkCode - 0x41 + 'a')
	case vkCode >= 0x30 && vkCode <= 0x39: // 0-9
		return rune(vkCode - 0x30 + '0')
	case vkCode == 0x08:
		return '\b'
	case vkCode == 0x20:
		return ' '
	case vkCode == 0x0D:
		return '\n'
	case vkCode == 0x09:
		return '\t'
	case vkCode == 0xBE:
		return '.'
	case vkCode == 0xBA:
		return ':'
	case vkCode == 0xBD:
		return '_'
	}
	return 0
}

func isKeyPressed(vk uint32) bool {
	ret, _, _ := getAsyncKeyState.Call(uintptr(vk))
	return (ret & 0x8000) != 0
}

func loadCursor() uintptr {
	ret, _, _ := loadCursorW.Call(0, IDC_ARROW)
	return ret
}

func stringToUTF16(s string) []uint16 {
	return syscall.StringToUTF16(s)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
