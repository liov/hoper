package modelconst

import "github.com/liov/hoper/go/v2/utils/strings2"

//确定的，极少量的枚举值，没必要map
type Gender uint8

const (
	GenderUnfilled Gender = iota
	GenderMale
	GenderFemale
)

func (x Gender) String() string {
	switch x {
	case GenderMale:
		return "男"
	case GenderFemale:
		return "女"
	default:
		return "未填"
	}
}

func (x Gender) MarshalJSON() ([]byte, error) {
	return strings2.QuoteToBytes(x.String()), nil
}
