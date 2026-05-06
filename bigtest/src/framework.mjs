import path from "node:path";
import { chromium } from "playwright";
import {
  areShapesCompatible,
  bodyShapeFromValue,
  buildRequestFingerprint,
  deepMerge,
  ensureDir,
  jsonStringify,
  logInfo,
  logWarn,
  normalizeUrlForFingerprint,
  getByPath,
  interpolate,
  listFilesRecursive,
  pathExists,
  relativeFrom,
  readText,
  readTextIfExists,
  sanitizeHeaders,
  sanitizeValue,
  summarizeJsonShape,
  timestampId,
  tryParseJson,
  readJson,
  writeJson,
  writeText,
} from "./utils.mjs";

const DEFAULT_CONFIG = {
  env: {
    frontendBaseUrl: "http://localhost:5173",
    replayBaseUrl: "http://localhost:8080",
    outputDir: "output",
  },
  browser: {
    headless: true,
    slowMoMs: 0,
    channel: "",
    executablePath: "",
  },
  debug: {
    maxScreenshots: 5,
    screenshotOnRouteFailure: true,
    screenshotOnButtonFailure: false,
  },
  capture: {
    urlIncludes: ["/api/v1/"],
    ignoreMethods: ["OPTIONS"],
    ignoreUrlPatterns: ["/api/v1/ai", "/api/v1/evaluations", "/upload"],
    ignoreResourceTypes: ["image", "media", "font", "stylesheet"],
    ignoreQueryKeys: ["timestamp", "_", "nonce"],
    maxEntries: 300,
    responseBodyLimit: 20000,
  },
  generation: {
    deduplicate: true,
    writeRawCaptures: true,
    writeGeneratedCases: true,
    autoGenerateOpenApiCases: true,
    generateNegativeCases: true,
    maxOpenApiCases: 120,
    maxNegativeCases: 40,
  },
  replay: {
    enabled: true,
    allowedMethods: ["GET", "POST", "PUT", "PATCH", "DELETE"],
    ignoreUrlPatterns: ["/api/v1/ai", "/api/v1/evaluations", "/upload"],
    timeoutMs: 15000,
    throttleMs: 50,
    compareResponseShape: true,
    compareContentType: false,
  },
  coverage: {
    openApiCandidates: [
      "../backend/gateway/docs/swagger.json",
      "../backend/docs/swagger.yaml",
    ],
    modules: {
      jobs: "职位",
      talents: "人才",
      resumes: "简历",
      applications: "投递",
      interviews: "面试",
      notices: "公告",
      messages: "消息",
      conversations: "消息",
      "online-status": "消息",
      logs: "日志",
      stats: "统计",
      recommendations: "推荐",
      users: "用户",
      profile: "用户",
      login: "认证",
      register: "认证",
    },
  },
  autoExplore: {
    enabled: true,
    projectRoot: "..",
    routeSource: "",
    routeSourceCandidates: [
      "frontend/src/router/index.ts",
      "frontend/src/router/index.js",
      "src/router/index.ts",
      "src/router/index.js",
      "src/routes.ts",
      "src/routes.js",
      "src/App.tsx",
      "src/App.jsx",
    ],
    pageDirCandidates: [
      "frontend/src/views",
      "frontend/src/pages",
      "src/views",
      "src/pages",
      "pages",
      "app",
    ],
    maxRoutes: 60,
    maxClicksPerPage: 5,
    maxSecondsPerPage: 8,
    maxDomLinksPerPage: 30,
    maxInputsPerPage: 40,
    clickRiskLevels: ["safe"],
    routeParamValues: {
      id: "1",
    },
    excludePathPatterns: [
      "/login",
      "/register",
      "/dev-logs",
      "/ai",
      "/evaluate",
      "/evaluation",
      "/upload",
      "/:pathMatch",
    ],
    safeClickTexts: ["搜索", "查看", "详情", "查看更多", "返回", "重置"],
    mediumRiskClickTexts: ["新增", "添加", "编辑", "保存", "提交", "发布", "确定", "发送", "安排面试", "更新状态", "标记已读", "全部已读"],
    dangerClickTexts: [
      "删除",
      "移除",
      "清空",
      "上传",
      "AI",
      "评估",
      "退出",
      "关闭职位",
      "开放职位",
    ],
    inputRules: {
      "搜索职位名称": "后端",
      "搜索公告标题": "校园招聘",
      "搜索消息内容": "面试",
      "搜索职位": "后端",
      "搜索": "测试",
    },
    defaultInputValue: "测试",
  },
  writeSafety: {
    enabled: true,
    inferFromCaptures: true,
    maxInferredScenarios: 5,
    testDataPrefix: "BIGTEST",
    allowMethods: ["POST", "PUT", "PATCH", "DELETE"],
    allowUnsafeReplayWrites: false,
    allowDeleteOnlyOwnedResources: true,
    autoCleanup: true,
    blockedUrlPatterns: ["/api/v1/ai/", "/api/v1/evaluations", "/ai/evaluate", "/evaluate", "/evaluation", "/upload"],
    scenarios: [],
  },
  profiles: [],
};

const PROFILE_DEFAULTS = {
  name: "default",
  auth: {
    loginPath: "/login",
    username: "",
    password: "",
    usernameSelector: "",
    passwordSelector: "",
    submitSelector: "",
    successUrlContains: "/dashboard",
    postLoginWaitMs: 1500,
  },
  flow: {
    steps: [],
  },
  capture: {},
  replay: {},
  browser: {},
};

export async function loadFrameworkConfig({ bigtestRoot, configPathArg }) {
  const configPath = configPathArg
    ? path.resolve(process.cwd(), configPathArg)
    : path.join(bigtestRoot, "config.local.json");
  const userConfig = await readJson(configPath);
  const normalizedUserConfig = normalizeUserConfig(userConfig);
  const config = deepMerge(DEFAULT_CONFIG, normalizedUserConfig);

  const profiles = config.profiles.map((profile, index) => {
    const merged = deepMerge(PROFILE_DEFAULTS, profile);
    if (!merged.name) {
      merged.name = `profile-${index + 1}`;
    }
    if (!merged.auth.username || !merged.auth.password) {
      throw new Error(`profile ${merged.name} 缺少登录账号或密码`);
    }
    return merged;
  });

  return {
    configPath,
    config: { ...config, profiles },
  };
}

function normalizeUserConfig(userConfig) {
  const normalized = { ...userConfig };

  if (!Array.isArray(normalized.profiles) || normalized.profiles.length === 0) {
    normalized.profiles = [
      {
        name: "auto-default",
        auth: normalized.auth || {},
        flow: normalized.flow || { steps: [] },
      },
    ];
  }

  delete normalized.auth;
  delete normalized.flow;
  return normalized;
}

export function createRunContext(bigtestRoot, config, configPath) {
  const runId = timestampId();
  const outputRoot = path.resolve(bigtestRoot, config.env.outputDir, runId);
  const testDataPrefix = `${config.writeSafety.testDataPrefix}_${runId}`;

  return {
    runId,
    bigtestRoot,
    outputRoot,
    config,
    configPath,
    testDataPrefix,
    registry: {
      resources: {},
    },
  };
}

export async function prepareRunContext(runContext) {
  await ensureDir(runContext.outputRoot);
}

export async function executeProfiles(runContext) {
  const results = [];
  for (const profile of runContext.config.profiles) {
    results.push(await executeProfile(runContext, profile));
  }
  return results;
}

async function executeProfile(runContext, profile) {
  logInfo(`开始执行 profile: ${profile.name}`);
  const rawCapture = await captureProfileTraffic(runContext, profile);
  const capturedCases = buildCasesFromCaptures(
    profile,
    runContext.config,
    rawCapture.captures,
  );
  const writeScenarioResults = runContext.config.writeSafety.enabled
    ? await executeWriteScenarios(runContext, profile, rawCapture.token, capturedCases)
    : [];
  const generatedCases = await buildAllGeneratedCases(
    runContext,
    profile,
    capturedCases,
    writeScenarioResults,
  );
  const replayResults = runContext.config.replay.enabled
    ? await replayGeneratedCases(runContext, profile, generatedCases, rawCapture.token)
    : [];

  return {
    profileName: profile.name,
    rawCapture,
    generatedCases,
    writeScenarioResults,
    replayResults,
  };
}

async function captureProfileTraffic(runContext, profile) {
  const browserConfig = { ...runContext.config.browser, ...profile.browser };
  const launchOptions = {
    headless: browserConfig.headless,
    slowMo: browserConfig.slowMoMs || 0,
  };
  if (browserConfig.channel) {
    launchOptions.channel = browserConfig.channel;
  }
  if (browserConfig.executablePath) {
    launchOptions.executablePath = browserConfig.executablePath;
  }

  const browser = await chromium.launch(launchOptions);
  const context = await browser.newContext();
  const page = await context.newPage();
  const captures = [];
  let captureOrder = 0;
  const captureConfig = { ...runContext.config.capture, ...profile.capture };
  const debugDir = path.join(runContext.outputRoot, profile.name, "debug");
  await ensureDir(debugDir);
  const debugArtifacts = {
    consoleErrors: [],
    networkFailures: [],
    screenshots: [],
  };

  page.on("console", (message) => {
    if (["error", "warning"].includes(message.type())) {
      debugArtifacts.consoleErrors.push({
        type: message.type(),
        text: message.text(),
        url: page.url(),
      });
    }
  });

  page.on("requestfailed", (request) => {
    debugArtifacts.networkFailures.push({
      method: request.method(),
      url: request.url(),
      failure: request.failure()?.errorText || "unknown",
    });
  });

  context.on("response", async (response) => {
    try {
      const entry = await buildCaptureEntry(response, captureConfig);
      if (!entry) {
        return;
      }
      captureOrder += 1;
      entry.id = captureOrder;
      entry.profile = profile.name;
      captures.push(entry);
      logInfo(
        `${profile.name} capture ${entry.method} ${new URL(entry.url).pathname} -> ${entry.response.status}`,
      );
    } catch (error) {
      logWarn(`${profile.name} 记录响应失败: ${error.message}`);
    }
  });

  try {
    const loginUrl = new URL(
      profile.auth.loginPath,
      runContext.config.env.frontendBaseUrl,
    ).toString();
    logInfo(`${profile.name} 打开登录页: ${loginUrl}`);
    await page.goto(loginUrl, { waitUntil: "networkidle" });

    const authSelectors = await resolveAuthSelectors(page, profile.auth);
    await page.locator(authSelectors.usernameSelector).fill(profile.auth.username);
    await page.locator(authSelectors.passwordSelector).fill(profile.auth.password);

    await Promise.all([
      page.waitForURL(
        (url) => url.href.includes(profile.auth.successUrlContains),
        { timeout: runContext.config.replay.timeoutMs },
      ),
      page.locator(authSelectors.submitSelector).click(),
    ]);

    await page.waitForTimeout(profile.auth.postLoginWaitMs || 1000);
    logInfo(`${profile.name} 登录成功: ${page.url()}`);

    for (const step of profile.flow.steps || []) {
      await runFlowStep(page, runContext.config.env.frontendBaseUrl, step);
    }

    const exploration = runContext.config.autoExplore.enabled
      ? await runAutoExplore(page, runContext, profile, debugArtifacts)
      : createEmptyExplorationReport();

    const token = await page.evaluate(() => localStorage.getItem("token"));
    const userRaw = await page.evaluate(() => localStorage.getItem("user"));

    return {
      token,
      user: tryParseJson(userRaw),
      captures,
      exploration,
      debugArtifacts,
    };
  } finally {
    await context.close();
    await browser.close();
  }
}

async function buildCaptureEntry(response, captureConfig) {
  const request = response.request();
  const url = request.url();
  const method = request.method().toUpperCase();
  const resourceType = request.resourceType();

  if (!shouldCapture(url, method, resourceType, captureConfig)) {
    return null;
  }

  const requestHeaders = await request.allHeaders();
  const postData = request.postData();
  const parsedRequestBody = tryParseJson(postData);
  const responseBodyText = await readResponseBody(
    response,
    captureConfig.responseBodyLimit,
  );
  const parsedResponseBody = tryParseJson(responseBodyText);
  const responseHeaders = response.headers();

  return {
    url,
    method,
    resourceType,
    request: {
      headers: requestHeaders,
      headersSanitized: sanitizeHeaders(requestHeaders),
      requiresAuth: Object.keys(requestHeaders).some(
        (key) => key.toLowerCase() === "authorization",
      ),
      postData,
      postDataSanitized: parsedRequestBody
        ? sanitizeValue(parsedRequestBody)
        : postData,
      bodyShape: bodyShapeFromValue(parsedRequestBody ?? postData),
    },
    response: {
      status: response.status(),
      headers: sanitizeHeaders(responseHeaders),
      contentType: responseHeaders["content-type"] || "",
      bodySanitized: parsedResponseBody
        ? sanitizeValue(parsedResponseBody)
        : responseBodyText,
      bodyShape: parsedResponseBody ? summarizeJsonShape(parsedResponseBody) : null,
    },
  };
}

function shouldCapture(url, method, resourceType, captureConfig) {
  if ((captureConfig.ignoreMethods || []).includes(method)) {
    return false;
  }
  if ((captureConfig.ignoreResourceTypes || []).includes(resourceType)) {
    return false;
  }
  if (!(captureConfig.urlIncludes || []).some((item) => url.includes(item))) {
    return false;
  }
  if ((captureConfig.ignoreUrlPatterns || []).some((item) => url.includes(item))) {
    return false;
  }
  return true;
}

async function readResponseBody(response, limit) {
  const contentType = response.headers()["content-type"] || "";
  if (
    contentType.includes("application/json") ||
    contentType.includes("text/") ||
    contentType.includes("application/problem+json")
  ) {
    const text = await response.text();
    return text.length > limit ? `${text.slice(0, limit)}\n...[truncated]` : text;
  }
  return `[omitted body: ${contentType || "unknown"}]`;
}

