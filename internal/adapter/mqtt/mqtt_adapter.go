package mqtt

import (
	"encoding/json"
	"fmt"
	"github.com/APCS20-Thesis/Backend/config"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"github.com/APCS20-Thesis/Backend/utils"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-logr/logr"
	"gorm.io/gorm"
)

type MqttSeverity string

const (
	MqttSeverity_Success MqttSeverity = "success"
	MqttSeverity_Error   MqttSeverity = "error"
)

type MqttAdapter interface {
	Connect()

	Publish(topic string, message interface{})
	Sub(topic string)
	Disconnect()
	PublishNotification(accountUuid string, message interface{})
}

var (
	messagePubHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	}

	connectHandler MQTT.OnConnectHandler = func(client MQTT.Client) {
		fmt.Println("Connected")
	}

	connectLostHandler MQTT.ConnectionLostHandler = func(client MQTT.Client, err error) {
		fmt.Printf("Connect lost: %v", err)
	}
)

type (
	Notification struct {
		Status     int32            `json:"status"`
		Message    string           `json:"message"`
		Severity   string           `json:"severity"`
		ActionType model.ActionType `json:"action_type"`
	}
)

type mqtt struct {
	log    logr.Logger
	client MQTT.Client
	db     *gorm.DB
}

func (m *mqtt) PublishNotification(accountUuid string, message interface{}) {
	topic := utils.GetMqttNotificationTopic(accountUuid)
	m.Publish(topic, message)
}

func (m *mqtt) Sub(topic string) {
	token := m.client.Subscribe(topic, 1, nil)
	token.Wait()
	m.log.Info("Subscribed to topic: %s", topic)

}

func (m *mqtt) Connect() {
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	topic := ""
	token := m.client.Subscribe(topic, 1, nil)
	token.Wait()
	m.log.Info("Mqtt Subscribe", "topic", topic)
}

func (m *mqtt) Publish(topic string, message interface{}) {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		m.log.Error(err, "Failed to publish message")
		return
	}
	token := m.client.Publish(topic, 0, false, jsonMessage)
	token.Wait()
	m.log.Info("Mqtt publish", "topic", topic)
}

func (m *mqtt) Disconnect() {
	m.client.Disconnect(250)
}

func NewMqttAdapter(config *config.Config, log logr.Logger, db *gorm.DB) (MqttAdapter, error) {
	opts := MQTT.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", config.MqttAdapterConfig.Host, config.MqttAdapterConfig.Port))
	opts.SetClientID(config.MqttAdapterConfig.ClientID)
	opts.SetUsername(config.MqttAdapterConfig.Username)
	opts.SetPassword(config.MqttAdapterConfig.Password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := MQTT.NewClient(opts)

	return &mqtt{
		log:    log,
		client: client,
		db:     db,
	}, nil
}
