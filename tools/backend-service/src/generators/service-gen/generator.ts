import {
  addProjectConfiguration,
  formatFiles,
  generateFiles,
  ProjectConfiguration,
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

  const defaultConfig: ProjectConfiguration = {
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
          commands: [
            `flyctl deploy -c apps/services/${options.serviceName}/fly.toml -y`,
            `flyctl scale count 1 -r iad -c apps/services/${options.serviceName}/fly.toml -y`,
            `flyctl scale count 0 -r ewr -c apps/services/${options.serviceName}/fly.toml -y`,
            `flyctl scale count 0 -r lax -c apps/services/${options.serviceName}/fly.toml -y`,
            `flyctl scale count 0 -r ord -c apps/services/${options.serviceName}/fly.toml -y`,
          ],
        },
      },
    },
  };

  let finalConfig: ProjectConfiguration = defaultConfig;
  let filesPath: string = path.join(__dirname, 'api-worker-files');

  switch (options.serviceType) {
    case 'graphql':
      Object.assign(finalConfig.targets, {
        generate: {
          executor: 'nx:run-commands',
          options: {
            cwd: '{projectRoot}',
            commands: ['go generate ./...', 'pnpm graphql:gen:dev'],
          },
        },
      });
      filesPath = path.join(__dirname, 'graphql-files');
      break;
    case 'api':
    case 'worker':
    default:
      break;
  }

  addProjectConfiguration(tree, options.serviceName, finalConfig);
  generateFiles(tree, filesPath, projectRoot, options);
  await formatFiles(tree);
  await execSync('pnpm go:tidy');
}

export default createBackendServiceGenerator;
