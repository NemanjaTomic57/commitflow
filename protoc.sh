#!/bin/bash -xeu

protoc --go_out=. --go_opt=paths=source_relative proto/commitflow.proto
