package logrus

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
	"strings"
)

func LogrusDemo() {
	var logger = logrus.New()

	logFormatter := &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return funcname, fmt.Sprintf("%s:%v", filename, f.Line)
		},
	}


	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(logFormatter)
	logger.SetReportCaller(true)

	logger.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")
}