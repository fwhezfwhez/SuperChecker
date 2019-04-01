package main

import (
	"fmt"
	"github.com/fwhezfwhez/SuperChecker"
	"github.com/fwhezfwhez/SuperChecker/example/cases/proto_case/message"
	"github.com/golang/protobuf/proto"
)

func main() {
	var req = messages.ProMiniGameShopPropListRequest{
		MessageId:  proto.Int32(660),
		GameId:     proto.Int32(9),
		PlatformId: proto.Int32(1),
		PageType:   messages.MiniGamePageType_MINIGAME_SHOP_PAGE_EXCHANGE.Enum(),
	}
	fmt.Println(req.PageType)
	sp := superChecker.GetChecker()
	ok, msg, e := sp.ValidateByTagKeyAndMapValue(req, "json", map[string]string{
		"page_type": "range,[1:7]",
	})
	fmt.Println(ok, msg, e)
}
