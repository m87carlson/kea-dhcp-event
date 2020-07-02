package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

type KeaHook struct {
	keaLease4Address   string
	keaQuery4Options60 string
}

type PostDhcpEvent struct {
	Host struct {
		IPAddress string `json:"ip_address"`
		Vclass    string `json:"vclass"`
	} `json:"host"`
}

type Option func(*Client) error

func main() {
	// if there are no args, just exit
	if len(os.Args) <= 1 {
		os.Exit(0)
	}

	//setup logging
	logfile := "/var/log/kea/dhcp_event.log"

	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	if err != nil {
		// Cannot open log file. Logging to stderr
		fmt.Println(err)
	} else {
		log.SetOutput(f)
	}
	defer f.Close()

	log.SetLevel(log.InfoLevel)

	debug, _ := strconv.ParseBool(os.Getenv("KEA_HOOK_DEBUG"))

	if debug {
		log.SetLevel(log.DebugLevel)
		log.Debugf("args: %v", os.Args)
		for _, msg := range os.Environ() {
			log.Debugf("env: %s", msg)
		}
	}

	// Unless we have lease4_select or lease4_renew, just exit and no nothing
	if (os.Args[1] == "lease4_select" && os.Getenv("KEA_FAKE_ALLOCATION") == "0") || os.Args[1] == "lease4_renew" {
		keaHooks := GetKeaHooks()

		p := Payload(keaHooks.keaLease4Address, keaHooks.keaQuery4Options60)
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(p)

		client, err := New()
		if err != nil {
			log.Error(err)
		}

		endPoint := client.baseURL + "/api/v1/hosts/discover"

		req, err := http.NewRequest("POST", endPoint, b)
		if err != nil {
			log.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", "Token token="+client.apiKey)

		c := http.Client{}
		resp, err := c.Do(req)
		if err != nil {
			log.Error(err)
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		if debug {
			log.Debugf("Server Response: %s", string(body))
			log.Debugf("Payload Sent: %v", p)
			log.Debugf("Response Code: %v", resp.StatusCode)
		}

		if resp.StatusCode != 202 {
			log.Errorf("Server returned non-200 status: %v : %s", resp.StatusCode, string(body))
		}
	} else {
		// if args[1] is anything else besides lease4_renew or lease4_select, do nothing
		os.Exit(0)
	}

}

func (c *Client) ParseOptions(opts ...Option) error {
	for _, option := range opts {
		err := option(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func BaseURL(baseURL string) Option {
	return func(c *Client) error {
		c.baseURL = baseURL
		return nil
	}
}

func New(opts ...Option) (*Client, error) {
	client := &Client{
		baseURL: os.Getenv("CMDB_URL"),
		apiKey:  os.Getenv("KEA_CMDB_TOKEN"),
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}

	if err := client.ParseOptions(opts...); err != nil {
		log.Error(err)
		return nil, err
	}

	return client, nil
}

func GetKeaHooks() *KeaHook {
	debug, _ := strconv.ParseBool(GetEnv("KEA_HOOK_DEBUG", "none"))

	if debug {
		log.Debug(GetEnv("KEA_LEASE4_ADDRESS", "none"))
		log.Debug(GetEnv("KEA_QUERY4_OPTION60", "none"))
	}

	return &KeaHook{
		keaLease4Address:   GetEnv("KEA_LEASE4_ADDRESS", "none"),
		keaQuery4Options60: GetEnv("KEA_QUERY4_OPTION60", "none"),
	}
}

func Payload(ip, vclass string) *PostDhcpEvent {
	s := &PostDhcpEvent{}
	s.Host.Vclass = vclass
	s.Host.IPAddress = ip

	return s
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
