# coffee-shop

# Project Overview
This project focuses on data storage solution on PostgreSQL relational database for managing business operations. The goal is to refactor existing handlers and the data access layer (repositories) to interact with PostgreSQL using SQL queries, improving scalability and maintainability. In addition to the refactor, features for aggregation and reporting will be implemented to leverage PostgreSQL's powerful querying capabilities.

By moving to PostgreSQL, the application can benefit from the robustness of relational databases, allowing for more complex operations, easier data consistency, and improved performance. It helps to gain hands-on experience with SQL, PostgreSQL, and the principles of database design, including table relationships and advanced data types.

# Usage
## to check:
```bash
docker ps
```
## launch psql:
```bash
docker-compose up --build
```
``` bash
docker exec -it frappuccino_db_1 bash
psql -U latte -d frappuccino
```
## to fix issues and ensure the table are created
```bash
docker-compose down -v
```
```bash
chmod 644 init.sql
```
```bash
docker-compose up --build
or
docker-compose up
```