package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func checkLogin(token, apiUrl string) error {
	req, err := http.NewRequest("POST", apiUrl, nil)
	if err != nil {
		return err
	}

	req.AddCookie(&http.Cookie{Name: "oa_token", Value: token})
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("调用判断是否登陆接口:/api/checklogin失败")
	}

	fmt.Println(string(body))

	/*
		type CheckloginData struct {
			UserId     int `json:"userId,omitempty"`
			LoginTime  string `json:"loginTime,omitempty"`
			ExpireTime string `json:"expireTime,omitempty"`
			LoginIp    string `json:"loginIp,omitempty"`
			Email      string `json:"email,omitempty"`
			Name       string `json:"name,omitempty"`
			Mobile     string `json:"mobile,omitempty"`
			Avatar     string `json:"avatar,omitempty"`
		}
	*/
	type CheckloginInfo struct {
		Code    int    `json:"code" binding:"required"`
		Message string `json:"message,omitempty"`
		//Data    CheckloginData `json:"data,omitempty"`
	}

	var checkloginInfo CheckloginInfo
	if err := json.Unmarshal(body, &checkloginInfo); err != nil {
		return errors.New("调用判断是否登陆接口:/api/checklogin json解析失败")
	}

	if checkloginInfo.Code != 0 {
		fmt.Printf("message: %s\n", checkloginInfo.Message)
		return errors.New("调用判断是否登陆接口:/api/checklogin失败")
	}
	
	return nil
}
