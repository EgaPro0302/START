package models

// будет чек
type Check struct {
	From_account_id int `json:"from_account_id"`
	To_account_id   int `json:"to_account_id"`
	Amount          int `json:"amount"`
}

// будет человек
type Users struct {
	ID    int `json:"id"`
	Money int `json:"money"`
}

type ErrorText struct {
	Error string
}
