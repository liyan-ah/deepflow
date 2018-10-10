package queue

import (
	"sync"
	"testing"

	rawqueue "gitlab.x.lan/yunshan/droplet-libs/queue"
)

func TestSingleQueueSingleUserPuts(t *testing.T) {
	queue := &MultiQueue{}
	queue.Init("whatever", 8, 1, 1)
	keys := []rawqueue.HashKey{0, 2, 2, 1, 2, 1}
	inBatch := []interface{}{10081, 10082, 10083, 10084, 10085}
	outBatch := make([]interface{}, 2)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if out := queue.Get(2); out != 10081 {
			t.Errorf("Expected 10081, actually %d", out)
		}
		if count := queue.Gets(1, outBatch); count != 2 || outBatch[0] != 10082 || outBatch[1] != 10083 {
			t.Errorf("Expected 2 values: 10082 & 10083, actually %d values: %d & %d", count, outBatch[0], outBatch[1])
		}
		if count := queue.Gets(2, outBatch); count != 2 || outBatch[0] != 10084 || outBatch[1] != 10085 {
			t.Errorf("Expected 2 values: 10084 & 10085, actually %d values: %d & %d", count, outBatch[0], outBatch[1])
		}
		wg.Done()
	}()
	queue.Puts(keys, inBatch)
	wg.Wait()
}

func TestMultipleQueueSingleUserPuts(t *testing.T) {
	queue := &MultiQueue{}
	queue.Init("whatever", 8, 3, 1)
	keys := []rawqueue.HashKey{0, 2, 2, 1, 2, 1}
	inBatch := []interface{}{10081, 10082, 10083, 10084, 10085}
	outBatch := make([]interface{}, 8)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if out := queue.Get(2); out != 10081 {
			t.Errorf("Expected 10081, actually %d", out)
		}
		if count := queue.Gets(1, outBatch); count != 2 || outBatch[0] != 10083 || outBatch[1] != 10085 {
			t.Errorf("Expected 2 values: 10083 & 10085, actually %d values: %d & %d", count, outBatch[0], outBatch[1])
		}
		if count := queue.Gets(2, outBatch); count != 2 || outBatch[0] != 10082 || outBatch[1] != 10084 {
			t.Errorf("Expected 2 values: 10082 & 10084, actually %d values: %d & %d", count, outBatch[0], outBatch[1])
		}
		wg.Done()
	}()
	queue.Puts(keys, inBatch)
	wg.Wait()
}

func TestMultipleQueueMultipleUserPuts(t *testing.T) {
	queue := &MultiQueue{}
	userCount := 8
	size := userCount * 8
	queue.Init("whatever", size, 3, userCount)
	wg := sync.WaitGroup{}
	wg.Add(userCount)
	userPuts := func(i int) {
		keys := []rawqueue.HashKey{rawqueue.HashKey(i), 0, 2, 2, 1, 2, 1}
		inBatch := []interface{}{10081, 10082, 10083, 10084, 10085, 10086}
		queue.Puts(keys, inBatch)
		wg.Done()
	}
	for i := 0; i < userCount; i++ {
		go userPuts(i)
	}
	outBatch := make([]interface{}, size)
	wg.Wait()
	if count := queue.Gets(0, outBatch); count != userCount {
		t.Errorf("Expected %d values, actually %d values", userCount, count)
	}
	if count := queue.Gets(1, outBatch); count != userCount*2 {
		t.Errorf("Expected %d values, actually %d values", userCount*2, count)
	}
	if count := queue.Gets(2, outBatch); count != userCount*3 {
		t.Errorf("Expected %d values, actually %d values", userCount*3, count)
	}
}
