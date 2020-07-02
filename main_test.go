package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

func setup() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client, _ = New(BaseURL(server.URL))

	return func() {
		server.Close()
	}
}

func TestGetKeaHooks(t *testing.T) {
	os.Setenv("KEA_LEASE4_ADDRESS", "192.168.100.100")
	os.Setenv("KEA_QUERY4_OPTION60", "idrac")
	os.Setenv("KEA_FAKE_ALLOCATION", "0")

	got := GetKeaHooks()
	want := &KeaHook{
		keaLease4Address:   "192.168.100.100",
		keaQuery4Options60: "idrac",
	}

	if got.keaLease4Address != want.keaLease4Address {
		t.Errorf("expected %s but got %s", want.keaLease4Address, got.keaLease4Address)
	}
}

func TestPayload(t *testing.T) {
	os.Setenv("KEA_LEASE4_ADDRESS", "192.168.100.200")
	os.Setenv("KEA_QUERY4_OPTION60", "idrac")
	os.Setenv("KEA_FAKE_ALLOCATION", "false")

	s := PostDhcpEvent{}
	s.Host.IPAddress = "192.168.100.200"
	s.Host.Vclass = "idrac"

	keaHooks := GetKeaHooks()
	got := Payload(keaHooks.keaLease4Address, keaHooks.keaQuery4Options60)
	want := s

	if got.Host.IPAddress != want.Host.IPAddress {
		t.Errorf("expected %s but got %s", want.Host.IPAddress, got.Host.IPAddress)
	}
}
func TestGetEnv(t *testing.T) {
	os.Setenv("FOO", "bar")

	got := GetEnv("FOO", "")
	want := "bar"

	if got != want {
		t.Errorf("expected %s, got %s", want, got)
	}

	os.Setenv("FOOZ", "")

	got = GetEnv("FOOZ", "")
	want = ""

	if got != want {
		t.Errorf("expected %s, got %s", want, got)
	}

}
