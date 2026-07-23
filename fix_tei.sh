#!/bin/bash
set -e
# api/embed/tei.go
sed -i 's/c.metrics.observeBatchFill(inputs)/c.metrics.ObserveBatchFillRatio(float64(inputs) \/ 32.0)/' api/embed/tei.go
