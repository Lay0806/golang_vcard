package golib_vcard

type BJSON struct {
	Profile string   `vdir:"json,profile"`
	SMSDATA []string `vdir:SMSDATA`
	Contact string
	DATE    string
}
