package utils

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

func CheckErrorMsg(s string, e error) {
	if e != nil {
		panic(s + e.Error())
	}
}
