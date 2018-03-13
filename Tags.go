package superChecker
type T struct {
	UserName string `superChecker:"userName"`
	Password string `superChecker:"password"`
	Phone string `superChecker:"mobilephone|telephone"`
	Text string `superChecker:"length,chineseOnly"`
}