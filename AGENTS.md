# AGENTS.md

This file provides guidance to AI agents (Claude Code, Cursor, Copilot, etc.) when working with code in this repository.

## Build and Development Commands

```bash
# Build and install
make build              # Compile provider to $GOPATH/bin
make install            # Build and install to local Terraform plugins directory
make uninstall          # Remove from local plugins directory

# Code quality
make fmt                # Format code with gofumpt
make fmtcheck           # Verify code formatting
make lint               # Run golangci-lint (1h timeout)

# Testing
make test               # Run unit tests (30s timeout, 4 parallel)
make testacc            # Run acceptance tests (120min timeout, requires credentials)

# Run specific test suite
TEST_SUITE=tf_acc_sysdig_monitor make testacc
TEST_SUITE=tf_acc_sysdig_secure make testacc
TEST_SUITE=tf_acc_ibm_monitor make testacc
TEST_SUITE=tf_acc_ibm_secure make testacc

# Run a single test
go test ./sysdig -v -run TestAccResourceSysdigUser -tags=tf_acc_sysdig_secure -timeout 120m

# Run with debug logging
TF_ACC=1 TF_LOG=DEBUG go test ./sysdig -v -tags=tf_acc_sysdig_secure -run TestAccMacro
```

## Required Environment Variables

For acceptance tests (set in `.envrc` or `.env`):
```bash
SYSDIG_SECURE_API_TOKEN     # Required for Secure tests
SYSDIG_MONITOR_API_TOKEN    # Required for Monitor tests
TF_ACC=1                    # Enable acceptance testing

# IBM Cloud (optional)
SYSDIG_IBM_MONITOR_API_KEY
SYSDIG_IBM_MONITOR_INSTANCE_ID
SYSDIG_IBM_SECURE_API_KEY
SYSDIG_IBM_SECURE_INSTANCE_ID
```

## Architecture Overview

Terraform provider for Sysdig Monitor and Sysdig Secure using Terraform Plugin SDK v2.

### Directory Structure

- `sysdig/` - Main provider package (resources, data sources, provider config)
- `sysdig/internal/client/v2/` - HTTP client layer for Sysdig API
- `website/docs/` - Terraform Registry documentation

### File Naming Conventions

- Resources: `sysdig/resource_sysdig_<service>_<entity>.go`
- Data sources: `sysdig/data_source_sysdig_<service>_<entity>.go`
- Tests: `sysdig/<source>_test.go` (same package, co-located)
- Client implementations: `sysdig/internal/client/v2/<entity>.go`

### Key Reference Files

When implementing new resources, reference these files for patterns and utilities:

| File | Purpose |
|------|---------|
| `sysdig/common.go` | ~60 predefined Schema key constants (`SchemaIDKey`, `SchemaTeamIDKey`, etc.) |
| `sysdig/schema.go` | Schema composition helpers, read-only schemas, `maps.Copy()` pattern |
| `sysdig/helpers.go` | Utility functions (`validateDiagFunc`, `parseAzureCreds`) |
| `sysdig/sysdig_clients.go` | Client interface and selection logic |
| `sysdig/resource_sysdig_secure_common_policy.go` | Secure policy patterns |
| `sysdig/resource_sysdig_monitor_alert_v2_common.go` | Alert V2 base implementation |
| `sysdig/resource_sysdig_monitor_notification_channel_common.go` | Notification channel patterns |

## Resource Implementation Patterns

### Schema Composition

Resources use composable schema functions for shared fields:

```go
// Alert resources compose multiple schema layers
Schema: createScopedSegmentedAlertV2Schema(
    createAlertV2Schema(
        map[string]*schema.Schema{
            "type_specific_field": { ... }
        }
    )
)

// Notification channels use shared base schema
Schema: createMonitorNotificationChannelSchema(
    map[string]*schema.Schema{
        "channel_specific_field": { ... }
    }
)
```

**Base schema functions** (in `sysdig/`):
- `createAlertV2Schema()` - Alert common fields
- `createScopedSegmentedAlertV2Schema()` - Scope and group_by
- `createRuleSchema()` - Secure rule common fields
- `createMonitorNotificationChannelSchema()` - Notification channel base
- `createSecureNotificationChannelSchema()` - Secure notification base

### CRUD Implementation Pattern

```go
func resourceSysdigXxxCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
    client, err := getXxxClient(meta.(SysdigClients))
    if err != nil {
        return diag.FromErr(err)
    }

    obj := buildXxxStruct(d)  // ResourceData → Go struct

    created, err := client.CreateXxx(ctx, obj)
    if err != nil {
        return diag.FromErr(err)
    }

    d.SetId(strconv.Itoa(created.ID))
    return updateXxxState(d, &created)  // Go struct → ResourceData
}
```

### State Management

Two-function pattern for bidirectional mapping:
- `buildXxxStruct(d *schema.ResourceData)` - Terraform state → Go struct
- `updateXxxState(d *schema.ResourceData, obj *Type)` - Go struct → Terraform state

