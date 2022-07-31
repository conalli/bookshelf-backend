package search

import (
	"flag"
	"strings"
)

func formatURL(url string) string {
	if strings.Contains(url, "http://") || strings.Contains(url, "https://") {
		return url
	}
	return "http://" + url
}

// LSFlag represents the possible flags for the ls command.
type LSFlag struct {
	fs *flag.FlagSet
	b  *bool
	c  *bool
	bf *string
}

// NewLSFlagset returns a new flag set for the ls command.
func NewLSFlagset() LSFlag {
	fs := flag.NewFlagSet("ls", flag.ContinueOnError)
	b := fs.Bool("b", false, "lists all bookmarks")
	c := fs.Bool("c", false, "lists all cmds")
	bf := fs.String("bf", "", "lists all bookmarks for given folder")
	ls := LSFlag{
		fs: fs,
		b:  b,
		c:  c,
		bf: bf,
	}
	return ls
}

// TouchFlag represents the possible flags for the touch command.
type TouchFlag struct {
	fs   *flag.FlagSet
	b    *bool
	c    *string
	url  *string
	path *string
	name *string
}

// NewTouchFlagset returns a new flag set for the touch command.
func NewTouchFlagset() TouchFlag {
	fs := flag.NewFlagSet("touch", flag.ContinueOnError)
	b := fs.Bool("b", false, "adds a bookmark")
	c := fs.String("c", "", "adds a cmd")
	url := fs.String("url", "", "url for new bookmark")
	path := fs.String("path", "", "folder path for new bookmark")
	name := fs.String("name", "", "name for new bookmark")
	ls := TouchFlag{
		fs:   fs,
		b:    b,
		c:    c,
		url:  url,
		path: path,
		name: name,
	}
	return ls
}
