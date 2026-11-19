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
)

type Taskhandler struct {
	serv *service.TaskService
}

func NewTaskHandler(serv *service.TaskService) *Taskhandler {
	return &Taskhandler{serv: serv}
}

func (h *Taskhandler) GetAll(w http.ResponseWriter, r *http.Request) {
	logger, loggerClose, err := core.NewLogs("DEBUG")
	if err != nil {
		panic(err)
	}
	defer loggerClose()
	task, err := h.serv.GetAll(r.Context())
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			logger.Info("Task Not Found")
			w.WriteHeader(404)
			if err := json.NewEncoder(w).Encode(task); err != nil {
				logger.Error(fmt.Sprintf("%sError encoding:", err))
			}
			return
		}

		logger.Error(fmt.Sprintf("%sError:", err))
		http.Error(w, fmt.Sprintf("%sError:", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		logger.Error(fmt.Sprintf("%sError encoding:", err))
	}
}
func (h *Taskhandler) GetbyID(w http.ResponseWriter, r *http.Request) {

	logger, loggerClose, err := core.NewLogs("DEBUG")
	if err != nil {
		panic(err)
	}
	id := mux.Vars(r)["id"]
	realdID, err := strconv.Atoi(id)
	if err != nil {
		logger.Error(fmt.Sprintf("%sError Id:", err))
		http.Error(w, "Error ID request", 400)
		return
	}
	defer loggerClose()
	task, err := h.serv.GetByID(r.Context(), realdID)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			logger.Info("Task Not Found")
			http.Error(w, "Not Found", 404)
			return
		}
		logger.Error(fmt.Sprintf("%sError:", err))
		http.Error(w, fmt.Sprintf("%sError:", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		logger.Error(fmt.Sprintf("%sError encoding:", err))
	}
}

func (h *Taskhandler) Create(w http.ResponseWriter, r *http.Request) {
	logger, loggerClose, err := core.NewLogs("DEBUG")
	if err != nil {
		panic(err)
	}
	defer loggerClose()

	var dtotask core.DTOTask
	if err := json.NewDecoder(r.Body).Decode(&dtotask); err != nil {
		logger.Info(fmt.Sprintf("%sBad Request:", err))
		http.Error(w, "Bad Request", 400)
		return
	}
	task := core.Tasks{
		Title:       dtotask.Title,
		Description: dtotask.Description,
	}
	realTask, err := h.serv.Create(r.Context(), task)
	if err != nil {
		logger.Error(fmt.Sprintf("%sError", err))
		http.Error(w, fmt.Sprintf("%sError", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(realTask); err != nil {
		logger.Error(fmt.Sprintf("%sError", err))
		return
	}
}

func (h *Taskhandler) Delete(w http.ResponseWriter, r *http.Request) {
	logger, loggerClose, err := core.NewLogs("DEBUG")
	if err != nil {
		panic(err)
	}
	defer loggerClose()
	id := mux.Vars(r)["id"]
	realID, err := strconv.Atoi(id)
	if err != nil {
		logger.Error(fmt.Sprintf("%sError Id:", err))
		http.Error(w, "Error ID request", 400)
		return
	}
	if err := h.serv.Delete(r.Context(), realID); err != nil {
		if errors.Is(err, core.ErrBadRequest) {
			logger.Info(fmt.Sprintf("%sBad Request:", err))
			http.Error(w, "Bad Request", 400)
			return
		}
		if errors.Is(err, core.ErrNotFound) {
			logger.Info("Task Not Found")
			http.Error(w, "Not Found", 404)
			return
		}
		logger.Error(fmt.Sprintf("%sError", err))
		http.Error(w, fmt.Sprintf("%sError", err), 500)
		return
	}
	w.WriteHeader(204)
}

func (h *Taskhandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	logger, loggerClose, err := core.NewLogs("DEBUG")
	if err != nil {
		panic(err)
	}
	defer loggerClose()
	id := mux.Vars(r)["id"]
	realID, err := strconv.Atoi(id)
	if err != nil {
		logger.Error(fmt.Sprintf("%sError Id:", err))
		http.Error(w, "Error ID request", 400)
		return
	}
	completed := true
	task, err := h.serv.UpdateStatus(r.Context(), completed, realID)
	if err != nil {

		if errors.Is(err, core.ErrBadRequest) {
			logger.Info(fmt.Sprintf("%sBad Request:", err))
			http.Error(w, "Bad Request", 400)
			return
		}
		if errors.Is(err, core.ErrNotFound) {
			logger.Info("Task Not Found")
			http.Error(w, "Not Found", 404)
			return
		}
		logger.Error(fmt.Sprintf("%sError", err))
		http.Error(w, fmt.Sprintf("%sError", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		logger.Error(fmt.Sprintf("%sError encoding:", err))
	}
}

func (h *Taskhandler) UpdateContent(w http.ResponseWriter, r *http.Request) {
	logger, loggerClose, err := core.NewLogs("DEBUG")
	if err != nil {
		panic(err)
	}
	defer loggerClose()
	id := mux.Vars(r)["id"]
	realID, err := strconv.Atoi(id)
	if err != nil {
		logger.Error(fmt.Sprintf("%sError Id:", err))
		http.Error(w, "Error ID request", 400)
		return
	}
	var dtotask core.DTOTask
	if err := json.NewDecoder(r.Body).Decode(&dtotask); err != nil {
		logger.Info(fmt.Sprintf("%sBad Request:", err))
		http.Error(w, "Bad Request", 400)
		return
	}
	task, err := h.serv.UpdateContent(r.Context(), realID, dtotask.Title, dtotask.Description)
	if err != nil {

		if errors.Is(err, core.ErrBadRequest) {
			logger.Info(fmt.Sprintf("%sBad Request:", err))
			http.Error(w, "Bad Request", 400)
			return
		}
		if errors.Is(err, core.ErrNotFound) {
			logger.Info("Task Not Found")
			http.Error(w, "Not Found", 404)
			return
		}
		logger.Error(fmt.Sprintf("%sError", err))
		http.Error(w, fmt.Sprintf("%sError", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		logger.Error(fmt.Sprintf("%sError encoding:", err))
	}
}
