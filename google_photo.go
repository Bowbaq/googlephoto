package googlephoto

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// An Album is the top level object in the Google Photo (Picasa) API. An Album is
// essentially a list of Photos with some accompagnying metadata
type Album struct {
	// The Google Photo ID of this album
	ID string `xml:"http://schemas.google.com/photos/2007 id"`

	Name string `xml:"title"`

	// This number is sometimes higher than the number of Photo objects
	// retrievable with ListPhotos
	NumPhotos int `xml:"http://schemas.google.com/photos/2007 numphotos"`
}

// FeedURL return the URL of the atom feed for this Album
func (a Album) FeedURL() string {
	return "https://picasaweb.google.com/data/feed/api/user/default/albumid/" + a.ID
}

// A Photo has all the metadata that Google Photo (Picasa) stores about a user's photo
type Photo struct {
	ID     string `xml:"http://schemas.google.com/photos/2007 id"`
	ExifID string `xml:"http://schemas.google.com/photos/exif/2007 tags>imageUniqueID"`

	URL     string `xml:"http://www.w3.org/2005/Atom id"`
	Content struct {
		URL string `xml:"src,attr"`
	} `xml:"http://www.w3.org/2005/Atom content"`

	Name string `xml:"title"`

	Timestamp int `xml:"http://schemas.google.com/photos/2007 timestamp"`
	Size      int `xml:"http://schemas.google.com/photos/2007 size"`

	Published time.Time `xml:"http://www.w3.org/2005/Atom published"`
	Updated   time.Time `xml:"http://www.w3.org/2005/Atom updated"`
}

// Client represents a Google Photo API client
type Client struct {
	c *http.Client
}

// NewClient returns a new Client. A properly authentified oauth2 client
// must be provided (see golang.org/x/oauth2/google)
func NewClient(httpClient *http.Client) *Client {
	return &Client{httpClient}
}

// ListAlbums returns a list of the authenticated user's albums
func (c *Client) ListAlbums() ([]*Album, error) {
	data, err := c.getFeed("https://picasaweb.google.com/data/feed/api/user/default")
	if err != nil {
		return nil, err
	}

	var albumFeed struct {
		Albums []*Album `xml:"entry"`
	}
	err = xml.Unmarshal(data, &albumFeed)

	return albumFeed.Albums, err
}

// ListPhotos returns a list of Photo given an Album
func (c *Client) ListPhotos(album *Album) ([]*Photo, error) {
	var photos []*Photo

	start, end := 1, 1000
	previousLen := -1
	// Sometimes, the number of retrievable photos is less than the NumPhotos property on the album
	for previousLen != len(photos) {
		previousLen = len(photos)

		page, err := c.listPhotos(album, start, end)
		if err != nil {
			return nil, err
		}

		photos = append(photos, page...)

		start, end = end, end+1000
	}

	return photos, nil
}

func (c *Client) listPhotos(album *Album, start, end int) ([]*Photo, error) {
	url := album.FeedURL() + fmt.Sprintf("?start-index=%d&max-results=%d", start, end)
	log.Println(url)

	data, err := c.getFeed(url)
	if err != nil {
		return nil, err
	}

	var photoFeed struct {
		Photos []*Photo `xml:"entry"`
	}

	err = xml.Unmarshal(data, &photoFeed)

	return photoFeed.Photos, err
}

func (c *Client) getFeed(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("GData-Version", "2")

	res, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}
