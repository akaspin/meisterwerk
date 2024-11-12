UP_ARGS ?=
up:
	docker compose up -d --remove-orphans --wait $(UP_ARGS)
.PHONY: stack-up

clean:
	docker compose down -v --remove-orphans
.PHONY: clean

test:
	go test ./...
.PHONY: test

gen/server/%:
	rm -rf api/gen/server/$*
	openapi-generator generate -g go-server -i api/openapi/$*.yaml -o api/gen/server \
    		--additional-properties="outputAsLibrary=true,sourceFolder=$*,router=mux,packageName=$*,enumClassPrefix=true" \
    		--openapi-generator-ignore-list "README.md,docs/*,api/*,.gitignore,.openapi-generator-ignore,.openapi-generator/*"
	goimports -w $(wildcard api/gen/server/$*/*)
.PHONY: gen/server/%

gen/client/%:
	rm -rf api/gen/client/$*
	openapi-generator generate -g go -i api/openapi/$*.yaml -o api/gen/client/$* \
		--additional-properties=isGoSubmodule=true,withGoMod=false,packageName=$* \
		--openapi-generator-ignore-list "README.md,.travis.yml,test/*,docs/*,git_push.sh,api/*,.gitignore,.openapi-generator-ignore,.openapi-generator/*"
	goimports -w $(wildcard api/gen/client/$*/*)
.PHONY: gen/client/%


gen: \
	gen/server/quotes \
	gen/client/quotes \
	gen/server/orders \
	gen/client/orders
.PHONY: gen


