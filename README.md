# CheckIP in Golang

----

This is a simple web service for detecting a user's remote IP address and
geolocation information.

The web service is implemented in Go, with no specific library required. 

----

## Design

The application is designed as a single thread web service which takes the client request
in sequence. After the service receives a HTTP request from port 80, it checks header
X-Forwarded-For and sets the remote IP to this value if this header is set. Otherwise,
http.Request.RemoteAddr is used to set the remote IP.

In order to get geolocation data (location and time zone), the web service makes
HTTP call to http://ipapi.co which should return the geolocation information.  This call
will be made only when the IP is public. In the case that an error is happened, the
call to ipapi.co will fail and the geolocation will return "NA" instead.

The web service will write the IP, location and time zone into html and send back
to the client. If no geolocation information is returned because of the error, the web
service will return IP address back with the location and time zone set to "NA".

The web service is listening on HTTP port 80. The server code can be changed to use
different port. If use port 80, server need to be run in sudo or administrator.

The web service is able to provide the service for both IPv4 and IPv6.

## To compile

The server code is in server.go and server_util.go. Run the following cmd to compile:

    go build -o checkip


## Unit Test

There are two test files:

### server_unit_test.go:

These tests can be executed without server is running. These unit tests verify the
server functions.

### server_sys_test.go:

These system tests are executed when server is running. The tests send HTTP requests to
server and check/verify the returned HTML result.

To run these tests, with/without server running, run the cmd:

        go test -v

All tests should pass:

       me@Ubuntu-VM:~/checkip$ go test -v
       === RUN   TestXForwardedForIPV6
       --- PASS: TestXForwardedForIPV6 (0.50s)
       === RUN   TestMain
       --- PASS: TestMain (0.00s)
       === RUN   TestGetValidIP
       --- PASS: TestGetValidIP (0.00s)
       === RUN   TestGetGeolocation
       --- PASS: TestGetGeolocation (0.86s)
       === RUN   TestGeteGeolocationErrInjection
       --- PASS: TestGeteGeolocationErrInjection (0.67s)
       === RUN   TestGenerateResponse
       --- PASS: TestGenerateResponse (0.00s)
       PASS
       ok      _/home/me/checkip     2.039s


## To start the service

Run the following cmd (e.g., on Linux):

        sudo ./checkip

## To stop the service

Kill the above cmd process.

## GUI

Once the web service is deployed, opens the brower and goes to "http://<server>", it should
display the this information:

      Current IP Address: 172.58.190.212
      Time Zone: America/New_York
      Location: Baltimore, MD, USA

which is an html data:

      <html><head><title></title></head><body><div>Current IP Address: 172.58.190.212</div><div>Time Zone: America/New_York</div><div>Location: Baltimore, MD, USA</div></body></html>

If for some reasons the server can not retrieve the geolocation, it will display the IP with time zone and location set to "NA":
  
      Current IP Address: ::1
      Time Zone: NA
      Location: NA

## Limitation

IP.IsPrivate() is in Golang latest std package, but not on Ubuntu 18.04 VM that I am testing.
So I am using an online source to handle this (in server_util.go).

