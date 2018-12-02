#!/usr/bin/env bash

curl -i -X POST http://localhost:8080/api/hangout/message -H "Content-Type: application/json" --data-binary "@./added-event.json"