# GhostRunner

This project aims to create a way to schedule commands to be run as soon as possible when they offline initialy.<br>
The way to accomplish this is to create a tracked task list, and keep track of it has been successfully done.

# Technical details.

Go(lang) backend server which exposes an HTTP API which can be used to add tasks to the process.<br>
Python executor/runner which actually executes the commands, Python was chosen because of the availability of: [LibMeshCtrl Python Library](https://pypi.org/project/libmeshctrl/).<br>
