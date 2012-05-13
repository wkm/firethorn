/*
	Firethorn configuration

	The goal is to parse JSON in this format:

	{ //- firethorn.Config 
		""
	}
*/

package main

type Config struct {
	// administration
	ListenAddress string
	PidFile       string
	LogDir        string

	// data schema
	Schemas map[string]Schema

	// redis configuration
	Storage Storage
}

type Schema struct {
	Dimensions map[string]Dimension
}

type Dimension struct {
	Id        uint
	Key       string
	Hierarchy string
	Enums     []string
}

type Storage struct {
	SamplingFactor int
	Redundancy     int
	Sharding       string
	Pools          []map[string]interface{}
}
