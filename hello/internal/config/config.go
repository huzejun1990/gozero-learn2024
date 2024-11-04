package config

import "github.com/zeromicro/go-zero/rest"

//type Config struct {
//	rest.RestConf
//	Address string
//}

//type Config struct {
//	//	rest.RestConf
//	Name    string
//	Host    string
//	Port    int
//	Address string `json:"address"`
//	Prop    string `json:"myProp"`
//	// optional 默认值是零值
//	NoConfStr        string `json:"noConfStr,optional"`
//	NoConfStrDefault string `json:"noConfStrDefault,default=默认值"`
//}

type Config struct {
	rest.RestConf
	DataBase DataBase
	// optional 默认值是零值
	NoConfStr        string `json:"noConfStr,optional"`
	NoConfStrDefault string `json:"noConfStrDefault,default=默认值"`
}

type DataBase struct {
	Url string `json:"url"`
}
