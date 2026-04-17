package poller

import "log"

type LogSink struct{}

func (LogSink) HandleResult(r Result) {
	if r.Err != nil {
		log.Printf("[DOWN] id=%s url=%s err=%v dur=%s", r.ID, r.Endpoint, r.Err, r.Duration)
		return
	}
	upDown := "UP"
	if r.StatusCode >= 400 || r.StatusCode == 0 {
		upDown = "DOWN"
	}
	log.Printf("[%s] id=%s url=%s code=%d dur=%s", upDown, r.ID, r.Endpoint, r.StatusCode, r.Duration)
}
