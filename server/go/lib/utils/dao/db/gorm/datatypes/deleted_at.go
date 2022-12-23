package datatypes

import (
	"encoding/json"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"time"
)

type DeletedAt time.Time

func (n DeletedAt) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(n))
}

func (n *DeletedAt) UnmarshalJSON(b []byte) error {
	return (*time.Time)(n).UnmarshalJSON(b)
}

func (n DeletedAt) MarshalText() ([]byte, error) {
	return (time.Time)(n).MarshalText()
}

func (n *DeletedAt) UnmarshalText(b []byte) error {
	return (*time.Time)(n).UnmarshalText(b)
}

func (n DeletedAt) String() string {
	return (time.Time)(n).String()
}

func (DeletedAt) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{clause.Where{
		Exprs: []clause.Expression{
			clause.Eq{
				Column: "deleted_at",
				Value:  "0001-01-01 08:05:43",
			},
		},
	}}
}

func (d DeletedAt) UpdateClauses(f *schema.Field) []clause.Interface {
	return d.QueryClauses(f)
}

func (d DeletedAt) DeleteClauses(f *schema.Field) []clause.Interface {
	return d.QueryClauses(f)
}
