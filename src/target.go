package nacos

import (
	"net/url"
	"os"
	"strings"
	"time"
    "fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/mapping"
)

type target struct {
	Addr        string        `key:",optional"`
	User        string        `key:",optional"`
	Password    string        `key:",optional"`
	AcccessKey  string        `key:",optional"`
	SecretKey    string       `key:",optional"`
	Service     string        `key:",optional"`
	GroupName   string        `key:",optional"`
	Clusters    []string      `key:",optional"`
	NamespaceID string        `key:"namespaceid,optional"`
	Timeout     time.Duration `key:"timeout,optional"`
	AppName     string        `key:"appName,optional"`
	LogLevel    string        `key:",optional"`
	LogDir      string        `key:",optional"`
	CacheDir    string        `key:",optional"`
	NotLoadCacheAtStart  bool `key:"notLoadCacheAtStart,optional"`
	UpdateCacheWhenEmpty bool `key:"updateCacheWhenEmpty,optional"`
}

// parseURL with parameters
func parseURL(rawURL url.URL) (target, error) {
	if rawURL.Scheme != schemeName ||
		len(rawURL.Host) == 0 || len(strings.TrimLeft(rawURL.Path, "/")) == 0 {
		return target{},
			errors.Errorf("Malformed URL('%s'). Must be in the next format: 'nacos://[accessKey:sectetKey]@host/service?param=value'", rawURL.String())
	}

	var tgt target
	fmt.Print("===>rawURL ", rawURL)

	fmt.Print("rawURL<=====")

	fmt.Print("===>rawURL.Query()", rawURL.Query())

	fmt.Print("rawURL.Query() <====")

	// rawURL.Query()值为： map[groupName:[dev] namespaceid:[862309f6-ed95-49a4-87f1-d7fafdc03bae] timeout:[13000ms]]

	params := make(map[string]interface{}, len(rawURL.Query()))
	var groupName string
	for name, value := range rawURL.Query() {
		params[name] = value[0]
		fmt.Print("name <====", name)
		fmt.Print("value <====", value[0])
		if name == "groupname"  {

		    groupName = value[0]
		    fmt.Print("获取groupName value",groupName)

		}
	}

	err := mapping.UnmarshalKey(params, &tgt)

    tgt.GroupName =  groupName

	if err != nil {
		return target{}, errors.Wrap(err, "Malformed URL parameters")
	}


	tgt.LogLevel = os.Getenv("NACOS_LOG_LEVEL")
	tgt.LogDir = os.Getenv("NACOS_LOG_DIR")
	tgt.CacheDir = os.Getenv("NACOS_CACHE_DIR")

	tgt.User = rawURL.User.Username()
	tgt.Password, _ = rawURL.User.Password()
	tgt.Addr = rawURL.Host
	tgt.Service = strings.TrimLeft(rawURL.Path, "/")

	if logLevel, exists := os.LookupEnv("NACOS_LOG_LEVEL"); exists {
		tgt.LogLevel = logLevel
	}

	if logDir, exists := os.LookupEnv("NACOS_LOG_DIR"); exists {
		tgt.LogDir = logDir
	}

	if notLoadCacheAtStart, exists := os.LookupEnv("NACOS_NOT_LOAD_CACHE_AT_START"); exists {
		tgt.NotLoadCacheAtStart = notLoadCacheAtStart == "true"
	}

	if updateCacheWhenEmpty, exists := os.LookupEnv("NACOS_UPDATE_CACHE_WHEN_EMPTY"); exists {
		tgt.UpdateCacheWhenEmpty = updateCacheWhenEmpty == "true"
	}
    fmt.Print("===>target:", tgt)
    fmt.Print("target<======")
	return tgt, nil
}

