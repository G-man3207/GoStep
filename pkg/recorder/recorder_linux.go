//go:build !windows
// +build !windows

package recorder

import (
	"errors"
	"image"
	"time"
)

type Step struct {
	Screenshot  image.Image
	Description string
	Timestamp   time.Time
	Action      string
	Coordinates image.Point
}

type Recorder struct {
	steps []Step
}

func NewRecorder() *Recorder {
	return &Recorder{
		steps: make([]Step, 0),
	}
}

func (r *Recorder) Start() error {
	return errors.New("recording is only supported on Windows")
}

func (r *Recorder) Stop() error {
	return errors.New("recording is only supported on Windows")
}

func (r *Recorder) GetSteps() []Step {
	return r.steps
}
