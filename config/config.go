package config

import (
	"github.com/fraenkel/candiedyaml"
	"io/ioutil"
)

// Clustered version of nats supported
type NatsConfig struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

type Config struct {
	Nats             []NatsConfig `yaml:"nats"`
	Uris             []UriConfig  `yaml:"uris"`
	LoggingLevel     string       `yaml:"logging_level"`
	RegisterInterval int          `yaml:"register_message_interval"`
	GoMaxProcs       int          `yaml:"max_go_procs"`
	Port             uint16       `yaml:"port"`
}

// Multiple url are possible
type UriConfig struct {
	Uri string `yaml:"uri"`
}

var defaultNatsConfig = NatsConfig{
	Host: "localhost",
	Port: 4222,
	User: "",
	Pass: "",
}

var defaultUriConfig = UriConfig{
	Uri: "http://localhost",
}

var defaultConfig = Config{
	Nats:             []NatsConfig{defaultNatsConfig},
	Uris:             []UriConfig{defaultUriConfig},
	LoggingLevel:     "info",
	RegisterInterval: 5,
	GoMaxProcs:       0,
	Port:             80,
}

func (c *Config) Initialize(configYAML []byte) error {
	c.Nats = []NatsConfig{}
	c.Uris = []UriConfig{}
	return candiedyaml.Unmarshal(configYAML, &c)
}

func DefaultConfig() *Config {
	c := defaultConfig
	return &c
}

func InitConfigFromFile(path string) *Config {
	var c *Config = DefaultConfig()
	var e error

	b, e := ioutil.ReadFile(path)
	if e != nil {
		panic(e.Error())
	}

	e = c.Initialize(b)
	if e != nil {
		panic(e.Error())
	}

	return c
}
