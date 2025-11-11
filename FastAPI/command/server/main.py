from fastapi import FastAPI

app = FastAPI()

@app.get("/")
def serve_home():
    return "this is websocket server"