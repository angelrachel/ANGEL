import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface ModuleEntry { id: string; layer: number; name: string; capability: string; safe: boolean; }
export interface Engagement { id: string; name: string; organization_id: string; authorized: boolean; starts_at: string; ends_at: string; policy_hash: string; stopped: boolean; }
export interface ScopeEntry { id: string; engagement_id: string; kind: string; value: string; ports: number[]; actions: string[]; excluded: boolean; }
export interface AssessmentJob { id: string; engagement_id: string; task_type: string; target: string; policy_hash: string; status: string; created_at: string; signature: string; }
export interface CheckResult { check_id: string; title: string; passed: boolean; severity: string; confidence: number; summary: string; remediation: string; }
export interface Finding { id: string; title: string; severity: string; confidence: number; target: string; status: string; }

@Injectable({ providedIn: 'root' })
export class ApiService {
  private readonly http = inject(HttpClient);
  private readonly headers = new HttpHeaders({ Authorization: 'Bearer local-development-operator' });
  health(): Observable<{ status: string }> { return this.http.get<{ status: string }>('/healthz'); }
  modules(): Observable<ModuleEntry[]> { return this.http.get<ModuleEntry[]>('/api/v1/modules', { headers: this.headers }); }
  createEngagement(payload: { engagement: Engagement; scope: ScopeEntry[]; budget: number }): Observable<Engagement> { return this.http.post<Engagement>('/api/v1/engagements', payload, { headers: this.headers }); }
  createJob(payload: { engagement_id: string; task_type: string; target: string; action: string }): Observable<AssessmentJob> { return this.http.post<AssessmentJob>('/api/v1/assessment-jobs', payload, { headers: this.headers }); }
  checks(): Observable<string[]> { return this.http.get<string[]>('/api/v1/checks', { headers: this.headers }); }
  jobs(): Observable<AssessmentJob[]> { return this.http.get<AssessmentJob[]>('/api/v1/jobs', { headers: this.headers }); }
  findings(): Observable<Finding[]> { return this.http.get<Finding[]>('/api/v1/findings', { headers: this.headers }); }
  evaluate(check_id: string, target: string, fixture: Record<string, unknown>): Observable<CheckResult> { return this.http.post<CheckResult>('/api/v1/checks', { check_id, input: { id: `ui-${Date.now()}`, target, ...fixture } }, { headers: this.headers }); }
}
