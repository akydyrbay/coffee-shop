package dal

import "frappuccino/models"

func SortItemID(orders []models.Order) []string {
	// m := make(map[string]int, 0)
	// for _, order := range orders {
	// 	if order.Status == "closed" {
	// 		for _, item := range order.Items {
	// 			m[item.ID] += item.Quantity
	// 		}
	// 	}
	// }
	keys := make([]string, 0)
	// for key := range m {
	// 	keys = append(keys, key)
	// }
	// sort.SliceStable(keys, func(i, j int) bool {
	// 	return m[keys[i]] > m[keys[j]]
	// })
	return keys
}
