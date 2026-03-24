const http = require('http')
const fs = require('fs')
const path = require('path')
const net = require('net')
const { URL } = require('url')

const rootDir = path.resolve(__dirname, '..', '..')
const manifestPath = process.env.DEV_DASHBOARD_MANIFEST || path.join(rootDir, 'logs', 'dev', 'services.json')
const dashboardPort = Number(process.env.DEV_DASHBOARD_PORT || 8091)
const staticDir = __dirname

function sendSSE(res, eventName, payload) {
  res.write(`event: ${eventName}\n`)
  res.write(`data: ${JSON.stringify(payload)}\n\n`)
}

function sendJson(res, statusCode, payload) {
  res.writeHead(statusCode, {
    'Content-Type': 'application/json; charset=utf-8',
    'Cache-Control': 'no-store'
  })
  res.end(JSON.stringify(payload))
}

function sendFile(res, filePath, contentType) {
  fs.readFile(filePath, (error, content) => {
    if (error) {
      sendJson(res, 404, { error: 'File not found' })
      return
    }

    res.writeHead(200, {
      'Content-Type': contentType,
      'Cache-Control': 'no-store'
    })
    res.end(content)
  })
}

function normalizeLogText(text) {
  return String(text || '').replace(/\0/g, '').replace(/\r\n/g, '\n')
}

async function loadManifest() {
  try {
    const raw = await fs.promises.readFile(manifestPath, 'utf8')
    const services = JSON.parse(raw)
    return Array.isArray(services) ? services : []
  } catch (error) {
    return []
  }
}

async function findService(serviceName) {
  const services = await loadManifest()
  return services.find((service) => service.name === serviceName) || null
}

async function readPid(pidFile) {
  try {
    const raw = await fs.promises.readFile(pidFile, 'utf8')
    const pid = Number.parseInt(raw.trim(), 10)
    return Number.isFinite(pid) ? pid : null
  } catch (error) {
    return null
  }
}

function isPidRunning(pid) {
  if (!pid) {
    return false
  }

  try {
    process.kill(pid, 0)
    return true
  } catch (error) {
    return false
  }
}

function isPortListening(port) {
  return new Promise((resolve) => {
    const socket = net.createConnection({ host: '127.0.0.1', port })
    let settled = false

    const finish = (value) => {
      if (settled) {
        return
      }
      settled = true
      socket.destroy()
      resolve(value)
    }

    socket.setTimeout(400)
    socket.once('connect', () => finish(true))
    socket.once('timeout', () => finish(false))
    socket.once('error', () => finish(false))
  })
}

async function readTail(logFile, maxLines = 18) {
  if (!Number.isFinite(maxLines) || maxLines <= 0) {
    return ''
  }

  try {
    const stats = await fs.promises.stat(logFile)
    const bytesToRead = Math.min(stats.size, 64 * 1024)
    const start = Math.max(0, stats.size - bytesToRead)
    const handle = await fs.promises.open(logFile, 'r')
    const buffer = Buffer.alloc(bytesToRead)
    await handle.read(buffer, 0, bytesToRead, start)
    await handle.close()

    const text = normalizeLogText(buffer.toString('utf8'))
    const lines = text.split(/\r?\n/)
    return lines.slice(-maxLines).join('\n').trim()
  } catch (error) {
    return ''
  }
}

async function readSlice(logFile, start, length) {
  if (length <= 0) {
    return ''
  }

  try {
    const handle = await fs.promises.open(logFile, 'r')
    const buffer = Buffer.alloc(length)
    const { bytesRead } = await handle.read(buffer, 0, length, start)
    await handle.close()
    return normalizeLogText(buffer.slice(0, bytesRead).toString('utf8'))
  } catch (error) {
    return ''
  }
}

async function buildServiceState(service, previewLines) {
  const pid = await readPid(service.pidFile)
  const pidRunning = isPidRunning(pid)
  const listening = await isPortListening(service.port)
  const preview = await readTail(service.logFile, previewLines)

  let status = 'stopped'
  if (listening) {
    status = 'running'
  } else if (pidRunning) {
    status = 'starting'
  }

  return {
    ...service,
    pid,
    pidRunning,
    listening,
    status,
    logPreview: preview,
    logPath: path.relative(rootDir, service.logFile)
  }
}

async function buildOverview(previewLines) {
  const services = await loadManifest()
  const items = await Promise.all(services.map((service) => buildServiceState(service, previewLines)))
  const running = items.filter((item) => item.status === 'running').length
  const starting = items.filter((item) => item.status === 'starting').length

  return {
    generatedAt: new Date().toISOString(),
    dashboardPort,
    total: items.length,
    running,
    starting,
    stopped: items.length - running - starting,
    services: items
  }
}

