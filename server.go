package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Order struct {
	ID          int    `json:"id"`
	CardHolder  string `json:"card_holder"`
	CardBrand   string `json:"card_brand"`
	AddressLine string `json:"address_line"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

func OrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query(`SELECT id, card_holder, card_brand, address_line, status, created_at FROM orders ORDER BY id DESC`)
	if err != nil {
		log.Printf("Erro ao consultar pedidos: %v", err)
		http.Error(w, "Erro ao buscar pedidos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.CardHolder, &o.CardBrand, &o.AddressLine, &o.Status, &o.CreatedAt); err != nil {
			log.Printf("Erro ao ler linha: %v", err)
			continue
		}
		orders = append(orders, o)
	}

	json.NewEncoder(w).Encode(orders)
}

func main() {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		log.Fatal("❌ DATABASE_URL não configurada.")
	}

	var err error
	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Erro ao abrir conexão com DB: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Erro ao conectar com DB: %v", err)
	}

	log.Println("✅ Conectado ao banco de dados com sucesso!")

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/api/orders", OrdersHandler)

	port := ":8080"
	log.Printf("🌍 Servidor rodando em http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
