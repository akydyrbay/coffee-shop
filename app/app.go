package app

import (
	"flag"
	"fmt"
	"frappuccino/internal/dal"
	"frappuccino/internal/handler"
	"frappuccino/internal/service"
	"frappuccino/postgres"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"

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
	// dir := flag.String("dir", "data", "The directory to serve")
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

	customerRepo := dal.NewCustomerRepo(db)
	customerService := service.NewServiceCustomer(customerRepo)
	customerHandler := handler.NewCustomerHandler(customerService)

	menuRepo := dal.NewMenuRepo(db)
	menuService := service.NewMenuService(menuRepo)
	menuHandler := handler.NewMenuHandler(menuService)

	orderRepo := dal.NewOrderRepo(db)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderService)
	basicService := service.NewService(orderRepo, inventoryRepo, menuRepo, customerRepo)
	basicHandler := handler.NewHandler(basicService)

	aggrRepo := dal.NewAggrRepo(db)
	aggService := service.NewAggragationService(aggrRepo, orderRepo, menuRepo)
	aggHandler := handler.NewAggragationHandler(aggService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /customer", customerHandler.PostItem)
	mux.HandleFunc("GET /customer", customerHandler.GetAllItem)
	mux.HandleFunc("GET /customer/{id}", customerHandler.GetItemById)
	mux.HandleFunc("DELETE /customer/{id}", customerHandler.DeleteItem)

	mux.HandleFunc("POST /orders", orderHandler.PostOrder)
	mux.HandleFunc("GET /orders", orderHandler.GetAllOrders)
	mux.HandleFunc("GET /orders/{id}", orderHandler.GetOrderByID)
	mux.HandleFunc("PUT /orders/{id}", orderHandler.PutOrderByID)
	mux.HandleFunc("DELETE /orders/{id}", orderHandler.DeleteOrderByID)
	mux.HandleFunc("POST /orders/{id}/close", basicHandler.Close)
	mux.HandleFunc("POST /orders/batch-process", basicHandler.BatchProcessOrders)

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

	mux.HandleFunc("GET /reports/total-sales", aggHandler.GetAllSales)
	mux.HandleFunc("GET /reports/popular-items", aggHandler.GetPopularSales)
	mux.HandleFunc("GET /orders/numberOfOrderedItems", aggHandler.GetOrderByDate)
	mux.HandleFunc("GET /reports/search", aggHandler.FullTextSearchReport)
	mux.HandleFunc("GET /reports/orderedItemsByPeriod", aggHandler.OrderedItemsByPeriod)
	mux.HandleFunc("GET /inventory/getLeftOvers", aggHandler.GetLeftOvers)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), mux))
}

func printHelpUsage() {
	fmt.Println("./hot-coffee --help\nCoffee Shop Management System\n\nUsage:\n  hot-coffee [--port <N>] [--dir <S>] \n  hot-coffee --help\n\nOptions:\n  --help       Show this screen.\n  --port N     Port number.\n  --dir S      Path to the data directory.")
}
