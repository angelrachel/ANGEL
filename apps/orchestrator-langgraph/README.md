# LangGraph/MCP Advisory Boundary

This directory reserves the planner boundary from the blueprint. The planner may produce structured simulation plans, but it must not execute commands, alter scope, create implants, delete evidence, or select live targets.

Inputs and outputs must follow the JSON contracts under `contracts/`. Human review and the Go policy gate remain mandatory before a simulation route is accepted.
