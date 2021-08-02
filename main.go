package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

// type JsonResponse struct {
// 	Bids []Bid `json:"bids"`
// 	Asks []Bid `json:"asks"`
// }

// type Bid struct {
// 	PriceLevel string
// 	Quantity   string
// }
//bids - спрос(покупатель заявляет цену за 1 ед и объем - актив который хочет купить) 	Desc, 40.9, 40.5, 40.3
//asks - предложение(продавец завяляет цену за 1 ед и объем - актив которы хочет продать) asc,  40.1, 40.2, 40,5

// how to render graphic ? get prev data , andalyze, then show new graph ?
// {"lastUpdateId":12690389951,"bids":[["41634.03000000","0.46010100"],["41632.20000000","0.03058500"],["41632.18000000","0.08230900"],["41632.00000000","0.00240200"],["41631.01000000","0.04946700"]],"asks":[["41634.04000000","0.80552300"],["41634.73000000","0.17000000"],["41636.31000000","0.59072500"],["41636.47000000","0.26180500"],["41636.48000000","0.14400000"]]}
type JsonResponse struct {
	// Mu       sync.Mutex
	Sequence int64           `json:"lastupdateid"`
	Bids     [][]interface{} `json:"bids"`
	Asks     [][]interface{} `json:"asks"`
	Spread   float64         `json:"spread"`
	Type     string          `json:"type"`
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

	wsType := WsType{}

	var seqDepth = []JsonResponse{}

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
			seqDepth = append(seqDepth, <-data)
			// b, err := json.Marshal(seqDepth)
			// if err != nil {
			// 	panic(err)
			// }
			// seqDepth.Type = "data"
			conn.WriteJSON(seqDepth)
		}
		defer conn.Close()
	}
}

//ticker - 1 sec, get new data from trade - aggTrade -> put data in ResponseData
func wsGetBtcUdts() {
	//client ws
	//@depth20@100ms get count 20 bids and 20 asks - bidth each 100ms
	c, _, err := websocket.DefaultDialer.Dial("wss://stream.binance.com:9443/ws/btcusdt@depth5@1000ms", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	done := make(chan struct{})

	// go func() {
	defer close(done)
	for {
		jsonData := JsonResponse{}
		err = c.ReadJSON(&jsonData) //bid, ask
		if err != nil {
			log.Println("read:", err)
			return
		}
		//jsonData.Mu.Lock()
		//bid := make(chan float64)
		// go func() {
		minBid := jsonData.Bids[len(jsonData.Asks)-1][0].(string)
		mb, _ := strconv.ParseFloat(minBid, 64)
		maxAsk := jsonData.Asks[len(jsonData.Asks)-1][0].(string)
		ma, _ := strconv.ParseFloat(maxAsk, 64)
		jsonData.Spread = ma - mb
		// }()
		jsonData.Type = "newdata"
		// log.Println(jsonData, jsonData.Spread, "data")
		data <- jsonData
	}
	// }()
}

func main() {
	//file server
	fs := http.FileServer(http.Dir("./client/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", Index)

	http.HandleFunc("/wsbirge", handleWsClient)
	log.Println(http.ListenAndServe(":8081", nil))
}