async function runFlowStep(page, frontendBaseUrl, step) {
  switch (step.type) {
    case "goto": {
      const target = step.url || new URL(step.path || "/", frontendBaseUrl).toString();
      logInfo(`step goto ${target}`);
      await page.goto(target, { waitUntil: step.waitUntil || "networkidle" });
      if (step.waitForMs) {
        await page.waitForTimeout(step.waitForMs);
      }
      return;
    }
    case "click":
      logInfo(`step click ${step.selector}`);
      await page.locator(step.selector).click();
      if (step.waitForMs) {
        await page.waitForTimeout(step.waitForMs);
      }
      return;
    case "fill":
      logInfo(`step fill ${step.selector}`);
      await page.locator(step.selector).fill(step.value || "");
      return;
    case "press":
      logInfo(`step press ${step.selector} -> ${step.key}`);
      await page.locator(step.selector).press(step.key);
      return;
    case "select":
      logInfo(`step select ${step.selector}`);
      await page.locator(step.selector).selectOption(step.value);
      return;
    case "check":
      logInfo(`step check ${step.selector}`);
      await page.locator(step.selector).check();
      return;
    case "uncheck":
      logInfo(`step uncheck ${step.selector}`);
      await page.locator(step.selector).uncheck();
      return;
    case "wait":
      logInfo(`step wait ${step.ms}ms`);
      await page.waitForTimeout(step.ms || 1000);
      return;
    default:
      throw new Error(`不支持的 flow step type: ${step.type}`);
  }
}

async function resolveAuthSelectors(page, auth) {
  return {
    usernameSelector:
      auth.usernameSelector || (await firstExistingSelector(page, [
        "input[autocomplete='username']",
        "input[name='username']",
        "input[name='email']",
        "input[type='email']",
        "input[placeholder*='用户名']",
        "input[placeholder*='邮箱']",
        "input[placeholder*='账号']",
        "input[type='text']",
      ])),
    passwordSelector:
      auth.passwordSelector || (await firstExistingSelector(page, [
        "input[autocomplete='current-password']",
        "input[name='password']",
        "input[type='password']",
        "input[placeholder*='密码']",
      ])),
    submitSelector:
      auth.submitSelector || (await firstExistingSelector(page, [
        "button[type='submit']",
        "button.login-btn",
        "button:has-text('登录')",
        "button:has-text('Login')",
        "input[type='submit']",
      ])),
  };
}

async function firstExistingSelector(page, selectors) {
  for (const selector of selectors) {
    try {
      if ((await page.locator(selector).count()) > 0) {
        return selector;
      }
    } catch {
      // 忽略当前框架不支持的选择器，继续尝试下一个。
    }
  }
  throw new Error(`无法自动识别登录选择器，候选: ${selectors.join(", ")}`);
}

async function runAutoExplore(page, runContext, profile, debugArtifacts) {
  const exploreConfig = {
    ...runContext.config.autoExplore,
    ...(profile.autoExplore || {}),
  };
  const report = createEmptyExplorationReport();
  const initialRoutes = await discoverRoutes(runContext.bigtestRoot, exploreConfig, report);
  const pendingRoutes = [...initialRoutes];
  const visitedRoutes = new Set();
  logInfo(`${profile.name} 自动探索路由数量: ${initialRoutes.length}`);

  page.on("dialog", async (dialog) => {
    logWarn(`自动探索捕获确认弹窗，已取消: ${dialog.message()}`);
    report.dialogsDismissed += 1;
    await dialog.dismiss();
  });

  while (pendingRoutes.length > 0 && visitedRoutes.size < (exploreConfig.maxRoutes || 60)) {
    const routePath = pendingRoutes.shift();
    if (visitedRoutes.has(routePath) || shouldExcludeRoute(routePath, exploreConfig)) {
      continue;
    }
    visitedRoutes.add(routePath);
    const startedAt = Date.now();
    const target = new URL(routePath, runContext.config.env.frontendBaseUrl).toString();
    try {
      logInfo(`${profile.name} autoExplore goto ${target}`);
      await page.goto(target, { waitUntil: "networkidle", timeout: runContext.config.replay.timeoutMs });
      await page.waitForTimeout(600);
      report.pagesVisited += 1;
      report.visitedRoutes.push(routePath);
      for (const domRoute of await discoverRoutesFromDom(page, exploreConfig, report)) {
        if (!visitedRoutes.has(domRoute) && !pendingRoutes.includes(domRoute)) {
          pendingRoutes.push(domRoute);
        }
      }
      await fillSafeInputs(page, exploreConfig, report, runContext);
      await clickSafeButtons(page, exploreConfig, startedAt, report, runContext, profile, debugArtifacts);
    } catch (error) {
      logWarn(`${profile.name} 自动探索 ${routePath} 失败: ${error.message}`);
      const screenshot = shouldSaveScreenshot(runContext, debugArtifacts, "route")
        ? await saveDebugScreenshot(
            page,
            path.join(runContext.outputRoot, profile.name, "debug"),
            `route-failed-${report.failedRoutes.length + 1}`,
          )
        : "";
      if (screenshot) {
        debugArtifacts?.screenshots.push({ type: "route-failed", path: routePath, file: screenshot });
      }
      report.failedRoutes.push({ path: routePath, error: error.message, screenshot });
    }
  }

  report.routesDiscovered = new Set([
    ...report.routesFromFiles,
    ...report.routesFromPageDirs,
    ...report.routesFromDom,
  ]).size;
  return report;
}

function createEmptyExplorationReport() {
  return {
    routesDiscovered: 0,
    routesFromFiles: [],
    routesFromPageDirs: [],
    routesFromDom: [],
    visitedRoutes: [],
    pagesVisited: 0,
    inputsFound: 0,
    inputsFilled: 0,
    buttonsFound: 0,
    buttonsClicked: 0,
    skippedButtons: [],
    buttonFailures: [],
    failedRoutes: [],
    dialogsDismissed: 0,
  };
}

async function discoverRoutes(bigtestRoot, exploreConfig, report) {
  const sourceRoutes = await discoverRoutesFromSourceFiles(bigtestRoot, exploreConfig, report);
  const pageDirRoutes =
    sourceRoutes.length === 0
      ? await discoverRoutesFromPageDirs(bigtestRoot, exploreConfig, report)
      : [];
  const routes = [...sourceRoutes, ...pageDirRoutes];
  const excluded = exploreConfig.excludePathPatterns || [];
  const included = exploreConfig.includePathPatterns || [];

  return routes
    .filter((routePath) => routePath && routePath !== "/")
    .filter((routePath) => !excluded.some((pattern) => routePath.includes(pattern)))
    .filter((routePath) =>
      included.length === 0 ? true : included.some((pattern) => routePath.includes(pattern)),
    )
    .filter((routePath, index, all) => all.indexOf(routePath) === index)
    .slice(0, exploreConfig.maxRoutes || 60);
}

async function discoverRoutesFromSourceFiles(bigtestRoot, exploreConfig, report) {
  const projectRoot = path.resolve(bigtestRoot, exploreConfig.projectRoot || "..");
  const candidatePaths = [];
  if (exploreConfig.routeSource) {
    candidatePaths.push(path.resolve(bigtestRoot, exploreConfig.routeSource));
  }
  for (const candidate of exploreConfig.routeSourceCandidates || []) {
    candidatePaths.push(path.resolve(projectRoot, candidate));
  }

  const routes = [];
  for (const sourcePath of [...new Set(candidatePaths)]) {
    const source = await readTextIfExists(sourcePath);
    if (!source) {
      continue;
    }
    const parsed = parseRouterPaths(source, exploreConfig.routeParamValues || {});
    routes.push(...parsed);
    report.routesFromFiles.push(...parsed);
  }

  return [...new Set(routes)];
}

async function discoverRoutesFromPageDirs(bigtestRoot, exploreConfig, report) {
  const projectRoot = path.resolve(bigtestRoot, exploreConfig.projectRoot || "..");
  const routes = [];
  for (const candidate of exploreConfig.pageDirCandidates || []) {
    const dirPath = path.resolve(projectRoot, candidate);
    if (!(await pathExists(dirPath))) {
      continue;
    }
    const files = await listFilesRecursive(dirPath, {
      maxDepth: 5,
      includeExtensions: [".vue", ".tsx", ".jsx", ".ts", ".js"],
    });
    for (const filePath of files) {
      const routePath = routeFromPageFile(filePath, dirPath);
      if (routePath) {
        routes.push(routePath);
        report.routesFromPageDirs.push(routePath);
      }
    }
  }
  return [...new Set(routes)];
}

function routeFromPageFile(filePath, rootDir) {
  const relativePath = path.relative(rootDir, filePath).replace(/\\/g, "/");
  const withoutExtension = relativePath.replace(/\.(vue|tsx|jsx|ts|js)$/i, "");
  const segments = withoutExtension
    .split("/")
    .filter((item) => !["index", "Index"].includes(item))
    .map((segment) =>
      segment
        .replace(/Page$/i, "")
        .replace(/View$/i, "")
        .replace(/List$/i, "")
        .replace(/Detail$/i, "/1")
        .replace(/\[([^\]]+)\]/g, (_, key) => (key.toLowerCase() === "id" ? "1" : "test")),
    )
    .filter(Boolean);

  if (segments.length === 0) {
    return "/";
  }
  return `/${segments.join("/")}`.replace(/\/+/g, "/").toLowerCase();
}

async function discoverRoutesFromDom(page, exploreConfig, report) {
  const hrefs = await page
    .locator("a[href], [role='link'][href]")
    .evaluateAll((items) =>
      items
        .map((item) => item.getAttribute("href"))
        .filter(Boolean),
    )
    .catch(() => []);
  const routes = hrefs
    .map((href) => routeFromHref(href, page.url()))
    .filter((routePath) => routePath && !shouldExcludeRoute(routePath, exploreConfig))
    .slice(0, exploreConfig.maxDomLinksPerPage || 30);

  report.routesFromDom.push(...routes);
  return [...new Set(routes)];
}

function routeFromHref(href, currentUrl) {
  try {
    const url = new URL(href, currentUrl);
    if (url.origin !== new URL(currentUrl).origin) {
      return "";
    }
    return `${url.pathname}${url.search}` || "/";
  } catch {
    return "";
  }
}

function shouldExcludeRoute(routePath, exploreConfig) {
  return (exploreConfig.excludePathPatterns || []).some((pattern) =>
    routePath.includes(pattern),
  );
}

function parseRouterPaths(source, routeParamValues) {
  const paths = [];
  const lines = source.split(/\r?\n/);
  let context = "";

  for (const line of lines) {
    if (line.includes("path: \"/portal\"")) {
      context = "/portal";
    } else if (line.includes("path: \"/\"")) {
      context = "";
    }

    const matches = [
      ...line.matchAll(/path:\s*["']([^"']*)["']/g),
      ...line.matchAll(/path=["']([^"']*)["']/g),
    ];
    if (matches.length === 0) {
      continue;
    }

    for (const match of matches) {
      const rawPath = match[1];
      if (rawPath === "" || rawPath === "/" || rawPath === "*") {
        continue;
      }
      if (rawPath.includes("pathMatch")) {
        continue;
      }

      const fullPath = rawPath.startsWith("/")
        ? rawPath
        : `${context}/${rawPath}`.replace(/\/+/g, "/");
      paths.push(materializeRouteParams(fullPath, routeParamValues));
    }
  }

  return [...new Set(paths)];
}

function materializeRouteParams(routePath, routeParamValues) {
  return routePath.replace(/:([A-Za-z0-9_]+)(\([^)]*\))?(\*)?/g, (_, key) => {
    return routeParamValues[key] || "1";
  });
}

async function fillSafeInputs(page, exploreConfig, report, runContext) {
  const inputs = await page
    .locator("input:not([type='hidden']):not([type='file']), textarea")
    .all();
  report.inputsFound += inputs.length;
  for (const input of inputs.slice(0, exploreConfig.maxInputsPerPage || 40)) {
    try {
      if (!(await input.isVisible()) || !(await input.isEditable())) {
        continue;
      }
      const meta = await readInputMeta(input);
      const value = chooseInputValue(meta, exploreConfig, runContext);
      if (!value) {
        continue;
      }
      await input.fill(value, { timeout: 1000 });
      report.inputsFilled += 1;
    } catch {
      // 忽略不可编辑或组件代理输入框。
    }
  }

  const selects = await page.locator("select").all();
  report.inputsFound += selects.length;
  for (const select of selects.slice(0, exploreConfig.maxInputsPerPage || 40)) {
    try {
      if (!(await select.isVisible()) || !(await select.isEditable())) {
        continue;
      }
      const value = await select
        .locator("option")
        .nth(1)
        .getAttribute("value", { timeout: 500 });
      if (value) {
        await select.selectOption(value, { timeout: 1000 });
        report.inputsFilled += 1;
      }
    } catch {
      // 忽略不可选的原生 select。
    }
  }
}

async function readInputMeta(input) {
  const placeholder = (await input.getAttribute("placeholder")) || "";
  const name = (await input.getAttribute("name")) || "";
  const type = (await input.getAttribute("type")) || "";
  const ariaLabel = (await input.getAttribute("aria-label")) || "";
  const text = [placeholder, name, type, ariaLabel].filter(Boolean).join(" ");
  return { placeholder, name, type, ariaLabel, text };
}

