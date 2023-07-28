package gostopwatch_test

import (
	"errors"
	"gostopwatch/gostopwatch"
	"testing"
	"time"
)

type TestCaseNewSW struct {
	input time.Duration
	want  error
}

func TestGStopwatch_NewGStopWatch(t *testing.T) {
	tcs := []TestCaseNewSW{
		{500 * time.Millisecond, errors.New("invalid duration")},
		{48 * time.Hour, errors.New("invalid duration")},
		{1 * time.Second, nil},
	}
	for _, tc := range tcs {
		_, got := gostopwatch.NewGStopwatch(tc.input)
		if got != nil && tc.want != nil {
			if got.Error() != tc.want.Error() {
				t.Errorf("NewGStopwatch(%v) = %v; want %v", tc.input, got, tc.want)
			}
		} else if got == nil && tc.want != nil || tc.want == nil && got != nil {
			t.Errorf("NewGStopwatch(%v) = %v; want %v", tc.input, got, tc.want)

		}
	}
}

func TestGStopwatch_Start(t *testing.T) {
	sw, _ := gostopwatch.NewGStopwatch(2 * time.Second)
	tick, d, err := sw.Start()
	if err != nil {
		t.Errorf("Start() = %v; want %v", err, nil)
	}
	tc := time.NewTicker(4 * time.Second)
	for {
		select {
		case <-tick:
			continue
		case <-tc.C:
			t.Errorf("StopWatch exceeded %v", 2*time.Second)
			tc.Stop()
			return
		case <-d:
			return
		}
	}

}

func TestGStopwatch_Stop(t *testing.T) {
	sw, _ := gostopwatch.NewGStopwatch(3 * time.Second)
	_, _, err := sw.Start()
	if err != nil {
		t.Errorf("Start() = %v; want %v", err, nil)
	}
	err = sw.Stop()
	if err != nil {
		t.Errorf("Stop() = %v; want %v", err, nil)
	}
}

func TestGStopwatch_PauseResume(t *testing.T) {
	sw, _ := gostopwatch.NewGStopwatch(3 * time.Second)
	tick, d, err := sw.Start()
	if err != nil {
		t.Errorf("Start() = %v; want %v", err, nil)
	}
	err = sw.Pause()
	if err != nil {
		t.Errorf("Pause() = %v; want %v", err, nil)
	}
	err = sw.Resume()
	if err != nil {
		t.Errorf("Resume() = %v; want %v", err, nil)
	}
	tc := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-tick:
			continue
		case <-tc.C:
			t.Errorf("StopWatch exceeded %v", 3*time.Second)
			tc.Stop()
			return
		case <-d:
			return
		}
	}
}
