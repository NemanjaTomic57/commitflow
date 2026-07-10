#!/bin/bash -xeu

golangci-lint run   
yamllint .          
hadolint Dockerfile 
