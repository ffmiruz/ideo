package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	json "github.com/json-iterator/go"
)

func main() {
	cfg := &Config{}
	err := cfg.loadConfig("assets/config.json")
	if err != nil {
		log.Println(err)
	}
	tpl, err := ioutil.ReadFile("assets/template.html")
	if err != nil {
		log.Fatal(err)
	}
	files, err := getFile("assets/content")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		err := writeHtml(f.Name(), tpl, cfg)
		if err != nil {
			log.Println(err)
		}
	}
	err = writeCss(cfg)
	if err != nil {
		log.Println(err)
	}
}

type Config struct {
	// Site base url
	BASE string `json:"base"`
	// Name or Title of the site
	NAME string `json:"name"`
	// Folder to write to. Default to repo root folder
	Outpath string `json:"outpath"`
}

func (c *Config) loadConfig(path string) error {
	b, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(b, c)
}

func getFile(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	//sort.Slice(list, func(i, j int) bool { return list[i].ModTime() > list[j].ModTime() })
	return filterMd(list), nil
}

func filterMd(a []os.FileInfo) []os.FileInfo {
	b := a[:0]
	for _, x := range a {
		if strings.HasSuffix(x.Name(), ".md") {
			b = append(b, x)
		}
	}
	return b
}

func writeHtml(file string, tpl []byte, c *Config) error {
	bufin, err := ioutil.ReadFile("assets/content/" + file)
	if err != nil {
		return err
	}
	body := markdown.ToHTML(bufin, nil, nil)

	path := file[:len(file)-3] + ".html"
	bufout := bytes.Replace(tpl, []byte("{{CONTENT}}"), body, 1)
	bufout = bytes.Replace(bufout, []byte("{{BASE}}"), []byte(c.BASE), -1)
	bufout = bytes.Replace(bufout, []byte("{{NAME}}"), []byte(c.NAME), -1)
	if strings.Contains(path, "index.html") {
		bufout = bytes.Replace(bufout, []byte("{{PATH}}"), []byte("index"), -1)
	} else {
		bufout = bytes.Replace(bufout, []byte("{{PATH}}"), []byte(path), -1)
	}
	err = ioutil.WriteFile(c.Outpath+path, bufout, 0644)
	if err != nil {
		return err
	}
	return err
}

func writeCss(c *Config) error {
	bufin, err := ioutil.ReadFile("assets/style.css")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(c.Outpath+"style.css", bufin, 0644)
	if err != nil {
		return err
	}
	return err
}
