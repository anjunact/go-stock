package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

type NotFoundError struct {
	Message string
}

func (e NotFoundError) Error() string {
	return e.Message
}

type RemoteError struct {
	Host string
	Err  error
}

func (e *RemoteError) Error() string {
	return e.Err.Error()
}

var UserAgent = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/29.0.1541.0 Safari/537.36"

// HttpGet gets the specified resource. ErrNotFound is returned if the
// server responds with status 404.
func HttpGet(client *http.Client, url string, header http.Header) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	for k, vs := range header {
		req.Header[k] = vs
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, &RemoteError{req.URL.Host, err}
	}
	if resp.StatusCode == 200 {
		return resp.Body, nil
	}
	resp.Body.Close()
	if resp.StatusCode == 404 { // 403 can be rate limit error.  || resp.StatusCode == 403 {
		err = NotFoundError{"Resource not found: " + url}
	} else {
		err = &RemoteError{req.URL.Host, fmt.Errorf("get %s -> %d", url, resp.StatusCode)}
	}
	return nil, err
}

// HttpGetBytes gets the specified resource. ErrNotFound is returned if the server
// responds with status 404.
func HttpGetBytes(client *http.Client, url string, header http.Header) ([]byte, error) {
	rc, err := HttpGet(client, url, header)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return ioutil.ReadAll(rc)
}

// HttpGetToFile gets the specified resource and writes to file.
// ErrNotFound is returned if the server responds with status 404.
func HttpGetToFile(client *http.Client, url string, header http.Header, fileName string) error {
	rc, err := HttpGet(client, url, header)
	if err != nil {
		return err
	}
	defer rc.Close()


	os.MkdirAll(path.Dir(fileName), os.ModePerm)
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, rc)
	return err
}


var img = regexp.MustCompile(`href=\"javascript:goView\((\d+)`)
var imgPattern = regexp.MustCompile(`id="mainImage" src=\"../upload(.*?).jpg`)
var totalTask int

func download(url string, num chan bool) {
	url = strings.TrimPrefix(url, `href="javascript:goView(`)
	page := "http://www.gdweb.co.kr/main/koreaWebView.asp?idx=%s&url=koreaWeb.asp"
	t, err := HttpGetBytes(&http.Client{}, fmt.Sprintf(page, url), nil)


	if err != nil {
		log.Fatalf("获取页面失败：%v", err)
	}
	matches := imgPattern.FindAll(t, -1)
	for _, match := range matches {
		url = "http://www.gdweb.co.kr" + strings.TrimPrefix(string(match), `id="mainImage" src="..`)
		log.Printf("正在下载：%s", url)
		err := HttpGetToFile(&http.Client{}, url, nil, "pics/"+path.Base(url))
		if err != nil {
			log.Printf("图片下载失败（%s）：%v", url, err)
		}
	}
	totalTask--
	<-num
}

func main() {
	// 控制同时下载数量
	num := make(chan bool, 5)


	// 主线程爬取页面，子线程下载图片
	//baseUrl := "http://nvmingxing.net/hotness/%d/"
	//abaseUrl := "http://www.gdweb.co.kr/main/koreaWebView.asp?idx=8200&url=koreaWeb.asp"
	baseUrl := "http://www.gdweb.co.kr/main/koreaWeb.asp?idx=&url=index.asp&lpage=124&page=%d"
	for i := 2; i < 124; i++ {
		log.Printf("抓取页面：%d", totalTask)
		data, err := HttpGetBytes(&http.Client{}, fmt.Sprintf(baseUrl, i+1), nil)
		if err != nil {
			log.Fatalf("获取页面失败（%d）：%v", i, err)
		}
		matches := img.FindAll(data, -1)
		for _, match := range matches {
			totalTask++
			num <- true
			go download(string(match), num)
		}
	}
}
