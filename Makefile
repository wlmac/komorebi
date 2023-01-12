prefix=/usr

proxy: ./cmd/proxy/main.go $(wildcard **.go)
	go build -o $@ $<

install: proxy
	cp ./proxy ${prefix}/bin/komorebi-proxy
	cp ./komorebi.service ${prefix}/lib/systemd/system/
	cp ./komorebi.socket ${prefix}/lib/systemd/system/

uninstall:
	rm ${prefix}/bin/komorebi-proxy
	rm ${prefix}/lib/systemd/system/komorebi.service
	rm ${prefix}/lib/systemd/system/komorebi.socket

.PHONY: install uninstall
