package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

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

//find average, send cleint
//or now - send to client ? - that sort - calc avg - show ?
func parseJsonData(data []byte) {
	// asks := []AsksBtc{}
	// json.Unmarshal(data, &bids)
	// fmt.Printf("%+v\n", string(data))
}

//5 id diff - for grapihc ?
// {"lastUpdateId":12690389951,"bids":[["41634.03000000","0.46010100"],["41632.20000000","0.03058500"],["41632.18000000","0.08230900"],["41632.00000000","0.00240200"],["41631.01000000","0.04946700"]],"asks":[["41634.04000000","0.80552300"],["41634.73000000","0.17000000"],["41636.31000000","0.59072500"],["41636.47000000","0.26180500"],["41636.48000000","0.14400000"]]}
type JsonResponse struct {
	Mu       sync.Mutex
	Sequence int64           `json:"lastupdateid"`
	Bids     [][]interface{} `json:"bids"`
	Asks     [][]interface{} `json:"asks"`
	Spread   float64         `json:"spread"`
}

//ticker - 1 sec, get new data from trade - aggTrade -> put data in ResponseData
//all trader get data,
func wsGetBtcUdts(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./client/index.html")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//@depth20@100ms get count 20 bids and 20 asks - bidth each 100ms
	//@aggTrade
	c, _, err := websocket.DefaultDialer.Dial("wss://stream.binance.com:9443/ws/btcusdt@depth5@1000ms", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	// go func() {
	defer close(done)
	for {
		// _, message, err := c.ReadMessage()
		// log.Printf("receive: %s", message) //get data from binance, every sec
		jsonData := JsonResponse{}
		err = c.ReadJSON(&jsonData) //bid, ask
		if err != nil {
			log.Println("read:", err)
			return
		}
		// jsonData.Mu.Lock()
		// ch := make(chan bool)

		//avg - > send avg
		// go func() {
		minBid := jsonData.Bids[len(jsonData.Asks)-1][0].(string)
		mb, _ := strconv.ParseFloat(minBid, 64)
		maxAsk := jsonData.Asks[len(jsonData.Asks)-1][0].(string)
		ma, _ := strconv.ParseFloat(maxAsk, 64)
		jsonData.Spread = ma - mb
		// }()
		//		bid := make(chan float64)
		log.Println(jsonData, jsonData.Spread)
		b, err := json.Marshal(jsonData)
		if err != nil {
			panic(err)
		}
		w.Write(b)
	}
	// }()
	//1 handshake; - get each sec from server bnance data - by chan Done
	// ticker := time.NewTicker(time.Second) // get data every 1 sec
	// defer ticker.Stop()
}

func main() {
	// runtime.GOMAXPROCS(4)
	// ws client config
	http.HandleFunc("/", wsGetBtcUdts)
	fs := http.FileServer(http.Dir("./client/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Println(http.ListenAndServe(":8080", nil))
}
