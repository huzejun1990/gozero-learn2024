// @Author huzejun 2024/11/5 17:51:00
package biz

const OK = 200

var (
	DBError         = NewError(10000, "数据库错误")
	AlreadyRegister = NewError(10100, "用户已注册")
)
