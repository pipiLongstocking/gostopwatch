package gostopwatch

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// NewGStopwatch creates and returns a GStopwatch.
func NewGStopwatch(d time.Duration) (*GStopwatch, error) {
	ticks := int(d.Seconds())
	if ticks < minStopwatchTicks || ticks > maxStopwatchTicks {
		return nil, errors.New("invalid duration")
	}
	pt := &GStopwatch{
		t:              nil,
		state:          stopwatchStatusStopped,
		done:           make(chan struct{}),
		monitorStopSig: make(chan struct{}),
		interrupt:      make(chan os.Signal, 1),
		tick:           make(chan time.Duration),
		d:              d,
		ticksLeft:      ticks,
	}
	signal.Notify(pt.interrupt, syscall.SIGINT, syscall.SIGTERM)
	return pt, nil
}

// getState returns the current state of the GStopwatch in the enum form.
// only for internal use.
func (sw *GStopwatch) getState() StopwatchState {
	sw.stateLock.RLock()
	defer sw.stateLock.RUnlock()
	return sw.state
}

func (sw *GStopwatch) setstate(s StopwatchState) {
	sw.stateLock.Lock()
	defer sw.stateLock.Unlock()
	sw.state = s
}

// GetState returns the current state of the GStopwatch in string form.
// The possible values are "stopped", "running" and "paused".
func (sw *GStopwatch) GetState() string {
	sw.stateLock.RLock()
	defer sw.stateLock.RUnlock()
	return sw.state.String()
}

// Start starts the GStopwatch. Errors if the GStopwatch is not in stopped state.
// The tick channel, done channel, and the error are returned.
// The tick channel returns the seconds left after every tick.
// The done channel returns an empty struct signal once the timer ends. Use this to stop the timer monitoring.
// Use GetState to check the state first before issuing the command.
func (sw *GStopwatch) Start() (<-chan time.Duration, <-chan struct{}, error) {
	sw.watchOp.Lock()
	defer sw.watchOp.Unlock()
	if sw.getState() != stopwatchStatusStopped {
		return nil, nil, errors.New("ticker is not in stopped state")
	}
	sw.setstate(stopwatchStatusRunning)
	sw.t = time.NewTicker(tickerFrequency)
	go sw.monitorProgress()
	return sw.tick, sw.done, nil
}

func (sw *GStopwatch) destroy() {
	sw.t.Stop()
	sw.done <- struct{}{}
	close(sw.done)
	close(sw.tick)
	close(sw.interrupt)
}

// GetTimeLeft returns the number of seconds left in the timer.
func (sw *GStopwatch) GetTimeLeft() time.Duration {
	sw.tlLock.RLock()
	defer sw.tlLock.RUnlock()
	return time.Duration(sw.ticksLeft) * time.Second
}

func (sw *GStopwatch) monitorProgress() {
	for {
		select {
		case <-sw.t.C:
			if sw.getState() == stopwatchStatusPaused {
				continue
			}
			sw.tlLock.Lock()
			sw.ticksLeft -= 1
			sw.tlLock.Unlock()
			sw.tick <- time.Duration(sw.ticksLeft) * time.Second
			if sw.ticksLeft == 0 {
				// end the timer
				sw.setstate(stopwatchStatusStopped)
				go func() { sw.destroy() }()
				return
			}
		case <-sw.monitorStopSig:
			// end the timer
			sw.setstate(stopwatchStatusStopped)
			go func() { sw.destroy() }()
			return
		case <-sw.interrupt:
			// Interrupt received
			sw.setstate(stopwatchStatusStopped)
			sw.destroy()
			return
		}
	}
}

// Stop stops the GStopwatch. Errors if the GStopwatch is not in running state.
// Use GetState to check the state first before issuing the command.
func (sw *GStopwatch) Stop() error {
	sw.watchOp.Lock()
	defer sw.watchOp.Unlock()
	if sw.getState() != stopwatchStatusRunning {
		return errors.New("ticker is not in running state")
	}
	sw.setstate(stopwatchStatusStopped)
	sw.monitorStopSig <- struct{}{}
	return nil
}

// Pause pauses the GStopwatch. Errors if the GStopwatch is not in running state.
// Use GetState to check the state first before issuing the command.
func (sw *GStopwatch) Pause() error {
	sw.watchOp.Lock()
	defer sw.watchOp.Unlock()
	if sw.getState() != stopwatchStatusRunning {
		return errors.New("ticker is not in running state")
	}
	sw.setstate(stopwatchStatusPaused)
	return nil
}

// Resume resumes the GStopwatch. Errors if the GStopwatch is not in paused state.
// Use GetState to check the state first before issuing the command.
func (sw *GStopwatch) Resume() error {
	sw.watchOp.Lock()
	defer sw.watchOp.Unlock()
	if sw.getState() != stopwatchStatusPaused {
		return errors.New("ticker is not in paused state")
	}
	sw.setstate(stopwatchStatusRunning)
	return nil
}
