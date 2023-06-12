package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// 定义一个全局数组变量status_data，包含status、code、doit、callback四个键
var status_data = make(map[string]interface{})

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	status_data["status"] = "error"
	status_data["code"] = "1001"
	status_data["doit"] = "NO_KEY"
	status_data["callback"] = "INVALID_KEY"

	// 定义url、type、header、send_data四个变量
	var url, typ string
	var header, send_data map[string]interface{}
	// 判断请求方法
	if r.Method == "GET" {
		// 从get请求中获取url、type、header、data参数值
		url = r.URL.Query().Get("url")
		typ = strings.ToUpper(r.URL.Query().Get("type"))
		header = JsonToMap(r.URL.Query().Get("header"))
		send_data = JsonToMap(r.URL.Query().Get("data"))
	} else if r.Method == "POST" {
		// 从post请求的body中获取url、type、header、data参数值
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		var data map[string]string
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println(err)
			return
		}
		url = data["url"]
		typ = strings.ToUpper(data["type"])
		header = JsonToMap(data["header"])
		send_data = JsonToMap(data["data"])
	}

	// 判断如果没有获取到url值，则返回错误信息并结束程序
	if url == "" {
		status_data["code"] = "1001"
		status_data["doit"] = "NO_URL"
		status_data["callback"] = "INVALID_HOOK"
        json_data, err := json.Marshal(status_data)
        if err != nil {
            fmt.Println(err)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.Write(json_data)
        return
    }

    // 将status_data["doit"]赋值为一个多维数组，其中包含"url"键值为url值，"type"键值为type值，"headers"键值为header数组，"data"键值为send_data数组
    status_data["doit"] = map[string]interface{}{
        "url": url,
        "type": typ,
        "headers": header,
        "data": send_data,
    }

    // 根据type的方式向url发送请求，携带header和send_data，并返回响应数据
    var resp *http.Response
    var err error
    if typ == "POST" {
        // 如果type是post，使用http.Post方法发送请求，将send_data转换为json格式的字节切片作为body参数
        json_data, err := json.Marshal(send_data)
        if err != nil {
            fmt.Println(err)
            return
        }
        // return
        // 创建一个自定义的请求对象，设置请求方法为post，请求的url为url，请求的数据为json_data
        req, err := http.NewRequest("POST", url, ioutil.NopCloser(bytes.NewReader(json_data)))
        if err != nil {
            fmt.Println(err)
            return
        }
        // 设置请求的标头，遍历header变量中的键值对，使用http.Header.Set方法设置对应的头部字段和值
        for k, v := range header {
            req.Header.Set(k, fmt.Sprint(v))
        }
        // 创建一个http.Client对象
        client := &http.Client{}
        // 使用http.Client.Do方法发送请求，并获取响应
        resp, err := client.Do(req)
        if err != nil {
            fmt.Println(err)
            return
        }
        defer resp.Body.Close()
        // 读取响应的状态码和数据，并输出到页面上
        status := resp.StatusCode
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            fmt.Println(err)
            return
        }
        data := string(body)
        // 设置status_data["status"]为"success"，status_data["code"]为"1102-响应状态码"
        status_data["status"] = "success"
        status_data["code"] = fmt.Sprintf("1102-%d", status)
        // 读取响应的body数据，并将其赋值给status_data["callback"]
        status_data["callback"] = data
    } else {
        // 如果type不是post，默认使用get方法发送请求，将send_data转换为查询字符串作为url参数
        query := make(map[string]string)
        for k, v := range send_data {
            query[k] = fmt.Sprint(v)
        }
        // 判断url中是否包含?字符，有则使用&连接，否则使用?连接
        if strings.Contains(url, "?") {
            resp, err = http.Get(url + "&" + StrEncode(query))
        } else {
            resp, err = http.Get(url + "?" + StrEncode(query))
        }
        if err != nil {
            fmt.Println(err)
            return
        }
        defer resp.Body.Close()
        // 读取响应的body数据，并将其赋值给status_data["callback"]
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            fmt.Println(err)
            return
        }
        status_data["callback"] = string(body)

        // 设置status_data["status"]为"success"，status_data["code"]为"1101-响应状态码"
        status_data["status"] = "success"
        status_data["code"] = fmt.Sprintf("1101-%d", resp.StatusCode)
    }
    // 将status_data转换为json数据并返回
    json_data, err := json.Marshal(status_data)
    if err != nil {
        fmt.Println(err)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(json_data)
    return
}

// 定义一个辅助函数，将json格式的字符串转换为map类型
func JsonToMap(s string) map[string]interface{} {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return m
}

func StrEncode(query map[string]string) string {
	// 定义一个字符串切片，用于存储编码后的键值对
	var pairs []string
	// 遍历query变量中的键值对
	for k, v := range query {
		// 使用url.QueryEscape函数对键和值进行编码，替换特殊字符为%XX序列
		k = url.QueryEscape(k)
		v = url.QueryEscape(v)
		// 将编码后的键和值用=连接起来，形成一个键值对字符串
		pair := k + "=" + v
		// 将键值对字符串追加到字符串切片中
		pairs = append(pairs, pair)
	}
	// 使用strings.Join函数将字符串切片中的元素用&连接起来，形成一个查询字符串
	queryStr := strings.Join(pairs, "&")
	// 返回查询字符串
	return queryStr
}