function chooseInputValue(meta, exploreConfig, runContext) {
  const normalized = meta.text.toLowerCase();
  if (["password", "file", "hidden"].includes(meta.type)) {
    return "";
  }
  for (const [keyword, value] of Object.entries(exploreConfig.inputRules || {})) {
    if (meta.text.includes(keyword)) {
      return value;
    }
  }
  if (/搜索|查询|关键字|search|keyword/.test(meta.text)) {
    return exploreConfig.defaultInputValue || "测试";
  }
  if (/email|邮箱|mail/.test(normalized)) {
    return `bigtest_${runContext.runId.replace(/[^a-zA-Z0-9]/g, "_")}@example.com`;
  }
  if (/phone|mobile|手机号|电话/.test(normalized)) {
    return "13800000000";
  }
  if (/date|time|日期|时间/.test(normalized)) {
    return new Date(Date.now() + 86400000).toISOString().slice(0, 10);
  }
  if (/title|name|名称|标题/.test(normalized)) {
    return `${runContext.testDataPrefix}_名称`;
  }
  if (/content|description|remark|备注|内容|描述/.test(normalized)) {
    return `${runContext.testDataPrefix}_自动化测试内容`;
  }
  return "";
}

async function clickSafeButtons(page, exploreConfig, startedAt, report, runContext, profile, debugArtifacts) {
  let clicked = 0;
  const maxClicks = exploreConfig.maxClicksPerPage || 5;
  const maxMs = (exploreConfig.maxSecondsPerPage || 8) * 1000;
  const clickableRiskLevels = new Set(exploreConfig.clickRiskLevels || ["safe"]);
  const buttons = await page.locator("button, a, .el-button, [role='button']").all();
  report.buttonsFound += buttons.length;

  for (const button of buttons) {
    if (clicked >= maxClicks || Date.now() - startedAt > maxMs) {
      return;
    }
    try {
      if (!(await button.isVisible())) {
        continue;
      }
      if (!(await button.isEnabled())) {
        report.skippedButtons.push({
          route: page.url(),
          text: (await readClickableText(button)) || "<disabled>",
          risk: "disabled",
        });
        continue;
      }
      const text = await readClickableText(button);
      const risk = classifyClickText(text, exploreConfig);
      if (!clickableRiskLevels.has(risk)) {
        report.skippedButtons.push({
          route: page.url(),
          text: text || "<empty>",
          risk,
        });
        continue;
      }
      await button.click({ timeout: 1000 });
      clicked += 1;
      report.buttonsClicked += 1;
      await page.waitForTimeout(700);
      await fillSafeInputs(page, exploreConfig, report, runContext);
    } catch (error) {
      const screenshot = shouldSaveScreenshot(runContext, debugArtifacts, "button")
        ? await saveDebugScreenshot(
            page,
            path.join(runContext.outputRoot, profile.name, "debug"),
            `button-failed-${report.buttonFailures.length + 1}`,
          )
        : "";
      if (screenshot) {
        debugArtifacts?.screenshots.push({ type: "button-failed", route: page.url(), file: screenshot });
      }
      report.buttonFailures.push({
        route: page.url(),
        error: error.message,
        screenshot,
      });
    }
  }
}

function isSafeClickText(text, exploreConfig) {
  return classifyClickText(text, exploreConfig) === "safe";
}

async function readClickableText(locator) {
  const values = [];
  try {
    const domMeta = await locator.evaluate((element) => {
      const attrs = ["aria-label", "title", "data-testid", "data-test", "class"]
        .map((name) => element.getAttribute(name))
        .filter(Boolean);
      const datasetValues = Object.values(element.dataset || {}).filter(Boolean);
      const describedBy = element.getAttribute("aria-describedby");
      const tooltipText = describedBy
        ? document.getElementById(describedBy)?.textContent || ""
        : "";
      const parentText = element.parentElement?.textContent || "";
      return [...attrs, ...datasetValues, tooltipText, parentText].join(" ");
    });
    if (domMeta) {
      values.push(domMeta);
    }
  } catch {
    // 忽略 DOM 上下文读取失败。
  }
  for (const reader of [
    () => locator.innerText({ timeout: 500 }),
    () => locator.getAttribute("aria-label", { timeout: 500 }),
    () => locator.getAttribute("title", { timeout: 500 }),
    () => locator.getAttribute("data-testid", { timeout: 500 }),
    () => locator.getAttribute("data-test", { timeout: 500 }),
    () => locator.getAttribute("class", { timeout: 500 }),
  ]) {
    try {
      const value = await reader();
      if (value) {
        values.push(value);
      }
    } catch {
      // 忽略不可读取属性。
    }
  }

  const combined = values.join(" ").trim();
  if (combined) {
    return normalizeClickableText(combined);
  }

  try {
    const parentText = await locator.locator("xpath=..").innerText({ timeout: 500 });
    return normalizeClickableText(parentText);
  } catch {
    return "";
  }
}

function normalizeClickableText(text) {
  const normalized = String(text || "").replace(/\s+/g, " ").trim();
  const lower = normalized.toLowerCase();
  const iconHints = [
    ["edit", "编辑"],
    ["delete", "删除"],
    ["remove", "删除"],
    ["view", "查看"],
    ["detail", "详情"],
    ["search", "搜索"],
    ["plus", "新增"],
    ["add", "新增"],
    ["close", "关闭"],
    ["upload", "上传"],
  ];
  for (const [keyword, label] of iconHints) {
    if (lower.includes(keyword)) {
      return label;
    }
  }
  return normalized;
}

function classifyClickText(text, exploreConfig) {
  if (!text) {
    return "empty";
  }
  if ((exploreConfig.dangerClickTexts || []).some((danger) => text.includes(danger))) {
    return "danger";
  }
  if ((exploreConfig.mediumRiskClickTexts || []).some((medium) => text.includes(medium))) {
    return "medium";
  }
  return (exploreConfig.safeClickTexts || []).some((safe) => text.includes(safe))
    ? "safe"
    : "unknown";
}

function buildCasesFromCaptures(profile, config, captures) {
  if (captures.length === 0) {
    return [];
  }

  const captureConfig = { ...config.capture, ...profile.capture };
  const groups = new Map();

  for (const capture of captures.slice(0, captureConfig.maxEntries)) {
    const fingerprint = buildRequestFingerprint(capture, captureConfig.ignoreQueryKeys);
    if (!config.generation.deduplicate) {
      groups.set(`${fingerprint}::${capture.id}`, [capture]);
      continue;
    }
    if (!groups.has(fingerprint)) {
      groups.set(fingerprint, []);
    }
    groups.get(fingerprint).push(capture);
  }

  return [...groups.entries()].map(([fingerprint, items], index) => {
    const first = items[0];
    const normalizedPath = normalizeUrlForFingerprint(
      first.url,
      captureConfig.ignoreQueryKeys,
    );

    return {
      id: `${profile.name}-${index + 1}`,
      profile: profile.name,
      fingerprint,
      occurrences: items.length,
      sourceRequestIds: items.map((item) => item.id),
      request: {
        method: first.method,
        originalUrl: first.url,
        normalizedPath,
        resourceType: first.resourceType,
        requiresAuth: first.request.requiresAuth,
        headersSanitized: first.request.headersSanitized,
        bodySanitized: first.request.postDataSanitized,
        bodyShape: first.request.bodyShape,
      },
      expected: {
        status: first.response.status,
        responseShape: first.response.bodyShape,
        contentType: first.response.contentType,
      },
    };
  });
}

async function buildAllGeneratedCases(runContext, profile, capturedCases, writeScenarioResults) {
  const cases = [...capturedCases];
  const apiSpec = runContext.config.generation.autoGenerateOpenApiCases
    ? await loadOpenApiSpec(runContext)
    : { operations: [] };
  const resourcePool = buildResourcePool(cases, writeScenarioResults);

  if (runContext.config.generation.autoGenerateOpenApiCases) {
    cases.push(
      ...buildOpenApiCases(
        runContext,
        profile,
        apiSpec.operations,
        cases,
        resourcePool,
      ),
    );
  }

  if (runContext.config.generation.generateNegativeCases) {
    cases.push(
      ...buildNegativeCases(
        runContext,
        profile,
        apiSpec.operations,
        cases,
        resourcePool,
      ),
    );
  }

  return cases;
}

function buildOpenApiCases(runContext, profile, operations, existingCases, resourcePool) {
  const maxCases = runContext.config.generation.maxOpenApiCases || 120;
  const output = [];
  for (const operation of operations) {
    if (output.length >= maxCases) {
      break;
    }
    if (isOperationAlreadyCovered(operation, [...existingCases, ...output])) {
      continue;
    }
    const request = buildOperationRequest(runContext, operation, resourcePool, {
      negative: false,
      profile,
    });
    if (!request) {
      continue;
    }
    output.push(createGeneratedCaseFromOperation(profile, operation, request, {
      index: existingCases.length + output.length + 1,
      source: "openapi",
      expectedStatuses: operation.expectedStatuses,
    }));
  }
  return output;
}

function buildNegativeCases(runContext, profile, operations, existingCases, resourcePool) {
  const maxCases = runContext.config.generation.maxNegativeCases || 40;
  const output = [];
  for (const operation of operations) {
    if (output.length >= maxCases) {
      break;
    }
    if (!["GET", "POST", "PUT", "PATCH", "DELETE"].includes(operation.method)) {
      continue;
    }
    if (operation.requiresAuth) {
      const unauthorizedRequest = buildOperationRequest(runContext, operation, resourcePool, {
        negative: false,
        profile,
      });
      if (unauthorizedRequest) {
        unauthorizedRequest.forceNoAuth = true;
        output.push(createGeneratedCaseFromOperation(profile, operation, unauthorizedRequest, {
          index: existingCases.length + output.length + 1,
          source: "negative",
          expectedStatuses: [401, 403],
          negative: true,
          subtype: "unauthorized",
        }));
      }
      if (output.length >= maxCases) {
        break;
      }
    }
    if (operation.method === "GET") {
      continue;
    }
    const request = buildOperationRequest(runContext, operation, resourcePool, {
      negative: true,
      profile,
    });
    if (!request) {
      continue;
    }
    output.push(createGeneratedCaseFromOperation(profile, operation, request, {
      index: existingCases.length + output.length + 1,
      source: "negative",
      expectedStatuses: chooseNegativeExpectedStatuses(operation),
      negative: true,
      subtype: "invalid-input",
    }));
  }
  return output;
}

function buildOperationRequest(runContext, operation, resourcePool, options) {
  const pathValue = materializeApiPath(operation.path, operation, resourcePool, options);
  if (!pathValue) {
    return null;
  }
  const query = buildQueryString(operation.parameters || [], options.negative);
  const normalizedPath = `${pathValue}${query}`;
  const resourceName = parseApiResource(operation.path)?.name || "";
  const body = buildBodyForOperation(runContext, operation, options);
  return {
    method: operation.method,
    pathTemplate: operation.path,
    queryString: query,
    resourceName,
    normalizedPath,
    originalUrl: new URL(normalizedPath, runContext.config.env.replayBaseUrl).toString(),
    resourceType: "fetch",
    requiresAuth: operation.requiresAuth,
    headersSanitized: {
      accept: "application/json",
      "content-type": "application/json",
    },
    bodySanitized: body,
    bodyShape: bodyShapeFromValue(body),
  };
}

function createGeneratedCaseFromOperation(profile, operation, request, options) {
  return {
    id: `${profile.name}-${options.source}-${options.index}`,
    profile: profile.name,
    fingerprint: `${options.source}:${request.method}:${request.normalizedPath}`,
    occurrences: 1,
    sourceRequestIds: [],
    source: options.source,
    negative: Boolean(options.negative),
    negativeSubtype: options.subtype || "",
    operationSummary: operation.summary || "",
    request,
    expected: {
      status: options.expectedStatuses,
      responseShape: null,
      contentType: "application/json",
    },
  };
}

function isOperationAlreadyCovered(operation, cases) {
  return cases.some((testCase) =>
    String(testCase.request?.method || "").toUpperCase() === operation.method &&
    pathMatchesTemplate(operation.path, testCase.request?.normalizedPath || ""),
  );
}

function buildResourcePool(cases, writeScenarioResults) {
  const pool = {};
  const add = (resourceName, id) => {
    if (!resourceName || !id) {
      return;
    }
    if (!pool[resourceName]) {
      pool[resourceName] = [];
    }
    if (!pool[resourceName].includes(String(id))) {
      pool[resourceName].push(String(id));
    }
  };

  for (const testCase of cases || []) {
    const resource = parseApiResource(testCase.request?.normalizedPath || "");
    const id = extractLastPathId(testCase.request?.normalizedPath || "");
    add(resource?.name, id);
  }

  for (const scenario of writeScenarioResults || []) {
    for (const step of scenario.steps || []) {
      const resource = parseApiResource(step.path || "");
      const id = getByPath(step.body, "data.id") || extractLastPathId(step.path || "");
      add(resource?.name || scenario.resourceType, id);
    }
  }

  return pool;
}

function materializeApiPath(pathTemplate, operation, resourcePool, options) {
  return String(pathTemplate || "").replace(/\{([^}]+)\}|:([A-Za-z0-9_]+)/g, (_, braceKey, colonKey) => {
    const key = braceKey || colonKey;
    if (options.negative && /id$/i.test(key)) {
      return "99999999";
    }
    return choosePathParamValue(key, operation, resourcePool);
  });
}

function choosePathParamValue(key, operation, resourcePool) {
  const lower = String(key || "").toLowerCase();
  const resource = parseApiResource(operation.path);
  if (lower === "id") {
    return resourcePool[resource?.name]?.[0] || "1";
  }
  if (lower.includes("job") || lower.includes("position")) {
    return resourcePool.jobs?.[0] || "1";
  }
  if (lower.includes("talent") || lower.includes("candidate")) {
    return resourcePool.talents?.[0] || "1";
  }
  if (lower.includes("interview")) {
    return resourcePool.interviews?.[0] || "1";
  }
  if (lower.includes("notice")) {
    return resourcePool.notices?.[0] || resourcePool.notice?.[0] || "1";
  }
  if (lower.includes("message")) {
    return resourcePool.messages?.[0] || "1";
  }
  if (lower.includes("user") || lower.includes("interviewer")) {
    return "1";
  }
  if (lower.includes("filename")) {
    return "test.pdf";
  }
  return "1";
}

