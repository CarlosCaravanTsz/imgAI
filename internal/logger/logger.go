package logger

import (
	"os"
	log "github.com/sirupsen/logrus"
)

var Logger *log.Logger // variable global, se puede usar globalmente

func init() {
	Logger = log.New()

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatal("No se pudo abrir el archivo de logs", err)
	}

	Logger.Out = file

	Logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	Logger.SetLevel(log.DebugLevel)
}

func LogInfo(msg string, fields log.Fields) {
	Logger.WithFields(fields).Info(msg)
}

func LogError(msg string, fields log.Fields) {
	Logger.WithFields(fields).Error(msg)
}
