SRC_PATH=src/
BIN_PATH=bin/
GO="/cygdrive/c/go/bin/go.exe"

CLIENT=$(BIN_PATH)logger-client.exe
SERVER=$(BIN_PATH)logger-server.exe

SERVER_ADDR="127.0.0.1:60000"

all: init $(CLIENT) $(SERVER)

test: $(CLIENT) $(SERVER)
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

