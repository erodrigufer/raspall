set shell := ["/bin/zsh", "-c"]

IMAGE_NAME := 'raspall'
CONTAINER_NAME := 'raspall'
BUILD_HASH_COMMIT := ` \
  if [ -n "$(git status --porcelain)" ]; then \
    echo "$(git log -1 --pretty=%h)+dirty"; \
  else \
    echo "$(git log -1 --pretty=%h)"; \
  fi`

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
	cd backend && docker build . --build-arg BUILD_HASH_COMMIT={{ BUILD_HASH_COMMIT }} --tag {{ IMAGE_NAME }}

[group('docker')]
build-linux:
	cd backend && docker build . --platform linux/amd64 --build-arg BUILD_HASH_COMMIT={{ BUILD_HASH_COMMIT }} --tag ${DOCKER_REPO}

# Build Linux Docker image and push to DockerHub.
[group('docker')]
push: build-linux
	docker push ${DOCKER_REPO}

# Locally execute the Docker image at localhost:80.
[group('docker')]
docker-run: build-mac
	docker run -d --rm -p 80:80 --env AUTH_USERNAME=${AUTH_USERNAME} --env AUTH_PASSWORD=${AUTH_PASSWORD} --name {{ CONTAINER_NAME }} {{ IMAGE_NAME }}

# Stop the local Docker container.
[group('docker')]
docker-stop:
  docker stop {{ CONTAINER_NAME }}

# Remove build and air artifacts.
clean:
  cd backend && rm -rf ./build ./tmp
