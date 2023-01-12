# Komorebi

Middleware that dynamically resizes image to (almost) a specified width/height.

On the [Metropolis](https://maclyonsden.com) website, static media was directly
served from the uploaded version, which were usually unnecessarily larger and
were not in a optimal format (e.g. PNG).

In order to speed up loading times, I developed this to transparently[^1]
reduce bandwidth used.  A quick search did not turn up a simple reverse proxy
that I could put transparently (i.e. disables itself if not explicitly used)
and has the least amount of features I need.

To run, compile `cmd/proxy` and run it with `-net` and `-addr` to specify the
network and address (e.g. `-net unix -addr /tmp/komorebi.sock`).

[^1]: We only reverse proxy to Komorebi when there is a `width=[0-9]+px`
      request query parameter.
