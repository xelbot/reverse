package reverse

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
)

func TestOldReverse(t *testing.T) {
	showError := func(info string) {
		t.Error(fmt.Sprintf("Error: %s. urlStore: %s", info, routes))
	}

	if routes.mustAdd("firstUrl", "/first") != "/first" {
		showError("0")
	}

	if routes.mustAdd("helloUrl", "/hello/:p1:p2") != "/hello/:p1:p2" {
		showError("0-1")
	}

	/* routes.mustAdd("secondUrl", "/second/:param/:param2")

	routes.mustAdd("thirdUrl", "/comment/:p1")

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
	} */
}

func TestAdd(t *testing.T) {
	clearRoutes()

	cases := []struct {
		name    string
		pattern string
	}{
		{
			name:    "first",
			pattern: "/first",
		},
		{
			name:    "second",
			pattern: "/second/{id:[0-9]+}",
		},
	}

	for idx, item := range cases {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			result := Add(item.name, item.pattern)

			if result != item.pattern {
				t.Errorf("%s : got %s; want %s", item.name, result, item.pattern)
			}
		})
	}
}

func TestAddDuplicate(t *testing.T) {
	clearRoutes()

	Add("first", "/first")
	Add("second", "/second")

	_, err := routes.add("first", "/first-second")
	if err == nil {
		t.Error("an error was expected")
	} else {
		if !errors.Is(err, RouteAlreadyExist) {
			t.Error("another error was expected")
		}
	}
}
