package main

import (
	"flag"
	"github.com/clemens97/scion-path-oracle/server"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		listen        string
		listenMonitor string
	)
	flag.StringVar(&listen, "listen", ":443", "Listening address for the path oracle.")
	flag.StringVar(&listenMonitor, "listenMonitoring", ":8080", "Non-SCION listening address for real time monitoring.")
	flag.Parse()

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("could not initialize logger: %s", err)
	}

	defer logger.Sync()
	slogger := logger.Sugar()

	go func() {
		handleSignals(signalChannel, slogger)
	}()

	oracle, err := server.New(listen, listenMonitor, slogger)
	if err != nil {
		slogger.Fatalw("error setting up path oracle", "error", err)
	}

	slogger.Infow("starting path oracle", "listen", listen, "listenMonitor", listenMonitor)
	err = oracle.Start()
	if err != nil {
		slogger.Fatalw("error running path oracle", "error", err)
	}

}

func handleSignals(sigChan chan os.Signal, logger *zap.SugaredLogger) {
	for sig := range sigChan {
		logger.Infow("captured signal, exiting", "signal", sig)
		os.Exit(0)
	}
}
