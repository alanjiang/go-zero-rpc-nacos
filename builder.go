package nacos

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/pkg/errors"
	"google.golang.org/grpc/resolver"
)

func init() {

    fmt.Print("----> nacos init <-----")
	resolver.Register(&builder{})
}

// schemeName for the urls
// All target URLs like 'nacos://.../...' will be resolved by this resolver
const schemeName = "nacos"

// builder implements resolver.Builder and use for constructing all consul resolvers
type builder struct{}

func (b *builder) Build(url resolver.Target, conn resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	fmt.Print("==> grpc-go/resolver/resolver.go Build 接口实现 <====")
	tgt, err := parseURL(url.URL)
	if err != nil {
		return nil, errors.Wrap(err, "Wrong nacos URL")
	}

	host, ports, err := net.SplitHostPort(tgt.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed parsing address error: %v", err)
	}
	port, _ := strconv.ParseUint(ports, 10, 16)

	fmt.Print("********Host, port ******")

	fmt.Print("**port:", port)

	fmt.Print("**host:", host)

	fmt.Print("********Host, port ******")

	sc := []constant.ServerConfig{
        		*constant.NewServerConfig(host, port, constant.WithContextPath("/nacos")),
     }


    /*
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(host, port),
	}

	cc := &constant.ClientConfig{
		AppName:     tgt.AppName,
		NamespaceId: tgt.NamespaceID,
		AccessKey:    tgt.User,
		SecretKey:    tgt.Password,
		TimeoutMs:   uint64(tgt.Timeout),
		NotLoadCacheAtStart:  tgt.NotLoadCacheAtStart,
		UpdateCacheWhenEmpty: tgt.UpdateCacheWhenEmpty,
	} */


	// start
	cc :=  constant.NewClientConfig(
                	constant.WithTimeoutMs(10*1000),
                	constant.WithBeatInterval(2*1000),
                	constant.WithNotLoadCacheAtStart(true),
                	constant.WithAccessKey(tgt.User),
                	constant.WithSecretKey(tgt.Password),
                	constant.WithNamespaceId(tgt.NamespaceID),
                	constant.WithOpenKMS(false),
                	//constant.WithKMSVersion(constant.KMSv3),
                	constant.WithRegionId("cn-hangzhou-e"),
      )

	// end

	if tgt.CacheDir != "" {
		cc.CacheDir = tgt.CacheDir
	}
	if tgt.LogDir != "" {
		cc.LogDir = tgt.LogDir
	}
	if tgt.LogLevel != "" {
		cc.LogLevel = tgt.LogLevel
	}


    fmt.Print("===>nacos:", tgt)

     fmt.Print("nacos<=====")

	cli, err := clients.NewNamingClient(vo.NacosClientParam{
		ServerConfigs: sc,
		ClientConfig:  cc,
	})




	if err != nil {
		return nil, errors.Wrap(err, "Couldn't connect to the nacos API")
	}

	ctx, cancel := context.WithCancel(context.Background())
	pipe := make(chan []string)
    tgt.GroupName = "prod"



	go cli.Subscribe(&vo.SubscribeParam{
		ServiceName:       tgt.Service,
		Clusters:          tgt.Clusters,
		GroupName:         tgt.GroupName,
		SubscribeCallback: newWatcher(ctx, cancel, pipe).CallBackHandle, // required
	})

	go populateEndpoints(ctx, conn, pipe)

	return &resolvr{cancelFunc: cancel}, nil
}

// Scheme returns the scheme supported by this resolver.
// Scheme is defined at https://github.com/grpc/grpc/blob/master/doc/naming.md.
func (b *builder) Scheme() string {
	return schemeName
}
