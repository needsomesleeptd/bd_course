package document_handler

import (
	service "annotater/internal/bl/documentService"
	response "annotater/internal/lib/api"
	"annotater/internal/middleware/auth_middleware"
	"annotater/internal/models"
	pdf_utils "annotater/internal/pkg/pdfUtils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

var (
	ErrDecodingJson     = errors.New("broken load document request")
	ErrLoadingDocument  = errors.New("error loading document")
	ErrCheckingDocument = errors.New("error checking document")
	ErrGettingID        = errors.New("error generating unqiue document ID")
	ErrInvalidFile      = errors.New("got invalid file")
	ErrGettingFile      = errors.New("error retriving file")
	ErrCreatingReport   = errors.New("error creating document report")
	ErrBrokenRequest    = errors.New("broken request")

	ErrSendingFile      = errors.New("error sending file")
	ErrGettingPageCount = errors.New("error getting pagecount")
)

const (
	FILE_HEADER_KEY = "file"

	ErrGetttingMetaData = "error getting metadataDocuments"
)

type IDocumentHandler interface {
	LoadDocument(documentService service.IDocumentService) http.HandlerFunc
	CheckDocument(documentService service.IDocumentService) http.HandlerFunc
}

type Documenthandler struct {
	logger     *slog.Logger
	docService service.IDocumentService
}

type RequestLoadDocument struct {
	Document []byte `json:"document_data"`
}
type RequestCheckDocument struct {
	Document []byte `json:"document_data"`
}

type RequestPassed struct {
	ID        uuid.UUID `json:"ID"`
	HasPassed bool      `json:"has_passed"`
}

type RequestID struct {
	ID uuid.UUID `json:"ID"`
}

type ResponseCheckDoucment struct {
	Response    response.Response
	Markups     []models.Markup     `json:"markups"`
	MarkupTypes []models.MarkupType `json:"markupTypes"`
}

type ResponseGettingMetaData struct {
	Response          response.Response
	DocumentsMetaData []models.DocumentMetaData `json:"documents_metadata"`
}

type ResponseGetReport struct {
	Response    response.Response
	Markups     []models.Markup     `json:"markups"`
	MarkupTypes []models.MarkupType `json:"markupTypes"`
}

func NewDocumentHandler(logSrc *slog.Logger, serv service.IDocumentService) Documenthandler {
	return Documenthandler{
		logger:     logSrc,
		docService: serv,
	}
}

func ExtractfileBytesHelper(file multipart.File) ([]byte, error) {

	defer file.Close()

	fileBytes, err := io.ReadAll(file)

	if err != nil {
		return nil, ErrInvalidFile
	}

	return fileBytes, nil

}

func writeBytesIntoResponse(w http.ResponseWriter, data []byte) error {
	w.Header().Set("Content-Type", http.DetectContentType(data))
	w.Header().Set("Content-Length", fmt.Sprintf("%v", len(data)))
	_, err := w.Write(data)
	if err != nil {
		return errors.Join(err, ErrSendingFile)
	}
	return nil

}

// @Summary Get document by ID
// @Description Fetches a document file without metadata by its ID
// @Security ApiKeyAuth
// @Tags Document
// @Accept json
// @Produce application/pdf,json
// @Param Document body RequestID true  "Document ID"
// @Success 200 {object} []byte
// @Failure 200 {object} response.Response
// @Router /document/getDocument [post]
func (h *Documenthandler) GetDocumentByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestID
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(ErrBrokenRequest.Error()))
			return
		}

		document, err := h.docService.GetDocumentByID(req.ID)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.logger.Error(err.Error())
			return
		}
		err = writeBytesIntoResponse(w, document.DocumentBytes)
		if err != nil {
			render.JSON(w, r, response.Error(ErrSendingFile.Error()))
			h.logger.Error(err.Error())
			return
		}
		render.JSON(w, r, response.OK())
	}
}

// @Summary Get report by ID
// @Description Fetches a report file without metadata by its ID
// @Security ApiKeyAuth
// @Tags Document
// @Accept json
// @Produce application/pdf,json
// @Param ReportID body RequestID true  "Report ID"
// @Success 200 {object} []byte
// @Failure 200 {object} response.Response
// @Router /document/getReport [post]
func (h *Documenthandler) GetReportByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestID
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(ErrBrokenRequest.Error()))
			return
		}

		report, err := h.docService.GetReportByID(req.ID)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.logger.Error(err.Error())
			return
		}

		err = writeBytesIntoResponse(w, report.ReportData)
		if err != nil {
			render.JSON(w, r, response.Error(ErrSendingFile.Error()))
			h.logger.Error(err.Error())
			return
		}
		render.JSON(w, r, response.OK())
	}
}

