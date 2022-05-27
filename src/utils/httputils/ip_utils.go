package httputils

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
)

// These are here to make sure that the IP is calculated only once.
var ownIPOnce = &sync.Once{}
var ownIPSingleton string

// GetOwnIP gets the IP address of the host machine.
func GetOwnIP() (string, error) {
	var err error
	ownIPOnce.Do(func() {
		ownIPSingleton, err = getOwnIP()
	})

	if err != nil {
		return "", err
	}
	return ownIPSingleton, nil
}

// getOwnIP gets the IP address of the host machine.
func getOwnIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", fmt.Errorf("error in net.Dial call: %w", err)
	}
	defer func() { _ = conn.Close() }()

	// This address is of the form <ip>:<port>
	address := conn.LocalAddr().String()

	// Removing the port from the address.
	ip, _, err := net.SplitHostPort(address)
	if err != nil {
		return "", fmt.Errorf("error in net.SplitHostPort call: %w", err)
	}
	return ip, nil
}

// GetClientIP extracts the client IP Address from the given HTTP request.
func GetClientIP(req *http.Request) (string, error) {
	// Using x-real-ip header.
	ip := req.Header.Get("x-real-ip")
	if parsedIP := net.ParseIP(ip); parsedIP != nil {
		return parsedIP.String(), nil
	}

	// Using x-forwarded-for header.
	ips := req.Header.Get("x-forwarded-for")
	ipArr := strings.Split(ips, ",")
	if len(ipArr) > 0 {
		if parsedIP := net.ParseIP(ipArr[0]); parsedIP != nil {
			return parsedIP.String(), nil
		}
	}

	// Using RemoteAddr property.
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return "", fmt.Errorf("error in net.SplitHostPort call: %w", err)
	}

	if parsedIP := net.ParseIP(ip); parsedIP != nil {
		return parsedIP.String(), nil
	}

	return "", fmt.Errorf("no method worked")
}
