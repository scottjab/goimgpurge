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
func TestPurge(t *testing.T) {
	for _, image := range badImgs {
		_, err := imgpurge.Purge(image)
		assert.Error(t, err, "Not imgur.")
	}

	// Test regular image
	url, _ := imgpurge.Purge("http://imgur.com/SyuKFYj")
	assert.Equal(t, url.Host, "i.imgur.com")
	fmt.Println(url.Path)
	assert.Equal(t, url.Path, "/SyuKFYj.jpg")

	// Test fall though jpg
	url, _ = imgpurge.Purge("https://i.imgur.com/ycArzxR.jpg")
	assert.Equal(t, url.Host, "i.imgur.com")
	fmt.Println(url.Path)
	assert.Equal(t, url.Path, "/ycArzxR.jpg")

	// Test gifv
	url, _ = imgpurge.Purge("http://i.imgur.com/EhO081n.gifv")
	assert.Equal(t, url.Host, "i.imgur.com")
	fmt.Println(url.Path)
	assert.Equal(t, url.Path, "/EhO081n.gifv")

	// Test gif to gifv
	url, _ = imgpurge.Purge("http://i.imgur.com/Ci07j.gif")
	assert.Equal(t, url.Host, "i.imgur.com")
	fmt.Println(url.Path)
	assert.Equal(t, url.Path, "/Ci07j.gifv")

}

func TestPurgeNotImgur(t *testing.T) {
	url, err := imgpurge.Purge("http://validurl.com/somekindofgif.jpg")
	if err != nil {
		if err == imgpurge.ErrNotImgur {
			assert.Equal(t, url.String(), "http://validurl.com/somekindofgif.jpg")
		}
	}

}
