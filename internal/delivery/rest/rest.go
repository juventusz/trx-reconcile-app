package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
	"trx-reconcile-app/internal/usecase"
)

type Handler struct {
	ReconcileUsecase usecase.Usecase
}

func New(reconUsecase usecase.Usecase) *Handler {
	return &Handler{
		ReconcileUsecase: reconUsecase,
	}
}

func (h *Handler) ReconcileTrx(w http.ResponseWriter, r *http.Request) {
	//part 1: parse input
	systemTrxFilePath := r.FormValue("system_trx")
	banksStatsFIlePath := strings.Split(r.FormValue("banks_stats_trx"), ",")

	if systemTrxFilePath == "" || len(banksStatsFIlePath) == 0 {
		log.Println("empty filepaths, rejected")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse("2006-01-02 15:04:05", r.FormValue("start_time"))
	if err != nil {
		log.Printf("Skipping row with invalid timestamp: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", r.FormValue("end_time"))
	if err != nil {
		log.Printf("Skipping row with invalid timestamp: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// part 2: proceed handler
	resp, err := h.ReconcileUsecase.ReconcileTransactionAndBanks(startTime, endTime, systemTrxFilePath, banksStatsFIlePath...)
	if err != nil {
		log.Println("error reconcile, err:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Println("error marshal, err:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}
