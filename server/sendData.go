package server

import (
	"SM9_Client/mySM9"
	"SM9_Client/server/pkg"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type SendData map[string]string

// Marshal 将数据序列化
func (sendData *SendData) Marshal() []byte {
	res, _ := json.Marshal(sendData)
	return res
}

// MarshalWithLen 将数据序列化并返回长度
func (sendData *SendData) MarshalWithLen() ([]byte, string) {
	res, _ := json.Marshal(sendData)
	return res, strconv.FormatInt(int64(len(res)), 10)
}

// SendData 数据发送到服务端
func (sendData *SendData) SendData(url string) (map[string]string, error) {
	data := sendData.Marshal()
	response, err := http.Post(url, "application/json", bytes.NewReader(data))
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}
	val, _ := io.ReadAll(response.Body)
	var s map[string]string
	err = json.Unmarshal(val, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// DataEncSm9 对数据进行SM9加密
func (sendata *SendData) DataEncSm9() ([]byte, error) {
	res := sendata.Marshal()
	// 	加密
	encData, err := mySM9.Encrypt(pkg.Pkg.PKGPK, res)
	if err != nil {
		return nil, err
	}
	return encData, nil
}
