# Komorebi

Middleware that dynamically resizes image to (almost) a specified width/height.

On the [Metropolis](https://maclyonsden.com) website, static media was directly
served from the uploaded version, which were usually unnecessarily larger and
were not in a optimal format (e.g. PNG).

In order to speed up loading times, I developed this to transparently[^1]
reduce bandwidth used.  A quick search did not turn up a simple reverse proxy
that I could put transparently (i.e. disables itself if not explicitly used)
and has the least amount of features I need.
(Actually, [Patrick](https://github.com/ApocalypseCalculator)
found https://github.com/cshum/imagor . Oh well...)

To run, compile `cmd/proxy` and run it with `-net` and `-addr` to specify the
network and address (e.g. `-net unix -addr /tmp/komorebi.sock`).

## API

`/<path>?w=&h=&fmt=`

where
`<path>` is the source path of the file
`w` is the width (0 to leave blank)
`h` is the height (0 to leave blank)
`fmt` is the format (only `webp` is supported)

[Example](https://maclyonsden.com/media/featured_image/b8a57d0a4ca4439987c735affd7f3859.png?fmt=webp&w=500)

## TODO

- [ ] clamp width and height of the resultant image
- [ ] store metrics and clear cache from most unused
- [ ] more formats (JPEG, PNG, AVIF)

[^1]: We only reverse proxy to Komorebi when there is a `width=[0-9]+px`
      request query parameter.
