// File Name: server.go
//
// This is a simple web service for detecting a user's remote IP address and geolocation.
// To run:  go build -o checkip; sudo ./checkip

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

// Listen to HTTP port
const port = ":80"

// Error code injection for unit test
const (
	HttpNoErr = iota
	HttpGetErr
	HttpReadAllErr
	HttpUnmarshalErr
	HttpReqErr
)

func main() {
	http.HandleFunc("/", handler)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// HTTP request handler
func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	logMsg("Info", fmt.Sprintf("Received request from host: %s, remoteAddr: %s", r.Host, r.RemoteAddr))

	ip := r.RemoteAddr
	if r.Header.Get("X-Forwarded-For") != "" {
		// Get() always returns the first IP which is the client's IP in X-Forwarded-For list
		ip = r.Header.Get("X-Forwarded-For")
	}
	logMsg("Debug", fmt.Sprintf("Client IP is: %s", ip))

	ip = getValidIP(ip)
	if ip == "" {
		// Just return if IP is invalid. Do not fatal the server
		logMsg("Warning", fmt.Sprintf("Invalid IP: %s", ip))
		fmt.Fprintf(w, "<html><head><title></title></head><body>Invalid IP: %s</body></html>", ip)
		return
	}

	timezone, location, errStr := "NA", "NA", ""

	// Use a "hackered" isPrivateIP() from server_util.go
	// Can not use IP.IsPrivate() because it is too new and not on many hosts
	if isPrivateIP(net.ParseIP(ip)) {
		logMsg("Info", fmt.Sprintf("IP %s is a private address", ip))
	} else {
		timezone, location, errStr = getGeolocation(ip, HttpNoErr)
		if errStr != "" {
			logMsg("Warning", errStr)
		}
	}

	var res = generateResponse(ip, timezone, location)
	logMsg("Info", fmt.Sprintf("Return response to client: %s", res))
	fmt.Fprintf(w, "%s", res)
}

// Return a valid IPv4/IPv6 or empty string
// Input IP is either a valid IP or valid IP+Port
func getValidIP(ip string) string {
	if net.ParseIP(ip) == nil {
		// IP is invalid. Check if IP has port in it
		// If IP has "IPv4:Port" or [IPv6]:Port" format, remove Port
		if strings.Contains(ip, "]") {
			split := strings.Split(ip, "]")
			ip = strings.TrimLeft(split[0], "[")
		} else if strings.Count(ip, ":") == 1 {
			split := strings.Split(ip, ":")
			ip = split[0]
		} else {
			return ""
		}
	}

	return ip
}

// Get location and timezone from IP via http://ipapi.co
// errorInsert is injection error for unit testing only.
// Return timezone, location, error string
func getGeolocation(ip string, errorInsert int) (string, string, string) {
	timezone, location := "NA", "NA"

	url := fmt.Sprintf("http://ipapi.co/%s/json", ip)
	resp, err := http.Get(url)
	if err != nil || errorInsert == HttpGetErr {
		return timezone, location, fmt.Sprintf("Failed to call Get: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil || errorInsert == HttpReadAllErr {
		return timezone, location, fmt.Sprintf("Failed to read resp body: %v", err)
	}

	var bd map[string]interface{}
	err = json.Unmarshal([]byte(body), &bd)
	if err != nil || errorInsert == HttpUnmarshalErr {
		return timezone, location, fmt.Sprintf("Failed to unmarshal: %v", err)
	}

	if bd["error"] != nil || errorInsert == HttpReqErr {
		return timezone, location, fmt.Sprintf("Failed to get geolocation: %s", bd["reason"])
	} else {
		timezone = bd["timezone"].(string)
		location = fmt.Sprintf("%s, %s, %s", bd["city"], bd["region_code"], bd["country_code_iso3"])
	}

	return timezone, location, ""
}

// Generate html format response for the client
func generateResponse(ip string, timezone string, location string) string {
	var res bytes.Buffer
	const resHtml = `<html><head><title></title></head><body>{{range .Items}}<div>{{ . }}</div>{{end}}</body></html>`
	errHtml := fmt.Sprintf("<html><head><title></title></head><body><div>Current IP Address: %s</div><div>Time Zone: NA</div><div>Location: NA</div></body></html>", ip)

	data := struct {
		Items []string
	}{
		Items: []string{
			"Current IP Address: " + ip,
			"Time Zone: " + timezone,
			"Location: " + location,
		},
	}

	parsedTemplate, err := template.New("webpage").Parse(resHtml)
	if err != nil {
		logMsg("Error", fmt.Sprintf("Failed to parse template: %v", err))
		return errHtml
	}

	err = parsedTemplate.Execute(&res, data)
	if err != nil {
		logMsg("Error", fmt.Sprintf("Failed to execute html template: %v", err))
		return errHtml
	}

	return res.String()
}

// Log message with different level keyword
// This function should be improved to log messages with real level
func logMsg(level string, msg string) {
	s := fmt.Sprintf("%s - %s", level, msg)
	log.Println(s)
}
