package home

import (
	"SM9_Client/utils"
	_ "embed"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"net/url"
)

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}
func Home(_ fyne.Window) fyne.CanvasObject {

	// 顶部组件
	img := widget.NewIcon(theme.HomeIcon())
	str1 := canvas.NewText("SM9密码系统", nil)
	str1.TextSize = 35
	top := container.NewCenter(container.NewHBox(img, str1))

	// 底部组件
	text1 := widget.NewLabel("-")
	text2 := widget.NewLabel("-")
	url1 := widget.NewHyperlink("百度百科", parseURL("https://baike.baidu.com/item/SM9/2353509?fr=aladdin"))
	url2 := widget.NewHyperlink("文档支持", parseURL("http://c.gb688.cn/bzgk/gb/showGb?type=online&hcno=02A8E87248BD500747D2CD484C034EB0"))
	url3 := widget.NewHyperlink("国家密码管理局", parseURL("http://www.nca.gov.cn/"))
	bottom := container.NewCenter(container.NewHBox(url1, text1, url2, text2, url3))

	// 主要内容组件
	home, err := utils.Img.GetImage("home")
	if err != nil {
		fmt.Println("图片 fail", err)
	}
	img2 := canvas.NewImageFromResource(fyne.NewStaticResource("home", home))
	img2.FillMode = canvas.ImageFillContain

	return container.NewBorder(top, bottom, nil, nil, img2)
}
