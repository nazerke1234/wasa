package routes

import (
	controllers2 "wasa/backend/service/api/controllers"
	"wasa/backend/service/api/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes the routes
func SetupRoutes(r *gin.Engine) {
	// Route for login
	r.POST("/login", controllers2.Login)
	r.POST("/refresh", controllers2.RefreshToken)

	api := r.Group("/api")
	{
		api.GET("/users", controllers2.SearchUsers) // Search users

		authorized := api.Group("/")
		authorized.Use(middleware.AuthMiddleware()) // JWT required
		{
			// All routes requiring authorization, including chat creation
			authorized.GET("/me", controllers2.GetProfile)               // Get user data
			authorized.PUT("/me/username", controllers2.UpdateUsername)  // Update username
			authorized.PUT("/me/photo", controllers2.UpdateProfilePhoto) // Update profile photo

			authorized.GET("/chats", controllers2.GetMyConversations)
			authorized.GET("/chats/:chat_id", controllers2.GetConversation)
			authorized.POST("/messages", controllers2.SendMessage)
			authorized.POST("/messages/forward", controllers2.ForwardMessage)
			authorized.PUT("/messages/:message_id/comment", controllers2.CommentMessage)
			authorized.DELETE("/messages/:message_id/comment", controllers2.UncommentMessage)
			authorized.DELETE("/messages/:message_id", controllers2.DeleteMessage)

			authorized.POST("/group/:groupId/add", controllers2.AddUserToGroup)          // Add user to group
			authorized.POST("/group/:groupId/leave", controllers2.LeaveGroup)            // Leave group
			authorized.PUT("/group/:groupId/updateName", controllers2.SetGroupName)      // Update group name
			authorized.PUT("/group/:groupId/updatePhoto", controllers2.UpdateGroupPhoto) // Update group photo
			authorized.GET("/group/:groupId/info", controllers2.GetGroupInfo)

			// Move chat creation here to require authorization
			authorized.POST("/personalChats", controllers2.CreatePersonalChat)
			authorized.POST("/groupChats", controllers2.CreateGroupChat) // For creating chat
		}
	}
}
