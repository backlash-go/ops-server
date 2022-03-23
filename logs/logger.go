package logs

import (
	"fmt"
	"log"

	go_logger "github.com/phachon/go-logger"
	"github.com/spf13/viper"
)

var logger *go_logger.Logger

// InitLogger 初始化 logger
func InitLogger(env string) {
	logger = go_logger.NewLogger()
	_ = logger.Detach("console")
	logFormat := "[%timestamp_format%] [%level_string%] %function%:%line% %body%"
	if env == "dev" {
		logger.Detach("console")
		// console adapter config
		consoleConfig := &go_logger.ConsoleConfig{
			Color:      true,      // Does the text display the color
			JsonFormat: false,     // Whether or not formatted into a JSON string
			Format:     logFormat, // JsonFormat is false, logger message output to console format string
		}
		// add output to the console
		logger.Attach("console", go_logger.LOGGER_LEVEL_INFO, consoleConfig)
		log.Println("log to : console")
		return
	}

	logDir := viper.GetString("logDir")
	// file adapter config
	fileConfig := &go_logger.FileConfig{
		LevelFileName: map[int]string{
			logger.LoggerLevel("error"): fmt.Sprintf("%s/ops-server.log", logDir),
			logger.LoggerLevel("info"):  fmt.Sprintf("%s/ops-server.log", logDir),
			logger.LoggerLevel("debug"): fmt.Sprintf("%s/ops-server.log", logDir),
		},
		MaxSize:    0,
		MaxLine:    0,
		DateSlice:  "y",
		JsonFormat: false,
		Format:     logFormat,
	}

	_ = logger.Attach("file", go_logger.LOGGER_LEVEL_DEBUG, fileConfig)
}

func GetLogger() *go_logger.Logger {
	return logger
}
