package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	flag "github.com/ogier/pflag"
)

func main() {
	var l string
	var m string
	var r string
	var s bool
	var c int
	flag.StringVarP(&l, "list", "l", "", "url list")
	flag.StringVarP(&m, "site-map", "m", "", "url list")
	flag.StringVarP(&r, "responce", "r", "", "responce file")
	flag.BoolVarP(&s, "silent", "s", false, "silent mode")
	flag.IntVarP(&c, "concurrency", "c", 20, "concurrency")
	flag.Parse()
	url := flag.Arg(0)

	const banner = `
	_____     _ _ ____                               
	|  ___|_ _| | |  _ \ __ _ _ __ __ _ _ __ ___  ___ 
	| |_ / _' | | | |_) / _' | '__/ _' | '_ ' _ \/ __|
	|  _| (_| | | |  __/ (_| | | | (_| | | | | | \__ \
	|_|  \__,_|_|_|_|   \__,_|_|  \__,_|_| |_| |_|___/  v0.0.1											  
   	by: mk990
`
	if !s {
		println(banner)
	}

	parameter := []string{}
	if url != "" {
		parameter = append(parameter, getHttpParams(url)[:]...)
	}

	if r != "" {
		file, err := ioutil.ReadFile(r)
		if err != nil {
			log.Fatal(err)
		}
		buf := bytes.NewBuffer(file)
		parameter = append(parameter, getTextParams(buf)[:]...)
	}
	if m != "" {
		getSitemapParams(m)
	}
	// sc := bufio.NewScanner(os.Stdin)
	// for sc.Scan() {
	// 	if sc.Text() != "" {
	// 		parameter = append(parameter, getHttpParams(sc.Text())[:]...)
	// 	}
	// }

	if len(parameter) == 0 {
		os.Exit(0)
	}

	for _, p := range uniqueArray(parameter) {
		fmt.Println(p)
	}

}

func getSitemapParams(url string) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		// return []string{}
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("loc").Each(func(i int, s *goquery.Selection) {
		url := s.Text()
		fmt.Println(url)
	})
}

func getHttpParams(url string) []string {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return []string{}
	}
	defer res.Body.Close()
	parameter := getTextParams(res.Body)
	parameter = append(parameter, getURLParameter(res)[:]...)
	return uniqueArray(parameter)
}

func getTextParams(text io.Reader) []string {
	parameter := []string{}

	doc, err := goquery.NewDocumentFromReader(text)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("form input").Each(func(i int, s *goquery.Selection) {
		title, _ := s.Attr("name")
		if title != "" {
			parameter = append(parameter, title)
		}
	})

	return parameter
}

func getURLParameter(res *http.Response) []string {
	parameter := []string{}
	for key, _ := range res.Request.URL.Query() {
		parameter = append(parameter, key)
	}
	return parameter
}

func uniqueArray(e []string) []string {
	r := []string{}

	for _, s := range e {
		if !inArray(r[:], s) {
			r = append(r, s)
		}
	}
	return r
}

func inArray(e []string, c string) bool {
	for _, s := range e {
		if s == c {
			return true
		}
	}
	return false
}
