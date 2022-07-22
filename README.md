# HightPerformancePaymentGateway-BalanceService
service balance for all partner, provider, end user, ...


- [Review characteristics balance service in payment gateway](#ReviewCharacteristicsBalanceServiceInPaymentGateway)
- [System design service balancer handle several billion transaction per day](#SystemDesignServiceBalancerHandleSeveralBillionTransactionPerDay)
  - [Link all chart](#LinkAllChart)
  - [Follow check balance avaible](#FollowCheckBalanceAvaible)
  - [When use this follow?](#WhenUseThisFollow)
  - [Bottleneck](#Bottleneck)
  - [Solution](#Solution)
  - [Detail solution](#DetailSolution)
  - [Why i don't use auto sharding](#WhyIDontUseAutoSharding)
  - [Expansion sharding for e-wallet problem](#ExpansionForEWalletProblem)
    - [Problem](#Problem)
    - [Bottleneck](#Bottleneck)
    - [Solution](#Solution)
  - [Save data](#ProblemSaveData)
    - [Problem save data](#ProblemSaveData)
    - [Solution save data](#SolutionSaveData)
    - [Detail solution save data](#DetailSolutionSaveData)

  - [Get Data Trans](#GetDataTrans)
    - [Problem get data trans](#ProblemGetDataTrans)
    - [Solution get data trans](#SolutionGetDataTrans)


- [System design system get data for several billion user](#SystemDesignSystemGetDataForSeveralBillionUser)
  - [Problem system high get data](#ProblemSystemHighGetData)
    - [Problem high qps for io DB](#ProblemHighQpsForIoDB)
    - [Problem overload request to DB when invalidate cache](#ProblemOverloadRequestToDBWhenInvalidateCache)
  - [Solution system high get data](#SolutionSystemHighGetData)
    - [Solution get from cache](#SolutionGetFromCache)
    - [Solution update cache](#SolutionUpdateCache)
    - [Solution only handle one request update cache for one key in race conditions](#SolutionOnlyHandleOneRequestUpdateCacheForOneKeyInRaceConditions)
    - [Solution smart select DB](#SolutionSmartSelectDB)

- [System divide infra by region](#SystemDivideInfraByRegion)
- [Problem divide infra by region](#ProblemDivideInfraByRegion)
  - [Problem DB divide infra by region](#ProblemDBivideInfraByRegion)
  - [Problem cache divide infra by region](#ProblemCacheivideInfraByRegion)
- [Solution divide infra by region](#SolutionDivideInfraByRegion)
  - [Solution DB divide infra by region](#SolutionDBivideInfraByRegion)
  - [Solution cache divide infra by region](#SolutionCacheivideInfraByRegion)
- [Service cache is independence between services](#ServiceCacheIsIndependenceBetweenServices)


## Review characteristics balance service in payment gateway <a name="ReviewCharacteristicsBalanceServiceInPaymentGateway"></a>
The balance service at payment gateways and e-wallets in general has several characteristics: </br>
1) Payment model is a hierarchical model, the top level payment will manage its child payment, its child payment will manage the payment grandchildren, .... so on to the enduser. </br>

2) Service balance will have to work with alternative payment models that require checking the partner balance before paying, most commonly found in the ebill model, with the service without balance checking, the cash flow has been circulated from the enduser to the bank. </br>

3) Due to the decentralized model, each payment usually has a small number of direct child partners, usually less than 10,000 partners. This feature determines how the DB is selected. </br>

4) Service balance requires extremely strict ACID. Every mistake costs a great deal of money. </br>


## Link all chart <a name="LinkAllChart"></a>

## Follow check balance avaible  <a name="FollowCheckBalanceAvaible"></a>
![](img_readme/fl_balance_avaible_1.png)
![](img_readme/fl_balance_avaible_2.png)


## When use this follow? <a name="WhenUseThisFollow"></a>
Follow is extended to payment services that need to check partner balances before making payments. </br>

## Bottleneck <a name="Bottleneck"></a>
When I want the system to handle several billion transactions per day, bottlenecks appear at many points, often points where cannot be scale horizon  . </br>

Specifically: Those are IO operations: check balance, create transaction, update total amount. </br>

## Solution <a name="Solution"></a>
![](img_readme/sharding_balance_db.png)

Solution: </br>
There are many solutions to this problem, the most common are: </br>
1) No blocking with Architecture LMAX </br>
2) DB sql scale horizon </br>

I chose solution 2. The reasons I chose it:

1) As you can see, the peculiarity of data balance is strict ACID. The payment of IO race conditions is inevitable. I was looking at several locking solutions: Redlock, lockDB, locking an instanse cache, I decided to go with the lockDB and sharding DB solution. I want to take full advantage of sql's integrity mechanism and system sql design of sharding to be able to have horizontal performance. </br>

2) Mysql(sql) is so famous and stable, it has reached the point of nice in terms of practical application technology. It is mature, stable and powerful. </br>

3) The cost of a sharding node is very small compared to its business scalability. <br>

## Detail solution <a name="DetailSolution"></a>
Characteristics of the number of partners < 10000 (this is a very large number, I need 2 parts). </br>
1) Part one is the part of managing partner sharding, it will regulate and manage which partners + region exist on which sharding </br>
2) Part two is the sharding data part, it will shard and contain balance partner information. Struct DB sharding requires enough basic fields, extension does not edit fields and adds data and extended json objects. This is an important thing with manual sharding systems. </br>

## Why i don't use auto sharding <a name="WhyIDontUseAutoSharding"></a>

I have reference some tools for auto sharding mysql like Vitess. It offers a lot of automated and convenient features. I ask myself, should I use it? </br>

Re-survey the feature I need, it's quite simple, including control partner sharding and creating sharding, the feature that needs to be automated is not too much. When I manually shard, I can control almost everything, every query, move partner to other shard. I won't be attached to a other layer of tool auto sharding anymore. After considering the problem, I chose manual sharding. </br>


## Expansion sharding for e-wallet problem <a name="ExpansionForEWalletProblem"></a>
## Problem <a name="Problem"></a>
The problem is similar to the payment gateway, with a slight difference. With e-wallets, the object is a user, not a partner. The number of users can be up to several hundred million, but the number of payments of a user is several tens of thousands of times smaller than that of a partner. </br>

## Bottleneck <a name="Bottleneck"></a>
Same with payment gateway </br>
## Solution  <a name="Solution"></a>
Same with payment gateway </br>

## Save data  <a name="SaveData"></a>
## Problem save data <a name="ProblemSaveData"></a>
With a number of several billion transactions a day, long-term data storage on mysql is a utopia. Even if mysql is available, it is not easy to ensure stable query and operation of mysql. </br>
## Solution save data <a name="SolutionSaveData"></a>
The key here is that I need mysql to be lightweight to ensure system performance. The queries with the balance information of the success order are usually simple queries, do not require complex aggregation, I choose Cassandra DB for the solution. Cassandra is a distributed DB, it is born to serve huge storage needs with extremely high read threshold and simple queries.

## Detail solution save data <a name="DetailSolutionSaveData"></a>
![](img_readme/move_db_from_mysql_to_cassandra.png)


## System design system get data for several billion user <a name="SystemDesignSystemGetDataForSeveralBillionUser"></a>
A system that pays api get data for billions of users is a difficult system, has many problems and needs to be calculated from an overview to meticulously each problem for the system to work stably. I will dissect each problem encountered and the solution in the following section. </br>

## Problem system high get data <a name="ProblemSystemHighGetData"></a>
## Problem high qps for io DB <a name="ProblemHighQpsForIoDB"></a>
With billions of users online and getting data, the amount of qps can exceed the processing capacity of any DB, which is the first problem with most high-load systems. </br>

## Problem overload request to DB when invalidate cache <a name="ProblemOverloadRequestToDBWhenInvalidateCache"></a>
![](img_readme/overload_cpu_when_invalidate_cache.png)
With  data cache, when a cache invalidate, the common solution is to get and update the cache. At the breaking load threshold, this is fine. But at high load thresholds, there are cache keys that can make a large number of requests to a DB in a short time (until a new cache is available). It is DB overload. </br>

This is easy to imagine when you have 10000 get commands almost simultaneously into 1 key cache, at the same time that key cache invalidate. If handled in the usual way, you will almost have 10000 requests to the DB at the same time. </br>



## Solution system high get data <a name="SolutionSystemHighGetData"></a>
## Solution get from cache <a name="SolutionGetFromCache"></a>
![](img_readme/get_data_from_cache.png)

In general, with billions of users and extremely large rqs, it is almost very limited to interact directly with the DB. Here, the receive will be taken from the cache. </br>

There exist 2 cases:
1) Data exist in cahce, get data and response </br>
2) Data not exits, call to update service and return the result. </br>

## Solution update cache <a name="SolutionUpdateCache"></a>
Service update cache must ensure: there are 1000 requests to update cache with the same cache key at the same time, only process one request, only perform 1 process query to DB and update cache and return data for that request. Other requests will have to wait until the next query. </br>

In this problem, the number of race conditions is not large. With a large number of race conditions, service get cache and service update cache must be isolated. They just transmit data to each other via the message queue. </br>

## Solution only handle one request update cache for one key in race conditions <a name="SolutionOnlyHandleOneRequestUpdateCacheForOneKeyInRaceConditions"></a>
![](img_readme/race_condition_update_cache.png)
I use mutext lock redis to handle it. With n update cache requests with 1 cache key, the first request to be handled will get the mutex and update the cache and return the result. Other requests will wait until the next update to get the latest results (at this point the cache is valid and the lock is also released). </br>

## Solution smart select DB <a name="SolutionSmartSelectDB"></a>
I have designed:  order, log, balance DB success to be moved from mysql to cassandra. Besides, Cassandra has much better load capacity with mysql. The solution is simply to always prioritize the query in cassandra first, if not exists, switch the query to mysql. </br>


