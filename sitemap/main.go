package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Links struct {
	Href string
	Text string
}

func Parse_links(n *html.Node, links []Links) []Links {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				text := ""
				parse_text(n, &text)
				links = append(links, Links{attr.Val, text})
				break
			}
		}
		return links

	}
	for i := n.FirstChild; i != nil; i = i.NextSibling {
		links = Parse_links(i, links)
	}

	return links
}

func parse_text(n *html.Node, t *string) {
	if n.Type == html.TextNode {
		*t += n.Data
	}
	for i := n.FirstChild; i != nil; i = i.NextSibling {
		parse_text(i, t)
	}
}
func errExit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func extractSiteMap(start_url string, depth uint) *map[string]bool {
	var url_set map[string]bool = make(map[string]bool, 0)

	d_start := strings.Index(start_url, "/") + 2
	d_end := strings.Index(start_url[d_start:], "/")
	//domain := start_url[d_start:d_end+1]
	scheme_len := len(start_url[:d_start])
	scheme_domain := start_url[:d_end+scheme_len]

	var recur func(url *string, depth uint)

	recur = func(url *string, depth uint) {
		if depth != 0 {
			if _, ok := url_set[*url]; ok {
				return
			}

			url_set[*url] = true

			resp, err := http.Get(*url)

			if err != nil {
				fmt.Println("Request failed!!")
			}

			defer resp.Body.Close()

			resp_data, err := io.ReadAll(resp.Body)

			if err != nil {
				fmt.Println("Error reading html response!!")
			}

			parsed_html, err := html.Parse(strings.NewReader(string(resp_data)))

			if err != nil {
				fmt.Println("Unable to parse html tree!!")
			}

			resp_links := Parse_links(parsed_html, make([]Links, 0))
			for _, link := range resp_links {
				n_url := link.Href
				if !strings.HasPrefix(n_url, scheme_domain) {
					if n_url[0] == '/' {
						n_url = scheme_domain + n_url
					} else if n_url[0] != 'h' {
						n_url = scheme_domain + "/" + n_url
					} else {
						continue
					}
				}

				recur(&n_url, depth-1)
			}

		}
	}
	recur(&start_url, depth)

	return &url_set
}

func build_sitemap_xml(url_set *map[string]bool) string {
	const xml_header string = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"

	type Url_t struct {
		Loc string `xml:"loc"`
	}

	type Urlset struct {
		XMLName   xml.Name `xml:"urlset"`
		Namespace string   `xml:"xmlns,attr"`
		Url       []Url_t  `xml:"url"`
	}

	urls := make([]Url_t, 0)
	for k := range *url_set {
		urls = append(urls, Url_t{k})
	}

	urlset := &Urlset{Namespace: "http://www.sitemaps.org/schemas/sitemap/0.9", Url: urls}

	output, err := xml.MarshalIndent(urlset, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	return xml_header + string(output)

}

func main() {
	url := flag.String("url", "https://www.calhoun.io/", "specify the start url to crawl")
	depth := flag.Uint("depth", 2, "specify the depth to crawl")
	flag.Parse()

	fmt.Println(build_sitemap_xml(extractSiteMap(*url, *depth)))

}
