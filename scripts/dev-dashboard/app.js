const state = {
  services: [],
  serviceFingerprint: '',
  logBuffers: new Map(),
  eventSources: new Map(),
  reconnectTimers: new Map(),
  streamStates: new Map(),
  autoScroll: true,
  overviewTimer: null
}

const elements = {
  totalCount: document.getElementById('totalCount'),
  runningCount: document.getElementById('runningCount'),
  startingCount: document.getElementById('startingCount'),
  stoppedCount: document.getElementById('stoppedCount'),
  dashboardPort: document.getElementById('dashboardPort'),
  lastUpdated: document.getElementById('lastUpdated'),
  streamStatus: document.getElementById('streamStatus'),
  serviceWallboard: document.getElementById('serviceWallboard'),
  refreshButton: document.getElementById('refreshButton'),
  reconnectButton: document.getElementById('reconnectButton'),
  autoScrollToggle: document.getElementById('autoScrollToggle')
}

function buildDashboardUrl(path) {
  return new URL(path, window.location.href).toString()
}

function escapeHtml(value) {
  return String(value || '')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

function serviceKey(name) {
  return String(name || '').replace(/[^a-zA-Z0-9_-]/g, '_')
}

function formatStatus(status) {
  if (status === 'running') {
    return '运行中'
  }
  if (status === 'starting') {
    return '启动中'
  }
  return '已停止'
}

function formatTime(isoString) {
  if (!isoString) {
    return '--'
  }

  return new Date(isoString).toLocaleString('zh-CN', {
    hour12: false
  })
}

function normalizeLogText(text) {
  return String(text || '').replace(/\0/g, '').replace(/\r\n/g, '\n')
}

function clampLogLines(text, maxLines = 700) {
  const lines = normalizeLogText(text).split('\n')
  if (lines.length <= maxLines) {
    return lines.join('\n').trim()
  }
  return lines.slice(-maxLines).join('\n').trim()
}

function classifyLogLine(line) {
  const value = String(line || '')
  const lower = value.toLowerCase()

  if (
    /\b(error|fatal|panic|exception|traceback|refused|timeout|timed out|failed|failure|denied|unavailable)\b/i.test(value) ||
    /\b5\d{2}\b/.test(lower)
  ) {
    return 'error'
  }

  if (/\b(warn|warning|deprecated)\b/i.test(value)) {
    return 'warn'
  }

  if (/\b(success|listening|running|started|ready|healthy)\b/i.test(value)) {
    return 'success'
  }

  return 'normal'
}

function renderHighlightedLog(text) {
  if (!text) {
    return '<span class="log-line log-empty">等待日志输出...</span>'
  }

  return normalizeLogText(text)
    .split('\n')
    .map((line) => {
      const variant = classifyLogLine(line)
      const content = line ? escapeHtml(line) : '&nbsp;'
      return `<span class="log-line log-${variant}">${content}</span>`
    })
    .join('')
}

function panelRefs(serviceName) {
  const key = serviceKey(serviceName)
  return {
    status: document.getElementById(`panel-status-${key}`),
    pid: document.getElementById(`panel-pid-${key}`),
    stream: document.getElementById(`panel-stream-${key}`),
    log: document.getElementById(`panel-log-${key}`),
    updated: document.getElementById(`panel-updated-${key}`)
  }
}

function setPanelLog(serviceName, text) {
  const refs = panelRefs(serviceName)
  if (!refs.log) {
    return
  }

  const normalized = clampLogLines(text)
  state.logBuffers.set(serviceName, normalized)
  refs.log.innerHTML = renderHighlightedLog(normalized)

  if (state.autoScroll) {
    refs.log.scrollTop = refs.log.scrollHeight
  }
}

function appendPanelLog(serviceName, chunk) {
  const current = state.logBuffers.get(serviceName) || ''
  const normalizedChunk = normalizeLogText(chunk).trim()
  if (!normalizedChunk) {
    return
  }

  const next = current ? `${current}\n${normalizedChunk}` : normalizedChunk
  setPanelLog(serviceName, next)
}

function setStreamState(serviceName, label, variant) {
  state.streamStates.set(serviceName, variant)

  const refs = panelRefs(serviceName)
  if (refs.stream) {
    refs.stream.className = `stream-pill stream-${variant}`
    refs.stream.textContent = label
  }

  updateGlobalStreamStatus()
}

function updateGlobalStreamStatus() {
  const total = state.services.length
  const connected = Array.from(state.streamStates.values()).filter((value) => value === 'live').length
  const reconnecting = Array.from(state.streamStates.values()).filter((value) => value === 'reconnecting').length

  if (!total) {
    elements.streamStatus.textContent = '实时日志：暂无服务'
    return
  }

  if (reconnecting > 0) {
    elements.streamStatus.textContent = `实时日志：${connected}/${total} 已连接，${reconnecting} 个重连中`
    return
  }

  elements.streamStatus.textContent = `实时日志：${connected}/${total} 已连接`
}

function updatePanelMeta(service) {
  const refs = panelRefs(service.name)
  if (!refs.status || !refs.pid) {
    return
  }

  refs.status.className = `status-pill status-${service.status}`
  refs.status.textContent = formatStatus(service.status)
  refs.pid.textContent = service.pid ? `PID ${service.pid}` : 'PID --'

  if (refs.updated) {
    refs.updated.textContent = `状态更新时间：${formatTime(new Date().toISOString())}`
  }
}

function renderWallboard() {
  if (!state.services.length) {
    elements.serviceWallboard.innerHTML = `
      <article class="log-panel empty-panel">
        <h2>还没有服务清单</h2>
        <p>请先执行 <code>./start-all.sh</code>，面板会自动开始监听所有日志。</p>
      </article>
    `
    return
  }

  elements.serviceWallboard.innerHTML = state.services
    .map((service) => {
      const key = serviceKey(service.name)
      const openLink = service.url
        ? `<a class="panel-link" href="${escapeHtml(service.url)}" target="_blank" rel="noreferrer">打开服务</a>`
        : '<span class="panel-link disabled-link">无服务地址</span>'

      return `
        <article class="log-panel">
          <div class="panel-head">
            <div>
              <p class="service-name">${escapeHtml(service.name)}</p>
              <p class="service-port">监听端口 :${escapeHtml(service.port)}</p>
            </div>
            <div class="panel-badges">
              <span id="panel-status-${key}" class="status-pill status-${escapeHtml(service.status)}">${formatStatus(service.status)}</span>
              <span id="panel-stream-${key}" class="stream-pill stream-connecting">日志流连接中</span>
            </div>
          </div>
          <div class="panel-meta">
            <span id="panel-pid-${key}">${service.pid ? `PID ${service.pid}` : 'PID --'}</span>
            <span>${escapeHtml(service.logPath)}</span>
          </div>
          <pre id="panel-log-${key}" class="live-log">等待日志流连接...</pre>
          <div class="panel-foot">
            <span id="panel-updated-${key}" class="foot-note">状态更新时间：--</span>
            <div class="foot-actions">
              ${openLink}
              <button class="ghost-button" type="button" data-clear-log="${escapeHtml(service.name)}">清空当前视图</button>
            </div>
          </div>
        </article>
      `
    })
    .join('')

  elements.serviceWallboard.querySelectorAll('[data-clear-log]').forEach((button) => {
    button.addEventListener('click', () => {
      const serviceName = button.getAttribute('data-clear-log')
      if (!serviceName) {
        return
      }

      setPanelLog(serviceName, '')
    })
  })

  state.services.forEach((service) => {
    updatePanelMeta(service)
    setPanelLog(service.name, state.logBuffers.get(service.name) || service.logPreview || '')
  })
}

async function fetchOverview(forceRender = false) {
  const response = await fetch(buildDashboardUrl('./api/overview?lines=0'), { cache: 'no-store' })
  if (!response.ok) {
    throw new Error(`Failed to load overview: ${response.status}`)
  }

  const data = await response.json()
  const services = data.services || []
  const nextFingerprint = services.map((service) => service.name).join('|')
  const shouldRender = forceRender || nextFingerprint !== state.serviceFingerprint

  state.services = services
  state.serviceFingerprint = nextFingerprint

  elements.totalCount.textContent = data.total || 0
  elements.runningCount.textContent = data.running || 0
  elements.startingCount.textContent = data.starting || 0
  elements.stoppedCount.textContent = data.stopped || 0
  elements.dashboardPort.textContent = `开发面板端口：${data.dashboardPort}`
  elements.lastUpdated.textContent = `最后刷新：${formatTime(data.generatedAt)}`

  if (shouldRender) {
    renderWallboard()
    reconnectAllStreams()
    return
  }

  state.services.forEach((service) => {
    updatePanelMeta(service)
  })
}

function closeStream(serviceName) {
  const existing = state.eventSources.get(serviceName)
  if (existing) {
    existing.close()
    state.eventSources.delete(serviceName)
  }

  const reconnectTimer = state.reconnectTimers.get(serviceName)
  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
    state.reconnectTimers.delete(serviceName)
  }
}

