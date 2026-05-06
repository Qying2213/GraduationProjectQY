import fs from "node:fs/promises";
import path from "node:path";

export const SENSITIVE_KEYS = [
  "authorization",
  "cookie",
  "set-cookie",
  "token",
  "access_token",
  "refresh_token",
  "password",
];

export function logInfo(message) {
  console.log(`[INFO] ${message}`);
}

export function logWarn(message) {
  console.warn(`[WARN] ${message}`);
}

export function logError(message) {
  console.error(`[ERROR] ${message}`);
}

export function timestampId() {
  const now = new Date();
  const pad = (value) => String(value).padStart(2, "0");
  return [
    now.getFullYear(),
    pad(now.getMonth() + 1),
    pad(now.getDate()),
  ].join("-") + "_" + [
    pad(now.getHours()),
    pad(now.getMinutes()),
    pad(now.getSeconds()),
  ].join("-");
}

export function isPlainObject(value) {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

export function deepMerge(base, override) {
  if (Array.isArray(base) || Array.isArray(override)) {
    return override ?? base;
  }
  if (!isPlainObject(base) || !isPlainObject(override)) {
    return override ?? base;
  }

  const output = { ...base };
  for (const key of Object.keys(override)) {
    output[key] = key in base ? deepMerge(base[key], override[key]) : override[key];
  }
  return output;
}

export function containsSensitiveKey(key) {
  const lower = String(key).toLowerCase();
  return SENSITIVE_KEYS.some((item) => lower.includes(item));
}

export function sanitizeValue(value) {
  if (Array.isArray(value)) {
    return value.map((item) => sanitizeValue(item));
  }
  if (isPlainObject(value)) {
    const result = {};
    for (const [key, item] of Object.entries(value)) {
      result[key] = containsSensitiveKey(key) ? "<redacted>" : sanitizeValue(item);
    }
    return result;
  }
  return value;
}

export function sanitizeHeaders(headers = {}) {
  const result = {};
  for (const [key, value] of Object.entries(headers)) {
    result[key] = containsSensitiveKey(key) ? "<redacted>" : value;
  }
  return result;
}

export function tryParseJson(text) {
  if (!text || typeof text !== "string") {
    return null;
  }
  try {
    return JSON.parse(text);
  } catch {
    return null;
  }
}

export function jsonStringify(value) {
  return JSON.stringify(value, null, 2);
}

export function summarizeJsonShape(value) {
  if (Array.isArray(value)) {
    const summary = { kind: "array", length: value.length };
    const first = value[0];
    if (isPlainObject(first)) {
      summary.sampleKeys = Object.keys(first).sort();
    } else if (first !== undefined) {
      summary.sampleType = typeof first;
    }
    return summary;
  }
  if (isPlainObject(value)) {
    return {
      kind: "object",
      keys: Object.keys(value).sort(),
    };
  }
  return {
    kind: typeof value,
  };
}

export function areShapesCompatible(expected, actual) {
  if (!expected || !actual) {
    return true;
  }
  if (expected.kind !== actual.kind) {
    return false;
  }
  if (expected.kind === "object") {
    const actualKeys = new Set(actual.keys || []);
    return (expected.keys || []).every((key) => actualKeys.has(key));
  }
  if (expected.kind === "array" && expected.sampleKeys?.length) {
    const actualKeys = new Set(actual.sampleKeys || []);
    return expected.sampleKeys.every((key) => actualKeys.has(key));
  }
  return true;
}

export async function ensureDir(dirPath) {
  await fs.mkdir(dirPath, { recursive: true });
}

export async function pathExists(filePath) {
  try {
    await fs.access(filePath);
    return true;
  } catch {
    return false;
  }
}

export async function listFilesRecursive(dirPath, options = {}) {
  const {
    maxDepth = 4,
    includeExtensions = [],
    ignoreDirs = ["node_modules", "dist", "build", ".git", "output"],
  } = options;
  const output = [];

  async function walk(currentDir, depth) {
    if (depth > maxDepth || !(await pathExists(currentDir))) {
      return;
    }

    const entries = await fs.readdir(currentDir, { withFileTypes: true });
    for (const entry of entries) {
      const entryPath = path.join(currentDir, entry.name);
      if (entry.isDirectory()) {
        if (!ignoreDirs.includes(entry.name)) {
          await walk(entryPath, depth + 1);
        }
        continue;
      }

      if (
        includeExtensions.length === 0 ||
        includeExtensions.some((extension) => entry.name.endsWith(extension))
      ) {
        output.push(entryPath);
      }
    }
  }

  await walk(dirPath, 0);
  return output;
}

export function normalizeUrlForFingerprint(rawUrl, ignoreQueryKeys = []) {
  const url = new URL(rawUrl);
  const params = new URLSearchParams();
  const ignored = new Set(ignoreQueryKeys.map((item) => String(item)));
  const entries = [...url.searchParams.entries()]
    .filter(([key, value]) => !ignored.has(key) && value !== "")
    .sort(([aKey, aValue], [bKey, bValue]) =>
      aKey === bKey ? aValue.localeCompare(bValue) : aKey.localeCompare(bKey),
    );

  for (const [key, value] of entries) {
    params.append(key, value);
  }

  const search = params.toString();
  return `${url.pathname}${search ? `?${search}` : ""}`;
}

export function buildRequestFingerprint(entry, ignoreQueryKeys = []) {
  const normalizedPath = normalizeUrlForFingerprint(entry.url, ignoreQueryKeys);
  const requiresAuth = entry.request.requiresAuth ? "auth" : "anon";
  const bodySignature = entry.request.bodyShape
    ? jsonStringify(entry.request.bodyShape)
    : entry.request.postData
      ? typeof entry.request.postData
      : "";

  return [
    entry.method.toUpperCase(),
    normalizedPath,
    requiresAuth,
    bodySignature,
  ].join("::");
}

export function bodyShapeFromValue(value) {
  if (value === null || value === undefined) {
    return null;
  }
  if (typeof value === "string") {
    const parsed = tryParseJson(value);
    return parsed ? summarizeJsonShape(parsed) : { kind: "string" };
  }
  return summarizeJsonShape(value);
}

export async function readJson(filePath) {
  const raw = await fs.readFile(filePath, "utf-8");
  return JSON.parse(raw);
}

export async function readText(filePath) {
  return fs.readFile(filePath, "utf-8");
}

export async function readTextIfExists(filePath) {
  if (!(await pathExists(filePath))) {
    return null;
  }
  return readText(filePath);
}

export async function writeJson(filePath, value) {
  await fs.writeFile(filePath, jsonStringify(value), "utf-8");
}

export async function writeText(filePath, value) {
  await fs.writeFile(filePath, value, "utf-8");
}

export function relativeFrom(rootPath, targetPath) {
  return path.relative(rootPath, targetPath) || ".";
}

export function getByPath(value, dottedPath) {
  if (!dottedPath) {
    return value;
  }
  return dottedPath.split(".").reduce((current, key) => {
    if (current === null || current === undefined) {
      return undefined;
    }
    return current[key];
  }, value);
}

export function interpolate(value, variables) {
  if (typeof value === "string") {
    return value.replace(/\$\{([^}]+)\}/g, (_, key) => {
      const resolved = getByPath(variables, key.trim());
      return resolved === undefined || resolved === null ? "" : String(resolved);
    });
  }
  if (Array.isArray(value)) {
    return value.map((item) => interpolate(item, variables));
  }
  if (isPlainObject(value)) {
    const result = {};
    for (const [key, item] of Object.entries(value)) {
      result[key] = interpolate(item, variables);
    }
    return result;
  }
  return value;
}
