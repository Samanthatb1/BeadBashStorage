// Unit tests for orders' CRUD operations

package tests

import (
	"context"
	"database/sql"
	"testing"

	sqlc "github.com/samanthatb1/beadBashStorage/db/sqlc"
	"github.com/samanthatb1/beadBashStorage/util"
	"github.com/stretchr/testify/require"
)

/* Helper Functions */

func allOrderFieldsEqual(t *testing.T, order1 sqlc.Order, order2 sqlc.Order, updateFields bool){
	require.Equal(t, order1.OrderID, order2.OrderID) 	// Order ID matches
	require.Equal(t, order1.AccountID, order2.AccountID) 	// Account ID matches
	require.Equal(t, order1.Username, order2.Username) 	// username matches
	require.Equal(t, order1.FullName, order2.FullName) // full name matches
	require.Equal(t, order1.Currency, order2.Currency) // currency matches
	require.Equal(t, order1.DateOrdered, order2.DateOrdered) // date ordered matches

	if (!updateFields){
		require.Equal(t, order1.PurchaseAmount, order2.PurchaseAmount) // purchase amount matches
		require.Equal(t, order1.PurchasedItem, order2.PurchasedItem) // purchased item matches
		require.Equal(t, order1.ShippingLocation, order2.ShippingLocation) // shipping location matches
	}
}

func createRandomOrder(t *testing.T, user sqlc.User) sqlc.Order{
	createOrderParams := sqlc.CreateOrderParams{
		AccountID: user.ID,
		Username: user.Username,
		FullName: util.RandomLongString(),
		PurchaseAmount: util.RandomCost(),
		PurchasedItem: util.RandomLongString(),
		ShippingLocation: util.RandomLongString(),
		Currency: util.RandomCurrency(),
		DateOrdered: util.RandomLongString(),
	}

	createdOrder, err := testQueries.CreateOrder(context.Background(), createOrderParams)
	require.NoError(t, err)
	require.NotEmpty(t, createdOrder)

	// Check that the inserted user fields match the fields we passed
	require.Equal(t, createOrderParams.AccountID, createdOrder.AccountID)
	require.Equal(t, createOrderParams.Username, createdOrder.Username)
	require.Equal(t, createOrderParams.FullName, createdOrder.FullName)
	require.Equal(t, createOrderParams.PurchaseAmount, createdOrder.PurchaseAmount)
	require.Equal(t, createOrderParams.PurchasedItem, createdOrder.PurchasedItem)
	require.Equal(t, createOrderParams.ShippingLocation, createdOrder.ShippingLocation)
	require.Equal(t, createOrderParams.Currency, createdOrder.Currency)
	require.Equal(t, createOrderParams.DateOrdered, createdOrder.DateOrdered)

	require.NotZero(t, createdOrder.OrderID)

	return createdOrder
}

/* Test Functions */

// Test Scenario: update users total order number
func TestCreateOrder(t *testing.T){
	user := createRandomUser(t)
	createRandomOrder(t, user)
}

// Test Scenario: get order by order id
func TestGetOrderById(t *testing.T){
	user := createRandomUser(t)
	createdOrder := createRandomOrder(t, user)

	fetchedOrder, err := testQueries.GetOrderById(context.Background(), createdOrder.OrderID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedOrder)

	allOrderFieldsEqual(t, createdOrder, fetchedOrder, false)
}

// Test Scenario: get order by order id
func TestListOrdersByUsername(t *testing.T){
	user := createRandomUser(t)

	// make 10 new orders based on the specific user previously created
	for i:= 0 ; i < 10 ; i++{
		createRandomOrder(t, user)
	}

	// make sure 10 non empty orders are returned
	orders, err := testQueries.ListOrdersByUsername(context.Background(), user.Username)
	require.NoError(t, err)
	require.Len(t, orders, 10)
	for _, order := range orders{
		require.NotEmpty(t, order)
	}
}

// Test Scenario: get all orders
func TestListAllOrders(t *testing.T){
	user := createRandomUser(t)
	// make 10 new orders
	for i:= 0 ; i < 10 ; i++{
		createRandomOrder(t, user)
	}
	// fetch 5 orders, skipping the first 5
	listAllOrdersParams := sqlc.ListAllOrdersParams{
		Limit: 5,
		Offset: 5,
	}
	// make sure 5 non empty users are returned
	orders, err := testQueries.ListAllOrders(context.Background(), listAllOrdersParams)
	require.NoError(t, err)
	require.Len(t, orders, 5)
	for _, order := range orders{
		require.NotEmpty(t, order)
	}
}

// Test Scenario: update order
func TestUpdateOrder(t *testing.T){
	user := createRandomUser(t)
	createdOrder := createRandomOrder(t, user)

	updateOrderParams := sqlc.UpdateOrderParams{
		OrderID:	createdOrder.OrderID,
		PurchasedItem: util.RandomLongString(),
		PurchaseAmount: util.RandomCost(),
    ShippingLocation: util.RandomLongString(),
	}

	updatedOrder, err := testQueries.UpdateOrder(context.Background(), updateOrderParams)
	require.NoError(t, err)
	allOrderFieldsEqual(t, createdOrder, updatedOrder, true)
	require.Equal(t, updateOrderParams.PurchasedItem, updatedOrder.PurchasedItem)
	require.Equal(t, updateOrderParams.PurchaseAmount, updatedOrder.PurchaseAmount)
	require.Equal(t, updateOrderParams.ShippingLocation, updatedOrder.ShippingLocation)
}

// Test Scenario: delete order
func TestDeleteOrder(t *testing.T){
	user := createRandomUser(t)
	createdOrder := createRandomOrder(t, user)
	err := testQueries.DeleteOrder(context.Background(), createdOrder.OrderID)
	require.NoError(t, err)

	// make sure it's truly gone by searching for it
	getOrder, err := testQueries.GetOrderById(context.Background(), createdOrder.OrderID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, getOrder)
}