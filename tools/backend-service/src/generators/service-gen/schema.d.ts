type ServiceType = 'api' | 'graphql' | 'worker';

export interface CreateBackendServiceGeneratorSchema {
  serviceName: string;
  serviceType: ServiceType;
}
