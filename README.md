ESP-OTA-Server
==============

Very simple OTA firmware server suitable for built-in [ESP8266 HTTP Updater][1].

Main purpose is to serve firmware files and passing MD5 hash -- to verify flashing.

Options:
- `-s` `--bind` listen address (default `:8092`)
- `-d` `--data-dir` data storage location. `<data-dir>/<project>/<file.bin>`

OTA URL: `http://<server-bind-host>/bin/<project>/<file.bin>`


TODO:
- Upload firmware (but for now rsync is enough)
- Automatic TLS via Lets Encrypt, with keeping same cert fingerprint (limitation of esp updater)
- Repository like, for multiple versions (if i really need that)
- Working md5-version check for SPIFFS images

[1]: https://github.com/esp8266/Arduino/tree/master/libraries/ESP8266httpUpdate
