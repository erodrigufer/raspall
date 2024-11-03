set shell := ["/bin/zsh", "-c"]

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

# go test.
[group('go')]
test:
  cd backend && gotest ./...

# generate templ files.
[group('go')]
templ:
  templ generate -path ./backend/internal/views

# Remove build and air artifacts.
clean:
  cd backend && rm -rf ./build ./tmp
