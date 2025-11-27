//go:build darwin

package system

/*
#cgo darwin LDFLAGS: -framework ApplicationServices
#include <ApplicationServices/ApplicationServices.h>

CGPoint GetMouseLocation() {
    CGEventRef event = CGEventCreate(NULL);
    CGPoint p = CGEventGetLocation(event);
    CFRelease(event);
    return p;
}
*/
import "C"

import (
	"backend/internal/config"
	"backend/internal/models"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	lastWindowCheck time.Time
	cachedWindows   []models.WindowInfo
	windowsCacheMu  sync.RWMutex
)

func GetCursorPosition() models.CursorInfo {
	p := C.GetMouseLocation()
	return models.CursorInfo{
		X: int(p.x),
		Y: int(p.y),
	}
}

func getChromeWindows() ([]models.WindowInfo, error) {
	script := `
tell application "Google Chrome"
	set out to ""
	repeat with w in windows
		set theTitle to title of active tab of w
		set {leftPos, topPos, rightPos, bottomPos} to bounds of w
		set out to out & theTitle & "|" & leftPos & "," & topPos & "," & rightPos & "," & bottomPos & linefeed
	end repeat
	return out
end tell
`

	cmd := exec.Command("osascript", "-e", script)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	var windows []models.WindowInfo

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			continue
		}
		title := parts[0]
		coordStr := parts[1]

		coords := strings.Split(coordStr, ",")
		if len(coords) != 4 {
			continue
		}

		left, err1 := strconv.Atoi(strings.TrimSpace(coords[0]))
		top, err2 := strconv.Atoi(strings.TrimSpace(coords[1]))
		right, err3 := strconv.Atoi(strings.TrimSpace(coords[2]))
		bottom, err4 := strconv.Atoi(strings.TrimSpace(coords[3]))
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			continue
		}

		w := right - left
		h := bottom - top

		if w <= 0 || h <= 0 {
			continue
		}

		windows = append(windows, models.WindowInfo{
			Title:  title,
			X:      left,
			Y:      top,
			Width:  w,
			Height: h,
		})
	}

	return windows, nil
}

func GetCachedChromeWindows() ([]models.WindowInfo, error) {
	windowsCacheMu.RLock()
	if time.Since(lastWindowCheck) < config.AppConfig.WindowCacheTTL && cachedWindows != nil {
		windows := cachedWindows
		windowsCacheMu.RUnlock()
		return windows, nil
	}
	windowsCacheMu.RUnlock()

	windows, err := getChromeWindows()
	if err != nil {
		return nil, err
	}

	windowsCacheMu.Lock()
	cachedWindows = windows
	lastWindowCheck = time.Now()
	windowsCacheMu.Unlock()
	return windows, nil
}
