package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

//bids - спрос(покупатель заявляет цену за 1 ед и объем - актив который хочет купить) 	Desc, 40.9, 40.5, 40.3
//asks - предложение(продавец завяляет цену за 1 ед и объем - актив которы хочет продать) asc,  40.1, 40.2, 40,5

// how to render graphic ? get prev data , andalyze, then show new graph ?
// {"lastUpdateId":12690389951,"bids":[["41634.03000000","0.46010100"],["41632.20000000","0.03058500"],["41632.18000000","0.08230900"],["41632.00000000","0.00240200"],["41631.01000000","0.04946700"]],"asks":[["41634.04000000","0.80552300"],["41634.73000000","0.17000000"],["41636.31000000","0.59072500"],["41636.47000000","0.26180500"],["41636.48000000","0.14400000"]]}
type JsonResponse struct {
	Mu            sync.Mutex
	LastId        int64   `json:"lastupdateid"`
	Spread        float64 `json:"spread"`
	Type          string  `json:"type"`
	Time          int64   `json:"time"`
	MinBid        float64 `json:"minbid"`
	MinAsk        float64 `json:"minask"`
	MaxAsk        float64 `json:"maxask"`
	MaxBid        float64 `json:"maxbid"`
	AmountBid     float64 `json:"amountbid"`
	AmountAsk     float64 `json:"amountask"`
	MaxDiffAskBid float64 `json:"maxdiff"`
}

type WsBinaceData struct {
	Bids [][]interface{} `json:"bids"`
	Asks [][]interface{} `json:"asks"`
}

type WsType struct {
	Name string `json:"type"`
}

var data = make(chan JsonResponse)

//server ws
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024, // read/write, count network call
	WriteBufferSize: 1024,
}

//entry point client
func Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./client/index.html")
	log.Println("serve file")
}

func handleWsClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err, "ERR Conn")
		return
	}
	go wsGetBtcUdts()

	var wsType WsType
	// var seqDepth = []JsonResponse{}
	// var depth = JsonResponse{}

	for {
		_, msg, errk := conn.ReadMessage()
		json.Unmarshal(msg, &wsType)

		if code, ok := errk.(*websocket.CloseError); ok {
			//logout, close tab -> leave
			if code.Code == 1001 {
				log.Println(code.Code)
				break
			}
			if code.Code == 1006 {
				log.Println(code.Code)
				break
			}
		}

		if wsType.Name == "getWsBinanceData" {
			// seqDepth = append(seqDepth, <-data)
			// log.Println(d.MinBid, d.MaxBid, d.Spread)
			conn.WriteJSON(<-data)
		}
		defer conn.Close()
	}
}

//last item - show binance, last bids
//ticker - 1 sec, get new data from trade - aggTrade -> put data in ResponseData
func wsGetBtcUdts() {
	//@depth20@100ms get count 20 bids and 20 asks - bidth each 100ms
	c, _, err := websocket.DefaultDialer.Dial("wss://stream.binance.com:9443/ws/btcusdt@depth5@1000ms", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	jsonData := JsonResponse{}
	binData := WsBinaceData{Asks: make([][]interface{}, 5), Bids: make([][]interface{}, 5)}

	for {
		err = c.ReadJSON(&binData) //bid, ask
		if err != nil {
			log.Println("read:", err)
			return
		}
		// log.Println(binData.Bids)
		go func() {
			maxBidStr := binData.Bids[0][0].(string)
			minBidStr := binData.Bids[len(binData.Bids)-1][0].(string)
			//bids - min ?
			maxBid, _ := strconv.ParseFloat(maxBidStr, 64)
			jsonData.MaxBid = maxBid

			minBid, _ := strconv.ParseFloat(minBidStr, 64)
			jsonData.MinBid = minBid

			maxAskStr := binData.Asks[len(binData.Asks)-1][0].(string)
			maxAsk, _ := strconv.ParseFloat(maxAskStr, 64)
			jsonData.MaxAsk = maxAsk

			minAskStr := binData.Asks[0][0].(string)
			minAsk, _ := strconv.ParseFloat(minAskStr, 64)
			jsonData.MinAsk = minAsk

			//range each ask, sum / count
			amountAskAvg := 0.0
			for _, ask := range binData.Asks {
				log.Print(ask[1])
				size, _ := strconv.ParseFloat(ask[1].(string), 64)
				amountAskAvg += size
			}
			log.Println(amountAskAvg, "avg ask count")
			jsonData.AmountAsk = float64(amountAskAvg)
			jsonData.MaxDiffAskBid = maxAsk - minBid
			// jsonData.Spread = maxAsk - minBid
			jsonData.Spread = minAsk - maxBid

			jsonData.Type = "newdata"
			jsonData.Time = time.Now().Local().Unix()
			//x, y = todo here ?
			data <- jsonData
		}()
	}
}

func main() {
	//file server
	fs := http.FileServer(http.Dir("./client/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", Index)

	http.HandleFunc("/wsbirge", handleWsClient)
	log.Println(http.ListenAndServe(":8081", nil))
}
