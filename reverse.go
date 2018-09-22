package reverse

import (
	"errors"
	"fmt"
	"strings"
)

var Urls *urlStore

func init() {
	Urls = &urlStore{store: make(map[string]url)}
}

func Clear() {
	for k := range Urls.store {
		delete(Urls.store, k)
	}
}

// Add url to store
func Add(urlName, urlAddr string, params ...string) string {
	return Urls.MustAdd(urlName, urlAddr, params...)
}

// AddGr url with concat group, but returns just urlAddr
func AddGr(urlName, group, urlAddr string, params ...string) string {
	return Urls.MustAddGr(urlName, group, urlAddr, params...)
}

// Rev url by name
func Rev(urlName string, params ...string) string {
	return Urls.MustReverse(urlName, params...)
}

// Get raw url by name
func Get(urlName string) string {
	return Urls.Get(urlName)
}

// GetAllURLs saved all urls
func GetAllURLs() map[string]string {
	out := map[string]string{}
	for key, value := range Urls.store {
		out[key] = value.url
	}
	return out
}

// GetAllParams all params
func GetAllParams() map[string][]string {
	out := map[string][]string{}
	for key, value := range Urls.store {
		out[key] = value.params
	}
	return out
}

type url struct {
	url    string
	params []string
}

type urlStore struct {
	store map[string]url
}

// Add a Url to the Store
func (us *urlStore) Add(urlName, urlAddr string, params ...string) (string, error) {
	return us.AddGr(urlName, "", urlAddr, params...)
}

// MustAdd a Url and panics if error
func (us urlStore) MustAdd(urlName, urlAddr string, params ...string) string {
	addr, err := us.Add(urlName, urlAddr, params...)
	if err != nil {
		panic(err)
	}
	return addr
}

// AddGr with group prefix
func (us *urlStore) AddGr(urlName, group, urlAddr string, params ...string) (string, error) {
	if _, ok := us.store[urlName]; ok {
		return "", errors.New("Url already exists. Try to use .Get() method.")
	}

	tmpUrl := url{group + urlAddr, params}
	us.store[urlName] = tmpUrl
	return urlAddr, nil
}

// MustAddGr a Url with group prefix
func (us urlStore) MustAddGr(urlName, group, urlAddr string, params ...string) string {
	addr, err := us.AddGr(urlName, group, urlAddr, params...)
	if err != nil {
		panic(err)
	}
	return addr
}

// Get raw url string
func (us urlStore) Get(urlName string) string {
	return us.store[urlName].url
}

// Reverse url
func (us urlStore) Reverse(urlName string, params ...string) (string, error) {
	if len(params) != len(us.store[urlName].params) {
		return "", errors.New("Bad Url Reverse: mismatch params for URL: " + urlName)
	}
	res := us.store[urlName].url
	for i, val := range params {
		res = strings.Replace(res, us.store[urlName].params[i], val, 1)
	}
	return res, nil
}

// MustReverse url and panics if error
func (us urlStore) MustReverse(urlName string, params ...string) string {
	res, err := us.Reverse(urlName, params...)
	if err != nil {
		panic(err)
	}
	return res
}

// Rev url, but returns empty string in case of error
func (us urlStore) Rev(urlName string, params ...string) string {
	return us.MustReverse(urlName, params...)
}

func (us urlStore) String() string {
	return fmt.Sprint(us.store)
}

// For testing
func (us urlStore) getParam(urlName string, num int) string {
	return us.store[urlName].params[num]
}
