
# Loom

[![CircleCI](https://circleci.com/gh/cheddartv/loom.svg?style=shield)](https://circleci.com/gh/cheddartv/loom)

Loom is a HLS inter-weaving tool designed to create a redundant manifest. This allows your stream to continue
even if one of your streams goes down. By weaving together any number of index manifests loom creates one manifest
referencing all of them.

## Installing

Download the latest release from [https://github.com/cheddartv/loom/releases](https://github.com/cheddartv/loom/releases)
on the server that generates your primary HLS stream. Then run:
```
dpkg -i <path-to-downloaded-deb-file>
```

## Configuration

Once installed, you need to configure loom for your system. This is all done through /etc/loom.yml

Here you can list multiple instances of output manifests. For each output manifest, there is a list of inputs that will be woven
together. Loom will create a thread for each output listed, weave the inputs together, and then watch the inputs for changes. If
the input files change, loom will likewise update the output.


```
manifests:
  - output: tmp/index.m3u8
    inputs:
      - example/primary/index.m3u8
      - example/backup/index.m3u8
  - output: tmp/index2.m3u8
    inputs:
      - example/primary/index2.m3u8
      - example/backup/index2.m3u8
```

An update to /etc/loom.yml requires that loom be restarted to initialize the new values in the process:
```
sudo /etc/init.d/loom restart
```

## Example Output
```
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=8171676,CODECS="avc1.4d4028,mp4a.40.5",RESOLUTION=1920x1080
../example/primary/1.m3u8
#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=8171676,CODECS="avc1.4d4028,mp4a.40.5",RESOLUTION=1920x1080
../example/backup/1.m3u8
#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=6332540,CODECS="avc1.4d401f,mp4a.40.5",RESOLUTION=1280x720
../example/primary/2.m3u8
#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=6332540,CODECS="avc1.4d401f,mp4a.40.5",RESOLUTION=1280x720
../example/backup/2.m3u8
```

### Prerequisites

Currently loom does not support Windows

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Running the tests

Tests can be run with either `go test` or `ginkgo`

### Current testing dependencies

The current test suite is dependent on the included loom.yml. Changes to the configuration will cause the test suite to fail.

## Authors

* [Ross Cooperman](https://github.com/rosscooperman) - *Initial work*, Organization: [Cheddar](https://github.com/cheddartv)
* [Paul Jones](https://github.com/paulijones) - *Initial work* Organization: [Cheddar](https://github.com/cheddartv)


## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
