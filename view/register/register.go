package register

import (
	"SM9_Client/common"
	"SM9_Client/common/createData"
	"SM9_Client/mySM9"
	"SM9_Client/server"
	"SM9_Client/server/pkg"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

func getParam(text *widget.Entry) {
	// 产生随机数R1
	common.R1 = createData.CreateRandem()
	// 产生时间戳T1
	common.T1 = createData.CreateTime()

	err := pkg.Pkg.GetPKGID()
	if err != nil {
		text.SetText(err.Error())
		return
	}
	err = pkg.Pkg.GetPKGPK()
	if err != nil {
		text.SetText(err.Error())
		return
	}
	text.SetText("已从服务端获取注册参数")
	text.SetPlaceHolder("已从服务端获取注册参数")
}

func SubmitRegister(w fyne.Window, userName string) {
	// KDF生成CK和Ik
	common.CK, common.IK = mySM9.Kdf(append([]byte(common.R1), []byte(common.R2)...))

	// 用CK通过加密(UID||T2||用户资料)
	common.T2 = createData.CreateTime()
	uidT2Data := server.SendData{
		"UID":      common.UID,
		"T2":       common.T2,
		"UserData": userName,
	}
	uidT2DataMar, Len_uidT2DataMar := uidT2Data.MarshalWithLen()
	CipherUidT2AndUserData, err := mySM9.EncryptSm4(common.CK, uidT2DataMar)
	if err != nil {
		fmt.Println("CipherUidT2AndUserData fail")
		return
	}

	// 哈希HMACIK(ECK(UID||T2||用户资料))
	HmacCipherUidT2AndUserData := mySM9.HmacData(CipherUidT2AndUserData, common.IK)

	// 发送给服务端
	data := server.SendData{
		"UID":      common.UID,
		"Cipher":   hex.EncodeToString(CipherUidT2AndUserData),
		"HmacData": hex.EncodeToString(HmacCipherUidT2AndUserData),
		"DecLen":   Len_uidT2DataMar,
	}
	// 	发送给服务端
	res, err := data.SendData("http://127.0.0.1:8080/RecClientData")
	if err != nil {
		fmt.Println("http err", err)
		return
	}
	fmt.Println("res = ", res)
	dialog.ShowInformation("Info", "注册成功", w)
}

func SubmitAuthentication(w fyne.Window, userName string) {

	// 加密SID
	encSID, err := mySM9.Encrypt(pkg.Pkg.PKGPK, []byte(common.SID))
	if err != nil {
		dialog.NewError(err, w)
		return
	}

	// 	对encSID加签
	signEncSID, err := mySM9.Sign(encSID)
	if err != nil {
		dialog.NewError(err, w)
		return
	}

	// 加密UID、PKGID、T1、R1
	data := server.SendData{
		"UID":   common.UID,
		"PKGID": pkg.Pkg.PKGID,
		"T1":    common.T1,
		"R1":    common.R1,
	}
	// 	加密
	encUidPkgIdT1R1, err := data.DataEncSm9()
	if err != nil {
		return
	}

	// 将加密的内容和签名的内容发送到服务端
	sendData := server.SendData{
		"CipherSID":              hex.EncodeToString(encSID),
		"SignCipherSID":          hex.EncodeToString(signEncSID),
		"CipherUidAndPkgIdAndTR": hex.EncodeToString(encUidPkgIdT1R1),
	}
	response, err := sendData.SendData("http://127.0.0.1:8080/RecEncData")
	if err != nil {
		fmt.Println("response")
		return
	}
	// 解密
	encData, err := hex.DecodeString(response["data"])
	if err != nil {
		fmt.Println("hex err", err)
		return
	}
	decDataLen, _ := strconv.Atoi(response["dataLen"])
	decData, err := mySM9.DecryptSm4([]byte(common.R1), encData)
	recData := map[string]string{}
	err = json.Unmarshal(decData[:decDataLen], &recData)
	if err != nil {
		fmt.Println("Unmarshal err", err)
		return
	}

	// 赋值
	if common.UID != recData["UID"] {
		fmt.Println("uid err")
		return
	}
	common.R2 = recData["R2"]
	common.Token = recData["TOKEN"]

	fmt.Println("认证通过")

	SubmitRegister(w, userName)
}

func Register(w fyne.Window) fyne.CanvasObject {

	name := widget.NewEntry()
	name.SetPlaceHolder("username")

	email := widget.NewEntry()
	email.SetPlaceHolder("email(test@example.com)")
	// email.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "请保持正确的输入")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("password")

	paramData := widget.NewEntry()
	if pkg.Pkg.PKGID != "" {
		paramData.SetPlaceHolder("已从服务端获取注册参数")
	} else {
		paramData.SetPlaceHolder("未向服务端请求注册参数")
	}
	param := container.NewGridWithColumns(2, paramData, widget.NewButtonWithIcon("点击获取", theme.DownloadIcon(), func() { getParam(paramData) }))

	sex := widget.NewRadioGroup([]string{"男", "女"}, func(string) {})
	sex.Horizontal = true

	largeText := widget.NewMultiLineEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "", Widget: widget.NewLabel("")},
			{Text: "请求参数", Widget: param, HintText: "Parameters requested from the server before registration"},
			{Text: "用户名", Widget: name, HintText: "Your full name"},
			{Text: "性别", Widget: sex},
			{Text: "邮箱", Widget: email, HintText: "A valid email address"},
			{Text: "用户密码", Widget: password, HintText: "Your full password"},
			{Text: "备注", Widget: largeText},
		},
		OnCancel: func() {
			name.SetText("")
			email.SetText("")
			sex.SetSelected("")
			password.SetText("")
			largeText.SetText("")
		},
		OnSubmit: func() {
			SubmitAuthentication(w, name.Text)
			name.SetText("")
			email.SetText("")
			sex.SetSelected("")
			password.SetText("")
			largeText.SetText("")
		},
	}

	return form
}