async function streamServiceLogs(req, res, service, initialLines) {
  res.writeHead(200, {
    'Content-Type': 'text/event-stream; charset=utf-8',
    'Cache-Control': 'no-cache, no-transform',
    Connection: 'keep-alive',
    'X-Accel-Buffering': 'no'
  })

  if (typeof res.flushHeaders === 'function') {
    res.flushHeaders()
  }

  let closed = false
  let currentSize = 0
  let currentStatus = ''
  let running = false
  let pollTimer = null
  let heartbeatTimer = null

  const cleanup = () => {
    if (closed) {
      return
    }

    closed = true
    if (pollTimer) {
      clearInterval(pollTimer)
    }
    if (heartbeatTimer) {
      clearInterval(heartbeatTimer)
    }
  }

  req.on('close', cleanup)
  req.on('aborted', cleanup)

  const emitStatus = async (force = false) => {
    const pid = await readPid(service.pidFile)
    const listening = await isPortListening(service.port)
    const status = listening ? 'running' : isPidRunning(pid) ? 'starting' : 'stopped'
    const payload = {
      name: service.name,
      pid,
      listening,
      status
    }
    const nextStatus = JSON.stringify(payload)

    if (force || currentStatus !== nextStatus) {
      currentStatus = nextStatus
      sendSSE(res, 'status', payload)
    }
  }

  const emitSnapshot = async () => {
    const content = await readTail(service.logFile, initialLines)

    try {
      const stats = await fs.promises.stat(service.logFile)
      currentSize = stats.size
    } catch (error) {
      currentSize = 0
    }

    await emitStatus(true)
    sendSSE(res, 'snapshot', {
      name: service.name,
      content
    })
  }

  await emitSnapshot()

  pollTimer = setInterval(async () => {
    if (closed || running) {
      return
    }

    running = true
    try {
      await emitStatus(false)

      let stats
      try {
        stats = await fs.promises.stat(service.logFile)
      } catch (error) {
        if (currentSize !== 0) {
          currentSize = 0
          sendSSE(res, 'snapshot', {
            name: service.name,
            content: ''
          })
        }
        return
      }

      if (stats.size < currentSize) {
        currentSize = stats.size
        sendSSE(res, 'snapshot', {
          name: service.name,
          content: await readTail(service.logFile, initialLines)
        })
        return
      }

      const diff = stats.size - currentSize
      if (diff <= 0) {
        return
      }

      if (diff > 128 * 1024) {
        currentSize = stats.size
        sendSSE(res, 'snapshot', {
          name: service.name,
          content: await readTail(service.logFile, initialLines)
        })
        return
      }

      const chunk = await readSlice(service.logFile, currentSize, diff)
      currentSize = stats.size

      if (chunk) {
        sendSSE(res, 'chunk', {
          name: service.name,
          content: chunk
        })
      }
    } finally {
      running = false
    }
  }, 1000)

  heartbeatTimer = setInterval(() => {
    if (!closed) {
      res.write(': ping\n\n')
    }
  }, 15000)
}

async function handleApi(req, res, parsedUrl) {
  if (parsedUrl.pathname === '/api/overview') {
    const lines = Number.parseInt(parsedUrl.searchParams.get('lines') || '18', 10)
    const overview = await buildOverview(Number.isFinite(lines) ? lines : 18)
    sendJson(res, 200, overview)
    return
  }

  if (parsedUrl.pathname === '/api/services') {
    const overview = await buildOverview(0)
    sendJson(res, 200, overview)
    return
  }

  if (parsedUrl.pathname.startsWith('/api/logs/') && parsedUrl.pathname.endsWith('/stream')) {
    const serviceName = decodeURIComponent(parsedUrl.pathname.slice('/api/logs/'.length, -'/stream'.length))
    const target = await findService(serviceName)

    if (!target) {
      sendJson(res, 404, { error: 'Service not found' })
      return
    }

    const lines = Number.parseInt(parsedUrl.searchParams.get('lines') || '120', 10)
    await streamServiceLogs(req, res, target, Number.isFinite(lines) ? lines : 120)
    return
  }

  if (parsedUrl.pathname.startsWith('/api/logs/')) {
    const serviceName = decodeURIComponent(parsedUrl.pathname.replace('/api/logs/', ''))
    const target = await findService(serviceName)

    if (!target) {
      sendJson(res, 404, { error: 'Service not found' })
      return
    }

    const lines = Number.parseInt(parsedUrl.searchParams.get('lines') || '120', 10)
    const content = await readTail(target.logFile, Number.isFinite(lines) ? lines : 120)
    const pid = await readPid(target.pidFile)
    const listening = await isPortListening(target.port)

    sendJson(res, 200, {
      ...target,
      pid,
      listening,
      status: listening ? 'running' : isPidRunning(pid) ? 'starting' : 'stopped',
      content
    })
    return
  }

  if (parsedUrl.pathname === '/healthz') {
    sendJson(res, 200, { ok: true, port: dashboardPort })
    return
  }

  sendJson(res, 404, { error: 'Not found' })
}

const server = http.createServer(async (req, res) => {
  const parsedUrl = new URL(req.url, `http://${req.headers.host || `127.0.0.1:${dashboardPort}`}`)

  if (parsedUrl.pathname.startsWith('/api/')) {
    await handleApi(req, res, parsedUrl)
    return
  }

  if (parsedUrl.pathname === '/' || parsedUrl.pathname === '/index.html') {
    sendFile(res, path.join(staticDir, 'index.html'), 'text/html; charset=utf-8')
    return
  }

  if (parsedUrl.pathname === '/app.js') {
    sendFile(res, path.join(staticDir, 'app.js'), 'application/javascript; charset=utf-8')
    return
  }

  if (parsedUrl.pathname === '/styles.css') {
    sendFile(res, path.join(staticDir, 'styles.css'), 'text/css; charset=utf-8')
    return
  }

  sendJson(res, 404, { error: 'Not found' })
})

server.listen(dashboardPort, '0.0.0.0', () => {
  console.log(`[dev-dashboard] listening on :${dashboardPort}`)
  console.log(`[dev-dashboard] manifest: ${manifestPath}`)
})
