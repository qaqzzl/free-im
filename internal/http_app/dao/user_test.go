package dao

import (
	"fmt"
	"testing"
)

func TestUser_Get(t *testing.T) {
	res, err := User.Get(1, "nickname", "avatar")
	fmt.Println(res, err)
}
