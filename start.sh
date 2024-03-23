#!/bin/bash

go mod tidy

(go run main.go &)
(go run ./GoScrapper/main.go &)