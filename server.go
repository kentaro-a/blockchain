package main

import (
	"encoding/json"
	"github.com/blockchain/modules"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sync"
)

var Blockchain []modules.Block
var Transactions []modules.Transaction
var mutex = &sync.Mutex{}

func main() {
	// Create Genesis Block and Push it to Chain
	gb := modules.GenerateGenesisBlock()
	Blockchain = append(Blockchain, gb)

	router := mux.NewRouter()

	router.PathPrefix("/chain/js/").Handler(http.StripPrefix("/chain/js/", http.FileServer(http.Dir("./webroot/js/"))))
	router.HandleFunc("/chain/index", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./webroot/index.html")
	}).Methods("GET")

	router.HandleFunc("/transaction/list", HandlerListTransaction).Methods("GET")
	router.HandleFunc("/transaction/add", HandlerAddTransaction).Methods("POST")

	router.HandleFunc("/chain/list", HandlerListChain).Methods("GET")
	router.HandleFunc("/chain/mining", HandlerAddChain).Methods("POST")

	svr := &http.Server{
		Addr:    ":8009",
		Handler: router,
	}
	if err := svr.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}

func HandlerListChain(w http.ResponseWriter, r *http.Request) {
	response(w, r, http.StatusOK, Blockchain)
}

func HandlerListTransaction(w http.ResponseWriter, r *http.Request) {
	response(w, r, http.StatusOK, Transactions)
}

func HandlerAddTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var transaction modules.Transaction
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&transaction); err != nil {
		response(w, r, http.StatusBadRequest, r.Body)
		log.Println(err)
		return
	}
	defer r.Body.Close()
	Transactions = append(Transactions, transaction)
	response(w, r, http.StatusCreated, map[string]string{"msg": "Successfully Paid"})
}

func HandlerAddChain(w http.ResponseWriter, r *http.Request) {
	if len(Transactions) > 0 {
		w.Header().Set("Content-Type", "application/json")

		m := map[string]string{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&m); err != nil {
			response(w, r, http.StatusBadRequest, r.Body)
			log.Println(err)
			return
		}

		defer r.Body.Close()
		mutex.Lock()
		newBlock := modules.GenerateBlock(Blockchain[len(Blockchain)-1], Transactions, m["Miner"])
		mutex.Unlock()
		if modules.ValidateBlock(newBlock, Blockchain[len(Blockchain)-1]) {
			Blockchain = append(Blockchain, newBlock)
			Transactions = []modules.Transaction{}
		}
		response(w, r, http.StatusCreated, map[string]string{"msg": "Successfully Added"})

	} else {
		response(w, r, http.StatusCreated, map[string]string{"msg": "No Transactions here"})
	}

}

func response(w http.ResponseWriter, r *http.Request, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	res, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("System Error."))
		return
	}
	w.WriteHeader(code)
	w.Write(res)
}
