package utils

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func Quit(w fyne.Window) func() {
	return func() {
		dialog.ShowConfirm("Quit", "确认退出？", func(b bool) {
			if b {
				w.Close()
			}
		}, w)
	}
}
