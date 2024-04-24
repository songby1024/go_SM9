package createData

import (
	"fyne.io/fyne/v2/widget"
	"strconv"
	"time"
)

func CreateTime() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func UpDate(oldTime *widget.Label) {
	oldTime.SetText(time.Now().Format("Time : 15:04:05"))
}
