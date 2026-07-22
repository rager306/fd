#!/bin/bash
sed -i 's/calls, totalTexts, durations := runCorpusBurst/totalTexts, durations := runCorpusBurst/' embed/coalesce_baseline_test.go
sed -i 's/_, total, durations := runCorpusBurst/total, durations := runCorpusBurst/' embed/coalesce_baseline_test.go
sed -i 's/_, totalTexts, durations := runCorpusBurst/totalTexts, durations := runCorpusBurst/' embed/coalesce_baseline_test.go
