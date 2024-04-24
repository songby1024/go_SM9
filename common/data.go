package common

import (
	"github.com/emmansun/gmsm/sm9"
)

const (
	// SID 私有标识
	SID = "SID001"
	// UID 公开标识
	UID = "UID001"
)

var (
	T1 string
	// 	生成新的时间戳T2
	T2 string
	T3 string
	R1 string
	// 	从服务端接收的R2
	R2 string
	R3 string
	// Token 从服务端接受的Token
	Token string
	// 	保存KDF生成的CK、IK
	CK []byte
	IK []byte
	// UIDSK 用户最终私钥
	UIDSK sm9.EncryptPrivateKey
)
