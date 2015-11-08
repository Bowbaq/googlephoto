Google Photo API client
=======================

[![GoDoc](https://godoc.org/github.com/Bowbaq/googlephoto?status.svg)](https://godoc.org/github.com/Bowbaq/googlephoto)

Example
-------
```go
package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"

  "github.com/Bowbaq/googlephoto"
  "golang.org/x/oauth2"
  "golang.org/x/oauth2/google"
)

func main() {
  client := googlephoto.NewClient(oauth2Client())

  photos := make(map[string]*googlephoto.Photo)

  albums, err := client.ListAlbums()
  check(err)

  for _, album := range albums {
    albumPhotos, err := client.ListPhotos(album)
    check(err)

    for _, photo := range albumPhotos {
      if _, seen := photos[photo.ID]; !seen {
        photos[photo.ID] = photo
      } else {
        photos[photo.ID].Albums = append(photos[photo.ID].Albums, album.ID)
      }
    }
  }
}

func oauth2Client() *http.Client {
  conf := &oauth2.Config{
    // Get Oauth2 credentials from https://console.developers.google.com
    ClientID:     "<client id>",
    ClientSecret: "<client secret>",
    RedirectURL:  "http://localhost/",
    Scopes: []string{
      "https://picasaweb.google.com/data/",
    },
    Endpoint: google.Endpoint,
  }

  var token oauth2.Token
  tokenData, err := ioutil.ReadFile("./token.json")
  if err == nil {
    err = json.Unmarshal(tokenData, &token)

    if err != nil {
      log.Println(err)
      token = refreshToken(conf)
    }
  } else {
    log.Println(err)
    token = refreshToken(conf)
  }

  return conf.Client(oauth2.NoContext, &token)
}

func refreshToken(conf *oauth2.Config) oauth2.Token {
  url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
  fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

  var authCode string
  fmt.Println("Your browser should have redirected to: http://localhost/?state=state&code=<code>")
  fmt.Print("Paste the code: ")
  fmt.Scanln(&authCode)

  token, err := conf.Exchange(oauth2.NoContext, authCode)
  check(err)

  data, err := json.Marshal(token)
  check(err)

  ioutil.WriteFile("./token.json", data, 0400)

  return *token
}

func check(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

```
