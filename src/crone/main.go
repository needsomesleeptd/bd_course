package main

import (
	nn_adapter "annotater/internal/bl/NN/NNAdapter"
	nn_model_handler "annotater/internal/bl/NN/NNAdapter/NNmodelhandler"
	repo_adapter "annotater/internal/bl/annotationService/annotattionRepo/anotattionRepoAdapter"
	annot_type_repo_adapter "annotater/internal/bl/anotattionTypeService/anottationTypeRepo/anotattionTypeRepoAdapter"
	document_service "annotater/internal/bl/documentService"
	service "annotater/internal/bl/documentService"
	doc_data_repo_adapter "annotater/internal/bl/documentService/documentDataRepo/documentDataRepo"
	document_repo_adapter "annotater/internal/bl/documentService/documentMetaDataRepo/documentMetaDataRepoAdapter"
	rep_data_repo_adapter "annotater/internal/bl/documentService/reportDataRepo/reportDataRepoAdapter"
	rep_creator_service "annotater/internal/bl/reportCreatorService"
	report_creator "annotater/internal/bl/reportCreatorService/reportCreator"
	models_da "annotater/internal/models/modelsDA"
	"fmt"
	"os"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	CONN_POSTGRES_STR    = "host=localhost user=andrew password=1 database=lab01db port=5432" //TODO:: export through parameters
	POSTGRES_CFG         = postgres.Config{DSN: CONN_POSTGRES_STR}
	MODEL_ROUTE          = "http://0.0.0.0:5000/pred"
	REPORTS_CREATOR_PATH = "../../../reportsCreator"
	REPORTS_PATH         = "../../../reports"
	DOCUMENTS_PATH       = "../../../documents"
	DOCUMENTS_EXT        = ".pdf"
	REPORTS_EXT          = ".pdf"
)
var (
	db   *gorm.DB
	doc  *models_da.DocumentQueue
	serv service.IDocumentService
)

func runTasks() error {
	fmt.Println("Started running tasks")
	var documentQueueDa []models_da.DocumentQueue
	err := db.Raw("SELECT * FROM getTasksReady()").Scan(&documentQueueDa).Error
	if len(documentQueueDa) == 0 {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "Error in getting not checked tasks")
	}
	fmt.Printf("Current running tasks %d\n", len(documentQueueDa))
	for _, doc := range documentQueueDa {
		go runTask(db, &doc, serv)
	}
	return nil
}

func runTask(db *gorm.DB, doc *models_da.DocumentQueue, serv service.IDocumentService) error {
	_, errWhileCreatingReport := serv.CreateReport(doc.DocID)
	if errWhileCreatingReport != nil {
		doc.Status = models_da.Error
		err := db.Model(&models_da.DocumentQueue{}).Where("doc_id", doc.DocID).Updates(*doc).Error
		if err != nil {
			fmt.Println(fmt.Errorf("error checking failed_task:%w-%w", errWhileCreatingReport, err).Error())
			return fmt.Errorf("error checking failed_task:%w-%w", errWhileCreatingReport, err)
		} else {
			doc.Status = models_da.Finished
			err := db.Model(&models_da.DocumentQueue{}).Where("doc_id", doc.DocID).Updates(*doc).Error
			if err != nil {
				fmt.Println(fmt.Errorf("Error checking success_task:%w", &err).Error())
				return fmt.Errorf("Error checking success_task:%w", &err)
			}
		}
	}
	fmt.Printf("Successfully served document %v\n", doc.DocID)
	return nil
}

func main() {
	var err error
	db, err = gorm.Open(postgres.New(POSTGRES_CFG), &gorm.Config{TranslateError: true})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	/*err = db.AutoMigrate(&models_da.DocumentQueue{})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}*/
	modelhandler := nn_model_handler.NewHttpModelHandler(MODEL_ROUTE)
	model := nn_adapter.NewDetectionModel(modelhandler)
	annotTypeRepo := annot_type_repo_adapter.NewAnotattionTypeRepositoryAdapter(db)
	reportCreator := report_creator.NewPDFReportCreator(REPORTS_CREATOR_PATH)
	annotRepo := repo_adapter.NewAnotattionRepositoryAdapter(db)
	reportCreatorService := rep_creator_service.NewDocumentService(model, annotTypeRepo, reportCreator, annotRepo)

	documentStorage := doc_data_repo_adapter.NewDocumentRepositoryAdapter(DOCUMENTS_PATH, DOCUMENTS_EXT)

	reportStorage := rep_data_repo_adapter.NewDocumentRepositoryAdapter(REPORTS_PATH, REPORTS_EXT)

	documentRepo := document_repo_adapter.NewDocumentRepositoryAdapter(db)
	documentService := document_service.NewDocumentService(documentRepo, documentStorage, reportStorage, reportCreatorService)
	serv = documentService
	s, err := gocron.NewScheduler()
	if err != nil {
		fmt.Print(err)
	}

	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			runTasks,
		),
	)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	// each job has a unique id
	fmt.Printf("Crone id:%v\n", j.ID())

	// start the scheduler
	s.Start()

	// block until you are ready to shut down
	select {
	case <-time.After(time.Second * 11):
	}

	// when you're done, shut it down
	err = s.Shutdown()
	if err != nil {
		fmt.Print(err.Error())
	}
}
