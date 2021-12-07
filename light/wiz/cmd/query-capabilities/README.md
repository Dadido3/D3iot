# query-capabilities

A simple tool to query a lamp and check its capabilities.

## Usage

Build the executable by running

``` shell
go build
```

from inside this directory.

Once compiled, use

``` shell
query-capabilities --address "wiz-123abc:38899"
```

or

``` shell
query-capabilities --address "192.168.1.123:38899"
```

with `123abc` replaced by the 6 last characters of your device's MAC address, or `192.168.1.123` replaced by your device's IP.
The result will be written into the `queried` sub-directory.
