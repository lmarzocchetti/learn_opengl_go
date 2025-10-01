GO := go

TARGETS := 	5_1_transformations \
			6_1_coordinate_system

BIN_DIR := bin

all: $(BIN_DIR) $(TARGETS)

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

# Build each Go target
# $@ is the target name (app1, app2, ...)
$(TARGETS): %:
	$(GO) build -o $(BIN_DIR)/$@ ./cmd/$@

clean:
	rm -rf $(BIN_DIR)