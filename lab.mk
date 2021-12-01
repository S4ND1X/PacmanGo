APP_NAME=main

deps:
	go get github.com/danicat/simpleansi

build:
	go build -o ${APP_NAME}.o ${APP_NAME}.go


test: deps build
	@echo Test1
	./${APP_NAME}.o || true

	@echo Test2
	./${APP_NAME}.o 4 || true

run: deps build
	go build -o ${APP_NAME}.o ${APP_NAME}.go && ./${APP_NAME}.o ${enemies}
	
clean:
	rm -rf main