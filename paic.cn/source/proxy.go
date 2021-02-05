package source

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	logrus "github.com/sirupsen/logrus"
)

var loga = logrus.New()

// RouteHandler 根据请求URL转发到服务
func RouteHandler(w http.ResponseWriter, r *http.Request) {
	// 判断请求类型
	var filepath string
	contentType := r.Header.Values("Content-type")[0]
	if strings.Contains(contentType, "multipart/form-data") {
		filepath = saveFile(w, r)
	}

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
	// resp, err := http.PostForm("http://" + url + path,)
	params := map[string]string{}
	resq, err := uploadRequest("http://"+url+path, params, filepath)
	if resq != nil {
		log.Println(resq.Header)
	}
	if err != nil {
		log.Println(err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(resq)
	if err != nil {
		log.Println(err.Error())
	}

	// resp, err := http.Get("http://" + url + path)
	// if err != nil {
	// 	log.Println(err.Error())
	// }
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

// 保存文件到本地
func saveFile(w http.ResponseWriter, r *http.Request) string {

	contentType := r.Header.Values("Content-type")

	loga.Info(contentType)
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return ""
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// save temp files
	tempFile, err := ioutil.TempFile("temp-images", "upload-*"+handler.Filename)
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")
	return tempFile.Name()
}

//参数为本地文件地址
func uploadRequest(url string, params map[string]string, path string) (*http.Request, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 实例化multipart
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 创建multipart 文件字段
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return nil, err
	}
	// 写入文件数据到multipart
	_, err = io.Copy(part, file)
	//将额外参数也写入到multipart
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	//创建请求
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	//不要忘记加上writer.FormDataContentType()，
	//该值等于content-type :multipart/form-data; boundary=xxxxx
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req, nil
}
