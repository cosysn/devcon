## 1. Implement OCI Feature Resolution

- [x] 1.0 Add E2E test to verify features are installed in built image (E2E-035)
- [x] 1.1 Add shorthand name to full OCI reference conversion in pkg/feature/resolver.go
- [x] 1.2 Add OCI feature download function
- [x] 1.3 Add feature definition parsing from OCI labels
- [x] 1.4 Test OCI resolution with a real feature

## 2. Implement Feature Injection

- [x] 2.1 Modify build.go to call feature resolution before build
- [x] 2.2 Generate Dockerfile with feature install steps
- [x] 2.3 Pass feature options as environment variables
- [x] 2.4 Handle feature dependencies (topological sort)

## 3. Integration Testing

- [x] 3.1 Test build with shorthand features (docker-in-docker, git, go)
- [x] 3.2 Verify features are actually installed in the built image
- [x] 3.3 Run full E2E test suite

## 4. Error Handling

- [x] 4.1 Handle OCI registry errors gracefully
- [x] 4.2 Handle missing feature errors
- [x] 4.3 Handle circular dependency errors
