package golib_vcard

import (
	"encoding/json"
	"fmt"
)

//Vcard 转为Json
func VcardToJson(tVcard string) string {
	var c Card
	if err := Unmarshal([]byte(tVcard), &c); err != nil {
		fmt.Println(err)
		fmt.Println("转码失败")
	}
	d, _ := json.Marshal(c)
	return string(d)
}

//Vcalendar转为Json
func VcalendarToJson(tVcalendar string) string {
	var ca Calendar
	if err := Unmarshal([]byte(tVcalendar), &ca); err != nil {
		fmt.Printf("转码失败")
	}
	d, _ := json.Marshal(ca)
	return string(d)
}

//短信格式转为Json
func BjsonToJson(tBjson string) string {
	var ca BJSON
	if err := Unmarshal([]byte(tBjson), &ca); err != nil {
		fmt.Printf("转码失败")
	}
	d, _ := json.Marshal(ca)
	return string(d)
}
