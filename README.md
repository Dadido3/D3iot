# D3iot

A software package to handle IoT devices.

Right now this repository doesn't contain much that can run on its own.
There are several sub-packages you still might be interested in:

- [github.com/Dadido3/D3iot/light/drivers/wiz/](light/drivers/wiz/): Everything you need to communicate with WiZ light devices. Also contains tools to debug or profile these devices.
- [github.com/Dadido3/D3iot/light/](light/): A more general interface to control and query light emitting devices.
- [github.com/Dadido3/D3iot/light/emission/](light/emission/): A library that contains color space math stuff for anything that emits light.

## Examples and tools

- [Bias-light](light/tools/bias-light/): A tool that helps you to control a bias light for your computer screen.
