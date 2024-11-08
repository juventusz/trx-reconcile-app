package main

import (
	"log"
	"net/http"
	"trx-reconcile-app/internal/delivery/rest"
	"trx-reconcile-app/internal/repository/trxfile/csv"
	"trx-reconcile-app/internal/usecase/impl"
)

func main() {
	repo := csv.New()
	reconcileUsecase := impl.New(repo)
	handler := rest.New(reconcileUsecase)

	initServer(handler)
}

func initServer(h *rest.Handler) {
	http.HandleFunc("/reconcile/trxtobanks", h.ReconcileTrx)

	log.Println("Server running :9000")
	http.ListenAndServe(":9000", nil)
}
