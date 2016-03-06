package utils

import (
    "os/user"
)

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

func CurrentUser() *user.User {
    usr, err := user.Current()
    if err != nil {
        panic(err)
    }
    return usr
}