package sys_model

import "time"

// MenuListRequest 菜单列表请求参数
type MenuListRequest struct {
	MenuName string `form:"menuName"` // 字典名称
	Status   string `form:"status"`   // 状态（0正常 1停用）
}

type SysMenu struct {
	MenuId     string                 `gorm:"primary_key;auto_increment;column:menu_id" json:"menuId" `           // 菜单ID
	MenuName   string                 `gorm:"column:menu_name;not null" json:"menuName" form:"menuName"`          // 菜单名称
	ParentId   string                 `gorm:"column:parent_id;default:0" json:"parentId"`                         // 父菜单ID
	OrderNum   int                    `gorm:"column:order_num;default:0" json:"orderNum"`                         // 显示顺序
	Path       string                 `gorm:"column:path;default:''" json:"path"`                                 // 路由地址
	Component  string                 `gorm:"column:component" json:"component"`                                  // 组件路径
	Query      string                 `gorm:"column:query" json:"query"`                                          // 路由参数
	IsFrame    string                 `gorm:"column:is_frame;default:1" json:"isFrame"`                           // 是否为外链（0是 1否）
	IsCache    string                 `gorm:"column:is_cache;default:0" json:"isCache"`                           // 是否缓存（0缓存 1不缓存）
	MenuType   string                 `gorm:"column:menu_type;type:char(1);default:''" json:"menuType"`           // 菜单类型（M目录 C菜单 F按钮）
	Visible    string                 `gorm:"column:visible;type:char(1);default:'0'" json:"visible"`             // 菜单状态（0显示 1隐藏）
	Status     string                 `gorm:"column:status;type:char(1);default:'0'" json:"status" form:"status"` // 菜单状态（0正常 1停用）
	Perms      string                 `gorm:"column:perms" json:"perms"`                                          // 权限标识
	Icon       string                 `gorm:"column:icon;default:'#'" json:"icon"`                                // 菜单图标
	CreateBy   string                 `gorm:"column:create_by;default:''" json:"createBy"`                        // 创建者
	CreateTime *time.Time             `gorm:"column:create_time" json:"createTime"`                               // 创建时间
	UpdateBy   string                 `gorm:"column:update_by;default:''" json:"updateBy"`                        // 更新者
	UpdateTime *time.Time             `gorm:"column:update_time" json:"updateTime"`                               // 更新时间
	Remark     string                 `gorm:"column:remark;default:''" json:"remark"`                             // 备注
	Params     map[string]interface{} `gorm:"-"`
}

// TableName 设置表名
func (m *SysMenu) TableName() string {
	return "sys_menu"
}

type SelectSysMenu struct {
	Name       string      `json:"name"`       //路由名字
	Path       string      `json:"path"`       //路由地址
	Hidden     bool        `json:"hidden"`     //否隐藏路由
	MenuType   string      `json:"menuType"`   // 菜单类型（M目录 C菜单 F按钮）
	Redirect   string      `json:"redirect"`   //重定向地址
	Component  string      `json:"component"`  //组件地址
	Query      string      `json:"query"`      //路由参数
	AlwaysShow bool        `json:"alwaysShow"` //当你一个路由下面的 children 声明的路由大于1个时，自动会变成嵌套的模式--如组件页面
	Meta       MetaVo      `json:"meta"`
	Children   []*RouterVo `json:"children"`
}

type RouterVo struct {
	Name       string      `json:"name"`       //路由名字
	Path       string      `json:"path"`       //路由地址
	Hidden     bool        `json:"hidden"`     //否隐藏路由
	MenuType   string      `json:"menuType"`   // 菜单类型（M目录 C菜单 F按钮）
	Redirect   string      `json:"redirect"`   //重定向地址
	Component  string      `json:"component"`  //组件地址
	Query      string      `json:"query"`      //路由参数
	AlwaysShow bool        `json:"alwaysShow"` //当你一个路由下面的 children 声明的路由大于1个时，自动会变成嵌套的模式--如组件页面
	Meta       MetaVo      `json:"meta"`
	Children   []*RouterVo `json:"children"`
}

type MetaVo struct {
	Title   string `json:"title"`   //设置该路由在侧边栏和面包屑中展示的名字
	Icon    string `json:"icon"`    //设置该路由的图标，对应路径src/assets/icons/svg
	NoCache bool   `json:"noCache"` //设置为true，则不会被 <keep-alive>缓存
	Link    string `json:"link"`    //重定向地址
}

// MenuTree 定义菜单树结构
type MenuTree struct {
	*SysMenu
	Children []*MenuTree `json:"children"`
}
type TreeSelect struct {
	Id       string        `json:"id"`       //节点ID
	Label    string        `json:"label"`    //节点名称
	Disabled bool          `json:"disabled"` //boolean
	Children []*TreeSelect `json:"children"`
}

//// SelectMenuTree 定义菜单树结构
//type SelectMenuTree struct {
//	*TreeSelect
//
//}
