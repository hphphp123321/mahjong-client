package cmd

import (
	"github.com/hphphp123321/mahjong-client/client"
	"github.com/hphphp123321/mahjong-common/osutils"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
	"strconv"
)

func SetupSignalHandler(c *client.MahjongClient) {
	go func() {
		exitChannel := osutils.NewShutdownSignal()
		for {
			osutils.WaitExit(exitChannel, func() {
				log.Println("Exit")
				if err := c.Logout(); err != nil {
					log.Fatalf("Logout failed: %v", err)
				}
			})
		}
	}()
}

func setupLogger() {
	switch logFormat {
	case "text":
		log.SetFormatter(&log.TextFormatter{
			ForceColors:               true,
			TimestampFormat:           "2006-01-02 15:04:05",
			FullTimestamp:             true,
			EnvironmentOverrideColors: true,
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				//处理文件名
				fileName := path.Base(frame.File)
				return ": " + strconv.Itoa(frame.Line), fileName
			},
		})
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.Error("set log format error")
	}

	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	default:
		log.Error("set log level error")
	}

	switch logOutput {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	default:
		log.Error("set log output error")
	}

	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(f)
		} else {
			log.Error("set log file error")
		}
	}

}