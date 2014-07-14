all: nstop
nstop :	cli.go
	go test
cli.go :
	cd arguments && go test
