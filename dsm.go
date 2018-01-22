// Copyright (c) 2018 Jean-François PHILIPPE
// DSM Client in Go

package godsm

import "fmt"

// GoDsm define interface
type Dsm interface {
	APIInfo(api string) (*DsmAPIInfo, error)
	Login(user string, passwd string, sid bool) error
	LoadAllAPIInfo() error
	Logout() error
	System() System
}

// System get System informations
type System interface {
	DsmInfo() (map[string]interface{}, error)
	Network() (map[string]interface{}, error)
	StorageInfo() (map[string]interface{}, error)
	SystemUtilization() (map[string]interface{}, error)
	Users(offset, limit int) (map[string]interface{}, error)
	User(name string) (map[string]interface{}, error)
}

// DsmError error
type DsmError struct {
	// Error Code
	Code int
	// Error Message
	Msg string
}

// Error Interface implementation
func (d *DsmError) Error() string {
	return fmt.Sprintf("Error '%s' (%d)", d.Msg, d.Code)
}

// DsmAPIInfo memorize Api info
type DsmAPIInfo struct {
	Key           string
	Path          string
	RequestFormat string
	MinVersion    int
	MaxVersion    int
}

// vi:set fileencoding=utf-8 tabstop=4 ai
