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
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func main() {
	var format string
	var p string

	flag.StringVar(&format, "time", time.Now().Format("2006.01.02"), "206.01.02")
	flag.StringVar(&p, "p", "", "unzip file")
	flag.Parse()

	path := "./"
	fileName := "clash-yaml"
	if p != "" {
		log.Println("p file name=>", p)
		if !strings.HasSuffix(p, ".zip") {
			log.Println("[error]", "not a zip file")
			return
		}

		exec.Command("Bandizip.exe", "x", "-aoa", "-y", "-p:59672", "-o:"+path+fileName, p).Run()
		walkDir := path + fileName
		m := make(map[int64]string)
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
		os.RemoveAll(path + p)
		return
	}

	f, err := os.ReadDir(path)
	if err != nil {
		log.Println(err)
	}
	for _, entry := range f {
		b := strings.Contains(entry.Name(), "每日更换链接更新"+format) && strings.Contains(entry.Name(), ".zip")
		if b {
			log.Println("t file name=>", entry.Name())
			exec.Command("Bandizip.exe", "x", "-aoa", "-y", "-p:59672", "-o:"+path+fileName, path+entry.Name()).Run()
			walkDir := path + fileName
			m := make(map[int64]string)
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
			os.RemoveAll(path + entry.Name())
			return
		}
	}
}

func uploadFile(m map[int64]string) {
	url := "http://127.0.0.1:8001/upload"

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
		log.Println("uploadFileName:", s)
		copyFile("upload[]", s, bodyWriter, &i)
		i += 1
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	request, _ := http.NewRequest("POST", url, bodyBuffer)
	request.Header.Add("Content-Type", contentType)

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
