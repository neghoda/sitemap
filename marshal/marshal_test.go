package marshal

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestSitemapToXML(t *testing.T) {
	links := map[string]bool{
		"https://test.com": true,
		"http://test.com":  true,
	}
	buff := new(bytes.Buffer)
	expected, error := ioutil.ReadFile("test.xml")
	if error != nil {
		t.Errorf("Error while reading test.xml - %v", error)
	}
	if error := SitemapToXML(links, buff); reflect.DeepEqual(buff, expected) || error != nil {
		t.Errorf("Expected %v to be equal %v", buff.String(), string(expected))
	}
}
