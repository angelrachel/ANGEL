using AngelGateway;

var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

app.MapGet("/healthz", () => Results.Ok(new { status = "ok" }));

app.MapPost("/v1/simulation/tasks", (SimulationTask? task) =>
{
    var error = SimulationTaskValidator.Validate(task);
    if (error is not null)
    {
        return Results.BadRequest(new { status = "rejected", simulation = true, authorized = false, message = error });
    }

    return Results.StatusCode(StatusCodes.Status403Forbidden);
});

app.Run();
