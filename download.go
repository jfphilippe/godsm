// Copyright (c) 2018 Jean-François PHILIPPE
// DSM Client in Go
// DownloadStation API

package godsm

import (
	"strconv"
)

// DlInfo get downloadStation info.
func (c *GoDsmImpl) DlInfo() (map[string]interface{}, error) {
	data, err := c.getJSON("SYNO.DownloadStation.Info", 1, "getinfo",
		nil,
		nil,
	)
	if nil != err {
		return nil, err
	}
	return data.(map[string]interface{}), nil
}

// DlInfoConfig get downloadStation info.
func (c *GoDsmImpl) DlInfoConfig() (map[string]interface{}, error) {
	data, err := c.getJSON("SYNO.DownloadStation.Info", 1, "getconfig",
		nil,
		nil,
	)
	if nil != err {
		return nil, err
	}
	return data.(map[string]interface{}), nil
}

// DlSetServer set downloadStation info.
// can only be used by admin users.
// available configs, one or more of :
// - bt_max_download (in KB/s, 0 means unlimited)
// - bt_max_upload (in KB/s, 0 means unlimited)
// - emule_max_download (in KB/s, 0 means unlimited)
// - emule_max_upload (in KB/s, 0 means unlimited)
// - nzb_max_download (in KB/s, 0 means unlimited)
// - http_max_download (in KB/s, 0 means unlimited)
// - ftp_max_download (in KB/s, 0 means unlimited)
// - emule_enabled
// - unzip_service_enabled
// - default_destination
// - emule_default_destination
//
func (c *GoDsmImpl) DlSetServer(cfg map[string]string) error {
	params := map[string]string{}
	// will only retains those two
	values := []string{"bt_max_download",
		"bt_max_upload",
		"emule_max_download",
		"emule_max_upload",
		"nzb_max_download",
		"http_max_download",
		"ftp_max_download",
		"emule_enabled",
		"unzip_service_enabled",
		"default_destination",
		"emule_default_destination"}
	for _, v := range values {
		val, ok := cfg[v]
		if ok {
			params[v] = val
		}
	}
	_, err := c.getJSON("SYNO.DownloadStation.Info", 1, "setserverconfig",
		params,
		nil,
	)
	return err
}

// ===========================================================================
// DownloadStation Schedule API
// ===========================================================================

// DlSchedConfig get downloadStation Schedule Config.
func (c *GoDsmImpl) DlSchedConfig() (map[string]interface{}, error) {
	data, err := c.getJSON("SYNO.DownloadStation.Schedule", 1, "getconfig",
		nil,
		nil,
	)
	if nil != err {
		return nil, err
	}
	return data.(map[string]interface{}), nil
}

// DlSchedSetConfig set server config
// can only be used by admin users.
// cfg should contains one or more of :
// - enabled if download schedule is enabled
// - emule_enabled if eMule download schedule is enabled.
func (c *GoDsmImpl) DlSchedSetConfig(cfg map[string]bool) error {
	// Convert params to string
	params := map[string]string{}
	// will only retains those two
	values := []string{"enabled", "emule_enabled"}
	for _, v := range values {
		val, ok := cfg[v]
		if ok {
			params[v] = strconv.FormatBool(val)
		}
	}
	_, err := c.getJSON("SYNO.DownloadStation.Schedule", 1, "setconfig",
		params,
		nil,
	)
	return err
}

// ===========================================================================
// DownloadStation Task API
// ===========================================================================

// Tasks get Tasks List
func (c *GoDsmImpl) Tasks(offset, limit int) (map[string]interface{}, error) {
	params := map[string]string{}
	if offset > 0 {
		params["offset"] = strconv.Itoa(offset)
	}
	if limit > 0 {
		params["limit"] = strconv.Itoa(limit)
	}
	// Details may be retrives through GetInfo
	// params["additional"] = "detail,file,transfer"
	data, err := c.getJSON("SYNO.DownloadStation.Task", 1, "list",
		params,
		nil,
	)
	if nil != err {
		return nil, err
	}
	return data.(map[string]interface{}), nil
}
