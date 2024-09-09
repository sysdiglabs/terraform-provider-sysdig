#!/usr/bin/env bash

API=""

# Parse input arguments
for arg in "$@"; do
    case $arg in
        --api=*)
            API="${arg#*=}"
            shift
            ;;
        *)
            ;;
    esac
done

# Check if the required parameter (api) is provided
if [ -z "$API" ]; then
    echo "Error: No api parameter provided. Usage: ./generate-from-oas.sh --api=<api_name>"
    exit 1
fi

# Define the base directory for OpenAPI files
BASE_OPENAPI_DIR="../openapi"

# Use the provided parameter to construct the OpenAPI directory path
API="$API"
OPENAPI_DIR="${BASE_OPENAPI_DIR}/${API}"

# Define paths based on the constructed OpenAPI directory
OPENAPI_SPEC="${OPENAPI_DIR}/oas.yaml"
GENERATOR_CONFIG_FILE="${OPENAPI_DIR}/generator_config.yml"
GENERATED_SCHEMA_FILE="${OPENAPI_DIR}/generated-schema/provider_code_spec.json"
OUTPUT_DIR="../sysdig/generated/${API}"

# Check if the OpenAPI spec file exists
if [ ! -f "${OPENAPI_SPEC}" ]; then
    echo "Error: OpenAPI specification file '${OPENAPI_SPEC}' not found in the directory '${OPENAPI_DIR}'."
    exit 1
fi

# Check if the generator config file exists
if [ ! -f "${GENERATOR_CONFIG_FILE}" ]; then
    echo "Error: Configuration file '${GENERATOR_CONFIG_FILE}' not found in the directory '${OPENAPI_DIR}'."
    exit 1
fi

# Check if the necessary tools are installed
if ! command -v tfplugingen-openapi &> /dev/null; then
    echo "terraform-plugin-codegen-openapi (tfplugingen-openapi) could not be found. Please install it using 'go install'."
    exit 1
fi

if ! command -v tfplugingen-framework &> /dev/null; then
    echo "terraform-plugin-codegen-framework (tfplugingen-framework) could not be found. Please install it using 'go install'."
    exit 1
fi

# Create directories if they do not exist
mkdir -p "$(dirname "${GENERATED_SCHEMA_FILE}")"
mkdir -p "${OUTPUT_DIR}"

# Generate OpenAPI Provider Schema
echo "Generating provider schema from OpenAPI spec..."
tfplugingen-openapi generate --config="${GENERATOR_CONFIG_FILE}" --output="${GENERATED_SCHEMA_FILE}" "${OPENAPI_SPEC}"
if [ $? -ne 0 ]; then
    echo "Error generating provider schema. Exiting."
    exit 1
fi
echo "Provider schema generated successfully at ${GENERATED_SCHEMA_FILE}"

# Generate Terraform Provider Code for data-sources
echo "Generating Terraform provider code for data-sources using the generated schema..."
tfplugingen-framework generate data-sources --input="${GENERATED_SCHEMA_FILE}" --output="${OUTPUT_DIR}" --package="${API}"
if [ $? -ne 0 ]; then
    echo "Error generating data-sources code. Exiting."
    exit 1
fi
echo "Terraform provider data-sources code generated successfully at ${OUTPUT_DIR}"

# Generate Terraform Provider Code for resources
echo "Generating Terraform provider code for resources using the generated schema..."
tfplugingen-framework generate resources --input="${GENERATED_SCHEMA_FILE}" --output="${OUTPUT_DIR}" --package="${API}"
if [ $? -ne 0 ]; then
    echo "Error generating resources code. Exiting."
    exit 1
fi
echo "Terraform provider code generated successfully at ${OUTPUT_DIR}"

# Run 'go mod tidy' to update go.mod and go.sum files
echo "Running 'go mod tidy' to clean up dependencies..."
(cd "${OUTPUT_DIR}" && go mod tidy)
if [ $? -ne 0 ]; then
    echo "Error running 'go mod tidy'. Exiting."
    exit 1
fi
echo "'go mod tidy' completed successfully."
