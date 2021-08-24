package mock

type FieldTag struct {
	Example string
	Type    Type
	Len     string
	Regexp  string
}

type Type int8

const (
	Phone Type = iota
	Mail
	DateTime
)

func (t Type) String() string {
	switch t {
	case Phone:
		return "1235678910"
	case Mail:
		return "123@qq.com"
	case DateTime:
		return "2006-01-02 15:04:05"
	}
	return ""
}
