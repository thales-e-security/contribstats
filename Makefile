.PHONY=build push


build:
	docker build -t thalesesecurity/contribstats:latest .

push: build
	docker push thalesesecurity/contribstats:latest

run: build
	docker run -ti --rm thalesesecurity/contribstats

helm: check-env push
	helm upgrade --install --set config.token=$(CONTRIBSTATS_TOKEN) --recreate-pods ghstats chart

check-env:
ifndef CONTRIBSTATS_TOKEN
  $(error CONTRIBSTATS_TOKEN is undefined, and is needed for proper API Access)
endif

test:
	docker build -t thalesesecurity/contribstats:test -f Dockerfile.test .
	docker run -ti --rm thalesesecurity/contribstats:test