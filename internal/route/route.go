package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	serviceConfig "github.com/yeom-c/golnag-dynamodb-api/config"
	HealthHandler "github.com/yeom-c/golnag-dynamodb-api/internal/handler/health"
	ProductHandler "github.com/yeom-c/golnag-dynamodb-api/internal/handler/product"
	"github.com/yeom-c/golnag-dynamodb-api/internal/repository/adapter"
)

type Router struct {
	config *Config
	router *chi.Mux
}

func NewRouter() *Router {
	return &Router{
		config: NewConfig().SetTimeout(serviceConfig.GetConfig().Timeout),
		router: chi.NewRouter(),
	}
}

func (r *Router) SetRouter(repository adapter.Interface) *chi.Mux {
	r.setConfigRouter()
	r.RouterHealth(repository)
	r.RouterProduct(repository)

	return r.router
}

func (r *Router) setConfigRouter() {
	r.EnableCORS()
	r.EnableLogger()
	r.EnableTimeout()
	r.EnableRecover()
	r.EnableRequestID()
	r.EnableRealIP()

}

func (r *Router) RouterHealth(repository adapter.Interface) {
	handler := HealthHandler.NewHandler(repository)

	r.router.Route("/health", func(router chi.Router) {
		router.Get("/", handler.Get)
		router.Post("/", handler.Post)
		router.Put("/", handler.Put)
		router.Delete("/", handler.Delete)
		router.Options("/", handler.Options)
	})
}

func (r *Router) RouterProduct(repository adapter.Interface) {
	handler := ProductHandler.NewHandler(repository)

	r.router.Route("/product", func(router chi.Router) {
		router.Get("/", handler.Get)
		router.Get("/{ID}", handler.Get)
		router.Post("/", handler.Post)
		router.Put("/{ID}", handler.Put)
		router.Delete("/{ID}", handler.Delete)
		router.Options("/", handler.Options)
	})
}

func (r *Router) EnableTimeout() *Router {
	r.router.Use(middleware.Timeout(r.config.GetTimeout()))
	return r
}

func (r *Router) EnableCORS() *Router {
	r.router.Use(r.config.Cors)
	return r
}

func (r *Router) EnableLogger() *Router {
	r.router.Use(middleware.Logger)
	return r
}

func (r *Router) EnableRecover() *Router {
	r.router.Use(middleware.Recoverer)
	return r
}

func (r *Router) EnableRequestID() *Router {
	r.router.Use(middleware.RequestID)
	return r
}

func (r *Router) EnableRealIP() *Router {
	r.router.Use(middleware.RealIP)
	return r
}
