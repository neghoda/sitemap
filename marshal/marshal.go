package marshal

import (
	"encoding/xml"
	"io"
)

// Sitemap represents individual <loc> tag of sitemap protocol
type Sitemap struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
}

// Urlset resresents top <urlset> tag of sitemap protocol
type Urlset struct {
	XMLName    xml.Name `xml:"urlset"`
	XMLns      string   `xml:"xmlns,attr"`
	SitemapXML []Sitemap
}

// SitemapToXML writes sitemap XML to given location
// Doesn't write anything if sitemap parameter is empty
func SitemapToXML(sitemap map[string]bool, w io.Writer) error {
	if len(sitemap) == 0 {
		return nil
	}
	s := make([]Sitemap, 0, len(sitemap))
	for k := range sitemap {
		s = append(s, Sitemap{Loc: k})
	}
	xmlns := "http://www.sitemaps.org/schemas/sitemap/0.9"
	sitemapXML, err := xml.MarshalIndent(Urlset{XMLns: xmlns, SitemapXML: s}, "  ", "  ")
	if err != nil {
		return err
	}
	w.Write([]byte(xml.Header))
	w.Write(sitemapXML)
	return nil
}
