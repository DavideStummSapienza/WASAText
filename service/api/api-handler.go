package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Login
	rt.router.POST("/session", rt.login)

	// Geschützte Route mit Middleware
	//protected := AuthMiddleware(http.HandlerFunc(yourProtectedHandler))
	//rt.router.Handler("GET", "/protected-endpoint", protected)

	// Search
	rt.router.GET("/users", rt.searchUsers)

	// User Profile
	rt.router.PUT("/user-profile", rt.changeUsername)
	rt.router.GET("/user-profile", rt.listConversations)
	rt.router.PUT("/profile-picture", rt.changeProfilePicture)

	// Conversation
	rt.router.GET("/conversations/:partner-username", rt.showConversation)
	rt.router.POST("/conversations/:partner-username", rt.sendMessage)
	rt.router.POST("/conversations/:partner-username/messages/:message-timestamp", rt.forwardMessage)
	rt.router.DELETE("/conversations/:partner-username/messages/:message-timestamp", rt.deleteMessage)

	// Comment
	rt.router.PUT("/conversations/:partner-username/messages/:message-timestamp/comment", rt.makeComment)
	rt.router.DELETE("/conversations/:partner-username/messages/:message-timestamp/comment", rt.deleteComment)

	// Groups
	rt.router.PUT("/groups", rt.addToGroup)
	rt.router.PUT("/groups/{groupname}", rt.changeGroupName)
	rt.router.DELETE("/groups/:groupname", rt.leaveGroup)
	rt.router.PUT("/groups/:groupname/group-photo", rt.changeGroupPicture)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
