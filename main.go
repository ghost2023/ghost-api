package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ghost2023/ghost-api/cmd/ghost-api/config"
	"github.com/ghost2023/ghost-api/cmd/ghost-api/response"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

var conf = config.Config{
	Port:    "7500",
	Latency: 250,
	Jitter:  50,
	Timeout: 1000,
}

type Router struct{}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for i, endpoint := range conf.Endpoints {
		if req.URL.Path == endpoint.Url {
			delay := rand.Intn(endpoint.Jitter*2) - endpoint.Jitter
			time.Sleep(time.Duration(delay+endpoint.Latency) * time.Millisecond)

			fmt.Print("\nRequest received:", req.URL.Path)
			fmt.Println(",", delay+endpoint.Latency, "ms")

			w.WriteHeader(conf.Endpoints[i].Response.StatusCode)
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			res := response.GenerateResponse(&endpoint.Response.Data)
			json.NewEncoder(w).Encode(res)
			return
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
