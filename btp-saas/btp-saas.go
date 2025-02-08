package main

import (
	"flag"
	"fmt"

	"btp-saas/dao"
	"btp-saas/global"
	"btp-saas/internal/config"
	"btp-saas/internal/handler"
	"btp-saas/internal/svc"
	"btp-saas/mq"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/btp-saas.yaml", "the config file")

func main() {
	flag.Parse()
	logx.DisableStat()
	logx.Disable()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	global.Conf = c

	dao.Start(c)
	mq.Start(c.RedisConf)
	defer mq.Stop()

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
