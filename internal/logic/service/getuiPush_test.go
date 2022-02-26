package service

import (
	"fmt"
	"testing"
)

func TestGetToken(t *testing.T) {
	fmt.Println(Push.getToken())
}

func TestPush(t *testing.T) {
	fmt.Println(Push.push("b496d6ff9d7e7c64db10c2a501abfd3f"))
}
