type ServiceType = 'app' | 'worker';

export interface CreateBackendServiceGeneratorSchema {
  serviceName: string;
  serviceType: ServiceType;
}
