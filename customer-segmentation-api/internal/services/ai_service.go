package services

import (
    "bytes"
    "encoding/json"
    "net/http"

    "github.com/DavutcanJ/customer-segmentation-api/internal/models"
)

func CallAIService(input models.SegmentRequest) (models.SegmentResponse, error) {
    inputJSON, _ := json.Marshal(input)
    resp, err := http.Post("http://localhost:5005/segment", "application/json", bytes.NewBuffer(inputJSON))
    if err != nil {
        return models.SegmentResponse{}, err
    }
    defer resp.Body.Close()

    var aiResponse models.SegmentResponse
    json.NewDecoder(resp.Body).Decode(&aiResponse)
    return aiResponse, nil
}