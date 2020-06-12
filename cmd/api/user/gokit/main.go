package main

import (
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"gophr.v2/config/configutil"
	"gophr.v2/user/api/v1/http/gokit"
	"gophr.v2/user/repository"
	"gophr.v2/user/service"
	"gophr.v2/user/service/decorators/logging"
	"log"
	"net/http"
)

func main() {
	conf, err := configutil.LoadDefault()
	if err != nil {
		log.Fatal(err)
	}
	repo, _ := repository.Get(conf, repository.FileRepo)
	svc := service.New(repo)

	loggingSvc := logging.Apply(svc)

	getUserByIDHandler := httptransport.NewServer(
		gokit.MakeGetUserIDEndpoint(loggingSvc),
		gokit.DecodeGetByUserIDRequest,
		gokit.EncodeResponse,
	)

	getAllUserHandler := httptransport.NewServer(
		gokit.MakeGetAllEndpoint(loggingSvc),
		gokit.DecodeGetAllRequest,
		gokit.EncodeResponse,
	)

	registerHandler := httptransport.NewServer(
		gokit.MakeRegisterEndpoint(loggingSvc),
		gokit.DecodeRegisterRequest,
		gokit.EncodeResponse,
	)

	updateHandler := httptransport.NewServer(
		gokit.MakeUpdateEndpoint(loggingSvc),
		gokit.DecodeUpdateRequest,
		gokit.EncodeResponse,
	)

	deleteHandler := httptransport.NewServer(
		gokit.MakeDeleteEndpoint(loggingSvc),
		gokit.DecodeDeleteRequest,
		gokit.EncodeResponse,
	)

	r := mux.NewRouter()
	r.Handle("/user/{id}", getUserByIDHandler).Methods(http.MethodGet)
	r.Handle("/user", getAllUserHandler).Methods(http.MethodGet)
	r.Handle("/user/new", registerHandler).Methods(http.MethodPut)
	r.Handle("/user/update", updateHandler).Methods(http.MethodPost)
	r.Handle("/user/{id}", deleteHandler).Methods(http.MethodDelete)

	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatal(err)
	}
}
