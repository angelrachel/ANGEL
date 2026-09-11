import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface ModuleEntry { id: string; layer: number; name: string; capability: string; safe: boolean; }
export interface CheckDescriptor { id: string; layer: number; }
export interface CapabilityPolicy { default_decision: 'deny'; allowed: string[]; denied: string[]; }
export interface Engagement { id: string; name: string; organization_id: string; authorized: boolean; starts_at: string; ends_at: string; policy_hash: string; stopped: boolean; }
export interface ScopeEntry { id: string; engagement_id: string; kind: string; value: string; ports: number[]; actions: string[]; excluded: boolean; }
export interface AssessmentJob { id: string; engagement_id: string; task_type: string; target: string; policy_hash: string; status: string; created_at: string; signature: string; }
export interface CheckResult { check_id: string; title: string; passed: boolean; severity: string; confidence: number; summary: string; remediation: string; }
export interface ExecutionOutput { job_id: string; check_id: string; status: string; result: CheckResult; evidence: string; collected_at: string; }
export interface Finding { id: string; title: string; severity: string; confidence: number; target: string; status: string; }
export interface Remediation { id: string; finding_id: string; owner: string; plan: string; due_at: string; status: string; }
export interface Retest { id: string; remediation_id: string; passed: boolean; evidence_id: string; notes: string; tested_at: string; }
export interface SnapshotValidation { engagements: number; jobs: number; findings: number; evidence: number; audit_events: number; digest: string; }

@Injectable({ providedIn: 'root' })
export class ApiService {
  private readonly http = inject(HttpClient);
  private readonly headers = new HttpHeaders({ Authorization: 'Bearer local-development-operator' });
  health(): Observable<{ status: string }> { return this.http.get<{ status: string }>('/healthz'); }
  modules(): Observable<ModuleEntry[]> { return this.http.get<ModuleEntry[]>('/api/v1/modules', { headers: this.headers }); }
  createEngagement(payload: { engagement: Engagement; scope: ScopeEntry[]; budget: number }): Observable<Engagement> { return this.http.post<Engagement>('/api/v1/engagements', payload, { headers: this.headers }); }
  createJob(payload: { engagement_id: string; task_type: string; target: string; action: string }): Observable<AssessmentJob> { return this.http.post<AssessmentJob>('/api/v1/assessment-jobs', payload, { headers: this.headers }); }
  checks(): Observable<string[]> { return this.http.get<string[]>('/api/v1/checks', { headers: this.headers }); }
  checkCatalog(): Observable<CheckDescriptor[]> { return this.http.get<CheckDescriptor[]>('/api/v1/check-catalog', { headers: this.headers }); }
  capabilities(): Observable<CapabilityPolicy> { return this.http.get<CapabilityPolicy>('/api/v1/capabilities', { headers: this.headers }); }
  execute(payload: { job_id: string; target: string; check_id: string; fixture: Record<string, unknown> }): Observable<ExecutionOutput> { return this.http.post<ExecutionOutput>('/api/v1/executions', payload, { headers: this.headers }); }
  jobs(): Observable<AssessmentJob[]> { return this.http.get<AssessmentJob[]>('/api/v1/jobs', { headers: this.headers }); }
  findings(): Observable<Finding[]> { return this.http.get<Finding[]>('/api/v1/findings', { headers: this.headers }); }
  evaluate(check_id: string, target: string, fixture: Record<string, unknown>): Observable<CheckResult> { return this.http.post<CheckResult>('/api/v1/checks', { check_id, input: { id: `ui-${Date.now()}`, target, ...fixture } }, { headers: this.headers }); }
  remediations(): Observable<Remediation[]> { return this.http.get<Remediation[]>('/api/v1/remediations', { headers: this.headers }); }
  updateRemediation(id: string, status: string): Observable<Remediation> { return this.http.patch<Remediation>(`/api/v1/remediations?id=${encodeURIComponent(id)}`, { status }, { headers: this.headers }); }
  retest(payload: { remediation_id: string; evidence_id: string; notes: string; passed: boolean }): Observable<Retest> { return this.http.post<Retest>('/api/v1/remediations/retest', payload, { headers: this.headers }); }
  snapshot(): Observable<unknown> { return this.http.get('/api/v1/snapshot', { headers: this.headers }); }
  verifyAudit(): Observable<{ valid: boolean; events: number }> { return this.http.get<{ valid: boolean; events: number }>('/api/v1/audit/verify', { headers: this.headers }); }
}
