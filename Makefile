.PHONY=build push


build:
	docker build -t thalesesecurity/contribstats:latest .

push: build
	docker push thalesesecurity/contribstats:latest

run: build
	docker run -ti --rm thalesesecurity/contribstats


circleci: check-env
	circleci local execute -e CONTRIBSTATS_TOKEN=$(CONTRIBSTATS_TOKEN)

helm: check-env
	helm upgrade --install --set config.token=$(CONTRIBSTATS_TOKEN) --set ingress.istio=true --set ingress.enabled=true --set ingress.hosts={contribstats.h2.tes-labs.technology} --set config.organizations={thales-e-security} --recreate-pods contribstats chart

check-env:
ifndef CONTRIBSTATS_TOKEN
  $(error CONTRIBSTATS_TOKEN is undefined, and is needed for proper API Access)
endif

test:
	docker build -t thalesesecurity/contribstats:test -f Dockerfile.test .
	docker run -ti --rm thalesesecurity/contribstats:test
