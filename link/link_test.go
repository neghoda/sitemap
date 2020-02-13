package link

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestExtractLinks(t *testing.T) {
	testData := "<html><body><h1>Hello!</h1><a href=\"/other-page\">A link to another page</a></body></html>"
	expected := []string{"/other-page"}
	if result, _ := ExtractLinks(strings.NewReader(testData)); !reflect.DeepEqual(result, expected) {
		t.Error(testFailMessage(expected, result))
	}

	testData = "<html><body><h1>Hello!</h1></body></html>"
	expected = []string{}
	if result, _ := ExtractLinks(strings.NewReader(testData)); len(result) != 0 {
		t.Error(testFailMessage(expected, result))
	}

	testData = `<div>
    				<a href="https://www.twitter.com/joncalhoun">
      					Check me out on twitter
      					<i class="fa fa-twitter" aria-hidden="true"></i>
    				</a>
    				<a href="https://github.com/gophercises">
      					Gophercises is on <strong>Github</strong>!
    				</a>
  				</div>`
	expected = []string{"https://www.twitter.com/joncalhoun", "https://github.com/gophercises"}
	if result, _ := ExtractLinks(strings.NewReader(testData)); !reflect.DeepEqual(result, expected) {
		t.Error(testFailMessage(expected, result))
	}
	testData = "<a href=\"/dog-cat\">dog cat <!-- commented text SHOULD NOT be included! --></a>"
	expected = []string{"/dog-cat"}
	if result, _ := ExtractLinks(strings.NewReader(testData)); !reflect.DeepEqual(result, expected) {
		t.Error(testFailMessage(expected, result))
	}
}

func TestLinksByHost(t *testing.T) {
	l := []string{
		"https://google.com",
		"http://test.com",
		"https://test.com",
	}
	expected := []string{"http://test.com", "https://test.com"}

	if result := LinksByHost(l, "test.com"); !reflect.DeepEqual(result, expected) {
		t.Error(testFailMessage(expected, result))
	}
}

func TestIsValidURL(t *testing.T) {
	l := "https:sdadestcom"
	if result := IsValidURL(l); result {
		t.Error(testFailMessage(result, false))
	}
	l = "https://test.com"
	if result := IsValidURL(l); !result {
		t.Error(testFailMessage(result, true))
	}
}

func TestStripEndingSlash(t *testing.T) {
	l := "https://test.com/"
	expected := "https://test.com"
	if result := StripEndingSlash(l); result != expected {
		t.Error(testFailMessage(expected, result))
	}
	if result := StripEndingSlash(expected); result != expected {
		t.Error(testFailMessage(expected, result))
	}
}

func testFailMessage(expected interface{}, result interface{}) string {
	return fmt.Sprintf("Expected %v to be equal %v", expected, result)
}
