package bookmarks

import (
	"errors"
	"io"
	"log"

	"golang.org/x/net/html"
)

// Bookmark represents a web bookmark.
type Bookmark struct {
	ID     string `json:"id" bson:"_id,omitempty"`
	APIKey string `json:"APIKey" bson:"APIKey"`
	Path   string `json:"path" bson:"path"`
	Name   string `json:"name" bson:"name"`
	URL    string `json:"url" bson:"url"`
}

func parseBookmarkFileHTML(APIKey string, tokenizer *html.Tokenizer) ([]Bookmark, error) {
	bookmarks := []Bookmark{}
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if tokenType == html.DoctypeToken {
			token := tokenizer.Token()
			if token.Data != "NETSCAPE-Bookmark-file-1" {
				// TODO: Change error
				return nil, errors.New("Bookmark file incorrect format")
			}
		}
		if tokenType == html.StartTagToken {
			token := tokenizer.Token()
			if token.Data == "dt" {
				err := parseFolder(&bookmarks, "", APIKey, tokenizer)
				if err != nil {
					// TODO: Change error
					return nil, errors.New("Failed to parse bookmarks")
				}
			}
		}
	}
	log.Println("parseHTMLBookmarks: ", bookmarks)
	return bookmarks, nil
}

func parseFolder(bookmarks *[]Bookmark, path, APIKey string, tokenizer *html.Tokenizer) error {
	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()
		data := token.Data
		attr := token.Attr
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}
			return err
		}
		if tokenType == html.EndTagToken && data == "dl" {
			break
		}
		if tokenType == html.StartTagToken {
			switch data {
			case "h3":
				nextTokenType := tokenizer.Next()
				if nextTokenType == html.TextToken {
					newPath := path
					if len(newPath) == 0 {
						newPath += ","
					}
					newPath += tokenizer.Token().Data + ","
					parseFolder(bookmarks, newPath, APIKey, tokenizer)
				}
			case "a":
				URL := findURL(attr)
				b, err := createBookmark(path, URL, APIKey, tokenizer)
				if err != nil {
					return err
				}
				*bookmarks = append(*bookmarks, b)
			}
		}
	}
	return nil
}

func createBookmark(path, URL, APIKey string, tokenizer *html.Tokenizer) (Bookmark, error) {
	b := Bookmark{
		APIKey: APIKey,
		Path:   path,
		URL:    URL,
	}
	tokenType := tokenizer.Next()
	if tokenType != html.TextToken {
		return Bookmark{}, errors.New("OMG")
	}
	b.Name = html.UnescapeString(tokenizer.Token().Data)
	return b, nil
}

func findURL(attr []html.Attribute) string {
	for _, a := range attr {
		if a.Key == "href" {
			return a.Val
		}
	}
	return ""
}
