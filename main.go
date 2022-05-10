package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/anaskhan96/soup"
	flag "github.com/ogier/pflag"
)

func main() {
	var m string
	var r string
	var s bool
	var c int
	var v bool
	var u bool
	flag.StringVarP(&m, "site-map", "m", "", "url list")
	flag.StringVarP(&r, "responce-file", "r", "", "responce file")
	flag.BoolVarP(&s, "silent", "s", false, "silent mode")
	flag.BoolVarP(&v, "verbose", "v", false, "verbose mode")
	flag.BoolVarP(&v, "urlparams", "u", false, "urlparams only")
	flag.IntVarP(&c, "concurrency", "c", 20, "concurrency")
	flag.Parse()
	url := flag.Arg(0)

	const banner = `
	_____     _ _ ____                               
	|  ___|_ _| | |  _ \ __ _ _ __ __ _ _ __ ___  ___ 
	| |_ / _' | | | |_) / _' | '__/ _' | '_ ' _ \/ __|
	|  _| (_| | | |  __/ (_| | | | (_| | | | | | \__ \
	|_|  \__,_|_|_|_|   \__,_|_|  \__,_|_| |_| |_|___/  v0.0.5											  
   	by: mk990
`
	if !s {
		println(banner)
	}

	parameter := []string{}
	if url != "" {
		parameter = append(parameter, getHttpParams(url)...)
		printParams(parameter)
		os.Exit(0)
	}

	if r != "" {
		file, err := ioutil.ReadFile(r)
		if err != nil {
			log.Fatal(err)
		}
		printParams(getTextParams(string(file)))
		os.Exit(0)
	}

	if m != "" {
		getSitemapParams(m)
	}

	sc := bufio.NewScanner(os.Stdin)
	var wg sync.WaitGroup
	ch := make(chan string)
	for i := 0; i < c; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range ch {
				parameter = append(parameter, getHttpParams(item)...)
			}
			if len(parameter) != 0 {
				printParams(parameter)
			}
		}()
	}
	for sc.Scan() {
		if sc.Text() != "" {
			if v {
				println(sc.Text())
			}
			if u {
				continue
			}

			ch <- sc.Text()
		}
	}
	close(ch)
	wg.Wait()
}

func getSitemapParams(url string) {
	res, _ := httpReq(url, "GET")
	doc := soup.HTMLParse(res)
	links := doc.FindAll("loc")
	for _, link := range links {
		fmt.Println(link.Text())
	}
}

func getHttpParams(url string) []string {
	page, urlparams := httpReq(url, "GET")
	parameter := getTextParams(page)
	parameter = append(parameter, urlparams...)
	return uniqueArray(parameter)
}

func getTextParams(text string) []string {
	parameter := []string{}
	doc := soup.HTMLParse(text)
	links := doc.FindAll("input")
	for _, link := range links {
		name := link.Attrs()["name"]
		id := link.Attrs()["id"]
		if name != "" {
			parameter = append(parameter, name)
		}
		if id != "" {
			parameter = append(parameter, id)
		}
	}

	js := ""
	jre := regexp.MustCompile(`(< *script *.*?>)((?:.*?\r?\n?)*?)(<\/script>)`)
	jss := jre.FindAllStringSubmatch(text, -1)
	for _, i := range jss {
		js += "\n\n" + i[2]
	}
	if isValidJSON(text) {
		js = text
	}
	text = strings.ReplaceAll(js, ";", ";\n")
	text = strings.ReplaceAll(text, "{", "{\n")
	text = strings.ReplaceAll(text, "[", "[\n")
	// fmt.Println(text)

	// js var
	re := regexp.MustCompile(`[^\w](let |var |const )([a-z|A-Z]*?)( *\=)`)
	jv := re.FindAllStringSubmatch(text, -1)
	for _, i := range jv {
		parameter = append(parameter, i[2])
	}

	// json
	reg := regexp.MustCompile(`"*([\w]+?)"*\s*:`)
	jj := reg.FindAllStringSubmatch(text, -1)
	for _, item := range jj {
		parameter = append(parameter, item[1])
	}

	return parameter
}

func getURLParameter(res *http.Response) []string {
	parameter := []string{}
	for key := range res.Request.URL.Query() {
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

func printParams(params []string) {
	for _, p := range uniqueArray(params) {
		fmt.Println(p)
	}
}

func CheckError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func isValidJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func httpReq(url string, method string) (string, []string) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body), getURLParameter(resp)[:]
}
