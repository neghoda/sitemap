package link

import (
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// ExtractLinks searches valid html for links tags
func ExtractLinks(r io.Reader) ([]string, error) {
	document, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	links := parseForLinks(document)
	return links, nil
}

// LinksByHost execudes all links to other hosts
func LinksByHost(links []string, host string) []string {
	if len(links) == 0 {
		return links
	}
	LinksByHost := make([]string, 0, len(links))
	for _, v := range links {
		if isOnTheSameHost(v, host) {
			LinksByHost = append(LinksByHost, v)
		}
	}

	return LinksByHost
}

// IsValidURL checks if URL can be parsed and containts host and scheme
// Not reliable. Returns true for something like httpps://test.com
func IsValidURL(link string) bool {
	u, err := url.Parse(link)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}

func isOnTheSameHost(link string, host string) bool {
	if strings.HasPrefix(link, "/") {
		return true
	}
	if !IsValidURL(link) {
		return false
	}
	u, _ := url.Parse(link)
	if u.Host != host {
		return false
	}

	return true
}

// PathsToUrls transforms all relative paths to URLS
func PathsToUrls(links []string, host string, scheme string) []string {
	if len(links) == 0 {
		return links
	}
	for i, v := range links {
		if strings.HasPrefix(v, "/") {
			links[i] = scheme + "://" + host + v
		}
	}
	return links
}

// StripEndingSlash returns link without clousing slash
func StripEndingSlash(link string) string {
	if strings.HasSuffix(link, "/") {
		return strings.TrimSuffix(link, "/")
	}
	return link
}

func parseForLinks(n *html.Node) []string {
	var links []string
	if n.Type == html.ElementNode && n.Data == "a" {
		links = append(links, extractAttributeValue(n, "href"))
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, parseForLinks(c)...)
	}
	return links
}

func extractAttributeValue(n *html.Node, attr string) string {
	for _, v := range n.Attr {
		if v.Key == attr {
			return v.Val
		}
	}
	return ""
}
