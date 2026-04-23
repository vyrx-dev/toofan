# Topic: Async IO

import asyncio

async def fetch_data():
    await asyncio.sleep(1)
    return "data received"

async def main():
    result = await fetch_data()
    print(result)

asyncio.run(main())