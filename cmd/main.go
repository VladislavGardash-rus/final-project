package main

import (
	"context"
	"flag"
	"github.com/gardashvs/final-project/cfg"
	"github.com/gardashvs/final-project/internal/logger"
	http_server "github.com/gardashvs/final-project/internal/transport/http"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		PrintVersion("0.0.1", "01.02.2024", "")
		os.Exit(0)
	}

	err := cfg.InitConfig(configFile)
	if err != nil {
		panic(err)
	}

	err = logger.InitLogger(cfg.Config().Logger.Level)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go watchExitSignals(cancel)

	httpServer := http_server.NewServer(net.JoinHostPort(cfg.Config().HttpServer.Host, cfg.Config().HttpServer.Port), "image_previewer")
	go func() {
		err := httpServer.Start()
		if err != nil {
			logger.UseLogger().Error(err)
			cancel()
		}
	}()

	logger.UseLogger().Info("image_previewer service is running...")

	<-ctx.Done()
	shutDownServers(ctx, httpServer)

	logger.UseLogger().Info("image_previewer service was stopped")
}

func watchExitSignals(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	cancel()
}

func shutDownServers(ctx context.Context, httpServer *http_server.Server) {
	err := httpServer.Stop(ctx)
	if err != nil {
		logger.UseLogger().Error(err)
	}
}
