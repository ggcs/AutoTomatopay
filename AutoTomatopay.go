package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	//"fmt"
)

var email = ""
var password = ""
var method = "1" 
var alipayLimit = 1000
var wechatLimit = 0

func main() {
	urlHome := "https://b.fanqieui.com/"
	client := &http.Client{}
	rand.Seed(time.Now().UnixNano())
	addLog(time.Now().String(), false)
///////////////////////////////////
//////////////////////////////////
	requestBody := url.Values{}
	requestBody.Set("email", email)
	requestBody.Set("password", password)
	request, _ := http.NewRequest("POST", urlHome+"_login.php", strings.NewReader(requestBody.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, _ := client.Do(request)
	cookies := response.Cookies()
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	var responseBodyMap map[string]interface{}
	responseBodyMap = make(map[string]interface{})
	errorLog := json.Unmarshal(responseBody, &responseBodyMap)
	if errorLog != nil {
		addLog(errorLog.Error(), true)
	}
	if responseBodyMap["code"] != "1" {
		addLog("Login fail", true)
	}
///////////////////////////////////
//////////////////////////////////
	request, _ = http.NewRequest("GET", urlHome+"dashboard/withdrawal.php", nil)
	for _, cookieIndex := range cookies {
		request.AddCookie(cookieIndex)
	}
	response, _ = client.Do(request)
	cookies1 := response.Cookies()
	defer response.Body.Close()
	responseBody, _ = ioutil.ReadAll(response.Body)
	responseBodyString := string(responseBody)
	token := GetBetweenStr(responseBodyString,"value=\"","\"")
	cny_s := GetBetweenStr(responseBodyString,"¥ ","\"")
	cny_i,_:=strconv.ParseFloat(cny_s, 64)
	cny := int(cny_i)
	if method == "1" && cny > alipayLimit {
		cny = alipayLimit
	}
	if  cny < 2 {
		addLog("Must be greater than 2", true)
	}
	fmt.Println(token)
///////////////////////////////////
//////////////////////////////////
	requestBody = url.Values{}
	requestBody.Set("token", token)
	requestBody.Set("cny", strconv.Itoa(cny))
	request, _ = http.NewRequest("POST", urlHome+"dashboard/_withdrawal.php", strings.NewReader(requestBody.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for _, cookieIndex := range cookies {
		request.AddCookie(cookieIndex)
	}
	for _, cookieIndex := range cookies1 {
		request.AddCookie(cookieIndex)
	}
	response, _ = client.Do(request)
	defer response.Body.Close()
	responseBody, _ = ioutil.ReadAll(response.Body)
	errorLog = json.Unmarshal(responseBody, &responseBodyMap)
	if errorLog != nil {
		addLog(errorLog.Error(), true)
	}
	if responseBodyMap["code"] != "1" {
		addLog("fail", true)
	}else{
		addLog("succee", true)
	}
}

var allLog = ""

func addLog(log string, exit bool) {
	allLog += log
	allLog += "\n"

	if exit {
		allLog += "\n"
		logFile, _ := os.OpenFile("1.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			os.ModeAppend)
		logFile.WriteString(allLog)
		os.Exit(0)
	}
}

func GetBetweenStr(str, start, end string) string {
    n := strings.Index(str, start)
    if n == -1 {
        n = 0
    } else {
        n = n + len(start)
    }
    str = string([]byte(str)[n:])
    m := strings.Index(str, end)
    if m == -1 {
        m = len(str)
    }
    str = string([]byte(str)[:m])
    return str
}
