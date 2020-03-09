#include .env

start:
	@clear
	go build
	cp decent-ft test/
	./decent-ft
