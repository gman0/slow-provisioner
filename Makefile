NAME=slow-provisioner
TAG=v0.3.0

all: provisioner

provisioner:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o  _output/$(NAME) ./cmd/slowprovisioner

image: provisioner
	cp _output/$(NAME) deploy/docker
	docker build -t $(NAME):$(TAG) deploy/docker

clean:
	go clean -r -x
	rm -f deploy/docker/$(NAME)
