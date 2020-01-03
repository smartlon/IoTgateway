package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Container struct {
	ContainerID string `json:"ContainerID"`
	Description string `json:"Description"`
	Timestamp string `json:"Timestamp"`
	Location  string `json:"Location"`
	Used      string   `json:"Used"`  //true is used and false is not used
	Holder  string `json:"Holder"`
}

type QueryContainerResult struct {
	Container           Container       `json:"container"`
	MamState        MAMState       `json:"mamstate"`
}

type QueryAllContainerResult struct {
	Key           string       `json:"Key"`
	Record        Container       `json:"Record"`
}

type MAMState struct {
	Root           string       `json:"root"`
	SideKey        string       `json:"sideKey"`
	Seed           string       `json:"seed"`

}
type Response struct {
	Code           string       `json:"code"`
	Message        string       `json:"message"`
	Result           string       `json:"result"`
}

func transmit(seed string, sideKey string, message *IoTData) {
	client := &http.Client{}
	data := make(map[string]interface{})
	data["Message"] = message
	data["Seed"] = seed
	data["SideKey"] = sideKey
	bytesData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	req, err := http.NewRequest("POST",SERVICEURL+"/iota/mamtransmit",bytes.NewReader(bytesData))
	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	var response Response
	err = json.Unmarshal(body,&response)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("code = %s, message = %s \n", response.Code, response.Message)

}

func queryContainer(container string) (string, string) {
	client := &http.Client{}
	data := make(map[string]interface{})
	data["Func"] = "QueryContainer"
	args := []string{container}
	data["Args"] = args
	bytesData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	req, err := http.NewRequest("POST",SERVICEURL+"/fabric/querycontainer",bytes.NewReader(bytesData))
	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	var response Response
	err = json.Unmarshal(body,&response)
	if err != nil {
		fmt.Println(err)
	}
	if response.Code == "200" {
		fmt.Println(response.Message)
		var result QueryContainerResult
		err :=json.Unmarshal([]byte(response.Result),&result)
		if err != nil {
			fmt.Println(err)
		}
		if result.Container.Used == "true" {
			return result.MamState.Seed,result.MamState.SideKey
		}
	}
	return "",""
}

func queryAllContainer() []QueryAllContainerResult {
	client := &http.Client{}
	data := make(map[string]interface{})
	data["Func"] = "QueryAllContainers"
	args := []string{}
	data["Args"] = args
	bytesData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	req, err := http.NewRequest("POST",SERVICEURL+"/fabric/queryallcontainers",bytes.NewReader(bytesData))
	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	var response Response
	err = json.Unmarshal(body,&response)
	if err != nil {
		fmt.Println(err)
	}
	if response.Code == "200" {
		//fmt.Println(response.Message)
		var result []QueryAllContainerResult
		err :=json.Unmarshal([]byte(response.Result),&result)
		if err != nil {
			fmt.Println(err)
		}
		return result
	}
	return []QueryAllContainerResult{}
}
