package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Reshnyak/innopolis/ratelimiter/limiter"
	mbeego "github.com/Reshnyak/innopolis/ratelimiter/middleware/beego"
	mgin "github.com/Reshnyak/innopolis/ratelimiter/middleware/gin"
	"github.com/Reshnyak/innopolis/ratelimiter/middleware/stdhttp"
	beepack "github.com/astaxie/beego"
	ginpack "github.com/gin-gonic/gin"
)

const (
	stdLim  = 5
	stdAddr = "localhost:7777"
	ginLim  = 4
	ginAddr = "localhost:8888"
	beeLim  = 3
	beeAddr = "localhost:9999"
)

type MainController struct {
	beepack.Controller
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	//standart http server
	go func() {
		stdMiddleware := stdhttp.NewStdHttpMiddleware(limiter.New(stdLim))
		http.Handle("/", stdMiddleware.Handler(http.HandlerFunc(HandleTask)))
		fmt.Printf("Standart server is running on %s...", stdAddr)
		log.Fatal(http.ListenAndServe(stdAddr, nil))
	}()
	// gin http server
	go func() {
		ginpack.SetMode(ginpack.DebugMode) // ReleaseMode
		ginMiddleware := mgin.NewGinMiddleware(limiter.New(ginLim))
		router := ginpack.Default()
		router.ForwardedByClientIP = true
		router.Use(ginMiddleware)
		router.GET("/", HandleGinTask)
		log.Fatal(router.Run(ginAddr))
	}()

	//beego http server
	go func() {
		beeMiddleware := mbeego.NewBeeMiddleware((limiter.New(beeLim)))
		beepack.Router("/", &MainController{})
		beepack.RunWithMiddleWares(beeAddr, beeMiddleware.Handler)
	}()
	<-ctx.Done()
}

func HandleGinTask(context *ginpack.Context) {
	// Simulate a resource-intensive task
	time.Sleep(2 * time.Second)
	type message struct {
		Message string `json:"message"`
	}
	resp := message{Message: fmt.Sprintf("Task from %s Completed", context.RemoteIP())}
	context.JSON(http.StatusOK, resp)
}
func HandleTask(writer http.ResponseWriter, request *http.Request) {
	// Simulate a resource-intensive task
	time.Sleep(2 * time.Second)
	_, _ = fmt.Fprintf(writer, "Task from %s Completed", request.RemoteAddr)
}
func (m *MainController) Get() {
	// Simulate a resource-intensive task
	time.Sleep(2 * time.Second)
	m.Ctx.WriteString(fmt.Sprintf("Task from %s Completed", m.Ctx.Request.RemoteAddr))
}
