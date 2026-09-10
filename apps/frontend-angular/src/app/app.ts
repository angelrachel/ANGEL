import { Component, computed, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterOutlet } from '@angular/router';
import { ApiService, AssessmentJob, Engagement, ModuleEntry, ScopeEntry } from './api.service';

@Component({ selector: 'app-root', imports: [FormsModule, RouterOutlet], templateUrl: './app.html', styleUrl: './app.scss' })
export class App {
  private readonly api = inject(ApiService);
  readonly modules = signal<ModuleEntry[]>([]);
  readonly jobs = signal<AssessmentJob[]>([]);
  readonly error = signal('');
  readonly notice = signal('');
  readonly connected = signal(false);
  readonly engagement = signal<Engagement>({ id: 'eng-local-001', name: 'Authorized Lab Assessment', organization_id: 'local-org', authorized: true, starts_at: new Date(Date.now() - 60000).toISOString(), ends_at: new Date(Date.now() + 3600000).toISOString(), policy_hash: 'local-policy-v1', stopped: false });
  readonly scope = signal<ScopeEntry>({ id: 'scope-local-001', engagement_id: 'eng-local-001', kind: 'fixture', value: 'fixture://lab/web', ports: [], actions: ['surface-map'], excluded: false });
  readonly job = signal({ task_type: 'surface-map', target: 'fixture://lab/web', action: 'surface-map' });
  readonly readiness = computed(() => this.connected() && this.modules().length === 650);

  constructor() { this.refresh(); }
  refresh(): void { this.error.set(''); this.api.health().subscribe({ next: () => { this.connected.set(true); this.api.modules().subscribe({ next: value => this.modules.set(value), error: err => this.error.set(err?.error ?? 'Module catalog unavailable') }); }, error: err => { this.connected.set(false); this.error.set(err?.error ?? 'Assessment API unavailable'); } }); }
  updateEngagement(field: keyof Engagement, value: string | boolean): void { this.engagement.update(current => ({ ...current, [field]: value } as Engagement)); }
  updateScope(field: keyof ScopeEntry, value: string): void { this.scope.update(current => ({ ...current, [field]: value } as ScopeEntry)); }
  updateJob(field: 'task_type' | 'target' | 'action', value: string): void { this.job.update(current => ({ ...current, [field]: value })); }
  register(): void { this.error.set(''); this.api.createEngagement({ engagement: this.engagement(), scope: [this.scope()], budget: 20 }).subscribe({ next: value => { this.engagement.set(value); this.notice.set('Engagement registered and scope policy is active.'); }, error: err => this.error.set(err?.error ?? 'Engagement rejected') }); }
  createJob(): void { this.error.set(''); this.api.createJob({ engagement_id: this.engagement().id, ...this.job() }).subscribe({ next: value => { this.jobs.update(current => [value, ...current]); this.notice.set('Signed assessment job queued.'); }, error: err => this.error.set(err?.error ?? 'Assessment job rejected') }); }
}
