#See: https://cloud.google.com/functions/docs/2nd-gen/getting-started#http

REGION=asia-northeast1
FUNC_NAME=zl-example-helloworld
PROJECT_ID=$(shell gcloud config get-value project)
PROJECT_NUMBER=$(shell gcloud projects list --filter="project_id:$(PROJECT_ID)" --format='value(project_number)')

# This is a test command. It is not necessary when actually using zl.
cp:
	mkdir -p ./pkg_test ./pkg_test/zl
	cp ../../../*.go ./pkg_test/zl
	rm ./pkg_test/zl/*_test.go

start: cp
	export FUNCTION_TARGET=HelloGet; go run cmd/main.go

deploy: cp
	gcloud beta functions deploy $(FUNC_NAME) \
    --gen2 \
    --runtime go116 \
    --trigger-http \
    --entry-point HelloGet \
    --source . \
    --allow-unauthenticated

show:
	gcloud beta functions describe $(FUNC_NAME) --gen2

log:
	gcloud beta functions logs read $(FUNC_NAME) --gen2 --limit=100

open:
	open https://console.cloud.google.com/functions/details/$(REGION)/$(FUNC_NAME)?env=gen2
