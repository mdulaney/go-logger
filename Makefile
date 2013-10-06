SRC_PATH=src/
BIN_PATH=bin/
GO=/cygdrive/c/go/bin/go.exe
KILL=taskkill
KILL_FLAGS=/F /PID

CLIENT=$(BIN_PATH)logger-client.exe
REPORTER=$(BIN_PATH)logger-reporter.exe
SERVER=$(BIN_PATH)logger-server.exe

SERVER_ADDR=127.0.0.1:6000

PID_FILE=logger-server.pid
PID=$(shell cat $(PID_FILE))

TARGETS=$(CLIENT) $(SERVER) $(REPORTER)

all: init $(TARGETS)

destroy: 
	if [ -e $(PID_FILE) ]; then  		\
		$(KILL) $(KILL_FLAGS) $(PID);	\
		rm $(PID_FILE);					\
	fi

clean: destroy
	@rm -f $(SERVER) $(REPORTER)

test: destroy all
	$(SERVER) -l $(SERVER_ADDR) &
	$(REPORTER) -t $(SERVER_ADDR)

init:
	if [ ! -e "$(BIN_PATH)" ]; then \
		mkdir $(BIN_PATH);	\
	fi

$(CLIENT): $(SRC_PATH)logger-client.go
	$(GO) build -o $@ $<

$(REPORTER): $(SRC_PATH)logger-reporter.go
	$(GO) build -o $@ $<

$(SERVER): $(SRC_PATH)logger-server.go
	$(GO) build  -o $@ $<

