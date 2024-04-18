package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	CheckinUrl = "https://glados.rocks/api/user/checkin"
	StatusUrl  = "https://glados.rocks/api/user/status"
)

type CheckinResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Points  int    `json:"points"`
}

type StatusResult struct {
	Code int `json:"code"`
	Data struct {
		Email    string `json:"email"`
		LeftDays string `json:"leftDays"`
	} `json:"data"`
}

func GetCookie() string {
	cookiePtr := flag.String("cookie", "", "GLaDOS Cookie")
	flag.Parse()
	if *cookiePtr != "" {
		return *cookiePtr
	}
	return os.Getenv("GLADOS_COOKIE")
}

func Request(url, method string, data, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	var req *http.Request
	var err error
	if method == "POST" {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}
	} else if method == "GET" {
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil

}

func main() {
	cookie := GetCookie()
	if cookie == "" {
		fmt.Println("请使用 -cookie 参数提供 GLaDOS Cookie 或者设置环境变量 GLADOS_COOKIE")
		return
	}
	data := map[string]string{"token": "glados.one"}
	headers := map[string]string{
		"cookie":       cookie,
		"content-type": "application/json;charset=UTF-8",
		"user-agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
	}
	body, err := Request(CheckinUrl, "POST", data, headers)
	if err != nil {
		fmt.Println("签到请求失败:", err)
		return
	}
	checkinResult := CheckinResult{}
	json.Unmarshal(body, &checkinResult)

	body, err = Request(StatusUrl, "GET", nil, headers)
	if err != nil {
		fmt.Println("状态请求失败:", err)
		return
	}
	statusResult := StatusResult{}
	json.Unmarshal(body, &statusResult)
	if checkinResult.Message != "" {
		fmt.Printf("签到成功，账号=%s,签到信息=%s, 剩余天数=%s\n", statusResult.Data.Email, checkinResult.Message, strings.Split(statusResult.Data.LeftDays, ".")[0])
	} else {
		fmt.Println("签到失败!")
	}
}
