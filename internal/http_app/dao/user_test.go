package dao

import (
	"fmt"
	"testing"
)

func TestUser_GetByUID(t *testing.T) {
	res, err := User.GetByUID(1, "nickname", "avatar")
	fmt.Println(res, err)
}
