# AI-assisted code
import websockets
import asyncio
import psutil
import json
import urllib.request

API_URL = "http://127.0.0.1:8080"
INGESTER_URL = "ws://127.0.0.1:9000"
ACCOUNT_NUMBER = ""
PASSWORD = ""


def login(account_number, password):
    req = urllib.request.Request(
        f"{API_URL}/api/login",
        data=json.dumps({"account_number": account_number, "password": password}).encode(),
        headers={"Content-Type": "application/json"},
    )
    with urllib.request.urlopen(req) as resp:
        return json.loads(resp.read())["token"]


async def ws_client():
    token = login(ACCOUNT_NUMBER, PASSWORD)
    url = f"{INGESTER_URL}/ingest?token={token}"

    async with websockets.connect(url) as ws:
        print("connected")
        while True:
            await asyncio.sleep(1)
            data = json.dumps(
                {
                    "log_set": "ram",
                    "data": {"perc": psutil.virtual_memory().percent},
                }
            )
            await ws.send(data)
            resp = await ws.recv()
            print(resp)


asyncio.run(ws_client())
