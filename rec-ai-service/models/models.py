from pydantic import BaseModel, ConfigDict, Field
from typing import List, Dict

class CustomerInput(BaseModel):
    customer_id: str = Field(..., min_length=1, description="Müşteri kimliği, boş olamaz")
    purchase_history: List[str] = Field(..., min_items=1, description="En az bir satın alma geçmişi gerekli")
    demographics: Dict = Field(..., description="Demografik bilgiler, yaş ve gelir seviyesi içermeli")
    last_interaction: str = Field(..., description="Son etkileşim tarihi, ISO formatında")

    model_config = ConfigDict(
        json_schema_extra={
            "example": {
                "customer_id": "12345",
                "purchase_history": ["laptop", "mouse", "headphones"],
                "demographics": {"age": 28, "location": "Istanbul", "income_level": "medium"},
                "last_interaction": "2025-10-04"
            }
        }
    )

class SegmentRequest(BaseModel):
    customer: CustomerInput
    business_type: str = Field(..., min_length=1, description="İşletme türü, boş olamaz")