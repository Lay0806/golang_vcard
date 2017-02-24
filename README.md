# Vcard/Vcalendar and Json conversion to 

本模块是使用golang制作的 Vard/Valendar格式与Json格式之间的相互转换。
##### 下载说明
  * go get远程git库
  `go get -insecure git@github.com:Lay0806/golang_vcard.git`

##### 使用说明
- import
    ```go
         import ("git@github.com:Lay0806/golang_vcard.git")
    ```
- vard/valendar转为Json
   ```go
   
    //Vcard 转为Json
    func VcardToJson(tVcard string) string {
	    var c Card
	    if err := Unmarshal([]byte(tVcard), &c); err != nil {
		    fmt.Printf("转码失败")
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
   ```
- Json转为Vcard/Valendar
   ```go
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
   ```
   

