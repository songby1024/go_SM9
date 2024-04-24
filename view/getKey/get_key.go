package getKey

import (
	"SM9_Client/common"
	"SM9_Client/common/createData"
	"SM9_Client/mySM9"
	"SM9_Client/server"
	"SM9_Client/server/pkg"
	"crypto/hmac"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/emmansun/gmsm/sm9"
	"strconv"
	"time"
)

func textShow(text *widget.Entry, msg string) {
	text.SetText(text.Text + "\n\n" + msg)
	time.Sleep(time.Millisecond * 450)
}

// 处理数据
func handleData(CipherAndHmacData map[string]string, w fyne.Window, text *widget.Entry) bool {
	EncUidT4UidSkHmac, _ := hex.DecodeString(CipherAndHmacData["EncUidT4UidSkHmac"])
	HmacEncUidT4UidSkHmac, _ := hex.DecodeString(CipherAndHmacData["HmacEncUidT4UidSkHmac"])
	len_EncUidT4UidSkHmac, _ := strconv.Atoi(CipherAndHmacData["len_EncUidT4UidSkHmac"])
	len_uidSkData, _ := strconv.Atoi(CipherAndHmacData["len_uidSkData"])
	len_uidSK, _ := strconv.Atoi(CipherAndHmacData["len_uidSK"])

	// 	利用IK校验消息完整性
	HmacIk := mySM9.HmacData(EncUidT4UidSkHmac, common.IK)
	if !hmac.Equal(HmacIk, HmacEncUidT4UidSkHmac) {
		fmt.Println("IK校验消息完整性失败")
	}
	textShow(text, "7 利用IK校验消息完整性 success")

	// 	解密DECK (ECK (UID||T4|| ER3(UIDSK ) ||HMACTOKEN))
	textShow(text, "8 解密DECK (ECK (UID||T4|| ER3(UIDSK ) ||HMACTOKEN))")
	decData, err := mySM9.DecryptSm4(common.CK, EncUidT4UidSkHmac)
	if err != nil {
		fmt.Println("EncUidT4UidSkHmac fail", err)
	}

	// 对数据进行反序列化
	textShow(text, "9 对数据进行反序列化")
	recData := map[string]string{}
	err = json.Unmarshal(decData[:len_EncUidT4UidSkHmac], &recData)
	if err != nil {
		fmt.Println("Unmarshal err", err)
	}
	uid := recData["UID"]
	t4 := recData["T4"]
	uidSkCipher, _ := hex.DecodeString(recData["EncUidSk"])
	hmacToken, _ := hex.DecodeString(recData["hmacToken"])

	// 	判断UID
	if uid != common.UID {
		fmt.Println("UID不一样")
	}

	// 用token检验完整性
	myHmac := mySM9.HmacData(uidSkCipher, []byte(common.Token))
	if !hmac.Equal(myHmac, hmacToken) {
		fmt.Println("哈希验证失败")
	}
	textShow(text, "10 用token检验完整性")

	// 用t4判断消息是否重放
	timeData := createData.CreateTime()
	time, _ := strconv.Atoi(timeData)
	t4Val, _ := strconv.Atoi(t4)
	if t4Val > time {
		fmt.Println("t4 时间不对")
	}
	textShow(text, "11 用t4判断消息是否重放 success")

	// 解密UIDSK（R3 sm4加密）
	uidSK, err := mySM9.DecryptSm4([]byte(common.R3), uidSkCipher)
	if err != nil {
		fmt.Println("uidSK err", err)
	}
	textShow(text, "12 解密UIDSK success")

	// 本地保存UIDSK
	parKey := uidSK[:len_uidSkData][:len_uidSK]
	pubKey := uidSK[:len_uidSkData][len_uidSK:]
	common.UIDSK.UnmarshalASN1(parKey)
	publicKey := new(sm9.EncryptMasterPublicKey)
	publicKey.UnmarshalASN1(pubKey)
	common.UIDSK.EncryptMasterPublicKey = *publicKey

	fmt.Println("UIDSK = ", common.UIDSK)
	return true
}

func SendResponse(w fyne.Window, text *widget.Entry) {
	// 产生随机数R3与时间戳T3并加密
	// 1产生R3与T3
	common.R3 = createData.CreateRandem()
	common.T3 = createData.CreateTime()
	textShow(text, "1 产生随机数R3与时间戳T3")

	// 	2用PKGPK对R3进行sm9加密
	cipherR3, err := mySM9.Encrypt(pkg.Pkg.PKGPK, []byte(common.R3))
	if err != nil {
		fmt.Println("R3 加密失败", err)
	}
	textShow(text, "2 对随机数R3进行SM9加密")

	// 	3用Ck对UID、T3、cipherR3进行sm4加密
	data1 := server.SendData{
		"UID":      common.UID,
		"T3":       common.T3,
		"CipherR3": hex.EncodeToString(cipherR3),
	}
	textShow(text, "3 用Ck对UID、T3、cipherR3进行sm4加密")
	data2, Len_uidT3R3 := data1.MarshalWithLen()
	UidT3CipherR3, err := mySM9.EncryptSm4(common.CK, data2)

	// 	4对UidT3CipherR3进行哈希
	hmacUidT3CipherR3 := mySM9.HmacData(UidT3CipherR3, common.IK)
	textShow(text, "4 对UidT3CipherR3进行哈希")

	// 将数据发送到服务端
	textShow(text, "5 将数据发送到服务端………………")
	textShow(text, "loading……………………………………")
	send := server.SendData{
		"UID":               common.UID,
		"UidT3CipherR3":     hex.EncodeToString(UidT3CipherR3),
		"hmacUidT3CipherR3": hex.EncodeToString(hmacUidT3CipherR3),
		"len_UidT3CipherR3": Len_uidT3R3,
	}
	response, err := send.SendData("http://127.0.0.1:8080/sendUIDSK")
	if err != nil {
		fmt.Println("服务端数据读取失败", err)
	}
	textShow(text, "6 成功读取服务器的返回消息")
	// 	对返回的数据进行处理
	if handleData(response, w, text) {
		fmt.Println("用户私钥获取成功")
		textShow(text, "用户私钥获取成功")
		text.SetPlaceHolder(text.Text)
		dialog.ShowInformation("Info", "已成功获取私钥", w)
	}
}

func GetKey(w fyne.Window) fyne.CanvasObject {
	text := widget.NewMultiLineEntry()
	text.SetPlaceHolder("此处显示获取私钥信息")
	button := widget.NewButton("点击获取私钥", func() {
		SendResponse(w, text)
	})
	return container.NewVSplit(container.NewCenter(button), text)
}
