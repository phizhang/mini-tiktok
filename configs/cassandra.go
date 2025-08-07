package configs

// CassandraConfig holds configuration for Cassandra connection
 type CassandraConfig struct {
	Hosts    []string
	Keyspace string
}
