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
	"sync"
)

var Nodes = make([]core.Node, 0)
var ID int
var mtx sync.RWMutex

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
		var MaxBytesReader *http.MaxBytesError
		if errors.As(err, &MaxBytesReader) {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
			return
		}
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
	mtx.Lock()
	ID++
	RealNode := core.Node{
		ID:          ID,
		Title:       node.Title,
		Description: node.Description,
	}
	Nodes = append(Nodes, RealNode)
	mtx.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	if err := json.NewEncoder(w).Encode(RealNode); err != nil {
		fmt.Println("Error Encoder {Post}:", err)
		return
	}

}

func GetNode(w http.ResponseWriter, r *http.Request) {
	mtx.RLock()
	nodes := make([]core.Node, len(Nodes))
	copy(nodes, Nodes)
	mtx.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(nodes); err != nil {
		fmt.Println("Error Encoder {Get}")
		return
	}
}

func GetByID(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("id")
	RealID, err := strconv.Atoi(ID)
	if err != nil {
		http.Error(w, "ID Bad Request!", 400)
		return
	}
	if RealID <= 0 {
		http.Error(w, "Error Bad Request", 400)
		return
	}
	mtx.RLock()
	nodes := make([]core.Node, len(Nodes))
	copy(nodes, Nodes)
	var node core.Node
	found := false
	for _, v := range nodes {
		if v.ID == RealID {
			node = v
			found = true
			break
		}
	}
	mtx.RUnlock()
	if !found {
		http.Error(w, "Error Not Found", 404)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(node); err != nil {
		fmt.Println("Error Encoder {GetByID}:", err)
		return
	}

}
