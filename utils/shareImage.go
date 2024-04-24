package utils

import (
	_ "embed"
	"errors"
)

var Img Images

type Images struct {
	Image map[string][]byte
}

func (img *Images) ReceiveImage(Images map[string][]byte) {
	img.Image = Images
}

func (img *Images) GetImage(s string) ([]byte, error) {
	data, ok := img.Image[s]
	if !ok { // 不存在这个图片
		return nil, errors.New("该图片不存在")
	}
	return data, nil
}

func (img *Images) GetAll() map[string][]byte {
	return img.Image
}
