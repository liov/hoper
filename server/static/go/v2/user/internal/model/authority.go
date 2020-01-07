package modelconst

//权限
type Authority uint64

const (
	DeleteUser Authority = 1 << iota
	EditUser
)

var authority = map[Authority]string{
	DeleteUser: "删除用户",
}

func (a Authority) String() string {
	return authority[a]
}
