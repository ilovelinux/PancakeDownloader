package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: ./main <URL>")
		os.Exit(1)
	}
	var threadURL string
	resp, err := http.Get(os.Args[1])
	var urlArray []string
	s, _ := ioutil.ReadAll(resp.Body)
	doc, err := html.Parse(strings.NewReader(string(s)))

	threadURL = os.Args[1]
	DirName := path.Base(threadURL)

	if err != nil {
		fmt.Println(err.Error())
	}
	if resp.StatusCode == 200 {
		os.Mkdir(DirName, os.FileMode(0777))

		var f func(*html.Node)
		f = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "a" {
				for _, a := range n.Attr {
					if a.Key == "class" && a.Val == "fileThumb" {
						for _, y := range n.Attr {
							if y.Key == "href" {
								urlArray = append(urlArray, strings.Replace(y.Val, "//", "http://", 1))
							}
						}
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(doc)
		for i := 0; i < len(urlArray); i++ {
			fmt.Println(urlArray[i])
			downloadFile(urlArray[i], DirName)
		}
	}
}

func downloadFile(picURL string, dir string) {
	fileURL, err := url.Parse(picURL)

	if err != nil {
		fmt.Println(err.Error())
	}

	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[2]
	file, err := os.Create(fmt.Sprintf("%s/%s", dir, fileName))

	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()

	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := check.Get(picURL)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	//fmt.Println(resp.Status)

	size, err := io.Copy(file, resp.Body)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("%s with %v bytes downloaded\n", fileName, size)
}