function scheduleReconnect(service, delay = 1800) {
  closeStream(service.name)
  setStreamState(service.name, '日志流重连中', 'reconnecting')

  const timer = window.setTimeout(() => {
    state.reconnectTimers.delete(service.name)
    connectStream(service)
  }, delay)

  state.reconnectTimers.set(service.name, timer)
}

function connectStream(service) {
  closeStream(service.name)
  setStreamState(service.name, '日志流连接中', 'connecting')

  const source = new EventSource(buildDashboardUrl(`./api/logs/${encodeURIComponent(service.name)}/stream?lines=120`))
  state.eventSources.set(service.name, source)

  source.onopen = () => {
    setStreamState(service.name, '实时流已连接', 'live')
  }

  source.addEventListener('snapshot', (event) => {
    try {
      const payload = JSON.parse(event.data)
      setPanelLog(service.name, payload.content || '')
    } catch (error) {
      console.error(`解析 ${service.name} 日志快照失败`, error)
    }
  })

  source.addEventListener('chunk', (event) => {
    try {
      const payload = JSON.parse(event.data)
      appendPanelLog(service.name, payload.content || '')
    } catch (error) {
      console.error(`解析 ${service.name} 日志增量失败`, error)
    }
  })

  source.addEventListener('status', (event) => {
    try {
      const payload = JSON.parse(event.data)
      const target = state.services.find((item) => item.name === service.name)
      if (target) {
        target.pid = payload.pid
        target.status = payload.status
      }
      updatePanelMeta({
        ...service,
        pid: payload.pid,
        status: payload.status
      })
    } catch (error) {
      console.error(`解析 ${service.name} 状态失败`, error)
    }
  })

  source.onerror = () => {
    if (state.eventSources.get(service.name) !== source) {
      return
    }

    scheduleReconnect(service)
  }
}

function reconnectAllStreams() {
  state.services.forEach((service) => {
    connectStream(service)
  })
}

function closeAllStreams() {
  Array.from(state.eventSources.keys()).forEach((serviceName) => {
    closeStream(serviceName)
  })
}

async function refreshAll(forceRender = false) {
  try {
    await fetchOverview(forceRender)
  } catch (error) {
    elements.lastUpdated.textContent = `刷新失败：${error.message}`
  }
}

function startOverviewPolling() {
  if (state.overviewTimer) {
    clearInterval(state.overviewTimer)
  }

  state.overviewTimer = window.setInterval(() => {
    refreshAll(false)
  }, 2500)
}

elements.refreshButton.addEventListener('click', () => {
  refreshAll(false)
})

elements.reconnectButton.addEventListener('click', () => {
  reconnectAllStreams()
})

elements.autoScrollToggle.addEventListener('change', (event) => {
  state.autoScroll = event.target.checked

  if (state.autoScroll) {
    state.services.forEach((service) => {
      const refs = panelRefs(service.name)
      if (refs.log) {
        refs.log.scrollTop = refs.log.scrollHeight
      }
    })
  }
})

window.addEventListener('beforeunload', closeAllStreams)

refreshAll(true)
startOverviewPolling()
