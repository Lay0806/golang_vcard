package golib_vcard

type Card struct {
	Profile       string `vdir:"vcard,profile"`
	Version       string
	FormattedName string `vdir:"fn"`
	Name          Name   `vdir:"n"`
	NickName      []string
	Birthday      string `vdir:"bday"`
	Anniversary   string
	Addresses     []Address `vdir:"adr"`
	Label         []TypedValue
	Telephones    []TypedValue `vdir:"tel"`
	Email         []TypedValue
	Url           []TypedValue
	Related       []TypedValue
	Cday          []Date
	Title         string
	Role          string
	Org           string
	Categories    []string
	Note          string
	URL           string
	Photo         Photo

	Rev    string
	ProdId string
	Uid    string

	XICQ    string `vdir:"x-icq"`
	XSkype  string `vdir:"x-skype"`
	XAIM    string `vdir:"x-aim"`
	XJABBER string `vdir:"x-jabber"`
	IMPP    []IMPP

	Sensitivity string `vdir:sensitivity`
	Folder      string `vdir:folder`
	Gender      string `vdir:"x-wab-gender"`
	DisplayName string `vdir:"display-name"`
	SelfURL     string `vdir:"selfurl"`
	Starred     string `vdir:starred`
	ShortName   string `vdir:short-name`
	IsShare     string `vdir:isshare`
	BrInterval  string `vdir:br-interval`
	ArInterval  string `vdir:ar-interval`
	IsRemind    string `vdir:isremind`
}

type IMPP struct {
	Type         []string `vdir:",param"`
	XServiceType string   `vdir:"x-service-type,param"`
	Value        string
}

type Name struct {
	FamilyName        []string
	GivenName         []string
	AdditionalNames   []string
	HonorificNames    []string
	HonorificSuffixes []string
}

type Photo struct {
	Encoding  string `vdir:",param"`
	MediaType string `vdir:",param"`
	Type      string `vdir:",param"`
	Value     string `vdir:",param"`
	Data      string
}

type Address struct {
	Type            []string `vdir:",param"`
	Label           string   `vdir:",param"`
	PostOfficeBox   string
	ExtendedAddress string
	Street          string
	Locality        string
	Region          string
	PostalCode      string
	CountryName     string
}

type Date struct {
	Type          []string `vdir:",param"`
	Value         string   `vdir:"value"`
	IsLunar       string   `vdir:"isLunar"`
	IsRemind      string   `vdir:"isRemind"`
	RemindContent string   `vdir:"remindContent"`
	Label         string   `vdir:",param"`
}

type TypedValue struct {
	Type  []string `vdir:",param"`
	Value string
}
