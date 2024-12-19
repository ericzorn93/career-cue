import {
  addProjectConfiguration,
  formatFiles,
  generateFiles,
  Tree,
} from '@nx/devkit';
import * as path from 'path';
import { CreateBackendServiceGeneratorSchema } from './schema';
import { spawn } from 'child_process';

export async function createBackendServiceGenerator(
  tree: Tree,
  options: CreateBackendServiceGeneratorSchema
) {
  const projectRoot = `apps/services/${options.serviceName}`;
  addProjectConfiguration(tree, options.serviceName, {
    root: projectRoot,
    projectType: 'application',
    sourceRoot: projectRoot,
    targets: {},
  });
  generateFiles(tree, path.join(__dirname, 'files'), projectRoot, options);
  await formatFiles(tree);
  await spawn('pnpm go:tidy');
}

export default createBackendServiceGenerator;
