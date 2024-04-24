package viewData

import (
	"SM9_Client/view/getKey"
	"SM9_Client/view/home"
	"SM9_Client/view/register"
	"SM9_Client/view/register/userInfo"
	"fyne.io/fyne/v2"
)

type TreeData struct {
	Title, Info string
	View        func(w fyne.Window) fyne.CanvasObject
	SupportWeb  bool
}

var (
	TreeDataS = map[string]TreeData{
		"home":        {"Home", "主页", home.Home, true},
		"userInfo":    {"用户注册", "用户信息", userInfo.UserInfo, true},
		"register":    {"注册", "用户注册", register.Register, true},
		"get_UserKey": {"私钥获取", "用户私钥下发", getKey.GetKey, true},
	}
	TreeDataIndex = map[string][]string{
		"":         {"home", "userInfo", "get_UserKey"},
		"userInfo": {"register"},
	}
)
