package dbi

var (
	tableName = [...]string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
)

func TableName(name string, id uint64) string {
	if id < 2000_00000 {
		return name
	}
	if id < 2_0000_00000 {
		return name + tableName[id/2000_00000-1]
	}
	return name + string(byte(id/2000_00000+49))
}
