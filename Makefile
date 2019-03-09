.PHONY: image
image:
	go build
	docker build . -t anfernee/k8s-issue-74839
