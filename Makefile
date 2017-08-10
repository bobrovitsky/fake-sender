PACKAGE  = gitlab.com/mailmachine/fake-sender
GOPATH   = $(CURDIR)/.gopath
BASE     = $(GOPATH)/src/$(PACKAGE)

BIN = sender

export GOPATH

all: | $(BASE)
	glide install
	cd $(BASE) && go build -o ${BIN} main.go

update:
	glide update

format:
	go fmt gitlab.com/...

pprof-mem:
	go tool pprof -inuse_space ${BIN} mem.pprof

pprof-cpu:
	go tool pprof ${BIN} cpu.pprof

clean:
	rm -rf ${GOPATH}
	rm -f ${BIN}

$(BASE):
	@mkdir -p $(dir $@)
	@ln -sf $(CURDIR) $@
