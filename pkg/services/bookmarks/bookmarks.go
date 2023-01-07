package bookmarks

import (
	"errors"
	"io"
	"mime/multipart"
	"net/url"

	"golang.org/x/net/html"
)

const (
	BookmarksFileKey     string = "bookmarks_file"
	BookmarksFileMaxSize int64  = 204800
)

// Bookmark represents a web bookmark.
type Bookmark struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	APIKey   string `json:"api_key" bson:"api_key"`
	Path     string `json:"path" bson:"path"`
	Name     string `json:"name" bson:"name"`
	URL      string `json:"url" bson:"url"`
	IsFolder bool   `json:"is_folder" bson:"is_folder"`
}

type HTMLBookmarkParser struct {
	tokenizer *html.Tokenizer
	APIKey    string
	bookmarks []Bookmark
}

func NewHTMLBookmarkParser(file multipart.File, APIKey string) *HTMLBookmarkParser {
	tokenizer := html.NewTokenizer(file)
	return &HTMLBookmarkParser{
		tokenizer: tokenizer,
		APIKey:    APIKey,
		bookmarks: []Bookmark{},
	}
}

func (h *HTMLBookmarkParser) parseBookmarkFileHTML() ([]Bookmark, error) {
	for {
		tokenType := h.tokenizer.Next()
		if tokenType == html.ErrorToken {
			err := h.tokenizer.Err()
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if tokenType == html.DoctypeToken {
			token := h.tokenizer.Token()
			if token.Data != "NETSCAPE-Bookmark-file-1" {
				return nil, errors.New("bookmark file incorrect format")
			}
		}
		if tokenType == html.StartTagToken {
			token := h.tokenizer.Token()
			if token.Data == "dt" {
				err := h.parseFolder("")
				if err != nil {
					return nil, errors.New("failed to parse bookmarks")
				}
			}
		}
	}
	return h.bookmarks, nil
}

func (h *HTMLBookmarkParser) parseFolder(path string) error {
	for {
		tokenType := h.tokenizer.Next()
		token := h.tokenizer.Token()
		data := token.Data
		attr := token.Attr
		if tokenType == html.ErrorToken {
			err := h.tokenizer.Err()
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
				f, err := h.createFolder(path)
				if err != nil {
					return err
				}
				h.bookmarks = append(h.bookmarks, f)
				newPath := path
				if len(newPath) == 0 {
					newPath += ","
				}
				newPath += f.Name + ","
				h.parseFolder(newPath)
			case "a":
				URL := findURL(attr)
				if len(URL) == 0 {
					break
				}
				b, err := h.createBookmark(path, URL)
				if err != nil {
					return err
				}
				h.bookmarks = append(h.bookmarks, b)
			}
		}
	}
	return nil
}

func (h *HTMLBookmarkParser) createFolder(path string) (Bookmark, error) {
	b := Bookmark{
		APIKey:   h.APIKey,
		Path:     path,
		URL:      "",
		Name:     "",
		IsFolder: true,
	}
	tokenType := h.tokenizer.Next()
	if tokenType != html.TextToken {
		return Bookmark{}, errors.New("bookmark folder does not have name text")
	}
	b.Name = html.UnescapeString(h.tokenizer.Token().Data)
	return b, nil
}

func (h *HTMLBookmarkParser) createBookmark(path string, URL string) (Bookmark, error) {
	b := Bookmark{
		APIKey:   h.APIKey,
		Path:     path,
		URL:      URL,
		Name:     "",
		IsFolder: false,
	}
	tokenType := h.tokenizer.Next()
	if tokenType != html.TextToken {
		return Bookmark{}, errors.New("bookmark does not have description text")
	}
	b.Name = html.UnescapeString(h.tokenizer.Token().Data)
	return b, nil
}

func findURL(attr []html.Attribute) string {
	for _, a := range attr {
		if a.Key == "href" {
			href, err := url.Parse(a.Val)
			if err != nil || href.Host == "" || href.Scheme == "" {
				break
			}
			return a.Val
		}
	}
	return ""
}
