import websockets
import asyncio
import psutil
import time


async def ws_client():
    url = "ws://127.0.0.1:9000/injest"

    # Connect to the server
    async with websockets.connect(url) as ws:
        print("WebSocket: Client Connected.")

        # Send memory usage every second
        while True:
            time.sleep(1)
            await ws.send(str(psutil.virtual_memory().percent))


# Start the connection
asyncio.run(ws_client())
