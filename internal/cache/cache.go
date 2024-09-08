package cache

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var cachePath string

func generateKey(r *http.Request) string {
	path := strings.Trim(r.URL.Path, "/")
	if r.URL.RawQuery != "" {
		path += "?" + r.URL.RawQuery
	}
	requestLine := fmt.Sprintf("%s /%s", r.Method, path)
	hash := sha1.Sum([]byte(requestLine))
	return hex.EncodeToString(hash[:])
}

func Store(r *http.Request, res *http.Response) error {
	key := generateKey(r)
	file, err := os.Create(cachePath + key)
	if err != nil {
		return err
	}
	defer file.Close()
	return res.Write(file)
}

func Load(r *http.Request) (*http.Response, error) {
	key := generateKey(r)
	file, err := os.Open(cachePath + key)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	return http.ReadResponse(reader, r)
}

func Has(r *http.Request) bool {
	key := generateKey(r)
	_, err := os.Stat(cachePath + key)
	return !os.IsNotExist(err)
}

func init() {
	var err error
	cachePath, err = os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	cachePath += "/caching-proxy/"
	os.Mkdir(cachePath, 0755)
}
