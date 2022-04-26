package code

import error2 "github.com/quanxiang-cloud/cabin/error"

func init() {
	error2.CodeTable = CodeTable
}

const (
	// ErrExistRoleNameState
	ErrExistRoleNameState = 90014000001
	// ErrExistPermitState ErrExistRoleNameState
	ErrExistPermitState = 90014000002
	//ErrItemConvert ErrItemConvert
	ErrItemConvert = 90054000001
	// ErrNotPermit ErrNotPermit
	ErrNotPermit = 90074000002
	// ErrParameter ErrParameter
	ErrParameter = 90074000003
)

// CodeTable 码表
var CodeTable = map[int64]string{
	ErrExistRoleNameState: "角色名称不能重复，请重新输入！",
	ErrExistPermitState:   "权限设置已设置 ,不能重复设置",
	ErrItemConvert:        "参数Items错误",
	ErrNotPermit:          "没有权限 ，权限为空",
	ErrParameter:          "类型转换错误",
}
