### Background（背景）

 在实际的项目中，我们以go-zero 作为gRpc 网关，由网关将请求路由至后端gRpc 微服务。采用nacos作为配置中心和服务发现。打造高可用的微服务网关。整体架构如下：



![https://github.com/alanjiang/go-zero-rpc-nacos/blob/master/grpc.jpg](https://github.com/alanjiang/go-zero-rpc-nacos/blob/master/grpc.jpg)



但面临的问题是： go-zero默认以 ETCD作为服务发现， 而nacos 作为服务发现需要自己去实现，文档在这块也是缺失的，参考的资料凤毛鳞角，给广大开发者带来困拢。

于是看到开源项目：

https://github.com/zeromicro/zero-contrib/tree/main/zrpc/registry

问题是开源项目只给出了 用户名/密码这种连接nacos的认证方式，并不支持 accessKey/secretKey这种方式。

而阿里云的accessKey/secretKey连接方式安全性较高，是推荐的一种连接方式。

故在开源项目https://github.com/zeromicro/zero-contrib/tree/main/zrpc/registry的基础上，我们作了修改

以布了新的服务。

## 如何引入？

` _ "github.com/alanjiang/go-zero-rpc-nacos"`

## 网关代码参见：

- main.go

```go
package main

import (
    "github.com/zeromicro/go-zero/core/logx"
    "net/http"
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/gateway"
	 _ "github.com/alanjiang/go-zero-rpc-nacos"
)

var configFile = flag.String("f", "etc/gateway.yaml", "config file")


func main() {
	flag.Parse()

	var c gateway.GatewayConf
	conf.MustLoad(*configFile, &c)
	gw := gateway.MustNewServer(c)
	defer gw.Stop()

	gw.Start()
}

	
```

## 



- etc/gateway.yaml

```yaml

Target: nacos://accessKey:secretKey@nacos服务器域名:8848/logistic.rpc?namespaceid=空间ID&timeout=13000ms
    Mappings:
        - Method: put
          Path: /logistic/query
          RpcPath: logistic.LogisticService/Query               
```



通过以上的配置就可实现查询 服务名为 ：logistic.rpc 的快递服务。

注：网关不需要注册客户端。 

引入  _ "github.com/alanjiang/go-zero-rpc-nacos" 会自动执行 builder.go 的 init 方法。

```
func init() {

    fmt.Print("----> nacos init <-----")
	  resolver.Register(&builder{})
}
```





## 关注作者：

抖音号： 80758963124



