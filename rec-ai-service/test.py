import pytest
from fastapi.testclient import TestClient
from main import app

client = TestClient(app)

@pytest.mark.asyncio
async def test_segment_endpoint_success():
    test_data = [
        {
            "customer": {
                "customer_id": "12345",
                "purchase_history": ["laptop", "mouse"],
                "demographics": {"age": 28, "income_level": "medium"},
                "last_interaction": "2025-10-04"
            },
            "business_type": "ecommerce"
        },
        {
            "customer": {
                "customer_id": "67890",
                "purchase_history": ["book", "pen"],
                "demographics": {"age": 35, "income_level": "low"},
                "last_interaction": "2025-10-03"
            },
            "business_type": "ecommerce"
        }
    ]

    for data in test_data:
        response = client.post("/segment", json=data)
        assert response.status_code == 200, f"Failed with status {response.status_code}: {response.json()}"
        response_json = response.json()
        assert "customer_id" in response_json, f"Missing customer_id: {response_json}"
        assert response_json["customer_id"] == data["customer"]["customer_id"]
        assert "segment" in response_json, f"Missing segment: {response_json}"
        assert "recommended_offer" in response_json, f"Missing recommended_offer: {response_json}"
        assert "insights" in response_json, f"Missing insights: {response_json}"
        assert isinstance(response_json["insights"], list), f"Insights not a list: {response_json}"

@pytest.mark.asyncio
async def test_segment_endpoint_invalid_input():
    response = client.post("/segment", json={
        "customer": {
            "customer_id": "",  # Geçersiz: Boş customer_id
            "purchase_history": ["laptop"],
            "demographics": {"age": 28},
            "last_interaction": "2025-10-04"
        },
        "business_type": "ecommerce"
    })
    assert response.status_code == 422, f"Expected 422, got {response.status_code}: {response.json()}"

@pytest.mark.asyncio
async def test_health_endpoint():
    response = client.get("/health")
    assert response.status_code in [200, 503], f"Unexpected status {response.status_code}: {response.json()}"
    if response.status_code == 200:
        assert response.json()["status"] == "healthy"