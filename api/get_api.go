package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)


func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// 定义一个全局数组变量status_data，包含status、code、doit、callback四个键
	var status_data = make(map[string]interface{})
	// 初始化status_data["status"]="error"、status_data["code"]="1001"、status_data["doit"]="NO_KEY"、status_data["callback"]="INVALID_KEY"
	status_data["status"] = "error"
	status_data["code"] = "1001"
	status_data["doit"] = "NO_KEY"
	status_data["callback"] = "INVALID_KEY"

	// 定义一个url变量
	var url string
	// 先获取get请求中设置的url参数赋值给url
	url = r.URL.Query().Get("url")
	// 再判断如果是post方式的请求
	if r.Method == "POST" {
		// 则判断请求的body数据中有无url参数
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
		// 如果有则将url参数值覆盖给url
		if data["url"] != "" {
			url = data["url"]
		}
	}
	// 最后再判断变量url是否有值
	if url != "" {
		// 增加将url值赋给status_data["doit"]
		status_data["doit"] = url
		// 有则向该url发送get请求
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		// 如果请求成功了直接输出返回的数据
		if resp.StatusCode == http.StatusOK {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
			w.Write(body)
		} else {
			// 否则输出json数据{"status":"error","code":"1002","doit":"获取到的url值","callback":"INVALID_URL_加请求状态"}
			w.Header().Set("Content-Type", "application/json")
			status_data["code"] = "1002"
			status_data["callback"] = fmt.Sprintf("INVALID_URL_%d", resp.StatusCode)
			json_data, err := json.Marshal(status_data)
			if err != nil {
				fmt.Println(err)
				return
			}
			w.Write(json_data)
		}
	} else {
		// 判断变量url没有值时，status_data["doit"]="NO_URL"、status_data["callback"]="INVALID_HOOK"
		status_data["doit"] = "NO_URL"
		status_data["callback"] = "INVALID_HOOK"
		// 然后将数组status_data转换为json数据后输出
		w.Header().Set("Content-Type", "application/json")
        json_data, err := json.Marshal(status_data)
        if err != nil {
            fmt.Println(err)
            return
        }
        w.Write(json_data)
    }
}