package callback

import "time"

// DeliveryMode defines the method for delivering messages to clients.
// It can be used to select a notification delivery strategy,
// whether it's sending a single message to one client or broadcasting to all clients.
type DeliveryMode string

var (
	// RoundRobin specifies a message delivery mode to one of the connected clients.
	// Messages are sent sequentially to each client using the Round Robin algorithm.
	// This means that messages are distributed in a cyclic manner, ensuring
	// each client receives a message in turn, rather than all at once.
	RoundRobin DeliveryMode = "round_robin"

	// Broadcast specifies a message delivery mode that sends to all connected clients at once.
	// Each notification is broadcast immediately to all clients connected to the server.
	// This mode ensures that every client receives the same message at the same time.
	Broadcast DeliveryMode = "broadcast"
)

// Transport defines the transport protocol used for message delivery.
type Transport string

var (
	// REST is the transport protocol that uses RESTful API for message delivery.
	// This is typically over HTTP and suitable for stateless communication.
	REST Transport = "REST"

	// QUIC is the transport protocol that uses QUIC (Quick UDP Internet Connections) for message delivery.
	// It provides low-latency, secure transport and is typically faster than traditional HTTP/HTTPS protocols.
	QUIC Transport = "QUIC"
)

// RetryMode defines how retry logic is handled when sending messages to endpoints.
type RetryMode string

var (
	// Repeat mode retries the failed request on the same endpoint where the error occurred.
	Repeat RetryMode = "repeat"

	// Next mode retries the request on the next available endpoint after a failed attempt.
	Next RetryMode = "next"
)

// Options contains configuration options for the message delivery system.
type Options struct {

	// Transport defines the transport protocol used for message delivery.
	Transport Transport

	// DeliveryMode defines the method for delivering messages to clients.
	// It can be used to select a notification delivery strategy,
	// whether it's sending a single message to one client or broadcasting to all clients.
	DeliveryMode DeliveryMode

	// RetryMode defines the behavior when retrying failed message delivery attempts.
	RetryMode RetryMode

	// EndPoints specifies the IP addresses or addresses of endpoints for message delivery.
	// This can be modified in real-time based on server settings, allowing dynamic control over the delivery targets.
	EndPoints []string

	// RetryLimit is the maximum number of retry attempts for message delivery.
	// If the server fails to deliver a message within the set limit, it will temporarily stop sending messages to this endpoint.
	// Default value: 5
	RetryLimit int

	// RetryTimeout is the duration for which the server will pause sending messages to an endpoint
	// after exceeding the retry limit. This helps to prevent overloading the endpoint with repeated failed attempts.
	// Default value: time.Second * 5
	RetryTimeout time.Duration

	// RetryWindow is the period of time during which retries will be counted toward the RetryLimit.
	// This window ensures that the RetryLimit is not exceeded within a short burst of attempts.
	RetryWindow time.Duration
}

// defaultOptions initializes default values for Options fields that are not set.
// It checks each field in the Options struct, and if any field has a zero value,
// it assigns a default value for that field.
func defaultOptions(opt *Options) *Options {

	// Set default transport to REST if none is specified
	if opt.Transport == "" {
		opt.Transport = REST
	}

	// Set default delivery mode to RoundRobin if none is specified
	if opt.DeliveryMode == "" {
		opt.DeliveryMode = RoundRobin
	}

	// Set default retry mode to Next if none is specified
	if opt.RetryMode == "" {
		opt.RetryMode = Next
	}

	// Set default retry limit to 5 if none is specified
	if opt.RetryLimit == 0 {
		opt.RetryLimit = 5
	}

	// Set default retry timeout to 5 seconds if none is specified
	if opt.RetryTimeout == 0 {
		opt.RetryTimeout = time.Second * 5
	}

	// Set default retry window to 3 seconds if none is specified
	if opt.RetryWindow == 0 {
		opt.RetryWindow = time.Second * 3
	}

	return opt
}
