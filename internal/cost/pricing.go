package cost

// PricingData contains hourly/monthly rates for various cloud resources
// Prices are approximate US East region on-demand pricing (USD)
type PricingData struct {
	// AWS EC2 instance types -> hourly rate
	EC2Instances map[string]float64

	// AWS RDS instance classes -> hourly rate
	RDSInstances map[string]float64

	// AWS EBS volume types -> per GB/month
	EBSStorage map[string]float64

	// AWS Load Balancers -> hourly rate
	LoadBalancers map[string]float64

	// NAT Gateway hourly rate
	NATGateway float64

	// AWS Elasticache node types -> hourly rate
	Elasticache map[string]float64

	// EKS cluster hourly rate
	EKSCluster float64

	// GCP machine types -> hourly rate
	GCPInstances map[string]float64

	// Azure VM sizes -> hourly rate
	AzureVMs map[string]float64
}

// NewDefaultPricing returns pricing data with approximate current rates
func NewDefaultPricing() *PricingData {
	return &PricingData{
		EC2Instances: map[string]float64{
			// General Purpose
			"t3.nano":     0.0052,
			"t3.micro":    0.0104,
			"t3.small":    0.0208,
			"t3.medium":   0.0416,
			"t3.large":    0.0832,
			"t3.xlarge":   0.1664,
			"t3.2xlarge":  0.3328,
			"t3a.nano":    0.0047,
			"t3a.micro":   0.0094,
			"t3a.small":   0.0188,
			"t3a.medium":  0.0376,
			"t3a.large":   0.0752,
			"t3a.xlarge":  0.1504,
			"t3a.2xlarge": 0.3008,
			"m5.large":    0.096,
			"m5.xlarge":   0.192,
			"m5.2xlarge":  0.384,
			"m5.4xlarge":  0.768,
			"m5.8xlarge":  1.536,
			"m5.12xlarge": 2.304,
			"m5.16xlarge": 3.072,
			"m5.24xlarge": 4.608,
			"m6i.large":   0.096,
			"m6i.xlarge":  0.192,
			"m6i.2xlarge": 0.384,
			"m6i.4xlarge": 0.768,
			// Compute Optimized
			"c5.large":    0.085,
			"c5.xlarge":   0.17,
			"c5.2xlarge":  0.34,
			"c5.4xlarge":  0.68,
			"c5.9xlarge":  1.53,
			"c5.18xlarge": 3.06,
			"c6i.large":   0.085,
			"c6i.xlarge":  0.17,
			"c6i.2xlarge": 0.34,
			// Memory Optimized
			"r5.large":    0.126,
			"r5.xlarge":   0.252,
			"r5.2xlarge":  0.504,
			"r5.4xlarge":  1.008,
			"r5.8xlarge":  2.016,
			"r5.12xlarge": 3.024,
			// GPU Instances
			"p3.2xlarge":  3.06,
			"p3.8xlarge":  12.24,
			"p3.16xlarge": 24.48,
			"g4dn.xlarge": 0.526,
			"g4dn.2xlarge": 0.752,
			"g4dn.4xlarge": 1.204,
		},

		RDSInstances: map[string]float64{
			"db.t3.micro":    0.017,
			"db.t3.small":    0.034,
			"db.t3.medium":   0.068,
			"db.t3.large":    0.136,
			"db.t3.xlarge":   0.272,
			"db.t3.2xlarge":  0.544,
			"db.m5.large":    0.171,
			"db.m5.xlarge":   0.342,
			"db.m5.2xlarge":  0.684,
			"db.m5.4xlarge":  1.368,
			"db.r5.large":    0.24,
			"db.r5.xlarge":   0.48,
			"db.r5.2xlarge":  0.96,
			"db.r5.4xlarge":  1.92,
		},

		EBSStorage: map[string]float64{
			"gp2":      0.10,  // per GB/month
			"gp3":      0.08,
			"io1":      0.125,
			"io2":      0.125,
			"st1":      0.045,
			"sc1":      0.015,
			"standard": 0.05,
		},

		LoadBalancers: map[string]float64{
			"alb":     0.0225,
			"nlb":     0.0225,
			"classic": 0.025,
		},

		NATGateway: 0.045,

		Elasticache: map[string]float64{
			"cache.t3.micro":   0.017,
			"cache.t3.small":   0.034,
			"cache.t3.medium":  0.068,
			"cache.m5.large":   0.156,
			"cache.m5.xlarge":  0.312,
			"cache.m5.2xlarge": 0.624,
			"cache.r5.large":   0.226,
			"cache.r5.xlarge":  0.452,
		},

		EKSCluster: 0.10, // per hour

		GCPInstances: map[string]float64{
			"e2-micro":      0.0084,
			"e2-small":      0.0168,
			"e2-medium":     0.0336,
			"e2-standard-2": 0.0672,
			"e2-standard-4": 0.1344,
			"e2-standard-8": 0.2688,
			"n1-standard-1": 0.0475,
			"n1-standard-2": 0.095,
			"n1-standard-4": 0.19,
			"n1-standard-8": 0.38,
			"n2-standard-2": 0.0971,
			"n2-standard-4": 0.1942,
			"n2-standard-8": 0.3884,
		},

		AzureVMs: map[string]float64{
			"Standard_B1s":   0.0104,
			"Standard_B1ms":  0.0207,
			"Standard_B2s":   0.0416,
			"Standard_B2ms":  0.0832,
			"Standard_D2s_v3": 0.096,
			"Standard_D4s_v3": 0.192,
			"Standard_D8s_v3": 0.384,
			"Standard_E2s_v3": 0.126,
			"Standard_E4s_v3": 0.252,
			"Standard_E8s_v3": 0.504,
			"Standard_F2s_v2": 0.085,
			"Standard_F4s_v2": 0.169,
			"Standard_F8s_v2": 0.338,
		},
	}
}
