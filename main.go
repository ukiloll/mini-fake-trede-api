package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github/ukilolll/trade/auth"
	"github/ukilolll/trade/pkg"
	"github/ukilolll/trade/service"
	"github/ukilolll/trade/test"

	"github.com/gin-gonic/gin"
)

var(
	_ = pkg.LoadEnv()
)

func main() {
	runServer()
}

func runServer(){
	r := gin.Default()
	r.GET("/home",auth.AuthMiddleware,func(ctx *gin.Context) {
		id,_ :=ctx.Get("id")
		username,_ := ctx.Get("username")
		ctx.String(http.StatusOK,fmt.Sprintf("hello %v %v",id,username))
		//ctx.String(http.StatusOK,fmt.Sprintf("%v"),ctx.Request)
	})
	r.GET("/login_page",auth.HandleMain)

	authR := r.Group("/auth")
	authR.GET("/google/login",auth.HandleGoogleLogin)
	authR.GET("/google/callback",auth.HandleGoogleCallback)

	trade := r.Group("/trade",auth.AuthMiddleware)
	trade.POST("/buy",service.CheckAsset,service.BuyAsset)
	trade.POST("/sell",service.CheckAsset,service.SellAsset)
	trade.GET("/check",service.LookAsset)

	go service.RunDashboard()

	srv := &http.Server{
		Addr: fmt.Sprintf(":%v",os.Getenv("SERVER_PORT")),
		Handler: r,
		ReadTimeout: 2*time.Second,
		WriteTimeout: 2*time.Second,
		MaxHeaderBytes: 1<<20,
	}


	if err := srv.ListenAndServe(); err != nil{
		log.Panic(err)
	}
}

func runTest(){
	test.Test0()
}