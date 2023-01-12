prefix=/usr

proxy: ./cmd/proxy/main.go $(wildcard **.go)
	go build -o $@ $<

install: proxy
	cp ./proxy ${prefix}/bin/komorebi-proxy
	cp ./komorebi.service ${prefix}/lib/systemd/system/
	cp ./komorebi.socket ${prefix}/lib/systemd/system/
	mkdir ${prefix}/etc/komorebi

uninstall:
	rm -f ${prefix}/bin/komorebi-proxy
	rm -f ${prefix}/lib/systemd/system/komorebi.service
	rm -f ${prefix}/lib/systemd/system/komorebi.socket
	rm -rf ${prefix}/etc/komorebi

.PHONY: install uninstall
