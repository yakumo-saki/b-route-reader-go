package watchdog

import "time"

type WatchdogTimer struct {
	timelimit time.Duration
	timer     *time.Timer
	handler   func()
}

func NewWatchdogTimer(timelimit time.Duration, handler func()) *WatchdogTimer {
	wdt := WatchdogTimer{
		timelimit: timelimit,
		handler:   handler,
	}
	wdt.timer = time.AfterFunc(timelimit, wdt.handleFire)
	return &wdt
}

func (wdt *WatchdogTimer) handleFire() {
	wdt.handler()
}
