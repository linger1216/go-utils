package pattern

import (
	"fmt"
	"sync"
	"time"
)

type Batch struct {
	size     int
	duration time.Duration
	stop     chan struct{}
	in       chan interface{}
	out      chan []interface{}
	flush    chan struct{}
	wg       *sync.WaitGroup
}

func NewBatch(batchSize, bufferBatchCount int, duration time.Duration) *Batch {
	return &Batch{
		size:     batchSize,
		duration: duration,
		stop:     make(chan struct{}),
		in:       make(chan interface{}, bufferBatchCount*batchSize),
		out:      make(chan []interface{}),
		flush:    make(chan struct{}),
	}
}

func (b *Batch) stopTimer() {

}

func (b *Batch) Start() {
	if b.wg != nil {
		return
	}

	timer := time.NewTimer(b.duration)
	timer.Stop()

	var batch []interface{}

	emit := func() {
		timer.Stop()
		select {
		case <-timer.C:
		default:
		}

		if len(batch) == 0 {
			return
		}
		b.out <- batch
		batch = nil
	}

	b.wg = &sync.WaitGroup{}
	b.wg.Add(1)

	go func() {
		defer b.wg.Done()
		for {
			select {
			case <-b.stop:
				emit()
				return
			case p := <-b.in:
				if batch == nil {
					if b.size > 0 {
						batch = make([]interface{}, 0, b.size)
					}
					if b.duration > 0 {
						timer.Reset(b.duration)
					}
				}

				batch = append(batch, p)
				if len(batch) >= b.size {
					emit()
				}

			case <-b.flush:
				emit()
			case <-timer.C:
				fmt.Printf("timer driver")
				emit()
			}
		}
	}()
}

func (b *Batch) Flush() {
	b.flush <- struct{}{}
}

func (b *Batch) In() chan<- interface{} {
	return b.in
}

func (b *Batch) Out() <-chan []interface{} {
	return b.out
}

func (b *Batch) Stop() {
	if b.wg == nil {
		return
	}
	// need wait
	time.Sleep(b.duration + time.Millisecond*100)
	close(b.stop)
	b.wg.Wait()
}
