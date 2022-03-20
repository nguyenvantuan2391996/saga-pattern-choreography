package grpc

import (
	"context"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	stock "saga-pattern-choreography/stock-service/api"
)

type Config struct {
	GRPCPort            string
	DatastoreDBHost     string
	DatastoreDBUser     string
	DatastoreDBPassword string
	DatastoreDBSchema   string
}

// RunServer runs gRPC server and HTTP gateway
func RunServerCMD() error {
	ctx := context.Background()

	// get configuration
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "3002", "gRPC port to bind")
	flag.StringVar(&cfg.DatastoreDBHost, "db-host", "127.0.0.1:3306", "Database host")
	flag.StringVar(&cfg.DatastoreDBUser, "db-user", "root", "Database user")
	flag.StringVar(&cfg.DatastoreDBPassword, "db-password", "root", "Database password")
	flag.StringVar(&cfg.DatastoreDBSchema, "db-schema", "stock_service", "Database schema")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	// add MySQL driver specific parameter to parse date/time
	// Drop it for another database
	param := "parseTime=true"

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		cfg.DatastoreDBUser,
		cfg.DatastoreDBPassword,
		cfg.DatastoreDBHost,
		cfg.DatastoreDBSchema,
		param)

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	v1API := stock.NewStockServiceServer(db)

	return RunServer(ctx, v1API, cfg.GRPCPort)
}
