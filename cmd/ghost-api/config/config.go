package config

type NodeType string

const (
	TypeDate   NodeType = "date"
	TypeString NodeType = "string"
	TypeNumber NodeType = "number"
	TypeObject NodeType = "object"
	TypeArray  NodeType = "array"
	TypeEnum   NodeType = "enum"
)

type Node struct {
	Type     NodeType         `json:"type"`
	Value    interface{}      `json:"value,omitempty"`
	Range    []string         `json:"range,omitempty"`
	Metadata string           `json:"metadata,omitempty"`
	Fields   map[string]*Node `json:"fields,omitempty"`
	Items    *Node            `json:"items,omitempty"`
}

type Response struct {
	StatusCode int    `yaml:"status_code"`
	DataType   string `yaml:"data_type"`
	Data       Node   `yaml:"data"`
}

type EndPoint struct {
	Name     string   `yaml:"name"`
	Url      string   `yaml:"url"`
	Latency  int      `yaml:"latency"`
	Jitter   int      `yaml:"jitter"`
	Timeout  int      `yaml:"timeout"`
	Response Response `yaml:"response"`
}

type Config struct {
	Port      string     `yaml:"port"`
	Latency   int        `yaml:"latency"`
	Jitter    int        `yaml:"jitter"`
	Timeout   int        `yaml:"timeout"`
	Endpoints []EndPoint `yaml:"endpoints"`
}

func ApplyDefaults(config *Config) {
	for i := range config.Endpoints {
		if config.Endpoints[i].Latency == 0 {
			config.Endpoints[i].Latency = config.Latency
		}
		if config.Endpoints[i].Jitter == 0 {
			config.Endpoints[i].Jitter = config.Jitter
		}
		if config.Endpoints[i].Timeout == 0 {
			config.Endpoints[i].Timeout = config.Timeout
		}
	}
}
