build:
	cd _examples && go build tailf.go

clean:
	cd _examples && rm -f tailf tailf.exe

.PHONY: build
