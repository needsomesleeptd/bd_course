package user_handler

import (
	service "annotater/internal/bl/userService"
	response "annotater/internal/lib/api"
	"annotater/internal/models"
	models_dto "annotater/internal/models/dto"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

var (
	ErrChangingRole    = errors.New("error changing role")
	ErrDecodingJson    = errors.New("broken request")
	ErrGettingAllUsers = errors.New("error getting all users")
)

type RequestChangeRole struct {
	Login   string      `json:"login"`
	ReqRole models.Role `json:"req_role"`
}

type ResponseGetAllUsers struct {
	response.Response
	Users []models_dto.User `json:"users"`
}

type UserHandler struct {
	logger      *slog.Logger
	userService service.IUserService
}

func NewDocumentHandler(logSrc *slog.Logger, serv service.IUserService) UserHandler {
	return UserHandler{
		logger:      logSrc,
		userService: serv,
	}
}

// @Summary Change user role (available if admin)
// @Description Set user role, only available if you are admin
// 0 -- sender, 1 --controller, 2-- admin
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param ChangeUserRoleParams body RequestChangeRole true "data to get the user and change his role"
// @Success 200 {object} response.Response
// @Failure 200 {object} response.Response "Annotation not found"
// @Router /user/role [post]
func (h *UserHandler) ChangeUserPerms() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestChangeRole
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(ErrDecodingJson.Error()))
			h.logger.Error(err.Error())
			return
		}
		err = h.userService.ChangeUserRoleByLogin(req.Login, req.ReqRole)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.logger.Error(err.Error())
			return
		}
		render.JSON(w, r, response.OK())
	}
}

// @Summary Get all users data (available if admin)
// @Description Get all users parametersm available if you are admin
// roles description: 0 -- sender, 1 --controller, 2-- admin
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} ResponseGetAllUsers
// @Failure 200 {object} response.Response "some internal user errror"
// @Router /user/getUsers [get]
func (h *UserHandler) GetAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		users, err := h.userService.GetAllUsers()
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.logger.Error(err.Error())
			return
		}
		usersDTO := models_dto.ToDtoUserSlice(users)
		resp := ResponseGetAllUsers{response.OK(), usersDTO}
		render.JSON(w, r, resp)
	}
}
