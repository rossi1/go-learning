package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// findFastestWeb will take a slice of url's represented as a []string,
// and return back the response of the webserver that returned the
// complete page first as respData
func findFastestWeb(url []string) []respData {
	var r []respData
	var wg sync.WaitGroup
	var mu sync.Mutex // mutex to protect the first to update the found variable.

	// Loop over all the url's given as input, and start a go routine
	// with a http.Get for each of them.
	for _, u := range url {
		wg.Add(1)
		go func(u string) {
			t := time.Now()
			resp, err := http.Get(u)
			if err != nil {
				log.Println("error: http.Get for one", err)
			}

			totalTime := time.Since(t)

			mu.Lock()
			r = append(r, respData{time: totalTime, resp: resp})
			mu.Unlock()

			wg.Done()
			return

		}(u)

	}

	wg.Wait()
	return r
}

type respData struct {
	time time.Duration
	resp *http.Response
}

func main() {

	urls := []string{"https://dagbladet.no", "https://aftenposten.no", "https://vg.no"}

	r := findFastestWeb(urls)

	for i, v := range r {
		fmt.Printf("%v : %v : , and in time it took %v\n", i+1, v.resp.Request.URL, v.time)
	}

}
