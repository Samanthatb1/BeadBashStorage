// Store to hold any transactions
package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Queries only supports individual DB operation functions
// Golang Composition: Extending Store to include Queries function plus it's own functions
type Store struct {
	*Queries
	db *sql.DB // required to create a new DB transaction
}

func NewStore(db *sql.DB) *Store {
	// Build the store object
	return &Store{
		db : db, // sql db
		Queries: New(db), // Queries from ./db.go
	}
}

// Ensure that all transactions are commited ONLY if there is no errors; rollback if error is found
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// Begin the transaction
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	queries := New(tx)
	err = fn(queries)
	if err != nil {
		// Rollback so that any DB actions are not commited
		if rbErr := tx.Rollback(); rbErr != nil{
			return fmt.Errorf("tx err: %v, rollbackk err: %v", err, rbErr)
		}
		return err
	}

	// If all operations within the transaction are succesful, commit the changes
	return tx.Commit()
} 

type NewOrderTxParams struct {
	Username         string  `json:"username"`
	FullName         string  `json:"full_name"`
	PurchaseAmount   float64 `json:"purchase_amount"`
	PurchasedItem    string  `json:"purchased_item"`
	ShippingLocation string  `json:"shipping_location"`
	Currency         string  `json:"currency"`
	DateOrdered      string  `json:"date_ordered"`
}

type newOrderResult struct {
	EditedUser User `json:"edited_user"`
	OrderMade Order `json:"order_made"`
}

// Add new order -> Must handle if the user for that order exists or not
func (store *Store) NewOrderTx(ctx context.Context, args NewOrderTxParams) (newOrderResult, error){
	var result newOrderResult

	// Create a new DB transaction
	err := store.execTx(ctx, func(q *Queries) error{

	// Check to see if a user with the inputted username already exists
	user, err := q.GetUserByUsername(ctx, args.Username)

		// If User doesnt exist
    if err == sql.ErrNoRows { 
			newUser, err := q.CreateUser(ctx, CreateUserParams{
				FullName: args.FullName,
				Username: args.Username,
				TotalOrders: 1,
			})
			if err != nil { return err }
			// Set the user to the result
			result.EditedUser = newUser

		// If theres an Error
    } else if err != nil{ 
				return err

		// User already exists
    } else { 
				// Increase users order number by 1
				updatedUser, err := q.UpdateUser(ctx, UpdateUserParams{
					ID: user.ID, 
					TotalOrders: (user.TotalOrders + 1),
				})
				if err != nil {return err}

				// Set the user to the result
				result.EditedUser = updatedUser
		}

	 // Create the new order
	 order, err := q.CreateOrder(ctx, CreateOrderParams{
		AccountID: result.EditedUser.ID,
		Username: result.EditedUser.Username,
		FullName: result.EditedUser.FullName,
		PurchaseAmount: args.PurchaseAmount,
		PurchasedItem: args.PurchasedItem,
		ShippingLocation: args.ShippingLocation,
		Currency: args.Currency,
		DateOrdered: args.DateOrdered,
	 })
	 if err != nil{ return err }

	 result.OrderMade = order
	 return nil // No Error
	})
	// Return the new order and the updated user
	return result, err
}

/***************************************/

type DeleteOrderTxParams struct {
	OrderID        int64   `json:"order_id"`
}

type deleteOrderResult struct {
	Status string `json:"deletion_status"`
	PurchasedItem string `json:"purchased_item"`
}

// Order is deleted -> Must update the associated user information
func (store *Store) DeleteOrderTx(ctx context.Context, args DeleteOrderTxParams) (deleteOrderResult, error){
	var result deleteOrderResult

	// Begin Transaction
	err := store.execTx(ctx, func(q *Queries) error{
		// Make sure the order exists before deleting it
		order, err := q.GetOrderById(ctx, args.OrderID)
		if err == sql.ErrNoRows { 
			return err
		}

		err = q.DeleteOrder(ctx, order.OrderID)
		if err != nil {return err}

		result.PurchasedItem = order.PurchasedItem
		result.Status = "Deleted"
		return nil // No Error
	})

	return result, err
}
