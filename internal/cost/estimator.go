package cost

import (
	"fmt"
	"strings"

	"github.com/ober/terraform-cost-guard/internal/plan"
)

// CostEstimate represents the estimated cost for a resource
type CostEstimate struct {
	ResourceAddress string
	ResourceType    string
	Action          string
	MonthlyCost     float64
	Details         string
}

// EstimationResult contains the total cost estimation results
type EstimationResult struct {
	Estimates           []CostEstimate
	TotalMonthlyCost    float64
	TotalMonthlyChange  float64 // positive = increase, negative = decrease
	CreatedResources    int
	DestroyedResources  int
	UpdatedResources    int
	UnsupportedTypes    []string
}

// Estimator calculates cost estimates for terraform plans
type Estimator struct {
	pricing *PricingData
}

// NewEstimator creates a new cost estimator
func NewEstimator() *Estimator {
	return &Estimator{
		pricing: NewDefaultPricing(),
	}
}

// Estimate calculates the cost impact of a terraform plan
func (e *Estimator) Estimate(p *plan.Plan) (*EstimationResult, error) {
	result := &EstimationResult{
		Estimates:        make([]CostEstimate, 0),
		UnsupportedTypes: make([]string, 0),
	}

	unsupportedSet := make(map[string]bool)

	for _, rc := range p.ResourceChanges {
		action := strings.Join(rc.Change.Actions, "+")

		// Skip no-op changes
		if action == "no-op" || action == "" {
			continue
		}

		estimate := CostEstimate{
			ResourceAddress: rc.Address,
			ResourceType:    rc.Type,
			Action:          action,
		}

		// Calculate cost based on action
		switch {
		case containsAction(rc.Change.Actions, "create") && !containsAction(rc.Change.Actions, "delete"):
			// New resource being created
			cost, details, supported := e.estimateResourceCost(rc.Type, rc.Change.After)
			if !supported && !unsupportedSet[rc.Type] {
				unsupportedSet[rc.Type] = true
				result.UnsupportedTypes = append(result.UnsupportedTypes, rc.Type)
			}
			estimate.MonthlyCost = cost
			estimate.Details = details
			result.TotalMonthlyChange += cost
			result.CreatedResources++

		case containsAction(rc.Change.Actions, "delete") && !containsAction(rc.Change.Actions, "create"):
			// Resource being destroyed
			cost, details, supported := e.estimateResourceCost(rc.Type, rc.Change.Before)
			if !supported && !unsupportedSet[rc.Type] {
				unsupportedSet[rc.Type] = true
				result.UnsupportedTypes = append(result.UnsupportedTypes, rc.Type)
			}
			estimate.MonthlyCost = -cost
			estimate.Details = details + " (removed)"
			result.TotalMonthlyChange -= cost
			result.DestroyedResources++

		case containsAction(rc.Change.Actions, "create") && containsAction(rc.Change.Actions, "delete"):
			// Resource being replaced
			oldCost, _, _ := e.estimateResourceCost(rc.Type, rc.Change.Before)
			newCost, details, supported := e.estimateResourceCost(rc.Type, rc.Change.After)
			if !supported && !unsupportedSet[rc.Type] {
				unsupportedSet[rc.Type] = true
				result.UnsupportedTypes = append(result.UnsupportedTypes, rc.Type)
			}
			estimate.MonthlyCost = newCost - oldCost
			estimate.Details = details + " (replaced)"
			result.TotalMonthlyChange += (newCost - oldCost)
			result.UpdatedResources++

		case containsAction(rc.Change.Actions, "update"):
			// In-place update
			oldCost, _, _ := e.estimateResourceCost(rc.Type, rc.Change.Before)
			newCost, details, supported := e.estimateResourceCost(rc.Type, rc.Change.After)
			if !supported && !unsupportedSet[rc.Type] {
				unsupportedSet[rc.Type] = true
				result.UnsupportedTypes = append(result.UnsupportedTypes, rc.Type)
			}
			estimate.MonthlyCost = newCost - oldCost
			estimate.Details = details + " (updated)"
			result.TotalMonthlyChange += (newCost - oldCost)
			result.UpdatedResources++
		}

		result.Estimates = append(result.Estimates, estimate)
	}

	result.TotalMonthlyCost = result.TotalMonthlyChange

	return result, nil
}

