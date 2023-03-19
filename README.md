# HightPerformancePaymentGateway-BalanceServiceCommand

service balance for all partner, provider, end user, ...

- [Context](#Context)
- [Soluion](#Solution)
- [System design](#SystemDesign)

## Context <a name="Context"></a>

Service Balance handle balancer of amount, plus, reduction, check amount valid. Service hasn't much partner but race
conditions is very large. One partner maybe have many users. Let's say we have several thousand partners globally but we
have several billion users in the world. Race conditions partner's balance is very large. Assuming the amount of RPS is
20 M. This is over the capacity of the DB, even in the case of sharding, because the number of partners is not large but
the race conditions are very large. </br>
We solved the problem and scale horizon service. </br>

## Soluion <a name="Solution"></a>

We use divide and conquer and LMAX architecture. </br>
We group partner to many group. Every group has one LMAX service. We push request of one group to one LMAX service
mapping respectively. We can infinitely scale the inventory, or balance problem. </br>

About LMAX architecture, has detail in link: https://martinfowler.com/articles/lmax.html

## System design <a name="SystemDesign"></a>

![](img_readme/system_design.png)

1) Request to WAF to LB </br>
2) Request Service group By Partner </br>
3) Request To LMAX architect matching and response with status pending </br>
4) Result of LMAX will publish to Kafka </br>
5) and 6) Corresponding services and DB subscribe kafka and update status </br>

Note: To scale out this service, we don't keep status request's order. This is a bottleneck. Ex: 1000 requests in step
1, but 1000 requests in step 4 may have a different order. As the realtime system is very fast, this is not necessary
for us. Even if person A requests 0.01 ms before person B, but requesting person B to the server first, it will be
handled first, and person A may not be able to buy goods. This is accepted by us with a very small tolerance, there will
be no bad user experience. </br>

## Benchmark <a name="benchmark"></a>

Run: make benchmark </br>
With a laptop i7 5600U 2.2ghz 8Gb ram, with an instance running a group partner we easily got 1000000 requests with
about 110s. Each request is closed into a bucket with 20 request orders. Thus, I handled 20 M request balance with about
110s. We handle about 200K orders per second. Running with 5 LMAX clusters, we can handle 1M orders per second. This
number does not make too much of a difference when running on Aws. Here we can scale the number of clusters to a very
large horizontally, but 1M orders per second is also a very large number about. </br>

```
 1000000 / 1000000 [==================================================================================================================================================================] 100.00% 9047/s 1m50s
Done!
Statistics        Avg      Stdev        Max
  Reqs/sec      9056.96    3418.32   28633.70
  Latency       66.27ms     8.86ms   294.84ms
  Latency Distribution
     50%    59.38ms
     75%    79.78ms
     90%    89.23ms
     95%    94.98ms
     99%   114.28ms
  HTTP codes:
    1xx - 0, 2xx - 1000000, 3xx - 0, 4xx - 0, 5xx - 0
    others - 0
  Throughput:    24.03MB/s

```

## Function <a name="function"></a>

1) Api request balance </br>
2) Auth, auth-internal </br>
3) Auto push job to queue, auto handle it </br>
4) Deploy, CICD with github action, docker, k8s, ssl, LB, auto-scale. </br>

## Todo <a name="todo"></a>

1) Load prepare from storage </br>
2) Publish result of balance request to kafka </br>
3) Backup server and backup info balance realtime, for recovery if crash occurs </br>
4) Auto recovery if system crash. </br>