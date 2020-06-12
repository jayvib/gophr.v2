package gokit

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"gophr.v2/user"
)

func MakeGetUserIDEndpoint(svc user.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getByUserIDRequest)
		usr, err := svc.GetByUserID(ctx, req.UserID)
		if err != nil {
			return getByUserIDResponse{Error: err.Error()}, nil
		}
		return getByUserIDResponse{
			User: usr,
		}, nil
	}
}

func MakeGetAllEndpoint(svc user.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getAllRequest)
		usrs, next, err := svc.GetAll(ctx, req.Cursor, req.Number)
		if err != nil {
			return getAllResponse{Error: err.Error()}, nil
		}
		return getAllResponse{
			Users:      usrs,
			NextCursor: next,
		}, nil
	}
}

func MakeRegisterEndpoint(svc user.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(registerRequest)
		err := svc.Register(ctx, req.User)
		if err != nil {
			return registerResponse{Error: err.Error()}, nil
		}
		return registerResponse{}, nil
	}
}

func MakeUpdateEndpoint(svc user.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateRequest)
		err := svc.Update(ctx, req.User)
		if err != nil {
			return updateResponse{Error: err.Error()}, nil
		}
		return updateResponse{}, nil
	}
}

func MakeDeleteEndpoint(svc user.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteRequest)
		err := svc.Delete(ctx, req.ID)
		if err != nil {
			return deleteResponse{Error: err.Error()}, nil
		}
		return deleteResponse{}, nil
	}
}
