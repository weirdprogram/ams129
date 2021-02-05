package main

import (
	"log"
	"sync"
	"os"
	"fmt"
	"golang.org/x/net/websocket"
)

func check(e error) {
    if e != nil {
        log.Fatal(e)
    }
}

func main() {
    var wg sync.WaitGroup
    
    wg.Add(1)
    go func(){
    	origin := "http://localhost/"
		ws, err := websocket.Dial("wss://stream.binance.com:9443/stream?streams=dogeeur@miniTicker/dogeeur@kline_15m/dogeeur@kline_30m/dogeeur@kline_1h", "", origin)

	    check(err)

	    defer ws.Close()

    	for {
		   	var message string
			websocket.Message.Receive(ws, &message)

			log.Println("Received: ", message)
			file, err := os.OpenFile("data/stream.binance.com.1.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend | 0644)
		    check(err)
		    defer file.Close()
		 
		    _, err = file.WriteString(fmt.Sprintln(message))
		    check(err)
			
		}
		wg.Done()
	}()
	
	wg.Wait()   
}