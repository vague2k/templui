# Dev Tools
templ:
	templ generate --watch --proxy="http://localhost:8090" --open-browser=false -v

server:
	air \
	--build.cmd "go build -o tmp/bin/main ./cmd/server" \
	--build.bin "tmp/bin/main" \
	--build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

tailwind-clean:
	tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css --clean

tailwind-watch:
	tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css --watch

dev:
	make tailwind-clean
	make -j3 templ server tailwind-watch

debug:
	make -j3 templ tailwind-app tailwind
	
generate-sitemap:
	go run ./cmd/sitemap/main.go --baseurl="https://templui.io" --routes="./cmd/server/main.go" --output="./static/sitemap.xml" --robots="./static/robots.txt"

generate-icons:
	go run cmd/icongen/main.go

install-compinstall:
	go install ./cmd/compinstall
