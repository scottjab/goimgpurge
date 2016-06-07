package imgpurge_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/scottjab/goimgpurge"
	"github.com/stretchr/testify/assert"
)

var goodAlbums = []string{"https://imgur.com/gallery/arP2Otg",
	"https://imgur.com/a/arP2Otg"}
var goodImgs = []string{
	"https://i.imgur.com/ycArzxR.jpg",
	"http://imgur.com/SyuKFYj"}
var badImgs = []string{
	"https://google.com/thing.jpg",
	"https://reddit.com/r/otherthing",
	"http://imgurl.com/imgur.com/.gif",
}

func TestIsAlbum(t *testing.T) {
	for _, image := range goodAlbums {
		url, _ := url.Parse(image)
		assert.True(t, imgpurge.IsAlbum(url))
	}
	for _, image := range goodImgs {
		url, _ := url.Parse(image)
		assert.False(t, imgpurge.IsAlbum(url))
	}
}

func TestIsImgur(t *testing.T) {
	// Test for false
	for _, image := range badImgs {
		url, _ := url.Parse(image)
		assert.False(t, imgpurge.IsImgur(url))
	}
	// Test for true
	for _, image := range goodImgs {
		url, _ := url.Parse(image)
		assert.True(t, imgpurge.IsImgur(url))
	}
}

func TestCleanMobile(t *testing.T) {
	mobileLink, _ := url.Parse("https://m.imgur.com/SyuKFYj.jpg")
	mobileLink = imgpurge.CleanMobileLink(mobileLink)
	assert.Equal(t, "i.imgur.com", mobileLink.Host)
}

func TestCleanMobileAlbum(t *testing.T) {
	mobileAlbum, _ := url.Parse("https://m.imgur.com/a/arP2Otg")
	mobileAlbum = imgpurge.CleanMobileLink(mobileAlbum)
	assert.Equal(t, "imgur.com", mobileAlbum.Host)
}

func TestCleanMobileFallThough(t *testing.T) {
	mobileLink, _ := url.Parse("https://m.catbot.io/image.jpg")
	mobileLink = imgpurge.CleanMobileLink(mobileLink)
	assert.Equal(t, "m.catbot.io", mobileLink.Host)
}

func TestPurge(t *testing.T) {
	for _, image := range badImgs {
		_, err := imgpurge.Purge(image)
		assert.Error(t, err, "Not imgur.")
	}

	// Test regular image
	url, _ := imgpurge.Purge("http://imgur.com/SyuKFYj")
	assert.Equal(t, "i.imgur.com", url.Host)
	fmt.Println(url.Path)
	assert.Equal(t, "/SyuKFYj.jpg", url.Path)

	// Test fall though jpg
	url, _ = imgpurge.Purge("https://i.imgur.com/ycArzxR.jpg")
	assert.Equal(t, "i.imgur.com", url.Host)
	fmt.Println(url.Path)
	assert.Equal(t, "/ycArzxR.jpg", url.Path)

	// Test gifv
	url, _ = imgpurge.Purge("http://i.imgur.com/EhO081n.gifv")
	assert.Equal(t, "i.imgur.com", url.Host)
	fmt.Println(url.Path)
	assert.Equal(t, "/EhO081n.gifv", url.Path)

	// Test gif to gifv
	url, _ = imgpurge.Purge("http://i.imgur.com/Ci07j.gif")
	assert.Equal(t, "i.imgur.com", url.Host)
	fmt.Println(url.Path)
	assert.Equal(t, "/Ci07j.gifv", url.Path)

}

func TestPurgeNotImgur(t *testing.T) {
	url, err := imgpurge.Purge("http://validurl.com/somekindofgif.jpg")
	if err != nil {
		if err == imgpurge.ErrNotImgur {
			assert.Equal(t, url.String(), "http://validurl.com/somekindofgif.jpg")
		}
	}

}
