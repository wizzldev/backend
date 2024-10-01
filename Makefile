# env copies the example env to the real place
env:
	cp .env.example .env

buildID:
	git show -s --format='{"build": "%h", "date": "%cD", "author": "%an" }' > ./pkg/configs/build.json

build: buildID
	@go build -o ./bin/wizzl

run: build
	@./bin/wizzl

# handler makes a new handler inside the handlers folder to speed up things
handler:
	@if [ -z "$(name)" ]; then \
		echo "Please specify the name of the handler! Usage: make handler name='handler_name'"; \
	elif [ -e "./app/handlers/$(name).go" ]; then \
      	echo "The handler already exists"; \
    else \
		pascal=$(call to_pascal, $(name)); \
		camel=$(call to_camel, $(name)); \
		printf "%s\n" \
		  "package handlers" \
		  "" \
		  "import \"github.com/gofiber/fiber/v2\"" \
		  "" \
		  "type $$camel struct{}" \
		  "" \
		  "var $$pascal $$camel" \
		  "" \
		  "func ($$camel) Index(*fiber.Ctx) error {" \
		  "	return nil" \
		  "}" > ./app/handlers/$(name).go; \
		  echo "Handler successfully created"; \
	fi

# some defined methods
define to_pascal
$(shell echo $(1) | sed -e 's/_\([a-zA-Z]\)/\U\1/g' -e 's/^./\U&/')
endef

define to_camel
$(shell echo $(1) | sed -e 's/_\([a-zA-Z]\)/\U\1/g' -e 's/^./\L&/' -e 's/_\([a-zA-Z]\)/\U\1/g')
endef
