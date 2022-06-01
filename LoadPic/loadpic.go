package Loadpic

import (
	"image"
	_ "image/png"
	"os"

	"github.com/pkg/errors"

	"github.com/faiface/pixel"
)

func LoadPic(picpath string) (pic pixel.Picture, err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "加载图片出错")
		}
	}()

	//获取图片
	picfile, err := os.Open(picpath)
	if err != nil {
		return nil, err
	}
	defer picfile.Close()
	picpng, _, err := image.Decode(picfile)
	if err != nil {
		return nil, err
	}
	pic = pixel.PictureDataFromImage(picpng)

	return pic, nil
}
