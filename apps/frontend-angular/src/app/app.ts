import { Component, computed, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterOutlet } from '@angular/router';

interface SimulationTask {
  id: string;
  agentType: string;
  targetRef: string;
  technique: string;
  mode: 'observe' | 'simulate';
  requestedBy: string;
}

@Component({
  selector: 'app-root',
  imports: [FormsModule, RouterOutlet],
  templateUrl: './app.html',
  styleUrl: './app.scss'
})
export class App {
  readonly task = signal<SimulationTask>({
    id: 'sim-task-recon-001',
    agentType: 'recon',
    targetRef: 'fixture://lab/web-app-01',
    technique: 'http-header-observation',
    mode: 'simulate',
    requestedBy: 'operator'
  });
  readonly submitted = signal(false);
  readonly error = signal('');
  readonly valid = computed(() => {
    const value = this.task();
    return Boolean(value.id.trim() && value.agentType.trim() && value.technique.trim() && value.requestedBy.trim() && value.targetRef.startsWith('fixture://'));
  });

  update(field: keyof SimulationTask, value: string): void {
    this.task.update((current) => ({ ...current, [field]: value } as SimulationTask));
    this.submitted.set(false);
    this.error.set('');
  }

  submit(): void {
    if (!this.valid()) {
      this.error.set('Target must be a non-empty fixture:// reference and all fields are required.');
      this.submitted.set(false);
      return;
    }
    this.error.set('');
    this.submitted.set(true);
  }
}
