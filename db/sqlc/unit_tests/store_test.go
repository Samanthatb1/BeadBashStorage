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

// Test Scenario: Add a new order if user already exists
func TestNewOrderWithExistingUserTx(t *testing.T){
	// Create a new store to use regular DB operations plus the additional operations
	store := sqlc.NewStore(testDB)

	// Create User
	user := createRandomUser(t)

	// New order info
  newOrderParam := sqlc.NewOrderTxParams{
		Username: user.Username,
		FullName: util.RandomLongString(),
		PurchaseAmount: util.RandomCost(),
		PurchasedItem: util.RandomLongString(),
		ShippingLocation: util.RandomLongString(),
		Currency: util.RandomCurrency(),
		DateOrdered: util.RandomLongString(),
	}

	// Add new order given that user exists
	result, err := store.NewOrderTx(context.Background(), newOrderParam)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, result.EditedUser.TotalOrders, user.TotalOrders + 1) // User's total orders must increase by one
	require.Equal(t, result.EditedUser.Username, newOrderParam.Username)
	require.NotZero(t, result.EditedUser.TotalOrders)
	require.Equal(t, user.FullName, result.OrderMade.FullName)
}

// Test Scenario: Add a new order if user doesnt exist
func TestNewOrderWithNonExistingUserTx(t *testing.T){
	// Create a new store to use regular DB operations plus the additional operations
	store := sqlc.NewStore(testDB)

	// New order info
  newOrderParam := sqlc.NewOrderTxParams{
		Username: util.RandomLongString(), // This username must not exist in the db
		FullName: util.RandomLongString(),
		PurchaseAmount: util.RandomCost(),
		PurchasedItem: util.RandomLongString(),
		ShippingLocation: util.RandomLongString(),
		Currency: util.RandomCurrency(),
		DateOrdered: util.RandomLongString(),
	}

	// Search for user with the same username
	user, err := testQueries.GetUserByUsername(context.Background(), newOrderParam.Username)
	// User should not exist
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user)

	// Create order
	result, err := store.NewOrderTx(context.Background(), newOrderParam)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, result.EditedUser.TotalOrders, int64(1)) // User's total orders must be 1
	require.Equal(t, result.EditedUser.Username, newOrderParam.Username) // User's total orders must increase by one
	require.Equal(t, newOrderParam.FullName, result.OrderMade.FullName)
}

// Test Scenario: Delete order if it exists
func TestDeleteExistingOrderTx(t *testing.T){
	// Create a new store to use regular DB operations plus the additional operations
	store := sqlc.NewStore(testDB)

	// Create order and user
	user := createRandomUser(t)
	order := createRandomOrder(t, user)

	// Delete the order
	result, err := store.DeleteOrderTx(context.Background(), sqlc.DeleteOrderTxParams{OrderID: order.OrderID})
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, result.PurchasedItem, order.PurchasedItem)
	require.Equal(t, result.Status,"Deleted")

	// Make sure it was deleted
	orderAfter, err := testQueries.GetOrderById(context.Background(), order.OrderID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, orderAfter)
}

// Test Scenario: Delete order if it doesnt exist
func TestDeleteNonExistingOrderTx(t *testing.T){
	// Create a new store to use regular DB operations plus the additional operations
	store := sqlc.NewStore(testDB)

	// Delete the order
	result, err := store.DeleteOrderTx(context.Background(), sqlc.DeleteOrderTxParams{OrderID: util.RandomID()})
	require.Error(t, err)
	require.Empty(t, result)
}