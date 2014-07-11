all: nstop
nstop :	cli.go
	go test
cli.go :
	cd cli && go test
