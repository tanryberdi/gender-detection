# Variables
APP_NAME=gender-detector
BUILD_DIR=build
DATA_DIR=data
TRAINING_DIR=training

lint:
	golangci-lint run

fmt:
	gofmt -s -w .

run:
	go run main.go


clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@rm -f $(DATA_DIR)/model.json
	@echo "Clean complete"