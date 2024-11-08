# trx-reconcile-app
Sample system of trx reconciliation written in Go


# Notes
- Code written using clean architecture, testing would be really easy using dependency injection / mocking on interfaces
- Transaction ID will have this pattern:
    [1-4 UNIQUE BANK CODE][5-36 TRX UUID]
- As for MVP service, all csv data will be load in memory
    Future plan: as data grow which can burden RAM, we can introduce batching mechanism
- csv format will put timestamp as first column
- security mechanism on filepath param not included


# Running
```
./run.sh
```

Sample Request
```
http://localhost:9000/reconcile/trxtobanks?start_time=2024-11-06%2000%3A00%3A00&end_time=2024-11-10%2000%3A00%3A00&system_trx=data%2Ftransaction%2Ftrx.csv&banks_stats_trx=data%2Fbank_ini%2Ftrx.csv,data%2Fbank_itu%2Ftrx.csv
```