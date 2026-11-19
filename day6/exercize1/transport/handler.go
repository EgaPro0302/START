package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"study/core"
)

var Nodes []core.Node
var ID int

func PostNode(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()

	var node core.DTONode
	if err := decode.Decode(&node); err != nil {
		var MaxBytesReader *http.MaxBytesError
		if errors.As(err, &MaxBytesReader) {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "Error Json", 400)
		return
	}
	if err := decode.Decode(&struct{}{}); err != io.EOF {
		http.Error(w, "body must contain only one JSON  object", 400)
		return
	}
	node.Title = strings.TrimSpace(node.Title)
	node.Description = strings.TrimSpace(node.Description)
	// учтем что бы в нашем блокноте было все обязательно, две строки!
	if node.Title == "" || node.Description == "" {
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
	if RealID <= 0 {
		http.Error(w, "Error Bad Request", 400)
		return
	}
	var f bool = false
	for _, v := range Nodes {
		if v.ID == RealID {
			node = v
			f = true
		}
	}
	if !f {
		http.Error(w, "Error Not Found", 404)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(node); err != nil {
		fmt.Println("Error Encoder {GetByID}:", err)
		return
	}

}
