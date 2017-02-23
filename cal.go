package golib_vcard

//日历
type Calendar struct {
	Profile  string     `vdir:"vcalendar,profile"` //类型：Vcar或vcalendar
	Version  string     //版本号
	ProdId   string     //ID
	Timezone []Timezone `vdir:"vtimezone,object"` //时区
	Events   []Event    `vdir:"vevent,object"`    //事件
	ToDos    []Todo     `vdir:"vtodo,object"`
	Journals []Journal  `vdir:"vjournal,object"`
	FreeBusy []FreeBusy `vdir:"vfreebusy,object"`
}

type Event struct {
	Profile      string `vdir:"vevent,profile"`
	UID          string
	DTStamp      DateTimeValue
	Organizer    Person
	DTStart      DateTimeValue
	DTEnd        DateTimeValue
	Location     string
	Summary      string
	Categories   []string
	Description  string
	Method       string
	Status       string
	Class        string
	Sequence     string
	Created      string
	LastModified string  `vdir:"last-modified"`
	Alarms       []Alarm `vdir:",object"`
	RRule        RecurrenceRule
}

type Person struct {
	CommonName string `vdir:"cn"`
	Url        string
}

type Alarm struct {
	Profile     string `vdir:"valarm,profile"`
	Trigger     string
	Action      string
	Description string
	Repeat      string
	Duration    string
}

type Todo struct {
	Profile  string `vdir:"vtodo,profile"`
	DTStamp  DateTimeValue
	Sequence string
	UID      string
	Due      string
	Status   string
	Summary  string
	Alarms   []Alarm `vdir:",object"`
}

type Journal struct {
	Profile     string `vdir:"vjournal,profile"`
	DTStamp     DateTimeValue
	UID         string
	Organizer   Person
	Status      string
	Class       string
	Categories  []string
	Description string
}

type FreeBusy struct {
	Profile   string `vdir:"vfreebusy,profile"`
	Organizer Person
	DTStart   DateTimeValue
	DTEnd     DateTimeValue
	FreeBusy  []string `vdir:",multiple"`
	Url       string
}

type Timezone struct {
	Profile  string `vdir:"vtimezone,profile"`
	TZId     string
	Daylight []TimeZoneInfo `vdir:",object"`
	Standard []TimeZoneInfo `vdir:",object"`
}

type TimeZoneInfo struct {
	TZOffsetFrom string
	TZOffsetTo   string
	TZName       string
	DTStart      string
	RRule        RecurrenceRule
}

type RecurrenceRule struct {
	Rule1 string
	Rule2 string
	Rule3 string
	Rule4 string
	Rule5 string
}

type DateTimeValue struct {
	TZId  string `vdir:",param"`
	Type  string `vdir:"value,param"`
	Value string
}