// estimateResourceCost returns the monthly cost for a resource type with given attributes
func (e *Estimator) estimateResourceCost(resourceType string, attrs map[string]interface{}) (float64, string, bool) {
	if attrs == nil {
		return 0, "no attributes", false
	}

	switch resourceType {
	// AWS EC2
	case "aws_instance":
		return e.estimateEC2Instance(attrs)

	// AWS RDS
	case "aws_db_instance":
		return e.estimateRDSInstance(attrs)

	// AWS EBS
	case "aws_ebs_volume":
		return e.estimateEBSVolume(attrs)

	// AWS ELB/ALB
	case "aws_lb", "aws_alb":
		return e.estimateALB(attrs)
	case "aws_elb":
		return e.estimateELB(attrs)

	// AWS NAT Gateway
	case "aws_nat_gateway":
		return e.estimateNATGateway(attrs)

	// AWS Elasticache
	case "aws_elasticache_cluster":
		return e.estimateElasticache(attrs)

	// AWS Lambda (compute time estimated)
	case "aws_lambda_function":
		return e.estimateLambda(attrs)

	// AWS S3
	case "aws_s3_bucket":
		return e.estimateS3Bucket(attrs)

	// AWS EKS
	case "aws_eks_cluster":
		return e.estimateEKSCluster(attrs)

	// AWS ECS
	case "aws_ecs_service":
		return e.estimateECSService(attrs)

	// GCP Compute
	case "google_compute_instance":
		return e.estimateGCPInstance(attrs)

	// Azure VM
	case "azurerm_virtual_machine", "azurerm_linux_virtual_machine", "azurerm_windows_virtual_machine":
		return e.estimateAzureVM(attrs)

	default:
		return 0, "unsupported resource type", false
	}
}

func (e *Estimator) estimateEC2Instance(attrs map[string]interface{}) (float64, string, bool) {
	instanceType := getStringAttr(attrs, "instance_type", "t3.micro")
	hourlyRate := e.pricing.EC2Instances[instanceType]
	if hourlyRate == 0 {
		hourlyRate = e.pricing.EC2Instances["t3.micro"] // fallback
	}
	monthlyCost := hourlyRate * 730 // average hours per month
	return monthlyCost, fmt.Sprintf("EC2 %s", instanceType), true
}

func (e *Estimator) estimateRDSInstance(attrs map[string]interface{}) (float64, string, bool) {
	instanceClass := getStringAttr(attrs, "instance_class", "db.t3.micro")
	hourlyRate := e.pricing.RDSInstances[instanceClass]
	if hourlyRate == 0 {
		hourlyRate = e.pricing.RDSInstances["db.t3.micro"]
	}

	// Add storage cost
	storageGB := getFloat64Attr(attrs, "allocated_storage", 20)
	storageCost := storageGB * e.pricing.EBSStorage["gp2"]

	monthlyCost := (hourlyRate * 730) + storageCost
	return monthlyCost, fmt.Sprintf("RDS %s + %.0fGB storage", instanceClass, storageGB), true
}

func (e *Estimator) estimateEBSVolume(attrs map[string]interface{}) (float64, string, bool) {
	volumeType := getStringAttr(attrs, "type", "gp2")
	sizeGB := getFloat64Attr(attrs, "size", 8)
	rate := e.pricing.EBSStorage[volumeType]
	if rate == 0 {
		rate = e.pricing.EBSStorage["gp2"]
	}
	monthlyCost := sizeGB * rate
	return monthlyCost, fmt.Sprintf("EBS %s %.0fGB", volumeType, sizeGB), true
}

func (e *Estimator) estimateALB(attrs map[string]interface{}) (float64, string, bool) {
	// ALB has hourly cost + LCU charges (we estimate base cost only)
	monthlyCost := e.pricing.LoadBalancers["alb"] * 730
	return monthlyCost, "Application Load Balancer", true
}

func (e *Estimator) estimateELB(attrs map[string]interface{}) (float64, string, bool) {
	monthlyCost := e.pricing.LoadBalancers["classic"] * 730
	return monthlyCost, "Classic Load Balancer", true
}

