package superChecker

import "regexp"

var compiledMap map[string]*regexp.Regexp

func init() {
	// for regex checking, inner basic
	compiledMap = make(map[string]*regexp.Regexp, 0)
	compiledMap["chineseonly"] = mustCompile("^[\u4E00-\u9FA5]*$")              // only chinese
	compiledMap["notnull"] = mustCompile("^[\\s\\S]+$")                         // not empty string
	compiledMap["username"] = mustCompile("^[\u4E00-\u9FA5a-zA-Z0-9_.]{0,40}$") // username,chinese,english character,'_','.'，lenggh in 40
	compiledMap["number"] = mustCompile("^[0-9]+$")                             // number of positive integer
	compiledMap["decimal"] = mustCompile("^\\d+\\.[0-9]+$")                     // decimal, 2.2
	compiledMap["mobilephone"] = mustCompile("^1[0-9]{10}$")                    // mobilephone, length is 10, 13802930292
	compiledMap["telephone"] = mustCompile("^[0-9]{8}$")                        // telephone,length is 8,consist of 0-9 numbers,88501918
}

func mustCompile(regex string) *regexp.Regexp {
	r, e := regexp.Compile(regex)
	if e != nil {
		panic(e)
	}
	return r
}
