package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"study/core"
)

var Nodes []core.Node
var ID int

func PostNode(w http.ResponseWriter, r *http.Request) {
	var node core.DTONode
	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}
	// учтем что бы в нашем блокноте было все обязательно, две строки!
	if node.Title == "" && node.Description == "" {
		http.Error(w, "Bad Request", 400)
		return
	}
	ID++
	RealNode := core.Node{
		ID:          ID,
		Title:       node.Title,
		Description: node.Description,
	}
	Nodes = append(Nodes, RealNode)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	if err := json.NewEncoder(w).Encode(RealNode); err != nil {
		fmt.Println("Error Encoder {Post}:", err)
		return
	}

}

func GetNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(Nodes); err != nil {
		fmt.Println("Error Encoder {Get}")
		return
	}
}

func GetByID(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("id")
	var node core.Node
	RealID, err := strconv.Atoi(ID)
	if err != nil {
		http.Error(w, "ID Bad Request!", 400)
		return
	}
	for _, v := range Nodes {
		if v.ID == RealID {
			node = v
		}
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(node); err != nil {
		fmt.Println("Error Encoder {GetByID}:", err)
		return
	}

}
