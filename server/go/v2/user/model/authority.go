package modelconst

//权限
type Authority uint64

const (
	DeleteUser Authority = 1 << iota
	EditUser
)

func (a Authority) String() string {
	switch a {
	case DeleteUser:
		return "删除用户"
	}
	return ""
}
