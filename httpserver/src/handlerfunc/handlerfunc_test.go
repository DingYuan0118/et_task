package handlerfunc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	pb "et-protobuf3/src/gomicroapi"
	"httpserver/src/auth"
)

type AuthResponse struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data DataResponse `json:"data"`
}

type DataResponse struct {
	Token string `json:"token"`
}


func TestParseToken(t *testing.T) {
	tests := []struct{
		token string
		want bool
	}{
		{token: Token, want: true},
		{token: Token + "1", want: false},
	}
	
	for _, test := range tests {
		got, err := auth.ParseToken(test.token)
		res := (err == nil)
		if res != test.want {
			t.Errorf("Token check func failed. Got: %v, Want %v", got, test.want)
		}

	}
}

const dst = "http://0.0.0.0:8080"

var jsonData = []byte(`{
	"username": "Ding",
	"password": "dingyuan"
}`)

const Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IkRpbmciLCJleHAiOjIwMjAwMzAyODMsImlzcyI6ImVudHJ5IHRhc2sifQ.RVV8jwO5PDC-CW8cu25ulnLGOFitE8Ibxe02k-tLqH0"

func body2string(response *http.Response) []byte {
	body, _ := ioutil.ReadAll(response.Body)
	return body
}

func TestUserLogin(t *testing.T) {
	request, err := http.NewRequest("POST", dst+"/login", bytes.NewBuffer(jsonData))
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
	request, err := http.NewRequest("GET", dst+"/query", nil)
	if err != nil {
		t.Error(err.Error())
	}
	q := request.URL.Query()
	q.Add("username", "Ding")
	request.URL.RawQuery = q.Encode()
	request.Header.Set("Authorization", "Bearer "+Token)

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
	request, err := http.NewRequest("POST", dst+"/update-nickname", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Error(err.Error())
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", "Bearer "+ Token)
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

func TestUploadPic(t *testing.T) {
	url := dst + "/upload-pic"
	params := map[string]string{"username":"Ding"}
	filenameField := "file"
	fileName := "test_small.png"
	filepath := "/Users/yuan.ding/Desktop/code/entry_task/images/" + fileName
	file, err := os.Open(filepath)
	if err != nil {
		t.Error(err)
	}
	// 返回 response
	response, err := UploadPicPOST(url, params, filenameField, fileName, file)
	if err != nil {
		t.Error(err)
	}

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


// 用户发送POST请求的辅助函数供测试 UploadPic 使用
func UploadPicPOST(url string, params map[string]string, filenameField string, fileName string, file io.Reader) (*http.Response, error) {
    body := new(bytes.Buffer)
    writer := multipart.NewWriter(body)
	var HttpClient = &http.Client{}

    formFile, err := writer.CreateFormFile(filenameField, fileName)
    if err != nil {
        return nil, err
    }

    _, err = io.Copy(formFile, file)
    if err != nil {
        return nil, err
    }

    for key, val := range params {
        _ = writer.WriteField(key, val)
    }

    err = writer.Close()
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Authorization", "Bearer "+ Token)
    if err != nil {
        return nil, err
    }
    //req.Header.Set("Content-Type","multipart/form-data")
    req.Header.Add("Content-Type", writer.FormDataContentType())

    resp, err := HttpClient.Do(req)
    if err != nil {
        return nil, err
    }
    return resp, nil
}
