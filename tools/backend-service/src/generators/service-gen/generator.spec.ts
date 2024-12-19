import { createTreeWithEmptyWorkspace } from '@nx/devkit/testing';
import { Tree, readProjectConfiguration } from '@nx/devkit';

import { createBackendServiceGenerator } from './generator';
import { CreateBackendServiceGeneratorSchema } from './schema';

describe('service-gen generator', () => {
  let tree: Tree;
  const options: CreateBackendServiceGeneratorSchema = { serviceName: 'test' };

  beforeEach(() => {
    tree = createTreeWithEmptyWorkspace();
  });

  it('should run successfully', async () => {
    await createBackendServiceGenerator(tree, options);
    const config = readProjectConfiguration(tree, 'test');
    expect(config).toBeDefined();
  });
});