function buildQueryString(parameters, negative) {
  const params = new URLSearchParams();
  for (const parameter of parameters.filter((item) => item.in === "query")) {
    if (!parameter.required && !["page", "page_size", "status", "search", "keyword"].includes(parameter.name)) {
      continue;
    }
    params.set(parameter.name, chooseParameterValue(parameter, negative));
  }
  const value = params.toString();
  return value ? `?${value}` : "";
}

function chooseParameterValue(parameter, negative) {
  if (negative) {
    return parameter.type === "integer" ? "-1" : "'; DROP TABLE bigtest; --";
  }
  const name = String(parameter.name || "").toLowerCase();
  if (name === "page") {
    return "1";
  }
  if (name === "page_size") {
    return "10";
  }
  if (name.includes("status")) {
    return "open";
  }
  if (name.includes("search") || name.includes("keyword")) {
    return "测试";
  }
  if (parameter.type === "integer") {
    return "1";
  }
  if (parameter.type === "array") {
    return "Go";
  }
  return "测试";
}

function buildBodyForOperation(runContext, operation, options) {
  if (!["POST", "PUT", "PATCH"].includes(operation.method)) {
    return null;
  }
  if (options.negative) {
    return buildNegativeBody(operation);
  }
  if (isLoginApiPath(operation.path)) {
    return {
      username: options.profile.auth.username,
      password: options.profile.auth.password,
    };
  }
  const resource = parseApiResource(operation.path)?.name || "";
  return buildResourceBody(runContext, resource, operation);
}

function buildResourceBody(runContext, resource, operation) {
  const prefix = "${prefix}";
  const common = {
    jobs: {
      title: `${prefix}_自动职位`,
      description: `${prefix}_自动生成职位描述`,
      requirements: ["Go", "Vue"],
      salary: "10k-20k",
      location: "测试城市",
      type: "full-time",
      status: "open",
      department: "测试部",
      level: "mid",
      education: "本科",
      skills: ["Go", "Vue"],
      benefits: ["五险一金"],
      headcount: 1,
    },
    talents: {
      name: `${prefix}_自动人才`,
      email: "bigtest_${runId}@example.com",
      phone: "13800000000",
      skills: ["Go", "Vue"],
      experience: 3,
      education: "本科",
      status: "active",
      tags: ["BIGTEST"],
      location: "测试城市",
      salary: "15k-25k",
      summary: `${prefix}_自动人才摘要`,
      gender: "未知",
      age: 28,
      current_company: "BIGTEST 公司",
      current_position: "后端工程师",
      source: "bigtest",
    },
    interviews: {
      candidate_id: 1,
      candidate_name: `${prefix}_候选人`,
      position_id: 1,
      position: `${prefix}_职位`,
      type: "initial",
      date: "${futureDate}",
      time: "${futureTime}",
      duration: 60,
      interviewer_id: 1,
      interviewer: "BIGTEST 面试官",
      method: "video",
      location: "线上会议",
      notes: `${prefix}_面试备注`,
      created_by: 1,
    },
    messages: {
      receiver_id: 1,
      title: `${prefix}_消息`,
      content: `${prefix}_自动消息内容`,
      type: "system",
      is_read: false,
    },
    notices: {
      title: `${prefix}_公告`,
      content: `${prefix}_自动公告内容`,
      status: "draft",
      is_pinned: false,
      priority: "normal",
    },
    conversations: {
      participant_id: 1,
      title: `${prefix}_会话`,
      content: `${prefix}_会话消息`,
    },
    recommendations: {
      job_id: 1,
      talent_id: 1,
      limit: 10,
      query: "Go 后端",
    },
    resumes: {
      name: `${prefix}_简历`,
      title: `${prefix}_简历`,
      content: `${prefix}_简历内容`,
      status: "active",
    },
    applications: {
      job_id: 1,
      talent_id: 1,
      resume_id: 1,
      status: "pending",
    },
    profile: {
      name: `${prefix}_用户`,
      phone: "13800000000",
    },
  };

  const body = common[resource] || { name: `${prefix}_${resource || "resource"}`, content: `${prefix}_自动测试内容` };
  if (operation.path.includes("/reschedule")) {
    return { date: "${futureDate2}", time: "${futureTime2}", reason: `${prefix}_自动改期` };
  }
  if (operation.path.includes("/complete")) {
    return { feedback: `${prefix}_自动完成`, rating: 4 };
  }
  if (operation.path.includes("/cancel")) {
    return { reason: `${prefix}_自动取消` };
  }
  if (operation.path.includes("/feedback")) {
    return { rating: 4, strengths: "稳定", weaknesses: "无", comments: `${prefix}_反馈`, recommendation: "pass" };
  }
  return interpolate(body, buildTemplateVariables(runContext, "auto"));
}

function buildNegativeBody(operation) {
  const resource = parseApiResource(operation.path)?.name || "";
  if (isLoginApiPath(operation.path)) {
    return { username: "", password: "" };
  }
  if (["jobs", "talents", "interviews", "notices", "messages"].includes(resource)) {
    return { name: "", title: "", email: "not-an-email", status: "__invalid__", id: "abc" };
  }
  return { invalid: true, id: "abc", status: "__invalid__" };
}

function chooseNegativeExpectedStatuses(operation) {
  if (operation.requiresAuth) {
    return [400, 401, 403, 404, 422];
  }
  return [400, 404, 422];
}

function buildTemplateVariables(runContext, profileName) {
  return {
    runId: runContext.runId,
    prefix: runContext.testDataPrefix,
    profile: profileName,
    now: new Date().toISOString(),
    futureDate: formatDateDaysFromNow(2),
    futureDate2: formatDateDaysFromNow(3),
    futureTime: "10:00",
    futureTime2: "15:00",
  };
}

async function replayGeneratedCases(runContext, profile, generatedCases, token) {
  const replayConfig = { ...runContext.config.replay, ...profile.replay };
  const allowedMethods = new Set(
    (replayConfig.allowedMethods || []).map((item) => item.toUpperCase()),
  );
  const results = [];
  const runtimeResourcePool = {};

  for (const testCase of generatedCases) {
    const method = testCase.request.method.toUpperCase();
    if (!allowedMethods.has(method)) {
      results.push({
        id: testCase.id,
        path: testCase.request.normalizedPath,
        method,
        skipped: true,
        reason: `method ${method} not enabled`,
      });
      continue;
    }
    const writeDecision = validateWriteReplay(runContext, replayConfig, testCase);
    if (!writeDecision.allowed) {
      results.push({
        id: testCase.id,
        path: testCase.request.normalizedPath,
        method,
        skipped: true,
        reason: writeDecision.reason,
      });
      continue;
    }
    const replayPath = resolveReplayPath(testCase, runtimeResourcePool);
    if (
      (replayConfig.ignoreUrlPatterns || []).some((item) =>
        replayPath.includes(item),
      )
    ) {
      results.push({
        id: testCase.id,
        path: replayPath,
        method,
        skipped: true,
        reason: "matched replay ignore pattern",
      });
      continue;
    }

    const replayUrl = new URL(
      replayPath,
      runContext.config.env.replayBaseUrl,
    ).toString();
    const headers = buildReplayHeaders(
      testCase.request.headersSanitized,
      testCase.request.requiresAuth && !testCase.request.forceNoAuth,
      token,
    );

    const init = {
      method,
      headers,
      body: buildReplayBody(testCase, method, profile),
    };

    const startedAt = Date.now();
    try {
      if (replayConfig.throttleMs > 0) {
        await delay(replayConfig.throttleMs);
      }
      const response = await fetch(replayUrl, init);
      const contentType = response.headers.get("content-type") || "";
      const text = await response.text();
      const parsedJson = tryParseJson(text);
      const actualShape = parsedJson ? summarizeJsonShape(parsedJson) : null;
      const expectedStatuses = Array.isArray(testCase.expected.status)
        ? testCase.expected.status
        : [testCase.expected.status];
      const statusMatches = evaluateStatusMatch(testCase, response.status, expectedStatuses);
      const inconclusiveReason = getInconclusiveReplayReason(testCase, response.status);
      const inconclusive = Boolean(inconclusiveReason);
      const shapeMatches = replayConfig.compareResponseShape
        ? areShapesCompatible(testCase.expected.responseShape, actualShape)
        : true;
      const contentTypeMatches = replayConfig.compareContentType
        ? matchContentType(testCase.expected.contentType, contentType)
        : true;
      const passed = !inconclusive && statusMatches && shapeMatches && contentTypeMatches;

      if (passed && parsedJson) {
        registerRuntimeResourceFromResponse(runtimeResourcePool, testCase, replayPath, parsedJson);
      }

      results.push({
        id: testCase.id,
        method,
        path: replayPath,
        replayUrl,
        expectedStatus: expectedStatuses.join(", "),
        actualStatus: response.status,
        durationMs: Date.now() - startedAt,
        statusMatches,
        shapeMatches,
        contentTypeMatches,
        passed,
        skipped: inconclusive,
        reason: inconclusiveReason,
        source: testCase.source || "capture",
        negative: Boolean(testCase.negative),
        responsePreview: parsedJson ? sanitizeValue(parsedJson) : text.slice(0, 2000),
        curlCommand: passed || inconclusive ? "" : buildCurlCommand(replayUrl, init),
      });
    } catch (error) {
      results.push({
        id: testCase.id,
        method,
        path: replayPath,
        replayUrl,
        expectedStatus: Array.isArray(testCase.expected.status)
          ? testCase.expected.status.join(", ")
          : testCase.expected.status,
        actualStatus: 0,
        durationMs: Date.now() - startedAt,
        statusMatches: false,
        shapeMatches: false,
        contentTypeMatches: false,
        passed: false,
        source: testCase.source || "capture",
        negative: Boolean(testCase.negative),
        error: error.message,
        curlCommand: buildCurlCommand(replayUrl, init),
      });
    }
  }

  return results;
}

function resolveReplayPath(testCase, runtimeResourcePool) {
  if (!testCase.request.pathTemplate) {
    return testCase.request.normalizedPath;
  }
  const pathValue = materializeApiPath(
    testCase.request.pathTemplate,
    { path: testCase.request.pathTemplate, method: testCase.request.method },
    runtimeResourcePool,
    { negative: Boolean(testCase.negative) },
  );
  return `${pathValue}${testCase.request.queryString || ""}`;
}

function evaluateStatusMatch(testCase, actualStatus, expectedStatuses) {
  if (testCase.negative) {
    if (testCase.negativeSubtype === "unauthorized") {
      return [401, 403].includes(actualStatus);
    }
    // 非法输入类负向用例：期望后端明确拒绝，不应成功，也不应 5xx。
    return actualStatus >= 400 && actualStatus < 500;
  }
  return expectedStatuses.includes(actualStatus);
}

function getInconclusiveReplayReason(testCase, actualStatus) {
  if (actualStatus === 429 && (testCase.source === "openapi" || testCase.negative)) {
    return "请求被限流，本轮无法判断接口行为";
  }
  if ((testCase.source || "capture") !== "openapi" || testCase.negative) {
    return "";
  }
  if ([400, 404, 409, 422].includes(actualStatus)) {
    return "OpenAPI 自动用例缺少业务前置数据或字段模板";
  }
  if ([401, 403].includes(actualStatus)) {
    return "OpenAPI 自动用例缺少对应角色或资源权限";
  }
  return "";
}