### Version-Based Optimistic Locking

All Update operations must include version for conflict detection:
```go
obj.Version = d.Get("version").(int)  // Critical for updates
```

### Error Handling

Handle "not found" gracefully in Read operations. When a resource no longer exists in the API, remove it from Terraform state instead of erroring (allows `terraform apply` to recreate it):
```go
// Pattern 1: Using typed errors
if err == v2.ErrAlertV2NotFound {
    d.SetId("")  // Remove from state
    return nil   // Don't error
}

// Pattern 2: Using HTTP status codes
if statusCode == http.StatusNotFound {
    d.SetId("")
    return nil
}
```

### Timeouts

Standard timeout configuration for CRUD operations:
```go
timeout := 5 * time.Minute

Timeouts: &schema.ResourceTimeout{
    Create: schema.DefaultTimeout(timeout),
    Update: schema.DefaultTimeout(timeout),
    Read:   schema.DefaultTimeout(timeout),
    Delete: schema.DefaultTimeout(timeout),
}
```

### Optional Fields with GetOk()

Use `GetOk()` pattern for optional fields to distinguish between unset and zero values:
```go
if val, ok := d.GetOk("optional_field"); ok {
    obj.OptionalField = val.(string)  // Only set if explicitly configured
}
```

### Deprecating Fields

Mark fields as deprecated with clear migration guidance:
```go
"old_field": {
    Type:       schema.TypeBool,
    Optional:   true,
    Deprecated: "Use 'new_field' in the alert resource instead. This will be removed in a future version.",
}
```

### Read-Only Schema Helpers

Use helper functions from `sysdig/schema.go` for consistent read-only fields:
```go
"computed_field": ReadOnlyIntSchema(),     // Computed int field
"status":         ReadOnlyStringSchema(),  // Computed string field
"enabled":        BoolSchema(),            // Optional bool with default false
"active":         BoolComputedSchema(),    // Computed bool
```

## API Client Architecture

### Interface Hierarchy

```
Common interface (users, teams, notification channels, agent access keys)

SysdigCommon interface
├── Common
├── CustomRoleInterface, CustomRolePermissionInterface
├── GroupMappingInterface, GroupMappingConfigInterface
├── IPFiltersInterface, IPFilteringSettingsInterface
└── TeamServiceAccountInterface

MonitorCommon interface (alerts, alertsV2, dashboards, silence rules, inhibition rules)

SecureCommon interface (posture policies, posture zones, posture controls, posture accept risk, zones)

SysdigMonitor interface
├── SysdigCommon
├── MonitorCommon
└── CloudAccountMonitorInterface

SysdigSecure interface
├── SysdigCommon
├── SecureCommon
├── PolicyInterface, CompositePolicyInterface
├── RuleInterface, MacroInterface, ListInterface
├── CloudauthAccountSecureInterface, CloudauthAccountComponentSecureInterface, CloudauthAccountFeatureSecureInterface
├── OnboardingSecureInterface, OrganizationSecureInterface
└── VulnerabilityPolicyClient, VulnerabilityRuleBundleClient
```

### Client Types

Provider auto-detects client type based on configured credentials:
- `SysdigMonitor` - Direct Sysdig Monitor API
- `SysdigSecure` - Direct Sysdig Secure API
- `IBMMonitor` - IBM Cloud Sysdig Monitor
- `IBMSecure` - IBM Cloud Sysdig Secure

### SysdigClients Interface

Resources access API clients through the `SysdigClients` interface (`sysdig/sysdig_clients.go`):

```go
type SysdigClients interface {
    sysdigMonitorClientV2() (v2.SysdigMonitor, error)
    sysdigSecureClientV2() (v2.SysdigSecure, error)
    ibmMonitorClient() (v2.IBMMonitor, error)
    ibmSecureClient() (v2.IBMSecure, error)
}
```

Each resource implements a client getter function that selects the appropriate client based on configuration:
```go
func getAlertV2Client(c SysdigClients) (v2.AlertV2Interface, error) {
    var client v2.AlertV2Interface
    var err error
    switch c.GetClientType() {
    case IBMMonitor:
        client, err = c.ibmMonitorClient()
    default:
        client, err = c.sysdigMonitorClientV2()
    }
    if err != nil {
        return nil, err
    }
    return client, nil
}
```

### HTTP Client Features

- Retryable HTTP with exponential backoff (max 5 retries)
- Retries on 5xx errors and 409 Conflict
- TLS verification can be disabled via `insecure_tls`
- Custom headers support via `extra_headers`

### Authentication

**Sysdig (Direct):** Bearer token in Authorization header
**IBM Cloud:** IAM token exchange with caching and auto-refresh

## Testing Patterns

### Build Tags

Tests are organized by build tags:
- `unit` - Fast unit tests
- `tf_acc_sysdig_monitor` - Monitor acceptance tests
- `tf_acc_sysdig_secure` - Secure acceptance tests
- `tf_acc_ibm_monitor` - IBM Monitor tests
- `tf_acc_ibm_secure` - IBM Secure tests
- `tf_acc_onprem_monitor`, `tf_acc_onprem_secure` - On-prem tests

