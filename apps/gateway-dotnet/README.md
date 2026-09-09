# .NET Gateway Boundary

This directory reserves the .NET 10 gateway boundary from the blueprint. The .NET SDK 10.0.401 is installed in the development environment, but no `.csproj` or generated runtime is claimed until the API contract and project ownership are approved.

Implement the gateway against `contracts/api/openapi.yaml` only after the SDK/toolchain is pinned. It must preserve the Go control-plane policy boundary and expose observe/simulate contracts without arbitrary command or live-target endpoints.
