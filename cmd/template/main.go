package main

import (
	"template/internal/params"
	"template/internal/status"
	"template/pkg/logger"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	time.Local = time.FixedZone("CST", 8*3600) // 东八

	status := status.GetStatus()
	if params.GetParams().Version {
		status.ShowVersion()
		return
	}

	logger.Setup("./log/", logrus.DebugLevel, false)

	status.LogVersion()
}
