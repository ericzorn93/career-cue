import {
  addProjectConfiguration,
  formatFiles,
  generateFiles,
  Tree,
} from '@nx/devkit';
import * as path from 'path';
import { CreateBackendServiceGeneratorSchema } from './schema';
import { execSync } from 'child_process';

export async function createBackendServiceGenerator(
  tree: Tree,
  options: CreateBackendServiceGeneratorSchema
) {
  const projectRoot = `apps/services/${options.serviceName}`;
  addProjectConfiguration(tree, options.serviceName, {
    root: projectRoot,
    projectType: 'application',
    sourceRoot: projectRoot,
    targets: {
      build: {
        executor: '@nx-go/nx-go:build',
        options: {
          main: '{projectRoot}/cmd/server/main.go',
        },
      },
      serve: {
        executor: '@nx-go/nx-go:serve',
        options: {
          main: '{projectRoot}/cmd/server/main.go',
        },
      },
      test: {
        executor: '@nx-go/nx-go:test',
        options: {
          race: true,
        },
      },
      lint: {
        executor: '@nx-go/nx-go:lint',
      },
      tidy: {
        executor: '@nx-go/nx-go:tidy',
      },
      'docker-build': {
        dependsOn: ['build'],
        command: `docker build -f apps/services/${options.serviceName}/Dockerfile . -t ${options.serviceName}:latest`,
      },
      deploy: {
        executor: 'nx:run-commands',
        options: {
          command: `flyctl deploy -c apps/services/${options.serviceName}/fly.toml -y`,
        },
      },
    },
  });
  generateFiles(tree, path.join(__dirname, 'files'), projectRoot, options);
  await formatFiles(tree);
  await execSync('pnpm go:tidy');
}

export default createBackendServiceGenerator;
