// File Name: server_unit_test.go
//
// This is test suite which unit tests the web service for IP checking.
// To run:  go test -v

package main

import (
	"strings"
	"testing"
)

// Test get valid IP
func TestGetValidIP(t *testing.T) {
	ips := map[string]string{
		"8.8.8.8":                        "8.8.8.8",
		"8.8.8.8:2000":                   "8.8.8.8",
		"2001:4860:4860::8888":           "2001:4860:4860::8888",
		"[2001:4860:4860::8888]:3000":    "2001:4860:4860::8888",
		"800.8.8.8":                      "",
		"2001:4860:4860::8888:200:55555": "",
	}

	for ip, expected := range ips {
		actual := getValidIP(ip)
		if actual != expected {
			t.Errorf("Failed: check valid IP %s failed. Expected: %s, Actual: %s\n", ip, expected, actual)
		}
	}
}

// Test get geolocation from a IPv4 and IPv6 address
func TestGetGeolocation(t *testing.T) {
	// Test a valid IPv4
	timezone, location, errStr := getGeolocation("8.8.2.8", HttpNoErr)
	if errStr != "" {
		t.Errorf("Fetch ipv4 geolocation failed: %s, %s\n", timezone, location)
	}

	// Test an invalid IP
	timezone, location, errStr = getGeolocation("8.8.2.888", HttpNoErr)
	if errStr == "" {
		t.Errorf("Fetch invalid ipv4 geolocation failed: %s, %s\n", timezone, location)
	}

	// Test a valid IPv6
	timezone, location, errStr = getGeolocation("2001:4860:4860::8888", HttpNoErr)
	if errStr != "" {
		t.Errorf("Fetch ipv6 geolocation failed: %s, %s\n", timezone, location)
	}

	// Test an invalid IPv6
	timezone, location, errStr = getGeolocation("2001:4860:4860::8888888", HttpNoErr)
	if errStr == "" {
		t.Errorf("Fetch invalid ipv6 geolocation failed: %s, %s\n", timezone, location)
	}
}

// Test get geolocation with different injection error code
func TestGeteGeolocationErrInjection(t *testing.T) {
	_, _, errStr := getGeolocation("8.8.2.8", HttpGetErr)
	if !strings.Contains(errStr, "Failed to call Get") {
		t.Errorf("Verification failed in GetGeologation Get call: %s\n", errStr)
	}

	_, _, errStr = getGeolocation("8.8.2.8", HttpReadAllErr)
	if !strings.Contains(errStr, "Failed to read") {
		t.Errorf("Verification failed in GetGeologation read: %s\n", errStr)
	}

	_, _, errStr = getGeolocation("2001:4860:4860::8888", HttpUnmarshalErr)
	if !strings.Contains(errStr, "Failed to unmarshal") {
		t.Errorf("Verification failed in GetGeologation unmarshal: %s\n", errStr)
	}

	_, _, errStr = getGeolocation("2001:4860:4860::8888", HttpReqErr)
	if !strings.Contains(errStr, "Failed to get") {
		t.Errorf("Verification failed in GetGeologation get geolocation: %s\n", errStr)
	}
}

// Test generated html response to client
func TestGenerateResponse(t *testing.T) {
	// Verify normal case
	expectedRes := "<html><head><title></title></head><body><div>Current IP Address: 8.8.2.8</div><div>Time Zone: America/New_York</div><div>Location: Newark, NJ, USA</div></body></html>"
	result := generateResponse("8.8.2.8", "America/New_York", "Newark, NJ, USA")
	if result != expectedRes {
		t.Errorf("Generate html response failed, Expected: %s. Actual: %s\n", expectedRes, result)
	}

	// Verify error case
	expectedRes = "<html><head><title></title></head><body><div>Current IP Address: 2001:4860:4860::8888</div><div>Time Zone: NA</div><div>Location: NA</div></body></html>"
	result = generateResponse("2001:4860:4860::8888", "NA", "NA")
	if result != expectedRes {
		t.Errorf("Generate html response failed, Expected: %s. Actual: %s\n", expectedRes, result)
	}
}
