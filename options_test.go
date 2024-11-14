package callback

import (
	"testing"
	"time"
)

// Helper function to compare slices of strings
func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestDefaultOptions tests the defaultOptions function.
func TestDefaultOptions(t *testing.T) {
	tests := []struct {
		name     string
		input    *Options
		expected *Options
	}{
		{
			name:  "All fields empty",
			input: &Options{},
			expected: &Options{
				Transport:    REST,
				DeliveryMode: RoundRobin,
				RetryMode:    Next,
				EndPoints:    nil,
				RetryLimit:   5,
				RetryTimeout: time.Second * 5,
				RetryWindow:  time.Second * 3,
			},
		},
		{
			name: "Transport set to QUIC",
			input: &Options{
				Transport: QUIC,
			},
			expected: &Options{
				Transport:    QUIC,
				DeliveryMode: RoundRobin,
				RetryMode:    Next,
				EndPoints:    nil,
				RetryLimit:   5,
				RetryTimeout: time.Second * 5,
				RetryWindow:  time.Second * 3,
			},
		},
		{
			name: "DeliveryMode set to Broadcast",
			input: &Options{
				DeliveryMode: Broadcast,
			},
			expected: &Options{
				Transport:    REST,
				DeliveryMode: Broadcast,
				RetryMode:    Next,
				EndPoints:    nil,
				RetryLimit:   5,
				RetryTimeout: time.Second * 5,
				RetryWindow:  time.Second * 3,
			},
		},
		{
			name: "RetryMode set to Repeat",
			input: &Options{
				RetryMode: Repeat,
			},
			expected: &Options{
				Transport:    REST,
				DeliveryMode: RoundRobin,
				RetryMode:    Repeat,
				EndPoints:    nil,
				RetryLimit:   5,
				RetryTimeout: time.Second * 5,
				RetryWindow:  time.Second * 3,
			},
		},
		{
			name: "RetryLimit set to 10",
			input: &Options{
				RetryLimit: 10,
			},
			expected: &Options{
				Transport:    REST,
				DeliveryMode: RoundRobin,
				RetryMode:    Next,
				EndPoints:    nil,
				RetryLimit:   10,
				RetryTimeout: time.Second * 5,
				RetryWindow:  time.Second * 3,
			},
		},
		{
			name: "RetryTimeout set to 10 seconds",
			input: &Options{
				RetryTimeout: time.Second * 10,
			},
			expected: &Options{
				Transport:    REST,
				DeliveryMode: RoundRobin,
				RetryMode:    Next,
				EndPoints:    nil,
				RetryLimit:   5,
				RetryTimeout: time.Second * 10,
				RetryWindow:  time.Second * 3,
			},
		},
		{
			name: "RetryWindow set to 2 seconds",
			input: &Options{
				RetryWindow: time.Second * 2,
			},
			expected: &Options{
				Transport:    REST,
				DeliveryMode: RoundRobin,
				RetryMode:    Next,
				EndPoints:    nil,
				RetryLimit:   5,
				RetryTimeout: time.Second * 5,
				RetryWindow:  time.Second * 2,
			},
		},
		{
			name: "EndPoints set",
			input: &Options{
				EndPoints: []string{"192.168.1.1", "192.168.1.2"},
			},
			expected: &Options{
				Transport:    REST,
				DeliveryMode: RoundRobin,
				RetryMode:    Next,
				EndPoints:    []string{"192.168.1.1", "192.168.1.2"},
				RetryLimit:   5,
				RetryTimeout: time.Second * 5,
				RetryWindow:  time.Second * 3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := defaultOptions(tt.input)

			// Compare Transport, DeliveryMode, RetryMode, RetryLimit, RetryTimeout, RetryWindow
			if got.Transport != tt.expected.Transport {
				t.Errorf("defaultOptions() Transport = %v, want %v", got.Transport, tt.expected.Transport)
			}
			if got.DeliveryMode != tt.expected.DeliveryMode {
				t.Errorf("defaultOptions() DeliveryMode = %v, want %v", got.DeliveryMode, tt.expected.DeliveryMode)
			}
			if got.RetryMode != tt.expected.RetryMode {
				t.Errorf("defaultOptions() RetryMode = %v, want %v", got.RetryMode, tt.expected.RetryMode)
			}
			if got.RetryLimit != tt.expected.RetryLimit {
				t.Errorf("defaultOptions() RetryLimit = %v, want %v", got.RetryLimit, tt.expected.RetryLimit)
			}
			if got.RetryTimeout != tt.expected.RetryTimeout {
				t.Errorf("defaultOptions() RetryTimeout = %v, want %v", got.RetryTimeout, tt.expected.RetryTimeout)
			}
			if got.RetryWindow != tt.expected.RetryWindow {
				t.Errorf("defaultOptions() RetryWindow = %v, want %v", got.RetryWindow, tt.expected.RetryWindow)
			}

			// Compare EndPoints using slicesEqual helper
			if !slicesEqual(got.EndPoints, tt.expected.EndPoints) {
				t.Errorf("defaultOptions() EndPoints = %v, want %v", got.EndPoints, tt.expected.EndPoints)
			}
		})
	}
}
