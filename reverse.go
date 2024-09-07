package reverse

import (
	"errors"
	"strings"
)

type urlStore struct {
	store map[string]route
}

var (
	routes urlStore

	RouteAlreadyExist = errors.New("reverse: route already exists")
	RouteNotFound     = errors.New("reverse: route not found")
)

func init() {
	routes = urlStore{
		store: make(map[string]route),
	}
}

// clear routes store for testing
func clearRoutes() {
	for k := range routes.store {
		delete(routes.store, k)
	}
}

// Add url to store
func Add(urlName, urlAddr string) string {
	return routes.mustAdd(urlName, urlAddr)
}

// AddGr url with concat group, but returns just urlAddr
func AddGr(urlName, group, urlAddr string) string {
	return routes.mustAddGr(urlName, group, urlAddr)
}

// Get url by name
func Get(urlName string, pairs ...string) (string, error) {
	return routes.reverse(urlName, pairs...)
}

// MustGet url by name
func MustGet(urlName string, pairs ...string) string {
	return routes.mustReverse(urlName, pairs...)
}

// GetAllURLs saved all urls
func GetAllURLs() map[string]string {
	out := map[string]string{}
	for key, value := range routes.store {
		out[key] = value.pattern
	}

	return out
}

func (us urlStore) add(urlName, urlAddr string) (string, error) {
	return us.addGr(urlName, "", urlAddr)
}

func (us urlStore) mustAdd(urlName, urlAddr string) string {
	addr, err := us.add(urlName, urlAddr)
	if err != nil {
		panic(err)
	}

	return addr
}

func (us urlStore) addGr(urlName, group, urlAddr string) (string, error) {
	if _, ok := us.store[urlName]; ok {
		return "", RouteAlreadyExist
	}

	tmpUrl := route{
		pattern: group + urlAddr,
	}
	us.store[urlName] = tmpUrl

	return urlAddr, nil
}

func (us urlStore) mustAddGr(urlName, group, urlAddr string) string {
	addr, err := us.addGr(urlName, group, urlAddr)
	if err != nil {
		panic(err)
	}

	return addr
}

func (us urlStore) reverse(urlName string, pairs ...string) (string, error) {
	if _, ok := us.store[urlName]; !ok {
		return "", RouteNotFound
	}

	if len(pairs) != len(us.store[urlName].params) {
		return "", errors.New("reverse: mismatch params for route: " + urlName)
	}

	res := us.store[urlName].pattern
	for i, val := range pairs {
		res = strings.Replace(res, us.store[urlName].params[i], val, 1)
	}

	return res, nil
}

func (us urlStore) mustReverse(urlName string, pairs ...string) string {
	res, err := us.reverse(urlName, pairs...)
	if err != nil {
		panic(err)
	}

	return res
}

// For testing
func (us urlStore) getParam(urlName string, num int) string {
	return us.store[urlName].params[num]
}
