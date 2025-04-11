package main

import (
	"flag"
	"fmt"
	"github.com/colinrs/prompthub/pkg/response"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/colinrs/prompthub/internal/config"
	"github.com/colinrs/prompthub/internal/handler"
	"github.com/colinrs/prompthub/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/prompthub-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithCors("*"))
	defer server.Stop()
	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	httpx.SetErrorHandlerCtx(response.ErrHandle)
	httpx.SetOkHandler(response.OKHandle)
	stat.DisableLog()
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
