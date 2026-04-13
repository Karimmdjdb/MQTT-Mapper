package driver

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kubeedge/mapper-framework/pkg/common"
	"k8s.io/klog/v2"
)

func NewClient(protocol ProtocolConfig) (*CustomizedClient, error) {
	fmt.Println("NEW CLIENT")
	client := &CustomizedClient{
		ProtocolConfig: protocol,
		deviceMutex:    sync.Mutex{},
		Data:           map[string]interface{}{},
		DataForPush:    map[string]interface{}{},
	}
	return client, nil
}

func (c *CustomizedClient) InitDevice() error {
	// options for MQTT client
	opts := mqtt.NewClientOptions()
	opts.AddBroker(c.ProtocolConfig.ConfigData.BrokerAddr)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		c.deviceMutex.Lock()
		defer c.deviceMutex.Unlock()
		// payload recovery
		var payload map[string]interface{}
		if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
			klog.Errorf("Failed to parse MQTT payload: %v", err)
			return
		}
		// Cooja data recovery
		if data, ok := payload["d"].(map[string]interface{}); ok {
			for k, v := range data {
				c.Data[k] = v
				c.DataForPush[k] = v
			}
		}
		klog.Infof("Received MQTT message on topic %s", msg.Topic())
	})

	// connection to the broker
	c.MqttClient = mqtt.NewClient(opts)
	if token := c.MqttClient.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to broker %s: %v",
			c.ProtocolConfig.ConfigData.BrokerAddr, token.Error())
	}
	klog.Infof("Connected to MQTT broker: %s", c.ProtocolConfig.ConfigData.BrokerAddr)

	// sub to mote topic
	if token := c.MqttClient.Subscribe(c.ProtocolConfig.Topic, 0, nil); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %v", c.ProtocolConfig.Topic, token.Error())
	}

	go func() {
		for {
			fmt.Print("D : ")
			fmt.Println(c.Data)
			fmt.Print("DP : ")
			fmt.Println(c.DataForPush)
			time.Sleep(200 * time.Millisecond)
		}
	}()

	klog.Infof("Subscribed to topic: %s", c.ProtocolConfig.Topic)
	return nil
}

func (c *CustomizedClient) GetDeviceData(visitor *VisitorConfig) (interface{}, error) {
	c.deviceMutex.Lock()
	defer c.deviceMutex.Unlock()

	fieldName := visitor.VisitorConfigData.FieldName
	value, ok := c.Data[fieldName]
	if !ok {
		return nil, fmt.Errorf("field '%s' not found in data", fieldName)
	}
	delete(c.Data, fieldName)
	return value, nil
}

func (c *CustomizedClient) GetDeviceDataForPush(visitor *VisitorConfig) (interface{}, error) {
	c.deviceMutex.Lock()
	defer c.deviceMutex.Unlock()

	fieldName := visitor.VisitorConfigData.FieldName
	value, ok := c.DataForPush[fieldName]
	if !ok {
		return nil, fmt.Errorf("field '%s' not found in data for push", fieldName)
	}
	delete(c.DataForPush, fieldName)
	return value, nil
}

func (c *CustomizedClient) DeviceDataWrite(visitor *VisitorConfig, deviceMethodName string, propertyName string, data interface{}) error {
	// TODO: add the code to write device's data
	// you can use c.ProtocolConfig and visitor to write data to device
	return nil
}

func (c *CustomizedClient) SetDeviceData(data interface{}, visitor *VisitorConfig) error {
	// TODO: set device's data
	// you can use c.ProtocolConfig and visitor
	return nil
}

func (c *CustomizedClient) StopDevice() error {
	if c.MqttClient != nil && c.MqttClient.IsConnected() {
		c.MqttClient.Disconnect(250)
		klog.Info("Disconnected from MQTT broker")
	}
	return nil
}

func (c *CustomizedClient) GetDeviceStates() (string, error) {
	if c.MqttClient != nil && c.MqttClient.IsConnected() {
		return common.DeviceStatusOK, nil
	}
	return common.DeviceStatusDisCONN, nil
}
