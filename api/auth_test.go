package api

import (
	"fmt"
	"testing"
)

func TestAuth_Sign(t *testing.T) {
	a := Auth{
		CP:     "cp",
		Seckey: "testSeckey",
	}
	ctime := int64(1630983600)
	sign := a.Sign(ctime)
	fmt.Println(a.Sign(ctime))
	if sign != "ef4132113104454034eae76f12c0ff89" {
		panic("sign error")
	}
}
