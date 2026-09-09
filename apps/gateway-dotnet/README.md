# .NET Gateway Boundary

This directory reserves the .NET 10 gateway boundary from the blueprint. The current repository does not contain the .NET SDK, so no `.csproj` or generated runtime is claimed here.

Implement the gateway against `contracts/api/openapi.yaml` only after the SDK/toolchain is pinned. It must preserve the Go control-plane policy boundary and expose observe/simulate contracts without arbitrary command or live-target endpoints.
