set shell := ["/bin/zsh", "-c"]

IMAGE_NAME := 'raspall'
CONTAINER_NAME := 'raspall'

default:
  @just --list

# Start the server with air.
[group('go')]
dev: templ vet test
  cd backend && air

# go vet.
[group('go')]
vet:
  cd backend && go vet ./...

# go build.
[group('go')]
build: vet test
  cd backend && go build -o ./build/raspall ./cmd/raspall

# run binary
[group('go')]
run: build
  ./backend/build/raspall

# go test.
[group('go')]
test:
  cd backend && gotest ./...

# generate templ files.
[group('go')]
templ:
  templ generate -path ./backend/internal/views

[group('docker')]
build-mac:
	cd backend && docker build . --tag {{ IMAGE_NAME }}

[group('docker')]
build-linux:
	cd backend && docker build . --platform linux/amd64 --tag ${DOCKER_REPO}

[group('docker')]
push: build-linux
	docker push ${DOCKER_REPO}

[group('docker')]
docker-run: build-mac
	docker run -d --rm -p 80:80 --env AUTH_USERNAME=${AUTH_USERNAME} --env AUTH_PASSWORD=${AUTH_PASSWORD} --name {{ CONTAINER_NAME }} {{ IMAGE_NAME }}

[group('docker')]
docker-stop:
  docker stop {{ CONTAINER_NAME }}

# Remove build and air artifacts.
clean:
  cd backend && rm -rf ./build ./tmp
