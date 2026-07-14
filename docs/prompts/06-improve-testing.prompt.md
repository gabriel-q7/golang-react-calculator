Refactor the project to separate all test code from the main application source code in both the frontend and backend while preserving the existing architecture, functionality, and coding standards. Review the current project structure and reorganize it to improve maintainability and test discoverability.

For the frontend, move all unit, integration, mocks, fixtures, and test utilities into a dedicated top-level `tests/` directory (or another centralized testing directory consistent with the project's conventions) instead of colocating test files with production components. Organize the tests by feature while keeping imports clean and avoiding circular dependencies.

For the backend, move all test files, fixtures, mocks, helper functions, and test utilities into a dedicated `tests/` directory separate from the production packages. Keep the production code focused solely on the application logic, while ensuring the tests continue to exercise the public interfaces of the application.

Update the project configuration (Vitest/Jest, Go test configuration, coverage configuration, path aliases, Docker, Makefile, CI scripts, and any other affected tooling) so the new directory structure is fully supported. Ensure the following commands continue to work without modification:

* Run frontend tests
* Run backend tests
* Generate frontend coverage reports
* Generate backend coverage reports
* Execute all tests from the repository root

Update the `README.md` to document the new testing structure and execution commands. Finally, create an ADR in `docs/adr` explaining the decision to separate production code from test code, including the motivation, alternatives considered, consequences, and the benefits for maintainability, scalability, readability, and CI/CD workflows.
