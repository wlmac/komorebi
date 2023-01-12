prefix=/usr

proxy: ./cmd/proxy/main.go $(wildcard **.go)
	go build -o $@ $<

install: proxy
	cp ./proxy ${prefix}/bin/komorebi-proxy

uninstall:
	rm ${prefix}/bin/komorebi-proxy

.PHONY: install uninstall
