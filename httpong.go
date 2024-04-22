package httpong

import (
	"fmt"
	"net"
)

type Config struct {
	Addr string
}

type App struct {
	router Router
	config Config
}

func New() App {
	return App{
		router: newRouter(),
		config: Config{Addr: ":8080"},
	}
}

func NewWithConfig(config Config) App {
	return App{
		router: newRouter(),
		config: config,
	}
}

func (a *App) Run() error {
	ln, err := net.Listen("tcp", a.config.Addr)
	if err != nil {
		return err
	}
	fmt.Println("running in port ", a.config.Addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handleConnection(conn, a.router)
	}
}

func handleConnection(conn net.Conn, router Router) {
	defer conn.Close()

	req := ReadReq(conn)
	route := router.FindRoute(req)
	context := newContext(req, conn)
	if route == nil {
		context.res.ResponseString(conn, "failed to access path", 404)
		return
	}
	fmt.Println(route, context)
	route.Function(context)
}

func (a *App) RegisterRoute(method string, path string, fn Handler) {
	route := Route{method, path, fn}
	a.router.registerRoute(route)
}

func (a *App) GET(path string, fn Handler) {
	a.RegisterRoute("GET", path, fn)
}

func (a App) POST(path string, fn Handler) {
	a.RegisterRoute("POST", path, fn)
}

func (a App) PUT(path string, fn Handler) {
	a.RegisterRoute("PUT", path, fn)
}

func (a App) DELETE(path string, fn Handler) {
	a.RegisterRoute("DELETE", path, fn)
}

func (a App) PATCH(path string, fn Handler) {
	a.RegisterRoute("PATCH", path, fn)
}
