# terraform-cost-guard

A CLI tool that wraps Terraform apply with cost estimation. It analyzes your terraform plan and prompts for confirmation before applying if costs will change.

```
Hey, these changes will cost an additional $2000.00/month. Proceed? [y/N]
```

## Installation

```bash
go install github.com/terraform-cost-guard/cmd/tfcost@latest
```

Or build from source:

```bash
git clone https://github.com/your-org/terraform-cost-guard
cd terraform-cost-guard
go build -o tfcost ./cmd/tfcost
```

## Usage

### Option 1: Wrap Command (Recommended)

Run plan and apply in one command with cost estimation:

```bash
tfcost wrap
```

This will:
1. Run `terraform plan -out=tfplan`
2. Convert the plan to JSON
3. Estimate monthly costs
4. Prompt for confirmation
5. Run `terraform apply tfplan` if confirmed

### Option 2: Apply with Pre-generated Plan

Generate a plan JSON first, then use tfcost:

```bash
# Generate terraform plan
terraform plan -out=tfplan

# Convert to JSON
terraform show -json tfplan > tfplan.json

# Run with cost guard
tfcost apply --plan tfplan.json
```

### Option 3: Estimate Only

Just see the cost estimate without applying:

```bash
terraform plan -out=tfplan
terraform show -json tfplan > tfplan.json
tfcost estimate tfplan.json
```

## Flags

### Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--verbose` | `-v` | Show detailed cost breakdown per resource |

### Apply/Wrap Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--plan` | `-p` | Path to terraform plan JSON file (required for apply) |
| `--threshold` | `-t` | Only prompt if cost exceeds threshold ($/month) |
| `--auto-approve` | `-y` | Skip confirmation prompt |

## Examples

### Basic usage

```bash
tfcost wrap
```

Output:
```
Running: terraform plan -out=/tmp/tfcost-xxx/tfplan

Terraform will perform the following actions:
  # aws_instance.web will be created
  ...

Converting plan to JSON...

============================================================
                    COST ESTIMATE SUMMARY
============================================================

  Resources to be created:   3
  Resources to be destroyed: 0
  Resources to be updated:   0

------------------------------------------------------------

  Estimated Monthly Cost Increase: +$156.95

============================================================

Hey, these changes will cost an additional $156.95/month. Proceed? [y/N]
```

### With threshold

Only prompt if monthly cost increase exceeds $100:

```bash
tfcost wrap --threshold 100
```

### Verbose output

Show per-resource cost breakdown:

```bash
tfcost wrap --verbose
```

Output includes:
```
  Detailed Cost Breakdown:
  Resource                                           Monthly Cost Details
  ----------------------------------------------------------------------
  aws_instance.web                                         $60.74 EC2 t3.large
  aws_db_instance.main                                     $63.51 RDS db.t3.medium + 100GB storage
  aws_nat_gateway.main                                     $32.85 NAT Gateway
```

### CI/CD Integration

Auto-approve with threshold for CI pipelines:

```bash
tfcost wrap --auto-approve --threshold 500
```

This will:
- Auto-approve if cost change is under $500/month
- Skip the interactive prompt entirely

## Supported Resources

### AWS
- EC2 Instances (`aws_instance`)
- RDS Instances (`aws_db_instance`)
- EBS Volumes (`aws_ebs_volume`)
- Application Load Balancer (`aws_lb`, `aws_alb`)
- Classic Load Balancer (`aws_elb`)
- NAT Gateway (`aws_nat_gateway`)
- ElastiCache (`aws_elasticache_cluster`)
- Lambda Functions (`aws_lambda_function`)
- S3 Buckets (`aws_s3_bucket`)
- EKS Clusters (`aws_eks_cluster`)
- ECS Services (`aws_ecs_service`)

### GCP
- Compute Instances (`google_compute_instance`)

### Azure
- Virtual Machines (`azurerm_virtual_machine`, `azurerm_linux_virtual_machine`, `azurerm_windows_virtual_machine`)

## Limitations

- Cost estimates are approximate and based on US region on-demand pricing
- Data transfer costs are not included
- Some resource types are not yet supported (will show as $0)
- Reserved instance pricing is not considered
- Spot/preemptible pricing is not considered

## How It Works

1. Parses the Terraform plan JSON output
2. Identifies resources being created, destroyed, updated, or replaced
3. Looks up approximate hourly/monthly rates from embedded pricing data
4. Calculates the net monthly cost change
5. Prompts for user confirmation before proceeding

## Contributing

To add support for new resource types, edit `internal/cost/estimator.go` and `internal/cost/pricing.go`.

## License

MIT
