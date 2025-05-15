# Используем bin в текущей директории для установки плагинов protoc
LOCAL_BIN := $(CURDIR)/bin

BUF_BUILD := $(LOCAL_BIN)/buf

# Устанавливаем необходимые плагины
.bin-deps: export GOBIN := $(LOCAL_BIN)
.bin-deps:
	$(info Installing binary dependencies...)

	go install github.com/bufbuild/buf/cmd/buf@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Генерация .pb файлов с помощью buf
.buf-generate:
	$(info Run buf generate...)
	PATH="$(LOCAL_BIN):$(PATH)" $(BUF_BUILD) generate

generate: .buf-generate .tidy

.tidy:
	go mod tidy
