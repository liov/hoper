package dbi

type Delete struct {
	DeleteAt,Time string
}

func DeleteAt(t string) *Delete {
	return &Delete{"delete_at",t}
}