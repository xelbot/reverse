package reverse

import (
	"errors"
)

type urlStore struct {
	store  map[string]route
	groups map[string]string
}

var (
	routes urlStore

	RouteAlreadyExist = errors.New("reverse: route already exists")
	RouteNotFound     = errors.New("reverse: route not found")
	MismatchParams    = errors.New("reverse: mismatch params for route")

	GroupAlreadyExist = errors.New("reverse: group already exists")
	GroupNotFound     = errors.New("reverse: group not found")
)

func init() {
	routes = urlStore{
		store:  make(map[string]route),
		groups: make(map[string]string),
	}
}

// clear routes store for testing
func clearRoutes() {
	for k := range routes.store {
		delete(routes.store, k)
	}
	for k := range routes.groups {
		delete(routes.groups, k)
	}
}

// Add url to store
func Add(routeName, pattern string) string {
	return routes.mustAdd(routeName, pattern)
}

// AddGr url with group, but returns just pattern
func AddGr(routeName, groupName, pattern string) string {
	return routes.mustAddGr(routeName, groupName, pattern)
}

// Group add group prefix to store
func Group(groupName, pattern string) string {
	return routes.mustGroup(groupName, pattern)
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
	var prefix string
	var ok bool

	if prefix, ok = us.groups[groupName]; !ok {
		return "", GroupNotFound
	}

	_, err := us.add(routeName, prefix+pattern)
	if err != nil {
		return "", err
	}

	return pattern, nil
}

func (us urlStore) mustAddGr(routeName, groupName, pattern string) string {
	addr, err := us.addGr(routeName, groupName, pattern)
	if err != nil {
		panic(err)
	}

	return addr
}

func (us urlStore) group(groupName, pattern string) (string, error) {
	if _, ok := us.groups[groupName]; ok {
		return "", GroupAlreadyExist
	}

	us.groups[groupName] = pattern

	return pattern, nil
}

func (us urlStore) mustGroup(groupName, pattern string) string {
	prefix, err := us.group(groupName, pattern)
	if err != nil {
		panic(err)
	}

	return prefix
}

func (us urlStore) reverse(routeName string, pairs ...string) (string, error) {
	var r route
	var ok bool
	if r, ok = us.store[routeName]; !ok {
		return "", RouteNotFound
	}

	return r.url(pairs...)
}

func (us urlStore) mustReverse(routeName string, pairs ...string) string {
	res, err := us.reverse(routeName, pairs...)
	if err != nil {
		panic(err)
	}

	return res
}
