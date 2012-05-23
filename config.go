/*
	Firethorn configuration

	The goal is to parse JSON in this format:

	{ //- firethorn.Config 
		""
	}
*/

package firethorn

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

type Storage struct {
	SamplingFactor     int
	ErrorRateThreshold int
	Sharding           string
	Pools              []map[string]interface{}
}
