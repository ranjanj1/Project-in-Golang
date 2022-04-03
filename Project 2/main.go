package main

import (
	"Project2/config"
	grpcPkg "Project2/grpc"
	"Project2/handler"
	"Project2/mapper"
	tokenv1 "Project2/proto"
	"Project2/repository"
	"Project2/service"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	appConfig, err := config.ProvideAppConfig()
	check(err)

	serverConfig := appConfig.ServerConfig
	v := grpcPkg.ProvideGrpcServerOptions(serverConfig)
	server := grpcPkg.ProvideServer(v...)
	listener, err := grpcPkg.ProvideListener(serverConfig)
	check(err)

	tokenRepository := repository.ProvideTokenRepository()
	hashService := service.ProvideHashService()
	tokenService := service.ProvideTokenService(tokenRepository, hashService)
	mapperMapper := mapper.ProvideMapper()
	tokenHandler := handler.ProvideTokenHandler(tokenService, mapperMapper)
	app := &App{
		Server:      server,
		Listener:    listener,
		TickHandler: tokenHandler,
	}

	app.Start(check)

	<-interrupt()
	app.Shutdown()
}

type App struct {
	Server      *grpc.Server
	Listener    net.Listener
	TickHandler *handler.TokenHandler
}

func (app *App) Start(checkErr func(err error)) {
	app.registerServers()
	go func() {
		log.Println(fmt.Sprintf("GRPC Server started at %s", app.Listener.Addr()))
		err := app.Server.Serve(app.Listener)
		checkErr(err)
	}()
}

func (app *App) registerServers() {
	reflection.Register(app.Server)
	tokenv1.RegisterTokenServiceServer(app.Server, app.TickHandler)
}

func (app *App) Shutdown() {
	app.Server.GracefulStop()
}

func interrupt() chan os.Signal {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	return interrupt
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
