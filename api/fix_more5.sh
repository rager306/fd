# Fix goconst properly
sed -i '/import (/a\
\
const objectList = "list"' api/handlers/embeddings.go

sed -i '/import (/a\
\
const objectList = "list"' api/handlers/embeddings_integration_test.go

sed -i '/import (/a\
\
const objectList = "list"' api/handlers/queue_handlers.go

sed -i '/import (/a\
\
const labelTier = "tier"' api/observability/metrics.go

sed -i 's/"list"/objectList/g' api/handlers/embeddings.go
sed -i 's/"list"/objectList/g' api/handlers/embeddings_integration_test.go
sed -i 's/"list"/objectList/g' api/handlers/queue_handlers.go
sed -i 's/"tier"/labelTier/g' api/observability/metrics.go
