package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	_ "modernc.org/sqlite" // or "github.com/mattn/go-sqlite3"
)

// Command-line flags
var (
	jsonOutput     = flag.Bool("json", false, "Generate JSON report files")
	customerFilter = flag.String("customer", "", "Filter report by customer name (comma-separated)")
	productFilter  = flag.String("product", "", "Filter report by product name (comma-separated)")
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	// Connect to SQLite database
	db, err := sql.Open("sqlite", "sales.db")
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer db.Close()

	// Create tables and populate
	createTables(db)
	populateCustomers(db, 10)
	populateProducts(db, 15)
	populateOrders(db, 50)

	fmt.Println("Database populated with random data.")

	// Parse filter lists
	customers := parseCSV(*customerFilter)
	products := parseCSV(*productFilter)

	// Generate filtered reports
	reportSalesPerCustomer(db, *jsonOutput, customers)
	reportPopularProducts(db, *jsonOutput, products)
	reportOrdersPerMonth(db, *jsonOutput, customers, products)
}

// -------------------- Helper --------------------

func parseCSV(input string) []string {
	if input == "" {
		return nil
	}
	parts := strings.Split(input, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func saveJSON(filename string, data interface{}) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Failed to create JSON file:", err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		log.Fatal("Failed to write JSON:", err)
	}

	fmt.Println("JSON report saved to", filename)
}

// -------------------- Database Setup --------------------

func createTables(db *sql.DB) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS customers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT
	);
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		price REAL
	);
	CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		customer_id INTEGER,
		product_id INTEGER,
		quantity INTEGER,
		order_date TEXT,
		FOREIGN KEY(customer_id) REFERENCES customers(id),
		FOREIGN KEY(product_id) REFERENCES products(id)
	);`)
	if err != nil {
		log.Fatal(err)
	}
}

// -------------------- Populate Tables --------------------

func populateCustomers(db *sql.DB, count int) {
	for i := 1; i <= count; i++ {
		name := fmt.Sprintf("Customer%d", i)
		email := fmt.Sprintf("customer%d@example.com", i)
		_, err := db.Exec(`INSERT INTO customers (name, email) VALUES (?, ?)`, name, email)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func populateProducts(db *sql.DB, count int) {
	for i := 1; i <= count; i++ {
		name := fmt.Sprintf("Product%d", i)
		price := rand.Float64()*100 + 1
		_, err := db.Exec(`INSERT INTO products (name, price) VALUES (?, ?)`, name, price)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func populateOrders(db *sql.DB, count int) {
	for i := 0; i < count; i++ {
		customerID := rand.Intn(10) + 1
		productID := rand.Intn(15) + 1
		quantity := rand.Intn(5) + 1
		orderDate := randomDate(2025, 2025)
		_, err := db.Exec(`INSERT INTO orders (customer_id, product_id, quantity, order_date) VALUES (?, ?, ?, ?)`,
			customerID, productID, quantity, orderDate)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func randomDate(startYear, endYear int) string {
	start := time.Date(startYear, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	end := time.Date(endYear, 12, 31, 23, 59, 59, 0, time.UTC).Unix()
	randTime := rand.Int63n(end-start) + start
	return time.Unix(randTime, 0).Format("2006-01-02")
}

// -------------------- Reporting Functions --------------------

func reportSalesPerCustomer(db *sql.DB, toJSON bool, customers []string) {
	type CustomerSales struct {
		Name  string  `json:"name"`
		Total float64 `json:"total_sales"`
	}

	query := `
	SELECT c.name, SUM(p.price * o.quantity) AS total_sales
	FROM orders o
	JOIN customers c ON o.customer_id = c.id
	JOIN products p ON o.product_id = p.id
	`
	var args []interface{}
	if len(customers) > 0 {
		query += "WHERE c.name IN (" + placeholders(len(customers)) + ") "
		for _, c := range customers {
			args = append(args, c)
		}
	}
	query += "GROUP BY c.name ORDER BY total_sales DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var report []CustomerSales
	fmt.Println("\n--- Total Sales per Customer ---")
	for rows.Next() {
		var cs CustomerSales
		if err := rows.Scan(&cs.Name, &cs.Total); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-12s | $%.2f\n", cs.Name, cs.Total)
		report = append(report, cs)
	}

	if toJSON {
		saveJSON("/tmp/sales_per_customer.json", report)
	}
}

func reportPopularProducts(db *sql.DB, toJSON bool, products []string) {
	type ProductReport struct {
		Name  string `json:"name"`
		Units int    `json:"units_sold"`
	}

	query := `
	SELECT p.name, SUM(o.quantity) AS total_ordered
	FROM orders o
	JOIN products p ON o.product_id = p.id
	`
	var args []interface{}
	if len(products) > 0 {
		query += "WHERE p.name IN (" + placeholders(len(products)) + ") "
		for _, p := range products {
			args = append(args, p)
		}
	}
	query += "GROUP BY p.name ORDER BY total_ordered DESC LIMIT 5"

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var report []ProductReport
	fmt.Println("\n--- Top 5 Most Popular Products ---")
	for rows.Next() {
		var pr ProductReport
		if err := rows.Scan(&pr.Name, &pr.Units); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-12s | %d units\n", pr.Name, pr.Units)
		report = append(report, pr)
	}

	if toJSON {
		saveJSON("/tmp/popular_products.json", report)
	}
}

func reportOrdersPerMonth(db *sql.DB, toJSON bool, customers, products []string) {
	type MonthReport struct {
		Month string `json:"month"`
		Count int    `json:"orders_count"`
	}

	query := `
	SELECT strftime('%Y-%m', order_date) AS month, COUNT(*) AS orders_count
	FROM orders o
	JOIN customers c ON o.customer_id = c.id
	JOIN products p ON o.product_id = p.id
	`
	var args []interface{}
	conditions := []string{}
	if len(customers) > 0 {
		conditions = append(conditions, "c.name IN ("+placeholders(len(customers))+")")
		for _, c := range customers {
			args = append(args, c)
		}
	}
	if len(products) > 0 {
		conditions = append(conditions, "p.name IN ("+placeholders(len(products))+")")
		for _, p := range products {
			args = append(args, p)
		}
	}
	if len(conditions) > 0 {
		query += "WHERE " + strings.Join(conditions, " AND ") + " "
	}
	query += "GROUP BY month ORDER BY month"

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var report []MonthReport
	fmt.Println("\n--- Orders per Month ---")
	for rows.Next() {
		var mr MonthReport
		if err := rows.Scan(&mr.Month, &mr.Count); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s | %d orders\n", mr.Month, mr.Count)
		report = append(report, mr)
	}

	if toJSON {
		saveJSON("/tmp/orders_per_month.json", report)
	}
}

// Generate placeholders for SQL IN clause
func placeholders(n int) string {
	s := make([]string, n)
	for i := range s {
		s[i] = "?"
	}
	return strings.Join(s, ",")
}
