import path from "node:path";
import { fileURLToPath } from "node:url";
import {
  createRunContext,
  executeProfiles,
  hasReplayFailures,
  loadFrameworkConfig,
  prepareRunContext,
  printConsoleSummary,
  writeArtifacts,
} from "./framework.mjs";
import { logError, logInfo } from "./utils.mjs";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const BIGTEST_ROOT = path.resolve(__dirname, "..");

async function main() {
  const configPathArg = process.argv[2];
  const { config, configPath } = await loadFrameworkConfig({
    bigtestRoot: BIGTEST_ROOT,
    configPathArg,
  });

  const runContext = createRunContext(BIGTEST_ROOT, config, configPath);
  await prepareRunContext(runContext);

  logInfo(`使用配置文件: ${configPath}`);
  logInfo(`输出目录: ${runContext.outputRoot}`);

  const results = await executeProfiles(runContext);
  const artifactResult = await writeArtifacts(runContext, results);
  printConsoleSummary(runContext, results, artifactResult);

  if (hasReplayFailures(results)) {
    process.exitCode = 1;
  }
}

main().catch((error) => {
  logError(error.stack || error.message);
  process.exit(1);
});
