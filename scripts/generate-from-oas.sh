#!/usr/bin/env bash

API=""
TAG=""
ASSET_NAME=""
OAS_FILE_NAME="oas.yaml"  # Default output file name for downloaded OpenAPI spec

# Parse input arguments
for arg in "$@"; do
    case $arg in
        --api=*)
            API="${arg#*=}"
            shift
            ;;
        --tag=*)
            TAG="${arg#*=}"
            shift
            ;;
        --assetName=*)
            ASSET_NAME="${arg#*=}"
            shift
            ;;
        *)
            ;;
    esac
done

# Check if the required --api parameter is provided
if [ -z "$API" ]; then
    echo "Error: The --api parameter is required. Usage: ./script.sh --api=<directory_name> [--tag=<tag>] [--assetName=<asset_name>]"
    exit 1
fi

# Define the base directory for OpenAPI files
BASE_OPENAPI_DIR="../openapi"

# Use the provided --api parameter to construct the open api directory path
OPENAPI_DIR="${BASE_OPENAPI_DIR}/${API}"

# Define paths based on the constructed open api directory
OPENAPI_SPEC="${OPENAPI_DIR}/${OAS_FILE_NAME}"
GENERATOR_CONFIG_FILE="${OPENAPI_DIR}/generator_config.yml"
GENERATED_SCHEMA_FILE="${OPENAPI_DIR}/generated-schema/provider_code_spec.json"

# If both tag and asset name are provided, download the OpenAPI spec file
if [ -n "$TAG" ] && [ -n "$ASSET_NAME" ]; then
    echo "Downloading OpenAPI spec file from GitHub release asset..."
    go run ../openapi/download-gh-asset.go -tag="$TAG" -assetName="$ASSET_NAME" -outputFile="${OPENAPI_SPEC}"
    if [ $? -ne 0 ]; then
        echo "Error: Failed to download OpenAPI specification file from GitHub. Exiting."
        exit 1
    fi
    echo "OpenAPI spec file downloaded successfully to ${OPENAPI_SPEC}."
fi


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
OUTPUT_DIR="../sysdig/generated/${API}"
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

echo "Code generation completed successfully."
