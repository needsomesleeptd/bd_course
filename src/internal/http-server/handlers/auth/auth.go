package auth_handler

import (
	auth_service "annotater/internal/bl/auth"
	response "annotater/internal/lib/api"
	"annotater/internal/models"
	models_dto "annotater/internal/models/dto"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

var (
	ErrDecodingJson = errors.New("broken request")
	ErrInternalServ = errors.New("internal server error")
)

const (
	COOKIE_NAME = "auth_jwt"
)

type RequestSignUp struct {
	User models_dto.User `json:"user"`
}
type RequestSignIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ResponseSignIn struct {
	Response response.Response
	Jwt      string `json:"jwt,omitempty"`
}

// @Summary SignUp as a user
// @Description creates a new user with given login and password the role on creation is sender
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body RequestSignUp true "User login and password details"
// @Success 200 {object} response.Response "OK"
// @Failure 200 {object} response.Response
// @Router /user/SignUp [post]
func SignUp(authService auth_service.IAuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestSignUp
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(ErrDecodingJson.Error())) //TODO:: add logging here
			return
		}
		req.User.Role = models.Sender
		candidate := models_dto.FromDtoUser(&req.User)
		err = authService.SignUp(&candidate)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			return
		}

		render.JSON(w, r, response.OK())
	}
}

// @Summary SignIn an existing user
// @Description Allows an existing user to sign in using their login and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param SignInParams body RequestSignIn true "Login and password details"
// @Success 200 {object} ResponseSignIn "OK"
// @Failure 200 {object} ResponseSignIn
// @Router /user/SignIn [post]
func SignIn(authService auth_service.IAuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestSignIn
		var tokenStr string
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, ResponseSignIn{Response: response.Error(ErrDecodingJson.Error())})
			return
		}
		candidate := models.User{Login: req.Login, Password: req.Password}
		tokenStr, err = authService.SignIn(&candidate)
		if err != nil {
			render.JSON(w, r, ResponseSignIn{Response: response.Error(models.GetUserError(err).Error())})
			return
		}

		resp := ResponseSignIn{Response: response.OK(), Jwt: tokenStr}

		render.JSON(w, r, resp)
	}
}
