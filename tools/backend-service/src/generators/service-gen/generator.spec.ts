import { createTreeWithEmptyWorkspace } from '@nx/devkit/testing';
import { Tree, readProjectConfiguration } from '@nx/devkit';

import { createBackendServiceGenerator } from './generator';
import { CreateBackendServiceGeneratorSchema } from './schema';

// Mocks
jest.mock('child_process', () => ({
  __esModule: true,
  ...jest.requireActual('child_process'),
  execSync: jest.fn(),
}));

describe('service-gen generator', () => {
  let tree: Tree;
  const options: CreateBackendServiceGeneratorSchema = {
    serviceName: 'test',
    serviceType: 'api',
  };

  beforeEach(() => {
    tree = createTreeWithEmptyWorkspace();
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('should run successfully', async () => {
    await createBackendServiceGenerator(tree, options);
    const config = readProjectConfiguration(tree, 'test');
    expect(config).toBeDefined();
  });
});
