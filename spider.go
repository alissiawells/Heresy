package main

import (
	"fmt"
	"os"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"io"
	"io/ioutil"
	"github.com/kljensen/snowball"
)

func Stemmer(word string, language string, stopwords bool) string {
	stemmed, err := snowball.Stem(word, language, stopwords)
	if err != nil {
		fmt.Println(err)
	}
	return stemmed
}

func FuncKeywords(word string, Args []string) string {
	keywords := "\"(" + word + "*"
	for i := 5; i < len(os.Args); i++ {
		keywords = keywords + ")|(" + Stemmer(os.Args[i], os.Args[3], true) + "*"
	}
	keywords = keywords + ")\""
	return keywords
}

type Parser interface {
	// returns the body of URL and a slice of URLs found on that page
	Parse(url string) (body string, urls []string, err error)
}

type Scraper struct {
	// safe to use concurrently
	visited map[string]bool
	muxLock sync.Mutex
}

func (s *Scraper) Parse(url string) (body string, urls []string, err error) {
	// excludes visited URL
	s.muxLock.Lock()
	defer s.muxLock.Unlock()

	_, ok := s.visited[url] 
	if ok {
		return
	}
	s.visited[url] = true

	// parsing the URL body
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	} 
	
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	body = string(bytes)
	// finds URLs in body
	UrlRegexp := regexp.MustCompile(`(http|https):\/\/([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`)
	urls = UrlRegexp.FindAllString(body, -1)
	return
}

func Scrape(keywords string, url string, depth int, parser Parser) {
	// recursive scraping in parallel
	var wg sync.WaitGroup

	if depth <= 0 {
		return
	}
	body, urls, err := parser.Parse(url)

	// finding key words
	matched, err := regexp.MatchString(keywords, body)
	if err != nil {
		fmt.Println(err)
		return
	} 

	if matched {
		fmt.Println("found:", url)
		file, err := os.Create("data.txt")
		if err != nil {
			fmt.Println(err)
			return
		}

		n, err := io.WriteString(file, body)
		fmt.Println("output >> data.txt")
		if err != nil {
			fmt.Println(n, err)
			return
		}
		file.Close()
	}

	wg.Add(len(urls))
	for _, u := range urls {
		go func(url string) {
			defer wg.Done()
			Scrape(keywords, url, depth-1, parser)
		}(u)
	}
	wg.Wait()
	return
}

func main() {

	var keywords string
	depth, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("%s: the value of depth after the URL must be integer\n", os.Args[0])
	}

	word := Stemmer(os.Args[4], os.Args[3], true)

	if len(os.Args) > 5 {
		keywords = FuncKeywords(word, os.Args)
	} else {
		keywords = "\"" + word + "*" + "\""
	}

	parser := Scraper {
		visited: make(map[string]bool),
	}

	Scrape(keywords, os.Args[1], depth, &parser)
}
