package gostopwatch

import (
	"os"
	"sync"
	"time"
)

// StopwatchState is the enum for the possible ticker states.
type StopwatchState int

func (ts StopwatchState) String() string {
	switch ts {
	case stopwatchStatusStopped:
		return "stopped"
	case stopwatchStatusRunning:
		return "running"
	default:
		return "paused"
	}
}

// Constants for the current state of a ticker.
const (
	stopwatchStatusStopped StopwatchState = iota
	stopwatchStatusRunning
	stopwatchStatusPaused
)

// ticker Limits
const (
	maxStopwatchTicks = 86400
	minStopwatchTicks = 1
	tickerFrequency   = time.Second
)

type GStopwatch struct {
	t         *time.Ticker   // An internal ticker for issuing periodic ticks.
	state     StopwatchState // state denotes the current state of GStopwatch.
	stateLock sync.RWMutex   // RWmutex for read/writer operations to state.

	done chan struct{} // done channel which issues the signal for the end of the ticker.
	// Poll on this blocking channel to check when the timer ends.

	monitorStopSig chan struct{}      // Channel to issue the signal to stop the monitoring routine.
	interrupt      chan os.Signal     // Channel for listening to interrupts.
	tick           chan time.Duration // Channel that returns the time left in seconds after every tick.
	// Poll on this blocking channel to check the time left in the GStopwatch.

	d         time.Duration // The duration for which the ticker should run.
	ticksLeft int           // Ticks left in the ticker.
	tlLock    sync.RWMutex  //RWMutex for issuing the GetTimeLeft() method.

	watchOp sync.Mutex // Mutex for issuing Resume(), Pause(), Start(), Stop() operations.
}
