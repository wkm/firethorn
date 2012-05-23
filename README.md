# Firethorn
A high availability, scalable, redis-backed counter service written in Go that offers a tradeoff between read and write performance and (slight) inaccuracy.

The trick of firethorn is randomly choosing a node to accept an incr/decr operation. In this way write throughput scales linearly with the number of nodes in the pool. Reads can either be estimates: by reading from a single node and multiplying by the number of nodes in the pool, or reads can be near exact by reading from all nodes and summing. At the same time, counts are effectively replicated across multiple redis instances.

If a redis instance is lost, it can be directly replicated from one of its peers. In the interim Firethorn effectively samples by silently failing writes to the failed pool.

## Drawbacks
* Data jitter: multiple requests are going to give slightly jittered results (where the values are sometimes more, sometimes less). The jitter should be insignificant for all but the smallest counts. This is particularly an issue for historical data which isn't being modified anymore.
* Not "elastic": Firethorn does not in any way automatically scale up or down as machines are added to the cluster.
* Relatively expensive space wise:
    * the sharding scheme is conceptually similar to a RAID 0+1. Double the memory is required to maintain the same number of data points. This is natural with replication, the benefit of Firethorn is the resulting performance increase for reads and writes.
    * the OLAP data model increases the number of keys required per insert exponentially on the number of dimensions (the number of keys per insert is equal to the product of, for each dimension, the precision of the dimension plus one)


## Data Model
The basic datamodel is similar to Twitter's Rainbird hiearchical keys with the addition of multiple dimensions.

For example, to increment the number of likes for a particular post:

    time=2012/5/13/14/22
    client=123/123123123
    activity=likes

This lets us construct the following queries:

### Constructing Queries

All activity over all time for all clients: [empty query]

```

```

The count of all activity over all time for a client:

	client=456

Per-day counters for all activity for a client:

	client=456
	time=2012/5/0..31

Per-month counters for all like activity for all posts for a client in May and June 2012:

	client=456
	time=2012/5..6
	activity=likes

Per-day likes and reblog activity for a client in May 2012:

	client=456/123123
	time=2012/5/0..31
	activity=likes|reblogs

Per-day likes and reblog activity for a client in May, June, and July 2012:

	client=456/123123
	time=2012/5..7/0..31
	activity=likes|reblogs

Same, but for two posts:

	client=456/123123,789789
	time=2012/5..7/0..31
	activity=likes|reblogs

Same, but for four posts from two clients:

	client=456/123123,789789;489/1231,32452
	time=2012/5..7/0..31
	activity=likes|reblogs

You could remove the client constraint in any of these queries to get aggregates across all clients.

#### Query format
[need to flesh out]
* `---` to seperate multiple queries and inserts
* ` 

### Output Format
Output is in JSON containing:

* the tensor of counter data
* the retrieved axes ticks for each dimension
* meta data about the query

A sample output:

```json
{
	"schema": "activity",
	"millis": 1,
	"keysReferenced": 1,
	"instanceCount": 1,
	"queryCount": 1,
	"missings": 0,
	"error": {
		".0001": 5,
		".001": 1
	},
	"dimensions": {
		"client": [456]
	},
	"data": [123123]
}
```

* `schema` -- the namespace for this keyset
* `millis` -- time spent on the server aggregating the result set
* `keysReferenced` -- the number of keys which were queried
* `instanceCount` -- the number of redis instances reached
* `queryCount` -- the number of redis queries evaluated
* `missings` -- the number of missing data points
* `error` -- [XXXX]
* `dimensions` -- the axes ticks for each dimension of result
* `data` -- a multi-dimensional array of return values

## Configuration
Firethorn configuration is written in JSON as four components:

### Data Schemas
In theory the data schema could be implicitly derived from data insertions and queries. However, the subtleties of reconstructing keys and the exponential explosion in keys across dimensions suggest some rigor would be beneficial. So we require a data schema to be specified in the configuration:

```json
{
	"dimensions": {
		"time": {
			"id": 0,
			"key": "t",
			"schema": "#/#/#/#/#/#"
		},
		"client": {
			"id": 100,
			"key": "c",
			"schema": "#/#"
		},
		"activity": {
			"id": 200,
			"key": "a",
			"schema": ["likes", "follows", "reblogs", "views"]
		}
	}
}
```



### Redis Instances: Replication, Partitioning, and Sampling

The primary feature of the redis configuration is a specification of the individual pools.

```json
{
	"storage": {
		"errorRateThreshold": 15,
		"samplingfactor": 1,
		"sharding": "hashing",
		"pools": [
			{
				"redis01:6379" : {},
				"redis01:6380" : {},
				"redis01:6381" : {}
			},
			{
				"redis02:6379" : {},
				"redis02:6380" : {},
				"redis02:6381" : {}
			}
		]
	}
}
```

Firethorn's random selection model operates on pools. That is, each operation will randomly choose a pool to execute against. Within a pool, operations are sharded against the instances within that pool using the chosen sharding algorithm. (hashing, basically) It's possible to have heterogenous pools, but the complexity dosen't seem worth it.


### Administration

Finally, as a service there are a few configurable settings:

```json
{
	"pidfile": "/var/run/firethorn_01.pid",
	"logdir": "7 /var/log/firethorn_01/"
}
```


## API
### Writing
The fundamental operation in firethorn is an increment against multiple counters. This is executed as a POST request against a firethorn endpoint:

	curl '127.0.0.1:8000'

### Reading
Reading is executed as a GET operation against a firethorn endpoint:


## Development

You need Go 1.1.

Otherwise Firethorn is a normal golang project:

    go build
    go test
    go run

=)



## To-dos
* redis instances should be allowed to filter by dimensions; this will allow horizontal capacity scaling
* http stats endpoint ("gostrich")
* compress archived data