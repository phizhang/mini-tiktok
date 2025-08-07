package db

import (
	"fmt"
	"github.com/gocql/gocql"
)

var Session *gocql.Session

// InitCassandra initializes the Cassandra session
func InitCassandra(hosts []string, keyspace string) error {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("failed to connect to Cassandra: %w", err)
	}
	Session = session
	return nil
}

// CloseCassandra closes the Cassandra session
func CloseCassandra() {
	if Session != nil {
		Session.Close()
	}
}
