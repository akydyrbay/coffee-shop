package server

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"coffee-shop/internal/dal"
	"coffee-shop/internal/handler"
	"coffee-shop/internal/service"

	"coffee-shop/postgres"

	_ "github.com/lib/pq"
)

func StartTheCafe() {
	db, err := postgres.CheckDB()
	if err != nil {
		slog.Error("Failed to start program", "CheckDB err:", err)
		log.Fatal(err)
	}
	defer db.Close()
	port := flag.Int("port", 8080, "The server port")
	dir := flag.String("dir", "data", "The directory to serve")
	help := flag.Bool("help", false, "Show help")
	flag.Parse()
	if *port <= 0 || *port > 65535 {
		fmt.Println("Invalid port")
		os.Exit(1)
	}
	if *help {
		printHelpUsage()
		os.Exit(0)
	}

	inventoryRepo := dal.NewInventoryRepo(db)
	inventoryService := service.NewInventoryService(inventoryRepo)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)

	menuRepo := dal.NewMenuRepo(filepath.Join(*dir, "menu_items.json"))
	menuService := service.NewMenuService(menuRepo)
	menuHandler := handler.NewMenuHandler(menuService)

	// orderRepo := dal.NewOrderRepo(filepath.Join(*dir, "orders.json"))
	// orderService := service.NewOrderService(orderRepo, menuRepo, inventoryRepo)
	// orderHandler := handler.NewOrderHandler(orderService)

	// aggService := service.NewAggragationService(orderRepo, menuRepo)
	// aggHandler := handler.NewAggragationHandler(aggService)
	mux := http.NewServeMux()

	// mux.HandleFunc("POST /orders", orderHandler.PostOrder)
	// mux.HandleFunc("GET /orders", orderHandler.GetAllOrders)
	// mux.HandleFunc("GET /orders/{id}", orderHandler.GetOrderByID)
	// mux.HandleFunc("PUT /orders/{id}", orderHandler.PutOrderByID)
	// mux.HandleFunc("DELETE /orders/{id}", orderHandler.DeleteOrderByID)
	// mux.HandleFunc("POST /orders/{id}/close", orderHandler.PostCloseOrder)

	mux.HandleFunc("POST /inventory", inventoryHandler.PostItem)
	mux.HandleFunc("GET /inventory", inventoryHandler.GetAllItem)
	mux.HandleFunc("GET /inventory/{id}", inventoryHandler.GetItemById)
	mux.HandleFunc("PUT /inventory/{id}", inventoryHandler.PutItem)
	mux.HandleFunc("DELETE /inventory/{id}", inventoryHandler.DeleteItem)

	mux.HandleFunc("POST /menu", menuHandler.PostMenuHandler)
	mux.HandleFunc("GET /menu", menuHandler.GetAllMenuHandler)
	mux.HandleFunc("GET /menu/{id}", menuHandler.GetMenuItemHandler)
	mux.HandleFunc("PUT /menu/{id}", menuHandler.PutMenuHandler)
	mux.HandleFunc("DELETE /menu/{id}", menuHandler.DeleteMenuHandler)

	// mux.HandleFunc("GET /reports/total-sales", aggHandler.GetAllSales)
	// mux.HandleFunc("GET /reports/popular-items", aggHandler.GetPopularSales)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), mux))
}

func printHelpUsage() {
	fmt.Println("./hot-coffee --help\nCoffee Shop Management System\n\nUsage:\n  hot-coffee [--port <N>] [--dir <S>] \n  hot-coffee --help\n\nOptions:\n  --help       Show this screen.\n  --port N     Port number.\n  --dir S      Path to the data directory.")
}
