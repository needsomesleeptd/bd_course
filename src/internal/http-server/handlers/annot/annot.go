package annot_handler

import (
	service "annotater/internal/bl/annotationService"
	response "annotater/internal/lib/api"
	"annotater/internal/middleware/auth_middleware"
	"annotater/internal/models"
	models_dto "annotater/internal/models/dto"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/render"
)

var (
	ErrDecodingRequest = errors.New("broken request")
	ErrGettingFile     = errors.New("no file got")
	ErrGettingBbs      = errors.New("no bbs got")
	ErrAddingAnnot     = errors.New("error adding annotattion")
	ErrGettingAnnot    = errors.New("error getting annotattion")
	ErrDeletingAnnot   = errors.New("error deleting annotattion")
)

const (
	AnnotFileFieldName = "annotFile"
	JsonBbsFieldName   = "jsonBbs"
)

type RequestAddAnnot struct { //other data is a file
	ErrorBB    []float32 `json:"error_bb"`
	ClassLabel uint64    `json:"class_label"`
}

type RequestID struct {
	ID uint64 `json:"id"`
}

type RequestUpdate struct {
	ID         uint64    `json:"id"`
	ErrorBB    []float32 `json:"error_bb"`
	ClassLabel uint64    `json:"class_label"`
	TypeLabel  int       `json:"type_label"`
	IsValid    bool      `json:"is_valid"`
}

type ResponseGetAnnot struct {
	response.Response
	models_dto.Markup
}

type ResponseGetAnnots struct {
	response.Response
	Markups []models_dto.Markup
}

// @Summary Add an annotation
// @Description Adds an annotation to a specific document
// @Tags Annotation
// @Security ApiKeyAuth
// @Accept mpfd
// @Produce json
// @Param annotFile formData file true "PNG image to add"
// @Param annotParams body RequestAddAnnot true "Annotation parameters (bboxes and class)"
// @Success 200 {object} response.Response "OK"
// @Failure 200 {object} response.Response
// @Router /annotations/add [post]
func AddAnnot(annotService service.IAnotattionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestAddAnnot
		var pageData []byte
		userID := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		file, _, err := r.FormFile(AnnotFileFieldName)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error())) //TODO:: add logging here
			fmt.Print(err.Error())
			return
		}
		pageData, err = io.ReadAll(file)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error())) //TODO:: add logging here
			return
		}
		bbsString := r.FormValue(JsonBbsFieldName)

		err = json.Unmarshal([]byte(bbsString), &req)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error())) //TODO:: add logging here
			return
		}
		annot := models.Markup{
			PageData:   pageData,
			ErrorBB:    req.ErrorBB,
			ClassLabel: req.ClassLabel,
			CreatorID:  userID,
		}
		err = annotService.AddAnottation(&annot)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			fmt.Print(err.Error())
			return
		}
		render.JSON(w, r, response.OK())
	}

}

// @Summary Get a specific annotation
// @Description Get the specific annotation by ID
// @Tags Annotations
// @Accept json
// @Produce json
// @Param id path string true "The ID of the annotation to get"
// @Success 200 {object} ResponseGetAnnot
// @Failure 404 {object} response.Response
// @Router /annot/get [get]
func GetAnnot(annotService service.IAnotattionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestID
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(ErrDecodingRequest.Error())) //TODO:: add logging here
			return
		}
		var markUp *models.Markup
		markUp, err = annotService.GetAnottationByID(req.ID)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			return
		}
		resp := ResponseGetAnnot{Markup: *models_dto.ToDtoMarkup(*markUp), Response: response.OK()}
		render.JSON(w, r, resp)
	}
}

// @Summary Getting all anottattions from a database
// @Description Getting all anottattions from a database, works only when there are not a lot of annotattions has no paging
// @Tags Annotations
// @Accept json
// @Produce json
// @Success 200 {object} ResponseGetAnnots
// @Failure 200 {object} response.Response
// @Router /annot/getsAll [get]
func GetAllAnnots(annotService service.IAnotattionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		markUps, err := annotService.GetAllAnottations()
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			return
		}
		resp := ResponseGetAnnots{Markups: models_dto.ToDtoMarkupSlice(markUps), Response: response.OK()}
		render.JSON(w, r, resp)
	}
}

// @Summary Getting all anotattions created by a user
// @Description Getting all anottattions from a database, which were created by currently logged user
// @Tags Annotations
// @Accept json
// @Produce json
// @Success 200 {object} ResponseGetAnnots
// @Failure 200 {object} response.Response
// @Router /annot/creatorID [get]
func GetAnnotsByUserID(annotService service.IAnotattionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID, ok := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		if !ok {
			render.JSON(w, r, response.Error(ErrDecodingRequest.Error())) //TODO:: add logging here
			return
		}

		markUps, err := annotService.GetAnottationByUserID(userID)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			return
		}
		resp := ResponseGetAnnots{Markups: models_dto.ToDtoMarkupSlice(markUps), Response: response.OK()}
		render.JSON(w, r, resp)
	}
}

func DeleteAnnot(annotService service.IAnotattionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestID
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(ErrDecodingRequest.Error())) //TODO:: add logging here
			return
		}

		err = annotService.DeleteAnotattion(req.ID)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			return
		}
		render.JSON(w, r, response.OK())
	}
}

// @Summary Modifies markup and marks it as checked
// @Description Updating chosed markup and marking it as checked as well as setting creator_id to the current logged user
// @Tags Annotations
// @Accept json
// @Produce json
// @Param NewMarkupParams body RequestUpdate true "data to fix broken markup"
// @Success 200 {object} response.Response
// @Failure 200 {object} response.Response
// @Router /annot/check [post]
func Check(annotService service.IAnotattionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestUpdate
		userID := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(ErrDecodingRequest.Error())) //TODO:: add logging here
			return
		}
		markup := models.Markup{
			ID:         req.ID,
			ClassLabel: req.ClassLabel,
			TypeLabel:  req.TypeLabel,
			ErrorBB:    req.ErrorBB,
			CreatorID:  userID,
			IsValid:    req.IsValid,
		}
		err = annotService.CheckAnotattion(&markup)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			return
		}
		render.JSON(w, r, response.OK())
	}
}
