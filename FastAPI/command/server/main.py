from fastapi import FastAPI, WebSocket, WebSocketDisconnect
import asyncio
from random import randint

app = FastAPI()

@app.get("/")
def serve_home():
    return "this is websocket server"

@app.websocket("/ws")
async def serve_ws(ws: WebSocket):
    await ws.accept()
    
    async def send():
        try:
            for i in range(1, randint(5, 10)):
                await asyncio.sleep(2)
                await ws.send_text(f"message from server {i}")
        except (WebSocketDisconnect, RuntimeError, asyncio.CancelledError):
            print("send切断！")
        finally:
            await ws.close()
    
    async def recv():
        try:
            while True:
                data = await ws.receive_text()
                print(data)
        except (WebSocketDisconnect, RuntimeError, asyncio.CancelledError):
            print("recv切断！")
    
    
    send_task = asyncio.create_task(send())
    recv_task = asyncio.create_task(recv())
    _, pending = await asyncio.wait(
        [send_task, recv_task],
        return_when=asyncio.FIRST_COMPLETED
    )
    for task in pending:
        task.cancel()
        try:
            await task
        except WebSocketDisconnect:
            pass
        except asyncio.CancelledError:
            pass