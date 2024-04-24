package userInfo

import (
	"SM9_Client/common"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/emmansun/gmsm/sm9"
)

func UserInfo(_ fyne.Window) fyne.CanvasObject {

	// 设置标题
	title := canvas.NewText("当前用户信息", nil)
	title.TextSize = 20
	top := container.NewCenter(container.NewHBox(widget.NewIcon(theme.ComputerIcon()), title))

	// 展示用户的信息（UID，HID）
	uidText := canvas.NewText("UID :", nil)
	uid := widget.NewEntry()
	uid.SetPlaceHolder(common.UID)

	sidText := canvas.NewText("SID :", nil)
	sid := widget.NewEntry()
	sid.SetPlaceHolder(common.SID)

	uidSkText := canvas.NewText("UIDSK :", nil)
	uidSk := widget.NewEntry()
	data := new(sm9.EncryptPrivateKey)
	if common.UIDSK == *data {
		uidSk.SetPlaceHolder("用户还未获取私钥")
	} else {
		uidSk.SetPlaceHolder("用户已获取私钥")
		uidSk.SetText("用户已获取私钥")
	}

	// 介绍SM9
	boby1 := widget.NewLabel("\n\n\n\n\n注：SM9密码算法简介\n\nSM9是中华人民共和国政府采用的一种标识密码标准，由国家密码管理局于2016年3月28日发布,\n" +
		"相关标准为“GM/T 0044-2016 SM9标识密码算法。在商用密码体系中，SM9主要用于用户的身份\n" +
		"认证.资料显示SM9的加密强度等同于3072位密钥的RSA加密算法。")

	boby := container.NewVBox(widget.NewLabel(""), container.New(layout.NewFormLayout(), uidText, uid, sidText, sid, uidSkText, uidSk), container.NewCenter(boby1))
	return container.NewBorder(top, nil, nil, nil, boby)
}
