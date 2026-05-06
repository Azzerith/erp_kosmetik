package handler

import (
	"net/http"

	"erp-cosmetics-backend/internal/service"
	"erp-cosmetics-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ReportHandler struct {
	reportService service.ReportService
	logger        *zap.Logger
}

func NewReportHandler(reportService service.ReportService, logger *zap.Logger) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
		logger:        logger,
	}
}

func (h *ReportHandler) GetSalesReport(c *gin.Context) {
	var req service.SalesReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}
	
	report, err := h.reportService.GetSalesReport(c.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate sales report", err)
		return
	}
	
	utils.SuccessResponse(c, report)
}

func (h *ReportHandler) GetInventoryReport(c *gin.Context) {
	report, err := h.reportService.GetInventoryReport(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate inventory report", err)
		return
	}
	
	utils.SuccessResponse(c, report)
}

func (h *ReportHandler) GetCustomerReport(c *gin.Context) {
	report, err := h.reportService.GetCustomerReport(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate customer report", err)
		return
	}
	
	utils.SuccessResponse(c, report)
}

func (h *ReportHandler) GetDashboardSummary(c *gin.Context) {
	summary, err := h.reportService.GetDashboardSummary(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch dashboard summary", err)
		return
	}
	
	utils.SuccessResponse(c, summary)
}