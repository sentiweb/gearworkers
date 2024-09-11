
gearman:
	docker run -p 4730:4730 artefactual/gearmand:1.1.21.2-alpine

run:
	go run ./cmd/gearworkers

tester:
	go run ./cmd/tester/

build-dummy:
	go build ./cmd/dummy

release:
	goreleaser release --skip=publish --clean