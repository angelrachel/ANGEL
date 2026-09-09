# ADR-001: Architecture Boundary for the Current ANGEL Repository

**Status:** Accepted for the current repository baseline  
**Date:** 2026-09-10

## Context

The supplied blueprint describes six layers using Angular, a .NET 10 API gateway, LangGraph/MCP orchestration, and a Go/Rust C2 layer. The repository currently contains a Go module with the control-plane, policy, evidence, reporting, and deployment foundations under `src/c2` and `deploy/`.

A filename-level match is not sufficient to claim that the blueprint runtime exists. The repository must state which architecture is implemented and which requested runtimes remain absent.

## Decision

The current repository is treated as a **Go-only control-plane and simulation foundation**. Its supported engineering scope is:

- fail-closed engagement scope and approval policy;
- control-plane request validation and lifecycle hardening;
- deterministic evidence and report foundations;
- safe simulation routing and passive validation primitives;
- reproducible container and Terraform/Ansible configuration validation.

Angular, .NET 10, and LangGraph/MCP are **not silently substituted** by similarly named Go packages. They remain explicit architecture gaps until the repository owner chooses to add those runtimes or formally revises the blueprint.

High-risk offensive modules remain non-operational and unverified. They are not promoted from skeleton/partial status because they compile.

## Consequences

The status matrix and audit report may classify the repository as foundation/partial without treating that classification as a failure of the Go baseline. Future work must either add the missing runtimes with contract tests or update the blueprint through a reviewed architecture decision. Every new dispatch path must pass the scope and approval gate, and every simulation result must state that no target action was executed.
