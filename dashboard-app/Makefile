.PHONY: ui-build
ui-build:
	@pushd ui > /dev/null; \
		pnpm i > /dev/null; \
		pnpm build > /dev/null; \
		popd > /dev/null

.PHONY: build
build: ui-build
	@mkdir -p bin
	go build -o bin/api ./cmd/api/...

.PHONY: run
run: build
	@bin/api
