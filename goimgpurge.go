package imgpurge

import (
	"errors"
	"net/url"
	"path"
	"strings"
)

func IsImgur(url *url.URL) bool {
	return strings.Contains(url.Host, "imgur.com")
}

func IsAlbum(url *url.URL) bool {
	found := strings.Contains("/a/", url.Path[:3])
	if !found && len(url.Path) > 9 {
		found = strings.Contains("/gallery/", url.Path[:9])
	}
	return found
}

func makeGifv(url *url.URL) *url.URL {
	if path.Ext(url.Path) == ".gif" {
		url.Path = url.Path + "v"
	}
	return url
}

func GetImage(url *url.URL) *url.URL {
	if strings.Contains(url.Host, "i.imgur.com") {
		return makeGifv(url)
	} else {
		if path.Ext(url.Path) == "" {
			url.Path = url.Path[:len(url.Path)] + ".jpg"
			url.Host = "i.imgur.com"

		}
	}
	return url
}

func Purge(rawUrl string) (*url.URL, error) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		panic(err)
	}

	if !IsImgur(parsedUrl) {
		return nil, errors.New("Not imgur.")
	}

	if IsAlbum(parsedUrl) {
		return parsedUrl, nil
	}

	return GetImage(parsedUrl), nil
}
