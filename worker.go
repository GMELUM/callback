package callback

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gmelum/callback/transport"
)

// Worker represents a process that handles incoming data and interacts with an external callback interface.
type Worker struct {

	// A reference to the callback object the worker interacts with.
	callback *Callback

	// The point (or identifier) this worker is associated with.
	point string

	// A channel for receiving messages to process.
	messageQueue chan []byte

	// A channel to return the result (Response or Error) after processing.
	returnChannel chan Data

	// A mutex for synchronizing access to error data and blocking state.
	mu sync.Mutex

	// A slice holding timestamps of errors for retry logic.
	errorTimestamps []time.Time

	// The time until which the worker will be blocked if retry limits are exceeded.
	blockedUntil time.Time

	// A channel for stopping the worker.
	stop chan struct{}
}

// NewWorker creates a new Worker object and starts the necessary goroutines for processing data and handling responses.
func NewWorker(c *Callback, point string) *Worker {
	// Initialize a new Worker with required fields.
	worker := &Worker{

		// Set the callback object.
		callback: c,

		// Set the worker's point.
		point: point,

		// A buffered channel for message queue with a size of 2.
		messageQueue: make(chan []byte, 100),

		// Set the returnChannel from the callback.
		returnChannel: c.returnChannel,

		// Initialize the error timestamps with a maximum capacity.
		errorTimestamps: make([]time.Time, 0, c.retryLimit),
	}

	// Start another goroutine for processing incoming messages.
	go worker.handler()
	return worker
}

// handler is the main method for processing messages that are received from the messageQueue.
func (w *Worker) handler() {

	// defer is used to recover from any panics, ensuring the worker continues operating.
	defer func() {
		if r := recover(); r != nil { // If a panic occurs.
			// Send an error back to the return channel with panic information.
			w.returnChannel <- w.sendReturn(
				&Error{
					Code:     0,
					Message:  fmt.Sprintf("[PANIC] %v", r), // Format the panic message.
					Critical: true,
				},
			)

			go w.handler()

		}
	}()

	for {
		select {
		case data := <-w.messageQueue: // If a message is received from the messageQueue.
			// Process the message.
			res, err := w.handlerRequest(data)
			if err != nil { // If an error occurs while processing.
				// Increment the error count and send an error response.
				w.Inc()
				w.returnChannel <- w.sendReturn(
					&Error{
						Code:     0,
						Message:  fmt.Sprintf("[ERROR] %v", err.Error()), // Format the error message.
						Critical: true,
					},
				)
				continue // Continue processing the next message.
			}

			// If the processing succeeds, reset error counters and return the successful result.
			w.Reset()
			w.returnChannel <- w.sendReturn(&Response{res}) // Send the successful response.

		case <-w.stop: // If the stop signal is received.
			return // Exit the handler goroutine.
		}
	}

}

// sendReturn formats and sends the result (Response or Error) to the returnChannel.
func (w *Worker) sendReturn(result interface{}) Data {
	var success bool
	var res *Response
	var err *Error

	// Determine the type of result (Response or Error) and prepare it for sending.
	switch result.(type) {
	case ResponseInterface:
		success = true
		res = result.(*Response) // If it's a Response, cast and store it.
	case ErrorInterface:
		success = false
		err = result.(*Error) // If it's an Error, cast and store it.
	default:
		success = false
		err = &Error{Message: "Unknown result type"} // If the type is unknown, return an error.
	}

	// Return a Data object containing the point, success flag, and the respective response or error.
	return Data{
		// The point associated with the worker.
		Point: w.point,
		// Flag indicating whether the result is successful or an error.
		Success: success,
		// The response data if successful.
		Response: res,
		// The error data if failure occurred.
		Error: err,
	}
}

// handlerRequest processes incoming data. Currently a stub, needs to be implemented with specific logic.
func (w *Worker) handlerRequest(data []byte) ([]byte, error) {

	if w.callback.transport == REST {
		return transport.Post(w.point, data)
	}

	// TODO: Implement the logic to handle the incoming data (e.g., process the byte slice).
	return nil, errors.New("transport is not support")
}

// Inc increments the error count and checks if the worker should be blocked due to too many errors.
func (w *Worker) Inc() bool {
	w.mu.Lock()         // Lock for thread-safe access to shared resources.
	defer w.mu.Unlock() // Ensure the mutex is released when the method finishes.

	// Remove errors that are outside of the retry window.
	now := time.Now()                               // Get the current time.
	windowStart := now.Add(-w.callback.retryWindow) // The start of the retry window.
	filteredErrors := w.errorTimestamps[:0]         // Create a new slice to hold only recent errors.

	// Keep only errors that occurred within the retry window.
	for _, timestamp := range w.errorTimestamps {
		if timestamp.After(windowStart) {
			filteredErrors = append(filteredErrors, timestamp) // Append valid errors to the filtered slice.
		}
	}
	w.errorTimestamps = filteredErrors // Update the error timestamps with the filtered ones.

	// Add the current error timestamp.
	w.errorTimestamps = append(w.errorTimestamps, now)

	// If the number of errors exceeds the retry limit, check if the worker should be blocked.
	if len(w.errorTimestamps) > w.callback.retryLimit {
		// If the worker is not currently blocked, set the block time.
		if now.After(w.blockedUntil) {
			w.blockedUntil = now.Add(w.callback.retryTimeout) // Set the blockedUntil time to retryTimeout after the current time.
		}
		return true // Indicate that the worker is blocked due to too many errors.
	}

	// If the worker is not blocked, return false.
	return now.Before(w.blockedUntil) // Return whether the worker is still within the blocked period.
}

// Reset clears the error count and unblocks the worker.
func (w *Worker) Reset() {
	w.mu.Lock()         // Lock for thread-safe modification of the state.
	defer w.mu.Unlock() // Ensure the mutex is released when the method finishes.

	// Clear the list of error timestamps and reset the blockedUntil time.
	w.errorTimestamps = w.errorTimestamps[:0]
	w.blockedUntil = time.Time{} // Reset the blockedUntil time to zero.
}

// Close stops the worker by closing the stop channel, signaling all goroutines to terminate.
func (w *Worker) Close() {
	close(w.stop) // Close the stop channel to signal worker termination.
}
