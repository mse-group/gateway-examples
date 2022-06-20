# Docker image
REGISTRY ?= registry.cn-hangzhou.aliyuncs.com/mse-ingress
IMAGE_NAME ?= httpbin
TAG ?= latest
IMAGE = $(REGISTRY)/$(IMAGE_NAME):$(TAG)

docker-build:
	docker build -f Dockerfile -t $(IMAGE) .

docker-push: docker-build
	docker push $(IMAGE)