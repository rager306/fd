sed -i 's/rng := rand.New(rand.NewSource(42))/rng := rand.New(rand.NewSource(42)) \/\/nolint:gosec \/\/ G404: test data/g' api/embed/coalesce_baseline_test.go
sed -i 's/data, err := os.ReadFile(root)/data, err := os.ReadFile(root) \/\/nolint:gosec \/\/ G304: test file inclusion/g' api/embed/coalesce_baseline_test.go
