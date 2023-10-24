package cli

import (
	"fmt"
	"github.com/purvisb179/tim/pkg"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	initDBCmd = &cobra.Command{
		Use:   "initdb",
		Short: "Initialize the database tables",
		Long:  "Connects to the PostgreSQL database and sets up the required tables based on the provided models.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return InitDatabase()
		},
	}
)

func InitDatabase() error {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return fmt.Errorf("Failed to connect to database: %s", err.Error())
	}

	// AutoMigrate will create the tables
	err = db.AutoMigrate(&pkg.MainTable{}, &pkg.InNetworkTable{}, &pkg.NegotiatedRatesTable{}, &pkg.ProviderGroupsTable{}, &pkg.NpiTable{}, &pkg.TinTable{}, &pkg.NegotiatedPricesTable{}, &pkg.ServiceCodeTable{})
	if err != nil {
		return fmt.Errorf("Failed to migrate the tables: %s", err.Error())
	}

	fmt.Println("Database initialized successfully!")
	return nil
}

func GetInitDBCmd() *cobra.Command {
	return initDBCmd
}
