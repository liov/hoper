package model

func (u *User) TableName() string {
	if u.Id < 1_000_000 {
		return "user"
	}
	return "user_" + string(byte(u.Id/1_000_000 + 49))
}