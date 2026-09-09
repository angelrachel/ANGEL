import unittest

from planner import build_plan


class PlannerTests(unittest.TestCase):
    def valid_task(self):
        return {
            "id": "task-1",
            "agent_type": "recon",
            "target_ref": "fixture://lab/web-app-01",
            "technique": "http-header-observation",
            "mode": "simulate",
            "requested_by": "operator",
        }

    def test_builds_unauthorized_simulation_plan(self):
        plan = build_plan(self.valid_task())
        self.assertEqual(plan["status"], "accepted")
        self.assertTrue(plan["simulation"])
        self.assertFalse(plan["authorized"])
        self.assertEqual(plan["evidence_refs"], [])

    def test_rejects_external_target(self):
        task = self.valid_task()
        task["target_ref"] = "https://example.test"
        with self.assertRaisesRegex(ValueError, "fixture://"):
            build_plan(task)

    def test_rejects_active_mode(self):
        task = self.valid_task()
        task["mode"] = "active"
        with self.assertRaisesRegex(ValueError, "observe or simulate"):
            build_plan(task)


if __name__ == "__main__":
    unittest.main()
