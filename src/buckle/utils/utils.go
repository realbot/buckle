package utils

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func CheckMsg(s string, e error) {
	if e != nil {
		panic(s + e.Error())
	}
}