func (e *Estimator) estimateNATGateway(attrs map[string]interface{}) (float64, string, bool) {
	// NAT Gateway hourly charge (data processing extra)
	monthlyCost := e.pricing.NATGateway * 730
	return monthlyCost, "NAT Gateway", true
}

func (e *Estimator) estimateElasticache(attrs map[string]interface{}) (float64, string, bool) {
	nodeType := getStringAttr(attrs, "node_type", "cache.t3.micro")
	numNodes := getFloat64Attr(attrs, "num_cache_nodes", 1)
	hourlyRate := e.pricing.Elasticache[nodeType]
	if hourlyRate == 0 {
		hourlyRate = e.pricing.Elasticache["cache.t3.micro"]
	}
	monthlyCost := hourlyRate * 730 * numNodes
	return monthlyCost, fmt.Sprintf("Elasticache %s x%.0f", nodeType, numNodes), true
}

func (e *Estimator) estimateLambda(attrs map[string]interface{}) (float64, string, bool) {
	// Lambda pricing is complex (requests + duration), estimate minimal
	memoryMB := getFloat64Attr(attrs, "memory_size", 128)
	// Rough estimate: 1M requests/month at 100ms each
	monthlyCost := (memoryMB / 1024) * 0.0000166667 * 100 * 1000000 / 1000
	return monthlyCost, fmt.Sprintf("Lambda %0.fMB (estimated)", memoryMB), true
}

func (e *Estimator) estimateS3Bucket(attrs map[string]interface{}) (float64, string, bool) {
	// S3 cost depends on storage used - estimate minimal for bucket creation
	return 0.023, "S3 Bucket (minimal estimate)", true
}

func (e *Estimator) estimateEKSCluster(attrs map[string]interface{}) (float64, string, bool) {
	// EKS cluster has flat hourly rate
	monthlyCost := e.pricing.EKSCluster * 730
	return monthlyCost, "EKS Cluster", true
}

func (e *Estimator) estimateECSService(attrs map[string]interface{}) (float64, string, bool) {
	// ECS itself is free, cost is in underlying EC2/Fargate
	// Estimate based on desired count if using Fargate
	desiredCount := getFloat64Attr(attrs, "desired_count", 1)
	// Rough Fargate estimate (0.25 vCPU, 0.5GB)
	monthlyCost := desiredCount * (0.25*0.04048 + 0.5*0.004445) * 730
	return monthlyCost, fmt.Sprintf("ECS Service (%.0f tasks, Fargate estimate)", desiredCount), true
}

func (e *Estimator) estimateGCPInstance(attrs map[string]interface{}) (float64, string, bool) {
	machineType := getStringAttr(attrs, "machine_type", "e2-micro")
	hourlyRate := e.pricing.GCPInstances[machineType]
	if hourlyRate == 0 {
		hourlyRate = e.pricing.GCPInstances["e2-micro"]
	}
	monthlyCost := hourlyRate * 730
	return monthlyCost, fmt.Sprintf("GCP %s", machineType), true
}

func (e *Estimator) estimateAzureVM(attrs map[string]interface{}) (float64, string, bool) {
	size := getStringAttr(attrs, "size", "Standard_B1s")
	if size == "" {
		size = getStringAttr(attrs, "vm_size", "Standard_B1s")
	}
	hourlyRate := e.pricing.AzureVMs[size]
	if hourlyRate == 0 {
		hourlyRate = e.pricing.AzureVMs["Standard_B1s"]
	}
	monthlyCost := hourlyRate * 730
	return monthlyCost, fmt.Sprintf("Azure %s", size), true
}

func containsAction(actions []string, target string) bool {
	for _, a := range actions {
		if a == target {
			return true
		}
	}
	return false
}

func getStringAttr(attrs map[string]interface{}, key, defaultVal string) string {
	if v, ok := attrs[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultVal
}

func getFloat64Attr(attrs map[string]interface{}, key string, defaultVal float64) float64 {
	if v, ok := attrs[key]; ok {
		switch n := v.(type) {
		case float64:
			return n
		case int:
			return float64(n)
		case int64:
			return float64(n)
		}
	}
	return defaultVal
}
