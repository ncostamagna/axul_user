package user

import (
	"context"

	"github.com/digitalhouse-dev/dh-kit/response"
)

type (
	StoreReq struct {
		UserName     string `json:"username"`
		FirstName    string `json:"firstname"`
		LastName     string `json:"lastname"`
		Password     string `json:"password"`
		Email        string `json:"email"`
		Phone        string `json:"phone"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Token        string `json:"token"`
	}

	GetAllReq struct {
		ID      []string `json:"id"`
		Preload string   `json:"preload"`
		Limit   int      `json:"limit"`
		Page    int      `json:"page"`
	}

	GetReq struct {
		ID      string `json:"id"`
		Preload string `json:"preload"`
	}

	LoginReq struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}
)

type Controller func(ctx context.Context, request interface{}) (interface{}, error)

//Endpoints struct
type Endpoints struct {
	Get    Controller
	GetAll Controller
	Store  Controller
	Login  Controller
	Update Controller
	Delete Controller
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
		Store:  makeStoreEndpoint(s),
		Login:  makeLoginEndpoint(s),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

func makeGetEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetReq)

		users, err := service.Get(ctx, req.ID, req.Preload)
		if err != nil {
			switch err {
			case NotFound:
				return nil, NotFound
			default:
				return nil, err
			}

		}

		return response.Success("", users, nil, nil), nil
	}
}

func makeGetAllEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAllReq)
		filters := Filters{ID: req.ID}

		users, err := service.GetAll(ctx, filters, 0, 0, req.Preload)
		if err != nil {
			return nil, response.BadRequest(err.Error())
		}

		return response.Success("", users, nil, nil), nil
	}
}

func makeStoreEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(StoreReq)
		user, err := service.Create(ctx, req.UserName, req.FirstName, req.LastName, req.Password, req.Email, req.Phone, req.ClientID, req.ClientSecret, req.Token)
		if err != nil {
			switch err {
			case FieldIsRequired:
				return nil, FieldIsRequired
			default:
				return nil, err
			}
		}

		return response.Success("success", user, nil, nil), nil
	}
}

func makeLoginEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginReq)
		user, err := service.Login(ctx, req.ID, req.Password)
		if err != nil {
			switch err {
			case NotFound:
				return nil, NotFound
			case InvalidAuthentication:
				return nil, InvalidAuthentication
			default:
				return nil, err
			}
		}
		return response.Success("success", user, nil, nil), nil
	}
}

func makeUpdateEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		return response.Success("success", nil, nil, nil), nil
	}
}

func makeDeleteEndpoint(service Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		return response.Success("", nil, nil, nil), nil
	}
}
