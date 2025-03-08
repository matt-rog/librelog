import websockets
import asyncio
import psutil
import time
import json


async def ws_client():
    url = "ws://127.0.0.1:9000/injest"

    # Connect to the server
    async with websockets.connect(url) as ws:
        print("WebSocket: Client Connected.")

        # Send memory usage every second
        while True:
            time.sleep(1)
            data = json.dumps(
                {
                    "user_id": "fba64ae9-154f-42bb-b040-cf15a3f53ea5",
                    "log_set": "ram",
                    "data": {"perc": psutil.virtual_memory().percent},
                }
            )
            await ws.send(data)


# Start the connection
asyncio.run(ws_client())
