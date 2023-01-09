package bookmarks

type Folder struct {
	ID        string     `json:"id,omitempty"`
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	Bookmarks []Bookmark `json:"bookmarks"`
	Folders   []Folder   `json:"folders"`
}

func organizeBookmarks(bookmarks []Bookmark, folderID, folderName, folderPath, path string) *Folder {
	length := len(bookmarks)
	if length == 0 {
		return &Folder{}
	}
	folder := &Folder{ID: folderID, Name: folderName, Path: folderPath}
	for _, b := range bookmarks {
		if b.Path != path {
			continue
		}
		if b.IsFolder {
			newPath := updatePath(path, b.Name)
			folder.Folders = append(folder.Folders, *organizeBookmarks(bookmarks, b.ID, b.Name, path, newPath))
		} else {
			folder.Bookmarks = append(folder.Bookmarks, b)
		}
	}
	return folder
}
