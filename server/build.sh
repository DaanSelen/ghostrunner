#!/bin/bash

echo "Building..."
cd src/ || exit 1
go build -o ../
cd ..
echo "Done."
