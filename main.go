package main

import (
	"github.com/high-performance-payment-gateway/balance-service/balance/pkg/external/log_init"
	log "github.com/sirupsen/logrus"
)

func main() {
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
}
