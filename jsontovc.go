package golib_vcard

import (
	"encoding/json"
	"fmt"
)

//Json格式转为Vcard
func JsonToVcard(jsonstr string) string {
	var c Card
	var jtovc string
	if err := json.Unmarshal([]byte(jsonstr), &c); err == nil {
		b, err := Marshal(c)
		if err != nil {
			fmt.Println("json转vcard错误")
		}
		jtovc = string(b)

	} else {
		fmt.Println("json转vcard错误")
	}
	return jtovc
}

//Json转为Valendar格式
func JsonToVcalendar(jsonstr string) string {
	var c Calendar
	var jtovc string
	if err := json.Unmarshal([]byte(jsonstr), &c); err == nil {
		b, err := Marshal(c)
		if err != nil {
			fmt.Println(err)
		}
		jtovc = string(b)

	} else {
		fmt.Println("json转Vcalendar错误2")
	}
	return jtovc
}

//Json转为短信格式
func JsonToBjson(jsonstr string) string {
	var c BJSON
	var jtovc string
	if err := json.Unmarshal([]byte(jsonstr), &c); err == nil {
		b, err := Marshal(c)
		if err != nil {
			fmt.Println(err)
		}
		jtovc = string(b)

	} else {
		fmt.Println("json转Vcalendar错误2")
	}
	return jtovc
}
