package modelconst

type Action uint32

const (
	Signup Action = iota
	Active
	RestPassword
	EditPassword
	CreateResume
	EditResume
	DELETEResume
)

var action = map[Action]string{
	Signup:       "注册",
	Active:       "激活",
	RestPassword: "重置密码",
	EditPassword: "修改密码",
	CreateResume: "新建简历",
	EditResume:   "修改简历",
	DELETEResume: "删除简历",
}

//这样不可随意更改枚举值
var actionArray = []struct {
	Action
	string
}{
	{EditPassword, "修改密码"},
	{CreateResume, "新建简历"},
	{EditResume, "修改简历"},
	{DELETEResume, "删除简历"},
}

func (a Action) String() string {
	return action[a]
}
