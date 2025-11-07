// bank-statement-viewer/main.go
package main

import (
	"log"
	"net/http"

	"github.com/novanm/bank-viewer/backend/domain"
	httpHandler "github.com/novanm/bank-viewer/backend/handler/http"
	"github.com/novanm/bank-viewer/backend/repository/memory"
	"github.com/novanm/bank-viewer/backend/service"
)

func main() {
	var repo domain.TransactionRepository = memory.NewMemoryRepository()

	var txService domain.TransactionService = service.NewTransactionService(repo)

	handler := httpHandler.NewTransactionHandler(txService)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	corsHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			h.ServeHTTP(w, r)
		})
	}

	port := ":9090"
	log.Printf("Starting backend server on http://localhost%s\n", port)

	if err := http.ListenAndServe(port, corsHandler(mux)); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
