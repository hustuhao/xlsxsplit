package logger

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"github.com/metafates/xlsxsplit/filesystem"
	"github.com/metafates/xlsxsplit/key"
	"github.com/metafates/xlsxsplit/where"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

func Init() error {
	logsPath := where.Logs()

	if logsPath == "" {
		return errors.New("logs path is not set")
	}

	today := time.Now().Format("2006-01-02")
	logFilePath := filepath.Join(logsPath, fmt.Sprintf("%s.log", today))
	if !lo.Must(filesystem.Api().Exists(logFilePath)) {
		lo.Must(filesystem.Api().Create(logFilePath))
	}

	logFile, err := filesystem.Api().OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	// if you want to write to a file and stdout, use this: multiWriter := io.MultiWriter(os.Stdout, logFile)
	multiWriter := io.MultiWriter(logFile)
	logger := log.NewWithOptions(multiWriter, log.Options{
		TimeFormat:      time.TimeOnly,
		ReportTimestamp: true,
		ReportCaller:    viper.GetBool(key.LogsReportCaller),
	})

	level, err := log.ParseLevel(viper.GetString(key.LogsLevel))
	if err != nil {
		log.Fatal(err)
	}
	logger.SetLevel(level)

	log.SetDefault(logger)
	return nil
}
