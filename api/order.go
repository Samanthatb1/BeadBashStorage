package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	sqlc "github.com/samanthatb1/beadBashStorage/db/sqlc"
)

/**** CREATE ORDER ****/
type createOrderRequest struct {
	Username         string  `json:"username" binding:"required"`
	FullName         string  `json:"full_name" binding:"required"`
	PurchaseAmount   float64 `json:"purchase_amount" binding:"required"`
	PurchasedItem    string  `json:"purchased_item" binding:"required"`
	ShippingLocation string  `json:"shipping_location" binding:"required"`
	Currency         string  `json:"currency" binding:"required,oneof=USD EUR CAD"`
	DateOrdered      string  `json:"date_ordered" binding:"required"`
}

// Add createOrder function to the server instance
func (server *Server) createOrder(ctx *gin.Context){
	var reqBody createOrderRequest;

	// If params are invalid
	if err := ctx.ShouldBindJSON(&reqBody); err != nil { 
		ctx.JSON(http.StatusBadRequest, errResponseToJSON(err))
		return
	}

	// Create DB params
	newOrderParams := sqlc.NewOrderTxParams{
		Username: reqBody.Username,
		FullName: reqBody.FullName,
		PurchaseAmount: reqBody.PurchaseAmount,
		PurchasedItem: reqBody.PurchasedItem,
		ShippingLocation: reqBody.ShippingLocation,
		Currency: reqBody.Currency,
		DateOrdered: reqBody.DateOrdered,
	}

	// Access the store we constructed through the server instance
	result, err := server.store.NewOrderTx(ctx, newOrderParams)

	// Check if the DB insertion was successful 
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponseToJSON(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

/**** DELETE ORDER ****/
type deleteOrderRequest struct {
	OrderId   int64  `uri:"order_id" binding:"required"`
}

// Add deleteOrder function to the server instance
func (server *Server) deleteOrderById(ctx *gin.Context){
	var reqBody deleteOrderRequest;

	// If params are invalid
	if err := ctx.ShouldBindUri(&reqBody); err != nil { 
		ctx.JSON(http.StatusBadRequest, errResponseToJSON(err))
		return
	}

	// Access the store we constructed through the server instance
	result, err := server.store.DeleteOrderTx(ctx, sqlc.DeleteOrderTxParams{OrderID: reqBody.OrderId})
	// Check if the DB deletion was successful 
	if err != nil {
		if err == sql.ErrNoRows { // If that id doesnt exist
			ctx.JSON(http.StatusNotFound, gin.H{"error" : "Order doesn't exist"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponseToJSON(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

/**** DELETE ORDER ****/
type updateOrderByIdRequest struct {
	OrderId  				 int64  `json:"order_id" binding:"required"`
	PurchaseAmount   float64 `json:"purchase_amount"`
	PurchasedItem    string  `json:"purchased_item"`
	ShippingLocation string  `json:"shipping_location"`
}

// Add deleteOrder function to the server instance
func (server *Server) updateOrderById(ctx *gin.Context){
	var reqBody updateOrderByIdRequest;

	// If params are invalid
	if err := ctx.ShouldBindJSON(&reqBody); err != nil { 
		ctx.JSON(http.StatusBadRequest, errResponseToJSON(err))
		return
	}

	// Get order to update
	orderToUpdate, err := server.store.GetOrderById(ctx, reqBody.OrderId)
	// Check if DB search was successful
	if err != nil {
		if err == sql.ErrNoRows { // If that id doesnt exist
			ctx.JSON(http.StatusNotFound, gin.H{"error" : "Order doesn't exist"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponseToJSON(err))
		return
	}

	updateOrderParams := sqlc.UpdateOrderParams{}
	updateOrderParams.OrderID = reqBody.OrderId

	// If user sent data to update, change it; if not, keep the same
	if reqBody.PurchaseAmount == 0 { 
		updateOrderParams.PurchaseAmount = orderToUpdate.PurchaseAmount
	} else { updateOrderParams.PurchaseAmount = reqBody.PurchaseAmount }

	if reqBody.PurchasedItem == "" { 
		updateOrderParams.PurchasedItem = orderToUpdate.PurchasedItem
	} else { updateOrderParams.PurchasedItem = reqBody.PurchasedItem }

	if reqBody.ShippingLocation == "" { 
		updateOrderParams.ShippingLocation = orderToUpdate.ShippingLocation
	} else { updateOrderParams.ShippingLocation = reqBody.ShippingLocation }

	result, err := server.store.UpdateOrder(ctx, updateOrderParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponseToJSON(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

/**** LIST ORDERS FROM USER ****/
type listOrdersOfUserRequest struct {
	Username   string  `uri:"username" binding:"required"`
}

// Add listOrdersOfUser function to the server instance
func (server *Server) listOrdersOfUser(ctx *gin.Context){
	var reqBody listOrdersOfUserRequest;

	// If params are invalid
	if err := ctx.ShouldBindUri(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponseToJSON(err))
		return
	}
	orders, err := server.store.ListOrdersByUsername(ctx, reqBody.Username)
	// Check if DB search was successful
	if err != nil {
		if err == sql.ErrNoRows { // If that id doesnt exist
			ctx.JSON(http.StatusNotFound, gin.H{"error" : "User doesn't exist"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponseToJSON(err))
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

/**** LIST ORDERS ****/
type listOrdersRequest struct {
	PageId    int32 `form:"page_id" binding:"required"`
	PageSize  int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// Add listUsers function to the server instance
func (server *Server) listAllOrders(ctx *gin.Context){
	var reqBody listUsersRequest;

	// If params are invalid
	if err := ctx.ShouldBindQuery(&reqBody); err != nil { 
		ctx.JSON(http.StatusBadRequest, errResponseToJSON(err))
		return
	}

	args := sqlc.ListAllOrdersParams{
		Limit: reqBody.PageSize,
		Offset: (reqBody.PageId - 1) * reqBody.PageSize,
	}

	// Access the store we constructed through the server instance
	orders, err := server.store.ListAllOrders(ctx, args)
	// Check if the DB fetch was successful 
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponseToJSON(err))
		return
	}

	// Success, send user back to client
	ctx.JSON(http.StatusOK, orders)
}