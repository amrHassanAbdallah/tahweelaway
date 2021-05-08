package api

import (
	"encoding/json"
	"github.com/MagalixTechnologies/kube-notifier-service/utils"
	"github.com/amrHassanAbdallah/tahweelaway/persistence"
	"github.com/amrHassanAbdallah/tahweelaway/service"
	"github.com/go-chi/render"
	"net/http"
)

type server struct {
	tahweelawayService *service.TahweelawayService
}

func (s *server) QueryBanks(w http.ResponseWriter, r *http.Request, params QueryBanksParams) {
	panic("implement me")
}

func (u *NewUser) toServiceUser() (*service.User, error) {
	return &service.User{
		Email:    u.Email,
		Name:     u.Name,
		Username: u.Username,
		Password: u.Password,
		Currency: u.Currency,
	}, nil
}
func (s *server) AddUser(w http.ResponseWriter, r *http.Request, params AddUserParams) {
	ctx := r.Context()

	var newUser NewUser
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		HandleError(w, r, &ValidationError{
			Cause:  err,
			Detail: nil,
			Status: 400,
		})
		return
	}
	serviceUser, err := newUser.toServiceUser()
	if err != nil {
		HandleError(w, r, &ValidationError{
			Cause:  err,
			Detail: nil,
			Status: http.StatusBadRequest,
		})
		return
	}
	err = utils.Validator.Struct(serviceUser)
	if err != nil {
		HandleError(w, r, &ValidationError{
			Cause:  err,
			Detail: nil,
			Status: http.StatusBadRequest,
		})
		return
	}
	user, err := s.tahweelawayService.AddUser(ctx, *serviceUser)
	if err != nil {
		switch err.(type) {
		case *persistence.UserConstraintException:
			HandleError(w, r, &service.ServiceError{
				Cause: err,
				Type:  http.StatusConflict,
			})
		default:
			HandleError(w, r, &service.ServiceError{
				Cause: err,
				Type:  http.StatusInternalServerError,
			})
		}
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, MapUserToResponse(user))
	return
}

func (s *server) GetUser(w http.ResponseWriter, r *http.Request, userId string, params GetUserParams) {
	ctx := r.Context()
	user, err := s.tahweelawayService.GetUser(ctx, userId)
	if err != nil {
		switch err.(type) {
		case *persistence.InvalidIDException:
			HandleError(w, r, &service.ServiceError{
				Cause: err,
				Type:  http.StatusConflict,
			})
		default:
			HandleError(w, r, &service.ServiceError{
				Cause: err,
				Type:  http.StatusInternalServerError,
			})
		}
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, MapUserToResponse(user))
	return
}

func MapUserToResponse(user *persistence.User) UserResponse {
	balance := user.Balance.Int64
	return UserResponse{
		Id:      user.ID.String(),
		Balance: int(balance),
		NewUser: NewUser{
			Email:    user.Email,
			Name:     user.Name,
			Password: user.Password,
			Username: user.Username,
			Currency: user.Currency,
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: nil,
	}
}

func NewServer(svc *service.TahweelawayService) ServerInterface {
	return &server{
		tahweelawayService: svc,
	}
}
