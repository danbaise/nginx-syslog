package nginx_syslog

import (
	"regexp"
	"strings"
)

type Log struct {
	Host     string
	Identity string
	User     string
	Time     string
	Request  string
	Status   string
	Size     string
	Referer  string
	Agent    string
}

type Request struct {
	Method   string
	Path     string
	Protocol string
}

func NewLog(content string) *Log {
	/*
	　九列,每列之间是用空格分割的，每列的含义分别是客户端访问IP、用户标识、用户、访问时间、请求页面、请求状态、返回文件的大小、跳转来源、浏览器UA
		([^ ]*) ([^ ]*) ([^ ]*) (\[.*\]) (\".*?\") (-|[0-9]*) (-|[0-9]*) (\".*?\") (\".*?\")
	 */
	re := regexp.MustCompile(`([^ ]*) ([^ ]*) ([^ ]*) (\[.*\]) (\".*?\") (-|[0-9]*) (-|[0-9]*) (\".*?\") (\".*?\")`)
	result := re.FindStringSubmatch(content)
	return &Log{result[1], result[2], result[3], result[4], result[5], result[6], result[7], result[8], result[9]}
}

func NewRequest(content string) *Request {
	result := strings.Split(content, " ")
	return &Request{result[0],result[1],result[2]}
}
