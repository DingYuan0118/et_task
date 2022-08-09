package handlerfunc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	pb "et-protobuf3/src/gomicroapi"
	"httpserver/src/auth"
)

type AuthResponse struct {
	Code		int 			`json:"code"`
	Msg			string 			`json:"msg"`
	Data		DataResponse 	`json:"data"`
}

type DataResponse struct {
	Token	string  `json:"token"`
}
func TestGenToken(t *testing.T) {
	_, err := auth.GenToken("ding")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestParseToken(t *testing.T) {
	token, _ := auth.GenToken("Ding")
	_, err := auth.ParseToken(token)
	if err != nil {
		t.Error(err.Error())
	}
}

const dst = "http://0.0.0.0:8080"
var jsonData = []byte(`{
	"username": "Ding",
	"password": "dingyuan"
}`)

const Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IkRpbmciLCJleHAiOjE2NTk2OTgwNDUsImlzcyI6ImVudHJ5IHRhc2sifQ.fzE9_PUJuigyUR5A6v-iDG1tGSXzi2jjDxlGH_2A8x0"

func GetAuthToken(t *testing.T) string {
	request, err := http.NewRequest("POST", dst + "/login", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Error(err.Error())
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// Get Token
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
	return token
}

func TestJWTToken(t *testing.T) {	
	token := GetAuthToken(t)

	// Token authentication
	get_request, err := http.NewRequest("GET", dst + "/home", nil)
	if err != nil {
        log.Print(err)
        t.Fatal(err.Error())
    }
	q := get_request.URL.Query()
	q.Add("username", "yuan")
	get_request.URL.RawQuery = q.Encode()
	get_request.Header.Set("Authorization", "Bearer " + token + "1")

	client := &http.Client{}
	get_response, err := client.Do(get_request)
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

func TestUserLogin(t *testing.T) {
	request, err := http.NewRequest("POST", dst + "/login", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Error(err.Error())
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// Get Token
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer response.Body.Close()
	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("return body: ", string(body))
}

func TestUserQuery(t *testing.T) {
	request, err := http.NewRequest("GET", dst + "/query", nil)
	if err != nil {
		t.Error(err.Error())
	}
	q := request.URL.Query()
	q.Add("username", "Ding")
	request.URL.RawQuery = q.Encode()
	request.Header.Set("Authorization", "Bearer " + Token)

	// use Token get query result
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer response.Body.Close()
	queryResponse := new(pb.QueryReturn)
	body := body2string(response)
	err = json.Unmarshal(body, queryResponse)
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Printf("Query return info: \n %+v", queryResponse)
}

func TestUpdateNickname(t *testing.T) {
	jsonData = []byte(`{
		"username": "Ding",
		"nickname": "dingyuan12345"
	}`)
	request, err := http.NewRequest("POST", dst + "/update-nickname", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Error(err.Error())
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", "Bearer " + Token)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer response.Body.Close()
	body := body2string(response)
	updateNicknameResponse := new(pb.UpdateNicknameReturn)
	err = json.Unmarshal(body, updateNicknameResponse)
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	fmt.Println("return body: ", string(body))
}