function delay(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

function registerRuntimeResourceFromResponse(runtimeResourcePool, testCase, replayPath, body) {
  const resourceName =
    testCase.request.resourceName ||
    parseApiResource(testCase.request.pathTemplate || replayPath)?.name;
  const id = getByPath(body, "data.id") || getByPath(body, "id");
  if (!resourceName || !id) {
    return;
  }
  if (!runtimeResourcePool[resourceName]) {
    runtimeResourcePool[resourceName] = [];
  }
  const stringId = String(id);
  if (!runtimeResourcePool[resourceName].includes(stringId)) {
    runtimeResourcePool[resourceName].unshift(stringId);
  }
}

function buildReplayHeaders(headersSanitized, requiresAuth, token) {
  const output = {};
  for (const [key, value] of Object.entries(headersSanitized || {})) {
    const lower = key.toLowerCase();
    if (
      lower === "host" ||
      lower === "content-length" ||
      lower === "connection" ||
      lower === "origin" ||
      lower === "referer" ||
      lower === "accept-encoding" ||
      lower.startsWith("sec-") ||
      lower === "cookie" ||
      lower === "authorization"
    ) {
      continue;
    }
    if (value !== "<redacted>") {
      output[key] = value;
    }
  }
  if (requiresAuth && token) {
    output.Authorization = `Bearer ${token}`;
  }
  return output;
}

function buildReplayBody(testCase, method, profile) {
  if (["GET", "HEAD"].includes(method) || !testCase.request.bodySanitized) {
    return undefined;
  }

  const body = !testCase.negative && isLoginApiPath(testCase.request.normalizedPath) && isPlainRecord(testCase.request.bodySanitized)
    ? {
        ...testCase.request.bodySanitized,
        username: profile.auth.username || testCase.request.bodySanitized.username,
        password: profile.auth.password || testCase.request.bodySanitized.password,
      }
    : testCase.request.bodySanitized;

  return typeof body === "string" ? body : jsonStringify(body);
}

function isLoginApiPath(pathValue) {
  return String(pathValue || "").split("?")[0] === "/api/v1/login";
}

function buildCurlCommand(url, init) {
  const parts = ["curl", "-i", "-X", shellQuote(init.method || "GET"), shellQuote(url)];
  for (const [key, value] of Object.entries(init.headers || {})) {
    parts.push("-H", shellQuote(`${key}: ${sanitizeHeaderForCurl(key, value)}`));
  }
  if (init.body !== undefined) {
    parts.push("--data", shellQuote(sanitizeCurlBody(init.body)));
  }
  return parts.join(" ");
}

function sanitizeHeaderForCurl(key, value) {
  return containsSensitiveKeyName(key) ? "<redacted>" : value;
}

function sanitizeCurlBody(body) {
  const parsed = tryParseJson(body);
  return parsed ? jsonStringify(sanitizeValue(parsed)) : String(body);
}

function shellQuote(value) {
  return `'${String(value).replace(/'/g, "'\\''")}'`;
}

function matchContentType(expected, actual) {
  if (!expected || !actual) {
    return true;
  }
  const expectedMain = expected.split(";")[0].trim();
  const actualMain = actual.split(";")[0].trim();
  return expectedMain === actualMain;
}

async function executeWriteScenarios(runContext, profile, token, generatedCases) {
  const safety = {
    ...runContext.config.writeSafety,
    ...(profile.writeSafety || {}),
  };
  const manualScenarios = safety.scenarios || [];
  const inferredScenarios = safety.inferFromCaptures
    ? inferWriteScenarios(generatedCases, safety)
    : [];
  const scenarios = [...manualScenarios, ...inferredScenarios].slice(
    0,
    safety.maxInferredScenarios
      ? manualScenarios.length + safety.maxInferredScenarios
      : undefined,
  );
  const results = [];

  for (const scenario of scenarios) {
    const result = await executeWriteScenario(runContext, profile, safety, scenario, token);
    results.push(result);
  }

  if (safety.autoCleanup) {
    for (const scenario of scenarios) {
      await cleanupScenarioResource(runContext, profile, safety, scenario, token, results);
    }
  }

  return results;
}

function inferWriteScenarios(generatedCases, safety) {
  const byResource = new Map();
  for (const testCase of generatedCases) {
    const method = testCase.request.method.toUpperCase();
    const resource = parseApiResource(testCase.request.normalizedPath);
    if (!resource || isUnsafeWriteCandidate(testCase.request.normalizedPath, safety)) {
      continue;
    }
    if (!byResource.has(resource.name)) {
      byResource.set(resource.name, { resource, cases: [] });
    }
    byResource.get(resource.name).cases.push(testCase);
  }

  const scenarios = [];
  for (const { resource, cases } of byResource.values()) {
    const createCase = cases.find(
      (item) =>
        item.request.method.toUpperCase() === "POST" &&
        parseApiResource(item.request.normalizedPath)?.isCollection,
    );
    if (!createCase || !isPlainRecord(createCase.request.bodySanitized)) {
      continue;
    }

    const body = buildBigtestBodyTemplate(createCase.request.bodySanitized);
    if (!containsValue(body, "${prefix}")) {
      continue;
    }

    scenarios.push({
      name: `${resource.name}-inferred-owned-data`,
      resourceType: resource.name,
      inferred: true,
      create: {
        method: "POST",
        path: resource.collectionPath,
        body,
        idPath: "data.id",
        expectedStatuses: [200, 201],
      },
      verify: {
        method: "GET",
        path: `${resource.collectionPath}/\${id}`,
        expectedStatuses: [200],
      },
      update: {
        method: "PUT",
        path: `${resource.collectionPath}/\${id}`,
        body: buildBigtestBodyTemplate(createCase.request.bodySanitized, "_已更新"),
        expectedStatuses: [200],
      },
      delete: {
        method: "DELETE",
        path: `${resource.collectionPath}/\${id}`,
        expectedStatuses: [200, 204],
      },
    });
  }

  return scenarios;
}

function parseApiResource(normalizedPath) {
  const pathname = normalizedPath.split("?")[0];
  const segments = pathname.split("/").filter(Boolean);
  const apiIndex = segments.findIndex((segment) => segment === "api");
  const startIndex = apiIndex >= 0 && segments[apiIndex + 1]?.startsWith("v")
    ? apiIndex + 2
    : apiIndex >= 0
      ? apiIndex + 1
      : 0;
  const resourceName = segments[startIndex];
  if (!resourceName || ["login", "logout", "auth", "token", "upload"].includes(resourceName)) {
    return null;
  }
  const collectionSegments = segments.slice(0, startIndex + 1);
  const rest = segments.slice(startIndex + 1);
  return {
    name: resourceName,
    collectionPath: `/${collectionSegments.join("/")}`,
    isCollection: rest.length === 0,
  };
}

function isUnsafeWriteCandidate(normalizedPath, safety) {
  return (safety.blockedUrlPatterns || []).some((pattern) => normalizedPath.includes(pattern));
}

function buildBigtestBodyTemplate(body, suffix = "") {
  const output = {};
  for (const [key, value] of Object.entries(body || {})) {
    if (containsSensitiveKeyName(key)) {
      continue;
    }
    if (typeof value === "string") {
      output[key] = chooseBodyTemplateValue(key, value, suffix);
      continue;
    }
    if (typeof value === "boolean" || typeof value === "number") {
      output[key] = value;
      continue;
    }
    if (Array.isArray(value)) {
      output[key] = [];
      continue;
    }
    if (isPlainRecord(value)) {
      output[key] = buildBigtestBodyTemplate(value, suffix);
    }
  }
  return output;
}

function chooseBodyTemplateValue(key, originalValue, suffix) {
  const lower = key.toLowerCase();
  if (/email|mail/.test(lower)) {
    return "bigtest_${runId}@example.com";
  }
  if (/phone|mobile|tel/.test(lower)) {
    return "13800000000";
  }
  if (/status|type|priority|role/.test(lower)) {
    return originalValue;
  }
  if (/date|time/.test(lower)) {
    return "${now}";
  }
  return `\${prefix}_${key}${suffix}`;
}

function isPlainRecord(value) {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function containsSensitiveKeyName(key) {
  return /password|token|secret|credential|authorization|cookie/i.test(key);
}

async function executeWriteScenario(runContext, profile, safety, scenario, token) {
  const result = {
    name: scenario.name,
    resourceType: scenario.resourceType || scenario.name,
    inferred: Boolean(scenario.inferred),
    steps: [],
  };
  const variables = {
    runId: runContext.runId,
    prefix: runContext.testDataPrefix,
    profile: profile.name,
    now: new Date().toISOString(),
    futureDate: formatDateDaysFromNow(2),
    futureDate2: formatDateDaysFromNow(3),
    futureTime: "10:00",
    futureTime2: "15:00",
  };

  try {
    if (!scenario.create) {
      throw new Error("写操作场景缺少 create 步骤");
    }

    const createResult = await runSafeApiStep(
      runContext,
      safety,
      scenario.create,
      token,
      variables,
      { allowCreate: true },
    );
    result.steps.push({ name: "create", ...createResult });

    const id = getByPath(createResult.body, scenario.create.idPath || "data.id");
    if (!id) {
      throw new Error(`create 响应中没有取到资源 ID: ${scenario.create.idPath || "data.id"}`);
    }
    variables.id = id;
    registerOwnedResource(runContext, result.resourceType, id, createResult.path);

    for (const stepName of ["verify", "update"]) {
      if (!scenario[stepName]) {
        continue;
      }
      const stepResult = await runSafeApiStep(
        runContext,
        safety,
        scenario[stepName],
        token,
        variables,
        { requireOwnedId: true },
      );
      result.steps.push({ name: stepName, ...stepResult });
    }

    for (const customStep of scenario.steps || []) {
      const stepResult = await runSafeApiStep(
        runContext,
        safety,
        customStep,
        token,
        variables,
        { requireOwnedId: true },
      );
      result.steps.push({ name: customStep.name || "step", ...stepResult });
    }

    result.passed = result.steps.every((step) => step.passed);
  } catch (error) {
    result.passed = false;
    result.error = error.message;
  }

  return result;
}

async function cleanupScenarioResource(runContext, profile, safety, scenario, token, results) {
  if (!scenario.delete) {
    return;
  }
  const resourceType = scenario.resourceType || scenario.name;
  const ids = [...(runContext.registry.resources[resourceType] || [])];

  for (const id of ids) {
    const variables = {
      runId: runContext.runId,
      prefix: runContext.testDataPrefix,
      profile: profile.name,
      id,
      now: new Date().toISOString(),
    };
    const targetResult = results.find((item) => item.name === scenario.name);
    try {
      const cleanup = await runSafeApiStep(
        runContext,
        safety,
        scenario.delete,
        token,
        variables,
        { requireOwnedId: true },
      );
      if (targetResult) {
        targetResult.steps.push({ name: "delete", ...cleanup });
        targetResult.cleanupPassed = cleanup.passed;
        targetResult.passed = targetResult.passed && cleanup.passed;
      }
      if (cleanup.passed) {
        unregisterOwnedResource(runContext, resourceType, id);
      }
    } catch (error) {
      if (targetResult) {
        targetResult.cleanupPassed = false;
        targetResult.passed = false;
        targetResult.steps.push({
          name: "delete",
          passed: false,
          error: error.message,
        });
      }
    }
  }
}

async function runSafeApiStep(runContext, safety, step, token, variables, guard) {
  const method = (step.method || "GET").toUpperCase();
  const pathTemplate = step.path || "/";
  const requestPath = interpolate(pathTemplate, variables);
  const body = interpolate(step.body || null, variables);

  assertSafeWriteRequest(runContext, safety, {
    method,
    path: requestPath,
    body,
    guard,
  });

  const url = new URL(requestPath, runContext.config.env.replayBaseUrl).toString();
  const headers = {
    "content-type": "application/json",
  };
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const startedAt = Date.now();
  const response = await fetch(url, {
    method,
    headers,
    body: ["GET", "HEAD"].includes(method) || body === null ? undefined : jsonStringify(body),
  });
  const text = await response.text();
  const parsed = tryParseJson(text);
  const expectedStatuses = step.expectedStatuses || step.expectedStatus || [200];
  const allowedStatuses = Array.isArray(expectedStatuses)
    ? expectedStatuses
    : [expectedStatuses];
  const passed = allowedStatuses.includes(response.status);

  return {
    method,
    path: requestPath,
    status: response.status,
    expectedStatuses: allowedStatuses,
    durationMs: Date.now() - startedAt,
    passed,
    body: parsed,
    bodyPreview: parsed ? sanitizeValue(parsed) : text.slice(0, 1000),
  };
}

function assertSafeWriteRequest(runContext, safety, request) {
  if ((safety.blockedUrlPatterns || []).some((pattern) => request.path.includes(pattern))) {
    throw new Error(`路径已排除，不执行: ${request.path}`);
  }

  if (!isWriteMethod(request.method)) {
    return;
  }

  const allowedMethods = new Set((safety.allowMethods || []).map((item) => item.toUpperCase()));
  if (!allowedMethods.has(request.method)) {
    throw new Error(`写操作场景未允许 ${request.method}`);
  }

  if (request.guard?.allowCreate) {
    if (!containsValue(request.body, runContext.testDataPrefix)) {
      throw new Error("创建数据必须包含本轮 BIGTEST_ 前缀");
    }
    return;
  }

  if (request.guard?.requireOwnedId) {
    const id = request.guard.id || extractLastPathId(request.path);
    if (!id || !isOwnedResource(runContext, id)) {
      throw new Error(`只能修改/删除本轮已登记资源: ${request.path}`);
    }
  }
}

function validateWriteReplay(runContext, replayConfig, testCase) {
  const method = testCase.request.method.toUpperCase();
  if (!isWriteMethod(method)) {
    return { allowed: true };
  }

  const safety = runContext.config.writeSafety || {};
  if (!safety.enabled) {
    return { allowed: false, reason: "writeSafety disabled" };
  }
  if ((safety.blockedUrlPatterns || []).some((pattern) => testCase.request.normalizedPath.includes(pattern))) {
    return { allowed: false, reason: "matched excluded path pattern" };
  }
  if (safety.allowUnsafeReplayWrites) {
    return { allowed: true };
  }

  if (method === "POST") {
    const hasPrefix = containsValue(testCase.request.bodySanitized, runContext.testDataPrefix);
    return hasPrefix
      ? { allowed: true }
      : { allowed: false, reason: "POST body is not BIGTEST-owned" };
  }

  const id = extractLastPathId(testCase.request.normalizedPath);
  if (!id || !isOwnedResource(runContext, id)) {
    return { allowed: false, reason: "PUT/DELETE target is not registered BIGTEST data" };
  }

  return { allowed: true };
}

function isWriteMethod(method) {
  return ["POST", "PUT", "PATCH", "DELETE"].includes(method);
}

function registerOwnedResource(runContext, resourceType, id, resourcePath) {
  const key = String(resourceType);
  if (!runContext.registry.resources[key]) {
    runContext.registry.resources[key] = [];
  }
  runContext.registry.resources[key].push(String(id));
  logInfo(`登记 BIGTEST 资源 ${key}#${id}: ${resourcePath}`);
}

function unregisterOwnedResource(runContext, resourceType, id) {
  const key = String(resourceType);
  runContext.registry.resources[key] = (runContext.registry.resources[key] || []).filter(
    (item) => item !== String(id),
  );
}

function isOwnedResource(runContext, id) {
  return Object.values(runContext.registry.resources).some((ids) =>
    ids.map(String).includes(String(id)),
  );
}

function extractLastPathId(pathValue) {
  const matches = String(pathValue).match(/\/([0-9a-fA-F-]{1,64})(?=\/|$)/g);
  if (!matches?.length) {
    return "";
  }
  return matches[matches.length - 1].replace("/", "");
}

function formatDateDaysFromNow(days) {
  const date = new Date(Date.now() + days * 86400000);
  const pad = (value) => String(value).padStart(2, "0");
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`;
}

function containsValue(value, expected) {
  if (value === null || value === undefined) {
    return false;
  }
  if (typeof value === "string") {
    return value.includes(expected);
  }
  if (Array.isArray(value)) {
    return value.some((item) => containsValue(item, expected));
  }
  if (typeof value === "object") {
    return Object.values(value).some((item) => containsValue(item, expected));
  }
  return false;
}

async function saveDebugScreenshot(page, debugDir, name) {
  try {
    await ensureDir(debugDir);
    const filePath = path.join(debugDir, `${safeFileName(name)}.png`);
    await page.screenshot({ path: filePath, fullPage: true });
    return filePath;
  } catch {
    return "";
  }
}

function shouldSaveScreenshot(runContext, debugArtifacts, type) {
  const debugConfig = runContext.config.debug || {};
  if (type === "button" && debugConfig.screenshotOnButtonFailure !== true) {
    return false;
  }
  if (type === "route" && debugConfig.screenshotOnRouteFailure === false) {
    return false;
  }
  const maxScreenshots = Number(debugConfig.maxScreenshots ?? 5);
  return maxScreenshots > 0 && (debugArtifacts?.screenshots?.length || 0) < maxScreenshots;
}

function safeFileName(value) {
  return String(value || "debug").replace(/[^a-zA-Z0-9_-]+/g, "-").slice(0, 80);
}

export async function writeArtifacts(runContext, results) {
  const outputFiles = [];
  for (const result of results) {
    const profileDir = path.join(runContext.outputRoot, result.profileName);
    await ensureDir(profileDir);

    if (runContext.config.generation.writeRawCaptures) {
      const rawCaptureFile = path.join(profileDir, "raw-captures.json");
      await writeJson(rawCaptureFile, {
        profile: result.profileName,
        generatedAt: new Date().toISOString(),
        user: sanitizeValue(result.rawCapture.user),
        exploration: result.rawCapture.exploration,
        captures: result.rawCapture.captures.map((item) => ({
          ...item,
          request: {
            ...item.request,
            headers: item.request.headersSanitized,
            postData: item.request.postDataSanitized,
          },
        })),
      });
      outputFiles.push(rawCaptureFile);
    }

    const debugFile = path.join(profileDir, "debug-artifacts.json");
    await writeJson(debugFile, result.rawCapture.debugArtifacts || {
      consoleErrors: [],
      networkFailures: [],
      screenshots: [],
    });
    outputFiles.push(debugFile);

    if (runContext.config.generation.writeGeneratedCases) {
      const casesFile = path.join(profileDir, "generated-cases.json");
      await writeJson(casesFile, result.generatedCases);
      outputFiles.push(casesFile);
    }

    const discoveryFile = path.join(profileDir, "discovery-results.json");
    await writeJson(discoveryFile, {
      exploration: result.rawCapture.exploration,
      apiCoverage: summarizeApiCoverage(result.generatedCases, result.writeScenarioResults),
      safetyBlocks: summarizeSafetyBlocks(result),
    });
    outputFiles.push(discoveryFile);

    if (runContext.config.writeSafety.enabled) {
      const writeSafetyFile = path.join(profileDir, "write-safety-results.json");
      await writeJson(writeSafetyFile, {
        testDataPrefix: runContext.testDataPrefix,
        registry: runContext.registry,
        scenarios: result.writeScenarioResults,
      });
      outputFiles.push(writeSafetyFile);
    }

    const replayFile = path.join(profileDir, "replay-results.json");
    await writeJson(replayFile, result.replayResults);
    outputFiles.push(replayFile);
  }

  const reportFile = path.join(runContext.outputRoot, "report.md");
  await writeText(reportFile, await buildMarkdownReport(runContext, results));
  outputFiles.push(reportFile);

  return { reportFile, outputFiles };
}

async function buildMarkdownReport(runContext, results) {
  const totalMetrics = summarizeRunMetrics(results);
  const coverageBaseline = await loadCoverageBaseline(runContext);
  const lines = [
    "# bigtest 真实流量自动化测试报告",
    "",
    `- Run ID: \`${runContext.runId}\``,
    `- 配置文件: \`${runContext.configPath}\``,
    `- 输出目录: \`${relativeFrom(runContext.bigtestRoot, runContext.outputRoot)}\``,
    `- 自动探索: ${runContext.config.autoExplore.enabled ? "已开启" : "未开启"}`,
    `- 写操作场景: ${runContext.config.writeSafety.enabled ? `已开启，测试数据前缀 \`${runContext.testDataPrefix}\`` : "未开启"}`,
    "",
    "## 测试结论",
    "",
    `本次测试共访问 ${totalMetrics.pagesVisited} 个页面，捕获 ${totalMetrics.rawCaptures} 条真实 API 请求，生成 ${totalMetrics.generatedCases} 条接口回放用例。`,
    `接口回放通过 ${totalMetrics.replayPass} 条，失败 ${totalMetrics.replayFail} 条，跳过 ${totalMetrics.replaySkip} 条；写操作场景通过 ${totalMetrics.writePass} 个，失败 ${totalMetrics.writeFail} 个。`,
    `用例来源：真实流量 ${totalMetrics.caseSources.capture || 0} 条，OpenAPI 自动生成 ${totalMetrics.caseSources.openapi || 0} 条，负向异常用例 ${totalMetrics.caseSources.negative || 0} 条。`,
    `失败分层：真实失败 ${totalMetrics.failureCategories.real || 0} 条，疑似依赖/环境失败 ${totalMetrics.failureCategories.dependency || 0} 条，契约不一致 ${totalMetrics.failureCategories.contract || 0} 条，负向用例失败 ${totalMetrics.failureCategories.negative || 0} 条，自动/限流用例跳过 ${totalMetrics.replaySkip} 条。`,
    "",
    totalMetrics.replayFail === 0 && totalMetrics.writeFail === 0
      ? "**结论：本次自动化测试未发现接口回放失败或写操作场景失败。**"
      : "**结论：本次自动化测试存在失败项，请优先查看“回放结果”和“写操作场景步骤”。**",
    "",
    "## 测试环境与策略",
    "",
    "| 项目 | 值 |",
    "| --- | --- |",
    `| 前端地址 | \`${runContext.config.env.frontendBaseUrl}\` |`,
    `| 后端回放地址 | \`${runContext.config.env.replayBaseUrl}\` |`,
    `| 捕获 URL 规则 | \`${(runContext.config.capture.urlIncludes || []).join(", ")}\` |`,
    `| 回放方法 | \`${(runContext.config.replay.allowedMethods || []).join(", ")}\` |`,
    `| 回放节流 | ${runContext.config.replay.throttleMs || 0} ms/请求 |`,
    `| 排除路径 | \`${[
      ...(runContext.config.capture.ignoreUrlPatterns || []),
      ...(runContext.config.replay.ignoreUrlPatterns || []),
    ].filter((item, index, all) => all.indexOf(item) === index).join(", ")}\` |`,
    `| 页面探索最大路由数 | ${runContext.config.autoExplore.maxRoutes} |`,
    `| 每页最大点击数 | ${runContext.config.autoExplore.maxClicksPerPage} |`,
    `| 截图策略 | 最多 ${runContext.config.debug?.maxScreenshots ?? 5} 张；按钮点击失败默认不截图，只记录 URL 和错误 |`,
    `| AI 评估策略 | 明确排除 \`/api/v1/ai\`、\`/api/v1/evaluations\`、\`/ai/evaluate\`、\`/evaluate\` 等路径，不触发真实 AI 评估、上传或外部模型调用 |`,
    "",
    "## 总览",
    "",
    "| Profile | 页面 | 输入 | 按钮 | 原始请求 | 生成用例 | 写操作通过 | 写操作失败 | 回放通过 | 回放失败 | 回放跳过 |",
    "| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |",
  ];

  for (const result of results) {
    const exploration = result.rawCapture.exploration || createEmptyExplorationReport();
    const writePass = (result.writeScenarioResults || []).filter((item) => item.passed).length;
    const writeFail = (result.writeScenarioResults || []).filter((item) => !item.passed).length;
    const replayPass = result.replayResults.filter((item) => item.passed).length;
    const replayFail = result.replayResults.filter(
      (item) => !item.skipped && item.passed === false,
    ).length;
    const replaySkip = result.replayResults.filter((item) => item.skipped).length;
    lines.push(
      `| ${result.profileName} | ${exploration.pagesVisited} | ${exploration.inputsFound}/${exploration.inputsFilled} | ${exploration.buttonsFound}/${exploration.buttonsClicked} | ${result.rawCapture.captures.length} | ${result.generatedCases.length} | ${writePass} | ${writeFail} | ${replayPass} | ${replayFail} | ${replaySkip} |`,
    );
  }

  appendModuleCoverageReport(lines, runContext, results, coverageBaseline);
  appendUncoveredReport(lines, runContext, results, coverageBaseline);

  for (const result of results) {
    lines.push("", `## ${result.profileName}`, "");
    appendDiscoveryReport(lines, result);
    if (runContext.config.writeSafety.enabled) {
      const deleteStepCount = countWriteSteps(result.writeScenarioResults, "DELETE");
      lines.push("### 写操作场景结果", "");
      lines.push(`本轮写操作场景已执行 ${deleteStepCount} 次 \`DELETE\` 删除请求。`, "");
      lines.push("| 场景 | 来源 | 状态 | 资源 | 步骤数 | 说明 |");
      lines.push("| --- | --- | --- | --- | ---: | --- |");
      for (const item of result.writeScenarioResults || []) {
        lines.push(
          `| ${mdCell(item.name)} | ${item.inferred ? "自动推断" : "配置"} | ${item.passed ? "PASS" : "FAIL"} | ${mdCell(item.resourceType)} | ${item.steps?.length || 0} | ${mdCell(item.error || "BIGTEST 数据创建、更新、删除完成")} |`,
        );
      }
      if ((result.writeScenarioResults || []).length === 0) {
        lines.push("| - | SKIP | - | 0 | 未配置写操作场景 |");
      }
      lines.push("");
      lines.push("### 写操作场景步骤", "");
      lines.push("| 场景 | 步骤 | 方法 | 路径 | 期望 | 实际 | 状态 |");
      lines.push("| --- | --- | --- | --- | --- | ---: | --- |");
      let hasWriteSteps = false;
      for (const item of result.writeScenarioResults || []) {
        for (const step of item.steps || []) {
          hasWriteSteps = true;
          lines.push(
            `| ${mdCell(item.name)} | ${mdCell(formatWriteStepName(step.name, step.method))} | ${step.method} | \`${mdCell(step.path)}\` | ${mdCell((step.expectedStatuses || []).join(", "))} | ${step.status ?? "-"} | ${step.passed ? "PASS" : "FAIL"} |`,
          );
        }
      }
      if (!hasWriteSteps) {
        lines.push("| - | - | - | - | - | - | SKIP |");
      }
      lines.push("");
    }
    lines.push("### 去重后生成的用例", "");
    lines.push("| 用例ID | 来源 | 方法 | 路径 | 归并次数 | 需鉴权 |");
    lines.push("| --- | --- | --- | --- | ---: | --- |");
    for (const testCase of result.generatedCases) {
      lines.push(
        `| ${mdCell(testCase.id)} | ${mdCell(describeCaseSource(testCase.source || "capture"))} | ${testCase.request.method} | \`${mdCell(testCase.request.normalizedPath)}\` | ${testCase.occurrences} | ${testCase.request.requiresAuth ? "是" : "否"} |`,
      );
    }

    lines.push("", "### 回放结果", "");
    lines.push("| 用例ID | 方法 | 路径 | 状态 | 期望 | 实际 | 耗时(ms) | 说明 |");
    lines.push("| --- | --- | --- | --- | ---: | ---: | ---: | --- |");
    for (const replay of result.replayResults) {
      const status = replay.skipped ? "SKIP" : replay.passed ? "PASS" : "FAIL";
      const note = replay.skipped
        ? replay.reason
        : replay.error ||
          [
            replay.statusMatches ? null : "status mismatch",
            replay.shapeMatches ? null : "shape mismatch",
            replay.contentTypeMatches ? null : "content-type mismatch",
          ]
            .filter(Boolean)
            .join(", ") ||
          "ok";

      lines.push(
        `| ${mdCell(replay.id)} | ${mdCell(replay.method || "-")} | \`${mdCell(replay.path || "-")}\` | ${status} | ${replay.expectedStatus ?? "-"} | ${replay.actualStatus ?? "-"} | ${replay.durationMs ?? "-"} | ${mdCell(note)} |`,
      );
    }
    appendFailureTriage(lines, result);
    appendProfileAnalysis(lines, result);
  }

  return `${lines.join("\n")}\n`;
}

