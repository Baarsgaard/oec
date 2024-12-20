package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/opsgenie/oec/conf"
	"github.com/opsgenie/oec/queue"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var metricAddr = flag.String("oec-metrics", "7070", "The address to listen on for HTTP requests.")

var OECVersion string
var OECCommitVersion string

func main() {
	logrus.SetFormatter(conf.PrepareLogFormat())
	logrus.SetOutput(os.Stdout)

	logrus.Infof("OEC version is %s", OECVersion)
	logrus.Infof("OEC commit version is %s", OECCommitVersion)

	configuration, err := conf.Read()
	if err != nil {
		logrus.Fatalf("Could not read configuration: %s", err)
	}

	logrus.SetLevel(configuration.LogrusLevel)

	flag.Parse()
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		logrus.Infof("OEC-metrics serves in http://localhost:%s/metrics.", *metricAddr)
		logrus.Error("OEC-metrics error: ", http.ListenAndServe(":"+*metricAddr, nil))
	}()

	queueProcessor := queue.NewProcessor(configuration)
	queue.UserAgentHeader = fmt.Sprintf("%s/%s %s (%s/%s)", OECVersion, OECCommitVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH)

	go func() {
		if configuration.AppName != "" {
			logrus.Infof("%s is starting.", configuration.AppName)
		}
		err = queueProcessor.Start()
		if err != nil {
			logrus.Fatalln(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-signals:
		logrus.Infof("OEC will be stopped gracefully.")
		err := queueProcessor.Stop()
		if err != nil {
			logrus.Fatalln(err)
		}
	}

	os.Exit(0)
}
