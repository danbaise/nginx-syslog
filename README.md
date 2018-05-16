

### go 处理 nginx syslog rfc3614

#### nginx 配置
> access_log syslog:server=192.168.1.44:514,facility=local7,tag=nginx01,severity=info;

#### Example

> ```go
> package main
>
> import (
>    "github.com/danbaise/nginx-syslog"
> )
>
> func main() {
>    nginx_syslog.NewParser().Handle()
> }
> ```

