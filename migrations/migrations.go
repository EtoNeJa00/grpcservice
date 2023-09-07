package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"GRPCService/config"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/migrate"
)

const migrationKeySpaceName = "migration"
const path = "/migrations/cql"

func main() {
	ctx := context.Background()

	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	cluster := gocql.NewCluster(conf.ScyllaAddr)
	cluster.Consistency = gocql.Quorum

	err = createMigrationKeyspace(ctx, cluster)
	if err != nil {
		log.Fatal(err)
	}

	cluster.Keyspace = migrationKeySpaceName

	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		log.Fatal(err)
	}

	defer session.Close()

	currentDir, err := os.Getwd()
	if err != nil {
		log.Println(err)

		return
	}

	if err = migrate.FromFS(ctx, session, os.DirFS(filepath.Join(currentDir, path))); err != nil {
		log.Println(err)

		return
	}

	log.Println("migration successfully up")
}

func createMigrationKeyspace(ctx context.Context, cluster *gocql.ClusterConfig) error {
	ses, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}

	defer ses.Close()

	err = ses.Query(fmt.Sprintf("create keyspace if not exists %s with replication = {'class': 'SimpleStrategy','replication_factor': 1};", migrationKeySpaceName)).WithContext(ctx).Exec()
	if err != nil {
		return err
	}

	return nil
}
