package main

import (
	"deepwildcard/internal/deepwildcard"
	"flag"
	"log"
	"os"
)

var defaultConfigFile = "/etc/deepwildcard/config.yaml"

func main() {
	printHeader()
	configFile := flag.String("config", defaultConfigFile, "")
	logTime := flag.Bool("log-time", false, "")

	flag.Usage = printUsage
	flag.Parse()

	logFlags := 0
	if *logTime {
		logFlags |= log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC
	}

	dw, err := deepwildcard.New(
		deepwildcard.WithConfigFile(*configFile),
		deepwildcard.WithLogger(log.New(os.Stdout, "", logFlags)),
	)
	if err != nil {
		panic(err)
	}
	dw.ListenAndServe()
}
