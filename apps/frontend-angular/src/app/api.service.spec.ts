import { TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient } from '@angular/common/http';
import { ApiService } from './api.service';

describe('ApiService', () => {
  let api: ApiService;
  let http: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({ providers: [ApiService, provideHttpClient(), provideHttpClientTesting()] });
    api = TestBed.inject(ApiService);
    http = TestBed.inject(HttpTestingController);
  });

  afterEach(() => http.verify());

  it('requests the layered check catalog with operator headers', () => {
    api.checkCatalog().subscribe(value => expect(value[0].layer).toBe(9));
    const request = http.expectOne('/api/v1/check-catalog');
    expect(request.request.headers.get('Authorization')).toContain('Bearer');
    request.flush([{ id: 'web-security-headers', layer: 9 }]);
  });

  it('requests capabilities and executes a fixture check', () => {
    api.capabilities().subscribe(value => expect(value.default_decision).toBe('deny'));
    http.expectOne('/api/v1/capabilities').flush({ default_decision: 'deny', allowed: [], denied: ['arbitrary-command'] });
    api.execute({ job_id: 'job-1', target: 'fixture://http/security', check_id: 'web-security-headers', fixture: {} }).subscribe(value => expect(value.status).toBe('COMPLETED'));
    const request = http.expectOne('/api/v1/executions');
    expect(request.request.body.target).toBe('fixture://http/security');
    request.flush({ status: 'COMPLETED', check_id: 'web-security-headers', result: { passed: true } });
  });
});
