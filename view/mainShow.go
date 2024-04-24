package view

import (
	"SM9_Client/common/createData"
	"SM9_Client/view/viewData"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"time"
)

const preferenceCurrentTreeData = "currentTreeData"

func unsupportedTreeData(t viewData.TreeData) bool {
	return !t.SupportWeb && fyne.CurrentDevice().IsBrowser()
}

func MainShow(w fyne.Window, a fyne.App) {
	topWindow := w
	content := container.NewMax()
	title := widget.NewLabel("Component name")
	setTreeData := func(t viewData.TreeData) {
		if fyne.CurrentDevice().IsMobile() {
			child := a.NewWindow(t.Title)
			topWindow = child
			child.SetContent(t.View(topWindow))
			child.Show()
			child.SetOnClosed(func() {
				topWindow = w
			})
			return
		}

		title.SetText(t.Info)

		content.Objects = []fyne.CanvasObject{t.View(w)}
		// content.Refresh()
	}
	tutorial := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator()), nil, nil, nil, content)
	if fyne.CurrentDevice().IsMobile() {
		w.SetContent(makeNav(setTreeData, false))
	} else {
		split := container.NewHSplit(makeNav(setTreeData, true), tutorial)
		split.Offset = 0.2
		w.SetContent(split)
	}

}

func makeNav(setTreeData func(tutorial viewData.TreeData), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()
	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return viewData.TreeDataIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := viewData.TreeDataIndex[uid]
			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := viewData.TreeDataS[uid]
			if !ok {
				fyne.LogError("Missing treeData panel:"+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
			if unsupportedTreeData(t) {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{Italic: true}
			} else {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{}
			}
		},
		OnSelected: func(uid string) {
			if t, ok := viewData.TreeDataS[uid]; ok {
				if unsupportedTreeData(t) {
					return
				}
				a.Preferences().SetString(preferenceCurrentTreeData, uid)
				setTreeData(t)
			}
		},
	}
	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentTreeData, "home")
		tree.Select(currentPref)
	}

	// 显示时间
	text := widget.NewLabelWithStyle("------------", fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Italic: true})
	timeData := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Italic: true})
	go func() {
		for range time.Tick(time.Second) {
			createData.UpDate(timeData)
		}
	}()
	return container.NewBorder(nil, container.NewVBox(text, timeData), nil, nil, tree)
}
