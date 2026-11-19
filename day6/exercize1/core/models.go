package core

type Node struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type DTONode struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
