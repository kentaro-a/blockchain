package modules

type Transaction struct {
	Amount int
	From   string
	To     string
}

type Block struct {
	Index        int
	Timestamp    string
	Transactions []Transaction
	Hash         string
	PrevHash     string
	Nonce        string
	Miner        string
}
