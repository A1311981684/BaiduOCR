package ocr

import (
	"github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
)

func TestOpenImageFileToBase64(t *testing.T) {
	convey.Convey("", t, func(c convey.C) {
		s , err := OpenImageFileToBase64(`C:\Users\MetalFish\Desktop\a.jpg`)
		c.So(err,convey.ShouldBeNil)
		c.So(s, convey.ShouldNotEqual, "")
	})
}

func TestOCR(t *testing.T) {
	convey.Convey("", t, func(c convey.C) {
		s, err := OCR(`C:\Users\MetalFish\Desktop\a.jpg`)
		if err != nil {
			panic(err)
		}
		for _, v := range s{
			for k1, v1 := range v {
				log.Println(k1, v1)
			}
		}
	})
}