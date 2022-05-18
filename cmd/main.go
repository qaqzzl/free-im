package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	var req struct {
		ChatroomID int64 `json:"ChatroomID" json:"ChatroomID"`
	}

	str := `{"ChatroomID":1}`
	data := []byte(str)

	if err := json.Unmarshal(data, &req); err != nil {
		fmt.Println(err)
	}

	fmt.Println(req.ChatroomID)
}
