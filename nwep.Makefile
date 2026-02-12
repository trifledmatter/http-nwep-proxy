NWEP_VENDOR := vendor/github.com/usenwep/nwep-go
STAMP := $(NWEP_VENDOR)/.nwep-setup

build: $(STAMP)
	go build -mod=vendor ./...

test: $(STAMP)
	go test -mod=vendor ./...

$(STAMP): go.mod go.sum
	go mod vendor
	cd $(NWEP_VENDOR) && bash setup.sh
	@touch $@

clean:
	rm -rf vendor
