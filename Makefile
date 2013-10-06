SRC_PATH=src/
BIN_PATH=bin/
GO=/cygdrive/c/go/bin/go.exe
KILL=taskkill
KILL_FLAGS=/F /PID

CLIENT=$(BIN_PATH)logger-client.exe
SERVER=$(BIN_PATH)logger-server.exe

SERVER_ADDR=127.0.0.1:6000

PID_FILE=logger-server.pid
PID=$(shell cat $(PID_FILE))

all: init $(CLIENT) $(SERVER)

destroy: 
	if [ -e $(PID_FILE) ]; then  	\
	$(KILL) $(KILL_FLAGS) $(PID);	\
	rm $(PID_FILE);					\
	fi

test: destroy all
	$(SERVER) -l $(SERVER_ADDR) &
	$(CLIENT) -t $(SERVER_ADDR)

init:
	if [ ! -e "$(BIN_PATH)" ]; then \
		mkdir $(BIN_PATH);	\
	fi

$(CLIENT): $(SRC_PATH)logger-client.go
	$(GO) build -o $@ $<

$(SERVER): $(SRC_PATH)logger-server.go
	$(GO) build  -o $@ $<

