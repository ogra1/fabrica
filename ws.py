#!/usr/bin/env python3

import asyncio
import functools
import select
import subprocess
import sys
import websockets

async def launcher(cmd, match, websocket, path):
    print('tree: ' + cmd)
    await websocket.send('Building tree: ' + cmd)
    await websocket.send('please wait while a container is spawned ... ')
    if path == match:
        while True:
            proc = subprocess.Popen(['build.py', cmd ],
                    stdout=subprocess.PIPE, bufsize=0 )
            while True:
                line = proc.stdout.readline().decode('UTF-8')
                if not line:
                  break
                await websocket.send(str(line.rstrip()))
                print(line.rstrip())
            break
        sys.exit(0)
    elif path == '/test':
        await websocket.send('Test success !')

def do_websocket(cmd, match):
    loop = asyncio.get_event_loop()
    start_server = websockets.serve(
            functools.partial(launcher, cmd, match),
            "10.0.2.15", 5678)

    print("starting ws")
    loop.run_until_complete(start_server)
    loop.run_forever()

def main():
    do_websocket(sys.argv[1], sys.argv[2])

main()
