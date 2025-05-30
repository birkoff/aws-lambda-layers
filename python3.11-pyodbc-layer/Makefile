RUNTIME = python3.11
REGION = us-east-1

LAYER_NAME = pyodbc-layer
DESCRIPTION = "pyodbc for Python 3.11"

.PHONY: all build run clean deploy

all: build clean

build:
	@echo "Building Docker image..."
	docker build --platform linux/amd64 -t $(LAYER_NAME) .
	docker run --rm --platform linux/amd64 \
		-v "$(PWD)":/var/task \
		--entrypoint /bin/bash \
		$(LAYER_NAME) \
		-c "cp /opt/$(LAYER_NAME).zip /var/task/ && chown $(shell id -u):$(shell id -g) /var/task/$(LAYER_NAME).zip"
	@echo "Layer created: $(LAYER_NAME).zip"

test:
	docker compose up --build

run:
	docker run -it --rm --platform linux/amd64 \
		-v "$(PWD)":/var/task \
		--entrypoint /bin/bash \
		$(LAYER_NAME)

clean:
	@echo "Cleaning up..."
	-docker rmi $(LAYER_NAME) 2>/dev/null || true
	@echo "Cleanup complete!"

deploy:
	@echo "Deploying $(LAYER_NAME).zip to AWS Lambda..."
	aws lambda publish-layer-version \
		--layer-name $(LAYER_NAME) \
		--zip-file fileb://$(LAYER_NAME).zip \
		--compatible-runtimes $(RUNTIME) \
		--region $(REGION) \
		--description $(DESCRIPTION)