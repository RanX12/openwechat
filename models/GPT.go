package models

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	Code    int      `json:"code"`
	Data    DataTime `json:"data"`
	Content string   `json:"content"`
	Sangbo  string   `json:"sangbo"`
}

type BaiduResponse struct {
	ID               string `json:"id"`
	Object           string `json:"object"`
	Created          int64  `json:"created"`
	Result           string `json:"result"`
	IsTruncated      bool   `json:"is_truncated"`
	NeedClearHistory bool   `json:"need_clear_history"`
	Usage            struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type DataTime struct {
	Output string `json:"output"`
}

type GptMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestData struct {
	Model    string       `json:"model"`
	Messages []GptMessage `json:"messages"`
}

//func GPTmsg(content string) string {
//	// Prepare the request data
//	data := RequestData{
//		Model: "gpt-3.5-turbo",
//		Messages: []GptMessage{
//			{
//				Role:    "system",
//				Content: "You are a poetic assistant, skilled in explaining complex programming concepts with creative flair.",
//			},
//			{
//				Role:    "user",
//				Content: "Compose a poem that explains the concept of recursion in programming.",
//			},
//		},
//	}
//
//	// Convert the data to JSON
//	jsonData, err := json.Marshal(data)
//	if err != nil {
//		log.Println("Error: %s", err)
//		return "json Marshal 错误"
//	}
//
//	// Create a new HTTP request
//	req, err := http.NewRequest("POST", "https://service-89bh1w6k-1310716354.sg.apigw.tencentcs.com/v1/chat/completions", bytes.NewBuffer(jsonData))
//	if err != nil {
//		log.Println("Error: %s", err)
//		return "htt new request 错误"
//	}
//
//	// Set the headers
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Authorization", "Bearer "+viper.GetString("GPTkey"))
//
//	// Send the request
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		log.Println("Error: %s", err)
//		return "client Do 错误"
//	}
//	defer resp.Body.Close()
//
//	// Print the response
//	log.Println("Response status:", resp.Status)
//	log.Println(">>>>>>>>>>>resp: ", resp)
//	return resp.Status
//}

// GPT回复
//func GPTmsg(content string) string {
//	fmt.Println(">>>>>>>>>msg", content)
//	//apiUrl := "https://api.lolimi.cn/API/AI/mfcat3.5.php?type=json"
//	apiUrl := "https://luckycola.com.cn/ai/openwxyy"
//
//	// 构造请求参数
//	data := map[string]string{
//		"ques":   content,
//		"appKey": viper.GetString("luckycola.appKey"),
//		"uid":    viper.GetString("luckycola.uid"),
//	}
//
//	jsonData, err := json.Marshal(data)
//	if err != nil {
//		log.Fatalf("JSON encoding failed: %s", err)
//	}
//
//	//params.Set("sx", "你是一个全能专家，会 Java、Ruby、PHP、Golang 等各种编程语言")
//	//params.Set("key", viper.GetString("GPTkey"))
//
//	log.Println(">>>>>>>>>>params: ", jsonData)
//	// 发送HTTP GET请求
//	//resp, err := http.Get(apiUrl + "?" + params.Encode())
//	//if err != nil {
//	//	log.Println(err)
//	//	return "http Get 未知错误！"
//	//}
//	//defer resp.Body.Close()
//	// 发送HTTP GET请求
//	resp, err := http.Post(apiUrl, "application/json", bytes.NewReader(jsonData))
//	if err != nil {
//		log.Println(err)
//		return "http Get 未知错误！"
//	}
//	defer resp.Body.Close()
//
//	// Check the HTTP status code
//	if resp.StatusCode != 200 {
//		log.Printf("Received HTTP error code %d\n", resp.StatusCode)
//		return "HTTP error!"
//	}
//
//	// 解析JSON响应数据
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		log.Println(err)
//		return "io ReadAll 未知错误！"
//	}
//	return string(body)
//	//fmt.Println("string body: ", string(body))
//	//var msg Response
//	//log.Println(">>>>>>body: ", body)
//	//log.Println(">>>>>>msg: ", &msg)
//	//err = json.Unmarshal(body, &msg)
//	//if err != nil {
//	//	fmt.Println(err)
//	//	return "json Unmarshal 未知错误！"
//	//}
//	//return msg.Data.Output
//}

func GPTmsg(content string) string {
	url := "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/completions?access_token=" + GetAccessToken()

	jsonData := fmt.Sprintf(`{
		"messages":
			[{
				"role":"user",
				"content":"%s"
			}],
		"system":"你是一个编程高手，精通 Java、Go、Ruby、C 等各种编程语言，还是一个情感专家，善解人意的 AI 助手，请你的回答不要太过于官方了，人性化一点，甚至可以调皮一点！"
	}`, content)

	payload := strings.NewReader(jsonData)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		fmt.Println(err)
		return "http.NewRequest error"
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "client.Do error"
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "io.ReadAll error"
	}
	fmt.Println(string(body))
	var resp BaiduResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println(err)
		return "json.Unmarshal error"
	}
	return resp.Result
}

/**
 * 使用 AK，SK 生成鉴权签名（Access Token）
 * @return string 鉴权签名信息（Access Token）
 */
func GetAccessToken() string {
	url := "https://aip.baidubce.com/oauth/2.0/token"
	postData := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", viper.GetString("baidu-ai.API_KEY"), viper.GetString("baidu-ai.SECRET_KEY"))
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(postData))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	accessTokenObj := map[string]string{}
	json.Unmarshal([]byte(body), &accessTokenObj)
	return accessTokenObj["access_token"]
}
