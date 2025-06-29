package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.GET("/users/photo", rt.wrap(rt.getMyPhoto))
	rt.router.PUT("/users/photo", rt.wrap(rt.setMyPhoto))
	rt.router.PUT("/users/name", rt.wrap(rt.setMyUserName))
	rt.router.GET("/chats", rt.wrap(rt.getMyChats))
	rt.router.POST("/chats", rt.wrap(rt.startChat))
	rt.router.GET("/groups", rt.wrap(rt.getMyGroups))
	rt.router.POST("/groups", rt.wrap(rt.createGroup))
	rt.router.GET("/search", rt.wrap(rt.searchUsers))
	rt.router.GET("/chats/:chatId", rt.wrap(rt.getChat))
	rt.router.POST("/chats/:chatId/message", rt.wrap(rt.sendMessage))
	rt.router.DELETE("/chats/:chatId/message/:messageId", rt.wrap(rt.deleteMessage))
	rt.router.POST("/chats/:chatId/message/:messageId/forward", rt.wrap(rt.forwardMessage))
	rt.router.POST("/chats/:chatId/message/:messageId/comment", rt.wrap(rt.commentMessage))
	rt.router.DELETE("/chats/:chatId/message/:messageId/comment", rt.wrap(rt.uncommentMessage))
	rt.router.GET("/groups/:groupId", rt.wrap(rt.getGroup))
	rt.router.DELETE("/groups/:groupId", rt.wrap(rt.leaveGroup))
	rt.router.POST("/groups/:groupId", rt.wrap(rt.addToGroup))
	rt.router.PUT("/groups/:groupId/name", rt.wrap(rt.setGroupName))
	rt.router.PUT("/groups/:groupId/photo", rt.wrap(rt.setGroupPhoto))
	rt.router.GET("/liveness", rt.liveness)
	return rt.router
}
