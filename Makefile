.PHONY: build

clean:
	@rm -Rf dist

install:
	@go get ./...
	@go mod vendor
	@echo "[OK] Installed dependencies"

generate:
	@go generate ./...
	@echo "[OK] Files added to embed box"

build: clean generate
	@go build -o ./dist/kisslists main.go
	@echo "[OK] App binary was created"

full: clean install generate
	@go build -a -o ./dist/kisslists main.go
	@echo "[OK] App binary was created"

run:
	@./dist/kisslists