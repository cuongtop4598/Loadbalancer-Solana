package core

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/davecgh/go-spew/spew"

	"go.uber.org/zap"
)

func StartProxy(port int, nodes []Node, currentNodeId *int) {
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]interface{})
		data["nodes"] = nodes
		data["current"] = nodes[*currentNodeId].Endpoint
		js, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	})
	currentNodeUrl := nodes[*currentNodeId].Endpoint
	url, err := url.Parse(currentNodeUrl)
	if err != nil {
		panic(err)
	}
	spew.Dump(url)
	proxy := &httputil.ReverseProxy{
		Transport: roundTripper(rt),
		Director: func(req *http.Request) {
			req.URL.Scheme = url.Scheme
			req.Host = url.Host
			req.URL.Path = url.Path
			req.URL.Host = url.Host
		},
	}
	Log.Info("Starting proxy", zap.Any("Port", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), proxy))
}

func rt(req *http.Request) (*http.Response, error) {
	log.Printf("request received. url=%s", req.URL)
	defer log.Printf("request complete. url=%s", req.URL)
	return http.DefaultTransport.RoundTrip(req)
}

type roundTripper func(*http.Request) (*http.Response, error)

func (f roundTripper) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }
