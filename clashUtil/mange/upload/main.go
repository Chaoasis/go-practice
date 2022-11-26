package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func main() {
	s := flag.String("time", time.Now().Format("2006.01.02"), "put time")
	flag.Parse()
	path := "./"
	walkDir := path + "每日更换链接更新，直到新的视频播出"
	m := make(map[int64]string)
	format := *s
	log.Println("time is ", format)
	filepath.WalkDir(walkDir, func(p string, d fs.DirEntry, err error) error {
		name := d.Name()
		if strings.Contains(name, format) && (strings.Contains(name, ".yml") || strings.Contains(name, ".yaml")) {
			info, _ := d.Info()
			modTime := info.ModTime()
			m[modTime.UnixNano()] = p
		}
		return nil
	})
	uploadFile(m)
}

func uploadFile(m map[int64]string) {
	url := "http://127.0.0.1:9001/upload"

	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	//------------------
	var keys []int64
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	i := 1
	for _, key := range keys {
		s := m[key]
		log.Println(s)
		copyFile("upload[]", s, bodyWriter, &i)
		i += 1
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	request, _ := http.NewRequest("POST", url, bodyBuffer)
	request.Header.Add("Content-Type", contentType)
	request.Header.Add("appKey", "hello world")

	resp, _ := http.DefaultClient.Do(request)
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)

	log.Println(resp.Status, string(response))

}

func copyFile(key string, path string, bodyWriter *multipart.Writer, i *int) {
	formFile, _ := bodyWriter.CreateFormFile(key, fmt.Sprintf("%d.yaml", *i))
	file, _ := os.Open(path)
	defer file.Close()
	io.Copy(formFile, file)
}
