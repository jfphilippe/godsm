package main

import (
	"encoding/json"
	"fmt"
	"github.com/jfphilippe/godsm"
)

func main() {
	dsm, err := godsm.NewDSM("http://nas.houtraits.lan:5000/webapi")
	if nil == err {
		//	dsm.GetAPIInfo("all")
		dsm.LoadAllAPIInfo()
		//
		info, err := dsm.APIInfo("SYNO.DSM.User")
		if nil == err {
			fmt.Println(info)
		} else {
			fmt.Println(err)
		}
		info, err = dsm.APIInfo("SYNO.FileStation.CheckExist")
		if nil == err {
			fmt.Println(info)
		} else {
			fmt.Println(err)
		}
		err = dsm.Login("jeff", "fr69au98co02", true)
		if nil != err {
			fmt.Println(err)
		}
		data, err := dsm.System().Users(1, 3)
		if nil != err {
			fmt.Println(err)
		} else {
			b, err := json.MarshalIndent(data, "", "  ")
			if nil == err {
				fmt.Println(string(b))
			}
		}
		data, err = dsm.System().User("corentin")
		if nil != err {
			fmt.Println(err)
		} else {
			b, err := json.MarshalIndent(data, "", "  ")
			if nil == err {
				fmt.Println(string(b))
			}
		}
		data, err = dsm.Download().Tasks(0, 0)
		if nil != err {
			fmt.Println(err)
		} else {
			b, err := json.MarshalIndent(data, "", "  ")
			if nil == err {
				fmt.Println(string(b))
			}
		}
		err = dsm.Logout()
		if nil != err {
			fmt.Println(err)
		}
	}
}
