package logger

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"time"
)

// thanks a lot to https://blog.sandipb.net/2018/05/03/using-zap-creating-custom-encoders/
func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("Jan  2 15:04:05"))
}

// NewLoggerWithFile returns a new logger with file output.
//
// Sample usage: NewLoggerWithFile("/var/log/myproject/myproject.log")
func NewLoggerWithFile(filepath string) (*zap.SugaredLogger, error) {
	// prepare the logger config
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = SyslogTimeEncoder
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.StacktraceKey = "stacktrace"
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.FunctionKey = "function"
	cfg.OutputPaths = []string{
		filepath,
		"stderr", // for also writing to terminal
	}
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}

// CreateNewFile checks filepath if it exists or not, if it does not exist, creates it.
//
// Use it like this: CreateNewFile("./logs/backend")
//
// This will create a backend_(today'date).log file.
//
// Returns filepath in the end, such as ./logs/backend_20210101.log
func CreateNewFile(filepath string) string {
	LogFilepath := fmt.Sprintf("%s_%s.log", filepath, time.Now().Format("20060102"))
	if _, err := os.Stat(LogFilepath); errors.Is(err, os.ErrNotExist) {
		_, err = os.Create(LogFilepath)
		if err != nil {
			log.Println(err)
		}
	}
	return LogFilepath
}

// CreateNewFileCollector checks filepath if it exists or not, if it does not exist, creates it.
//
// Use it like this: CreateNewFileCollector("./logs/canercetin", "canercetin")
//
// This will create a collector_canercetin_(today's date).log file under the canercetin folder in the logs folder.
//
// Returns filepath in the end, such as ./logs/backend_20210101.log
//
// Only difference between this and CreateNewFile is, if file exists, we will put a number at the end of the file.
func CreateNewFileCollector(filepath string, username string) (string, int) {
	fileNumber := 1
	LogFilepath := fmt.Sprintf("%s/collector_%s_%s_%d.log", filepath, username, time.Now().Format("20060102"), 0)
	if _, err := os.Stat(LogFilepath); errors.Is(err, os.ErrNotExist) {
		_, err = os.Create(LogFilepath)
		if err != nil {
			log.Println(err)
		}
	} else {
		for {
			// put a number at the end of the file, such as collector_canercetin_20210101_1.log
			// increment the number until we find a file that does not exist.
			LogFilepath = fmt.Sprintf("%s/collector_%s_%s_%d.log", filepath, username, time.Now().Format("20060102"), fileNumber)
			if _, err = os.Stat(LogFilepath); errors.Is(err, os.ErrNotExist) {
				_, err = os.Create(LogFilepath)
				if err != nil {
					log.Println(err)
				}
				break
			}
			fileNumber++
		}
	}
	return LogFilepath, fileNumber
}

// CreateNewFileError is same as CreateNewFileCollector, but it creates a file called error_(today's date).log
//
// Seperating error logs from collector logs, it will help.
func CreateNewFileError(filepath string, username string) (string, int) {
	fileNumber := 1
	LogFilepath := fmt.Sprintf("%s/error_%s_%s_%d.log", filepath, username, time.Now().Format("20060102"), 0)
	if _, err := os.Stat(LogFilepath); errors.Is(err, os.ErrNotExist) {
		_, err = os.Create(LogFilepath)
		if err != nil {
			log.Println(err)
		}
	} else {
		for {
			// put a number at the end of the file, such as collector_canercetin_20210101_1.log
			// increment the number until we find a file that does not exist.
			LogFilepath = fmt.Sprintf("%s/error_%s_%s_%d.log", filepath, username, time.Now().Format("20060102"), fileNumber)
			if _, err = os.Stat(LogFilepath); errors.Is(err, os.ErrNotExist) {
				_, err = os.Create(LogFilepath)
				if err != nil {
					log.Println(err)
				}
				break
			}
			fileNumber++
		}
	}
	return LogFilepath, fileNumber
}

// CreateNewFolder is used for creating a new folder under ./logs.
//
// Use this like, CreateNewFolder(canercetin) and it will create a folder called canercetin under ./logs.
//
// Then use CreateNewFileCollector to create a new file under that folder, by CreateNewdFileCollector("./logs/canercetin", "canercetin")
//
// Or CreateNewFile("./logs/canercetin") whatever.
func CreateNewFolder(username string) error {
	if _, err := os.Stat("./logs/" + username); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir("./logs/"+username, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
