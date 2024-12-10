package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
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

var configFile *string

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for i, endpoint := range conf.Endpoints {
		if req.URL.Path == endpoint.Url {
			delay := rand.Intn(endpoint.Jitter*2) - endpoint.Jitter
			time.Sleep(time.Duration(delay+endpoint.Latency) * time.Millisecond)

			fmt.Print("Request received:", req.URL.Path)
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
	w.WriteHeader(404)
	w.Write([]byte("404 Not Found"))
	return
}

func main() {
	router := Router{}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// If the file was modified, reload the configuration
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("Config file changed:", event.Name)
					err := reloadConfig(*configFile)
					if err != nil {
						log.Printf("Failed to reload config: %v", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Watcher error: %v", err)
			}
		}
	}()

	err = watcher.Add(*configFile)
	if err != nil {
		log.Fatalf("Failed to watch config file: %v", err)
	}

	fmt.Printf("Starting server on port %s\n", conf.Port)
	err = http.ListenAndServe(":"+conf.Port, &router)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func reloadConfig(configFile string) error {
	configFileContent, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	err = yaml.Unmarshal(configFileContent, &conf)
	if err != nil {
		return fmt.Errorf("error unmarshalling config file: %v", err)
	}

	// Apply defaults to ensure valid configuration
	config.ApplyDefaults(&conf)
	return nil
}

func init() {
	configFile = flag.String("c", ".ghostrc", "Config file")
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
