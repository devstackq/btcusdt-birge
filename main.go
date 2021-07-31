package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type BidsBtc struct {
	ID          int     `json:"idbid"`
	PriceBid    float64 `json:"pricebid"`
	QuantityBid float64 `json:"qauntitybid"`
}
type AsksBtc struct {
	ID          int     `json:"idask"`
	PriceBid    float64 `json:"priceask"`
	QuantityBid float64 `json:"qauntityask"`
}

//append each time, then marhsal -> send client
type JsonResponse struct {
	Bids []BidsBtc `json:"bids"`
	Asks []AsksBtc `json:"asks"`
}

//find average, send cleint
//or now - send to client ? - that sort - calc avg - show ?
func parseJsonData(data []byte) {
	//create new struct ? || prev data nil ?
	jr := JsonResponse{Bids: make([]BidsBtc, 5), Asks: make([]AsksBtc, 5)}
	// for
	json.Marshal(data, &jr.Bids)
	// fmt.Printf("%+v\n", string(data))
}

//json struct - refactor
// {"lastUpdateId":12690389951,"bids":[["41634.03000000","0.46010100"],["41632.20000000","0.03058500"],["41632.18000000","0.08230900"],["41632.00000000","0.00240200"],["41631.01000000","0.04946700"]],"asks":[["41634.04000000","0.80552300"],["41634.73000000","0.17000000"],["41636.31000000","0.59072500"],["41636.47000000","0.26180500"],["41636.48000000","0.14400000"]]}

//ticker - 1 sec, get new data from trade - aggTrade -> put data in ResponseData
//all trader get data,
func wsGetBtcUdts() {
	//@depth20@100ms get count 20 bids and 20 asks - bidth each 100ms
	//@aggTrade
	c, _, err := websocket.DefaultDialer.Dial("wss://stream.binance.com:9443/ws/btcusdt@depth5@1000ms", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})
	//gorutine - new data from connection
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("receive: %s", message) //get data from binance, every sec
			// parseJsonData(message)
		}
	}()

	//1 handshake; - get each sec from server bnance data - by chan Done
	ticker := time.NewTicker(time.Second) // get data every 1 sec
	defer ticker.Stop()
	for {
		select {
		case <-done:
			log.Println("done programm")
			return
		case t := <-ticker.C:
			fmt.Printf("new tick: " + t.String())
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}

func main() {
	wsGetBtcUdts()
}
