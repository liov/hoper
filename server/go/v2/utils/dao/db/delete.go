package dbi

type Deleted struct {
	DeleteAt,Time string
}

func DeletedAt(t string) *Deleted {
	return &Deleted{"deleted_at",t}
}