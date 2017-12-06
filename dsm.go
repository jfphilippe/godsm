// Copyright (c) 2017 Jean-François PHILIPPE
// DSM Client in Go

package godsm

import "fmt"

// GoDsm define interface
type GoDsm interface {
	APIInfo(api string) (*DsmAPIInfo, error)
	Login(user string, passwd string, sid bool) error
	LoadAllAPIInfo() error
	Logout() error
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

// DsmApiInfo memorize Api info
type DsmAPIInfo struct {
	Key           string
	Path          string
	RequestFormat string
	MinVersion    int
	MaxVersion    int
}

// vi:set fileencoding=utf-8 tabstop=4 ai
