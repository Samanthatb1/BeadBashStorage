package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	sqlc "github.com/samanthatb1/beadBashStorage/db/sqlc"
)

/**** CREATE USER ****/

type createUserRequest struct {
	FullName    string `json:"full_name" binding:"required"`
	Username    string `json:"username" binding:"required"`
}

// Add createUser function to the server instance
func (server *Server) createUser(ctx *gin.Context){
	var reqBody createUserRequest;

	// If params are invalid
	if err := ctx.ShouldBindJSON(&reqBody); err != nil { 
		ctx.JSON(http.StatusBadRequest, errResponseToJSON(err))
		return
	}

	// Data is valid
	args := sqlc.CreateUserParams{
		FullName: reqBody.FullName,
		Username: reqBody.Username,
		TotalOrders: 0, // New user defaults to 0 orders
	}

	// Access the store we constructed through the server instance
	user, err := server.store.CreateUser(ctx, args)
	// Check if the DB insertion was successful 
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponseToJSON(err))
		return
	}

	// Success, send user back to client
	ctx.JSON(http.StatusOK, user)
}

/**** GET USER BY USERNAME OR ID ****/
type getUserUsernameRequest struct {
	Identifier    string   `uri:"identifier" binding:"required"`
}

// Add getUserByUsername function to the server instance
func (server *Server) getUserByUsername(ctx *gin.Context){
	var reqBody getUserUsernameRequest

	// If params are invalid
	if err := ctx.ShouldBindUri(&reqBody); err != nil { 
		ctx.JSON(http.StatusBadRequest, errResponseToJSON(err))
		return
	}

	// Check if client inputed an Id or Username
	id, err := strconv.ParseInt(reqBody.Identifier,10,64)
	var user sqlc.User;

	if err == nil { // If its a number
		// Access the store we constructed through the server instance
		user, err = server.store.GetUserById(ctx, id)
		// Check if the DB fetch was successful 
		if err != nil {
			if err == sql.ErrNoRows { // If that id doesnt exist
				ctx.JSON(http.StatusNotFound, gin.H{"error" : "User with that id doesn't exist"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, errResponseToJSON(err))
			return
		}
	} else { // If its a username
		// Access the store we constructed through the server instance
		user, err = server.store.GetUserByUsername(ctx, reqBody.Identifier)
		// Check if the DB fetch was successful 
		if err != nil {
			if err == sql.ErrNoRows { // If that id doesnt exist
				ctx.JSON(http.StatusNotFound, gin.H{"error" : "User with that username doesn't exist"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, errResponseToJSON(err))
			return
		}
	}

	// Success, send user to client
	ctx.JSON(http.StatusOK, user)
}

/**** LIST USERS ****/

type listUsersRequest struct {
	PageId    int32 `form:"page_id" binding:"required"`
	PageSize    int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// Add listUsers function to the server instance
func (server *Server) listUsers(ctx *gin.Context){
	var reqBody listUsersRequest;

	// If params are invalid
	if err := ctx.ShouldBindQuery(&reqBody); err != nil { 
		ctx.JSON(http.StatusBadRequest, errResponseToJSON(err))
		return
	}

	args := sqlc.ListUsersParams{
		Limit: reqBody.PageSize,
		Offset: (reqBody.PageId - 1) * reqBody.PageSize,
	}

	// Access the store we constructed through the server instance
	users, err := server.store.ListUsers(ctx, args)
	// Check if the DB fetch was successful 
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponseToJSON(err))
		return
	}

	// Success, send user back to client
	ctx.JSON(http.StatusOK, users)
}

/**** DELETE USER BY USERNAME ****/
type deleteUserByUsernameRequest struct {
	Username    string `uri:"username" binding:"required"`
}

// Add deleteUserById function to the server instance
func (server *Server) deleteUserByUsername(ctx *gin.Context){
	var reqBody deleteUserByUsernameRequest

	// If params are invalid
	if err := ctx.ShouldBindUri(&reqBody); err != nil { 
		ctx.JSON(http.StatusBadRequest, errResponseToJSON(err))
		return
	}

	// Make sure user exists
	user, err := server.store.GetUserByUsername(ctx, reqBody.Username)
	// Check if the DB fetch was successful 
	if err != nil {
		if err == sql.ErrNoRows { // If that id doesnt exist
			ctx.JSON(http.StatusNotFound, gin.H{"error" : "User doesn't exist"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponseToJSON(err))
		return
	}

	result, err := server.store.DeleteUserTx(ctx, sqlc.DeleteUserTxParams{ID : user.ID})
	// Check if the DB Delete was successful 
	if err != nil || result.Status != "Deleted" {
		ctx.JSON(http.StatusInternalServerError, errResponseToJSON(err))
		return
	}

	// Success, send user to client
	ctx.JSON(http.StatusOK, gin.H{"Deleted" : user.Username})
}