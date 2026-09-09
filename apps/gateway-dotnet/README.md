# .NET Gateway Boundary

This directory contains the .NET 10.0.401 empty gateway scaffold for the blueprint API boundary. It intentionally exposes no live-target, arbitrary command, implant, persistence, exploit, or destructive endpoints.

## Local commands

```bash
export DOTNET_ROOT="$HOME/.dotnet"
export PATH="$DOTNET_ROOT:$PATH"
dotnet restore AngelGateway.csproj
dotnet build AngelGateway.csproj --no-restore
```

Implement future gateway routes against `contracts/api/openapi.yaml` only. Preserve the Go control-plane policy boundary and expose observe/simulate contracts without bypassing authorization or evidence requirements.
