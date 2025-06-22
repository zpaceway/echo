package main

import (
	"net/http/httptest"
	"testing"
)

func TestGetClientIpFromRequest(t *testing.T) {
	tests := []struct {
		name        string
		headers     map[string]string
		remoteAddr  string
		expectedIP  string
		description string
	}{
		{
			name: "X-Original-Forwarded-For header present",
			headers: map[string]string{
				"X-Original-Forwarded-For": "192.168.1.100",
			},
			remoteAddr:  "10.0.0.1:12345",
			expectedIP:  "192.168.1.100",
			description: "Should prioritize X-Original-Forwarded-For header",
		},
		{
			name: "X-Real-IP header present (no X-Original-Forwarded-For)",
			headers: map[string]string{
				"X-Real-IP": "203.0.113.5",
			},
			remoteAddr:  "10.0.0.1:12345",
			expectedIP:  "203.0.113.5",
			description: "Should use X-Real-IP when X-Original-Forwarded-For is not present",
		},
		{
			name: "X-Forwarded-For header present (no higher priority headers)",
			headers: map[string]string{
				"X-Forwarded-For": "198.51.100.10",
			},
			remoteAddr:  "10.0.0.1:12345",
			expectedIP:  "198.51.100.10",
			description: "Should use X-Forwarded-For when higher priority headers are not present",
		},
		{
			name: "Multiple IPs in X-Forwarded-For (comma separated)",
			headers: map[string]string{
				"X-Forwarded-For": "203.0.113.1, 198.51.100.2, 192.168.1.50",
			},
			remoteAddr:  "10.0.0.1:12345",
			expectedIP:  "192.168.1.50",
			description: "Should return the last IP in comma-separated X-Forwarded-For list",
		},
		{
			name: "X-Forwarded-For with spaces around IPs",
			headers: map[string]string{
				"X-Forwarded-For": " 203.0.113.1 , 198.51.100.2 , 192.168.1.75 ",
			},
			remoteAddr:  "10.0.0.1:12345",
			expectedIP:  "192.168.1.75",
			description: "Should trim spaces from IPs in X-Forwarded-For",
		},
		{
			name:        "No headers present - fallback to RemoteAddr",
			headers:     map[string]string{},
			remoteAddr:  "172.16.0.100:8080",
			expectedIP:  "172.16.0.100",
			description: "Should extract IP from RemoteAddr when no headers are present",
		},
		{
			name: "Header priority test - X-Original-Forwarded-For wins",
			headers: map[string]string{
				"X-Original-Forwarded-For": "192.168.1.200",
				"X-Real-IP":                "203.0.113.15",
				"X-Forwarded-For":          "198.51.100.25",
			},
			remoteAddr:  "10.0.0.1:12345",
			expectedIP:  "192.168.1.200",
			description: "Should prioritize X-Original-Forwarded-For over other headers",
		},
		{
			name: "Header priority test - X-Real-IP wins over X-Forwarded-For",
			headers: map[string]string{
				"X-Real-IP":       "203.0.113.25",
				"X-Forwarded-For": "198.51.100.35",
			},
			remoteAddr:  "10.0.0.1:12345",
			expectedIP:  "203.0.113.25",
			description: "Should prioritize X-Real-IP over X-Forwarded-For",
		},
		{
			name: "Empty header values should be skipped",
			headers: map[string]string{
				"X-Original-Forwarded-For": "",
				"X-Real-IP":                "",
				"X-Forwarded-For":          "198.51.100.45",
			},
			remoteAddr:  "10.0.0.1:12345",
			expectedIP:  "198.51.100.45",
			description: "Should skip empty headers and use the first non-empty one",
		},
		{
			name:        "IPv6 address in RemoteAddr",
			headers:     map[string]string{},
			remoteAddr:  "[2001:db8::1]:8080",
			expectedIP:  "2001:db8::1",
			description: "Should handle IPv6 addresses in RemoteAddr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = tt.remoteAddr

			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			result := GetClientIpFromRequest(req)

			if result != tt.expectedIP {
				t.Errorf("%s: expected %s, got %s", tt.description, tt.expectedIP, result)
			}
		})
	}
}
