package request

type ListDir struct {
	Dir string `uri:"dir" json:"dir"`
}

type Exists struct {
	Md5  string `uri:"md5" query:"md5" json:"md5"`
	Size string `uri:"size" query:"size" json:"size"`
}

type Path struct {
	Path string `uri:"path" json:"path"`
}
