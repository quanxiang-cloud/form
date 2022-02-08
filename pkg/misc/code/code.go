package code

import error2 "github.com/quanxiang-cloud/cabin/error"

func init() {
	error2.CodeTable = CodeTable
}

const (
	// ErrExistGroupNameState 用户组名已经存在
	ErrExistGroupNameState = 90014000001

	// ErrGroupNotExit 权限组不存在
	ErrGroupNotExit = 90014000002

	// ErrRepeatMenuName 菜单名重复
	ErrRepeatMenuName = 90024000001

	// ErrDeleteMenu 删除菜单失败
	ErrDeleteMenu = 90024000002

	// InvalidCondition invalid input condition
	InvalidCondition = 90034000001

	// ErrExistDataSetNameState 数据集名已经存在
	ErrExistDataSetNameState = 90044000001

	// ErrNODataSetNameState 数据不存在
	ErrNODataSetNameState = 90044000002

	//ErrItemConvert ErrItemConvert
	ErrItemConvert = 90054000001

	//ErrValueConvert ErrValueConvert
	ErrValueConvert = 90054000002
	//ErrConvert ErrConvert
	ErrConvert = 90054000003
	//ErrRepeatTableID ErrRepeatTableID
	ErrRepeatTableID = 90064000001
	//ErrRepeatTableTitle ErrRepeatTableTitle
	ErrRepeatTableTitle = 90064000002

	// ErrCustomFileNotExist ErrCustomFileNotExist
	ErrCustomFileNotExist = 90074000001

	// ErrNotPermit ErrNotPermit
	ErrNotPermit = 90074000002
)

// CodeTable 码表
var CodeTable = map[int64]string{
	ErrExistGroupNameState: "角色名称不能重复，请重新输入！",
	ErrGroupNotExit:        "权限组不存在！",

	ErrRepeatMenuName:        "同一分组下的页面名称不能重复，请重新输入！",
	ErrDeleteMenu:            "删除分组或菜单失败！该分组存在页面时不能进行删除！",
	InvalidCondition:         "修改或删除失败！当前操作未授权或传参错误！",
	ErrExistDataSetNameState: "数据集名称不能重复，请重新输入！",
	ErrNODataSetNameState:    "数据不存在!",

	ErrItemConvert:  "参数Items错误",
	ErrValueConvert: "参数 value转换异常",
	ErrConvert:      "参数错误",

	ErrRepeatTableID:    "模型编码已经存在",
	ErrRepeatTableTitle: "模型名字已经存在",

	ErrCustomFileNotExist: "自定义页面文件不存在",

	ErrNotPermit: "没有权限 ，权限为空",
}
