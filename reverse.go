package reverse

import (
	"errors"
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
func Add(routeName, pattern string) string {
	return routes.mustAdd(routeName, pattern)
}

// AddGr url with concat group, but returns just pattern
func AddGr(routeName, group, pattern string) string {
	return routes.mustAddGr(routeName, group, pattern)
}

// Get url by name
func Get(routeName string, pairs ...string) (string, error) {
	return routes.reverse(routeName, pairs...)
}

// MustGet url by name
func MustGet(routeName string, pairs ...string) string {
	return routes.mustReverse(routeName, pairs...)
}

// GetAllURLs saved all urls
func GetAllURLs() map[string]string {
	out := map[string]string{}
	for key, value := range routes.store {
		out[key] = value.pattern
	}

	return out
}

func (us urlStore) add(routeName, pattern string) (string, error) {
	if _, ok := us.store[routeName]; ok {
		return "", RouteAlreadyExist
	}

	us.store[routeName] = createRoute(pattern)

	return pattern, nil
}

func (us urlStore) mustAdd(routeName, pattern string) string {
	addr, err := us.add(routeName, pattern)
	if err != nil {
		panic(err)
	}

	return addr
}

func (us urlStore) addGr(routeName, groupName, pattern string) (string, error) {
	return us.add(routeName, pattern)
}

func (us urlStore) mustAddGr(routeName, groupName, pattern string) string {
	addr, err := us.addGr(routeName, groupName, pattern)
	if err != nil {
		panic(err)
	}

	return addr
}

func (us urlStore) reverse(routeName string, pairs ...string) (string, error) {
	var r route
	var ok bool
	if r, ok = us.store[routeName]; !ok {
		return "", RouteNotFound
	}

	/* if len(pairs) != len(us.store[routeName].params) {
		return "", errors.New("reverse: mismatch params for route: " + routeName)
	}

	res := us.store[routeName].pattern
	for i, val := range pairs {
		res = strings.Replace(res, us.store[routeName].params[i], val, 1)
	} */

	return r.url(pairs...)
}

func (us urlStore) mustReverse(routeName string, pairs ...string) string {
	res, err := us.reverse(routeName, pairs...)
	if err != nil {
		panic(err)
	}

	return res
}

// For testing
func (us urlStore) getParam(routeName string, num int) parameter {
	return us.store[routeName].params[num]
}