function appendModuleCoverageReport(lines, runContext, results, coverageBaseline) {
  const modules = buildModuleCoverage(runContext, results, coverageBaseline);
  lines.push("", "## 模块覆盖统计", "");
  lines.push("| 模块 | 页面访问 | API 用例 | GET | POST | PUT | PATCH | DELETE | 未覆盖接口 |");
  lines.push("| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |");
  for (const item of modules) {
    lines.push(
      `| ${mdCell(item.label)} | ${item.pagesVisited} | ${item.apiCases} | ${item.methods.GET || 0} | ${item.methods.POST || 0} | ${item.methods.PUT || 0} | ${item.methods.PATCH || 0} | ${item.methods.DELETE || 0} | ${item.uncoveredApiCount} |`,
    );
  }
}

function appendUncoveredReport(lines, runContext, results, coverageBaseline) {
  const uncoveredPages = getUncoveredPages(runContext, results);
  const coveredApiOperations = getCoveredApiOperations(results);
  const uncoveredApis = coverageBaseline.operations.filter(
    (operation) => !isOperationCovered(operation, coveredApiOperations),
  );
  const uncoveredWrites = uncoveredApis.filter((operation) =>
    ["POST", "PUT", "PATCH", "DELETE"].includes(operation.method),
  );
  const excludedPatterns = [
    ...(runContext.config.capture.ignoreUrlPatterns || []),
    ...(runContext.config.replay.ignoreUrlPatterns || []),
    ...(runContext.config.writeSafety?.blockedUrlPatterns || []),
  ].filter((item, index, all) => item && all.indexOf(item) === index);

  lines.push("", "## 未覆盖清单", "");

  lines.push("### 未覆盖页面", "");
  if (uncoveredPages.length === 0) {
    lines.push("- 无。");
  } else {
    lines.push("| 页面 | 模块 | 来源 |");
    lines.push("| --- | --- | --- |");
    for (const page of uncoveredPages.slice(0, 80)) {
      lines.push(`| \`${mdCell(page.path)}\` | ${mdCell(page.module)} | ${mdCell(page.source)} |`);
    }
    if (uncoveredPages.length > 80) {
      lines.push(`| ... | ... | 还有 ${uncoveredPages.length - 80} 个未展示 |`);
    }
  }

  lines.push("", "### 未覆盖接口", "");
  if (uncoveredApis.length === 0) {
    lines.push("- 无。");
  } else {
    lines.push("| 模块 | 方法 | 路径 |");
    lines.push("| --- | --- | --- |");
    for (const operation of uncoveredApis.slice(0, 120)) {
      lines.push(`| ${mdCell(operation.module)} | ${operation.method} | \`${mdCell(operation.path)}\` |`);
    }
    if (uncoveredApis.length > 120) {
      lines.push(`| ... | ... | 还有 ${uncoveredApis.length - 120} 个未展示 |`);
    }
  }

  lines.push("", "### 未覆盖写操作", "");
  if (uncoveredWrites.length === 0) {
    lines.push("- 无。");
  } else {
    lines.push("| 模块 | 方法 | 路径 | 建议 |");
    lines.push("| --- | --- | --- | --- |");
    for (const operation of uncoveredWrites.slice(0, 80)) {
      lines.push(`| ${mdCell(operation.module)} | ${operation.method} | \`${mdCell(operation.path)}\` | 建议检查 OpenAPI 自动用例是否缺少必要前置资源或字段模板 |`);
    }
    if (uncoveredWrites.length > 80) {
      lines.push(`| ... | ... | ... | 还有 ${uncoveredWrites.length - 80} 个未展示 |`);
    }
  }

  lines.push("", "### 主动排除功能（AI/上传）", "");
  if (excludedPatterns.length === 0) {
    lines.push("- 无。");
  } else {
    lines.push("| 排除规则 | 原因 |");
    lines.push("| --- | --- |");
    for (const pattern of excludedPatterns) {
      lines.push(`| \`${mdCell(pattern)}\` | 不触发 AI 评估、上传或外部模型调用 |`);
    }
  }
}

