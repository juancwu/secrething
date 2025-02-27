package routes

import (
	"github.com/juancwu/konbini/server/handlers"
	"github.com/juancwu/konbini/server/middlewares"
	"reflect"
)

func setupGroupRoutes(routeConfig *RouteConfig) {
	e := routeConfig.Echo

	e.POST(
		"/group/new",
		handlers.NewGroup(routeConfig.DBConnector),
		// only allow request that comes with a full token
		middlewares.ProtectFull(routeConfig.DBConnector),
		middlewares.ValidateJson(reflect.TypeOf(handlers.NewGroupRequest{})),
	)

	e.DELETE(
		"/group/:id",
		handlers.DeleteGroup(routeConfig.DBConnector),
		middlewares.ProtectFull(routeConfig.DBConnector),
	)

	e.POST(
		"/group/invite",
		handlers.InviteUsersToJoinGroup(routeConfig.DBConnector),
		middlewares.ProtectFull(routeConfig.DBConnector),
		middlewares.ValidateJson(reflect.TypeOf(handlers.InviteUsersToJoinGroupRequest{})),
	)

	e.GET(
		"/group/invitation/accept",
		handlers.AcceptGroupInvitation(routeConfig.DBConnector),
	)
}
