package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"study/core"
	"study/features/service"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Taskhandler struct {
	serv *service.TaskService
	zap  *zap.Logger
}

func NewTaskHandler(serv *service.TaskService, zap *zap.Logger) *Taskhandler {
	return &Taskhandler{serv: serv, zap: zap}
}

func (h *Taskhandler) GetAll(w http.ResponseWriter, r *http.Request) {
	task, err := h.serv.GetAll(r.Context())
	if err != nil {
		h.zap.Error(fmt.Sprintf("%sError:", err))
		http.Error(w, fmt.Sprintf("%sError:", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		h.zap.Error(fmt.Sprintf("%sError encoding:", err))
	}
}
func (h *Taskhandler) GetbyID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	realdID, err := strconv.Atoi(id)
	if err != nil {
		h.zap.Error(fmt.Sprintf("%sError Id:", err))
		http.Error(w, "Error ID request", 400)
		return
	}
	task, err := h.serv.GetByID(r.Context(), realdID)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			h.zap.Info("Task Not Found")
			http.Error(w, "Not Found", 404)
			return
		}
		if errors.Is(err, core.ErrBadRequest) {
			h.zap.Info("user Bad Request")
			http.Error(w, "Bad Request", 400)
			return
		}
		h.zap.Error(fmt.Sprintf("%sError:", err))
		http.Error(w, fmt.Sprintf("%sError:", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		h.zap.Error(fmt.Sprintf("%sError encoding:", err))
	}
}

func (h *Taskhandler) Create(w http.ResponseWriter, r *http.Request) {
	var dtotask core.DTOTask
	if err := json.NewDecoder(r.Body).Decode(&dtotask); err != nil {
		h.zap.Info(fmt.Sprintf("%sBad Request:", err))
		http.Error(w, "Bad Request", 400)
		return
	}
	task := core.Tasks{
		Title:       dtotask.Title,
		Description: dtotask.Description,
	}
	realTask, err := h.serv.Create(r.Context(), task)
	if err != nil {
		if errors.Is(err, core.ErrBadRequest) {
			http.Error(w, "Bad Request", 400)
			return
		}
		h.zap.Error(fmt.Sprintf("%sError", err))
		http.Error(w, "Internal Server Error", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	if err := json.NewEncoder(w).Encode(realTask); err != nil {
		h.zap.Error(fmt.Sprintf("%sError", err))
		return
	}
}

func (h *Taskhandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	realID, err := strconv.Atoi(id)
	if err != nil {
		h.zap.Error(fmt.Sprintf("%sError Id:", err))
		http.Error(w, "Error ID request", 400)
		return
	}
	if err := h.serv.Delete(r.Context(), realID); err != nil {
		if errors.Is(err, core.ErrBadRequest) {
			h.zap.Info(fmt.Sprintf("%sBad Request:", err))
			http.Error(w, "Bad Request", 400)
			return
		}
		if errors.Is(err, core.ErrNotFound) {
			h.zap.Info("Task Not Found")
			http.Error(w, "Not Found", 404)
			return
		}
		h.zap.Error(fmt.Sprintf("%sError", err))
		http.Error(w, fmt.Sprintf("%sError", err), 500)
		return
	}
	w.WriteHeader(204)
}

func (h *Taskhandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	realID, err := strconv.Atoi(id)
	if err != nil {
		h.zap.Error(fmt.Sprintf("%sError Id:", err))
		http.Error(w, "Error ID request", 400)
		return
	}
	completed := true
	task, err := h.serv.UpdateStatus(r.Context(), completed, realID)
	if err != nil {

		if errors.Is(err, core.ErrBadRequest) {
			h.zap.Info(fmt.Sprintf("%sBad Request:", err))
			http.Error(w, "Bad Request", 400)
			return
		}
		if errors.Is(err, core.ErrNotFound) {
			h.zap.Info("Task Not Found")
			http.Error(w, "Not Found", 404)
			return
		}
		h.zap.Error(fmt.Sprintf("%sError", err))
		http.Error(w, fmt.Sprintf("%sError", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		h.zap.Error(fmt.Sprintf("%sError encoding:", err))
	}
}

func (h *Taskhandler) UpdateContent(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	realID, err := strconv.Atoi(id)
	if err != nil {
		h.zap.Error(fmt.Sprintf("%sError Id:", err))
		http.Error(w, "Error ID request", 400)
		return
	}
	var dtotask core.DTOTask
	if err := json.NewDecoder(r.Body).Decode(&dtotask); err != nil {
		h.zap.Info(fmt.Sprintf("%sBad Request:", err))
		http.Error(w, "Bad Request", 400)
		return
	}
	task, err := h.serv.UpdateContent(r.Context(), realID, dtotask.Title, dtotask.Description)
	if err != nil {

		if errors.Is(err, core.ErrBadRequest) {
			h.zap.Info(fmt.Sprintf("%sBad Request:", err))
			http.Error(w, "Bad Request", 400)
			return
		}
		if errors.Is(err, core.ErrNotFound) {
			h.zap.Info("Task Not Found")
			http.Error(w, "Not Found", 404)
			return
		}
		h.zap.Error(fmt.Sprintf("%sError", err))
		http.Error(w, fmt.Sprintf("%sError", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		h.zap.Error(fmt.Sprintf("%sError encoding:", err))
	}
}
