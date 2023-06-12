package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/apex/gateway"
)

// 定义一个全局数组变量status_data，包含status、code、doit、callback四个键
var status_data = make(map[string]interface{})

var (
	port = flag.Int("port", -1, "specify a port")
)

func main() {
	flag.Parse()

	http.HandleFunc("/get_api", getApiHandler)
	http.HandleFunc("/get_img", getImgHandler)
	http.HandleFunc("/readme", readmeHandler)
	http.HandleFunc("/send_api", sendApiHandler)
	listener := gateway.ListenAndServe
	portStr := "n/a"

	if *port != -1 {
		portStr = fmt.Sprintf(":%d", *port)
		listener = http.ListenAndServe
		http.Handle("/", http.FileServer(http.Dir("./static")))
	}

	log.Fatal(listener(portStr, nil))
}

func getApiHandler(w http.ResponseWriter, r *http.Request) {
	// 初始化status_data["status"]="error"、status_data["code"]="1001"、status_data["doit"]="NO_KEY"、status_data["callback"]="INVALID_KEY"
	status_data["callback"] = "INVALID_KEY"
	status_data["doit"] = "NO_KEY"
	status_data["code"] = "1001"
	status_data["status"] = "error"

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

func getImgHandler(w http.ResponseWriter, r *http.Request) {
	// 初始化status_data["status"]="error"、status_data["code"]="1001"、status_data["doit"]="NO_KEY"、status_data["callback"]="INVALID_KEY"
	status_data["callback"] = "INVALID_KEY"
	status_data["doit"] = "NO_KEY"
	status_data["code"] = "1001"
	status_data["status"] = "error"

	// 设置url、type、filename三个字符串变量
	var url, typ, filename string
	// 从get请求中获取url、type参数值
	url = r.URL.Query().Get("url")
	typ = r.URL.Query().Get("type")
	// filename初始化为"get_img"
	filename = "get_img"
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
			if data["type"] != "" {
					typ = data["type"]
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
	// 如果获取到url值，则将url中最后一个/字符之后的内容赋值给filename
	filename = url[strings.LastIndex(url, "/")+1:]
	// 再判断type是否有值，没有则将type赋值为"jpg"
	if typ == "" {
			typ = "jpg"
			// 再继续判断url结尾如果有.字符，则将从.之后开始到非字母结束的部分赋值给type，将filename中.字符及以后的部分删除
			if strings.Contains(filename, ".") {
					re := regexp.MustCompile(`\.[a-zA-Z]+`)
					match := re.FindString(filename)
					if match != "" {
							typ = match[1:]
							// filename = strings.TrimSuffix(filename, match)
							filename = filename[:strings.Index(filename, ".")]
					}
			}
	}
	// 然后向url发送get请求，将返回的图片显示到屏幕上，mimetype为image/type，设置ctrl+s保存图片时名字为filename.type
	resp, err := http.Get(url)
	if err != nil {
			fmt.Println(err)
			return
	}
	defer resp.Body.Close()
	
	w.Header().Set("Content-Type", fmt.Sprintf("image/%s", typ))
	
	// 设置ctrl+s保存图片时名字为filename.type
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s.%s", filename, typ))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(body)
}

func readmeHandler(w http.ResponseWriter, r *http.Request) {
	// 重定向页面到/
	http.Redirect(w, r, "/", http.StatusFound)
}

func sendApiHandler(w http.ResponseWriter, r *http.Request) {
	// 初始化status_data["status"]="error"、status_data["code"]="1001"、status_data["doit"]="NO_KEY"、status_data["callback"]="INVALID_KEY"
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
			header = jsonToMap(r.URL.Query().Get("header"))
			send_data = jsonToMap(r.URL.Query().Get("data"))
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
			header = jsonToMap(data["header"])
			send_data = jsonToMap(data["data"])
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
			"url":     url,
			"type":    typ,
			"headers": header,
			"data":    send_data,
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
			w.Write([]byte(fmt.Sprintf("Header: [%s]\nData: [%s]", header, send_data)))
			w.Write([]byte(fmt.Sprintf("convertMap(header): [%s]\njson.Marshal(send_data): [%s]", convertMap(header), json_data)))
			return
			// 创建一个自定义的请求对象，设置请求方法为post，请求的url为url，请求的数据为json_data
			req, err := http.NewRequest("POST", url, ioutil.NopCloser(bytes.NewReader(json_data)))
			if err != nil {
				fmt.Println(err)
				return
			}
			// 设置请求的标头，遍历header变量中的键值对，使用http.Header.Set方法设置对应的头部字段和值
			for k, v := range convertMap(header) {
				req.Header.Set(k, v)
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
			
	} else {
			// 如果type不是post，默认使用get方法发送请求，将send_data转换为查询字符串作为url参数
			query := make(map[string]string)
			for k, v := range send_data {
				query[k] = fmt.Sprint(v)
			}
			// 判断url中是否包含?字符，有则使用&连接，否则使用?连接
			if strings.Contains(url, "?") {
				resp, err = http.Get(url + "&" + strEncode(query))
			} else {
				resp, err = http.Get(url + "?" + strEncode(query))
			}
			if err != nil {
				fmt.Println(err)
				return
			}
			defer resp.Body.Close()
			
			// 设置请求的header，遍历header变量中的键值对，使用http.Header.Set方法设置对应的头部字段和值
			for k, v := range header {
				resp.Header.Set(k, fmt.Sprint(v))
			}
	
	}

	
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
	
	// 将status_data转换为json数据并返回
	json_data, err := json.Marshal(status_data)
	if err != nil {
			fmt.Println(err)
			return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_data)
}

// 定义一个辅助函数，将json格式的字符串转换为map类型
func jsonToMap(s string) map[string]interface{} {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
			fmt.Println(err)
			return nil
	}
	return m
}

// 定义一个函数，将map[string]interface{}类型转换为map[string]string类型
func convertMap(m map[string]interface{}) map[string]string {
	// 定义一个新的map[string]string类型的变量
	n := make(map[string]string)

	// 遍历原来的map，对每个键值对进行转换
	for k, v := range m {
		// 判断值的类型，使用不同的方法转换为字符串
		switch v := v.(type) {
		case string:
			// 如果是字符串，直接赋值
			n[k] = v
		case int:
			// 如果是整数，使用strconv.Itoa转换为字符串
			n[k] = strconv.Itoa(v)
		case bool:
			// 如果是布尔值，使用strconv.FormatBool转换为字符串
			n[k] = strconv.FormatBool(v)
		case float64:
			// 如果是浮点数，使用strconv.FormatFloat转换为字符串
			n[k] = strconv.FormatFloat(v, 'f', -1, 64)
		default:
			// 如果是其他类型，使用fmt.Sprint转换为字符串
			n[k] = fmt.Sprint(v)
		}
	}

	// 返回新的map
	return n
}

func strEncode(query map[string]string) string {
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