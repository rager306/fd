sed -i 's/data, err := os.ReadFile(root)/data, err := os.ReadFile(root) \/\/nolint:gosec/g' api/embed/coalesce_baseline_test.go
sed -i 's/"math\/rand"/"math\/rand\/v2"/g' api/embed/coalesce_baseline_test.go
sed -i 's/rng := rand.New(rand.NewSource(42))/rng := rand.New(rand.NewPCG(42, 42)) \/\/nolint:gosec/g' api/embed/coalesce_baseline_test.go
sed -i 's/(calls int, totalTexts int, durations \[\]time.Duration)/(calls, totalTexts int, durations \[\]time.Duration)/g' api/embed/coalesce_baseline_test.go
