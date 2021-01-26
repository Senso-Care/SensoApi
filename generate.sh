docker run --rm --user=$(id -u):$(id -g) -v "${PWD}:/local" openapitools/openapi-generator-cli generate \
    -i /local/api/openapi.yaml \
    -g go-server \
    -o /local/go-gen