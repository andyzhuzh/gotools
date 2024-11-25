package log
import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	JsonFormat     = "JSON"
	TextFormat     = "TEXT"
	OutputTypeFile = "LOGFILE"
	OutputTypeShow = "LOGShow"
	InfoLog        = "Info"
	DebugLog       = "Debug"
	TraceLog       = "Trace"
	WarnLog        = "Warn"
	ErrorLog       = "Error"
)

type Logger struct {
	FileName string
	FileObj  *os.File
	Logger   *logrus.Logger
}

// LogAdd(JsonFormat,OutputTypeFile,InfoLog,"Test Message","",map[string]interface{}{"name": "Jane", "age": 25})
// LogAdd(JsonFormat, OutputTypeShow, InfoLog, "Test Message", "", map[string]interface{}{"name": "Jane", "age": 25})
// LogAdd(JsonFormat, OutputTypeShow, InfoLog, "Test Message", "", nil)
func LogAdd(formatType, outputType, logType, logMessage, fileName string, fields map[string]interface{}) {
	logger, _ := NewLogger(formatType)
	if outputType == OutputTypeFile {
		logger.SetOutputFile(formatType, fileName)
	}
	switch logType {
	case ErrorLog:
		logger.AddError(fields, logMessage)
	case DebugLog:
		logger.AddDebug(fields, logMessage)
	case TraceLog:
		logger.AddTrace(fields, logMessage)
	case WarnLog:
		logger.AddWarn(fields, logMessage)
	default:
		logger.AddInfo(fields, logMessage)
	}
	logger.LogClose()
}

func AddShow(logMessage string, arg ...string) {
	logType := InfoLog
	if len(arg) > 0 {
		logType = arg[0]
	}
	LogAdd(TextFormat, OutputTypeShow, logType, logMessage, "", nil)
}

func AddFile(logMessage string, arg ...string) {
	logType := InfoLog
	fileName := ""
	if len(arg) > 0 {
		logType = arg[0]
	}
	if len(arg) > 1 {
		fileName = arg[1]
	}
	LogAdd(TextFormat, OutputTypeFile, logType, logMessage, fileName, nil)
}

func NewLogger(logFormmat string) (logger Logger, err error) {
	if logFormmat == JsonFormat {
		logger.Logger = logrus.New()
		logger.Logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.DateTime})
	} else {
		logger.Logger = logrus.New()
		logger.Logger.SetFormatter(&logrus.TextFormatter{DisableColors: false, FullTimestamp: true, TimestampFormat: time.DateTime})
	}
	logger.Logger.SetLevel(logrus.InfoLevel)
	// logrus.Info()
	// logrus.Debug()
	// logrus.Trace()
	// logrus.Warn()
	// logrus.Error()
	return
}

func (logObj *Logger) SetOutputFile(logFormmat string, fileName string) error {
	if logObj.Logger == nil {
		return fmt.Errorf("没有日志对象")
	}
	now := time.Now()
	fileFullNane := fileName
	if fileName == "" {
		logFilePath := ""
		if dir, err := os.Getwd(); err == nil {
			logFilePath = path.Join(dir, "logs")

		}
		_, err := os.Stat(logFilePath)
		if os.IsNotExist(err) {
			if err := os.MkdirAll(logFilePath, 0777); err != nil {
				return err
			}
		}
		fileExt := ".log"
		if logFormmat == JsonFormat {
			fileExt = ".json"
		}

		logFileName := "Log_" + now.Format("20060102") + "_" + now.Format("150405") + "_" + strconv.Itoa(now.Nanosecond()) + fileExt
		fileFullNane = path.Join(logFilePath, logFileName)
	}
	logObj.FileName = fileFullNane
	src, err := os.OpenFile(fileFullNane, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	logObj.FileObj = src
	logObj.Logger.Out = src
	if err != nil {
		return err
	}
	return nil
	// logger.SetOutput(file)
	// logger.Out = file
	// .SetOutput(os.Stdout)
}

func (logObj *Logger) LogClose() (err error) {
	logObj.Logger.Writer().Close()
	logObj.FileObj.Close()
	logObj = nil
	return
}

func (logObj *Logger) AddInfo(fields logrus.Fields, args ...interface{}) {
	if fields != nil {
		logEntry := logObj.Logger.WithFields(fields)
		logEntry.Info(args...)
		return
	}
	logObj.Logger.Info(args...)
}

func (logObj *Logger) AddDebug(fields logrus.Fields, args ...interface{}) {
	if fields != nil {
		logEntry := logObj.Logger.WithFields(fields)
		logEntry.Debug(args...)
		return
	}
	logObj.Logger.Debug(args...)
}

func (logObj *Logger) AddError(fields logrus.Fields, args ...interface{}) {
	if fields != nil {
		logEntry := logObj.Logger.WithFields(fields)
		logEntry.Error(args...)
		return
	}
	logObj.Logger.Error(args...)
}

func (logObj *Logger) AddWarn(fields logrus.Fields, args ...interface{}) {
	if fields != nil {
		logEntry := logObj.Logger.WithFields(fields)
		logEntry.Warn(args...)
		return
	}
	logObj.Logger.Warn(args...)
}

func (logObj *Logger) AddTrace(fields logrus.Fields, args ...interface{}) {
	if fields != nil {
		logEntry := logObj.Logger.WithFields(fields)
		logEntry.Trace(args...)
		return
	}
	logObj.Logger.Trace(args...)
}
