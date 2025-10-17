import numpy as np
import ollama
import json
import logging
import re
from sklearn.cluster import MiniBatchKMeans
from fastapi import HTTPException
import os

logger = logging.getLogger(__name__)

# MiniBatchKMeans modeli (global, başlatma main.py'den çağrılmaz)
kmeans = MiniBatchKMeans(n_clusters=int(os.getenv("NUM_CLUSTERS", 5)), random_state=42)

def initialize_kmeans():
    """Modeli başlatmak için dummy veri ile ön eğitim yapar."""
    dummy_data = np.array([
        [1, 25, 2],  # purchase_count, age, income_level
        [3, 30, 3],
        [2, 40, 1],
        [5, 35, 2],
        [4, 28, 3]
    ])
    kmeans.fit(dummy_data)
    logger.info("MiniBatchKMeans modeli dummy veri ile başlatıldı.")

initialize_kmeans()

def extract_features(input) -> np.ndarray:
    """Müşteri verilerinden özellik vektörü çıkarır."""
    try:
        purchase_count = len(input.purchase_history)
        age = input.demographics.get("age", 0)
        income_map = {"low": 1, "medium": 2, "high": 3}
        income_level = income_map.get(input.demographics.get("income_level", "medium"), 2)
        logger.info(f"Özellikler çıkarıldı: customer_id={input.customer_id}, features=[{purchase_count}, {age}, {income_level}]")
        return np.array([[purchase_count, age, income_level]])
    except Exception as e:
        logger.error(f"Özellik çıkarma hatası: {str(e)}")
        raise HTTPException(status_code=400, detail=f"Invalid demographics data: {str(e)}")

def calculate_confidence(features: np.ndarray, kmeans: MiniBatchKMeans, segment_id: int) -> float:
    """Veri noktasının kümesine olan confidence'ını hesaplar."""
    try:
        centers = kmeans.cluster_centers_
        own_distance = np.linalg.norm(features - centers[segment_id])
        other_distances = [np.linalg.norm(features - centers[i]) for i in range(len(centers)) if i != segment_id]
        if not other_distances:
            return 0.5
        min_other_distance = min(other_distances)
        confidence = 1 - (own_distance / (own_distance + min_other_distance))
        confidence = max(0.0, min(1.0, confidence))
        logger.info(f"Confidence hesaplandı: segment_id={segment_id}, confidence={confidence:.2f}")
        return confidence
    except Exception as e:
        logger.error(f"Confidence hesaplama hatası: {str(e)}")
        return 0.5

def extract_json_from_response(response: str) -> str:
    """Ollama yanıtından sadece JSON kısmını ayıklar."""
    try:
        json_match = re.search(r'\{[\s\S]*\}', response, re.MULTILINE)
        if json_match:
            json_str = json_match.group(0)
            json.loads(json_str)
            return json_str
        logger.error(f"No valid JSON found in response: {response}")
        return '{"segment": "unknown", "recommended_offer": "Generic offer", "insights": [{"priority": "low", "action": "Review customer data"}]}'
    except json.JSONDecodeError:
        logger.error(f"Invalid JSON format in response: {response}")
        return '{"segment": "unknown", "recommended_offer": "Generic offer", "insights": [{"priority": "low", "action": "Review customer data"}]}'

def validate_llm_response(response: str, model_name: str) -> dict:
    """Ollama yanıtını doğrular ve JSON'a çevirir."""
    try:
        json_str = extract_json_from_response(response)
        llm_response = json.loads(json_str)
        required_keys = {"segment", "recommended_offer", "insights"}
        if not all(key in llm_response for key in required_keys):
            raise ValueError(f"LLM response missing required keys: {required_keys}")
        return llm_response
    except json.JSONDecodeError as e:
        logger.error(f"Invalid JSON: {json_str}")
        return {
            "segment": "unknown",
            "recommended_offer": "Generic offer",
            "insights": [{"priority": "low", "action": "Review customer data"}]
        }
    except ValueError as e:
        logger.error(f"LLM response validation error: {str(e)}")
        return {
            "segment": "unknown",
            "recommended_offer": "Generic offer",
            "insights": [{"priority": "low", "action": "Review customer data"}]
        }