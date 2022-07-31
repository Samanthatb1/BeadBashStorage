package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/samanthatb1/beadBashStorage/db/sqlc"
	_ "github.com/lib/pq" // provides the DB driver
)

type Server struct {
	store *db.Store // Defined in store.go: allows us to access inherent and additional DB operations
	router *gin.Engine // Router from gin
}

// New server instance
func NewServer(store *db.Store) *Server {
	// Instance
	server := &Server{store: store} // Assign store
	router := gin.Default()

	// Routes

	/* User */
	router.POST("/users", server.createUser) // Params: full_name, username
	router.GET("/users/:identifier", server.getUserByUsername) // Params: username
	router.GET("/users/all", server.listUsers) // Params: page_id, page_size
	router.DELETE("/users/:username", server.deleteUserByUsername) // Params: id
	router.PATCH("/users/:id", server.updateUserById)
	
	/* Order */
	router.POST("/orders", server.createOrder) // Params: username, name, all purchase info
	router.DELETE("/orders/:order_id", server.deleteOrderById) // Params: order_id
	router.PATCH("/orders", server.updateOrderById) // Params: order_id
	router.GET("/orders/:username", server.listOrdersOfUser) // Params: username
	router.GET("/orders/all", server.listAllOrders)

	server.router = router // Assign router
	return server
}

// Runs the Http server on a specified address port
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// Converts error message to a map
func errResponseToJSON(err error) gin.H {
	return gin.H{"error" : err.Error()}
}