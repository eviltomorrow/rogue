package procutil

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitForSigterm() os.Signal {
	var ch = make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTSTP, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		sig := <-ch
		switch sig {
		case os.Interrupt, syscall.SIGTSTP, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			return sig
		default:
		}
	}
}
