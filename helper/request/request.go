package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

// 发送GET请求
// url：         请求地址
// response：    请求返回的内容
func Get(url string, timeout int) (string, error) {

	// 超时时间：5秒
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	//var buffer [512]byte
	//result := bytes.NewBuffer(nil)
	//for {
	//	n, err := resp.Body.Read(buffer[0:])
	//	result.Write(buffer[0:n])
	//	if err != nil && err == io.EOF {
	//		break
	//	} else if err != nil {
	//		panic(err)
	//	}
	//}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return string(result), errors.New(fmt.Sprintf("请求响应报错，错误码：%d", resp.StatusCode))
	}
	return string(result), nil
}

// 发送POST请求
// url：         请求地址
// data：        POST请求提交的数据
// contentType： 请求体格式，如：application/json
// content：     请求放回的内容
func Post(url string, data interface{}, contentType string, timeout int) (string, error) {

	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	var jsonBytes []byte
	if strings.ToLower(reflect.TypeOf(data).Name()) != "string" {
		jsonBytes, _ = json.Marshal(data)
	} else {
		jsonBytes = []byte(data.(string))
	}
	resp, err := client.Post(url, contentType, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return string(result), errors.New(fmt.Sprintf("请求响应报错，错误码：%d", resp.StatusCode))
	}
	return string(result), nil
}
