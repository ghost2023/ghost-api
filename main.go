package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type Response struct {
	StatusCode int         `yaml:"status_code"`
	DataType   string      `yaml:"data_type"`
	Data       interface{} `yaml:"data"`
}

type EndPoint struct {
	Name        string   `yaml:"name"`
	Url         string   `yaml:"url"`
	Latency     int      `yaml:"latency"`
	Jitter      int      `yaml:"jitter"`
	Timeout     int      `yaml:"timeout"`
	SuccessRate float64  `yaml:"success_rate"`
	Response    Response `yaml:"response"`
}

type Config struct {
	Port               string     `yaml:"port"`
	DefaultLatency     int        `yaml:"default_latency"`
	DefaultJitter      int        `yaml:"default_jitter"`
	DefaultTimeout     int        `yaml:"default_timeout"`
	DefaultSuccessRate float64    `yaml:"default_success_rate"`
	Endpoints          []EndPoint `yaml:"endpoints"`
}

var config Config = Config{
	Port:               "7500",
	DefaultLatency:     100,
	DefaultJitter:      10,
	DefaultTimeout:     1000,
	DefaultSuccessRate: 0.5,
}

type Router struct{}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Request received:", req.URL.Path)
	fmt.Printf("Endpoint found: %#v\n", req.URL)
	fmt.Println("")
	for i := range config.Endpoints {
		if req.URL.Path == config.Endpoints[i].Url {

			fmt.Println("Endpoint found:", config.Endpoints[i].Name)
			fmt.Println("Response:", config.Endpoints[i].Response)
			w.WriteHeader(config.Endpoints[i].Response.StatusCode)
			w.Header().Set("Content-Type", config.Endpoints[i].Response.DataType)
			w.Write([]byte("hello world"))
		}
	}
}

func main() {

	router := Router{}

	fmt.Printf("Starting server on port %s\n", config.Port)
	http.ListenAndServe(":"+config.Port, &router)

}

func init() {
	configFile := flag.String("c", ".ghostrc", "Config file")
	flag.Parse()

	configFileContent, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal(configFileContent, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if config.Endpoints == nil {
		log.Fatal("No endpoints defined in config file")
	}

	applyDefaults(&config)
}

func applyDefaults(config *Config) {
	for i := range config.Endpoints {
		if config.Endpoints[i].Latency == 0 {
			config.Endpoints[i].Latency = config.DefaultLatency
		}
		if config.Endpoints[i].Jitter == 0 {
			config.Endpoints[i].Jitter = config.DefaultJitter
		}
		if config.Endpoints[i].Timeout == 0 {
			config.Endpoints[i].Timeout = config.DefaultTimeout
		}
		if config.Endpoints[i].SuccessRate == 0 {
			config.Endpoints[i].SuccessRate = config.DefaultSuccessRate
		}
	}
}
