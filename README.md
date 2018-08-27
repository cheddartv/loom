# Loom
Loom is a HLS inter-weaving tool designed to create a redundant manifest. This alows your stream to continue even if one of your streams goes down. By weaving together any number of index manifests loom creates one manifest referencing all of them.

## Installing

On your server that generates your primary HLS stream, install loom:
```
sudo apt update
sudo apt install loom
```

## Deployment

Once installed, you need to configure loom for your system. This is all done through loom.yml

Here you can list multiple instances of output manifests. For each output manifest, there is a list of inputs that will be woven together. Loom will create a thread for each output listed, weave the inputs together, and then watch the inputs for changes. If the input files change, loom will likewise update the output.


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

An update to loom.yml requires that loom be restarted to initilaze the new values in the process

### Prerequisites <i don't think we have any?>

What things you need to install the software and how to install them

```
Give examples
```

## Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us.


## Running the tests

Tests can be run with either `go test` or `ginkgo`

### Current testing dependencies

The current test suite is dependent on the included loom.yml. Changes to the configuration will cause the test suite to fail.

## Authors

* [Ross Cooperman](https://github.com/rosscooperman) - *Initial work*
* [Paul Jones](https://github.com/paulijones) - *Initial work*


## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Hat tip to anyone whose code was used
* Inspiration
* etc
