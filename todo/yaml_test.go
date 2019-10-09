package todo

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	CONF_CONSUMER_FILE_PATH        = "CONF_CONSUMER_FILE_PATH"
	CONF_PROVIDER_FILE_PATH        = "CONF_PROVIDER_FILE_PATH"
	APP_LOG_CONF_FILE       string = "APP_LOG_CONF_FILE"
	//$ export CONF_CONSUMER_FILE_PATH="config/client.yml"
	//$ export APP_LOG_CONF_FILE="config/log.yml"
)

func TestYaml(t *testing.T) {
	path, _ := filepath.Abs("../testdata/consumer_config.yml")
	pc := &ProviderConfig{}
	Yaml2Struct(path, pc)

	fmt.Printf("path:%s,pc:%v \n", path, pc.ProtocolConf)
}

func TestGetenv(t *testing.T) {
	val := os.Getenv("CONF_CONSUMER_FILE_PATH")
	fmt.Println(val)
}

type ServerConfig struct {
	// session
	SessionTimeout string `default:"60s" yaml:"session_timeout" json:"session_timeout,omitempty"`
	sessionTimeout time.Duration
	SessionNumber  int `default:"1000" yaml:"session_number" json:"session_number,omitempty"`

	// grpool
	GrPoolSize  int `default:"0" yaml:"gr_pool_size" json:"gr_pool_size,omitempty"`
	QueueLen    int `default:"0" yaml:"queue_len" json:"queue_len,omitempty"`
	QueueNumber int `default:"0" yaml:"queue_number" json:"queue_number,omitempty"`

	// session tcp parameters
	//GettySessionParam GettySessionParam `required:"true" yaml:"getty_session_param" json:"getty_session_param,omitempty"`
}

type ProviderConfig struct {
	//BaseConfig   `yaml:",inline"`
	Filter       string `yaml:"filter" json:"filter,omitempty" property:"filter"`
	ProxyFactory string `yaml:"proxy_factory" default:"default" json:"proxy_factory,omitempty" property:"proxy_factory"`

	//ApplicationConfig *ApplicationConfig         `yaml:"application_config" json:"application_config,omitempty" property:"application_config"`
	//Registries        map[string]*RegistryConfig `yaml:"registries" json:"registries,omitempty" property:"registries"`
	//Services          map[string]*ServiceConfig  `yaml:"services" json:"services,omitempty" property:"services"`
	//Protocols         map[string]*ProtocolConfig `yaml:"protocols" json:"protocols,omitempty" property:"protocols"`
	ProtocolConf interface{} `yaml:"protocol_conf" json:"protocol_conf,omitempty" property:"protocol_conf" `
	FilterConf   interface{} `yaml:"filter_conf" json:"filter_conf,omitempty" property:"filter_conf" `
}
