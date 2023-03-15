package main

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()
	rand.Seed(time.Now().Unix())

	lbFuncOne := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("req to ServiceOne")
		time.Sleep(time.Second * 3)
		z, _ := io.ReadAll(r.Body)
		_, _ = w.Write(z)
	}
	lbFuncTwo := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("req to ServiceTwo")
		time.Sleep(time.Second * 2)
		z, _ := io.ReadAll(r.Body)
		_, _ = w.Write(z)
	}
	lbFuncThree := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("req to ServiceThree")
		time.Sleep(time.Second * 1)
		z, _ := io.ReadAll(r.Body)
		_, _ = w.Write(z)
	}

	lbFuncFour := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("req to ServiceFour")
		time.Sleep(time.Millisecond * 800)
		z, _ := io.ReadAll(r.Body)
		_, _ = w.Write(z)
	}
	lbFuncFive := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("req to ServiceFive")
		time.Sleep(time.Millisecond * 600)
		z, _ := io.ReadAll(r.Body)
		_, _ = w.Write(z)
	}
	lbFuncSix := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("req to ServiceSix")
		time.Sleep(time.Millisecond * 400)
		z, _ := io.ReadAll(r.Body)
		_, _ = w.Write(z)
	}

	serverOne := http.Server{
		Addr:    fmt.Sprintf(":%d", 9001),
		Handler: http.HandlerFunc(lbFuncOne),
	}
	serverTwo := http.Server{
		Addr:    fmt.Sprintf(":%d", 9002),
		Handler: http.HandlerFunc(lbFuncTwo),
	}
	serverThree := http.Server{
		Addr:    fmt.Sprintf(":%d", 9003),
		Handler: http.HandlerFunc(lbFuncThree),
	}

	serverFour := http.Server{
		Addr:    fmt.Sprintf(":%d", 9004),
		Handler: http.HandlerFunc(lbFuncFour),
	}
	serverFive := http.Server{
		Addr:    fmt.Sprintf(":%d", 9005),
		Handler: http.HandlerFunc(lbFuncFive),
	}
	serverSix := http.Server{
		Addr:    fmt.Sprintf(":%d", 9006),
		Handler: http.HandlerFunc(lbFuncSix),
	}
	go func() {
		err := serverOne.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		err := serverTwo.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		err := serverThree.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		err := serverFour.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		err := serverFive.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		err := serverSix.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()

	<-ctx.Done()
}
