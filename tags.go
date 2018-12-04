package superChecker

import (
	"github.com/shopspring/decimal"
	"time"
)

type T struct {
	UserName string `superChecker:"userName"`
	Password string `superChecker:"password"`
	Phone    string `superChecker:"mobilephone|telephone"`
	Text     string `superChecker:"length,chineseOnly"`
}
type T2 struct {
	CreatedAt         time.Time       `validate:"time.time"`
	UpdatedAt         string          `validate:"time.time,2006/01/02 15:04:05"`
	Count             int             `validate:"int,0:200"`
	MaxCount          int             `validate:"int,:200"`
	MinCount          int             `validate:"int,10:"`
	Count2            int             `validate:"int64,0:200"`
	RewardRate        float64         `validate:"float,0:0.4"`
	MaxRewardRate     float64         `validate:"float,:0.4"`
	MinRewradRate     float64         `validate:"float,0:"`
	RewardRate2       float64         `validate:"float64,0:0.4"`
	RewardRate3       decimal.Decimal `validate:"decimal,0:0.4"`
	OrderUsername     string          `validate:"regex,^[\u4E00-\u9FA5a-zA-Z0-9_.]{0,40}$"`
	OrderUsername2    string          `validate:"regex,username"`
	OrderStatus       int             `validate:"range,[1,2,3,4]"`
	OrderStatusName   string          `validate:"in,[unpaid,paid,closed]"`
	MailTypeCheckBox  string          `validate:"func,inAndLength"`
	MailTypeCheckBox2 string          `validate:"function,inAndLength"`
}
