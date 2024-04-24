package pkg

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var Pkg PkgData

type PkgData struct {
	PKGID  string `json:"pkgid"`
	STATUS int    `json:"status"`
	PKGPK  []byte `json:"pkgpk"`
	PKGSK  []byte `json:"pkgsk"`
}

func (pkg *PkgData) GetPKGID() error {
	// 	从PPS获取PKGID
	response, err := http.Get("http://127.0.0.1:8080/get-pkgid")
	if err != nil {
		return errors.New("PKGID fail")
	}
	val, _ := io.ReadAll(response.Body)
	defer response.Body.Close()
	err = json.Unmarshal(val, &pkg)
	if err != nil {
		return errors.New("PKGID Unmarshal fail")
	}
	return nil
}
func (pkg *PkgData) GetPKGPK() error {
	// 	从PPS获取PKGPk
	file, err := http.Get("http://127.0.0.1:8080/get-enc-public-key")
	if err != nil {
		return errors.New("PKGPK fail")
	}
	defer file.Body.Close()
	pkg.PKGPK, _ = io.ReadAll(file.Body)
	return nil
}
