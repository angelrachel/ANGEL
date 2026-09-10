import { Component, computed, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterOutlet } from '@angular/router';
import { ApiService } from './api.service';

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
  private readonly api = inject(ApiService);
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
    this.api.submitSimulation({
      id: this.task().id,
      agent_type: this.task().agentType,
      target_ref: this.task().targetRef,
      technique: this.task().technique,
      mode: this.task().mode,
      requested_by: this.task().requestedBy
    }).subscribe({
      next: (response) => {
        if (response.authorized === false) {
          this.error.set(response.message ?? 'Gateway policy approval is required.');
          this.submitted.set(false);
          return;
        }
        this.submitted.set(true);
      },
      error: (response) => {
        this.error.set(response?.error?.message ?? 'Gateway rejected the task request.');
        this.submitted.set(false);
      }
    });
  }
}
