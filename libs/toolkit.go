package libs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"reflect"
	"strings"
)

type (
	CookieJson []struct {
		Domain         string  `json:"domain"`
		ExpirationDate float64 `json:"expirationDate,omitempty"`
		HostOnly       bool    `json:"hostOnly"`
		HTTPOnly       bool    `json:"httpOnly"`
		Name           string  `json:"name"`
		Path           string  `json:"path"`
		SameSite       string  `json:"sameSite"`
		Secure         bool    `json:"secure"`
		Session        bool    `json:"session"`
		StoreID        string  `json:"storeId"`
		Value          string  `json:"value"`
		ID             int     `json:"id"`
	}
	RequestInfo struct {
		CookieFile string
		Headers    interface{}
		Method     string
		Url        string
		Data       interface{}
	}
)

func fillBySettings(st interface{}, settings map[string]interface{}) {

	if reflect.TypeOf(st).Kind() != reflect.Ptr {
		log.Println("the first param should be a pointer to the struct type.")
	}
	// Elem() 获取指针指向的值
	if reflect.TypeOf(st).Elem().Kind() != reflect.Struct {
		log.Println("the first param should be a pointer to the struct type.")
	}

	if settings == nil {
		log.Println("settings is nil.")
	}

	var (
		field reflect.StructField
		ok    bool
	)

	for k, v := range settings {
		if field, ok = (reflect.ValueOf(st)).Elem().Type().FieldByName(k); !ok {
			continue
		}
		if field.Type == reflect.TypeOf(v) {
			vstr := reflect.ValueOf(st)
			vstr = vstr.Elem()
			vstr.FieldByName(k).Set(reflect.ValueOf(v))
		}

	}
	return
}

func Sign(signData RequestInfo, callBack func()) {
	pwd, _ := os.Getwd()
	jsonFile := pwd + "/cookies/" + signData.CookieFile

	if CheckFile(jsonFile) == false {
		fmt.Printf("%s 不存在", jsonFile)
	}
	data, _ := ioutil.ReadFile(jsonFile)

	var v CookieJson
	if err := json.Unmarshal(data, &v); err != nil {
		panic(err)
	}

	req, _ := http.NewRequest(signData.Method, signData.Url, nil)

	//这里是从json添加多个cookie，所以需要使用cookiejar
	client := &http.Client{}
	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie
	for _, row := range v {
		cookies = append(cookies, &http.Cookie{
			Domain: row.Domain,
			Name:   row.Name,
			Path:   row.Path,
			Value:  row.Value,
		})
	}
	jar.SetCookies(req.URL, cookies)
	client.Jar = jar

	//解析并使用signData.Headers
	st := reflect.TypeOf(signData.Headers)
	sv := reflect.ValueOf(signData.Headers)

	if st.Kind() == reflect.Map {
		for _, key := range sv.MapKeys() {
			strct := sv.MapIndex(key)
			kk := fmt.Sprintf("%s", key.Interface())
			kv := fmt.Sprintf("%s", strct.Interface())
			//判断kk是否等于 UA
			if strings.EqualFold(kk, "UA") {
				kk = "User-Agent"
			}
			req.Header.Add(kk, kv)
		}
	}

	//开始请求
	res, _ := client.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("\n%s 请求结果: %s\n", signData.CookieFile, string(body))

	callBack()
}
