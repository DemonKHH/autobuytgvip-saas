package main

import (
	"flag"
	"fmt"

	"btp-agent/dao"
	"btp-agent/global"
	"btp-agent/internal/config"
	"btp-agent/internal/handler"
	"btp-agent/internal/svc"
	"btp-agent/tg"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/btp-agent.yaml", "the config file")

func main() {
	flag.Parse()
	logx.DisableStat()
	logx.Disable()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	global.Conf = c

	dao.Start(c)
	tg.Start(c.AppConf)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
