package main

import (
	"os"
	"os/signal"
	"syscall"
	"template/internal/params"
	"template/internal/status"
	"template/pkg/logger"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	defer saveDump()

	time.Local = time.FixedZone("CST", 8*3600) // 东八

	status := status.GetStatus()
	if params.GetParams().Version {
		status.ShowVersion()
		return
	}

	logger.Setup("./log/", logrus.DebugLevel, false)

	status.LogVersion()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	sign := <-signalChan
	logrus.Infof("get os signal=%v", sign)
	logrus.Infof("get os signal=%v", sign.String())
}

func saveDump() {
	if err := recover(); err != nil {
		logrus.Errorf("panic=%v", err)
		time.Sleep(1 * time.Second)
		os.Exit(1)
	}
}
