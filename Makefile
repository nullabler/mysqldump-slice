SHELL := /bin/bash
MAKEFLAGS += --silent
ARGS = $(filter-out $@,$(MAKECMDGOALS))

.default: help

include .build/.mk/*/*.mk
include .build/.mk/*.mk
