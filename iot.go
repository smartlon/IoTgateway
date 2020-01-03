package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Device struct {
	SensorID    string `json:"SensorID"`
	Temperature int `json:"Temperature"`
	Humidity int `json:"Humidity"`
}

type IoTData struct {
	ContainerID        string `json:"ContainerID"`
	Temperature string `json:"Temperature"`
	Location    string `json:"Location"`
	Time        string `json:"Time"`
	Status        string `json:"Status"`
}
func timeStamp() string {
	return strconv.FormatInt(time.Now().UnixNano() / 1000000, 10)
}
//接口函数，当收到订阅的消息时会启用这个回调函数
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	//fmt.Printf("TOPIC: %s\n", msg.Topic())
	//fmt.Printf("MSG: %s\n", msg.Payload())
	var device Device
	err := json.Unmarshal(msg.Payload(),&device)
	if err != nil {
		fmt.Println(err.Error())
	}
	data := &IoTData{
		device.SensorID,
		strconv.Itoa(device.Temperature),
		"",
		timeStamp(),
		"run",
	}
	time.Sleep(time.Duration(10)*time.Second)
	transmit(seed,sideKey,data)
}
func startReciver(stop  chan string, clientID string) {
	opts := MQTT.NewClientOptions().AddBroker("tcp://2y4x807794.wicp.vip:56866")
	//连接的客户端名字
	opts.SetClientID(clientID)
	opts.SetDefaultPublishHandler(f)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	seed, sideKey = queryContainer(clientID)
	fmt.Printf("seed = %s,   sideKey = %s \n",seed,sideKey)
	//订阅消息
	if token := c.Subscribe("IOTA/"+clientID, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	defer func() {
		//取消订阅和断开连接
		if token := c.Unsubscribe("IOTA/"+clientID); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}
		c.Disconnect(250)
	}()

	containerListen[clientID] = true
	flag := <- stop
	fmt.Printf("%s is %s \n",clientID,flag)
	containerListen[clientID] = false
}