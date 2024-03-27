run: build
	@./bin/website

install:
	@go install github.com/a-h/templ/cmd/templ@latest
	@go get ./...
	@go mod vendor
	@go mod tidy
	@go mod download
	@npm install -D tailwindcss postcss autoprefixer
	@npm install -D daisyui@latest

build:
	@tailwindcss -i view/css/app.css -o public/styles.css
	@templ generate view
	@go build -o bin/website main.go