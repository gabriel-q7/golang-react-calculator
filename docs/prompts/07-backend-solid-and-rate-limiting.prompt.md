The backend is already functional and integrated with the frontend. Refactor and improve the existing implementation to make it extensible, maintainable, and aligned with SOLID principles while preserving all current functionality and API contracts.

Review the current architecture and identify opportunities to improve the separation of responsibilities. Refactor the code so that adding a new calculator operation requires minimal changes to the existing codebase, following the Open/Closed Principle. Prefer composition over conditionals or large switch statements, and avoid designs that require modifying core business logic whenever a new operation is introduced.

Apply the SOLID principles throughout the backend:

* **Single Responsibility Principle:** Each package, service, and component should have one clear responsibility.
* **Open/Closed Principle:** New calculator operations should be added by implementing a new operation class/type rather than modifying existing business logic.
* **Liskov Substitution Principle:** Ensure all operation implementations conform to a common contract and are interchangeable.
* **Interface Segregation Principle:** Define small, focused interfaces instead of large generic ones.
* **Dependency Inversion Principle:** High-level services should depend on abstractions rather than concrete implementations, using dependency injection where appropriate.

Introduce a common abstraction for calculator operations (for example, an `Operation` interface) and implement each existing operation (subtraction, multiplication, division, exponentiation, square root, and percentage) as an independent implementation of that abstraction. Create a registry or factory responsible for discovering and resolving operations without requiring changes to the service layer when new operations are added.

Improve error handling by defining domain-specific error types and returning consistent JSON error responses. Keep business logic completely independent from the HTTP layer to maximize unit testability.

Since the frontend does not implement authentication or user identification, add backend rate limiting to provide a minimum level of protection against abuse and denial-of-service attempts. Implement IP-based rate limiting as middleware using a well-established Go library (for example, `golang.org/x/time/rate` or another widely adopted solution). Configure reasonable defaults (such as requests per minute with configurable burst capacity), make the limits configurable through environment variables, and return HTTP 429 (Too Many Requests) with a consistent JSON error response when the limit is exceeded. Ensure the middleware is easy to replace or extend in the future if authentication or API keys are introduced.

Update the dependency injection wiring, project documentation, and tests to reflect the new architecture. Add or update unit tests for the operation implementations, registry/factory, services, middleware, and error handling. Verify that all existing functionality continues to work without regression.

Finally, create an ADR in `docs/adr` documenting:

1. The architectural refactoring to support SOLID principles and extensible calculator operations.
2. The decision to implement IP-based rate limiting in lieu of authentication, including the rationale, alternatives considered, limitations, and how the design can evolve to authenticated rate limiting in the future.
