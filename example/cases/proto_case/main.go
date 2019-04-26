package main

import (
	"fmt"
	"github.com/fwhezfwhez/SuperChecker"
)

//
//func main() {
//	var req = messages.ProMiniGameShopPropListRequest{
//		MessageId:  proto.Int32(660),
//		GameId:     proto.Int32(9),
//		PlatformId: proto.Int32(1),
//		PageType:   messages.MiniGamePageType_MINIGAME_SHOP_PAGE_EXCHANGE.Enum(),
//	}
//	fmt.Println(req.PageType)
//	sp := superChecker.GetChecker()
//	ok, msg, e := sp.ValidateByTagKeyAndMapValue(req, "json", map[string]string{
//		"page_type": "range,[1:7]",
//	})
//	fmt.Println(ok, msg, e)
//}
func main() {
	type User struct{
		Username string `json:"username"`
		Age int `json:"age"`
	}
	var req = User{
		Username: "superchecker",
		Age:-1,
	}
	sp := superChecker.GetChecker()
	ok, msg, e := sp.ValidateByTagKeyAndMapValue(req, "json", map[string]string{
		"username": "regex,^[\u4E00-\u9FA5a-zA-Z0-9_.]{0,40}$",
		"age": "int,0:200",
	})
	fmt.Println(ok, msg, e)
}
