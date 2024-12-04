# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

tidy:
	go mod tidy

dev-up:
	minikube start
dev-down:
	minikube delete

DB_DSN				:= "postgresql://postgres:admin1234@localhost:5433/buy-better-core?sslmode=disable"
DB_NAME				:= "buy-better-core"
DB_USERNAME			:= "postgres"
CONTAINER_NAME		:= "pg-core-dev-db"

BASE_IMAGE_NAME 	:= opplieam
SERVICE_NAME    	:= bb-product-server
VERSION         	:= "0.0.1-$(shell git rev-parse --short HEAD)"
VERSION_DEV         := "cluster-dev"
SERVICE_IMAGE   	:= $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)
SERVICE_IMAGE_DEV   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION_DEV)

DEPLOYMENT_NAME		:= product-server-deployment
SECRET_NAME			:= product-server-secret
NAMESPACE			:= buy-better


docker-build-dev:
	@eval $$(minikube docker-env);\
	docker build \
		-t $(SERVICE_IMAGE_DEV) \
    	--build-arg BUILD_REF=$(VERSION_DEV) \
    	--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
    	-f dev.Dockerfile \
    	.
docker-build-prod:
	docker build \
		-t $(SERVICE_IMAGE) \
    	--build-arg BUILD_REF=$(VERSION) \
    	--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
    	.

docker-push:
	docker push $(SERVICE_IMAGE)

docker-build-push: docker-build-prod docker-push

gen-prod-chart:
	rm -rf .genmanifest
	helm template $(SERVICE_NAME) chart -f chart/dev.values.yaml \
		--set image=$(SERVICE_IMAGE) \
		--set imagePullPolicy=Always \
		--output-dir .genmanifest

helm-prod:
	helm upgrade --install -f ./chart/prod.values.yaml \
	--set image=$(SERVICE_IMAGE) \
	--set imagePullPolicy=Always \
	bb-product-server ./chart

kus-dev:
	kubectl apply -k k8s/dev/
helm-dev:
	helm upgrade --install -f ./chart/dev.values.yaml bb-product-server ./chart
dev-restart:
	kubectl rollout restart deployment $(DEPLOYMENT_NAME) --namespace=$(NAMESPACE)
dev-stop:
	kubectl delete -k k8s/dev/
dev-helm-stop:
	helm uninstall bb-product-server

dev-apply: tidy docker-build-dev kus-dev

dev-apply-helm: tidy docker-build-dev helm-dev

jet-gen:
	jet -dsn=$(DB_DSN) -path=./.gen