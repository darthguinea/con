# con
Will search and SSH into cloud servers, by tag or name

## Example:
```
darthguinea$ con cassandra
+-------------------------------+-------------+---------------------+------------+------------------------------+---------+------------+
|             NAME              | ENVIRONMENT |     INSTANCE ID     | PRIVATE IP |         LAUNCH TIME          |  STATE  |   REGION   |
+-------------------------------+-------------+---------------------+------------+------------------------------+---------+------------+
| s-cassandra-abcd-event-abe845 | production  | i-05xxxxxxxxxxxxx45 | 10.0.5.29  | Wed Aug 1 15:30:00 AEST 2018 | running | eu-west-1a |
+-------------------------------+-------------+---------------------+------------+------------------------------+---------+------------+
| s-cassandra-abcd-event-abe533 | production  | i-0c9xxxxxxxxxxx533 | 10.0.5.62  | Wed Aug 1 15:30:00 AEST 2018 | running | eu-west-1a |
+-------------------------------+-------------+---------------------+------------+------------------------------+---------+------------+
| s-cassandra-abcd-event-abe574 | production  | i-026xxxxxxxxxxe574 | 10.0.5.125 | Wed Aug 1 15:30:00 AEST 2018 | running | eu-west-1b |
+-------------------------------+-------------+---------------------+------------+------------------------------+---------+------------+
| s-cassandra-abcd-event-abe16d | production  | i-0fxxxxxxxxxxc016d | 10.0.5.117 | Wed Aug 1 15:28:01 AEST 2018 | running | eu-west-1b |
+-------------------------------+-------------+---------------------+------------+------------------------------+---------+------------+
```

## Connect:
```
darthguinea$ con i-05xxxxxxxxxxxxx45
```

or:
```
darthguinea$ con 10.0.5.29
```
