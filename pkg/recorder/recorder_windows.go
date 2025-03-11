//go:build windows
// +build windows

package recorder

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"sync"
	"syscall"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
)

// Windows API constants
const (
	VK_CONTROL = 0x11
	VK_MENU    = 0x12 // ALT key
	VK_R       = 0x52
	VK_LBUTTON = 0x01
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	getAsyncKeyState = user32.NewProc("GetAsyncKeyState")
)

type Step struct {
	Screenshot  image.Image
	Description string
	Timestamp   time.Time
	Action      string
	Coordinates image.Point
	Highlighted bool
}

type Recorder struct {
	steps        []Step
	isRecording  bool
	mu           sync.Mutex
	lastX, lastY int
	stopChan     chan struct{}
}

func NewRecorder() *Recorder {
	return &Recorder{
		steps:    make([]Step, 0),
		stopChan: make(chan struct{}),
	}
}

func (r *Recorder) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isRecording {
		return fmt.Errorf("recording is already in progress")
	}

	r.stopChan = make(chan struct{})
	r.steps = make([]Step, 0)
	r.isRecording = true

	go r.monitorEvents()
	return nil
}

func (r *Recorder) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.isRecording {
		return fmt.Errorf("no recording in progress")
	}

	r.isRecording = false
	if r.stopChan != nil {
		close(r.stopChan)
		r.stopChan = nil
	}
	log.Printf("Recording stopped. Captured %d steps", len(r.steps))
	return nil
}

func (r *Recorder) GetSteps() []Step {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.steps
}

func (r *Recorder) monitorEvents() {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	var lastMouseState bool

	for {
		select {
		case <-r.stopChan:
			return
		case <-ticker.C:
			r.mu.Lock()
			if (!r.isRecording) {
				r.mu.Unlock()
				continue
			}

			mouseState, _, _ := getAsyncKeyState.Call(uintptr(VK_LBUTTON))
			isMouseDown := mouseState&0x8000 != 0

			if isMouseDown && !lastMouseState {
				x, y := robotgo.GetMousePos()
				n := screenshot.NumActiveDisplays()

				var targetBounds image.Rectangle
				var displayIndex int = -1

				for i := 0; i < n; i++ {
					bounds := screenshot.GetDisplayBounds(i)
					if x >= bounds.Min.X && x < bounds.Max.X && y >= bounds.Min.Y && y < bounds.Max.Y {
						targetBounds = bounds
						displayIndex = i
						break
					}
				}

				if displayIndex == -1 {
					log.Printf("Click position (%d, %d) not found in any display", x, y)
					r.mu.Unlock()
					continue
				}

				img, err := screenshot.CaptureRect(targetBounds)
				if err != nil {
					log.Printf("Failed to capture screenshot: %v", err)
					r.mu.Unlock()
					continue
				}

				localX := x - targetBounds.Min.X
				localY := y - targetBounds.Min.Y

				highlightedImg := r.addHighlightCircle(img, localX, localY)

				r.steps = append(r.steps, Step{
					Screenshot:  highlightedImg,
					Timestamp:   time.Now(),
					Action:      "Mouse Click",
					Coordinates: image.Point{X: x, Y: y},
					Highlighted: true,
				})
			}
			lastMouseState = isMouseDown
			r.mu.Unlock()
		}
	}
}

// addHighlightCircle adds a red circle around the click point using Bresenham's circle algorithm
func (r *Recorder) addHighlightCircle(img image.Image, x, y int) image.Image {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)

	for px := bounds.Min.X; px < bounds.Max.X; px++ {
		for py := bounds.Min.Y; py < bounds.Max.Y; py++ {
			rgba.Set(px, py, img.At(px, py))
		}
	}

	radius := 20
	red := color.RGBA{R: 255, A: 255}

	for angle := 0; angle < 360; angle++ {
		radian := float64(angle) * math.Pi / 180
		px := int(float64(x) + float64(radius)*math.Cos(radian))
		py := int(float64(y) + float64(radius)*math.Sin(radian))

		for i := -2; i <= 2; i++ {
			for j := -2; j <= 2; j++ {
				if px+i >= 0 && px+i < bounds.Max.X && py+j >= 0 && py+j < bounds.Max.Y {
					rgba.Set(px+i, py+j, red)
				}
			}
		}
	}

	return rgba
}
