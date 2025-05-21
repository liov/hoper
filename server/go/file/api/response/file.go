package response

type ListDir struct {
	Items []*Entry `json:"items"`
}

type Entry struct {
	Name    string `json:"name"`
	Type    int    `json:"type"`
	ModTime int64  `json:"modTime"`
	Url     string `json:"url"`
	Thumb   string `json:"thumb"`
	Size    int64  `json:"size"`
}

type File struct {
	Id  string `json:"id"`
	URL string `json:"url"`
}
