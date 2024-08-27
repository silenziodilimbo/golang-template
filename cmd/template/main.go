package main

import (
	"os"
	"os/signal"
	"runtime"
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
		stackBuf := make([]byte, 1024)
		for {
			stackSize := runtime.Stack(stackBuf, true)
			if stackSize < len(stackBuf) {
				stackBuf = stackBuf[:stackSize]
				break
			}
			stackBuf = make([]byte, 2*len(stackBuf))
		}

		// 打印堆栈信息
		logrus.Fatalf("exit with panic\n%s", string(stackBuf))
	}
}
