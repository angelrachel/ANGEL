using System.Text.Json.Serialization;

namespace AngelGateway;

public sealed record SimulationTask(
    [property: JsonPropertyName("id")] string Id,
    [property: JsonPropertyName("agent_type")] string AgentType,
    [property: JsonPropertyName("target_ref")] string TargetRef,
    [property: JsonPropertyName("technique")] string Technique,
    [property: JsonPropertyName("mode")] string Mode,
    [property: JsonPropertyName("requested_by")] string RequestedBy,
    [property: JsonPropertyName("approval_id")] string? ApprovalId
);

public static class SimulationTaskValidator
{
    public static string? Validate(SimulationTask? task)
    {
        if (task is null)
        {
            return "task body is required";
        }

        if (string.IsNullOrWhiteSpace(task.Id) ||
            string.IsNullOrWhiteSpace(task.AgentType) ||
            string.IsNullOrWhiteSpace(task.Technique) ||
            string.IsNullOrWhiteSpace(task.RequestedBy))
        {
            return "task identity, agent type, technique, and requester are required";
        }

        if (!task.TargetRef.StartsWith("fixture://", StringComparison.Ordinal) || task.TargetRef.Length <= "fixture://".Length)
        {
            return "target_ref must use a fixture reference";
        }

        if (task.Mode is not ("observe" or "simulate"))
        {
            return "mode must be observe or simulate";
        }

        return null;
    }
}
