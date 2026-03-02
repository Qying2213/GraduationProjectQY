import requests
import json
import base64

BASE_URL = "http://localhost:8080/api/v1"


def test_update_profile():
    # 1. Login
    print("Logging in...")
    resp = requests.post(
        f"{BASE_URL}/login", json={"username": "admin", "password": "admin123"}
    )
    if resp.status_code != 200:
        print(f"Login failed: {resp.status_code} {resp.text}")
        return

    token = resp.json()["data"]["token"]
    headers = {"Authorization": f"Bearer {token}"}
    print("Login successful. Token obtained.")

    # 2. Update Profile with small avatar string
    print("\nTest 1: Updating profile with small string...")
    data = {"avatar": "small_string"}
    resp = requests.put(f"{BASE_URL}/profile", json=data, headers=headers)
    print(f"Status: {resp.status_code}")
    print(f"Response: {resp.text}")

    # 3. Update Profile with large Base64 string (simulating image)
    print("\nTest 2: Updating profile with large Base64 string (10KB)...")
    large_str = "data:image/png;base64," + "A" * 10240
    data = {"avatar": large_str}
    resp = requests.put(f"{BASE_URL}/profile", json=data, headers=headers)
    print(f"Status: {resp.status_code}")
    if resp.status_code != 200:
        print(f"Response: {resp.text}")
    else:
        print("Success (response truncated)")

    # 4. Update Profile with extra fields
    print("\nTest 3: Updating profile with extra fields...")
    data = {"avatar": "test", "extra": "field"}
    resp = requests.put(f"{BASE_URL}/profile", json=data, headers=headers)
    print(f"Status: {resp.status_code}")
    print(f"Response: {resp.text}")


if __name__ == "__main__":
    test_update_profile()
