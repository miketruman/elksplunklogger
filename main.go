package main

import (
	"bytes"
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var splunkLoggerService = SplunkLoggerService{}

func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(url), nil
}

func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("Request:", r)
		if r.Body != nil {
			body, _ := ioutil.ReadAll(r.Body)
			//		fmt.Println("BODY:", string(body))
			requests := bytes.Split(body, []byte("\n"))
			for _, skod := range requests {
				value, _, _, err := jsonparser.Get(skod, "query", "bool")
				if err == nil {
					value1, _, _, _ := jsonparser.Get(value, "must")
					value2, _, _, _ := jsonparser.Get(value, "should")
					value3, _, _, _ := jsonparser.Get(value, "must_not")
					filter, _, _, _ := jsonparser.Get(value, "filter")
					if len(value1) > 0 && len(value2) > 0 && len(value3) > 0 && len(filter) > 0 {
						fmt.Println("BOOL:", string(value))
						splunkLoggerService.writerbytes(value)
					}
				}
			}
			r.Body = ioutil.NopCloser(bytes.NewReader(body))
		}
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	splunkLoggerService.Init("https://splunk:8088/services/collector", "abcd1234")
	proxy, err := NewProxy("http://elasticsearch:9200")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", ProxyRequestHandler(proxy))
	log.Fatal(http.ListenAndServe(":9200", nil))
}

/*
var splunkLoggerService = SplunkLoggerService{}

func main() {
        splunkLoggerService.Init("https://localhost:8088/services/collector", "abcd1234")

        type User struct {
                Name string
        }
        user := &User{Name: "Frank"}
        for {
                fmt.Println("sending message")
                splunkLoggerService.writer2(user)
                time.Sleep(1 * time.Second)
        }
}

*/