function appendDiscoveryReport(lines, result) {
  const exploration = result.rawCapture.exploration || createEmptyExplorationReport();
  const apiCoverage = summarizeApiCoverage(result.generatedCases, result.writeScenarioResults);
  const safetyBlocks = summarizeSafetyBlocks(result);

  lines.push("### 自动发现覆盖", "");
  lines.push(`- 发现路由: ${exploration.routesDiscovered}`);
  lines.push(`- 实际访问页面: ${exploration.pagesVisited}`);
  lines.push(`- 输入框: ${exploration.inputsFound} 个，已自动填写 ${exploration.inputsFilled} 个`);
  lines.push(`- 按钮/链接: ${exploration.buttonsFound} 个，已点击 ${exploration.buttonsClicked} 个`);
  lines.push(`- 自动取消确认弹窗: ${exploration.dialogsDismissed}`);
  lines.push(`- API 资源数: ${apiCoverage.resources.length}`);
  lines.push(`- 跳过/排除: ${safetyBlocks.total}`);

  if (apiCoverage.resources.length > 0) {
    lines.push("", "### API 资源覆盖", "");
    lines.push("| 资源 | 方法 | 用例数 |");
    lines.push("| --- | --- | ---: |");
    for (const resource of apiCoverage.resources) {
      lines.push(
        `| ${mdCell(resource.name)} | ${mdCell(resource.methods.join(", "))} | ${resource.count} |`,
      );
    }
  }

  if (exploration.skippedButtons.length > 0) {
    lines.push("", "### 被跳过的页面操作", "");
    lines.push("| 风险 | 文案 | 页面 |");
    lines.push("| --- | --- | --- |");
    for (const item of exploration.skippedButtons.slice(0, 30)) {
      lines.push(`| ${item.risk} | ${mdCell(item.text)} | \`${mdCell(item.route)}\` |`);
    }
  }
}

function appendProfileAnalysis(lines, result) {
  const exploration = result.rawCapture.exploration || createEmptyExplorationReport();
  const skippedByRisk = summarizeSkippedButtons(exploration.skippedButtons || []);
  const failedReplay = result.replayResults.filter((item) => !item.skipped && item.passed === false);
  const skippedReplay = result.replayResults.filter((item) => item.skipped);

  lines.push("", "### 结果分析", "");
  lines.push("| 指标 | 结果 | 说明 |");
  lines.push("| --- | ---: | --- |");
  lines.push(`| 页面访问数 | ${exploration.pagesVisited} | 自动探索实际打开的前端页面数量 |`);
  lines.push(`| 输入框填充 | ${exploration.inputsFilled}/${exploration.inputsFound} | 仅填写可见、可编辑且能识别字段含义的输入框 |`);
  lines.push(`| 按钮点击 | ${exploration.buttonsClicked}/${exploration.buttonsFound} | 按 clickRiskLevels 配置执行，测试环境可允许 danger/empty 按钮 |`);
  lines.push(`| API 资源数 | ${summarizeApiCoverage(result.generatedCases, result.writeScenarioResults).resources.length} | 按 /api/v1/{resource} 归类统计 |`);
  const caseSources = summarizeCaseSources(result.generatedCases);
  lines.push(`| 用例来源 | ${caseSources.capture || 0}/${caseSources.openapi || 0}/${caseSources.negative || 0} | 真实流量/OpenAPI 自动生成/负向异常用例 |`);
  const failureCategories = summarizeFailureCategories(result.replayResults);
  lines.push(`| 失败分层 | ${failureCategories.real || 0}/${failureCategories.dependency || 0}/${failureCategories.contract || 0}/${failureCategories.negative || 0} | 真实失败/依赖环境失败/契约不一致/负向用例失败 |`);
  lines.push(`| 回放失败 | ${failedReplay.length} | 状态码、响应结构或请求错误不符合预期 |`);
  lines.push(`| 回放跳过 | ${skippedReplay.length} | 多为自动用例缺少前置数据、权限上下文或触发限流，当前轮无法判断 |`);

  if (Object.keys(skippedByRisk).length > 0) {
    lines.push("", "### 跳过原因统计", "");
    lines.push("| 风险类型 | 数量 | 说明 |");
    lines.push("| --- | ---: | --- |");
    for (const [risk, count] of Object.entries(skippedByRisk)) {
      lines.push(`| ${mdCell(risk)} | ${count} | ${mdCell(describeRisk(risk))} |`);
    }
  }

  if (failedReplay.length > 0) {
    lines.push("", "### 失败接口复现命令", "");
    for (const replay of failedReplay.slice(0, 20)) {
      lines.push(`#### ${mdCell(replay.id)} ${mdCell(replay.method)} ${mdCell(replay.path)}`, "");
      lines.push("```bash");
      lines.push(replay.curlCommand || "# 未生成 curl 命令");
      lines.push("```", "");
    }
  }

  appendDebugArtifactReport(lines, result);

  lines.push("", "### 后续建议", "");
  if (failedReplay.length === 0 && (result.writeScenarioResults || []).every((item) => item.passed)) {
    lines.push("- 当前报告未发现失败接口，可作为本轮回归测试通过依据。");
  } else {
    lines.push("- 优先处理回放失败和写操作失败项，查看对应 JSON 明细定位请求与响应。");
  }
  if (skippedByRisk.empty) {
    lines.push("- 存在空文本图标按钮未执行；可继续给按钮补充 `aria-label`、`title` 或 `data-*` 来提升识别率。");
  }
}

function appendFailureTriage(lines, result) {
  const failed = result.replayResults.filter((item) => !item.skipped && item.passed === false);
  if (failed.length === 0) {
    return;
  }

  lines.push("", "### 失败分层明细", "");
  lines.push("| 分类 | 用例 | 方法 | 路径 | 实际状态 | 原因 |");
  lines.push("| --- | --- | --- | --- | ---: | --- |");
  for (const item of failed.slice(0, 40)) {
    const category = classifyReplayFailure(item);
    lines.push(
      `| ${mdCell(describeFailureCategory(category))} | ${mdCell(item.id)} | ${mdCell(item.method)} | \`${mdCell(item.path)}\` | ${item.actualStatus ?? "-"} | ${mdCell(explainReplayFailure(item, category))} |`,
    );
  }
}

function classifyReplayFailure(replay) {
  if (replay.negative) {
    return "negative";
  }
  const text = JSON.stringify(replay.responsePreview || replay.error || "");
  if (/connection refused|connect:|no such host|timeout|does not exist|relation .* does not exist|column .* does not exist/i.test(text)) {
    return "dependency";
  }
  if (Number(replay.actualStatus) >= 500 || Number(replay.actualStatus) === 0) {
    return "real";
  }
  return "contract";
}

function describeFailureCategory(category) {
  const labels = {
    real: "真实失败",
    dependency: "依赖/环境",
    negative: "负向失败",
    contract: "契约不一致",
  };
  return labels[category] || category;
}

function explainReplayFailure(replay, category) {
  if (category === "negative") {
    return "负向用例没有被明确拒绝，或触发了 5xx";
  }
  if (category === "dependency") {
    return "依赖服务、数据库结构或外部组件未就绪";
  }
  if (category === "contract") {
    return "OpenAPI 期望状态和实际业务返回不一致";
  }
  return replay.error || "接口返回 5xx 或请求执行失败";
}

function appendDebugArtifactReport(lines, result) {
  const exploration = result.rawCapture.exploration || createEmptyExplorationReport();
  const debugArtifacts = result.rawCapture.debugArtifacts || {
    consoleErrors: [],
    networkFailures: [],
    screenshots: [],
  };
  const failedRoutes = exploration.failedRoutes || [];
  const buttonFailures = exploration.buttonFailures || [];
  const hasDebugInfo =
    failedRoutes.length > 0 ||
    buttonFailures.length > 0 ||
    debugArtifacts.consoleErrors.length > 0 ||
    debugArtifacts.networkFailures.length > 0 ||
    debugArtifacts.screenshots.length > 0;
  if (!hasDebugInfo) {
    return;
  }

  lines.push("", "### 错误现场", "");
  lines.push("| 类型 | 数量 | 说明 |");
  lines.push("| --- | ---: | --- |");
  lines.push(`| 页面失败 | ${failedRoutes.length} | 自动探索打开页面失败时保存截图和错误 |`);
  lines.push(`| 按钮失败 | ${buttonFailures.length} | 页面按钮点击失败时记录 URL、错误，按截图策略可选保存图片 |`);
  lines.push(`| 控制台错误/警告 | ${debugArtifacts.consoleErrors.length} | 来源于 Playwright console 事件 |`);
  lines.push(`| 网络失败请求 | ${debugArtifacts.networkFailures.length} | 来源于 Playwright requestfailed 事件 |`);

  if (failedRoutes.length > 0 || buttonFailures.length > 0) {
    lines.push("", "| 类型 | URL/页面 | 错误 | 截图 |");
    lines.push("| --- | --- | --- | --- |");
    for (const item of failedRoutes.slice(0, 20)) {
      lines.push(`| 页面失败 | \`${mdCell(item.path)}\` | ${mdCell(item.error)} | ${item.screenshot ? `\`${mdCell(item.screenshot)}\`` : "-"} |`);
    }
    for (const item of buttonFailures.slice(0, 20)) {
      lines.push(`| 按钮失败 | \`${mdCell(item.route)}\` | ${mdCell(item.error)} | ${item.screenshot ? `\`${mdCell(item.screenshot)}\`` : "-"} |`);
    }
  }

  if (debugArtifacts.consoleErrors.length > 0) {
    lines.push("", "| 控制台类型 | 页面 | 内容 |");
    lines.push("| --- | --- | --- |");
    for (const item of summarizeConsoleMessages(debugArtifacts.consoleErrors).slice(0, 20)) {
      const suffix = item.count > 1 ? `（重复 ${item.count} 次）` : "";
      lines.push(`| ${mdCell(item.type)} | \`${mdCell(item.url)}\` | ${mdCell(item.text)}${suffix} |`);
    }
  }

  if (debugArtifacts.networkFailures.length > 0) {
    lines.push("", "| 方法 | URL | 失败原因 |");
    lines.push("| --- | --- | --- |");
    for (const item of debugArtifacts.networkFailures.slice(0, 20)) {
      lines.push(`| ${mdCell(item.method)} | \`${mdCell(item.url)}\` | ${mdCell(item.failure)} |`);
    }
  }
}

function summarizeConsoleMessages(messages) {
  const grouped = new Map();
  for (const message of messages || []) {
    const normalizedText = String(message.text || "").replace(/\s+/g, " ").slice(0, 220);
    const key = `${message.type}|${message.url}|${normalizedText}`;
    if (!grouped.has(key)) {
      grouped.set(key, { ...message, text: normalizedText, count: 0 });
    }
    grouped.get(key).count += 1;
  }
  return [...grouped.values()];
}

function mdCell(value) {
  return String(value ?? "")
    .replace(/\|/g, "\\|")
    .replace(/\r?\n/g, " ")
    .replace(/\s+/g, " ")
    .trim();
}

function describeCaseSource(source) {
  const labels = {
    capture: "真实流量",
    openapi: "OpenAPI",
    negative: "负向",
  };
  return labels[source] || source;
}

