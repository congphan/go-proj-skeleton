// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type SubRequest struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	Id     int      `json:"id"`
}

var addr = flag.String("addr", "localhost:11011", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	group := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < 1; i++ {
		go mainA(ctx, group)
	}
	<-interrupt
	cancel()
	group.Wait()
}

func mainA(ctx context.Context, gr *sync.WaitGroup) {
	gr.Add(1)
	defer gr.Done()

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws/btcusdt@depth"}

	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	type AuthRequest struct {
		Method string `json:"method"`
		Params struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"params"`
		Id int `json:"id"`
	}

	authReq := AuthRequest{}
	authReq.Method = "AUTH"
	authReq.Id = 1
	authReq.Params.Type = "token"

	go func() {
		defer close(done)
		for {
			var message interface{}
			err := c.ReadJSON(&message)
			if err != nil {
				log.Println("read:", err)
				return
			}

			bytes, _ := json.Marshal(message)
			fmt.Printf("received message: %s\n", string(bytes))
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	tickerList := time.NewTicker(time.Second * 3)
	defer tickerList.Stop()

	tickerUnsub := time.NewTicker(time.Second * 10)
	defer tickerUnsub.Stop()

	go func() {
		// defer wait.Done()
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			sub := SubRequest{}
			if err := json.NewDecoder(bytes.NewReader(scanner.Bytes())).Decode(&sub); err != nil {
				fmt.Printf("====== got error: %v\n", err)
			}

			if sub.Method == "AUTH" {
				// 123456
				authReq.Params.Value = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzNDU2In0.JucMlcxvcClLoFKvZaLygMvDqueUkgaW-SZ9xlrBZgo"
				if err := c.WriteJSON(authReq); err != nil {
					fmt.Printf("====== writeJson AUTH got error: %v\n", err)
				}
				continue
			}

			if sub.Method == "AUTH1" {
				//12345
				authReq.Params.Value = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzNDUifQ.DPXS1Xi7fyVG9aW2jMKTWavV-Tv4VgTvbwSf8DaJ0fg"
				if err := c.WriteJSON(authReq); err != nil {
					fmt.Printf("====== writeJson AUTH1 got error: %v\n", err)
				}
				continue
			}

			if err := c.WriteJSON(sub); err != nil {
				fmt.Printf("====== writeJson got error: %v\n", err)
			}
		}

	}()

	for {
		select {
		case <-done:
			return
		case <-tickerList.C:

		case <-ctx.Done():
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}

}
