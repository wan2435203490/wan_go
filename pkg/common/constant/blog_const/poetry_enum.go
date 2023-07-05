package blog_const

type PoetryEnum struct {
	Code int8
	Msg  string
}

var (
	// STATUS_ENABLE
	STATUS_ENABLE  = &PoetryEnum{1, "启用"}
	STATUS_DISABLE = &PoetryEnum{0, "禁用"}

	// PUBLIC
	PUBLIC  = &PoetryEnum{1, "所有人可见"}
	PRIVATE = &PoetryEnum{0, "仅自己可见"}

	// USER_GENDER_BOY
	USER_GENDER_BOY  = &PoetryEnum{1, "男"}
	USER_GENDER_GIRL = &PoetryEnum{2, "女"}
	USER_GENDER_NONE = &PoetryEnum{0, "保密"}

	// SORT_TYPE_BAR
	SORT_TYPE_BAR    = &PoetryEnum{0, "导航栏分类"}
	SORT_TYPE_NORMAL = &PoetryEnum{1, "普通分类"}

	// UserRoleAdmin UserRole  roleId
	UserRoleAdmin    = &PoetryEnum{1, "站长"}
	UserRoleDev      = &PoetryEnum{2, "管理员"}
	UserRoleOperator = &PoetryEnum{3, "用户"}
)
