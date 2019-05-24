package token

import (
	"github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
)

func TestGetTokenStr(t *testing.T) {
	convey.Convey("", t, func(c convey.C) {
		s, n, e := getToken()
		c.So(e, convey.ShouldBeNil)
		log.Println(s, n)
	})
}
