package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/hashkeypro/api-src/golang/hashkey"
)

func main() {

	apiPath := "/v1/info/time"

	// PUB-PRIV
	// 公钥-私钥

	// preview
	baseURL := "https://preview-pro.hashkey.com/APITrade"
	apiKey := "MTU0MjEwNDAwMTA1NjAwMDAwMDAwNTQ="
	privateKey := "uvX6WIUzE5jJLMszT7elkTMKgRZEoYkx7X7mTpPWyXo="

	// HMAC
	// 消息验证码
	apiKeyHMAC := "MTU0NDUwODQ0NjI0NTAwMDAwMDAwNTQ="
	secretKey := "vprggEasLOksdmut6WcFvuv4oUuAbewdkGJY1fgAvBw="

	// generate the message to sign
	// 生成待签消息
	timestamp := strconv.FormatInt(time.Now().Unix()*1000, 10)
	originData := []byte(timestamp + "GET" + apiPath)

	// signature
	// 签名
	signStr, err := hashkey.ECCSignature(originData, privateKey)

	// hmac
	// 消息验证码
	hmacStr := hashkey.SHA256HMAC(originData, secretKey)

	req, err := http.NewRequest("GET", baseURL+apiPath, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// request header for PUB-PRIV
	// PUB-PRIV 请求消息头
	req.Header["API-KEY"] = []string{apiKey}
	req.Header["API-SIGNATURE"] = []string{signStr}
	req.Header["API-TIMESTAMP"] = []string{timestamp}
	req.Header["AUTH-TYPE"] = []string{"PUB-PRIV"}
	req.Header["Content-Type"] = []string{"application/json"}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("PUB-PRIV\n" + string(b))

	// request header for HMAC
	// HMAC 请求消息头
	req.Header["API-KEY"] = []string{apiKeyHMAC}
	req.Header["API-SIGNATURE"] = []string{hmacStr}
	req.Header["API-TIMESTAMP"] = []string{timestamp}
	req.Header["AUTH-TYPE"] = []string{"HMAC"}
	req.Header["Content-Type"] = []string{"application/json"}

	respHMAC, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer respHMAC.Body.Close()
	bHMAC, err := ioutil.ReadAll(respHMAC.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("HMAC\n" + string(bHMAC))

}
