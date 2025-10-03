# GeoCatalogo + GRO Catalog Integration

This repository has been extended with an in-memory repository backend to enable local testing with the GRO (Global Risk Observatory) unified geospatial catalog.

## What's New

### In-Memory Repository Backend

A new `memory.go` backend has been added that:
- Loads catalog records from a JSON file instead of Elasticsearch
- Supports all standard repository operations (Insert, Get, Query, Delete)
- Implements full text search across title, abstract, and identifiers
- Supports bounding box filtering
- Supports time-based filtering
- Supports collection filtering
- Enables local testing without Docker or Elasticsearch

### Repository Interface

Created a standard `Repository` interface that both Elasticsearch and Memory backends implement, allowing seamless backend switching via configuration.

### GRO Catalog Conversion

Added a converter script that transforms the GRO unified catalog CSV into GeoCatalogo JSON format:
- Location: `/Users/jjohnson/projects/geosure/catalog/scripts/convert_to_geocatalogo.py`
- Input: GRO unified catalog CSV (1,661 entries)
- Output: GeoCatalogo-compatible JSON records

## Running GeoCatalogo with GRO Catalog

### Quick Start

```bash
cd /Users/jjohnson/projects/geocatalogo
./run_gro_catalog.sh
```

The server will start on http://localhost:8000 with 1,599 GRO catalog records loaded.

### Manual Start

```bash
cd /Users/jjohnson/projects/geocatalogo/cmd/geocatalogo

export GEOCATALOGO_REPOSITORY_TYPE=memory
export GEOCATALOGO_REPOSITORY_URL="file:///Users/jjohnson/projects/geosure/catalog/data/geocatalogo_records.json"
export GEOCATALOGO_SERVER_URL=http://localhost:8000

./geocatalogo serve
```

## Example Queries

### Search by keyword
```bash
curl 'http://localhost:8000/?q=wildfire'
```

### Search for specific location
```bash
curl 'http://localhost:8000/?q=california'
```

### Get specific record by ID
```bash
curl 'http://localhost:8000/?recordids=db_h3_l8_union_cities_urban'
```

### Search with pagination
```bash
curl 'http://localhost:8000/?q=database&from=0&size=5'
```

### Multiple record IDs
```bash
curl 'http://localhost:8000/?recordids=db_adm_0,db_adm_1,db_adm_2'
```

## GRO Catalog Structure

The GRO catalog includes:
- **14 database tables** (existing_db) - Production PostgreSQL tables
- **469 local files** (existing_local) - Shapefiles, GeoJSON, CSVs, etc.
- **496 v6 jobs** (potential_v6) - Data collection job profiles
- **682 external sources** (external_*) - APIs, news feeds, academic sources, etc.

Each record includes:
- Standard STAC fields (id, type, bbox, geometry, properties)
- GRO-specific metadata (implementation_status, data_format, geographic_scope)
- Links to S3 storage and local file paths
- Collection categorization by source type

## Files Added/Modified

### New Files
- `repository/repository.go` - Repository interface definition
- `repository/memory.go` - In-memory repository implementation
- `run_gro_catalog.sh` - Helper script to start server with GRO catalog
- `GRO_INTEGRATION.md` - This documentation

### Modified Files
- `geocatalogo.go` - Updated to support backend selection (memory vs elasticsearch)
- `repository/elasticsearch.go` - Added interface implementation comment

### GRO Catalog Files
- `/Users/jjohnson/projects/geosure/catalog/data/geocatalogo_records.json` - Converted GRO catalog (1,599 records)
- `/Users/jjohnson/projects/geosure/catalog/scripts/convert_to_geocatalogo.py` - Conversion script

## Testing

The integration has been tested with:
- ✅ Text search (keyword matching in title, abstract, ID)
- ✅ Record retrieval by ID
- ✅ Pagination (from/size parameters)
- ✅ Multiple record IDs
- ✅ Server startup and record loading (1,599 records loaded successfully)

## Next Steps

To deploy this in production:
1. Review the full deployment plan: `/Users/jjohnson/projects/geosure/catalog/reports/GEOCATALOGO_INTEGRATION_PLAN.md`
2. Consider using AWS OpenSearch instead of Elasticsearch
3. Deploy as Lambda function with API Gateway
4. Set up automated catalog sync from GRO sources
5. Implement STAC API compliance for broader interoperability

## Architecture

```
┌─────────────────────┐
│   GeoCatalogo API   │
│   (Go HTTP Server)  │
└──────────┬──────────┘
           │
           │ Repository Interface
           │
      ┌────┴─────┐
      │          │
┌─────▼────┐ ┌──▼──────────┐
│  Memory  │ │ Elasticsearch│
│ Backend  │ │   Backend    │
└─────┬────┘ └──────────────┘
      │
      │ Loads from
      │
┌─────▼────────────────────────┐
│  geocatalogo_records.json    │
│  (1,599 GRO catalog entries) │
└──────────────────────────────┘
```

## Notes

- The in-memory backend loads all records at startup
- For production use with large catalogs (>10K records), use Elasticsearch/OpenSearch
- The converter can be re-run anytime the GRO catalog is updated
- Collection filtering works differently between backends (memory uses `properties.collection`, elasticsearch uses `properties.product_info.collection`)
