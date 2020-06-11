package restful

import (
	// log "github.com/sirupsen/logrus"
	"net/http"

	goji "goji.io/v3"
	"goji.io/v3/pat"

	"go-prj-skeleton/appprj/interface/restful/handler"
	"go-prj-skeleton/appprj/interface/restful/middleware"
	"go-prj-skeleton/appprj/jsonutil"
	"go-prj-skeleton/appprj/registry"
	"go-prj-skeleton/appprj/usecase"
)

// Handlers ...
func Handlers(ctn *registry.Container) *goji.Mux {
	mux := goji.NewMux()
	mux.Use(middleware.JSON)

	mux.HandleFunc(pat.Get("/"), Info)
	apiRoute := goji.SubMux()
	mux.Handle(pat.New("/api/*"), apiRoute)

	userHandler := handler.NewUserHandler(ctn.Resolve("user-usecase").(usecase.UserUsecase))

	apiRoute.HandleFunc(pat.Get("/users/:user_id/transactions"), userHandler.FindTransactions)
	apiRoute.HandleFunc(pat.Post("/users/:user_id/transactions"), userHandler.CreateTransaction)
	apiRoute.HandleFunc(pat.Put("/users/:user_id/transactions/:transaction_id"), userHandler.UpdateTransaction)
	apiRoute.HandleFunc(pat.Delete("/users/:user_id/transactions/:transaction_id"), userHandler.DeleteTransaction)

	return mux
}

func Info(w http.ResponseWriter, request *http.Request) {
	type svcInfo struct {
		JSONAPI struct {
			Version string `json:"version,omitempty"`
			Name    string `json:"name,omitempty"`
		} `json:"jsonapi"`
	}

	info := svcInfo{}
	info.JSONAPI.Version = "v1"
	info.JSONAPI.Name = "HRM API"

	w.Write(jsonutil.Marshal(info))
}
