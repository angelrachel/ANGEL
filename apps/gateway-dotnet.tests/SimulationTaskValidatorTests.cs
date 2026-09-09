using AngelGateway;

namespace AngelGatewayTests;

public sealed class SimulationTaskValidatorTests
{
    [Fact]
    public void AcceptsFixtureSimulationTask()
    {
        var task = new SimulationTask("task-1", "recon", "fixture://lab/web-app-01", "http-header-observation", "simulate", "operator", null);
        Assert.Null(SimulationTaskValidator.Validate(task));
    }

    [Theory]
    [InlineData("https://example.test", "simulate")]
    [InlineData("fixture://lab/web-app-01", "active")]
    public void RejectsUnsafeTask(string target, string mode)
    {
        var task = new SimulationTask("task-1", "recon", target, "http-header-observation", mode, "operator", null);
        Assert.NotNull(SimulationTaskValidator.Validate(task));
    }
}
