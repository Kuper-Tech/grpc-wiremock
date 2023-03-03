## Benchmarks

Suppose you need to deploy a lot of mock APIs for your service.

Let's determine how many APIs you can deploy without any noticeable performance degradation.

I tested on a MacBook Pro 16" with these characteristics:
- 2.6 GHz 6-Core Intel Core i7;
- 16 GB 2667 MHz DDR4;
- Intel UHD Graphics 630 1536 MB.

For testing we will choose a simple mock, which we will request from each mock API. We will send requests using the ab utility.
An example of such a mock:

```json 
{
    "request" : {
        "urlPath" : "/HealthCheck",
        "method" : "GET"
    },
    "response" : {
        "status" : 200,
        "body" : "success",
        "headers" : {
            "Content-Type" : "application/json"
        }
    }
}
```

Command to run tests:
```bash
ab -c 10 -n 100 "http://wiremock:${PORT}/HealthCheck"
```

Test parameters:
- the number of requests is 100;
- 10 requests can be executed simultaneously.

Test Scenario:
- preparing N Wiremock Standalone processes;
- waiting for all processes to be ready to serve clients;
- running the ```ab``` utility for all processes simultaneously.

### Results:

In the graph you can see that with the addition of a Wiremock instance, the maximum response time doubles. 

![comparison](images/comparison.png)

