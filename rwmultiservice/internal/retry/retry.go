package retry

import (
	"log"
	"sync"
	"time"

	"github.com/Reshnyak/innopolis/rwmultiservice/internal/models"
)

type Action int

const (
	stop    time.Duration = -1
	Succeed Action        = iota
	Fail
	Retry
)

type RetryWorker func(msg models.Message) <-chan error
type RetryPolice func(err error) Action
type FnBackoff func(attemp int) time.Duration
type Retrier struct {
	timeout     time.Duration
	attempMax   int
	attemptNum  int
	backoff     FnBackoff
	retryPolicy RetryPolice
	mu          *sync.Mutex
}

func simpleRetryPolicy(err error) Action {
	if err != nil {
		return Retry
	}
	return Succeed
}

// Возвращает интервал который увеличивается каждый раз.
// По аналогии можно создать фунrции с иной логикой.
func LinerBackoff(timeout time.Duration) FnBackoff {
	return func(attemp int) time.Duration {
		return time.Duration(attemp) * timeout
	}
}

// NewBackoff constructor
func NewRetrie(timeOut time.Duration, attempMax int) *Retrier {
	return &Retrier{
		timeout:     timeOut,
		attempMax:   attempMax,
		backoff:     LinerBackoff(timeOut),
		retryPolicy: simpleRetryPolicy,
		mu:          new(sync.Mutex),
	}
}

// Next duration time to retry
func (r *Retrier) Run(msg models.Message, retryWorker RetryWorker) error {
	// r.mu.Lock()
	// defer r.mu.Unlock()
	defer r.Reset()
	for {
		err := <-retryWorker(msg)
		switch r.retryPolicy(err) {
		case Succeed, Fail:
			return err
		case Retry:
			log.Printf("retry write %s attemp: %d", msg.FileId, r.attemptNum)
			var delay time.Duration
			if delay = r.next(); stop == delay {
				return err
			}
			<-time.After(delay)
		default:
			return err
		}
	}
}

// Next duration time to retry
func (r *Retrier) next() time.Duration {
	if r.attemptNum >= r.attempMax {
		return stop
	}
	r.attemptNum++
	return r.backoff(r.attemptNum)
}

// reset retry
func (r *Retrier) Reset() {
	r.attemptNum = 0
}
