package driver

import (
	"sync"
	"time"

	"github.com/kubeedge/mapper-framework/pkg/common"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// CustomizedDev is the customized device configuration and client information.
type CustomizedDev struct {
	Instance         common.DeviceInstance
	CustomizedClient *CustomizedClient
}

type CustomizedClient struct {
	deviceMutex sync.Mutex
	ProtocolConfig
	// Additional vars
	MqttClient mqtt.Client
	Data       map[string]interface{}
	LastTime   time.Time
}

type ProtocolConfig struct {
	ProtocolName string `json:"protocolName"`
	ConfigData   `json:"configData"`
}

type ConfigData struct {
	// protocol custom config data
	BrokerAddr string `json:"brokerAddr"`
	Topic      string `json:"topic"`
}

type VisitorConfig struct {
	ProtocolName      string `json:"protocolName"`
	VisitorConfigData `json:"configData"`
}

type VisitorConfigData struct {
	DataType string `json:"dataType"`
	// visitor custom config data
	FieldName string `json:"fieldName"`
}
