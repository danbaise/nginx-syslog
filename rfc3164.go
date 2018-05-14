package nginx_syslog

import (
	"regexp"
	"strconv"
)

type Rfc3164 struct {
	Pri int
	Header
	Msg
}

type Header struct {
	Timestamp string
	Hostname  string
}

type Msg struct {
	Tag     string
	Content string
}


func NewRfc3164(data string) *Rfc3164 {
	re := regexp.MustCompile(`^<(\d{1,3})>(\w{3}\s\d{2}\s\d{1,2}:\d{1,2}:\d{1,2})\s(\S+)\s(\S.*?):\s(.*)`)
	result := re.FindStringSubmatch(data)
	pri, err := strconv.Atoi(result[1])
	checkError(err)
	return &Rfc3164{Pri:pri,Header:Header{result[2],result[3]},Msg:Msg{result[4], result[5]}}
}
