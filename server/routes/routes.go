package routes

// SetupRoutes is a global setup function for all the routes.
func SetupRoutesV1(cfg *RouteConfig) {
	setupAuthRoutes(cfg)
	setupHealthRoutes(cfg)
}
