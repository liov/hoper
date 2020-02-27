package model

func (m *User) TableName() string {
	if m.Id < 1_000_000 {
		return "user"
	}
	return "user_" + string(byte(m.Id/1_000_000+49))
}

func (m *UserMainInfo) TableName() string {
	if m.Id < 1_000_000 {
		return "user"
	}
	return "user_" + string(byte(m.Id/1_000_000+49))
}

func (m *UserBaseInfo) TableName() string {
	if m.Id < 1_000_000 {
		return "user"
	}
	return "user_" + string(byte(m.Id/1_000_000+49))
}
