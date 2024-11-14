package callback

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Callback manages the sending of messages to multiple worker endpoints with configurable retry settings and delivery modes.
type Callback struct {
	transport       Transport     // Transport defines the method of communication with workers.
	deliveryMode    DeliveryMode  // DeliveryMode controls how messages are sent: RoundRobin or Broadcast.
	endPoints       []*Worker     // List of worker endpoints that handle message delivery.
	retryLimit      int           // Number of retry attempts allowed before giving up.
	retryTimeout    time.Duration // Wait time between retry attempts.
	retryWindow     time.Duration // Time window in which retries are allowed.
	roundRobinIndex atomic.Int32  // Index used for RoundRobin delivery mode to track the last worker.
	returnChannel   chan Data     // Channel for returning data back to the callback function.
	mu              sync.Mutex    // Mutex for concurrent access to endpoints.

	callback func(data *Data) // User-defined callback function to handle processed data.
}

// New initializes a new Callback instance with the provided options and sets up worker endpoints.
func New(opt *Options) *Callback {
	opt = defaultOptions(opt) // Apply default options if not provided.

	// Create a Callback instance and initialize fields with options.
	callback := &Callback{
		transport:     opt.Transport,
		deliveryMode:  opt.DeliveryMode,
		retryLimit:    opt.RetryLimit,
		retryTimeout:  opt.RetryTimeout,
		retryWindow:   opt.RetryWindow,
		returnChannel: make(chan Data, 100),
	}
	// Sync the initial set of endpoints provided in options.
	callback.SyncEndPoint(opt.EndPoints)

	// Launch a handler goroutine to listen on the return channel for incoming data.
	go callback.handler()

	return callback
}

// handler listens on the returnChannel and invokes the callback function when data is received.
// This handler recovers from panics to ensure it remains active.
func (c *Callback) handler() {
	// Recover from any panics to keep the handler operational.
	defer func() {
		if r := recover(); r != nil {

			// Log or handle the panic information.
			if c.callback != nil {
				c.callback(&Data{
					Point:   "",
					Success: false,
					Error: &Error{
						Code:     0,
						Message:  fmt.Sprintf("[PANIC] global error: %v", r),
						Critical: true,
					},
				})
			}

			go c.handler()

		}
	}()

	// Read and process each item from returnChannel
	for data := range c.returnChannel {
		if c.callback != nil {
			c.callback(&data) // Execute the user-defined callback function.
		}
	}
}

// findWorkerIndex locates the index of a worker based on its endpoint.
func (c *Callback) findWorkerIndex(endPoint string) int {
	for i, worker := range c.endPoints {
		if worker.point == endPoint {
			return i // Return index if worker is found.
		}
	}
	return -1 // Return -1 if worker is not found.
}

// AddEndpoint adds a new worker for the given endpoint.
func (c *Callback) AddEndpoint(host string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If the worker already exists, do not add again.
	if c.findWorkerIndex(host) != -1 {
		return
	}

	// Create and add a new worker for the endpoint.
	worker := NewWorker(c, host)
	c.endPoints = append(c.endPoints, worker)
}

// DeleteEndpoint removes and closes the worker for the given endpoint.
func (c *Callback) DeleteEndpoint(host string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Find the index of the worker to delete.
	index := c.findWorkerIndex(host)
	if index == -1 {
		return // Worker not found, exit without action.
	}

	// Close the worker and remove it from the slice.
	c.endPoints[index].Close()
	c.endPoints = append(c.endPoints[:index], c.endPoints[index+1:]...)
}

// SyncEndPoint synchronizes the current list of endpoints with a new list.
// It removes outdated workers and adds new ones.
func (c *Callback) SyncEndPoint(hosts []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Map for tracking new hosts and their presence.
	newHosts := make(map[string]struct{}, len(hosts))
	for _, host := range hosts {
		newHosts[host] = struct{}{}
	}

	// Remove outdated workers that are not in the new list of hosts.
	for i := len(c.endPoints) - 1; i >= 0; i-- {
		worker := c.endPoints[i]
		if _, exists := newHosts[worker.point]; !exists {
			worker.Close()
			c.endPoints = append(c.endPoints[:i], c.endPoints[i+1:]...) // Remove outdated worker.
		}
	}

	// Add new hosts that are not yet in the worker list.
	for _, host := range hosts {
		if c.findWorkerIndex(host) == -1 {
			c.endPoints = append(c.endPoints, NewWorker(c, host))
		}
	}
}

// Emit sends data to the workers based on the delivery mode.
func (c *Callback) Emit(data []byte) error {
	switch c.deliveryMode {
	case RoundRobin:
		c.roundRobin(data)
		return nil
	case Broadcast:
		// Broadcast to all workers (functionality to be implemented).
		return nil
	}
	return nil
}

// On sets a callback function to handle processed data received from the returnChannel.
func (c *Callback) On(clb func(data *Data)) {
	c.callback = clb
}
