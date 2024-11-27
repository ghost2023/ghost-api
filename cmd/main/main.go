package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"ghost-api/cmd/main/config"
	"ghost-api/cmd/main/response"

	// "ghost-api/cmd/main/response"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

var conf = config.Config{
	Port:    "7500",
	Latency: 100,
	Jitter:  10,
	Timeout: 1000,
}

type Router struct{}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Request received:", req.URL.Path)
	for i, endpoint := range conf.Endpoints {
		if req.URL.Path == endpoint.Url {

			fmt.Println("Endpoint found:", conf.Endpoints[i].Name)
			w.WriteHeader(conf.Endpoints[i].Response.StatusCode)
			res := response.GenerateResponse(&endpoint.Response.Data)
			json.NewEncoder(w).Encode(res)
		}
	}
}

func main() {
	router := Router{}

	fmt.Printf("Starting server on port %s\n", conf.Port)
	http.ListenAndServe(":"+conf.Port, &router)
}

func init() {
	configFile := flag.String("c", ".ghostrc", "Config file")
	flag.Parse()

	configFileContent, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal(configFileContent, &conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if conf.Endpoints == nil {
		log.Fatal("No endpoints defined in config file")
	}

	config.ApplyDefaults(&conf)
}
