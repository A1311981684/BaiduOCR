package en_decrypt

import (
	"github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
)

func TestEncryptText(t *testing.T) {
	convey.Convey("", t, func(c convey.C) {
		data := "abcd1234"
		encrypted, err := EncryptText(data)
		c.So(err, convey.ShouldBeNil)
		log.Println(encrypted)
	})
}

func TestDecryptText(t *testing.T) {
	convey.Convey("", t, func(c convey.C) {
		data := "abcd1234"
		encrypted, err := EncryptText(data)
		c.So(err, convey.ShouldBeNil)
		log.Println("encrypted to:", encrypted)

		res := DecryptText(encrypted)
		c.So(res, convey.ShouldEqual, data)
	})
}