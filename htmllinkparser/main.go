package main

import (
	"fmt"
	"os"
	"flag"
	"io/ioutil"
	"strings"

	"golang.org/x/net/html"
)

type Links struct {
	Href string
	Text string
}

func Parse_links(n *html.Node, links []Links) ([] Links){
	if n.Type == html.ElementNode && n.Data == "a"{
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				text := ""
				parse_text(n,&text)
				links = append(links, Links{attr.Val,text})
				break
			}
		}
		return links
	
	}
	for i := n.FirstChild; i != nil; i=i.NextSibling {
		 links = Parse_links(i, links)
	}

	return links
}

func parse_text(n *html.Node, t *string) {
	if n.Type == html.TextNode {
		 *t += n.Data
	}
	for i := n.FirstChild; i != nil; i= i.NextSibling {
		parse_text(i,t)
	}
}
func errExit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}


func main() {
	html_file := flag.String("file","index.html","pass the file path which contains html content")
	flag.Parse()

	f, err := ioutil.ReadFile(*html_file)
	if err != nil {
		errExit("Unable to read HTML file!!!")
	}

	html_data := string(f)
	html_root, err := html.Parse(strings.NewReader(html_data))

	if err != nil {
		errExit("Unable to Parse HTML!!!")
	}

	

	fmt.Println(Parse_links(html_root, make([]Links, 0)))
}