import requests

resp = requests.get("http://server:8000")
print(resp.content.decode('utf-8'))