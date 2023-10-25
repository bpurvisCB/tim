package cli

import (
	"encoding/json"
	"fmt"
	"github.com/purvisb179/tim/pkg"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"time"
)

var (
	uploadJSONCmd = &cobra.Command{
		Use:   "uploadjson",
		Short: "Upload data from JSON file to the database",
		Long:  "Reads the content of the specified JSON file and uploads it to the PostgreSQL database.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return uploadJSONToDatabase()
		},
	}
)

func uploadJSONToDatabase() error {
	startTime := time.Now() // Record the start time

	// 1. Read the JSON file
	filePath := "./2023-10-01_TheAlliance_AlliancePremierNetwork_in-network-rates.json"
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("Error reading JSON file: %s", err.Error())
	}

	// 2. Parse the content into Go structures
	var data pkg.JSONData
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		return fmt.Errorf("Error parsing JSON: %s", err.Error())
	}

	// 3. Connect to the PostgreSQL database using GORM
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("Failed to connect to database: %s", err.Error())
	}

	// 5. Insert the data into the respective tables
	mainEntry := pkg.MainTable{
		ReportingEntityName: data.ReportingEntityName,
		ReportingEntityType: data.ReportingEntityType,
		PlanName:            data.PlanName,
		PlanIDType:          data.PlanIDType,
		PlanID:              data.PlanID,
		PlanMarketType:      data.PlanMarketType,
		LastUpdatedOn:       data.LastUpdatedOn,
		Version:             data.Version,
	}
	result := db.Create(&mainEntry)
	if result.Error != nil {
		fmt.Printf("Error inserting into MainTable: %s\n", result.Error)
	}

	// 6. Loop through in_network entries
	for _, inNetwork := range data.InNetwork {
		inNetworkEntry := pkg.InNetworkTable{
			// Fill fields from inNetwork
			NegotiationArrangement: inNetwork.NegotiationArrangement,
			Name:                   inNetwork.Name,
			BillingCodeType:        inNetwork.BillingCodeType,
			BillingCodeTypeVersion: inNetwork.BillingCodeTypeVersion,
			BillingCode:            inNetwork.BillingCode,
			Description:            inNetwork.Description,
			// Assuming you have a foreign key to reference MainTable
			MainTableID: mainEntry.ID,
		}
		result := db.Create(&inNetworkEntry)
		if result.Error != nil {
			fmt.Printf("Error inserting into InNetworkTable: %s\n", result.Error)
			continue // Skip processing this inNetwork entry and move to the next one
		}

		// 7. For each in_network entry, loop through negotiated_rates
		for _, negotiatedRate := range inNetwork.NegotiatedRates {
			negotiatedRateEntry := pkg.NegotiatedRatesTable{
				// TODO: Fill fields from negotiatedRate if there are any extra fields
				InNetworkTableID: inNetworkEntry.ID,
			}
			result := db.Create(&negotiatedRateEntry)
			if result.Error != nil {
				fmt.Printf("Error inserting into NegotiatedRatesTable: %s\n", result.Error)
				continue // Skip this negotiatedRate entry and move to the next one
			}

			// Inserting provider_groups
			for _, providerGroup := range negotiatedRate.ProviderGroups {
				providerGroupEntry := pkg.ProviderGroupsTable{
					// Assuming you have a foreign key to reference NegotiatedRatesTable
					NegotiatedRatesTableID: negotiatedRateEntry.ID,
				}
				result := db.Create(&providerGroupEntry)
				if result.Error != nil {
					fmt.Printf("Error inserting into ProviderGroupsTable: %s\n", result.Error)
					continue // Skip this providerGroup and move to the next one
				}

				// Inserting npi entries
				for _, npi := range providerGroup.NPI {
					npiEntry := pkg.NpiTable{
						NpiValue: npi,
						// Assuming you have a foreign key to reference ProviderGroupsTable
						ProviderGroupsTableID: providerGroupEntry.ID,
					}
					result := db.Create(&npiEntry)
					if result.Error != nil {
						fmt.Printf("Error inserting into NpiTable: %s\n", result.Error)
						continue
					}
				}

				// Inserting tin
				tinEntry := pkg.TinTable{
					Type:  providerGroup.TIN.Type,
					Value: providerGroup.TIN.Value,
					// Assuming you have a foreign key to reference ProviderGroupsTable
					ProviderGroupsTableID: providerGroupEntry.ID,
				}
				result = db.Create(&tinEntry)
				if result.Error != nil {
					fmt.Printf("Error inserting into TinTable: %s\n", result.Error)
				}
			}

			// Inserting negotiated_prices entries
			for _, negotiatedPrice := range negotiatedRate.NegotiatedPrices {
				negotiatedPriceEntry := pkg.NegotiatedPricesTable{
					NegotiatedType: negotiatedPrice.NegotiatedType,
					NegotiatedRate: negotiatedPrice.NegotiatedRate,
					ExpirationDate: negotiatedPrice.ExpirationDate,
					BillingClass:   negotiatedPrice.BillingClass,
					// Assuming you have a foreign key to reference NegotiatedRatesTable
					NegotiatedRatesTableID: negotiatedRateEntry.ID,
				}
				result := db.Create(&negotiatedPriceEntry)
				if result.Error != nil {
					fmt.Printf("Error inserting into NegotiatedPricesTable: %s\n", result.Error)
					continue
				}

				// Inserting service_code entries
				for _, serviceCode := range negotiatedPrice.ServiceCode {
					serviceCodeEntry := pkg.ServiceCodeTable{
						ServiceCode: serviceCode,
						// Assuming you have a foreign key to reference NegotiatedPricesTable
						NegotiatedPricesTableID: negotiatedPriceEntry.ID,
					}
					result := db.Create(&serviceCodeEntry)
					if result.Error != nil {
						fmt.Printf("Error inserting into ServiceCodeTable: %s\n", result.Error)
						continue
					}
				}
			}

		}
	}

	elapsedTime := time.Since(startTime)              // Calculate the elapsed time
	fmt.Printf("Total time taken: %s\n", elapsedTime) // Print the elapsed time

	return nil
}

func GetUploadJSONCmd() *cobra.Command {
	return uploadJSONCmd
}
