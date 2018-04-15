package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
	resp, err := http.Get(os.Args[1])
	var urlArray []string
	s, _ := ioutil.ReadAll(resp.Body)
	doc, err := html.Parse(strings.NewReader(string(s)))

	dirName := path.Base(resp.Request.URL.Path)

	if err != nil {
		fmt.Println(err.Error())
	}
	if resp.StatusCode == 200 {
		if _, err := os.Stat(dirName); os.IsNotExist(err) {
			os.Mkdir(dirName, os.FileMode(0777))
		}

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
			fmt.Printf("Downloading %d of %d: %s\n", i+1, len(urlArray), urlArray[i])
			downloadFile(urlArray[i], dirName)
		}
	}
}

func downloadFile(picURL string, dir string) {
	fileName := path.Base(picURL)
	if _, err := os.Stat(fmt.Sprintf("%s/%s", dir, fileName)); err == nil {
		fmt.Println("File already downloaded, skipping...")
		return
	}
	file, err := os.Create(fmt.Sprintf("%s/%s", dir, fileName))

	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()

	resp, err := http.Get(picURL)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("%s with %v bytes downloaded\n", fileName, size)
}
