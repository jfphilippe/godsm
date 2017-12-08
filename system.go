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
		map[string]string{},
		map[int]string{},
	)
	if nil != err {
		return nil, err
	}
	return data.(map[string]interface{}), nil
}

// SystemUtilization get system utilization infos
func (c *GoDsmImpl) StorageInfo() (map[string]interface{}, error) {
	data, err := c.getJSON("SYNO.Storage.CGI.Storage", 1, "load_info",
		map[string]string{},
		map[int]string{},
	)
	if nil != err {
		return nil, err
	}
	return data.(map[string]interface{}), nil
}

// DsmInfo return Generic System info
func (c *GoDsmImpl) DsmInfo() (map[string]interface{}, error) {
	data, err := c.getJSON("SYNO.DSM.Info", 2, "getinfo",
		map[string]string{},
		map[int]string{},
	)
	if nil != err {
		return nil, err
	}
	return data.(map[string]interface{}), nil
}

// Networks return Network Configuration
func (c *GoDsmImpl) Network() (map[string]interface{}, error) {
	data, err := c.getJSON("SYNO.DSM.Network", 2, "list",
		map[string]string{},
		map[int]string{},
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
		map[int]string{},
	)
	if nil != err {
		return nil, err
	}
	return data.(map[string]interface{}), nil
}

// vi:set fileencoding=utf-8 tabstop=4 ai
