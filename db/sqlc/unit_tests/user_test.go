// Unit tests for Users' CRUD operations

package tests

import (
	"context"
	"database/sql"
	"testing"
	"time"

	sqlc "github.com/samanthatb1/beadBashStorage/db/sqlc"
	"github.com/samanthatb1/beadBashStorage/util"
	"github.com/stretchr/testify/require"
)

/* Helper Functions */

func createRandomUser(t *testing.T) sqlc.User {
	// Set fields of new user
	newUserParams := sqlc.CreateUserParams{
		FullName: util.RandomLongString(),
		Username: util.RandomLongString(),
		TotalOrders: util.RandomOrders(),
	}

	// Create the new user
	user, err := testQueries.CreateUser(context.Background(), newUserParams)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	// Check that the inserted user fields match the fields we passed
	require.Equal(t, newUserParams.FullName, user.FullName)
	require.Equal(t, newUserParams.Username, user.Username)
	require.Equal(t, newUserParams.TotalOrders, user.TotalOrders)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func allUserFieldsEqual(t *testing.T, user1 sqlc.User, user2 sqlc.User, testOrder bool){
	require.Equal(t, user1.ID, user2.ID) 	// ID matches
	require.Equal(t, user1.Username, user2.Username) 	// username matches
	require.Equal(t, user1.FullName, user2.FullName) // full name matches
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second) // created at matches within the second
	if testOrder {
		require.Equal(t, user1.TotalOrders, user2.TotalOrders) // full name matches
	}
}

/* Test Functions */

// Test Scenario: create a new user entry
func TestCreateUser(t *testing.T){
	createRandomUser(t)
}

// Test Scenario: get a user that exists
func TestGetUser(t *testing.T){
	createdUser := createRandomUser(t)

	// Should be able to get user by both ID and Username
	userById, err := testQueries.GetUserById(context.Background(), createdUser.ID)
	require.NoError(t, err)
	userByUsername, err := testQueries.GetUserByUsername(context.Background(), createdUser.Username)
	require.NoError(t, err)

	allUserFieldsEqual(t, createdUser, userById, true)
	allUserFieldsEqual(t, createdUser, userByUsername, true)
}

// Test Scenario: update user total order number
func TestUpdateUser(t *testing.T){
	createdUser := createRandomUser(t)

	userParams := sqlc.UpdateUserParams{
		ID:	createdUser.ID,
		TotalOrders: util.RandomOrders(),
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), userParams)
	require.NoError(t, err)
	allUserFieldsEqual(t, createdUser, updatedUser, false)
	require.Equal(t, userParams.TotalOrders, updatedUser.TotalOrders)
}

// Test Scenario: delete user
func TestDeleteUser(t *testing.T){
	createdUser := createRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), createdUser.Username)
	require.NoError(t, err)

	// make sure it's truly gone by searching for it
	getUser, err := testQueries.GetUserByUsername(context.Background(), createdUser.Username)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, getUser)
}

// Test Scenario: List Users
func TestListUser(t *testing.T){
	// make 10 new users
	for i:= 0 ; i < 10 ; i++{
		createRandomUser(t)
	}
	// fetch 5 users, skipping the first 5
	listUserParams := sqlc.ListUsersParams{
		Limit: 5,
		Offset: 5,
	}
	// make sure 5 non empty users are returned
	users, err := testQueries.ListUsers(context.Background(), listUserParams)
	require.NoError(t, err)
	require.Len(t, users, 5)
	for _, user := range users{
		require.NotEmpty(t, user)
	}
}

func TestNonExistUser(t *testing.T){
	getUser, err := testQueries.GetUserByUsername(context.Background(), util.RandomLongString()) // Username doesnt exist
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, getUser) 
}