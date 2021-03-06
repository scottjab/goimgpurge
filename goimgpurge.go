package imgpurge

import (
	"errors"
	"net/url"
	"path"
	"strings"
)

var ErrNotImgur = errors.New("NotImgur")

func IsImgur(url *url.URL) bool {
	return strings.Contains(url.Host, "imgur.com")
}

func CleanMobileLink(url *url.URL) *url.URL {
	if url.Host == "m.imgur.com" {
		if IsAlbum(url) {
			url.Host = "imgur.com"
		} else {
			url.Host = "i.imgur.com"
		}
	}
	return url
}

func IsAlbum(url *url.URL) bool {
	return strings.Contains(url.Path, "/a/") || strings.Contains(url.Path, "/gallery")
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
		return nil, err
	}

	if !IsImgur(parsedUrl) {
		return parsedUrl, ErrNotImgur
	}

	if IsAlbum(parsedUrl) {
		return parsedUrl, nil
	}

	return GetImage(parsedUrl), nil
}
