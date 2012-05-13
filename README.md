# Firethorn
A high availability, scalable, redis-backed counter service written in Go that offers a tradeoff between read and write performance and (slight) inaccuracy.

The trick of firethorn is randomly choosing a node to accept a counter operation. In this way write throughput scales linearly with the number of nodes in the pool. Reads can either be estimates: by reading from a single node and multiplying by the number of nodes in the pool, or reads can be near exact by reading from all nodes and summing. At the same time, counts are effectively replicated across multiple redis instances.

## Drawbacks
* Data jitter: multiple requests are going to give slightly jittered results (where the values are sometimes more, sometimes less). This is particularly an issue for historical data which isn"t being modified anymore.
* Not "elastic": 


## Data Model
The basic datamodel is similar to Twitter"s Rainbird hiearchical keys. However, Firethorn also supports multi-dimensional keys.

For example, to increment the number of likes for a particular post:

    time:2012/5/13/14/22
    client:123/123123123
    activity:likes

This lets us construct the following queries:

### Constructing Queries


The count of all activity over all time for a client:

	client=456

Per-day counters for all activity for a client:

	client=456
	time=2012/5/0..31

Per-month counters for all like activity for a client in May and June 2012:

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

### Output Format
Output is in JSON containing:

* the tensor of counter data
* the retrieved axes ticks for each dimension
* meta data about the query

A sample output:

```json
{
	"dimensions": {
		"client": [456]
	},
	"data": [123123],
	"stats": {
		"millis": 12,
		"keys": 1231,
		"instances": xxxx
	}
}
```

## Configuration
There are four aspects to configuring Firethorn:

### Data Schema
In theory the data schema could be implicitly derived from data insertions and queries. However, the subtleties of reconstructing keys and the exponential explosion in keys as dimensions are added seem like potential problems. So we require a data schema to be specified in the configuration:

```json
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
		"activity: {
			"id": 200,
			"key": "a",
			"schema": ["likes", "follows", "reblogs", "views"]
		}
	}
```



### Redis Instances: Replication, Partitioning, and Sampling

```json
	"storage": {
		"samplingfactor": 1,
		"redundancy": 2,
		"sharding": "hashing",
		"pools": [
			[
				"redis01:6379" : {},
				"redis01:6380" : {},
				"redis01:6381" : {}
			],
			[
				"redis02:6379" : {},
				"redis02:6380" : {},
				"redis02:6381" : {}
			]
		]
	}
```


### Administration

Finally, as a service there are a few configurable things:

```json
	"pidfile": "/var/run/firethorn_01.pid",
	"logdir": "/var/log/firethorn_01/"
```


## API
### Writing
The fundamental operation in firethorn is an increment against multiple counters.

### Reading


## Development

You need Go 1.1.

Otherwise Firethorn is a normal golang project:

    go build
    go test
    go run

=)



## To-dos
* http stats endpoint ("gostrich")
* archived data compression