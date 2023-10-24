package pkg

type MainTable struct {
	ID                  uint `gorm:"primaryKey"`
	ReportingEntityName string
	ReportingEntityType string
	PlanName            string
	PlanIDType          string
	PlanID              string
	PlanMarketType      string
	LastUpdatedOn       string
	Version             string
	InNetwork           []InNetworkTable
}

type InNetworkTable struct {
	ID                     uint `gorm:"primaryKey"`
	MainTableID            uint
	NegotiationArrangement string
	Name                   string
	BillingCodeType        string
	BillingCodeTypeVersion string
	BillingCode            string
	Description            string
	NegotiatedRates        []NegotiatedRatesTable
}

type NegotiatedRatesTable struct {
	ID               uint `gorm:"primaryKey"`
	InNetworkTableID uint
	ProviderGroups   []ProviderGroupsTable
	NegotiatedPrices []NegotiatedPricesTable
}

type ProviderGroupsTable struct {
	ID                     uint `gorm:"primaryKey"`
	NegotiatedRatesTableID uint
	NPI                    []NpiTable
	TIN                    TinTable
}

type NpiTable struct {
	ID                    uint `gorm:"primaryKey"`
	ProviderGroupsTableID uint
	NpiValue              int
}

type TinTable struct {
	ID                    uint `gorm:"primaryKey"`
	ProviderGroupsTableID uint
	Type                  string
	Value                 string
}

type NegotiatedPricesTable struct {
	ID                     uint `gorm:"primaryKey"`
	NegotiatedRatesTableID uint
	NegotiatedType         string
	NegotiatedRate         float64
	ExpirationDate         string
	BillingClass           string
	ServiceCodes           []ServiceCodeTable
}

type ServiceCodeTable struct {
	ID                      uint `gorm:"primaryKey"`
	NegotiatedPricesTableID uint
	ServiceCode             string
}

type JSONData struct {
	ReportingEntityName string `json:"reporting_entity_name"`
	ReportingEntityType string `json:"reporting_entity_type"`
	PlanName            string `json:"plan_name"`
	PlanIDType          string `json:"plan_id_type"`
	PlanID              string `json:"plan_id"`
	PlanMarketType      string `json:"plan_market_type"`
	InNetwork           []struct {
		NegotiationArrangement string `json:"negotiation_arrangement"`
		Name                   string `json:"name"`
		BillingCodeType        string `json:"billing_code_type"`
		BillingCodeTypeVersion string `json:"billing_code_type_version"`
		BillingCode            string `json:"billing_code"`
		Description            string `json:"description"`
		NegotiatedRates        []struct {
			ProviderGroups []struct {
				NPI []int `json:"npi"`
				TIN struct {
					Type  string `json:"type"`
					Value string `json:"value"`
				} `json:"tin"`
			} `json:"provider_groups"`
			NegotiatedPrices []struct {
				NegotiatedType string   `json:"negotiated_type"`
				NegotiatedRate float64  `json:"negotiated_rate"`
				ExpirationDate string   `json:"expiration_date"`
				ServiceCode    []string `json:"service_code"`
				BillingClass   string   `json:"billing_class"`
			} `json:"negotiated_prices"`
		} `json:"negotiated_rates"`
	} `json:"in_network"`
	LastUpdatedOn string `json:"last_updated_on"`
	Version       string `json:"version"`
}
