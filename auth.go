// Copyright (c) 2017 Jean-François PHILIPPE
// DSM Client in Go

package godsm

import (
	"net/http/cookiejar"
)

// Login Try to connect given user.
// if sid is true, use sid for session tracking, otherwise use cookie
func (c *GoDsmImpl) Login(account string, passwd string, sid bool) error {
	format := "cookie"
	if sid {
		format = "sid"
	} else {
		// Set a store for cookies
		cookieJar, err := cookiejar.New(nil)
		if nil != err {
			return err
		}
		c.httpClient.Jar = cookieJar
	}
	// TODO : create a uniq session
	if "" == c.session {
		c.session = "TEST"
	}
	data, err := c.getJSON("SYNO.API.Auth", 2, "login",
		map[string]string{
			"account": account,
			"passwd":  passwd,
			"session": c.session,
			"format":  format,
		},
		map[int]string{
			400: "No such account or incorrect password",
			401: "Account disabled",
			402: "Permission denied",
			403: "2-step verification code required",
			404: "Failed to authenticate 2-step verification code",
		},
	)
	if nil == err {
		// fetch sid
		jsonMap, ok := data.(map[string]interface{})
		if ok {
			c.sid = jsonMap["sid"].(string)
		}
	} else {
		// clear session ID
		c.session = ""
		c.httpClient.Jar = nil
		c.sid = ""
	}
	return err
}

// Logout logout current session.
func (c *GoDsmImpl) Logout() error {
	_, err := c.getJSON("SYNO.API.Auth", 1, "logout",
		map[string]string{
			"session": c.session,
		},
		map[int]string{
			400: "No such account or incorrect password",
			401: "Account disabled",
			402: "Permission denied",
			403: "2-step verification code required",
			404: "Failed to authenticate 2-step verification code",
		},
	)
	//clear session
	c.session = ""
	c.httpClient.Jar = nil
	c.sid = ""
	return err
}

// vi:set fileencoding=utf-8 tabstop=4 ai
