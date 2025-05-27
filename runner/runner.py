#!/bin/python3

import argparse
from ast import literal_eval
import asyncio
from json import dumps

from modules.connect import connect
from modules.utilities import utilities

def cmd_flags() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Process command-line arguments")

    parser.add_argument("-lo", "--list-online", action='store_true', help="Specify if the program needs to list online devices.")
    parser.add_argument("--run", action='store_true', help="Make the program run a command.")
    parser.add_argument("--command", type=str, help="Specify the actual command that is going to run.")
    parser.add_argument('--nodeid', nargs='+', help='List of node IDs')

    parser.add_argument("-i", "--indent", action='store_true', help="Specify whether the output needs to be indented.")

    return parser.parse_args()

async def prepare_command(command: str, nodeid: str) -> str: # Have some checks so it happens correctly.
    if len(nodeid) < 1 or len(command) < 1:
        print("No nodeid or command passed... quiting.")
        return ""
    
    return nodeid

async def main() -> None:
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
        return await connect.quit(session) # Exit gracefully. Because python.

    if args.run:
        if not args.command or not args.nodeid:
            print("When using run, also use --command and --nodeid")
            return await connect.quit(session) # Exit gracefully. Because python.
    
        command = args.command
        nodeid = args.nodeid
        nodeid = await prepare_command(command, nodeid)

        await connect.run(session, command, nodeid)

    await session.close()

if __name__ == "__main__":
    asyncio.run(main())