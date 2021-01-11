package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const HTML5 = `!DOCTYPE html>`
const HttpMinStatus = 200
const HttpMaxStatus = 299

type Output struct {
	HtmlVersion string
	PageTitle string
	Headings map[string]int
	Links map[string]int
	Login int
}

func getHTMLVersion(url string) string {

	fmt.Printf("HTML code of %s ...\n", url)

	resp, err := http.Get(url)
	// handle the error if there is one
	if err != nil {
		panic(err)
	}

	// do this now so it won't be forgotten
	defer resp.Body.Close()

	// reads html as a slice of bytes
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	s:=strings.Split(string(html), "<")

	matchedV5, _:=regexp.MatchString(HTML5, s[1])
	matchedV4, _:=regexp.MatchString(`.HTML 4.`, s[1])
	matchedVX, _:=regexp.MatchString(`.XHTML.`, s[1])

	if matchedV5 {
		return fmt.Sprint("HTML Version 5")
	}else if matchedV4 {
		return fmt.Sprint("HTML Version 4")
	}else if matchedVX {
		return fmt.Sprint("XHTML Version")
	}else {
		return fmt.Sprint("Version uncheckable")
	}
}

func scrapeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method!="GET"{
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	params:=r.URL.Query()
	var URL= params["url"][0]
	if URL==""{
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	output:=&Output{
		HtmlVersion: "",
		PageTitle: "",
		Links: make(map[string]int),
		Headings: make(map[string]int),
	}

	output.HtmlVersion = getHTMLVersion(URL)
	fmt.Printf("Version: %s\n",output.HtmlVersion)

	// Instantiate default collector
	c := colly.NewCollector()

	c.OnHTML("title", func(e *colly.HTMLElement) {
		fmt.Printf("Title: %s\n", e.Text)
		if e.Name != ""{
			output.PageTitle=e.Text
		}
	})

	c.OnHTML("h1", func(e *colly.HTMLElement) {
		// Print link
		fmt.Printf("H1 found: %s\n", e.Name)
		if e.Name!="" {
			output.Headings[e.Name]++
		}
	})
	c.OnHTML("h2", func(e *colly.HTMLElement) {
		// Print link
		fmt.Printf("H2 found: %s\n", e.Name)
		if e.Name!="" {
			output.Headings[e.Name]++
		}
	})
	c.OnHTML("h3", func(e *colly.HTMLElement) {
		// Print link
		fmt.Printf("H3 found: %s\n", e.Name)
		if e.Name!="" {
			output.Headings[e.Name]++
		}
	})
	c.OnHTML("h4", func(e *colly.HTMLElement) {
		// Print link
		fmt.Printf("H4 found: %s\n", e.Name)
		if e.Name!="" {
			output.Headings[e.Name]++
		}
	})
	c.OnHTML("h5", func(e *colly.HTMLElement) {
		// Print link
		fmt.Printf("H5 found: %s\n", e.Name)
		if e.Name!="" {
			output.Headings[e.Name]++
		}
	})
	c.OnHTML("h6", func(e *colly.HTMLElement) {
		// Print link
		fmt.Printf("H6 found: %s\n", e.Name)
		if e.Name!="" {
			output.Headings[e.Name]++
		}
	})



	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)

		// Parse link found on page
		val, _ := url.Parse(link)
		if val.Hostname()=="" {
			output.Links["internal"]++
		} else {
			output.Links["external"]++

			resp, err := http.Get(link)
			// handle the error if there is one
			if err != nil || !(resp.StatusCode >= HttpMinStatus && resp.StatusCode <= HttpMaxStatus){
				output.Links["inaccessible"]++
			}

			defer resp.Body.Close()
		}
		if val.Path=="/login"{
			output.Login++
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping
	c.Visit(URL)

	// dump results
	b, err :=  json.Marshal(output)
	if err != nil {
		log.Println("failed to serialize response:", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func main(){
	addr := ":7171"

	http.HandleFunc("/", scrapeHandler)

	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}