import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface SimulationTaskRequest {
  id: string;
  agent_type: string;
  target_ref: string;
  technique: string;
  mode: 'observe' | 'simulate';
  requested_by: string;
}

export interface SimulationTaskResponse {
  status: string;
  simulation: boolean;
  authorized: boolean;
  message?: string;
}

@Injectable({ providedIn: 'root' })
export class ApiService {
  private readonly http = inject(HttpClient);

  health(): Observable<{ status: string }> {
    return this.http.get<{ status: string }>('/healthz');
  }

  submitSimulation(task: SimulationTaskRequest): Observable<SimulationTaskResponse> {
    return this.http.post<SimulationTaskResponse>('/v1/simulation/tasks', task);
  }
}
