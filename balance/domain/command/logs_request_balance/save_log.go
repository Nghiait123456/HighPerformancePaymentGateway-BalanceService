package logs_request_balance

import (
	"math/rand"
	"os"
	"time"
)

type (
	OrderLog struct {
		OrderId uint64
		Amount  uint64
		Status  string
	}

	LogsChan = chan OrderLog

	Logs struct {
		allLog                          LogsChan
		maxNumberMustSaveLog            uint64
		maxRangeTimeAllowTwoTimesUpdate uint64
	}

	LogsInterface interface {
		Init()
		Push(o OrderLog)
		saveLogs() error
		sendLogsToServiceSave(s SendLogsToOtherService) error
	}

	SendLogsToOtherService struct {
		logs      []OrderLog
		requestId uint64
		maxRetry  uint
	}
)

const (
	MAXIMUM_LOG_MUST_SAVE_LOG_TO_OTHER_SERVICE = 5000
	MAX_DISTANCE_TIME_FOR_TWO_TIMES_UPDATE     = 400 // ms
	MAX_RETRY_SEND_LOGS                        = 5
)

func (l *Logs) Init() {
	l.allLog = make(LogsChan)
	l.maxNumberMustSaveLog = MAXIMUM_LOG_MUST_SAVE_LOG_TO_OTHER_SERVICE
	l.maxRangeTimeAllowTwoTimesUpdate = MAX_DISTANCE_TIME_FOR_TWO_TIMES_UPDATE
	go func() {
		l.saveLogs()
	}()
}

func (l *Logs) Push(o OrderLog) {
	l.allLog <- o
}

func (l *Logs) saveLogs() error {
	var logs []OrderLog
	timeout := time.After(time.Duration(l.maxRangeTimeAllowTwoTimesUpdate) * time.Millisecond)

	for {
		select {
		case log := <-l.allLog:
			logs = append(logs, log)

			if uint64(len(logs)) >= l.maxNumberMustSaveLog {
				s := SendLogsToOtherService{
					logs:      logs,
					requestId: rand.Uint64(),
					maxRetry:  MAX_RETRY_SEND_LOGS,
				}
				err := l.sendLogsToServiceSave(s)
				if err != nil {
					panic(err)
					os.Exit(0)
				}

				//reset logs
				logs = nil
				// reset and init timeout next
				timeout = time.After(time.Duration(l.maxRangeTimeAllowTwoTimesUpdate) * time.Millisecond)
			}

		case <-timeout:
			s := SendLogsToOtherService{
				logs:      logs,
				requestId: rand.Uint64(),
				maxRetry:  MAX_RETRY_SEND_LOGS,
			}

			err := l.sendLogsToServiceSave(s)
			if err != nil {
				panic(err)
				os.Exit(0)
			}

			//reset logs
			logs = nil
			// init timeout next
			timeout = time.After(time.Duration(l.maxRangeTimeAllowTwoTimesUpdate) * time.Millisecond)
		}
	}

}

func (l *Logs) sendLogsToServiceSave(s SendLogsToOtherService) error {
	for i := uint(0); i < s.maxRetry; i++ {
		// todo send logs
		// todo if success return
		// todo sleep x ms if send fail
	}

	return nil
}

func NewLog() LogsInterface {
	l := Logs{}
	l.Init()
	return &l
}
