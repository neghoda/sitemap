package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/neghoda/sitemap/link"
	"github.com/neghoda/sitemap/marshal"
)

func main() {
	targetURL := flag.String("url", "https://www.calhoun.io", "Target URL for sitemap")
	location := flag.String("l", "sitemap.xml", "Output file location")
	flag.Parse()
	if !link.IsValidURL(*targetURL) {
		log.Fatal("Providet URL is not valid")
	}
	*targetURL = link.StripEndingSlash(*targetURL)
	sitemap := make(map[string]bool)
	sitemap[*targetURL] = false
	for {
		sitemapLen := len(sitemap)
		log.Printf("Parsing %v for links. %v unique paths", *targetURL, sitemapLen)
		sitemap, err := processSitemap(sitemap)
		if err != nil {
			log.Fatal(err)
		}
		if sitemapLen == len(sitemap) {
			break
		}
	}
	outputXML, err := os.Create(*location)
	if err != nil {
		log.Fatal(err)
	}
	defer outputXML.Close()
	err = marshal.SitemapToXML(sitemap, outputXML)
	if err != nil {
		log.Fatal(err)
	}
}

func processSitemap(sitemap map[string]bool) (map[string]bool, error) {
	if len(sitemap) == 0 {
		return sitemap, nil
	}
	for k := range sitemap {
		if sitemap[k] != true {
			urls, err := getUrlsForPage(k)
			if err != nil {
				return sitemap, err
			}
			sitemap = writeToSitemap(sitemap, urls)
			sitemap[k] = true
		}
	}
	return sitemap, nil
}

func writeToSitemap(sitemap map[string]bool, urls []string) map[string]bool {
	for _, v := range urls {
		if _, ok := sitemap[v]; !ok {
			sitemap[v] = false
		}
	}

	return sitemap
}

func getUrlsForPage(pageURL string) ([]string, error) {
	u, err := url.Parse(pageURL)
	if err != nil {
		return nil, err
	}
	host := u.Hostname()
	// If part of the target site is on http and part on https
	// this will cause duplicate links with different scheme
	scheme := u.Scheme
	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	links, err := link.ExtractLinks(resp.Body)
	if err != nil {
		return nil, err
	}
	links = link.LinksByHost(links, host)
	links = link.PathsToUrls(links, host, scheme)
	for i, v := range links {
		links[i] = link.StripEndingSlash(v)
	}
	return links, nil
}
