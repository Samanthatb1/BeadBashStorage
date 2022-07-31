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

	// Create user and add order to that user
	user := createRandomUser(t)
	resultOrder, err := store.NewOrderTx(context.Background(), sqlc.NewOrderTxParams{
		Username: user.Username,
		FullName: util.RandomLongString(),
		PurchaseAmount: util.RandomCost(),
		PurchasedItem: util.RandomLongString(),
		ShippingLocation: util.RandomLongString(),
		Currency: util.RandomCurrency(),
		DateOrdered: util.RandomLongString(),
	})
	require.Equal(t, user.TotalOrders + 1, resultOrder.EditedUser.TotalOrders)

	// Delete the order
	result, err := store.DeleteOrderTx(context.Background(), sqlc.DeleteOrderTxParams{OrderID: resultOrder.OrderMade.OrderID})
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, result.DeletedItem, resultOrder.OrderMade.PurchasedItem)
	require.Equal(t, result.Status,"Deleted")

	// Make sure it was deleted
	deletedOrder, err := testQueries.GetOrderById(context.Background(), resultOrder.OrderMade.OrderID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, deletedOrder)

	// Make sure user's total order decreased by 1
	updatedUser, err := testQueries.GetUserById(context.Background(), user.ID)
	require.NoError(t, err)
	require.Equal(t, user.TotalOrders, updatedUser.TotalOrders)
}

// Test Scenario: Delete order if it doesnt exist
func TestDeleteNonExistingOrderTx(t *testing.T){
	// Create a new store to use regular DB operations plus the additional operations
	store := sqlc.NewStore(testDB)

	// Delete the order
	result, err := store.DeleteOrderTx(context.Background(), sqlc.DeleteOrderTxParams{OrderID: (util.RandomID() + util.RandomID() + util.RandomID())})
	require.Error(t, err)
	require.Empty(t, result)
}

// Test Scenario: Delete User if it exists
func TestDeleteUserIfExistsTx(t *testing.T){
	// Create a new store to use regular DB operations plus the additional operations
	store := sqlc.NewStore(testDB)

	// Create user with two orders
	user1 := createRandomUser(t)
	order1 := createRandomOrder(t, user1)
	order2 := createRandomOrder(t, user1)

	// Delete the user
	result1, err := store.DeleteUserTx(context.Background(), sqlc.DeleteUserTxParams{ID: user1.ID})
	require.NoError(t, err)
	require.NotEmpty(t, result1)

	// Check if their orders are gone too
	getOrder1, err := testQueries.GetOrderById(context.Background(), order1.OrderID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, getOrder1)
	getOrder2, err := testQueries.GetOrderById(context.Background(), order2.OrderID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, getOrder2)

	// Create User with no orders
	user2 := createRandomUser(t)

	// Delete the user
	result2, err := store.DeleteUserTx(context.Background(), sqlc.DeleteUserTxParams{ID: user2.ID})
	require.NoError(t, err)
	require.NotEmpty(t, result2)
}