package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestGenToken(t *testing.T) {
	_, err := GenToken("ding")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestParseToken(t *testing.T) {
	token, _ := GenToken("ding")
	_, err := ParseToken(token)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestJWTToken(t *testing.T) {	
	dst := "http://0.0.0.0:8080"
	var jsonData = []byte(`{
		"username": "ding",
		"password": "123"
	}`)
	request, err := http.NewRequest("POST", dst + "/auth", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Error(err.Error())
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer response.Body.Close()
	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	var target = AuthResponse{}
	err = json.Unmarshal(body, &target)
	if err != nil {
		t.Error(err.Error())
	}
	token := target.Data.Token
	get_request, err := http.NewRequest("GET", dst + "/home", nil)
	if err != nil {
        log.Print(err)
        t.Fatal(err.Error())
    }
	q := get_request.URL.Query()
	q.Add("name", "yuan")
	get_request.URL.RawQuery = q.Encode()
	get_request.Header.Set("Authorization", "Bearer " + token)
	client2 := &http.Client{}
	get_response, err := client2.Do(get_request)
	if err != nil {
		t.Fatal("get request failed")
	}
	defer get_response.Body.Close()
	get_body := body2string(get_response)
	fmt.Println("auth result:\n", string(get_body))
}

func body2string(response *http.Response) []byte {
	body, _ := ioutil.ReadAll(response.Body)
	return body
}