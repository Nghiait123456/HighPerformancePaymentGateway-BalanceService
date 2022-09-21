package log_init

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

type (
	Log struct {
		TypeOutput string
		TypeFormat string
		LinkFile   string
		LinkFolder string
	}
)

const (
	TYPE_OUTPUT_FILE    = "file"
	TYPE_OUTPUT_CONSOLE = "json"
	TYPE_FORMAT_JSON    = "json"
	TYPE_FORMAT_TEXT    = "text"

	PACKET_LOG = "github.com/sirupsen/logrus"
)

func setOutPut(l Log) {
	switch l.TypeOutput {
	case TYPE_OUTPUT_FILE:
		errCFL := os.MkdirAll(l.LinkFolder, 0666)
		if errCFL != nil {
			panic(fmt.Sprintf("Failed to create folder log, error : %s", errCFL.Error()))
		}

		file, err := os.OpenFile(l.LinkFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("Failed to open log file, error : %s", err.Error()))
		}

		log.SetOutput(file)
		return

	case TYPE_OUTPUT_CONSOLE:
		log.SetOutput(os.Stdout)
		return

	default:
		log.SetOutput(os.Stdout)
	}
}

func setFormat(l Log) {
	switch l.TypeFormat {
	case TYPE_FORMAT_JSON:
		log.SetFormatter(&log.JSONFormatter{})
	case TYPE_FORMAT_TEXT:
		log.SetFormatter(&log.TextFormatter{})
	default:
		log.SetFormatter(&log.TextFormatter{})
	}
}

// Init setup log for github.com/sirupsen/logrus
func Init(l Log) {
	setOutPut(l)
	setFormat(l)
}

/**
   example :

   log_init.Init(log_init.Log{
		TypeFormat: log_init.TYPE_FORMAT_TEXT,
		TypeOutput: log_init.TYPE_OUTPUT_FILE,
		LinkFile:   "balance/infrastructure/log/log_file/log.log",
	})

	log.WithFields(log.Fields{
		"animal": "walrus",
		"size":   10,
	}).Warningf("A group of walrus emerges from the ocean")
	log.WithFields(log.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")
*/
