#!/bin/python3

import argparse
import asyncio
from json import dumps

from modules.connect import connect
from modules.utilities import utilities

def cmd_flags() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Process command-line arguments")

    parser.add_argument("-lo", "--list-online", action='store_true', help="Specify if the program needs to list online devices.")
    parser.add_argument("-rc", "--run", action='store_true', help="Make the program run a command.")
    parser.add_argument("--command", type=str, help="Specify the actual command that is going to run.")
    parser.add_argument("--nodeids", type=str, help="Specify which nodes the command is going to be run on.")

    parser.add_argument("-i", "--indent", action='store_true', help="Specify whether the output needs to be indented.")

    return parser.parse_args()

async def main():
    args = cmd_flags()
    credentials = utilities.load_config()
    session = await connect.connect(credentials["hostname"],
                                       credentials["username"],
                                       credentials["password"])
    
    if args.list_online:
        online_devices = await connect.list_online(session)
        if args.indent:
            print(dumps(online_devices,indent=4))
        else:
            print(dumps(online_devices))

    if args.run:
        if not args.command or not args.nodeids:
            return
        
        print(args.nodeids)

    await session.close()

if __name__ == "__main__":
    asyncio.run(main())