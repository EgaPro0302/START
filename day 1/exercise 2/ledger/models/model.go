package models

type Health struct {
	Status string `json:"status"`
}

type Transaction struct {
	ID          int    `json:"id"`          //айди
	Amount      int    `json:"amount"`      //сумма
	Category    string `json:"category"`    //категория
	Description string `json:"description"` //описание
}

// не ста делать так что бы айди свмо увеличевалось без sql очень лень
