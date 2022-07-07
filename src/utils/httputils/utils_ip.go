package httputils

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
)

// These are here to make sure that the IP is calculated only once.
var (
	ownIPOnce      = &sync.Once{}
	ownIPSingleton string
)

// GetOwnIP gets the IP address of the host machine.
func GetOwnIP() (string, error) {
	var err error
	// This call executes only once.
	ownIPOnce.Do(func() {
		ownIPSingleton, err = getOwnIP()
	})
	// If there's an error, we return it along with an empty string.
	if err != nil {
		return "", err
	}
	// Returning valid IP.
	return ownIPSingleton, nil
}

// getOwnIP gets the IP address of the host machine.
func getOwnIP() (string, error) {
	// Creating dummy connection to get our own IP.
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", fmt.Errorf("error in net.Dial call: %w", err)
	}
	// Closing connection upon function return.
	defer func() { _ = conn.Close() }()

	// This address is of the form <ip>:<port>
	address := conn.LocalAddr().String()

	// Removing the port from the address.
	ip, _, err := net.SplitHostPort(address)
	if err != nil {
		return "", fmt.Errorf("error in net.SplitHostPort call: %w", err)
	}
	// Final IP.
	return ip, nil
}

// GetClientIP extracts the client IP Address from the given HTTP request.
func GetClientIP(req *http.Request) (string, error) {
	// Using x-real-ip header.
	ipAddr := req.Header.Get("x-real-ip")
	if parsedIP := net.ParseIP(ipAddr); parsedIP != nil {
		return parsedIP.String(), nil
	}

	// Using x-forwarded-for header.
	ips := req.Header.Get("x-forwarded-for")
	// IP addresses can be comma separated here. We'll use the first one.
	ipArr := strings.Split(ips, ",")
	// Checking if there is at least one IP.
	if len(ipArr) > 0 {
		// Parsing and returning the first IP.
		if parsedIP := net.ParseIP(ipArr[0]); parsedIP != nil {
			return parsedIP.String(), nil
		}
	}

	// Using RemoteAddr property.
	ipAddr, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return "", fmt.Errorf("error in net.SplitHostPort call: %w", err)
	}

	// Parsing and returning the IP.
	if parsedIP := net.ParseIP(ipAddr); parsedIP != nil {
		return parsedIP.String(), nil
	}

	return "", errors.New("no method worked")
}
