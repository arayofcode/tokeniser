package common

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

type testDbConfig struct {
	host     string
	port     string
	user     string
	db       string
	password string
}

var (
	initOnce    sync.Once
	testNetwork *testcontainers.DockerNetwork
)

type Log struct {
	LogType string
	Content []byte
}

type LogConsumer interface {
	Accept(Log)
}

type StdoutLogConsumer struct{}

func (lc *StdoutLogConsumer) Accept(l Log) {
	log.Info().Msg(string(l.Content))
}

func StartPostgres(ctx context.Context) (config testDbConfig, err error) {
	initOnce.Do(func() {
		testNetwork, err = NewTestcontainerNetwork()
		if err != nil {
			log.Fatal().Err(err).Msg("Could not create a new network")
			return
		}
		config, err = initPostgresContainer(ctx, testNetwork)
		if err != nil {
			log.Fatal().Err(err).Msg("Could not start postgres container")
			return
		}
	})
	return
}

func initPostgresContainer(ctx context.Context, newNetwork *testcontainers.DockerNetwork) (config testDbConfig, err error) {
	config = testDbConfig{
		user:     "test_user",
		password: "test_password",
		db:       "test_db",
	}

	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16.2-alpine3.19"),
		postgres.WithDatabase(config.db),
		postgres.WithUsername(config.user),
		postgres.WithPassword(config.password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
		testcontainers.WithLogConsumers(&testcontainers.StdoutLogConsumer{}),
		network.WithNetwork([]string{"postgres-flyway-network"}, newNetwork),
	)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to start container")
		return
	}

	config.host, _ = postgresContainer.Host(ctx)
	portData, _ := postgresContainer.MappedPort(ctx, "5432")
	config.port = portData.Port()
	log.Info().Str("Host", config.host).Str("Port", config.port).Msg("Created postgres container")

	// Clean up the container
	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatal().Err(err).Msg("failed to terminate container")
			return
		}
	}()
	return
}

func NewTestcontainerNetwork() (net *testcontainers.DockerNetwork, err error) {
	ctx := context.Background()

	net, err = network.New(ctx,
		network.WithCheckDuplicate(),
		network.WithAttachable(),
		network.WithInternal(),
	)

	log.Info().Str("name", net.Name).Str("ID", net.ID).Str("Driver", net.Driver).Msg("New network created")

	if err != nil {
		log.Fatal().Err(err).Msg("Error while creating network")
		return
	}

	defer func() {
		if err := net.Remove(ctx); err != nil {
			log.Fatal().Err(err).Msg("Failed to remove network")
		}
		log.Info().Str("ID", net.ID).Str("name", net.Name).Msg("Removing network")
	}()

	return
}
