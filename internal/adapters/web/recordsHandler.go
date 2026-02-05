package web

import (
	"fmt"
	"log"
	"net/http"
	"pi/internal/adapters/payloads"
	"pi/internal/app/entities"
	"pi/internal/app/interfaces/services"
	"pi/pkg/utils/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

type recordsHandler struct {
	recordsService services.RecordService
}

func NewRecordsHandler(recordsService services.RecordService) Handler {
	return &recordsHandler{recordsService: recordsService}
}

func (h *recordsHandler) Register(router *gin.Engine) {
	records := router.Group("/record")
	records.Use(middleware.AuthMiddleware())
	{
		records.POST("/add", h.AddRecordHandler)
		records.POST("/export", h.ExportRecordHandler)
	}
}

func (h *recordsHandler) AddRecordHandler(c *gin.Context) {
	var record payloads.RecordsPayload

	if err := c.BindJSON(&record); err != nil {
		log.Printf("error while binding record payload: %v\n", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad json body",
		})

		return
	}

	if record.ExpectedReturnDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "expected return date is required (YYYY-MM-DD)"})
		return
	}
	expectedReturnDate, parseErr := time.Parse("2006-01-02", record.ExpectedReturnDate)
	if parseErr != nil {
		log.Printf("error parsing expected_return_date %q: %v\n", record.ExpectedReturnDate, parseErr)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid date format (expected YYYY-MM-DD)"})
		return
	}

	workerID := record.WorkerId
	workerName := record.WorkerName
	departmentID := record.DepartmentId
	departmentName := record.DepartmentName
	if workerID == "" || workerName == "" || departmentID == "" || departmentName == "" {
		// fallback to logged-in user
		workerID = c.GetString("worker_id")
		workerName = c.GetString("worker_name")
		departmentID = c.GetString("department_id")
		departmentName = c.GetString("department_name")
	}
	if len(workerID) == 0 || len(workerName) == 0 || len(departmentID) == 0 || len(departmentName) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "select department and worker, or ensure you are logged in",
		})
		return
	}

	newRecord, err := h.recordsService.CreateRecord(record.EquipmentId, workerID, workerName, departmentID, departmentName, expectedReturnDate, record.Status)
	if err != nil {
		log.Printf("error creating record: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "new record was added",
		"record":  newRecord,
	})
}

func (h *recordsHandler) ExportRecordHandler(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.ParseInLocation("2006-01-02", startDateStr, time.UTC)
		if err != nil {
			log.Printf("Ошибка парсинга даты начала: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты начала"})
			return
		}
	} else {
		startDate = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	if endDateStr != "" {
		endDate, err = time.ParseInLocation("2006-01-02", endDateStr, time.UTC)
		if err != nil {
			log.Printf("Ошибка парсинга даты окончания: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты окончания"})
			return
		}
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	} else {
		endDate = time.Now().UTC().Add(24 * time.Hour) // включая весь сегодняшний день
	}

	records, err := h.recordsService.GetRecordsUpTo(startDate, endDate)
	if err != nil {
		log.Printf("Ошибка получения записей: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении данных"})
		return
	}

	pdf := generateRecordsPDF(records, startDate, endDate)

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=records_%s_to_%s.pdf",
		startDate.Format("2006-01-02"), endDate.Format("2006-01-02")))
	pdf.Output(c.Writer)
}

func generateRecordsPDF(records []entities.Records, startDate, endDate time.Time) *gofpdf.Fpdf {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(10, 15, 10)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Equipment Issuance Report")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 6, fmt.Sprintf("Period: %s - %s",
		startDate.Format("02.01.2006"),
		endDate.Format("02.01.2006")),
		"", 1, "L", false, 0, "")
	pdf.Ln(5)

	pdf.CellFormat(0, 6, fmt.Sprintf("Total records: %d", len(records)),
		"", 1, "L", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 9)

	widths := []float64{10, 35, 35, 28, 28, 28, 28, 22}

	headers := []string{"#", "Equipment", "Worker", "Department", "Issued", "Returned", "Expected", "Status"}

	for i, header := range headers {
		pdf.CellFormat(widths[i], 8, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 8)
	for i, record := range records {
		equipName := record.EquipmentName
		if equipName == "" {
			equipName = record.EquipmentId
		}
		pdf.CellFormat(widths[0], 7, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[1], 7, equipName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], 7, record.WorkerName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[3], 7, record.DepartmentName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[4], 7, record.IssuedAt.Format("02.01.2006"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[5], 7, formatDate(record.ReturnedAt), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[6], 7, record.ExpectedReturnDate.Format("02.01.2006"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[7], 7, formatStatusEn(record.Status), "1", 0, "C", false, 0, "")
		pdf.Ln(7)
	}

	pdf.SetY(-15)
	pdf.SetFont("Arial", "I", 8)
	pdf.CellFormat(0, 10, fmt.Sprintf("Generated: %s", time.Now().Format("02.01.2006 15:04")),
		"T", 0, "R", false, 0, "")

	return pdf
}

func formatDate(t time.Time) string {
	zeroTime := time.Time{}
	if t == zeroTime {
		return "-"
	}
	return t.Format("02.01.2006")
}

func formatStatusEn(status string) string {
	statusMap := map[string]string{
		"issued":   "Issued",
		"returned": "Returned",
		"overdue":  "Overdue",
	}
	if val, ok := statusMap[status]; ok {
		return val
	}
	return status
}
