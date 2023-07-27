package gostopwatch

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
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
		interrupt:      make(chan os.Signal, 1),
		Tick:           make(chan time.Duration),
		d:              d,
		ticksLeft:      ticks,
	}
	signal.Notify(pt.interrupt, syscall.SIGINT, syscall.SIGTERM)
	return pt, nil
}

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

func (sw *GStopwatch) GetState() string {
	sw.stateLock.RLock()
	defer sw.stateLock.RUnlock()
	return sw.state.String()
}

func (sw *GStopwatch) Start() error {
	sw.watchOp.Lock()
	defer sw.watchOp.Unlock()
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
	sw.Done <- struct{}{}
	close(sw.Done)
	close(sw.Tick)
	close(sw.interrupt)
}

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
			sw.Tick <- time.Duration(sw.ticksLeft) * time.Second
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
			go func() { sw.destroy() }()
			return
		}
	}
}

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

func (sw *GStopwatch) Pause() error {
	sw.watchOp.Lock()
	defer sw.watchOp.Unlock()
	if sw.getState() != stopwatchStatusRunning {
		return errors.New("ticker is not in running state")
	}
	sw.setstate(stopwatchStatusPaused)
	return nil
}

func (sw *GStopwatch) Resume() error {
	sw.watchOp.Lock()
	defer sw.watchOp.Unlock()
	if sw.getState() != stopwatchStatusPaused {
		return errors.New("ticker is not in paused state")
	}
	sw.setstate(stopwatchStatusRunning)
	return nil
}