### Acceptance Test Structure

```go
//go:build tf_acc_sysdig_secure

func TestAccResourceXxx(t *testing.T) {
    resource.ParallelTest(t, resource.TestCase{
        PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv),
        ProviderFactories: map[string]func() (*schema.Provider, error){
            "sysdig": func() (*schema.Provider, error) {
                return sysdig.Provider(), nil
            },
        },
        Steps: []resource.TestStep{
            {
                Config: testConfig(randomName()),
            },
            {
                ResourceName:      "sysdig_secure_xxx.test",
                ImportState:       true,
                ImportStateVerify: true,
            },
        },
    })
}
```

### Test Helpers (sysdig/common_test.go)

- `preCheckAnyEnv(t, envs...)` - Validate environment variables
- `randomText(len)` - Generate random strings for unique names
- `sysdigOrIBMMonitorPreCheck(t)` - Check for either credential type

### Avoiding Name Collisions

Multiple CI runs share the same Sysdig environment, so hardcoded resource names cause race conditions. Follow these rules:

- **Never hardcode resource names** in test configs — always use unique random names via `rText()` (`acctest.RandStringFromCharSet`) or `randomText()`
- **Prefix with `terraform_test_`** for easy identification and cleanup
- **Never reference built-in Sysdig resources by name** (e.g., the `"container"` macro) — instead, create your own base resource with a unique name and reference that

## Development Workflow (TDD)

Follow the **Red-Green-Refactor** cycle:

1. **Red:** Write acceptance test first (it will fail)
   - Create test file with appropriate build tag
   - Define test cases covering create, update, import
   - Run test to confirm it fails: `go test ./sysdig -v -run TestAccNewResource -tags=tf_acc_sysdig_secure`

2. **Green:** Write minimum code to make test pass
   - Implement resource schema and CRUD operations
   - Add client interface if new API endpoints needed
   - Register in `provider.go`
   - Run test until it passes

3. **Refactor:** Clean up if needed
   - Extract common patterns to shared functions
   - Ensure code follows existing conventions
   - Run `make fmt` and `make lint`

**Prefer small commits** - commit after each complete Red-Green-Refactor cycle. Each commit should represent a single, atomic, working change. This makes code review easier and keeps a clean git history.

## Adding New Resources

1. **Write acceptance test first** in `sysdig/resource_sysdig_<service>_<name>_test.go`
2. **Create resource file:** `sysdig/resource_sysdig_<service>_<name>.go`
3. **Implement schema** using composition pattern with existing base schemas
4. **Implement CRUD:** `CreateContext`, `ReadContext`, `UpdateContext`, `DeleteContext`
5. **Add client interface** in `sysdig/internal/client/v2/` if new API endpoints needed
6. **Register in provider:** Add to `ResourcesMap` in `sysdig/provider.go`
7. **Run tests until green:** Iterate until all tests pass
8. **Add documentation** in `website/docs/r/<name>.md`

## Documentation Conventions

Documentation files in `website/docs/` follow a standard structure:

### Resources (`website/docs/r/<name>.md`)

```markdown
---
subcategory: "Sysdig Secure"  # or "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_<name>"
description: |-
  Short description of the resource.
---

# Resource: sysdig_secure_<name>

Longer description explaining what the resource does.

## Example Usage

\`\`\`terraform
resource "sysdig_secure_<name>" "example" {
  name = "example"
  # ... required and common optional fields
}
\`\`\`

## Argument Reference

### Common Arguments
* `name` - (Required) The name of the resource.

### Type-Specific Arguments
* `specific_field` - (Optional) Description. Default: `value`.

### Nested Schema for `block_name`
* `nested_field` - (Required) Description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:
* `id` - The ID of the resource.
* `version` - The current version of the resource.

## Import

Resource can be imported using the ID:

\`\`\`
$ terraform import sysdig_secure_<name>.example 12345
\`\`\`
```

### Data Sources (`website/docs/d/<name>.md`)

Same structure but without Import section and with read-only field descriptions.

## CI/CD Pipeline

### PR Checks
- Multi-architecture build (darwin, linux, windows, freebsd, openbsd, solaris)
- golangci-lint (30min timeout in CI, 1h timeout in local GNUmakefile)
- Unit tests
- Acceptance tests (Monitor, Secure, IBM suites in parallel)
- Provider documentation validation (tfproviderdocs)
- CodeQL security analysis

### Release Process
- Tag with `v*` pattern triggers release
- GoReleaser builds multi-platform binaries
- GPG signing of checksums
- Automatic GitHub Release creation

## PR Guidelines

- Use [Conventional Commit](https://www.conventionalcommits.org/) format
  - Example: `feat(secure-policy): Add runbook to policy resources`
- Update CODEOWNERS for new resource areas
- Include acceptance tests
- Update documentation in `website/docs/`
