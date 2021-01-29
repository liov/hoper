package content


var MomentserviceServicedesc = &_MomentService_serviceDesc

/*----------------------------ORM-------------------------------*/
func momentTableName(id uint64) string {
	if id < 1_000_000 {
		return "moment"
	}
	return "moment_" + string(byte(id/1_000_000+49))
}
func (x *AddMomentReq) TableName() string {
	return "moment"
}