#!/bin/bash
#
# Start geocatalogo with GRO catalog in-memory backend
#

cd "$(dirname "$0")/cmd/geocatalogo"

export GEOCATALOGO_REPOSITORY_TYPE=memory
export GEOCATALOGO_REPOSITORY_URL="file:///Users/jjohnson/projects/geosure/catalog/data/geocatalogo_records.json"
export GEOCATALOGO_SERVER_URL=http://localhost:8000
export GEOCATALOGO_LOGGING_LEVEL=INFO
export GEOCATALOGO_METADATA_IDENTIFICATION_TITLE="GRO Geospatial Data Catalog"
export GEOCATALOGO_METADATA_IDENTIFICATION_ABSTRACT="Global Risk Observatory unified geospatial data catalog"
export GEOCATALOGO_METADATA_PROVIDER_NAME="Geosure"
export GEOCATALOGO_METADATA_PROVIDER_URL="https://geosure.ai"

echo "Starting GeoCatalogo with GRO catalog..."
echo "Server will be available at http://localhost:8000"
echo ""
echo "Example queries:"
echo "  curl 'http://localhost:8000/?q=wildfire'"
echo "  curl 'http://localhost:8000/?q=california'"
echo "  curl 'http://localhost:8000/?recordids=db_h3_l8_union_cities_urban'"
echo ""

./geocatalogo serve
