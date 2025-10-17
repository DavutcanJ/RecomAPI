package handlers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "time"

    "github.com/DavutcanJ/customer-segmentation-api/internal/models"
    "github.com/gin-gonic/gin"
)

type SegmentHandler struct{}

func NewSegmentHandler() *SegmentHandler {
    return &SegmentHandler{}
}

// @Summary Müşteri segmentasyonu
// @Description Müşteri verilerini segmentler ve Python servisine yönlendirir
// @Tags segment
// @Accept json
// @Produce json
// @Param request body models.SegmentRequest true "Müşteri segmentasyon isteği"
// @Success 200 {object} models.SegmentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /api/segment [post]
func (h *SegmentHandler) SegmentCustomer(c *gin.Context) {
    var req models.SegmentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid input: " + err.Error()})
        return
    }

    // Get Python service URL from environment
    pythonServiceURL := os.Getenv("PYTHON_SERVICE_URL")
    if pythonServiceURL == "" {
        pythonServiceURL = "http://localhost:5005"
    }

    // Log the request for debugging
    fmt.Printf("🔗 Python Service URL: %s\n", pythonServiceURL)
    fmt.Printf("📤 Request data: %+v\n", req)

    // Python servisine istek gönder
    httpClient := &http.Client{
        Timeout: 30 * time.Second,
        Transport: &http.Transport{
            DisableKeepAlives: true,
        },
    }
    
    jsonData, err := json.Marshal(req)
    if err != nil {
        fmt.Printf("❌ JSON Marshal error: %v\n", err)
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to marshal request: " + err.Error()})
        return
    }

    pythonReq, err := http.NewRequest("POST", pythonServiceURL+"/segment", bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Printf("❌ Request creation error: %v\n", err)
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create request: " + err.Error()})
        return
    }
    pythonReq.Header.Set("Content-Type", "application/json")
    pythonReq.Header.Set("Accept", "application/json")

    fmt.Printf("🚀 Sending request to: %s\n", pythonReq.URL.String())
    
    resp, err := httpClient.Do(pythonReq)
    if err != nil {
        fmt.Printf("❌ HTTP request error: %v\n", err)
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{
            Error: fmt.Sprintf("Python service connection error: %v. Make sure Python service is running on %s", err, pythonServiceURL),
        })
        return
    }
    defer resp.Body.Close()

    fmt.Printf("📥 Response status: %d\n", resp.StatusCode)
    fmt.Printf("📥 Response headers: %+v\n", resp.Header)

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Printf("❌ Response read error: %v\n", err)
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to read response: " + err.Error()})
        return
    }

    fmt.Printf("📥 Response body: %s\n", string(body))

    if resp.StatusCode != 200 {
        fmt.Printf("❌ Non-200 status code: %d, body: %s\n", resp.StatusCode, string(body))
        c.JSON(resp.StatusCode, models.ErrorResponse{Error: string(body)})
        return
    }

    var result models.SegmentResponse
    if err := json.Unmarshal(body, &result); err != nil {
        fmt.Printf("❌ JSON Unmarshal error: %v, body: %s\n", err, string(body))
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to parse response: " + err.Error()})
        return
    }

    fmt.Printf("✅ Successful response: %+v\n", result)
    c.JSON(http.StatusOK, result)
}