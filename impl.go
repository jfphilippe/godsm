// Copyright (c) 2017 Jean-François PHILIPPE
// DSM Client in Go

package godsm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

// GoDsmImpl implements interface
type GoDsmImpl struct {
	sid        string // current sid , "" if not authenticated
	session    string // current session name
	dsmURL     *url.URL
	httpClient *http.Client
	apis       map[string]*DsmAPIInfo // Cache of DsmAPIInfo
}

// NewDSM Build a new DSM
func NewDSM(dsmUrl string) (GoDsm, error) {
	u, err := url.ParseRequestURI(dsmUrl)
	if nil != err {
		return nil, err
	}
	// Force / at end of path.
	if !strings.HasSuffix(u.Path, "/") {
		u.Path = u.Path + "/"
	}
	dsm := GoDsmImpl{sid: "", dsmURL: u, httpClient: &http.Client{}, apis: make(map[string]*DsmAPIInfo)}
	// Bootstrap SYNO.API.Info
	dsm.apis["SYNO.API.Info"] = &DsmAPIInfo{Key: "SYNO.API.Info", Path: "query.cgi", RequestFormat: "JSON", MinVersion: 1, MaxVersion: 1}
	return &dsm, nil
}

// LoadAllAPIInfo load all DsmAPIInfo in cache.
func (c *GoDsmImpl) LoadAllAPIInfo() error {
	data, err := c.get("SYNO.API.Info", 1, "query",
		map[string]string{
			"query": "all",
		},
		map[int]string{},
	)
	if nil == err {
		// To json and Back to array of structs.
		stInfo, _ := json.Marshal(data.(map[string]interface{}))
		keys := make(map[string]DsmAPIInfo, 0)
		err = json.Unmarshal(stInfo, &keys)
		// fmt.Println(keys)
		if nil == err {
			// store/replace items in cache
			for k, v := range keys {
				v.Key = k
				c.apis[k] = &DsmAPIInfo{Key: k, Path: v.Path, RequestFormat: v.RequestFormat, MinVersion: v.MinVersion, MaxVersion: v.MaxVersion}
			}
		}
		//fmt.Println(c.apis)

	}
	return err
}

// APIInfo return an API Info, call LoadAllAPIInfo if needed
func (c *GoDsmImpl) APIInfo(api string) (*DsmAPIInfo, error) {
	info, found := c.apis[api]
	if !found {
		err := c.LoadAllAPIInfo()
		if nil != err {
			return nil, err
		}
		info, found = c.apis[api]
		if !found {
			return nil, &DsmError{Code: 00, Msg: fmt.Sprintf("Unknown API '%s'", api)}
		}
	}
	return info, nil
}

// send send a query
func (c *GoDsmImpl) get(api string, version int, method string, params map[string]string, respErrors map[int]string) (interface{}, error) {
	apiInfo, err := c.APIInfo(api)
	if nil != err {
		return nil, err
	}
	// Build URL
	u := url.URL{Scheme: c.dsmURL.Scheme, Host: c.dsmURL.Host, Path: c.dsmURL.Path + apiInfo.Path}
	query := u.Query()
	// set commons params
	query.Set("api", api)
	query.Set("version", strconv.Itoa(version))
	query.Set("method", method)
	// complete with given params.
	for k, v := range params {
		query.Set(k, v)
	}
	u.RawQuery = strings.Replace(query.Encode(), "+", "%20", -1)

	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse json
	var data interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	// analyse response
	// TODO check data type (map[string]interface{})
	jsonMap := data.(map[string]interface{})
	success := jsonMap["success"].(bool)
	if success {
		return jsonMap["data"], nil
	}
	// check error code and convert it to DsmError
	code := c.errCode(jsonMap)

	return nil, c.errorFromCode(code, respErrors)
}

// errCode extract error code from response
func (c *GoDsmImpl) errCode(data map[string]interface{}) int {
	code := 0
	err, found := data["error"]
	if found {
		// TODO : check err type !
		code = int(err.(map[string]interface{})["code"].(float64))
	}
	return code
}

// errorFromCode convert Error code into an DsmError.
func (c *GoDsmImpl) errorFromCode(code int, respErrors map[int]string) error {
	// Try errors for given service
	msg, found := respErrors[code]
	if !found {
		// Try commons errors
		switch code {
		case 100:
			msg = "Unknown error"
		case 101:
			msg = "Invalid parameter"
		case 102:
			msg = "The requested API does not exist"
		case 103:
			msg = "The requested method does not exist"
		case 104:
			msg = "The requested version does not support the functionality"
		case 105:
			msg = "The logged in session does not have permission"
		case 106:
			msg = "Session timeout"
		case 107:
			msg = "Session interrupted by duplicate login"
		default:
			msg = "Unknown Error"
		}
	}
	return &DsmError{Code: code, Msg: msg}
}

// Login Try to connect given user.
// if sid is true, use sid for session tracking, otherwise use cookie
func (c *GoDsmImpl) Login(user string, passwd string, sid bool) error {
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
	c.session = "TEST"
	data, err := c.get("SYNO.API.Auth", 2, "login",
		map[string]string{
			"account": user,
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
		fmt.Println(data)
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
	_, err := c.get("SYNO.API.Auth", 1, "logout",
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
