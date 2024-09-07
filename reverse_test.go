package reverse

import (
	"fmt"
	"testing"
)

func TestReverse(t *testing.T) {
	showError := func(info string) {
		t.Error(fmt.Sprintf("Error: %s. urlStore: %s", info, routes))
	}

	if routes.mustAdd("firstUrl", "/first") != "/first" {
		showError("0")
	}

	if routes.mustAdd("helloUrl", "/hello/:p1:p2", "1", "2") != "/hello/:p1:p2" {
		showError("0-1")
	}

	routes.mustAdd("secondUrl", "/second/:param/:param2", ":param", ":param2")

	// re := regexp.MustCompile("^/comment/(?P<id>\d+)$")
	routes.mustAdd("thirdUrl", "/comment/:p1", ":p1")

	if routes.getParam("helloUrl", 1) != "2" {
		showError("1")
	}

	if routes.getParam("secondUrl", 0) != ":param" {
		showError("3")
	}

	if routes.mustReverse("firstUrl") != "/first" {
		showError("4")
	}

	if routes.mustReverse("secondUrl", "123", "ABC") != "/second/123/ABC" {
		showError("5")
	}

	if routes.mustReverse("thirdUrl", "123") != "/comment/123" {
		t.Error(routes.reverse("thirdUrl", "123"))
		showError("6")
	}

	clearRoutes()
	if len(routes.store) != 0 {
		showError("7")
	}
}
