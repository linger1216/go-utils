package pattern

import (
	"github.com/bmizerany/assert"
	"testing"
	"time"
)

func TestBatchWithoutTimer(t *testing.T) {
	batch := NewBatch(10, 1, 0)
	batch.Start()
	in, out := batch.In(), batch.Out()

	go func() {
		defer batch.Stop()
		for i := 0; i < 101; i++ {
			in <- i
		}
	}()

	for e := range out {
		t.Log(e)
	}
}

func TestBatchWithTimer(t *testing.T) {
	batch := NewBatch(10, 1, time.Second)
	batch.Start()
	in, out := batch.In(), batch.Out()

	go func() {
		defer batch.Stop()
		for i := 0; i < 10; i++ {
			in <- i*3 + 0
			in <- i*3 + 1
			in <- i*3 + 2
			time.Sleep(time.Second)
		}
	}()

	for e := range out {
		t.Log(e)
	}
}

func TestBatchStopWhenBlocked_1(t *testing.T) {
	batch := NewBatch(10, 1, time.Second)
	batch.Start()
	in, _ := batch.In(), batch.Out()

	go func() {
		defer batch.Stop()
		for i := 0; i < 101; i++ {
			t.Log("trying to send message: ", i)
			in <- i
		}
	}()

	time.Sleep(time.Second * 2)
	batch.Stop()
}

func TestBatchStopWhenBlocked_2(t *testing.T) {
	batch := NewBatch(10, 0, 0)
	batch.Start()
	in, _ := batch.In(), batch.Out()

	for i := 0; i < 3; i++ {
		in <- i
	}

	batch.Stop()
}

func TestBatchStopWhenSendingData(t *testing.T) {
	batch := NewBatch(100, 0, time.Second)
	batch.Start()
	in, out := batch.In(), batch.Out()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				assert.Equal(t, "send on closed channel", err.(error).Error())
			}
		}()
		for {
			in <- 1
		}
	}()

	go func() {
		time.Sleep(time.Second * 3)
		batch.Stop()
	}()

	for _ = range out {

	}
}
