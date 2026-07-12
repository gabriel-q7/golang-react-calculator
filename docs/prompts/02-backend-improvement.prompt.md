The monorepo structure, shared tooling, Docker configuration, `Makefile`, and documentation established during the project bootstrap phase already exist. Review and follow those architectural and organizational decisions instead of creating alternative structures or duplicating configuration. 

Implement the Go REST API within the existing backend application using a clean, testable architecture that separates HTTP handlers, services, and pure calculation logic. Expose JSON endpoints for subtraction, multiplication, division, exponentiation, square root, and percentage operations. Validate all inputs, gracefully handle edge cases such as division by zero and invalid numeric values, and return consistent JSON responses. 

Integrate the backend with the existing Docker and Makefile workflows, add comprehensive unit tests with coverage reporting, and update the existing API documentation in `docs/api.md` and the `README.md` as needed.
