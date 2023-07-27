package gostopwatch

import (
	"errors"
	"log"
	"time"
)

// NewGStopwatch creates and returns a GStopwatch
func NewGStopwatch(d time.Duration) (*GStopwatch, error) {
	ticks := int(d.Seconds())
	if ticks < minStopwatchTicks || ticks > maxStopwatchTicks {
		return nil, errors.New("invalid duration")
	}
	pt := &GStopwatch{
		t:              nil,
		state:          stopwatchStatusStopped,
		Done:           make(chan struct{}),
		monitorStopSig: make(chan struct{}),
		Tick:           make(chan time.Duration),
		d:              d,
		ticksLeft:      ticks,
	}
	return pt, nil
}

func (sw *GStopwatch) getState() StopwatchState {
	sw.stateRwMutex.Lock()
	defer sw.stateRwMutex.Unlock()
	return sw.state
}

func (sw *GStopwatch) setstate(s StopwatchState) {
	sw.stateRwMutex.Lock()
	defer sw.stateRwMutex.Unlock()
	sw.state = s
}

func (sw *GStopwatch) GetState() string {
	sw.stateRwMutex.Lock()
	defer sw.stateRwMutex.Unlock()
	return sw.state.String()
}

func (sw *GStopwatch) Start() error {
	if sw.getState() != stopwatchStatusStopped {
		return errors.New("ticker is not in stopped state")
	}
	sw.setstate(stopwatchStatusRunning)
	sw.t = time.NewTicker(tickerFrequency)
	go sw.monitorProgress()
	return nil
}

func (sw *GStopwatch) destroy() {
	sw.t.Stop()
}

func (sw *GStopwatch) monitorProgress() {
	for {
		if sw.getState() == stopwatchStatusPaused {
			continue
		}
		select {
		case <-sw.t.C:
			sw.ticksLeft -= 1
			sw.Tick <- time.Duration(sw.ticksLeft)
			log.Printf("ticks left %d", sw.ticksLeft)
			if sw.ticksLeft == 0 {
				// end the timer
				log.Printf("ended the monitor routine")
				sw.setstate(stopwatchStatusStopped)
				sw.Done <- struct{}{}
				sw.destroy()
				return
			}
		case <-sw.monitorStopSig:
			// end the timer
			sw.setstate(stopwatchStatusStopped)
			sw.Done <- struct{}{}
			sw.destroy()
			return
		}
	}
}

func (sw *GStopwatch) Stop() error {
	if sw.getState() != stopwatchStatusRunning {
		return errors.New("ticker is not in running state")
	}
	sw.setstate(stopwatchStatusStopped)
	sw.monitorStopSig <- struct{}{}
	return nil
}

func (sw *GStopwatch) Pause() error {
	if sw.getState() != stopwatchStatusRunning {
		return errors.New("ticker is not in running state")
	}
	sw.setstate(stopwatchStatusPaused)
	return nil
}

func (sw *GStopwatch) Resume() error {
	if sw.getState() != stopwatchStatusPaused {
		return errors.New("ticker is not in paused state")
	}
	sw.setstate(stopwatchStatusRunning)
	return nil
}