// @Summary Get document meta data
// @Description Gets all document metadata which was created by this user
// @Security ApiKeyAuth
// @Tags Document
// @Accept json
// @Produce json
// @Success 200 {object} ResponseGettingMetaData
// @Failure 200 {object} response.Response
// @Router /document/getDocumentsMeta [get]
func (h *Documenthandler) GetDocumentsMetaData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		documentsMetaData, err := h.docService.GetDocumentsByCreatorID(userID)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.logger.Error(err.Error())
			return
		}
		resp := ResponseGettingMetaData{Response: response.OK(), DocumentsMetaData: documentsMetaData}
		render.JSON(w, r, resp)
	}
}

// @Summary Gets an unchecked file and marks it as being checked
// @Description Gets a document and it's metaData and returns it to the user, inly available for controllers and admin
// @Security ApiKeyAuth
// @Tags Document
// @Accept mpfd
// @Produce application/pdf,json
// @Success 200 {object} []byte
// @Failure 200 {object} response.Response
// @Router /document/getDocumentForCheck [get]
func (h *Documenthandler) CreateReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			render.JSON(w, r, response.Error(ErrGettingFile.Error()))
			h.logger.Error(err.Error())
			return
		}
		file, handler, err := r.FormFile(FILE_HEADER_KEY)
		if err != nil {
			render.JSON(w, r, response.Error(ErrGettingFile.Error()))
			h.logger.Error(err.Error())
		}

		var fileBytes []byte
		fileBytes, err = ExtractfileBytesHelper(file)

		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			h.logger.Error(err.Error())
			return
		}

		var pagesCount int
		var checksCount int64
		pagesCount, err = pdf_utils.GetPdfPageCount(fileBytes)
		if err != nil {
			h.logger.Error(errors.Join(err, ErrGettingPageCount).Error())
			pagesCount = -1
		}

		checksCount, err = h.docService.GetDocumentCountByDocumentName(handler.Filename)
		if err != nil {
			h.logger.Error(errors.Join(err, ErrGettingPageCount).Error())
			checksCount = -1
		}
		documentMetaData := models.DocumentMetaData{
			ID:           uuid.New(),
			CreatorID:    userID,
			DocumentName: handler.Filename,
			CreationTime: time.Now(),
			PageCount:    pagesCount,
			ChecksCount:  uint64(checksCount),
		}
		documentData := models.DocumentData{
			DocumentBytes: fileBytes,
			ID:            documentMetaData.ID,
		}

		var report *models.ErrorReport
		report, err = h.docService.LoadDocument(documentMetaData, documentData, userID)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.logger.Error(err.Error())
			return
		}

		w.Header().Set("Content-Type", http.DetectContentType(report.ReportData))
		w.Header().Set("Content-Length", fmt.Sprintf("%v", len(report.ReportData)))
		_, err = w.Write(report.ReportData)
		if err != nil {
			render.JSON(w, r, response.Error(ErrCreatingReport.Error()))
			h.logger.Error(err.Error())
			return
		}

	}
}

// @Summary Mark a document to pass by normocontroller (requires a role to be a normocontroller)
// @Description Set's a field of the document has passed to true and marks it as checked
// @Security ApiKeyAuth
// @Tags Document
// @Accept json
// @Produce json
// @Param HasPassedParams body RequestPassed true "id and bool value of whether the lab was merged"
// @Success 200 {object} response.Response
// @Failure 200 {object} response.Response
// @Router /document/makeDecision [post]
func (h *Documenthandler) MakeDecisionPassed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestPassed

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(models.ErrDecodingRequest.Error()))
			return
		}
		document := models.DocumentMetaData{HasPassed: req.HasPassed, CheckedStatus: models.WasChecked}
		err = h.docService.UpdateDocumentData(req.ID, document)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.logger.Error(err.Error())
			return
		}
		render.JSON(w, r, response.OK())
	}
}

