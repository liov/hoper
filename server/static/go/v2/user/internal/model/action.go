package modelconst

type Action uint32

const (
	EditPassWord Action = iota
	CreateResume
	EditResume
)

var action = map[Action]string{
	EditPassWord: "修改密码",
	CreateResume: "新建简历",
	EditResume:   "修改简历",
}

func (a Action) String() string {
	return action[a]
}
