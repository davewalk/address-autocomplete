#!/bin/bash

printf "Deleting the index if it already exists..."
curl -X DELETE http://localhost:9200/addresses
printf "\n"

printf "Creating index..."
curl -X PUT http://localhost:9200/addresses
printf "\n"

sleep 1
printf "Closing the index so that we can edit it..."
curl -X POST http://localhost:9200/addresses/_close
printf "\n"

printf "Adding custom 'address' analyzer..."
curl -X PUT -d '{"index": { "analysis": { "tokenizer": { "whitespace": { "type": "whitespace" }}, "analyzer": { "address": { "type": "custom", "tokenizer": "whitespace", "filter": ["trim", "lowercase"]}}}}}' http://localhost:9200/addresses/_settings
printf "\n"

printf "Adding the mapping to the index..."
curl -X PUT -d '{"address": { "properties": { "name": { "type": "string" }, "suggest": { "type": "completion", "index_analyzer": "standard", "search_analyzer": "address", "payloads": true, "preserve_separators": false }}}}' http://localhost:9200/addresses/address/_mapping
printf "\n"

echo "Opening back up for business..."
curl -X POST http://localhost:9200/addresses/_open
