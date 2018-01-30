package app

import (
	"database/sql"
	_ "github.com/lib/pq"
	"gopkg.in/gorp.v1"
	"log"
	"fmt"
	"os"
	"encoding/json"
)

//Configuration type corresponds to the DB Configs JSON from the 'dbconf.json' file
type Configuration struct {
	Host string `json:"host"`
	Port int `json:"port"`
	User string `json:"user"`
	Password string `json:"password"`
	Dbname string `json:"dbname"`
}

//Order type corresponds to the 'orders' table in the DB and defines its schema
type Order struct {
	// db tag lets you specify the column name if it differs from the struct field
	Id      int64  `db:"order_id"`
	CustomerName	string	`db:"customer_name"`
	Status	string	`db:"status"`
	CreatedAt int64	`db:"created_at"`
	UpdatedAt int64	`db:"updated_at"`
}

//OrderProduct type corresponds to the 'order_products' table in the DB and defines its schema.
//Can't specify foreign key constraints as currently, there is no support for foreign keys in gorp.
type OrderProduct struct {
	// db tag lets you specify the column name if it differs from the struct field
	Id      int64  `db:"order_product_id"`
	OrderId	int64	`db:"order_id"`
	ProductId	int64	`db:"product_id"`
	ProductEan	string	`db:"product_ean"`
	CustomerName	string	`db:"customer_name"`
	CreatedAt int64	`db:"created_at"`
	UpdatedAt int64	`db:"updated_at"`
}

//Customer type corresponds to the 'customers' table in the DB and defines its schema
type Customer struct {
	Id      int64  `db:"customer_id"`
	CustomerName	string	`db:"customer_name"`
	Address	string	`db:"address"`
	Phone	int64	`db:"phone"`
	CreatedAt int64	`db:"created_at"`
	UpdatedAt int64	`db:"updated_at"`
}

//Product type corresponds to the 'products' table in the DB and defines its schema
type Product struct {
	Id      int64  `db:"product_id" json:"product_id"`
	ProductName	string	`db:"product_name"`
	EanBarcode	string	`db:"product_ean" json:"product_ean"`
	Category	string	`db:"category"`
	SubCategory	string	`db:"sub_category"`
	Description	string	`db:"description"`
	CreatedAt int64	`db:"created_at"`
	UpdatedAt int64	`db:"updated_at"`
}

//InitDB establishes a connection to the required DB (based on SQL connection string).
//Initializes the DB with default tables required for running Order Microservice.
func InitDB() {

	//Invoke method to connect to DB and add table definitions to dbmap
	dbmap := connectToDB()
	defer dbmap.Db.Close()

	// Create the tables. In a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err := dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")
}

//connectToDB establishes a connection to the required DB (based on SQL connection string).
//Returns the DbMap for CRUD operations on the DB.
//Used by 'app.go' where all the APIs handlers implementation exists.
func connectToDB() *gorp.DbMap {

	//Passing a path relative to current directory. If run through Intellij by Run configs, fails as it expects absolute
	//file path from Project root such as --> "configs/dbconf.json"
	config := loadDBConfigs("../../configs/dbconf.json")

	//Build the DB connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Dbname)

	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("postgres", psqlInfo)
	checkErr(err, "sql.Open failed")
	//defer db.Close()

	err = db.Ping()
	checkErr(err, "Connection to DB or ping failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	// Add a table, setting the table name to 'orders', 'customers' and 'products' respectively, and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(Order{}, "orders").SetKeys(true, "Id")
	dbmap.AddTableWithName(Customer{}, "customers").SetKeys(true, "Id")
	dbmap.AddTableWithName(Product{}, "products").SetKeys(true, "Id")
	dbmap.AddTableWithName(OrderProduct{}, "order_products").SetKeys(true, "Id")

	return dbmap
}


//loadDBConfigs loads all the DB configs from JSON file.
//After loading, it returns the config json object.
func loadDBConfigs(filepath string) Configuration {

	configFile, err := os.Open(filepath)
	defer configFile.Close()
	checkErr(err, "Error reading DB configs from JSON file")
	jsonParser := json.NewDecoder(configFile)
	config := Configuration{}
	jsonParser.Decode(&config)
	return config
}

//checkErr checks for error and logs when present.
func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}