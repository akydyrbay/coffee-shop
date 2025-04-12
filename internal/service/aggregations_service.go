package service

import (
	"errors"
	"frappuccino/internal/dal"
	"frappuccino/models"
	"sort"
	"strconv"
	"strings"
)

type AggragationService interface {
	GetTotalSales() (float64, models.Status)
	GetPopularMenuItems() ([]models.OrderItem, error)
	GetOrderQuantityByDate(startDate, endDate string) ([]models.NumberOfOrderItems, error)
	FullSearchReport(q, filter, minPrice, maxPrice string) (models.FullSearchStruct, error)
	GetOrderedItemsByDay(monthstr string) (models.OrderByDayRequest, error)
	GetOrderedItemsByMonth(yeartstr string) (models.OrderByMonthRequest, error)
	GetLeftOvers(sortBy string, page, pageSize int) (models.LeftoversResponse, error)
}

type aggragationService struct {
	orderRepo dal.OrderRepository
	menuRepo  dal.MenuRepository
	aggrRepo  dal.AggregationRepository
}

func NewAggragationService(aggrRepo dal.AggregationRepository, orderRepo dal.OrderRepository, menuRepo dal.MenuRepository) *aggragationService {
	return &aggragationService{aggrRepo: aggrRepo, menuRepo: menuRepo, orderRepo: orderRepo}
}

func (s *aggragationService) GetTotalSales() (float64, models.Status) {
	allMenuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return 0, models.Status{ErrorMessage: err, Code: 500}
	}
	menuMap := make(map[int]models.MenuItem)
	for _, menuItem := range allMenuItems {
		menuMap[menuItem.ID] = menuItem
	}
	allOrderItems, err := s.orderRepo.GetAll()

	var totalSales float64
	for _, orderItem := range allOrderItems {
		if orderItem.Status == "active" {
			continue
		}
		for _, ingr := range orderItem.Items {
			itemMenu := menuMap[ingr.OrderID]
			if itemMenu.ID == 0 {
				return 0, models.NotFound
			} else {
				items := orderItem.Items
				for _, item := range items {
					totalSales += float64(item.Quantity) * itemMenu.Price
				}
			}
		}
	}
	return totalSales, models.Success
}

func (s *aggragationService) GetPopularMenuItems() ([]models.OrderItem, error) {
	orderItems, err := s.orderRepo.GetAll()
	if err != nil {
		return nil, err
	}

	itemCount := make(map[int]int)

	for _, order := range orderItems {
		if order.Status == "closed" {
			for _, item := range order.Items {
				itemCount[item.OrderID] += int(item.Quantity)
			}
		}
	}

	var popularItems []models.OrderItem
	for itemID, quantity := range itemCount {
		popularItems = append(popularItems, models.OrderItem{
			OrderID:  itemID,
			Quantity: float64(quantity),
		})
	}

	sort.Slice(popularItems, func(i, j int) bool {
		return popularItems[i].Quantity > popularItems[j].Quantity
	})

	return popularItems, nil
}

func (s *aggragationService) GetOrderQuantityByDate(startDate, endDate string) ([]models.NumberOfOrderItems, error) {
	return s.aggrRepo.GetOrderQuantitybyDate(startDate, endDate)
}

func (s *aggragationService) FullSearchReport(q, filter, minPriceStr, maxPriceStr string) (models.FullSearchStruct, error) {
	var result models.FullSearchStruct
	var menuExist, orderExist bool
	var minPrice, maxPrice float64
	if filter != "" {
		filters := strings.Split(filter, ",")
		for i := 0; i < len(filters); i++ {
			switch filters[i] {
			case "menu":
				menuExist = true
			case "orders":
				orderExist = true
			default:
				return models.FullSearchStruct{}, errors.New("filter value is incorrect")
			}
		}
	}
	if !menuExist && !orderExist {
		menuExist, orderExist = true, true
	}
	var err error
	if minPriceStr != "" {
		minPrice, err = strconv.ParseFloat(minPriceStr, 256)
		if err != nil {
			return result, err
		}
	} else {
		minPrice = 0
	}
	if maxPriceStr != "" {
		maxPrice, err = strconv.ParseFloat(maxPriceStr, 256)
		if err != nil {
			return result, err
		}
	} else {
		maxPrice = 0
	}
	if orderExist {
		orders, err := s.aggrRepo.FullSearchOrder(q, minPrice, maxPrice)
		if err != nil {
			return result, err
		}
		result.Orders = orders
	}
	if menuExist {
		menu, err := s.aggrRepo.FullSearchMenu(q, minPrice, maxPrice)
		if err != nil {
			return result, err
		}
		result.MenuItems = menu
	}
	result.TotalMatches = len(result.MenuItems) + len(result.Orders)
	if result.TotalMatches == 0 {
		return models.FullSearchStruct{}, errors.New("total matches is 0")
	}
	return result, nil
}

func (s *aggragationService) GetOrderedItemsByDay(monthstr string) (models.OrderByDayRequest, error) {
	var month int
	var err error
	var orderByDay models.OrderByDayRequest
	monthMap := map[string]int{
		"january":   1,
		"february":  2,
		"march":     3,
		"april":     4,
		"may":       5,
		"june":      6,
		"july":      7,
		"august":    8,
		"september": 9,
		"october":   10,
		"november":  11,
		"december":  12,
	}
	var exists bool
	month, exists = monthMap[monthstr]
	if !exists {
		return models.OrderByDayRequest{}, errors.New("month parameter is incorrect")
	}
	orderByDay.Period = "day"
	orderByDay.Month = monthstr
	orderByDay.OrderedItems = []models.OrderedItemDay{}
	err = s.aggrRepo.GetDayPeriod(month, &orderByDay)
	if err != nil {
		return models.OrderByDayRequest{}, err
	}
	return orderByDay, nil
}

func (s *aggragationService) GetOrderedItemsByMonth(yearstr string) (models.OrderByMonthRequest, error) {
	var year int
	var err error
	var orderByMonth models.OrderByMonthRequest

	year, err = strconv.Atoi(yearstr)
	if err != nil {
		return models.OrderByMonthRequest{}, err
	}
	if year <= 0 {
		return models.OrderByMonthRequest{}, errors.New("year parameter must be greater than zero")
	}
	orderByMonth.Period = "month"
	orderByMonth.Year = yearstr
	orderByMonth.OrderedItems = []models.OrderedItemMonth{}
	err = s.aggrRepo.GetMonthPeriod(year, &orderByMonth)
	if err != nil {
		return models.OrderByMonthRequest{}, err
	}
	return orderByMonth, nil
}

func (s *aggragationService) GetLeftOvers(sortBy string, page, pageSize int) (models.LeftoversResponse, error) {
	leftovers, totalItems, err := s.aggrRepo.GetLeftOvers(sortBy, page, pageSize)
	if err != nil {
		return models.LeftoversResponse{}, err
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	if page > totalPages {
		return models.LeftoversResponse{}, errors.New("page does not exist")
	}

	if len(leftovers) == 0 {
		return models.LeftoversResponse{}, errors.New("no inventory leftovers found")
	}

	hasNextPage := page < totalPages

	response := models.LeftoversResponse{
		CurrentPage: page,
		HasNextPage: hasNextPage,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		Data:        leftovers,
	}

	return response, nil
}
