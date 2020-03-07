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

type MAMState struct {
	Root           string       `json:"root"`
	SideKey        string       `json:"sideKey"`
	Seed           string       `json:"seed"`

}
type Response struct {
	Code           int       `json:"code"`
	Message        string       `json:"msg"`
	Result           string       `json:"data"`
}

type AllResponse struct {
	Code           int             `json:"code"`
	Message        string          `json:"msg"`
	Count          int             `json:"count"`
	Result         []Container       `json:"data"`
}

func transmit(seed string, sideKey string, message *IoTData,token string) {
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
	req.Header.Set("Authorization", "bearer " + token)
	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	var response Response
	err = json.Unmarshal(body,&response)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("code = %s, message = %s \n", response.Code, response.Message)

}

func queryContainer(container string,token string) (string, string) {
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
	req.Header.Set("Authorization", "bearer " + token)
	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	var response Response
	err = json.Unmarshal(body,&response)
	if err != nil {
		fmt.Println(err)
	}
	if response.Code == 200 {
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

func queryAllContainer(token string) []Container {
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
	req.Header.Set("Authorization", "bearer " + token)
	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	var response AllResponse
	err = json.Unmarshal(body,&response)
	if err != nil {
		fmt.Println(err)
	}
	if response.Code == 200 {

		return response.Result
	}
	return []Container{}
}

type userResp struct {
	Message string `json:"msg"`
	Success bool `json:"success"`
	Token string `json:"token"`
}

func enrollUser(userName,password,orgName string) string {
	client := &http.Client{}
	data := make(map[string]interface{})
	data["UserName"] = userName
	data["PassWord"] = password
	data["OrgName"] = orgName
	bytesData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	req, err := http.NewRequest("POST",SERVICEURL+"/enrolluser",bytes.NewReader(bytesData))
	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	var response userResp
	err = json.Unmarshal(body,&response)
	if err != nil {
		fmt.Println(err)
	}
	if response.Success  {

		return response.Token
	}
	return ""
}
