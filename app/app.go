//Package app specifies the models (struct) and the business logic for the application (Order microservice).
//Abstracts the complexity from executables (inside /cmd) and keeps them as clean as possible.
package app

import (
	"github.com/gin-gonic/gin"
	"time"
	"net/http"
	"log"
)

//Constants to specify the order status when adding/updating in 'orders' table
//The statuses also correspond to a typical Order lifecycle.
const (
	received = "RECEIVED"
	in_progress = "IN PROGRESS"
	shipped = "SHIPPED"
	delivered = "DELIVERED"
	cancelled = "CANCELLED"
)

//NewOrderReq type corresponds to the incoming POST request data for Add Order API.
type NewOrderReq struct {
	CustomerName	string	`json:"customer_name"`
	Products	[]Product	`json:"products"`
}

//OrderResponse type corresponds to the outgoing response for the GET request for fetching an order details.
type OrderResponse struct {
	OrderId	int64	`json:"order_id"`
	ProductId	int64	`json:"product_id"`
	ProductEan	string	`json:"product_ean"`
	CustomerName	string	`json:"customer_name"`
}

//Connect to DB and get the DbMap
var dbmap = connectToDB()

//Ping is used to check the health of the ORDER service.
//If service is up and running, it returns a status of '200 OK'
func Ping(c *gin.Context) {}

//CreateOrder accepts the create order request and stores order details in the DB.
func CreateOrder(c *gin.Context) {
	var orderReq NewOrderReq
	c.Bind(&orderReq)

	//Validate if the customer_name is empty or not from request data.
	if isEmpty(orderReq.CustomerName){
		c.JSON(http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest,"error": "Customer name cannot be empty. Pass a valid string value and try again !"})
		return
	}
	//Validate if the products array is empty or not.
	if len(orderReq.Products) == 0 {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest,"error": "Products cannot be empty. An order need to have atleast 1 product. Add a product and try again !"})
		return
	}

	order := &Order{
		CustomerName:	orderReq.CustomerName,
		Status:	received,
		CreatedAt:	time.Now().UnixNano(),
		UpdatedAt:	time.Now().UnixNano(),
	}

	//Insert the new order single record in the 'orders' table.
	err := dbmap.Insert(order)
	checkErr(err, "Add new order failed in orders table")

	//Iterate over the request data 'products array' and for each product
	//in this new order, add an entry into the 'order_products' table.
	for _, product := range orderReq.Products {

		//Validate if the product_ean is a valid EAN-13 string
		if isEAN(product.EanBarcode) {
			orderProduct := &OrderProduct{
				OrderId:      order.Id,
				ProductId:    product.Id,
				ProductEan:   product.EanBarcode,
				CustomerName: orderReq.CustomerName,
				CreatedAt:    time.Now().UnixNano(),
				UpdatedAt:    time.Now().UnixNano(),
			}
			//Insert the order -> product mapping record in the 'order_products' table.
			err := dbmap.Insert(orderProduct)
			checkErr(err, "Add new order_product mapping failed in order_products table")
		} else {
			errMsg := "Product EAN is invalid : " + product.EanBarcode
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest,"error": errMsg})
			return
		}
	}//End of for-loop

	c.JSON(http.StatusCreated,
		gin.H{"status": http.StatusCreated, "message": "Order Created Successfully!", "resourceId": order.Id})
}

