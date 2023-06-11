package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"

	"github.com/apex/gateway"
)

var (
	port = flag.Int("port", -1, "specify a port")
)

func main() {
	flag.Parse()

	http.HandleFunc("/api/bili", bili)
	http.HandleFunc("/api/get_api", getApiHandler)
	listener := gateway.ListenAndServe
	portStr := "n/a"

	if *port != -1 {
		portStr = fmt.Sprintf(":%d", *port)
		listener = http.ListenAndServe
		http.Handle("/", http.FileServer(http.Dir("./static")))
	}

	log.Fatal(listener(portStr, nil))
}

func bili(w http.ResponseWriter, r *http.Request) {
	url := "https://api.bilibili.com/x/web-interface/popular?ps=20&pn=1"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(body)
}

func getApiHandler(w http.ResponseWriter, r *http.Request) {
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
			// 有则输出json数据{"status":"success"}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"success"}`))
	} else {
			// 否则输出json数据{"status":"error"}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"error"}`))
	}
}