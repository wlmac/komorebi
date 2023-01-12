prefix=

proxy: ./cmd/proxy/main.go $(wildcard **.go)
	go build -o $@ $<

install: proxy
	cp ./proxy ${prefix}/usr/bin/komorebi-proxy
	cp ./komorebi@.service ${prefix}/usr/lib/systemd/system/
	cp ./komorebi@.socket ${prefix}/usr/lib/systemd/system/
	mkdir ${prefix}/etc/komorebi

uninstall:
	rm -f ${prefix}/usr/bin/komorebi-proxy
	rm -f ${prefix}/usr/lib/systemd/system/komorebi.service
	rm -f ${prefix}/usr/lib/systemd/system/komorebi.socket
	rm -f ${prefix}/usr/lib/systemd/system/komorebi@.service
	rm -f ${prefix}/usr/lib/systemd/system/komorebi@.socket
	rm -rf ${etc_prefix}/etc/komorebi

.PHONY: install uninstall
