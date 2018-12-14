package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/tkivisik/go-wild-with-web/web5/hn"
)

// TASK
// * Get Hacker Noon news.
// * Keep original order
// * Always print 30, not more
// * Cache for speed

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	fmt.Println("Sweet")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var client hn.Client
		ids, err := client.TopItems()
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}
		//var stories []item
		stories := [45]item{}
		mutex := sync.Mutex{}
		wg := sync.WaitGroup{}
		for i, id := range ids {
			if i >= 45 {
				break
			}
			go func(ii, id int, storiesp *[45]item) {
				fmt.Println(ii, id)
				wg.Add(1)
				hnItem, err := client.GetItem(id)
				if err != nil {
					return
				}
				item := parseHNItem(hnItem)
				if isStoryLink(item) {
					fmt.Println("Item: ", item)
					mutex.Lock()
					(*storiesp)[ii] = item
					mutex.Unlock()
				}
				wg.Done()
			}(i, id, &stories)
		}
		wg.Wait()
		fmt.Println("Items:", stories)
		storiesSlice := []item(stories[:])
		data := templateData{
			Stories: storiesSlice,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
