from fastapi import FastAPI, HTTPException
import logging
import os
from dotenv import load_dotenv
from uvicorn import Config, Server
from models.models import CustomerInput, SegmentRequest
from service.segment import extract_features, calculate_confidence, validate_llm_response, kmeans
import ollama

# Logging konfigurasyonu
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Çevresel değişkenleri yükle
load_dotenv()
MODEL_NAME = os.getenv("OLLAMA_MODEL", "llama3.1")

app = FastAPI(
    title="Customer Segmentation API",
    description="MiniBatchKMeans ve Ollama LLM ile müşteri segmentasyonu",
    version="1.0.0"
)

@app.get("/health", summary="Servis sağlık kontrolü")
async def health_check():
    """Ollama servisinin çalıştığını kontrol eder."""
    try:
        ollama.list()
        return {"status": "healthy", "ollama_model": MODEL_NAME}
    except Exception as e:
        logger.error(f"Health check failed: {str(e)}")
        raise HTTPException(status_code=503, detail="LLM service unavailable")

@app.post("/segment", response_model=dict, summary="Müşteri segmentasyonu yapar")
async def segment_customer(request: SegmentRequest):
    """Müşteri verilerini segmentler ve LLM ile öneriler üretir."""
    try:
        input = request.customer
        business_type = request.business_type
        logger.info(f"Segmentasyon isteği alındı: customer_id={input.customer_id}, business_type={business_type}")

        # Özellik vektörü oluştur ve segmentasyon yap
        features = extract_features(input)
        try:
            kmeans.partial_fit(features)
            segment_id = kmeans.predict(features)[0]
            logger.info(f"Segment ID: {segment_id} for customer_id={input.customer_id}")
        except ValueError as e:
            logger.error(f"KMeans hatası: {str(e)}")
            raise HTTPException(status_code=400, detail=f"KMeans error: {str(e)}")

        # Confidence hesapla
        confidence = calculate_confidence(features, kmeans, segment_id)

        # Ollama ile etiket ve öneri al
        prompt = f"""
        !!! Tüm cevap formatın JSON halinde olmalı !!!
        İşletme türü: {business_type}
        Müşteri ID: {input.customer_id}
        Satın alma geçmişi: {input.purchase_history}
        Demografik bilgiler: {input.demographics}
        Segment ID: {segment_id}

        Bu müşteri için {business_type} sektörüne uygun bir segment etiketi öner (örneğin, 'tech-savvy', 'price-sensitive').
        İşletmeye özel actionable öneriler üret.
        Yanıtı SADECE geçerli JSON formatında dön, örneğin:
        {{"segment": "tech-savvy", "recommended_offer": "15% discount", "insights": [{{"priority": "high", "action": "Email campaign"}}]}}
        Başka hiçbir metin, açıklama veya ```json``` bloğu ekleme, sadece saf JSON! Sadece JSON mesaj döndürmen kabul ediliyor.
        """
        try:
            response = ollama.generate(model=MODEL_NAME, prompt=prompt)
            logger.debug(f"Ollama raw response: {response['response']}")
            llm_response = validate_llm_response(response["response"], MODEL_NAME)
            logger.info(f"Ollama yanıtı alındı: segment={llm_response['segment']}")
        except Exception as e:
            logger.error(f"Ollama hatası: {str(e)}")
            raise HTTPException(status_code=500, detail=f"LLM service error: {str(e)}")

        return {
            "customer_id": input.customer_id,
            "segment": llm_response["segment"],
            "recommended_offer": llm_response["recommended_offer"],
            "confidence": confidence,
            "insights": llm_response["insights"]
        }
    except Exception as e:
        logger.error(f"Segmentasyon hatası: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Internal server error: {str(e)}")

if __name__ == "__main__":
    config = Config(
        app=app,
        host="0.0.0.0",
        port=5005,
        log_level="info"
    )
    server = Server(config)
    server.run()