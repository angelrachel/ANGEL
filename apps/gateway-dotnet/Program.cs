using AngelGateway;

var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

app.MapGet("/healthz", () => Results.Ok(new { status = "ok" }));
app.MapGet("/readyz", () => Results.Ok(new { status = "ready", policy = "deny-by-default", simulation = true }));
app.MapGet("/v1/capabilities", () => Results.Ok(new
{
    default_decision = "deny",
    allowed = new[] { "surface-map", "tls-assessment", "api-contract", "authorization-matrix", "evidence-collection", "dependency-inventory", "detection-validation", "synthetic-canary", "lab-proof", "lab-agent-simulation" },
    denied = new[] { "credential-collection", "credential-extraction", "persistence", "destructive-write", "log-deletion", "covert-channel", "process-injection", "evasion", "data-exfiltration", "arbitrary-command" }
}));

app.MapPost("/v1/simulation/tasks", (SimulationTask? task) =>
{
    var error = SimulationTaskValidator.Validate(task);
    if (error is not null)
    {
        return Results.BadRequest(new { status = "rejected", simulation = true, authorized = false, message = error });
    }

    if (string.IsNullOrWhiteSpace(task!.ApprovalId))
    {
        return Results.StatusCode(StatusCodes.Status403Forbidden);
    }

    return Results.Accepted(value: new { task_id = task.Id, status = "accepted", simulation = true, authorized = true, evidence_refs = Array.Empty<string>() });
});

app.Run();
