# GhostRunner

This project aims to create a way to schedule commands to be run as soon as possible when they offline initialy.<br>
The way to accomplish this is to create a tracked task list, and keep track of it has been successfully done.

# Technical details.

Go(lang) backend server which exposes an HTTP API which can be used to add tasks to the process.<br>
Python executor/runner which actually executes the commands, Python was chosen because of the availability of: [LibMeshCtrl Python Library](https://pypi.org/project/libmeshctrl/).<br>

Create a python virtual environment inside the `runner` folder.

# JSON Templates:

Following is mock data.

### `InfoResponse`

```json
{
  "status": 200,
  "message": "Request successful",
  "data": {
    "example": "This is some mock data"
  }
}
```

---

### `TokenCreateDetails`

```json
{
  "name": "DeploymentToken123"
}
```

---

### `TokenCreateBody`

```json
{
  "authtoken": "abc123securetoken",
  "details": {
    "name": "DeploymentToken123"
  }
}
```

---

### `TokenListBody`

```json
{
  "authtoken": "abc123securetoken"
}
```

---

### `TaskData`

```json
{
  "name": "UpdateScript",
  "command": "sudo apt update && sudo apt upgrade -y",
  "nodeids": ["node-001", "node-002", "node-003"],
  "creation": "2025-05-27T10:15:30Z",
  "status": "pending"
}
```

---

### `TaskBody`

```json
{
  "authtoken": "abc123securetoken",
  "details": {
    "name": "UpdateScript",
    "command": "sudo apt update && sudo apt upgrade -y",
    "nodeids": ["node-001", "node-002", "node-003"],
    "creation": "2025-05-27T10:15:30Z",
    "status": "pending"
  }
}
```

---

### `Device`

```json
{
  "name": "RaspberryPi-01",
  "nodeid": "node-001"
}
```

---

### `PyOnlineDevices`

```json
{
  "online_devices": [
    {
      "name": "RaspberryPi-01",
      "nodeid": "node-001"
    },
    {
      "name": "Server-02",
      "nodeid": "node-005"
    }
  ],
  "offline_devices": [
    {
      "name": "IoT-Gateway",
      "nodeid": "node-003"
    }
  ],
  "total_devices": 3
}
```