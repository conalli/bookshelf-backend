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

type LSFlag struct {
	fs *flag.FlagSet
	b  *bool
	c  *bool
	bf *string
}

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

type TouchFlag struct {
	fs   *flag.FlagSet
	b    *bool
	c    *bool
	url  *string
	path *string
	name *string
}

func NewTouchFlagset() TouchFlag {
	fs := flag.NewFlagSet("touch", flag.ContinueOnError)
	b := fs.Bool("b", false, "adds a bookmark")
	c := fs.Bool("c", false, "adds a cmds")
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

// cmds:flags
// ls: -c | -b | -bf
// touch | new: -c -url | -b -url -path? -name?
// man | help
