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
1) Part one is the part of managing partner sharding, it will regulate and manage which partners exist on which sharding </br>
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



