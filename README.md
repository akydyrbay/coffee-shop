# coffee-shop

A coffee shop management system built with Go and PostgreSQL, featuring advanced database operations and RESTful APIs.

## Features

- **Order Management**: Complete CRUD operations for orders
- **Menu Management**: Menu item catalog with pricing
- **Inventory Control**: Stock tracking and management
- **Advanced Reporting**: Sales analytics and search functionality
- **Containerized**: Easy deployment with Docker

## Tech Stack

- **Backend**: Go
- **Database**: PostgreSQL
- **Containerization**: Docker & Docker Compose

## Quick Start

```bash
docker compose up
```

### Core Database Tables

1. **orders** - Main order information with status tracking
2. **order_items** - Individual items within each order
3. **menu_items** - Product catalog with pricing
4. **menu_item_ingredients** - Recipe and ingredient relationships
5. **inventory** - Stock level management
6. **order_status_history** - Complete order lifecycle tracking
7. **price_history** - Historical pricing data
8. **inventory_transactions** - Stock movement records

### Database Connection
- **Host**: db
- **Port**: 5432
- **User**: latte
- **Password**: latte
- **Database**: frappuccino

## API Endpoints

### Order Management
- `POST /orders` - Create new order
- `GET /orders` - Retrieve all orders
- `GET /orders/{id}` - Get specific order
- `PUT /orders/{id}` - Update order
- `DELETE /orders/{id}` - Delete order
- `POST /orders/{id}/close` - Close order
- `POST /orders/batch-process` - Bulk order processing

### Menu Management
- `POST /menu` - Add menu item
- `GET /menu` - Get all menu items
- `GET /menu/{id}` - Get specific menu item
- `PUT /menu/{id}` - Update menu item
- `DELETE /menu/{id}` - Delete menu item

### Inventory Management
- `POST /inventory` - Add inventory item
- `GET /inventory` - Get all inventory items
- `GET /inventory/{id}` - Get specific inventory item
- `PUT /inventory/{id}` - Update inventory item
- `DELETE /inventory/{id}` - Delete inventory item
- `GET /inventory/getLeftOvers` - Get inventory leftovers with pagination

### Reporting & Analytics
- `GET /reports/total-sales` - Total sales amount
- `GET /reports/popular-items` - Popular menu items
- `GET /reports/search` - Full text search across orders and menu
- `GET /orders/numberOfOrderedItems` - Ordered items by time period
- `GET /reports/orderedItemsByPeriod` - Order statistics by day/month
