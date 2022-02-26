// File Name: server_sys_test.go
//
// This is test suite which tests the web service for IP checking.
// These are system tests which send requests to sever and get response back,
// so they need the server be running.
//
// To run:  go build -o checkip; sudo ./checkip;  go test -v

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

const (
	serverURL = "http://localhost"
	serverIP  = "127.0.0.1"
)

// Test X-Forwarded-For.
// Server should return the IP set from X-Forwarded-For header
func TestXForwardedForIPV6(t *testing.T) {
	expectedResult := "<html><head><title></title></head><body><div>Current IP Address: 2001:4860:4860::8888</div><div>Time Zone: America/Los_Angeles</div><div>Location: Mountain View, CA, USA</div></body></html>"

	client := &http.Client{}
	req, err := http.NewRequest("Get", serverURL, nil)
	if err != nil {
		t.Fatalf("Connect server: %v\n", err)
	}

	req.Header = http.Header{
		"X-Forwarded-For": []string{"2001:4860:4860::8888", "2001:4860:4860::1234", "2001:4860:4860::3456"},
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do call failed: %v\n", err)
	}
	defer resp.Body.Close()

	bd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Fetch body failed: %v\n", err)
	}

	bds := fmt.Sprintf("%s", bd)
	if bds == "" || bds != expectedResult {
		t.Fatalf("Result does not match. Expected: %s. Got: %s\n", expectedResult, bd)
	}
}

// Test the normal response from server
func TestMain(t *testing.T) {
	expectedResult := fmt.Sprintf("<html><head><title></title></head><body><div>Current IP Address: %s</div><div>Time Zone: NA</div><div>Location: NA</div></body></html>", serverIP)

	resp, err := http.Get(serverURL)
	if err != nil {
		t.Fatalf("fretch IP failed: %v\n", err)
	}
	defer resp.Body.Close()

	bd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Fetch body failed: %v\n", err)
	}

	bds := fmt.Sprintf("%s", bd)
	if bds != expectedResult {
		t.Fatalf("Result does not match. Expected: %s. Got: %s\n", expectedResult, bd)
	}
}
