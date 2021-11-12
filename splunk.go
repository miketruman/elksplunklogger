package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/ZachtimusPrime/Go-Splunk-HTTP/splunk/v2"
	"net/http"
	"time"
)

type SplunkLoggerService struct {
	writer splunk.Writer
}
	type splunkLoggerStatus struct {
		Event string
		Errors int
		Uptime float64
	}

func (splunkLoggerService *SplunkLoggerService) splunkErrorHandling() {
	bootTime := time.Now()
	errorCount := 0
	for {
		select {
		case message := <-splunkLoggerService.writer.Errors():
			fmt.Println(message)
			errorCount = errorCount + 1 
		case <-time.After(1 * time.Second):
			nowTime := time.Now()
			uptime := nowTime.Sub(bootTime).Seconds()
			user := &splunkLoggerStatus{Event: "Uptime", Uptime: uptime, Errors: errorCount}
			splunkLoggerService.writer2(user)
		}
	}
}
func (splunkLoggerService *SplunkLoggerService) writer2(message interface{}) error {
	b, err := json.Marshal(message)
	splunkLoggerService.writer.Write(b)
	return err
}
func (splunkLoggerService *SplunkLoggerService) writerbytes(message []byte) {
	splunkLoggerService.writer.Write(message)
}

func (splunkLoggerService *SplunkLoggerService) Init(server string, token string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}

	splunkClient := splunk.NewClient(
		httpClient,
		server,
		token,
		"",
		"",
		"",
	)

	writer := splunk.Writer{
		Client:         splunkClient,
		FlushInterval:  1 * time.Second, // How often we'll flush our buffer
		FlushThreshold: 1000000,              // Max messages we'll keep in our buffer, regardless of FlushInterval
		MaxRetries:     10,              // Number of times we'll retry a failed send
	}
	splunkLoggerService.writer = writer
	go splunkLoggerService.splunkErrorHandling()
}

