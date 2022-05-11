package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/no-src/log"
)

func main() {
	stop := make(chan bool, 1)
	log.Info("no lock screen starting...")

	go notify(func() error {
		stop <- true
		return nil
	})

	active(stop)

	log.Info("no lock screen stopped!")
}

func active(stop chan bool) {
	for {
		select {
		case <-stop:
			return
		case <-time.After(time.Second * 10):
		}
		robotgo.MoveRelative(0, 0)
	}
}

func notify(shutdown func() error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGTERM)
	for {
		s := <-c
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGTERM:
			log.Debug("received a signal [%s], waiting to exit", s.String())
			err := shutdown()
			if err != nil {
				log.Error(err, "shutdown error")
			} else {
				signal.Stop(c)
				close(c)
				return
			}
		default:
			log.Debug("received a signal [%s], ignore it", s.String())
		}
	}
}
