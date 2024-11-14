package callback

import (
	"errors"
	"time"
)

// roundRobin distributes data to available workers in a round-robin manner.
// It iterates over all endpoints (workers) in the c.endPoints slice and
// sends data to the first available worker's message queue. If a worker
// is blocked, it skips to the next one. If all workers are blocked,
// it returns an error indicating unavailability.
func (c *Callback) roundRobin(data []byte) error {
	// Loop through all endpoints in c.endPoints to find an available worker.
	for i := 0; i < len(c.endPoints); i++ {

		// Calculate the index of the current worker based on roundRobinIndex.
		// Increment roundRobinIndex by 1, subtract 1 to match the zero-based
		// indexing in arrays, then use modulo to cycle through endpoints
		// continuously in a round-robin manner.
		index := int(c.roundRobinIndex.Add(1)-1) % len(c.endPoints)

		// Retrieve the worker at the calculated index.
		worker := c.endPoints[index]

		// Check if this worker is available by comparing the current time with
		// worker.blockedUntil. If blockedUntil is in the future, the worker is
		// considered unavailable, so we continue to the next worker.
		if !time.Now().After(worker.blockedUntil) {
			continue
		}

		// If the worker is available, send the data to the worker's message queue.
		// worker.messageQueue is assumed to be a channel that the worker uses to receive tasks.
		worker.messageQueue <- data

		// Exit the function after sending data to one worker, ensuring that only
		// one worker processes this particular data payload in each roundRobin call.
		return nil
	}

	// If no worker was available, return an error indicating that all workers
	// are currently blocked and unable to process the data.
	return errors.New("all endpoints are blocked due to unavailability")
}
