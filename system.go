// Copyright (c) 2017 Jean-François PHILIPPE
// DSM Client in Go
// This file contains undocumented Web API
// To stay safe, most API are read only.

package godsm

import (
	"strconv"
)

// SystemUtilization get system utilization infos
func (c *GoDsmImpl) SystemUtilization() (map[string]interface{}, error) {
	data, err := c.getJSON("SYNO.Core.System.Utilization", 1, "get",
		nil,
		nil,
	)
	if nil != err {
		return nil, err
	}
	return data.(map[string]interface{}), nil
}

// StorageInfo get system utilization infos
func (c *GoDsmImpl) StorageInfo() (map[string]interface{}, error) {
	data, err := c.getJSON("SYNO.Storage.CGI.Storage", 1, "load_info",
		nil,
		nil,
	)
	if nil != err {
		return nil, err
	}
	return data.(map[string]interface{}), nil
}

// DsmInfo return Generic System info
func (c *GoDsmImpl) DsmInfo() (map[string]interface{}, error) {
	data, err := c.getJSON("SYNO.DSM.Info", 2, "getinfo",
		nil,
		nil,
	)
	if nil != err {
		return nil, err
	}
	return data.(map[string]interface{}), nil
}

// Network return Network Configuration
func (c *GoDsmImpl) Network() (map[string]interface{}, error) {
	data, err := c.getJSON("SYNO.DSM.Network", 2, "list",
		nil,
		nil,
	)
	if nil != err {
		return nil, err
	}
	return data.(map[string]interface{}), nil
}

// Users return Users list
// set offset and limit to 0 to get all users
func (c *GoDsmImpl) Users(offset, limit int) (map[string]interface{}, error) {
	params := map[string]string{}
	if offset > 0 {
		params["offset"] = strconv.Itoa(offset)
	}
	if limit > 0 {
		params["limit"] = strconv.Itoa(limit)
	}
	data, err := c.getJSON("SYNO.Core.User", 1, "list",
		params,
		nil,
	)
	if nil != err {
		return nil, err
	}
	return data.(map[string]interface{}), nil
}

// User return User id
func (c *GoDsmImpl) User(name string) (map[string]interface{}, error) {
	data, err := c.getJSON("SYNO.Core.User", 1, "get",
		map[string]string{"name": name},
		nil,
	)
	if nil != err {
		return nil, err
	}
	return data.(map[string]interface{}), nil
}

// vi:set fileencoding=utf-8 tabstop=4 ai
