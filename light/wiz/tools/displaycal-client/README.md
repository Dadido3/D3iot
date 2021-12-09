# DisplayCAL client

This program uses the `Web @ localhost` "Display" of DisplayCAL to use a WiZ light as displaying device.
DisplayCAL can then be used to create a color profile of the light device.

## State

It works in theory, but a colorimeter isn't the best device to calibrate those type of light devices.
A spectrometer would be better suited, but i don't have access to such thing.

At least i was able to confirm that the light intensity are mapped nearly linear to the RGB values.
Which is to be expected, as they probably use the values to directly set the PWM duty cycle of each color.

## Usage

Build the executable by running

``` shell
go build
```

from inside this directory.

Once compiled, use

``` shell
displaycal-client --address "wiz-123abc:38899" --http-server "http://localhost:8080/"
```

or

``` shell
displaycal-client --address "192.168.1.123:38899" --http-server "http://localhost:8080/"
```

with `123abc` replaced by the 6 last characters of your device's MAC address, or `192.168.1.123` replaced by your device's IP.

### Additional parameters

As some devices may limit their total light output, it may not be possible to output the full value on all RGB channels at the same time.
The device will normalize the RGB (and cold white/warm white) values so that their sum is below some specific value (~512 for one WiZ light i tested).
To circumvent this limitation, you can scale the channels down individually by using

``` shell
displaycal-client --address "192.168.1.123:38899" --max-r 170 --max-g 170 --max-b 170
```
