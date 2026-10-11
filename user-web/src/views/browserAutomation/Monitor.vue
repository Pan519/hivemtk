<template>
  <div class="page">
    <div class="header">
      <h2>执行监控 #{{ sessionId }}</h2>
      <el-space>
        <el-tag v-if="session" :type="{ completed: 'success', failed: 'danger', stopped: 'info', active: 'warning' }[session.status] || 'info'">
          {{ session.status }}
        </el-tag>
        <el-tag v-if="session?.confirm_pending" type="warning">{{ $t('auto.Monitor.3f7a87f1') }}</el-tag>
        <el-button v-if="session && ['created','active'].includes(session.status)" type="danger" @click="onStop">{{ $t('auto.Monitor.095e938e') }}</el-button>
        <el-button v-if="session?.confirm_pending" type="primary" :loading="confirming" :disabled="!gate?.payload_hash" @click="onConfirm">{{ $t('auto.Monitor.fb155863') }}</el-button>
        <el-button v-if="session" @click="onExport">{{ $t('auto.Monitor.9f7b3162') }}</el-button>
      </el-space>
    </div>

    <el-alert v-if="loadError" type="error" :closable="false" show-icon style="margin-top: 12px">
      <template #title>
        {{ session ? '最新状态读取失败，下方是刷新前的残留内容' : '这条会话读不到，本页没有在监控任何东西' }}：{{ loadError }}
      </template>
      <el-space style="margin-top: 6px">
        <el-button size="small" type="primary" @click="retryLoad">{{ $t('auto.Monitor.132c5cdc') }}</el-button>
        <el-button size="small" @click="goTasks">{{ $t('auto.Monitor.90e22f5e') }}</el-button>
      </el-space>
    </el-alert>

    <el-card v-if="gate" :header="$t('auto.Monitor.97a8546c')" style="margin-top: 12px">
      <el-descriptions :column="2" border>
        <el-descriptions-item :label="$t('auto.Monitor.a3e36110')">第 {{ gate?.step_index }} 步</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.Monitor.df826a24')">{{ formatTime(gate?.expires_at) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.Monitor.509feaf0')" :span="2">
          <span style="white-space: pre-wrap">{{ gate?.preview || '（无预览）' }}</span>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card v-if="session" style="margin-top: 12px">
      <el-descriptions :column="4" border>
        <el-descriptions-item :label="$t('auto.Monitor.39f1374d')">{{ sessionDurationText(session) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.Monitor.b664352f')">{{ session.total_steps ? `${session.success_steps}/${session.total_steps}` : '—' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.Monitor.592c5958')">{{ session.started_at ? new Date(session.started_at).toLocaleString('zh-CN') : '—' }}</el-descriptions-item>
        <el-descriptions-item v-if="session.error_msg" :label="$t('auto.Monitor.7030ff64')" :span="4">
          <span style="color: #ef4444">{{ session.error_msg }}</span>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card v-if="session?.llm_summary" :header="$t('auto.Monitor.3cab3678')" style="margin-top: 12px">
      {{ session.llm_summary }}
    </el-card>

    <el-card v-if="hasExtracted" :header="$t('auto.Monitor.4ac114e2')" style="margin-top: 12px">
      <div v-for="item in dataEntries(session.extracted_data)" :key="item.key" class="detail-item">
        <span class="detail-key">{{ item.key }}</span>
        <pre v-if="item.kind === 'block'" class="detail-block">{{ item.text }}</pre>
        <span v-else class="detail-inline">{{ item.text }}</span>
      </div>
    </el-card>

    <el-card v-if="session" :header="$t('auto.Monitor.2b6fe614')" style="margin-top: 12px">
      <el-table :data="steps" v-loading="loading">
        <el-table-column type="expand">
          <template #default="{ row }">
            <div class="step-detail">
              <div v-if="row.submit_state" class="submit-ledger">
                <span class="detail-key">{{ $t('auto.Monitor.7da5f255') }}</span>
                <el-tag size="small" :type="SUBMIT_COPY[row.submit_state]?.type || 'info'">{{ row.submit_state }}</el-tag>
                <span class="detail-inline">{{ SUBMIT_COPY[row.submit_state]?.text || '未知提交态（服务端未给结论，按未提交对待）' }}</span>
                <span v-if="row.text_hash" class="detail-hash">正文指纹 {{ row.text_hash }}</span>
              </div>
              <div v-for="g in stepDataGroups(row)" :key="g.title" class="detail-group">
                <div class="detail-title">{{ g.title }}</div>
                <div v-for="item in g.entries" :key="item.key" class="detail-item">
                  <span class="detail-key">{{ item.key }}</span>
                  <pre v-if="item.kind === 'block'" class="detail-block">{{ item.text }}</pre>
                  <span v-else class="detail-inline">{{ item.text }}</span>
                </div>
              </div>
              <div v-if="!row.submit_state && !stepDataGroups(row).length" class="detail-empty">
                这一步没有回包数据（动作只改变页面状态，不返回内容）
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column type="index" label="#" width="60" />
        <el-table-column prop="action" :label="$t('auto.Monitor.c500cfab')" width="130">
          <template #default="{ row }">
            <span>{{ row.action }}</span>
            <el-tag v-if="row.is_write" size="small" type="warning" class="write-tag">{{ $t('auto.Monitor.4d7dc6c5') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="target" :label="$t('auto.Monitor.73e82552')" min-width="160" show-overflow-tooltip />
        <el-table-column :label="$t('auto.Monitor.fe7509e0')" min-width="140" show-overflow-tooltip>
          <template #default="{ row }">{{ row.value || '—' }}</template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('auto.Monitor.3fea7ca7')" width="110">
          <template #default="{ row }">
            <el-tag :type="{ success: 'success', failed: 'danger', running: 'warning', skipped: 'info', pending: 'info' }[row.status] || 'info'">
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="duration_ms" :label="$t('auto.Monitor.39f1374d')" width="100">
          <template #default="{ row }">{{ durationText(row.duration_ms) }}</template>
        </el-table-column>
        <el-table-column prop="error_msg" :label="$t('auto.Monitor.7030ff64')" min-width="200" show-overflow-tooltip />
      </el-table>
    </el-card>

    <el-card v-if="session" style="margin-top: 12px">
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <span>{{ $t('auto.Monitor.c5d8a5e6') }}</span>
          <el-select v-model="logDirection" size="small" style="width: 110px" @change="loadLogs">
            <el-option :label="$t('auto.Monitor.a8b0c204')" value="" />
            <el-option label="command" value="command" />
            <el-option label="event" value="event" />
            <el-option label="judge" value="judge" />
          </el-select>
        </div>
      </template>
      <el-table :data="logs" size="small" max-height="420">
        <el-table-column prop="seq" label="seq" width="60" />
        <el-table-column prop="direction" :label="$t('auto.Monitor.a465db53')" width="90">
          <template #default="{ row }">
            <el-tag size="small" :type="{ command: '', event: 'success', judge: 'warning' }[row.direction] || 'info'">{{ row.direction }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="action" :label="$t('auto.Monitor.c500cfab')" width="130" />
        <el-table-column prop="ok" :label="$t('auto.Monitor.dd4d3660')" width="80">
          <template #default="{ row }">{{ outcomeMark(row.ok) }}</template>
        </el-table-column>
        <el-table-column prop="duration_ms" :label="$t('auto.Monitor.39f1374d')" width="80">
          <template #default="{ row }">{{ durationText(row.duration_ms) }}</template>
        </el-table-column>
        <el-table-column label="payload" min-width="260" show-overflow-tooltip>
          <template #default="{ row }"><span style="font-family: monospace; font-size: 12px">{{ payloadPreview(row.payload) }}</span></template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card v-if="session?.snapshot" style="margin-top: 12px">
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center">
          <span>页面快照（可访问性树文本，{{ session.snapshot.length }} 字符——定位失败时拿它对照真实结构）</span>
          <el-button size="small" @click="snapshotOpen = !snapshotOpen">{{ snapshotOpen ? '收起' : '展开全部' }}</el-button>
        </div>
      </template>
      <pre class="detail-block" :class="{ collapsed: !snapshotOpen }">{{ session.snapshot }}</pre>
    </el-card>

    <el-card v-if="session?.final_screenshot_url" :header="$t('auto.Monitor.c92a834b')" style="margin-top: 12px">
      <el-image
        :src="session.final_screenshot_url"
        fit="contain"
        style="max-width: 100%"
        :preview-src-list="[session.final_screenshot_url]"
      >
        <template #error>
          <div class="img-fallback">
            截图取不到（文件可能不在本机 uploads 目录里，或服务端存的是相对路径而非可访问 URL）：{{ session.final_screenshot_url }}
          </div>
        </template>
      </el-image>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getBrowserSession, getBrowserSessionSteps, stopBrowserSession, confirmBrowserSession, getBrowserConfirmGate, getBrowserSessionLogs, exportBrowserSessionAudit, interpretConfirmResult } from '@/api/browserAutomation'
import { durationText, sessionDurationText } from './durationText'

const route = useRoute()
const router = useRouter()
const goTasks = () => router.push('/browser-automation/tasks')
const sessionId = computed(() => route.params.id)
const { t } = useI18n()

const session = ref(null)
const gate = ref(null)
const steps = ref([])
const logs = ref([])
const logDirection = ref('')
const loading = ref(false)
const confirming = ref(false)

// 四个态的文案必须各说各的话：把「闸门挂在别的进程」说成「没有待确认提交点」，
// 用户会去改编排，而该改的是实例数（服务端 A10 分流的就是这一条）。
const CONFIRM_COPY = {
  'granted': { type: 'success', text: t('auto.Monitor.d6f8677f') },
  'no_gate': { type: 'warning', text: t('auto.Monitor.be333951') },
  'payload_mismatch': { type: 'error', text: t('auto.Monitor.40aad76e') },
  'gate_on_another_instance': { type: 'error', text: t('auto.Monitor.7561dc99') },
  'unknown': { type: 'error', text: t('auto.Monitor.0c94a2f1') },
}

const formatTime = (v) => {
  if (!v) return '—'
  const d = new Date(v)
  return Number.isNaN(d.getTime()) ? String(v) : d.toLocaleString('zh-CN')
}

const payloadPreview = (p) => {
  if (p == null) return '—'
  const str = typeof p === 'string' ? p : JSON.stringify(p)
  return str.length > 160 ? str.slice(0, 160) + '…' : str
}

// ok 是三态，不是布尔。null = 这一帧不携带结论（下发帧写下时 Host 还没回执；
// 之后更进一步——那一帧甚至可能根本没上线）。把 null 折进 falsy 分支就等于把每一条
// 「下发」都画成失败，而旧模板 `row.ok ? '✓' : '✗'` 在服务端改掉常量 true 的同一批就会犯这个错。
// 未知形状（缺字段/字符串）一律走 '—'：宁可说"这帧没结论"，也不替用户编一个红或绿。
const outcomeMark = (ok) => (ok === true ? '✓' : ok === false ? '✗' : '—')

// 写台账四态的文案。步的 status 记「这一步跑成什么样」，submit_state 记「这条内容的提交
// 是否可能发生」，两者在「send 到达但回查未见」时必然分叉（status=failed 而提交已生效）。
// 所以这一行不许并进状态标签里：把 failed 直接读成「没发出去」就会去重发，而重发＝双发。
const SUBMIT_COPY = {
  prepared: { type: 'info', text: t('auto.Monitor.5bacaac9') },
  sent: { type: 'warning', text: t('auto.Monitor.3c40c951') },
  verified: { type: 'success', text: t('auto.Monitor.4d0e462f') },
  unattributed: { type: 'danger', text: t('auto.Monitor.36b7ee82') },
}

const hasExtracted = computed(() => {
  const d = session.value?.extracted_data
  return d && typeof d === 'object' && Object.keys(d).length > 0
})

// 快照/正文类长字段在 a11y 树里动辄数千字符，默认折起来，展开才给全高。
const snapshotOpen = ref(false)

// 把 {k: v} 摊成 [{key, kind, text}]：kind=block 走等宽 pre（长文本/数组/对象），
// kind=inline 走一行短文本。extract 的回包形态（键=选择器名、值=字符串或字符串数组）
// 与 query/markdown 的标量形态共用这一份渲染，避免出现第二个「JSON 怎么显示」的口径。
const describeValue = (val) => {
  if (Array.isArray(val)) {
    return { kind: 'block', text: `${val.length} 项\n${val.map((v, i) => `${i + 1}. ${typeof v === 'string' ? v : JSON.stringify(v)}`).join('\n')}` }
  }
  if (val && typeof val === 'object') return { kind: 'block', text: JSON.stringify(val, null, 2) }
  if (typeof val === 'string') {
    return val.length > 60 ? { kind: 'block', text: val } : { kind: 'inline', text: val === '' ? '（空字符串）' : val }
  }
  if (val == null) return { kind: 'inline', text: 'null' }
  return { kind: 'inline', text: String(val) }
}

const dataEntries = (obj) => {
  if (obj == null) return []
  const src = typeof obj === 'string' ? (() => { try { return JSON.parse(obj) } catch { return null } })() : obj
  if (src == null || typeof src !== 'object') return [{ key: t('auto.Monitor.fe7509e0'), ...describeValue(obj) }]
  return Object.entries(src).map(([key, val]) => ({ key, ...describeValue(val) }))
}

// 步骤展开面板里的数据组：params 是编排声明的扩展参数（步「想怎么做」），
// result 是原语回包（步「实际拿到什么」）。两者都要能被看见——只给状态行不给回包，
// extract/snapshot/query 这类「产出物就是数据」的步骤等于什么都没执行。
const stepDataGroups = (row) => {
  const groups = []
  const params = dataEntries(row.params)
  if (params.length) groups.push({ title: t('auto.Monitor.1631942f'), entries: params })
  const result = dataEntries(row.result)
  if (result.length) groups.push({ title: t('auto.Monitor.b8b78ed8'), entries: result })
  return groups
}

async function loadLogs() {
  try {
    const res = await getBrowserSessionLogs(sessionId.value, logDirection.value)
    const list = unpack(res)
    logs.value = Array.isArray(list) ? list : list?.list || []
  } catch {
    logs.value = []
  }
}
let pollTimer = null

const unpack = (res) => res?.data ?? res

// 会话读不到时这一页本来会永远空转：轮询的终止条件写在 session 上，
// 而 session 正因为读不到才是空的 —— 于是每 2s 两次 404、每轮一条未捕获异常，
// 页面空白，用户看不出它到底在监控什么。404 的会话不会自己变回来；
// 网络抖光是暂时的，所以只连续失败若干轮才收，收掉之后还给用户一个「重试」。
const POLL_FAILURE_LIMIT = 3
const loadError = ref('')
let pollFailures = 0

const stopPolling = () => { if (pollTimer) { clearInterval(pollTimer); pollTimer = null } }
const startPolling = () => { if (!pollTimer) pollTimer = setInterval(load, 2000) }

function retryLoad() {
  pollFailures = 0
  loadError.value = ''
  startPolling()
  load()
}

async function load() {
  loading.value = steps.value.length === 0
  try {
    const [sRes, stRes] = await Promise.all([
      getBrowserSession(sessionId.value),
      getBrowserSessionSteps(sessionId.value),
    ])
    loadLogs()
    session.value = unpack(sRes)
    const list = unpack(stRes)
    steps.value = Array.isArray(list) ? list : list?.list || []
    loadGate()
    pollFailures = 0
    loadError.value = ''
    // 终态停止轮询
    if (session.value && ['completed', 'failed', 'stopped'].includes(session.value.status)) stopPolling()
  } catch (e) {
    // 拦截器已经把服务端原因弹成 toast（同文案 2.5s 内去重），这里只补两件它不做的：
    // 把「本页当前读不到」写在页面上，以及让空转停下来。异常不再往外抛，
    // 否则每个轮询周期都在控制台留一条未捕获 rejection。
    loadError.value = String(e?.message || e)
    pollFailures += 1
    if (e?.status === 404 || pollFailures >= POLL_FAILURE_LIMIT) stopPolling()
  } finally {
    loading.value = false
  }
}

// 闸门详情：只在挂起期间取，且取失败一律置空——没有载荷就没有放行入口可点（fail-closed）。
// 与 load() 并发而不串行：会话详情是主数据，闸门读侧慢/失败都不该把步流水一起拖住。
async function loadGate() {
  if (!session.value?.confirm_pending) {
    if (gate.value) gate.value = null
    return
  }
  if (gateLoading) return
  gateLoading = true
  try {
    const res = unpack(await getBrowserConfirmGate(sessionId.value))
    gate.value = res?.pending ? res.gate : null
  } catch {
    gate.value = null
  } finally {
    gateLoading = false
  }
}
let gateLoading = false

async function onStop() {
  try {
    const res = await stopBrowserSession(sessionId.value, t('auto.Monitor.484b13f7'))
    const data = unpack(res)
    if (data?.stopped) {
      ElMessage.success(t('auto.Monitor.8d1bcd4b'))
    } else {
      // 服务端对同一条命令回答的是「没有在跑的执行协程」（200 + stopped:false）。
      // 把它播报成「请求已发送」，用户就会等一次永远不会来的收口——而这条会话仍按
      // created/active 占着并发闸（CountRunningByUser），下一次下发会被判「已有任务执行中」。
      // 既然没有协程在跑，轮询也等不来任何变化，一并停掉。
      stopPolling()
      ElMessage.warning(t('auto.Monitor.6d363ebf'))
    }
  } catch {
    // 拦截器已经把服务端原因弹成 toast；这里接住它，只为了让停止失败不留下未捕获异常
  }
  load()
}

// D7：人工放行挂起的写操作提交点。放行绑载荷——只能批准页面正在显示的那一份；
// 结论按 status 分流（confirmed 只说「闸门认没认」，三种「没放行」必须各说各的原因）。
async function onConfirm() {
  const hash = gate.value?.payload_hash
  if (!hash) {
    ElMessage.warning(t('auto.Monitor.68e8b920'))
    load()
    return
  }
  confirming.value = true
  try {
    let raw
    try {
      raw = await confirmBrowserSession(sessionId.value, hash)
    } catch (err) {
      raw = err // 409 的结论同样写在响应体里，与 200 走同一个判读函数
    }
    const verdict = interpretConfirmResult(raw)
    const copy = CONFIRM_COPY[verdict.status] || CONFIRM_COPY.unknown
    if (copy.type === 'success') ElMessage.success(copy.text)
    else if (copy.type === 'warning') ElMessage.warning(copy.text)
    else ElMessage.error(copy.text)
    // 载荷不符时服务端把当前挂起的闸门带回，立刻换掉预览：用户下一步是照新内容重批，不是刷新碰运气
    if (verdict.gate) gate.value = verdict.gate
    load()
  } finally {
    confirming.value = false
  }
}

// I5：审计包导出（会话+步+命令流+LLM 成本账）落盘 JSON
async function onExport() {
  try {
    const res = await exportBrowserSessionAudit(sessionId.value)
    const payload = unpack(res)
    const blob = new Blob([JSON.stringify(payload, null, 2)], { type: 'application/json' })
    const a = document.createElement('a')
    a.href = URL.createObjectURL(blob)
    a.download = `browser_session_${sessionId.value}_audit.json`
    a.click()
    URL.revokeObjectURL(a.href)
    ElMessage.success(t('auto.Monitor.2f3c0f75'))
  } catch {
    ElMessage.error(t('auto.Monitor.dd51ab50'))
  }
}

onMounted(() => {
  load()
  startPolling()
})
onBeforeUnmount(stopPolling)
</script>

<style scoped>
.page { padding: 16px; }
.header { display: flex; justify-content: space-between; align-items: center; }
.step-detail { padding: 8px 16px; }
.detail-group { margin-bottom: 12px; }
.detail-title { font-size: 12px; color: #909399; margin-bottom: 6px; }
.detail-item { display: flex; gap: 8px; align-items: flex-start; margin-bottom: 6px; }
.detail-key { min-width: 110px; color: #606266; font-size: 12px; }
.detail-inline { font-size: 13px; }
.detail-hash { color: #909399; font-size: 12px; margin-left: 8px; }
.detail-block {
  margin: 0;
  flex: 1;
  max-height: 220px;
  overflow: auto;
  background: #f6f7f9;
  border-radius: 4px;
  padding: 8px;
  font-size: 12px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
}
.detail-block.collapsed { max-height: 120px; }
.submit-ledger { display: flex; align-items: center; gap: 8px; margin-bottom: 10px; flex-wrap: wrap; }
.write-tag { margin-left: 6px; }
.detail-empty { color: #909399; font-size: 12px; }
.img-fallback { display: flex; align-items: center; justify-content: center; height: 100%; min-height: 120px; color: #909399; font-size: 12px; padding: 12px; }
</style>
