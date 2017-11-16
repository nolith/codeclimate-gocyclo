.PHONY: image

IMAGE_NAME ?= codeclimate/codeclimate-gocyclo

image:
	docker build --rm -t $(IMAGE_NAME) .

test: image
	CODECLIMATE_DEBUG=1 codeclimate analyze --dev
