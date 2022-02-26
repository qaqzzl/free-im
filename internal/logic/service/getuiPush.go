package service

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type push struct {
	apiURL          string
	appID           string
	appKey          string
	appSecret       string
	masterSecret    string
	token           string
	tokenExpireTime string
}

var Push = push{
	apiURL:       "https://restapi.getui.com/v2/A98ePyLcrI7B80nTQXs364",
	appID:        "A98ePyLcrI7B80nTQXs364",
	appKey:       "Hk6KEAlMZU7M6xEnuwJej5",
	appSecret:    "gSVvfn4oKS8sBWjB7IsHa7",
	masterSecret: "j4NItxn7jPACocP0LCwSc8",
}

func (p *push) push(cid string) error {
	var api = "/push/single/cid"
	var dataStruct struct {
		RequestID string `json:"request_id"`
		Settings  struct {
			TTL int `json:"ttl"`
		} `json:"settings"`
		Audience struct {
			Cid []string `json:"cid"`
		} `json:"audience"`
		PushMessage struct {
			Notification struct {
				Title     string `json:"title"`
				Body      string `json:"body"`
				ClickType string `json:"click_type"`
				URL       string `json:"url"`
			} `json:"notification"`
		} `json:"push_message"`
	}
	dataStruct.RequestID = uuid.New().String()
	dataStruct.Audience.Cid = append(dataStruct.Audience.Cid, cid)

	dataStr, err := json.Marshal(dataStruct)
	if err != nil {
		return err
	}
	var data = bytes.NewBuffer(dataStr)
	resp, err := http.Post(p.apiURL+api, "application/json;charset=utf-8", data)
	token, err := p.getToken()
	if err != nil {
		return err
	}
	resp.Header.Set("token", token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var respData struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if json.Unmarshal(body, &respData) != nil {
		return errors.New("json 解析失败")
	}
	if respData.Code != 0 {
		return errors.New(respData.Msg)
	}
	return nil
}

func (p *push) getToken() (token string, err error) {
	var api = "/auth"
	var timestamp = int(time.Now().UnixNano() / 1e6)
	sign := sha256.Sum256([]byte(p.appKey + strconv.Itoa(timestamp) + p.masterSecret))
	signStr := hex.EncodeToString(sign[:])
	dataStr := "{\"sign\": \"" + signStr + "\",\"timestamp\": \"" + strconv.Itoa(timestamp) + "\", \"appkey\": \"" + p.appKey + "\"}"
	var data = bytes.NewBuffer([]byte(dataStr))
	resp, err := http.Post(p.apiURL+api, "application/json;charset=utf-8", data)
	if err != nil {
		return token, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var respData struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			ExpireTime string `json:"expire_time"`
			Token      string `json:"token"`
		} `json:"data"`
	}
	if json.Unmarshal(body, &respData) != nil {
		return token, errors.New("json 解析失败")
	}
	if respData.Code != 0 {
		return token, errors.New(respData.Msg)
	}
	p.token = respData.Data.Token
	p.tokenExpireTime = respData.Data.ExpireTime
	return p.token, err
}
