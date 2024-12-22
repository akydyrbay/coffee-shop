# coffee

## Project Overview
This project focuses on data storage solution on PostgreSQL relational database for managing business operations. The goal is to refactor existing handlers and the data access layer (repositories) to interact with PostgreSQL using SQL queries, improving scalability and maintainability. In addition to the refactor, new features for aggregation and reporting will be implemented to leverage PostgreSQL's powerful querying capabilities.

By moving to PostgreSQL, the application can benefit from the robustness of relational databases, allowing for more complex operations, easier data consistency, and improved performance. You will gain hands-on experience with SQL, PostgreSQL, and the principles of database design, including table relationships and advanced data types.

## Learning Objectives
1. SQL
Develop proficiency in writing SQL queries to interact with relational databases.
Implement basic CRUD (Create, Read, Update, Delete) operations using SQL.
Use more complex SQL queries involving JOINs, subqueries, and aggregation functions.
Learn to optimize queries for better performance and scalability.
2. PostgreSQL
Understand PostgreSQL’s unique features like JSONB, Arrays, ENUMs, and Timestamps.
Learn how to design efficient and normalized database schemas.
Leverage PostgreSQL’s advanced data types for storing and querying data efficiently.
3. CRUD Operations
Refactor and implement CRUD operations for managing business data (e.g., orders, customers, products).
Ensure data integrity and optimize relational operations for scalability.
4. ERD (Entity-Relationship Diagram)
Design and document an Entity-Relationship Diagram (ERD) to model the relationships between key business entities.
Implement the database schema based on the ERD, ensuring proper relationships and constraints are defined.
Abstract
This project involves refactoring an application to transition from a JSON-based data store to a relational database using PostgreSQL. The primary task is to replace existing JSON-based data handlers with SQL queries to interact with a relational database. This transition involves designing and implementing a database schema with appropriate tables, relationships, and constraints.

In addition, new features will be implemented for reporting and aggregation, such as tracking inventory changes, calculating sales data, and monitoring status changes in business operations. PostgreSQL’s advanced SQL functions, including aggregation tools, will be leveraged to enable complex reporting and analysis.

## Context
As systems scale, databases that rely on JSON for data storage often become difficult to manage and maintain. Relational databases, like PostgreSQL, provide more robust data management, scalability, and advanced querying capabilities, making them a better solution for growing applications. This project aims to address the challenges of using a JSON-based system by transitioning to a more structured, performant, and scalable PostgreSQL database.

The project will involve creating a normalized database schema, defining the relationships between entities, and refactoring the application to use SQL queries for CRUD operations. Additionally, it will introduce new functionality for aggregating and reporting business metrics, leveraging PostgreSQL’s powerful aggregation capabilities.


### Relationships of ERD
orders is linked to customers (via customer_id).
order_items connects orders to menu_items, storing quantity and price at the time of the order.
menu_items is connected to menu_item_ingredients, which links the items to the ingredients, including quantities for recipes.
ingredients is used in menu_item_ingredients and inventory_transactions to track ingredient usage and stock levels.
orders has a one-to-many relationship with order_status_history, storing order status changes over time.
menu_items is related to price_history to track pricing over time.
staff can be connected to orders (if needed) to track which staff member handled a particular order.