// @Summary Create an error report by given document
// @Description Gets a document, saves it with metadata on the system, then creates a report,
// saves it in the system and gives it back to the sender
// @Security ApiKeyAuth
// @Tags Document
// @Accept multipart/form-data
// @Produce application/pdf,json
// @Success 200 {object} []byte
// @Failure 200 {object} response.Response
// @Router /document/report [post]
func (h *Documenthandler) GetDocumentForCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		docMeta, doc, err := h.docService.GetDocumentForChecking()
		if err != nil {
			render.JSON(w, r, response.Error(ErrGettingFile.Error()))
			h.logger.Error(err.Error())
			return
		}
		docMetaMarshalled, err := json.Marshal(docMeta)
		if err != nil {
			render.JSON(w, r, response.Error("error marshalling  docMetaData"))
			h.logger.Error(err.Error())
			return
		}
		var requestBody bytes.Buffer
		mpWriter := multipart.NewWriter(&requestBody)

		jsonPart, err := mpWriter.CreateFormField("metadata")
		if err != nil {
			render.JSON(w, r, response.Error("error creating formField"))
			h.logger.Error(err.Error())
			return
		}
		_, err = jsonPart.Write(docMetaMarshalled)
		if err != nil {
			render.JSON(w, r, response.Error("error writing marshlled data"))
			h.logger.Error(err.Error())
			return
		}
		filePart, err := mpWriter.CreateFormFile("file", docMeta.DocumentName)
		if err != nil {
			render.JSON(w, r, response.Error("error getting filePart"))
			h.logger.Error(err.Error())
		}

		// Assuming fileData contains the file data in memory
		_, err = filePart.Write(doc.DocumentBytes)
		if err != nil {
			render.JSON(w, r, response.Error("error getting filePart"))
			h.logger.Error(err.Error())
		}
		mpWriter.Close()
		w.Header().Set("Content-Type", mpWriter.FormDataContentType())
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(requestBody.Bytes())
		if err != nil {
			render.JSON(w, r, response.Error("error sending response"))
			h.logger.Error(err.Error())
		}

	}
}

func (h *Documenthandler) MakeVerdict() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			render.JSON(w, r, response.Error(ErrGettingFile.Error()))
			h.logger.Error(err.Error())
			return
		}
		file, handler, err := r.FormFile(FILE_HEADER_KEY)
		if err != nil {
			render.JSON(w, r, response.Error(ErrGettingFile.Error()))
			h.logger.Error(err.Error())
		}

		var fileBytes []byte
		fileBytes, err = ExtractfileBytesHelper(file)

		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			h.logger.Error(err.Error())
			return
		}

		var pagesCount int
		pagesCount, err = pdf_utils.GetPdfPageCount(fileBytes)

		if err != nil {
			h.logger.Error(errors.Join(err, ErrGettingPageCount).Error())
			pagesCount = -1
		}

		documentMetaData := models.DocumentMetaData{
			ID:           uuid.New(),
			CreatorID:    userID,
			DocumentName: handler.Filename,
			CreationTime: time.Now(),
			PageCount:    pagesCount,
		}
		documentData := models.DocumentData{
			DocumentBytes: fileBytes,
			ID:            documentMetaData.ID,
		}

		var report *models.ErrorReport
		report, err = h.docService.LoadDocument(documentMetaData, documentData, userID)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.logger.Error(err.Error())
			return
		}

		w.Header().Set("Content-Type", http.DetectContentType(report.ReportData))
		w.Header().Set("Content-Length", fmt.Sprintf("%v", len(report.ReportData)))
		_, err = w.Write(report.ReportData)
		if err != nil {
			render.JSON(w, r, response.Error(ErrCreatingReport.Error()))
			h.logger.Error(err.Error())
			return
		}

	}
}

func (h *Documenthandler) SaveDocument() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			render.JSON(w, r, response.Error(ErrGettingFile.Error()))
			h.logger.Error(err.Error())
			return
		}
		file, handler, err := r.FormFile(FILE_HEADER_KEY)
		if err != nil {
			render.JSON(w, r, response.Error(ErrGettingFile.Error()))
			h.logger.Error(err.Error())
		}

		var fileBytes []byte
		fileBytes, err = ExtractfileBytesHelper(file)

		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			h.logger.Error(err.Error())
			return
		}

		var pagesCount int
		var checksCount int64
		pagesCount, err = pdf_utils.GetPdfPageCount(fileBytes)
		if err != nil {
			h.logger.Error(errors.Join(err, ErrGettingPageCount).Error())
			pagesCount = -1
		}

		checksCount, err = h.docService.GetDocumentCountByDocumentName(handler.Filename)
		if err != nil {
			h.logger.Error(errors.Join(err, ErrGettingPageCount).Error())
			checksCount = -1
		}
		documentMetaData := models.DocumentMetaData{
			ID:           uuid.New(),
			CreatorID:    userID,
			DocumentName: handler.Filename,
			CreationTime: time.Now(),
			PageCount:    pagesCount,
			ChecksCount:  uint64(checksCount),
		}

		err = h.docService.SaveMetaData(documentMetaData)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.logger.Error(err.Error())
			return
		}
		resp := response.OK()
		render.JSON(w, r, resp)
	}
}
