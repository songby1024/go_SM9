package topView

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TopView(a fyne.App, window fyne.Window) *fyne.MainMenu {
	it1 := fyne.NewMenuItem("新的窗口", func() {
		w2 := a.NewWindow("new")
		text := widget.NewLabel("new")
		w2.SetContent(text)
		w2.Show()
	})

	dark := fyne.NewMenuItem("Dark(黑色)", func() {
		a.Settings().SetTheme(theme.DarkTheme())
	})
	light := fyne.NewMenuItem("Light(白色)", func() {
		a.Settings().SetTheme(theme.LightTheme())
	})
	it2 := fyne.NewMenuItem("选择窗口背景", func() {})
	it2.ChildMenu = fyne.NewMenu("", dark, light)

	openSettings := func() {
		w := a.NewWindow("设置")
		w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(480, 480))
		w.Show()
	}
	settingsItem := fyne.NewMenuItem("面板设置", openSettings)

	m1 := fyne.NewMenu("设置", it1, it2, settingsItem)

	it4 := fyne.NewMenuItem("支持", func() {
		dialog.ShowError(errors.New("本项目暂不提供支持"), window)
	})
	it5 := fyne.NewMenuItem("联系作者", func() {
		dialog.ShowInformation("Info", "作者地址：桂林电子科技大学图书馆415\n"+
			"作者名称：桂电图书馆415无名氏       ", window)
	})
	m2 := fyne.NewMenu("帮助", it4, it5)
	return fyne.NewMainMenu(m1, m2)
}
