package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/polipopoliko/go-block/bootstrap"
)

func Serve(c *bootstrap.Container) {
	defer c.Close()

	go func() {

	}()

	var (
		itrpSignal = make(chan os.Signal, 1)
	)

	signal.Notify(itrpSignal, syscall.SIGTERM, os.Interrupt, syscall.SIGINT)

	<-itrpSignal
	log.Println("gracefully shutdown server...")
}
