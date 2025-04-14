import requests
import random
import time
import hmac
import hashlib
import os
import json
from time import time as timestamp

API_BASE_URL = "http://localhost:8080/api/leaderboard"
API_SECRET_KEY = os.getenv('API_SECRET_KEY', 'top-secret-api-key')

def generate_nonce():
    return os.urandom(16).hex()

def calculate_hmac(message):
    return hmac.new(
        API_SECRET_KEY.encode(),
        message.encode(),
        hashlib.sha256
    ).hexdigest()

# Simulate score submission
def submit_score(user_id):
    score = random.randint(100, 10000)
    current_timestamp = str(int(timestamp()))
    nonce = generate_nonce()
    
    # Prepare request body
    body = json.dumps({"user_id": user_id, "score": score})
    
    # Calculate signature
    message = f"{current_timestamp}:{nonce}:{body}"
    signature = calculate_hmac(message)
    
    # Prepare headers
    headers = {
        'Content-Type': 'application/json',
        'X-Timestamp': current_timestamp,
        'X-Nonce': nonce,
        'X-Signature': signature
    }
    
    res = requests.post(
        f"{API_BASE_URL}/submit",
        headers=headers,
        json={"user_id": user_id, "score": score}
    )
    res.raise_for_status()
    print("successful post: ", res.json())

# Fetch top players
def get_top_players():
    response = requests.get(f"{API_BASE_URL}/top")
    return response.json()

# Fetch user rank
def get_user_rank(user_id):
    response = requests.get(f"{API_BASE_URL}/rank/{user_id}")
    return response.json()

if __name__ == "__main__":
    while True:
        user_id = random.randint(1, 1000000)
        submit_score(user_id)
        print(get_top_players())
        print(get_user_rank(user_id))
        time.sleep(random.uniform(0.5, 2)) # Simulate real user interaction