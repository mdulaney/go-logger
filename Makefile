SRC_PATH=src/
BIN_PATH=bin/
GO=/cygdrive/c/go/bin/go.exe
KILL=taskkill
KILL_FLAGS=/F /PID

REPORTER=$(BIN_PATH)logger-reporter.exe
SERVER=$(BIN_PATH)logger-server.exe

SERVER_ADDR=127.0.0.1:6000

PID_FILE=logger-server.pid
PID=$(shell cat $(PID_FILE))

all: init $(REPORTER) $(SERVER)

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

$(REPORTER): $(SRC_PATH)logger-reporter.go
	$(GO) build -o $@ $<

$(SERVER): $(SRC_PATH)logger-server.go
	$(GO) build  -o $@ $<

