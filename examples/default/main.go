package main

import (
	"fmt"
	"time"

	"github.com/gmelum/callback"
)

func main() {

	clb := callback.New(&callback.Options{
		Transport:    callback.REST,
		DeliveryMode: callback.RoundRobin,
		RetryMode:    callback.Next,

		RetryLimit:   2,
		RetryTimeout: time.Second * 5,
		RetryWindow:  time.Second * 1,

		EndPoints: []string{
			"http://127.0.0.1:18301",
			"http://127.0.0.1:18302",
			"http://127.0.0.1:18303",
		},
	})

	clb.On(func(data *callback.Data) {
		println(data.Point, string(data.Response.Data))
	})

	go func() {

		count := 0

		for {

			count++

			clb.Emit([]byte(fmt.Sprintf("%v", count)))
			time.Sleep(time.Second)
		}

	}()

	select {}

}
