# Bias-light

Minimal application to synchronize a single (WiZ) light with your screen content. The idea is that you set up a single light bulb behind your computer monitor and let it be controlled by this application.

## Build

Build the executable with the go compiler from inside this directory:

``` shell
go build
```

## Usage

Start the application and pass the to be controlled device as parameter:

``` shell
bias-light --device-wiz "wiz-123abc:38899" --max-luminance 1500
```

With the following parameters:

- `--device-wiz "wiz-123abc:38899"`: The WiZ device by its hostname and port, you can also use its IP.
- `--max-luminance`: The maximum luminance in lumens, if omitted the emission value will be scaled to the full dynamic range of the light device.
- `--no-white-optimization`: Disables optimization for high CRI and high luminance by disabling white emitters. This may help to get a lower latency with some light devices, due to weaker low-pass filtering in the light device's primary color emitters.
- `--brighten`: Brightens up all colors by the given factor.
