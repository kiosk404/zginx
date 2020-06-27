/**
* @Author: kiosk
* @Mail: weijiaxiang007@foxmail.com
* @Date: 2020/6/20
**/
package zlog

import (
	"github.com/sirupsen/logrus"
	"os"
)


var Logger = logrus.New()

func init() {
	switch os.Getenv("ZginxLog") {
	case "DEBUG":
		Logger.SetLevel(logrus.DebugLevel)
	case "INFO":
		Logger.SetLevel(logrus.InfoLevel)
	case "WARN":
		Logger.SetLevel(logrus.WarnLevel)
	case "ERROR":
		Logger.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		Logger.SetLevel(logrus.FatalLevel)
	case "PANIC":
		Logger.SetLevel(logrus.PanicLevel)
	default:
		Logger.SetLevel(logrus.DebugLevel)
	}
	//Logger.SetReportCaller(true)
	Logger.SetReportCaller(false)
	Logger.SetOutput(os.Stdout)
	Logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
}

func ConnLog(connId uint32) *logrus.Entry {
	return Logger.WithFields(logrus.Fields{
		"ConnId": connId,
	})
}