# grpc-wiremock

[Wiremock](https://wiremock.org/docs) is a great way to test your connected services.
But it has one drawback. And that is support for **proto** contracts.

**grpc-wiremock** is designed to solve this problem,
and also provide a handy *tool for generating mocks*. And their *automatic reloading*. 

[Here](docs/comparsion.md) you can compare the functionality with existing solutions.

## Getting started

### Quick start
Check out our [wearable](https://github.com/nktch1/wearable) repo 
with an example service. 

There you will find how the service interacts with the mock and 
a docker-compose file for quick startup.

### Interface

```bash
MOCKS_PATH="$(PWD)/test/wiremock"
CERTS_PATH="$(PWD)/certs"
CONTRACTS_PATH="$(PWD)/deps"
WIREMOCK_GUI_PORT=9000

# grpc-wiremock supports multiple APIs simultaneously.
# This means that the another APIs will go up
# on port 8001, 8002, etc.

YOUR_MOCK_API=8000  

docker run \
  -p ${WIREMOCK_GUI_PORT}:${WIREMOCK_GUI_PORT} \
  -p ${YOUR_MOCK_API}:${YOUR_MOCK_API} \
  -v ${MOCKS_PATH}:/home/mock \
  -v ${CERTS_PATH}:/etc/ssl/mock/share \
  -v ${CONTRACTS_PATH}:/proto \
  SberMarket-Tech/grpc-wiremock@latest
```
## Overview

In general, **grpc-wiremock** contains two main components. 

You can read more about each of them here:
- [mocks generator](docs/mocks.md) (**COMING SOON**);
- [grpc-to-http-proxy generator](docs/proxy.md).

In the diagram you can see how your requests are distributed within the **grpc-wiremock** container.

![grpc-wiremock](docs/images/grpc-wiremock.png)

### Benchmarks
You can also read about performance with multiple mock APIs [here](docs/benchmarks.md).

### License
**grpc-wiremock** is under the Apache License, Version 2.0. See the LICENSE file for details.