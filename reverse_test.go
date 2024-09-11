package reverse

import (
	"errors"
	"strconv"
	"strings"
	"testing"
)

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

func TestSimpleURLGeneration(t *testing.T) {
	clearRoutes()

	Add("first", "/first")

	url, err := Get("first")
	if err != nil {
		t.Errorf("an error occured: %s", err.Error())

		return
	}

	if url != "/first" {
		t.Errorf("got %s; want /first", url)
	}
}

func TestNotExistURLGeneration(t *testing.T) {
	clearRoutes()

	Add("first", "/first")

	_, err := Get("second")
	if err == nil {
		t.Error("an error was expected")
	} else {
		if !errors.Is(err, RouteNotFound) {
			t.Error("another error was expected")
		}
	}
}

func TestURLGenerationWithParams(t *testing.T) {
	clearRoutes()

	cases := []struct {
		name    string
		pattern string
		pairs   []string
		want    string
	}{
		{
			name:    "article",
			pattern: "/article/{slug}",
			pairs:   []string{"slug", "test"},
			want:    "/article/test",
		},
		{
			name:    "avatar",
			pattern: "/avatar/{hash:[0-9A-Z]+}",
			pairs:   []string{"hash", "123ABC", "quality", "95"},
			want:    "/avatar/123ABC?quality=95",
		},
		{
			name:    "articles",
			pattern: "/articles/{month}-{day}-{year}",
			pairs:   []string{"day", "09", "year", "2024", "month", "05"},
			want:    "/articles/05-09-2024",
		},
		{
			name:    "books",
			pattern: "/books/{rid:^[0-9]{5,6}}",
			pairs:   []string{"rid", "20200109-this-is-so-cool"},
			want:    "/books/20200109-this-is-so-cool",
		},
		{
			name:    "test-query",
			pattern: "/query",
			pairs:   []string{"a[1]", "1", "a[2]", "A B+C", "z", "123", "foo", "bar"},
			want:    "/query?a%5B1%5D=1&a%5B2%5D=A+B%2BC&foo=bar&z=123",
		},
	}

	for _, item := range cases {
		Add(item.name, item.pattern)
	}

	for _, item := range cases {
		url, err := Get(item.name, item.pairs...)
		if err != nil {
			t.Errorf("an error occured: %s", err.Error())
		}

		if url != item.want {
			t.Errorf("got %s; want %s", url, item.want)
		}
	}
}

func TestInvalidNumberParams(t *testing.T) {
	clearRoutes()

	Add("first", "/first")

	_, err := Get("first", "A", "B", "C")
	if err == nil {
		t.Error("an error was expected")
	} else {
		if !strings.Contains(err.Error(), "the number of parameters must be even") {
			t.Error("another error was expected")
		}
	}
}

func TestMismatchParams(t *testing.T) {
	clearRoutes()

	Add("photos", "/photos/{year}/{month}")

	_, err := Get("photos", "year", "2019")
	if err == nil {
		t.Error("an error was expected")
	} else {
		if !errors.Is(err, MismatchParams) {
			t.Error("another error was expected")
		}
	}
}
