package pattern

import (
	"fmt"
	"sync"
	"time"
)

type Batch struct {
	batch    []interface{}
	size     int
	timer    *time.Timer
	duration time.Duration
	stop     chan struct{}
	in       chan interface{}
	out      chan []interface{}
	flush    chan struct{}
	wg       *sync.WaitGroup
	closed   bool
}

func NewBatch(batchSize, chanSize int, duration time.Duration) *Batch {
	return &Batch{
		batch:    make([]interface{}, 0, batchSize),
		size:     batchSize,
		duration: duration,
		stop:     make(chan struct{}),
		in:       make(chan interface{}, chanSize),
		out:      make(chan []interface{}),
		flush:    make(chan struct{}),
	}
}

func (b *Batch) Start() {
	if b.wg != nil {
		return
	}

	b.timer = time.NewTimer(b.duration)
	b.stopTimer(true)

	b.wg = &sync.WaitGroup{}
	b.wg.Add(1)

	go func() {
		defer b.wg.Done()
		for {
			if b.closed {
				return
			}
			select {
			case <-b.stop:
				close(b.in)
				b.tryEmit()
				close(b.out)
				return
			case p := <-b.in:
				if len(b.batch) == 0 && b.duration > 0 {
					b.timer.Reset(b.duration)
				}

				b.batch = append(b.batch, p)
				if len(b.batch) >= b.size {
					b.emit()
				}
			case <-b.flush:
				b.emit()
			case <-b.timer.C:
				fmt.Printf("timer driver\n")
				b.emit()
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

func (b *Batch) stopTimer(force bool) {
	if b.duration > 0 || force {
		b.timer.Stop()
		select {
		case <-b.timer.C:
		default:
		}
	}
}

func (b *Batch) emit() {
	b.stopTimer(false)
	if len(b.batch) == 0 {
		return
	}
	select {
	case b.out <- b.batch:
		b.batch = make([]interface{}, 0, b.size)
	case <-b.stop:
		fmt.Printf("close batch when sending data!")
		close(b.in)
		close(b.out)
		b.closed = true
	}
}

func (b *Batch) tryEmit() {
	b.stopTimer(false)
	if len(b.batch) == 0 {
		return
	}
	fmt.Printf("waiting %v to send data\n", b.duration+time.Millisecond*100)
	select {
	case b.out <- b.batch:
		b.batch = make([]interface{}, 0, b.size)
	case <-time.After(b.duration + time.Millisecond*100):
		fmt.Println("sending data timeout, close forcely")
	}
}
