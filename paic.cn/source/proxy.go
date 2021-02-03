package source

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	logrus "github.com/sirupsen/logrus"
)

var loga = logrus.New()

// RouteHandler 根据请求URL转发到服务
func RouteHandler(w http.ResponseWriter, r *http.Request) {

	path := r.RequestURI
	log.Println(path)
	// 根据请求URL获取服务名
	serviceName := strings.Split(path, "/")[1]
	// 去掉URL上的服务名
	path = strings.Replace(path, "/"+serviceName, "", 1)
	// 从consul注册中心获取服务地址
	servs := GetService(serviceName)
	// 获取随机地址
	url := randMapValue(servs)

	if url == "" {
		logrus.Error("url is empty!")
		return
	}

	// 请求转发
	resp, err := http.Get("http://" + url + path)
	if err != nil {
		log.Println(err.Error())
	}
	// 延迟释放请求
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
	}

	loga.WithFields(logrus.Fields{
		"url":  url,
		"path": path,
	}).Info(string(body))

	log.Println(string(body))
	fmt.Fprintln(w, string(body))

}

// 随机获取地址
func randMapValue(m map[string]struct{}) string {
	for key := range m {
		return key
	}
	return ""
}
