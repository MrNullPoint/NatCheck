# NAT 类型检测

### 使用方式
编译之后，指定 ip 或不指定使用内置的 stun 服务器 ip 运行

```shell
$ go build -o natcheck main.go
$ ./natcheck 或者 ./natcheck stun_server_ip 
```

### API
使用 `go-get`
