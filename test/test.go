package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type result_CreateEncryptPublicKey struct {
	Data       string `json:"data"`
	StatusCode int    `json:"status_code"`
}
type result_SM4Check struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

type result_ExchangeKey struct {
	Code int    `json:"code"`
	Data string `json:"data"`
	Key1 string `json:"key1"`
	Key2 string `json:"key2"`
}

func main() {
	a := app.New()
	w := a.NewWindow("Client Demo01")
	MainShow(w)
	w.Resize(fyne.Size{Width: 500, Height: 225})
	// w.CenterOnScreen()
	w.ShowAndRun()
}

func MainShow(w fyne.Window) {
	text := widget.NewMultiLineEntry()
	bt1 := widget.NewButton("CreateEncryptPublicKey", func() {
		go func() {
			resp, err := http.Get("http://127.0.0.1:8080/create-encrypt-public-key")
			if err != nil {
				fmt.Println("get failed, err:", err)
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("read from resp.Body failed,err:", err)
				return
			}
			var res result_CreateEncryptPublicKey
			json.Unmarshal(body, &res)
			code := fmt.Sprintf("%d", res.StatusCode)
			text.SetText("CreateEncryptPublicKey\nloading···\n")
			time.Sleep(time.Second)
			text.SetText("CreateEncryptPublicKey\nloading···\n" + "data : " + res.Data + "\n" + "status_code : " + code)
		}()
	})
	bt2 := widget.NewButton("SM4Check", func() {
		go func() {
			resp, err := http.Get("http://127.0.0.1:8080/sm4-check")
			if err != nil {
				fmt.Println("get failed, err:", err)
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("read from resp.Body failed,err:", err)
				return
			}
			var res result_SM4Check
			json.Unmarshal(body, &res)
			code := fmt.Sprintf("%d", res.Code)
			text.SetText("SM4Check\nloading···\n")
			time.Sleep(time.Second)
			text.SetText("SM4Check\nloading···\n" + "code : " + code + "\n" + "data : " + res.Data)
		}()
	})
	bt3 := widget.NewButton("ExchangeKey", func() {
		go func() {
			resp, err := http.Get("http://127.0.0.1:8080/exchange-key")
			if err != nil {
				fmt.Println("get failed, err:", err)
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("read from resp.Body failed,err:", err)
				return
			}
			var res result_ExchangeKey
			json.Unmarshal(body, &res)
			code := fmt.Sprintf("%d", res.Code)
			text.SetText("ExchangeKey\nloading···\n")
			time.Sleep(time.Second)
			text.SetText("ExchangeKey\nloading···\n" + "code : " + code + "\n" + "data : " + res.Data + "\n" + "key1 : " + res.Key1 + "\n" + "key2 : " + res.Key2)
		}()
	})
	v1 := container.NewHBox(bt1, bt2, bt3)
	content := container.NewBorder(v1, nil, nil, nil, text)
	w.SetContent(content)
}