function summarizeApiCoverage(generatedCases, writeScenarioResults = []) {
  const resources = new Map();
  const addResource = (pathValue, methodValue) => {
    const resource = parseApiResource(pathValue);
    if (!resource) {
      return;
    }
    if (!resources.has(resource.name)) {
      resources.set(resource.name, { name: resource.name, methods: new Set(), count: 0 });
    }
    const item = resources.get(resource.name);
    item.methods.add(String(methodValue).toUpperCase());
    item.count += 1;
  };

  for (const testCase of generatedCases || []) {
    addResource(testCase.request.normalizedPath, testCase.request.method);
  }

  for (const scenario of writeScenarioResults || []) {
    for (const step of scenario.steps || []) {
      addResource(step.path, step.method);
    }
  }

  return {
    resources: [...resources.values()].map((item) => ({
      ...item,
      methods: [...item.methods].sort(),
    })),
  };
}

async function loadCoverageBaseline(runContext) {
  const spec = await loadOpenApiSpec(runContext);
  return {
    source: spec.source,
    operations: spec.operations,
  };
}

async function loadOpenApiSpec(runContext) {
  const candidates = runContext.config.coverage?.openApiCandidates || [];
  for (const candidate of candidates) {
    const filePath = path.resolve(runContext.bigtestRoot, candidate);
    const operations = filePath.endsWith(".json")
      ? await loadOpenApiJsonOperations(filePath, runContext)
      : await loadOpenApiYamlOperations(filePath, runContext);
    if (operations.length > 0) {
      return { source: filePath, operations };
    }
  }
  return { source: "", operations: [] };
}

async function loadOpenApiJsonOperations(filePath, runContext) {
  if (!(await pathExists(filePath))) {
    return [];
  }
  try {
    const doc = await readJson(filePath);
    return openApiPathsToOperations(doc.paths || {}, runContext);
  } catch (error) {
    logWarn(`读取 OpenAPI JSON 失败 ${filePath}: ${error.message}`);
    return [];
  }
}

async function loadOpenApiYamlOperations(filePath, runContext) {
  const source = await readTextIfExists(filePath);
  if (!source) {
    return [];
  }

  const paths = {};
  let currentPath = "";
  let inPaths = false;
  for (const line of source.split(/\r?\n/)) {
    if (/^paths:\s*$/.test(line)) {
      inPaths = true;
      continue;
    }
    if (inPaths && /^[A-Za-z0-9_-]+:\s*$/.test(line)) {
      break;
    }
    const pathMatch = line.match(/^ {2}(\/[^:]+):\s*$/);
    if (pathMatch) {
      currentPath = pathMatch[1];
      paths[currentPath] = {};
      continue;
    }
    const methodMatch = line.match(/^ {4}(get|post|put|patch|delete):\s*$/i);
    if (currentPath && methodMatch) {
      paths[currentPath][methodMatch[1].toLowerCase()] = {};
    }
  }
  return openApiPathsToOperations(paths, runContext);
}

function openApiPathsToOperations(paths, runContext) {
  const methods = new Set(["get", "post", "put", "patch", "delete"]);
  const ignorePatterns = [
    ...(runContext.config.capture.ignoreUrlPatterns || []),
    ...(runContext.config.replay.ignoreUrlPatterns || []),
  ];
  const operations = [];
  for (const [pathValue, pathItem] of Object.entries(paths || {})) {
    if (!pathValue.startsWith("/api/") || ignorePatterns.some((pattern) => pathValue.includes(pattern))) {
      continue;
    }
    if (pathValue.includes("/ws")) {
      continue;
    }
    for (const method of Object.keys(pathItem || {})) {
      if (!methods.has(method.toLowerCase())) {
        continue;
      }
      const operation = pathItem[method] || {};
      operations.push({
        method: method.toUpperCase(),
        path: pathValue,
        summary: operation.summary || "",
        description: operation.description || "",
        parameters: [
          ...((pathItem.parameters || []).filter(Boolean)),
          ...((operation.parameters || []).filter(Boolean)),
        ],
        requiresAuth: Boolean(operation.security?.length),
        expectedStatuses: extractExpectedStatuses(operation.responses),
        module: getModuleLabelForApi(pathValue, runContext.config.coverage?.modules || {}),
      });
    }
  }
  const methodOrder = { GET: 1, POST: 2, PUT: 3, PATCH: 4, DELETE: 5 };
  return operations.sort((a, b) => {
    const byModulePath = `${a.module}${a.path}`.localeCompare(`${b.module}${b.path}`);
    return byModulePath || (methodOrder[a.method] || 9) - (methodOrder[b.method] || 9);
  });
}

function extractExpectedStatuses(responses) {
  const statuses = Object.keys(responses || {})
    .map((status) => Number(status))
    .filter((status) => Number.isInteger(status) && status >= 200 && status < 300);
  return statuses.length > 0 ? statuses : [200, 201, 204];
}

function buildModuleCoverage(runContext, results, coverageBaseline) {
  const modules = new Map();
  const ensureModule = (label) => {
    if (!modules.has(label)) {
      modules.set(label, {
        label,
        pagesVisitedSet: new Set(),
        apiCases: 0,
        methods: { GET: 0, POST: 0, PUT: 0, PATCH: 0, DELETE: 0 },
        uncoveredApiCount: 0,
      });
    }
    return modules.get(label);
  };

  for (const result of results) {
    const exploration = result.rawCapture.exploration || createEmptyExplorationReport();
    for (const routePath of exploration.visitedRoutes || []) {
      ensureModule(getModuleLabelForPage(routePath, runContext.config.coverage?.modules || {}))
        .pagesVisitedSet.add(routePath);
    }
    for (const operation of getCoveredApiOperations([result])) {
      const module = ensureModule(getModuleLabelForApi(operation.path, runContext.config.coverage?.modules || {}));
      module.apiCases += 1;
      module.methods[operation.method] = (module.methods[operation.method] || 0) + 1;
    }
  }

  const coveredApiOperations = getCoveredApiOperations(results);
  for (const operation of coverageBaseline.operations || []) {
    const module = ensureModule(operation.module);
    if (!isOperationCovered(operation, coveredApiOperations)) {
      module.uncoveredApiCount += 1;
    }
  }

  return [...modules.values()]
    .map((item) => ({
      ...item,
      pagesVisited: item.pagesVisitedSet.size,
    }))
    .sort((a, b) => a.label.localeCompare(b.label, "zh-CN"));
}

function getCoveredApiOperations(results) {
  const operations = [];
  for (const result of results || []) {
    for (const testCase of result.generatedCases || []) {
      operations.push({
        method: String(testCase.request.method || "").toUpperCase(),
        path: testCase.request.normalizedPath || "",
      });
    }
    for (const scenario of result.writeScenarioResults || []) {
      for (const step of scenario.steps || []) {
        operations.push({
          method: String(step.method || "").toUpperCase(),
          path: step.path || "",
        });
      }
    }
  }
  return operations;
}

function isOperationCovered(operation, coveredOperations) {
  return coveredOperations.some((covered) =>
    covered.method === operation.method && pathMatchesTemplate(operation.path, covered.path),
  );
}

function pathMatchesTemplate(template, actualPath) {
  const actual = String(actualPath || "").split("?")[0];
  const escaped = String(template || "")
    .replace(/[.*+?^${}()|[\]\\]/g, "\\$&")
    .replace(/\\\{[^}]+\\\}/g, "[^/]+")
    .replace(/:([A-Za-z0-9_]+)/g, "[^/]+");
  return new RegExp(`^${escaped}$`).test(actual);
}

function getUncoveredPages(runContext, results) {
  const discovered = new Map();
  const visited = new Set();
  for (const result of results || []) {
    const exploration = result.rawCapture.exploration || createEmptyExplorationReport();
    for (const routePath of exploration.visitedRoutes || []) {
      visited.add(routePath);
    }
    for (const [source, routes] of Object.entries({
      router: exploration.routesFromFiles || [],
      pages: exploration.routesFromPageDirs || [],
      dom: exploration.routesFromDom || [],
    })) {
      for (const routePath of routes) {
        if (!discovered.has(routePath)) {
          discovered.set(routePath, {
            path: routePath,
            source,
            module: getModuleLabelForPage(routePath, runContext.config.coverage?.modules || {}),
          });
        }
      }
    }
  }
  return [...discovered.values()].filter((item) => !visited.has(item.path));
}

function getModuleLabelForPage(routePath, moduleMap) {
  const normalized = String(routePath || "").toLowerCase();
  const hit = Object.entries(moduleMap || {}).find(([key]) =>
    normalized.includes(`/${key}`) || normalized.includes(key.replace(/s$/, "")),
  );
  return hit?.[1] || "其他";
}

function getModuleLabelForApi(pathValue, moduleMap) {
  const resource = parseApiResource(pathValue);
  if (!resource) {
    return "其他";
  }
  return moduleMap?.[resource.name] || moduleMap?.[resource.name.replace(/s$/, "")] || resource.name;
}

function summarizeRunMetrics(results) {
  return results.reduce(
    (summary, result) => {
      const exploration = result.rawCapture.exploration || createEmptyExplorationReport();
      summary.pagesVisited += exploration.pagesVisited;
      summary.rawCaptures += result.rawCapture.captures.length;
      summary.generatedCases += result.generatedCases.length;
      summary.writePass += (result.writeScenarioResults || []).filter((item) => item.passed).length;
      summary.writeFail += (result.writeScenarioResults || []).filter((item) => !item.passed).length;
      summary.replayPass += result.replayResults.filter((item) => item.passed).length;
      summary.replayFail += result.replayResults.filter((item) => !item.skipped && item.passed === false).length;
      summary.replaySkip += result.replayResults.filter((item) => item.skipped).length;
      const caseSources = summarizeCaseSources(result.generatedCases);
      for (const [source, count] of Object.entries(caseSources)) {
        summary.caseSources[source] = (summary.caseSources[source] || 0) + count;
      }
      const failureCategories = summarizeFailureCategories(result.replayResults);
      for (const [category, count] of Object.entries(failureCategories)) {
        summary.failureCategories[category] = (summary.failureCategories[category] || 0) + count;
      }
      return summary;
    },
    {
      pagesVisited: 0,
      rawCaptures: 0,
      generatedCases: 0,
      writePass: 0,
      writeFail: 0,
      replayPass: 0,
      replayFail: 0,
      replaySkip: 0,
      caseSources: {},
      failureCategories: {},
    },
  );
}

function summarizeFailureCategories(replayResults) {
  const summary = {};
  for (const replay of replayResults || []) {
    if (replay.skipped || replay.passed !== false) {
      continue;
    }
    const category = classifyReplayFailure(replay);
    summary[category] = (summary[category] || 0) + 1;
  }
  return summary;
}

function summarizeCaseSources(generatedCases) {
  const summary = {};
  for (const testCase of generatedCases || []) {
    const source = testCase.source || "capture";
    summary[source] = (summary[source] || 0) + 1;
  }
  return summary;
}

function summarizeSkippedButtons(skippedButtons) {
  const summary = {};
  for (const item of skippedButtons || []) {
    summary[item.risk || "unknown"] = (summary[item.risk || "unknown"] || 0) + 1;
  }
  return summary;
}

function countWriteSteps(writeScenarioResults, method) {
  const targetMethod = String(method).toUpperCase();
  return (writeScenarioResults || []).reduce((count, scenario) => {
    return count + (scenario.steps || []).filter((step) => String(step.method).toUpperCase() === targetMethod).length;
  }, 0);
}

function formatWriteStepName(name, method) {
  const normalizedName = String(name || "");
  const normalizedMethod = String(method || "").toUpperCase();
  if (normalizedName === "delete" || normalizedMethod === "DELETE") {
    return "delete 删除测试数据";
  }
  const labels = {
    create: "create 创建测试数据",
    verify: "verify 验证测试数据",
    update: "update 更新测试数据",
  };
  return labels[normalizedName] || normalizedName;
}

function describeRisk(risk) {
  const descriptions = {
    danger: "高风险页面操作，如删除、上传、AI评估、真实投递等",
    empty: "按钮没有可识别文本，通常是图标按钮",
    unknown: "无法根据文案归类的按钮或链接",
    medium: "中风险操作，如新增、编辑、保存、提交",
    safe: "低风险操作，如查询、查看、返回、重置",
  };
  return descriptions[risk] || "未分类操作";
}

function summarizeSafetyBlocks(result) {
  const skippedReplay = (result.replayResults || []).filter((item) => item.skipped).length;
  const skippedButtons = result.rawCapture.exploration?.skippedButtons?.length || 0;
  return {
    total: skippedReplay + skippedButtons,
    skippedReplay,
    skippedButtons,
  };
}

export function printConsoleSummary(runContext, results, artifactResult) {
  console.log("");
  console.log("=".repeat(72));
  console.log("bigtest 真实流量自动化测试完成");
  console.log("=".repeat(72));
  console.log(`Run ID: ${runContext.runId}`);
  console.log(`输出目录: ${artifactResult.reportFile.replace(/\/report\.md$/, "")}`);

  for (const result of results) {
    const exploration = result.rawCapture.exploration || createEmptyExplorationReport();
    const writePass = (result.writeScenarioResults || []).filter((item) => item.passed).length;
    const writeFail = (result.writeScenarioResults || []).filter((item) => !item.passed).length;
    const replayPass = result.replayResults.filter((item) => item.passed).length;
    const replayFail = result.replayResults.filter(
      (item) => !item.skipped && item.passed === false,
    ).length;
    const replaySkip = result.replayResults.filter((item) => item.skipped).length;

    console.log(
      `- ${result.profileName}: 页面 ${exploration.pagesVisited}, 原始请求 ${result.rawCapture.captures.length}, 生成用例 ${result.generatedCases.length}, 写操作通过 ${writePass}, 写操作失败 ${writeFail}, 回放通过 ${replayPass}, 失败 ${replayFail}, 跳过 ${replaySkip}`,
    );
  }

  console.log(`报告文件: ${artifactResult.reportFile}`);
  console.log("=".repeat(72));
}

export function hasReplayFailures(results) {
  return results.some((result) =>
    result.replayResults.some((item) => !item.skipped && item.passed === false) ||
    (result.writeScenarioResults || []).some((item) => item.passed === false),
  );
}
