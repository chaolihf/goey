CLANG_FORMAT ?= clang-format
CLANG_TIDY ?= clang-tidy
CFLAGS_GTK := `pkg-config --cflags -- gtk+-3.0`

all:
	go build

fmt:
	$(CLANG_FORMAT) -i *.c *.h
	gofmt -w *.go

analysis: cgotmp
	$(CLANG_TIDY) --checks=* *.c *.h -- $(CFLAGS_GTK) -I./cgotmp -I.

cgotmp: $(filter-out doc.go,$(wildcard *.go))
	mkdir -p cgotmp
	go tool cgo -objdir ./cgotmp/ -importpath github.com/chaolihf/goey/internal/gtkloop -- -I ./cgotmp/ -g -O2 $(CFLAGS_GTK) $^

clean:
	rm -rf cgotmp
