package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

// Config is Api server configuration
type Config struct {
	Listen                   string       `json:"listen"`
	AccessControlAllowOrigin string       `json:"accessControllAllowOrigin"`
	URIPrefix                string       `json:"uriPrefix"`
	HTTPServiceTimeout       *HTTPTimeout `json:"httpServiceTimeout,omitempty"`
	HTTPSServiceTimeout      *HTTPTimeout `json:"httpsServiceTimeout,omitempty"`
	AppName                  string       `json:"appName"`
	PostgresSQL              *SQLConn     `json:"sqlConn"`
}

type SQLConn struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Server   string `json:"server"`
	Database string `json:"database"`
}

// type AuthAPIConfig struct {
// 	URLPrefix  string `json:"urlPrefix"`
// 	HostHeader string `json:"hostHeader"`
// }

// HTTPTimeout is timeout configuration
type HTTPTimeout struct {
	ReadTimeout  int `json:"readTimeout,omitempty"`
	WriteTimeout int `json:"writeTimeout,omitempty"`
}

// APISrvDynConfigurator is interface for API configuration loader
type APISrvDynConfigurator interface {
	GetConf() (*Config, error)
}

// APISrvDynConf is APISrvDynConfigurator implementation using etcd as backing store.
type APISrvDynConf struct {

	// path to config file
	configFilePath string

	// conf keeps last good configuration
	conf *Config

	// lock
	lock sync.RWMutex
}

func NewAPISrvDynConf(key string) APISrvDynConfigurator {
	return &APISrvDynConf{
		configFilePath: key,
	}
}

func (c *APISrvDynConf) GetConf() (*Config, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.conf == nil {
		return c.getConf()
	}
	return c.conf, nil

}

func (c *APISrvDynConf) getConf() (*Config, error) {

	var cfg Config
	b, err := os.ReadFile(c.configFilePath)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = json.Unmarshal(b, &cfg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &cfg, nil

}
