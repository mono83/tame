package page

import (
	"fmt"
	"time"
)

// NewTrace returns new trace struct with initialized initial time.
func NewTrace() *Trace {
	return &Trace{Init: time.Now()}
}

// Trace contains HTTP tracing information for page
type Trace struct {
	Init time.Time

	DNS         time.Duration // DNS resolve time
	Connection  time.Duration // TCP connection time
	RequestSent time.Duration // Time until whole request was sent
	Total       time.Duration // Total time
	App         time.Duration // Application time (Total - RequestSent)
}

// Close closes trace and writes total time.
// Works only once
func (t *Trace) Close() {
	if t.Total.Nanoseconds() == 0 {
		t.Total = time.Now().Sub(t.Init)
		t.App = time.Duration(int64(t.Total) - int64(t.RequestSent))
	}
}

func (t Trace) String() string {
	return fmt.Sprintf(
		"DNS: %.3f, Connect: %.3f, Sent: %.3f, Total: %.3f, App: %.3f",
		t.DNS.Seconds(),
		t.Connection.Seconds(),
		t.RequestSent.Seconds(),
		t.Total.Seconds(),
		t.App.Seconds(),
	)
}
