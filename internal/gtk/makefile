CLANG_FORMAT ?= clang-format
CLANG_TIDY ?= clang-tidy
CFLAGS_GTK := `pkg-config --cflags -- gtk+-3.0`

all: gtk.go
	go build

gtk.go: thunks.h thunks.i
	swig -go -cgo -intgosize 64 thunks.i
	gofmt -w -s gtk.go

fmt:
	-$(RM) gtk.go
	$(CLANG_FORMAT) -i *.c *.h
	gofmt -w *.go

analysis: cgotmp
	$(CLANG_TIDY) --checks=* $(filter-out thunks_wrap.c,$(wildcard *.c)) *.h -- $(CFLAGS_GTK) -I./cgotmp -I.

cgotmp: $(filter-out doc.go,$(wildcard *.go))
	mkdir -p cgotmp
	go tool cgo -objdir ./cgotmp/ -importpath github.com/chaolihf/goey/internal/gtkloop -- -I ./cgotmp/ -g -O2 $(CFLAGS_GTK) $^

clean:
	rm -rf cgotmp
