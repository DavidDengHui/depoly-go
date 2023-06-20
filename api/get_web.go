package handler

import (
	"encoding/json"
	"fmt"
	"io"
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

	// 设置url、type、filename三个字符串变量
	var url string
	// 从get请求中获取url、type参数值
	url = r.URL.Query().Get("url")
	// 判断请求方式如果是POST，则将请求body部分中设置的url和type参数值覆盖
	if r.Method == "POST" {
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
		if data["url"] != "" {
			url = data["url"]
		}
	}
	// 如果没有获取到url值，则status_data["doit"]="NO_URL"、status_data["callback"]="INVALID_HOOK"，将status_data转换为json输出
	if url == "" {
		status_data["doit"] = "NO_URL"
		status_data["callback"] = "INVALID_HOOK"
		w.Header().Set("Content-Type", "application/json")
        json_data, err := json.Marshal(status_data)
        if err != nil {
            fmt.Println(err)
            return
        }
        w.Write(json_data)
        return
    }

		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// 将返回的数据直接写入到页面
		io.Copy(w, resp.Body)
		
}