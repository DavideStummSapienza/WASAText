package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// | Unprotected Routes |

	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	// Image serving route
	rt.router.ServeFiles("/uploads/*filepath", http.Dir("./uploads")) // make images publicly accessable

	// Login
	rt.router.POST("/session", rt.login)

	// | Protected Routes |

	// Image Upload
	rt.router.POST("/upload", rt.wrapWithAuth(rt.uploadImage))

	// Search
	rt.router.GET("/users", rt.wrapWithAuth(rt.searchUsers))

	// User Profile
	rt.router.PUT("/user-profile", rt.wrapWithAuth(rt.changeUsername))
	rt.router.GET("/user-profile", rt.wrapWithAuth(rt.listConversations))
	rt.router.PUT("/profile-picture", rt.wrapWithAuth(rt.changeProfilePicture))

	// Conversation
	rt.router.GET("/conversations/:partner-username", rt.wrapWithAuth(rt.showConversation))
	rt.router.POST("/conversations/:partner-username", rt.wrapWithAuth(rt.sendMessage))
	rt.router.POST("/conversations/:partner-username/messages/:message-id", rt.wrapWithAuth(rt.forwardMessage))
	rt.router.DELETE("/conversations/messages/:message-id", rt.wrapWithAuth(rt.deleteMessage))

	// Comment
	rt.router.PUT("/conversations/messages/:message-id/comment", rt.wrapWithAuth(rt.makeComment))
	rt.router.DELETE("/conversations/messages/:message-id/comment", rt.wrapWithAuth(rt.deleteComment))

	// Groups
	rt.router.POST("/groups", rt.wrapWithAuth(rt.addToGroup))
	rt.router.PUT("/groups/:groupname", rt.wrapWithAuth(rt.changeGroupName))
	rt.router.DELETE("/groups/:groupname", rt.wrapWithAuth(rt.leaveGroup))
	rt.router.PUT("/groups/:groupname/group-photo", rt.wrapWithAuth(rt.changeGroupPicture))

	return rt.router
}