//FetchOrder retrieves the order details from DB based on the 'order_id' passed in path param.
func FetchOrder(c *gin.Context) {
	orderId := c.Params.ByName("id")

	//Validate if the input ID is a valid number or not
	if isNumber(orderId) {
		//Do nothing as order ID is a valid int
	} else {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest,"error": "Order ID passed is not a valid number."})
		return
	}

	log.Println("Fetching order details for order ID : " + orderId)
	var orderProducts []OrderProduct
	var query = "SELECT * FROM order_products where order_id=" + orderId + " ORDER BY order_product_id"

	_, err := dbmap.Select(&orderProducts, query)

	if len(orderProducts)==0 || err != nil {
		c.JSON(http.StatusNotFound,
			gin.H{"status": http.StatusNotFound,"error": "No order with requested ID exists in the table. Invalid ID."})
	} else {
		//Iterate over the orderProducts array and for each record,
		//transform it into response model type.
		var orderResArr []OrderResponse
		for _, op := range orderProducts {
			orderRes := OrderResponse {
				OrderId:	op.OrderId,
				ProductId:	op.ProductId,
				ProductEan:	op.ProductEan,
				CustomerName:	op.CustomerName,
			}
			orderResArr = append(orderResArr, orderRes)
		}//End of for-loop

		c.JSON(http.StatusOK,
			gin.H{"status": http.StatusOK, "message": "Order Details Fetched Successfully!", "order": orderResArr})
	}//End of else-block
}

//UpdateOrder retrieves an order details from DB based on the 'order_id' passed in path param.
//After fetching order details, it updates the old details with new details for the same order record.
func UpdateOrder(c *gin.Context) {
	orderId := c.Params.ByName("id")

	//Validate if the input ID is a valid number or not
	if isNumber(orderId) {
		//Do nothing as order ID is a valid int
	} else {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest,"error": "Order ID passed is not a valid number."})
		return
	}

	log.Println("Updating order details for order ID :", orderId)
	var query_orders = "SELECT * FROM orders where order_id=" + orderId
	var order Order
	err := dbmap.SelectOne(&order, query_orders)

	//Case when error occurred, or no record with given 'order_id' exists in the 'orders' table.
	if err != nil || (Order{}) == order || len(order.CustomerName) == 0 {
		c.JSON(http.StatusNotFound,
			gin.H{"status": http.StatusNotFound,"error": "No order with requested ID exists in the table. Invalid ID."})
	} else {
		//Bind the request data to struct
		var orderReq NewOrderReq
		c.Bind(&orderReq)

		//Update the 'customer_name' (if changed) and 'UpdatedAt' columns in the 'orders' table
		if orderReq.CustomerName == order.CustomerName {
			order.UpdatedAt = time.Now().UnixNano()
		} else {
			order.CustomerName = orderReq.CustomerName
			order.UpdatedAt = time.Now().UnixNano()
		}

		//Updating the order record in the 'orders' table
		_, err = dbmap.Update(&order)

		var orderProducts []OrderProduct
		var query_order_products = "SELECT * FROM order_products where order_id=" + orderId + " ORDER BY order_product_id"
		_, err := dbmap.Select(&orderProducts, query_order_products)

		if len(orderProducts)==0 || err != nil {
			c.JSON(http.StatusNotFound,
				gin.H{"status": http.StatusNotFound,"error": "No products found with requested order ID. Aborting!"})
		} else if len(orderProducts) != len(orderReq.Products) {
			c.JSON(http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest,"error": "v1 Update API supports only product details updation in the existing order. New products addition and existing products deletion from existing order will be supported in future API ver. Aborting!"})
		} else {
			//Iterate over the fetched 'OrderProduct array' from 'order_products' table for given order.
			//For each record, update data from the 'products array' from request data.
			//Write updated data back into the 'order_products' table.
			for index, orderProduct := range orderProducts {

				product := orderReq.Products[index]

				orderProduct.CustomerName = orderReq.CustomerName
				orderProduct.ProductId = product.Id
				orderProduct.ProductEan = product.EanBarcode
				orderProduct.UpdatedAt = time.Now().UnixNano()

				//Updating the order -> product mapping record in the 'order_products' table.
				_, err := dbmap.Update(&orderProduct)
				checkErr(err, "Updating order_product mapping failed in order_products table")
			} //End of for-loop

			c.JSON(http.StatusOK,
				gin.H{"status": http.StatusOK, "message": "Order Updated Successfully!", "resourceId": order.Id})
		}
	}
}