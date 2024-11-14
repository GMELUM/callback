package callback

import (
	"testing"
	"time"
)

// TestRoundRobin_Success tests that data is successfully sent to an available worker.
func TestRoundRobin_Success(t *testing.T) {
	// Create a message queue and an available worker (blockedUntil is in the past).
	messageQueue := make(chan []byte, 1)
	worker := &Worker{
		messageQueue: messageQueue,
		blockedUntil: time.Now().Add(-time.Minute), // worker is immediately available
	}
	callback := &Callback{
		endPoints: []*Worker{worker},
	}

	// Attempt to send data and check that no error is returned.
	data := []byte("test data")
	err := callback.roundRobin(data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify that the data was sent to the worker's message queue.
	select {
	case result := <-worker.messageQueue:
		if string(result) != string(data) {
			t.Errorf("expected data %s, got %s", data, result)
		}
	default:
		t.Error("expected data to be sent to the worker, but queue was empty")
	}
}

// TestRoundRobin_AllBlocked tests that an error is returned when all workers are blocked.
func TestRoundRobin_AllBlocked(t *testing.T) {
	// Create a worker that is blocked (blockedUntil is in the future).
	messageQueue := make(chan []byte, 1)
	worker := &Worker{
		messageQueue: messageQueue,
		blockedUntil: time.Now().Add(time.Minute), // worker is blocked
	}
	callback := &Callback{
		endPoints: []*Worker{worker},
	}

	// Attempt to send data and check for the expected error.
	data := []byte("test data")
	err := callback.roundRobin(data)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "all endpoints are blocked due to unavailability" {
		t.Errorf("expected error 'all endpoints are blocked due to unavailability', got %v", err)
	}
}

// TestRoundRobin_RoundRobinOrder tests the round-robin ordering, ensuring data
// is sent to each worker in sequence as they become available.
func TestRoundRobin_RoundRobinOrder(t *testing.T) {
	// Create two workers. The first worker is available immediately,
	// while the second worker is initially blocked.
	messageQueue1 := make(chan []byte, 1)
	messageQueue2 := make(chan []byte, 1)

	worker1 := &Worker{
		messageQueue: messageQueue1,
		blockedUntil: time.Now().Add(-time.Minute), // worker1 is available immediately
	}
	worker2 := &Worker{
		messageQueue: messageQueue2,
		blockedUntil: time.Now().Add(time.Minute), // worker2 is initially blocked
	}
	callback := &Callback{
		endPoints: []*Worker{worker1, worker2},
	}

	// First call should send data to the first available worker (worker1).
	data1 := []byte("test data 1")
	err := callback.roundRobin(data1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	select {
	case result := <-worker1.messageQueue:
		if string(result) != string(data1) {
			t.Errorf("expected data %s for worker1, got %s", data1, result)
		}
	default:
		t.Error("expected data to be sent to worker1, but queue was empty")
	}

	// Now make worker2 available by setting blockedUntil to a past time.
	worker2.blockedUntil = time.Now().Add(-time.Minute) // worker2 becomes available
	data2 := []byte("test data 2")

	// Second call should now send data to worker2 in a round-robin sequence.
	err = callback.roundRobin(data2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	select {
	case result := <-worker2.messageQueue:
		if string(result) != string(data2) {
			t.Errorf("expected data %s for worker2, got %s", data2, result)
		}
	default:
		t.Error("expected data to be sent to worker2, but queue was empty")
	}
}