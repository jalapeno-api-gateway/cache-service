package helpers

import (
	"os"
	"os/signal"
)

func WatchInterruptSignals() chan os.Signal {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
    return signals
}