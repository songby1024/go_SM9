package main

import (
	"SM9_Client/utils"
	"SM9_Client/view"
	"SM9_Client/view/topView"
	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"os"
)

//go:embed common/image/Sm9.png
var sm9Icon []byte

//go:embed common/image/home.jpg
var home []byte

// 传送图片
func init() {
	data := map[string][]byte{
		"sm9Icon": sm9Icon,
		"home":    home,
	}
	utils.Img.ReceiveImage(data)
}

func main() {
	// 引入中文字体文件
	os.Setenv("FYNE_FONT", "common/font/方正仿宋简体.ttf")
	// 设置字体大小
	os.Setenv("FYNE_SCALE", "1.3")
	// 产生一个新的窗口
	a := app.New()
	// 设置默认背景为黑色
	a.Settings().SetTheme(theme.DarkTheme())
	// 设置窗口名称
	w := a.NewWindow("SM9 Client")
	// 设置窗口大小
	w.Resize(fyne.NewSize(750, 500))
	// 设置退出窗口要调用的函数
	w.SetCloseIntercept(utils.Quit(w))
	// 设置窗口在中间位置显示
	w.CenterOnScreen()
	// 设置窗体图标
	w.SetIcon(fyne.NewStaticResource("SM9", sm9Icon))
	// 顶部视图
	w.SetMainMenu(topView.TopView(a, w))
	// 主体视图
	view.MainShow(w, a)
	// 展示与运行
	w.ShowAndRun()
}
