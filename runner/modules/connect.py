import meshctrl
from json import dumps

class connect:
    @staticmethod
    async def quit(session: meshctrl.Session) -> None: # Function to use when quitting, but gracefully.
        await session.close()

    @staticmethod
    async def connect(hostname: str, username: str, password: str) -> meshctrl.Session:
        session = meshctrl.Session(
            hostname,
            user=username,
            password=password
        )
        await session.initialized.wait()
        return session
    
    @staticmethod
    async def run(session: meshctrl.Session, command: str, nodeids: list[str]) -> None:
        try:
            response = await session.run_command(nodeids=nodeids,
                                                 command=command,
                                                 ignore_output=False,
                                                 timeout=900)
        except Exception as error:
            print("Run Command failed.", error)
            return

        for device in response:
            print(dumps(response[device]["result"]))

    @staticmethod
    async def list_online(session: meshctrl.Session) -> dict: # Default is return online devices, but function can also return the offline devices if specified.

        raw_device_list = await session.list_devices()

        complete_list = {}
        complete_list["online_devices"] = []
        complete_list["offline_devices"] = []

        for raw_device in raw_device_list:
            if raw_device.connected:
                complete_list["online_devices"].append({
                    "name": raw_device.name,
                    "nodeid": raw_device.nodeid
                })
            else:
                complete_list["offline_devices"].append({
                    "name": raw_device.name,
                    "nodeid": raw_device.nodeid
                })
        complete_list["total_devices"] = len(complete_list["online_devices"]) + len(complete_list["offline_devices"])
                
        return complete_list