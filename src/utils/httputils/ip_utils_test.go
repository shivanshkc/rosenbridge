package httputils_test

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shivanshkc/rosenbridge/src/utils/httputils"
)

const (
	mockIP1 = "127.0.0.1"
	mockIP2 = "127.0.0.2"
	mockIP3 = "127.0.0.3"
)

// TestGetOwnIP tests if the GetOwnIP function returns a valid IP without failing.
func TestGetOwnIP(t *testing.T) {
	ip, err := httputils.GetOwnIP()
	if err != nil {
		t.Errorf("expected error: %+v, but got: %+v", nil, err)
		return
	}

	// Validating the received IP address.
	if net.ParseIP(ip) == nil {
		t.Errorf("expected valid ip address, but got invalid")
		return
	}
}

// TestGetClientIP_UsingXRealIP tests if the GetClientIP function prioritizes the "x-real-ip" header above others.
func TestGetClientIP_UsingXRealIP(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("x-real-ip", mockIP1)
	request.Header.Set("x-forwarded-for", mockIP2)
	request.RemoteAddr = mockIP3

	ip, err := httputils.GetClientIP(request)
	if err != nil {
		t.Errorf("expected error: %+v, but got: %+v", nil, err)
		return
	}

	// The IP addresses should match.
	if ip != mockIP1 {
		t.Errorf("expected returned ip to be: %s, but got: %s", mockIP1, ip)
		return
	}
}

// TestGetClientIP_UsingXForwardedFor tests if the GetClientIP function prioritizes the "x-forwarded-for" header
// when "x-real-ip" is not provided.
func TestGetClientIP_UsingXForwardedFor(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("x-forwarded-for", mockIP2)
	request.RemoteAddr = mockIP3

	ip, err := httputils.GetClientIP(request)
	if err != nil {
		t.Errorf("expected error: %+v, but got: %+v", nil, err)
		return
	}

	// The IP addresses should match.
	if ip != mockIP2 {
		t.Errorf("expected returned ip to be: %s, but got: %s", mockIP2, ip)
		return
	}
}

// TestGetClientIP_UsingRemoteAddr tests if the GetClientIP function uses request.RemoteAddr property when neither
// "x-real-ip" nor "x-forwarded-for" are provided.
func TestGetClientIP_UsingRemoteAddr(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.RemoteAddr = mockIP3 + ":80" // Port is required in RemoteAddr.

	ip, err := httputils.GetClientIP(request)
	if err != nil {
		t.Errorf("expected error: %+v, but got: %+v", nil, err)
		return
	}

	// The IP addresses should match.
	if ip != mockIP3 {
		t.Errorf("expected returned ip to be: %s, but got: %s", mockIP3, ip)
		return
	}
}

// TestGetClientIP_NoMethod tests if the GetClientIP function returns an error when all methods to retrieve
// IP address fail.
func TestGetClientIP_NoMethod(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.RemoteAddr = ""

	if _, err := httputils.GetClientIP(request); err == nil {
		t.Errorf("expected error but got nil.")
		return
	}
}
