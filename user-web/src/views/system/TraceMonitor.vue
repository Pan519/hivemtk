<template>
  <div class="app-container trace-monitor">
    <el-tabs v-model="activeTab" class="trace-tabs">
      
      <el-tab-pane :label="$t('auto.TraceMonitor.1.33668d')" name="trace">
        <el-form :inline="true" class="trace-query" @submit.prevent>
          <el-form-item label="Trace ID">
            <el-input v-model.trim="form.trace_id" :placeholder="$t('auto.TraceMonitor.2.7eac93')" clearable style="width: 260px" />
          </el-form-item>
          <el-form-item :label="$t('auto.TraceMonitor.3.427bc5')">
            <el-input v-model.trim="form.conversation_id" :placeholder="$t('auto.TraceMonitor.4.8b6554')" clearable style="width: 260px" />
          </el-form-item>
          <el-form-item :label="$t('auto.TraceMonitor.5.795e00')">
            <el-input v-model.trim="form.msg_id" :placeholder="$t('auto.TraceMonitor.6.bccd98')" clearable style="width: 220px" />
          </el-form-item>
          <el-form-item :label="$t('auto.TraceMonitor.7.ebee05')">
            <el-select v-model="form.channel" :placeholder="$t('auto.TraceMonitor.8.a8b0c2')" clearable style="width: 160px">
              <el-option v-for="c in channels" :key="c" :label="c" :value="c" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="loadingTree" @click="loadTree">{{ $t('auto.TraceMonitor.62.bee912') }}</el-button>
            <el-button @click="resetForm">{{ $t('auto.TraceMonitor.63.4b9c32') }}</el-button>
          </el-form-item>
        </el-form>

        <el-alert
          v-if="treeError"
          type="error"
          :closable="false"
          :title="treeError"
          show-icon
          class="mb-12"
        />

        <div v-else-if="tree" class="trace-summary mb-12">
          <el-tag type="info" effect="plain">Trace: {{ tree.trace_id || '—' }}</el-tag>
          <el-tag type="info" effect="plain">会话: {{ tree.conversation_id || '—' }}</el-tag>
          <el-tag type="info" effect="plain">渠道: {{ tree.channel || '—' }}</el-tag>
          <span class="trace-count">共 {{ tree.spans ? tree.spans.length : 0 }} 个 Span</span>
        </div>

        <el-empty v-if="!loadingTree && !tree && !treeError" :description="$t('auto.TraceMonitor.9.e4417d')" />
        <el-skeleton v-else-if="loadingTree" :rows="8" animated />

        <div v-else class="trace-tree">
          <template v-for="(item, idx) in flattened" :key="idx">
            
            <el-card
              v-if="!item._turn && !item._tool"
              class="span-card"
              :class="{ abnormal: isAbnormal(item), agent: item.node === 'ai_dispatch' }"
              shadow="hover"
            >
              <div class="span-head" role="button" tabindex="0" :aria-expanded="!!expanded[item.id]" @click="toggle(item)" @keydown.enter.prevent="toggle(item)" @keydown.space.prevent="toggle(item)">
                <span class="caret">{{ expanded[item.id] ? '▾' : '▸' }}</span>
                <span class="span-kind kind-lifecycle">{{ $t('auto.TraceMonitor.64.176808') }}</span>
                <span class="span-title">{{ nodeLabel(item.node) }}</span>
                <span class="span-node">{{ item.node }}</span>
                <el-tag size="small" :type="statusType(item)">{{ statusText(item) }}</el-tag>
                <span class="span-dur">{{ item.duration_ms }} ms</span>
                <span class="span-time">{{ item.created_at }}</span>
              </div>
              <div v-show="expanded[item.id]" class="span-body">
                <el-descriptions :column="1" border size="small">
                  <el-descriptions-item :label="$t('auto.TraceMonitor.10.c6b523')">
                    <pre class="code-block">{{ fmt(item.input) }}</pre>
                  </el-descriptions-item>
                  <el-descriptions-item :label="$t('auto.TraceMonitor.11.bd30be')">
                    <pre class="code-block">{{ fmt(item.output) }}</pre>
                  </el-descriptions-item>
                  <el-descriptions-item :label="$t('auto.TraceMonitor.12.384c09')">{{ item.expected || '—' }}</el-descriptions-item>
                  <el-descriptions-item v-if="item.abnormal || item.error" :label="$t('auto.TraceMonitor.13.c195df')">
                    <span class="err-text">{{ item.error || '异常' }}</span>
                  </el-descriptions-item>
                </el-descriptions>
              </div>
            </el-card>

            
            <el-card v-else-if="item._turn" class="span-card turn-card" shadow="hover">
              <div class="span-head" role="button" tabindex="0" :aria-expanded="!!expanded[item.turn.id]" @click="toggle(item.turn)" @keydown.enter.prevent="toggle(item.turn)" @keydown.space.prevent="toggle(item.turn)">
                <span class="caret">{{ expanded[item.turn.id] ? '▾' : '▸' }}</span>
                <span class="span-kind kind-turn">{{ $t('auto.TraceMonitor.65.35dcb8') }}</span>
                <span class="span-title">第 {{ item.ti }} 轮</span>
                <span v-if="item.turn.agent_id" class="span-node">{{ item.turn.agent_id }}</span>
                <el-tag size="small" :type="statusType(item.turn)">{{ statusText(item.turn) }}</el-tag>
                <span class="span-dur">{{ item.turn.duration_ms }} ms</span>
                <span class="span-time">{{ item.turn.created_at }}</span>
              </div>
              <div v-show="expanded[item.turn.id]" class="span-body">
                <el-descriptions :column="1" border size="small">
                  <el-descriptions-item :label="$t('auto.TraceMonitor.10.c6b523')">
                    <pre class="code-block">{{ fmt(item.turn.input) }}</pre>
                  </el-descriptions-item>
                  <el-descriptions-item :label="$t('auto.TraceMonitor.11.bd30be')">
                    <pre class="code-block">{{ fmt(item.turn.output) }}</pre>
                  </el-descriptions-item>
                  <el-descriptions-item v-if="item.turn.abnormal || item.turn.error" :label="$t('auto.TraceMonitor.13.c195df')">
                    <span class="err-text">{{ item.turn.error || '异常' }}</span>
                  </el-descriptions-item>
                </el-descriptions>
              </div>
            </el-card>

            
            <el-card v-else class="span-card tool-card" shadow="never">
              <div class="span-head">
                <span class="span-kind kind-tool">{{ $t('auto.TraceMonitor.66.20dce2') }}</span>
                <span class="span-title">{{ item.tool.tool_name }}</span>
                <el-tag size="small" :type="statusType(item.tool)">{{ statusText(item.tool) }}</el-tag>
                <span class="span-dur">{{ item.tool.duration_ms }} ms</span>
                <span class="span-time">{{ item.tool.created_at }}</span>
              </div>
              <div class="span-body">
                <el-descriptions :column="1" border size="small">
                  <el-descriptions-item :label="$t('auto.TraceMonitor.10.c6b523')">
                    <pre class="code-block">{{ fmt(item.tool.input) }}</pre>
                  </el-descriptions-item>
                  <el-descriptions-item :label="$t('auto.TraceMonitor.11.bd30be')">
                    <pre class="code-block">{{ fmt(item.tool.output) }}</pre>
                  </el-descriptions-item>
                  <el-descriptions-item v-if="item.tool.abnormal || item.tool.error" :label="$t('auto.TraceMonitor.13.c195df')">
                    <span class="err-text">{{ item.tool.error || '异常' }}</span>
                  </el-descriptions-item>
                </el-descriptions>
              </div>
            </el-card>
          </template>
        </div>
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.TraceMonitor.14.fce779')" name="health">
        <div class="stat-grid">
          <el-card v-for="s in healthStats" :key="s.label" shadow="hover" class="stat-card">
            <div class="stat-label">{{ s.label }}</div>
            <div class="stat-value" :class="s.cls">{{ s.value }}</div>
            <div class="stat-unit">{{ s.unit }}</div>
          </el-card>
        </div>
        <el-alert
          v-if="health && ((health.stuck_reachable || 0) > 0 || (health.stuck_unreachable || 0) > 0 || (health.sync_gap_count || 0) > 0)"
          type="warning"
          :closable="false"
          :title="$t('auto.TraceMonitor.15.8bd911')"
          class="mt-12"
        />
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.TraceMonitor.16.cdf0e4')" name="node">
        <el-alert v-if="nodeWindow" type="info" :closable="false" ::title="$t('auto.TraceMonitor.17.e9e4cd')" class="mb-12" />
        <el-table :data="nodeRows" border stripe size="small" :empty-text="$t('auto.TraceMonitor.18.21efd8')">
          <el-table-column prop="channel" :label="$t('auto.TraceMonitor.7.ebee05')" width="160" />
          <el-table-column prop="node" :label="$t('auto.TraceMonitor.19.3bf3c0')" width="200" />
          <el-table-column prop="total" :label="$t('auto.TraceMonitor.20.c1ce23')" width="100" />
          <el-table-column prop="avg_duration_ms" :label="$t('auto.TraceMonitor.21.188b05')" width="130" />
          <el-table-column prop="p95_duration_ms" label="P95(ms)" width="120" />
          <el-table-column :label="$t('auto.TraceMonitor.22.ab80c8')" width="120">
            <template #default="{ row }">{{ (row.abnormal_rate * 100).toFixed(1) }}%</template>
          </el-table-column>
          <el-table-column :label="$t('auto.TraceMonitor.23.fbc026')" width="100">
            <template #default="{ row }">
              <el-tag :type="row.abnormal_rate > 0.3 ? 'danger' : row.abnormal_rate > 0.05 ? 'warning' : 'success'" size="small">
                {{ row.abnormal_rate > 0.3 ? '差' : row.abnormal_rate > 0.05 ? '中' : '良' }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.TraceMonitor.24.b445ee')" name="latency">
        <el-table :data="latencyRows" border stripe size="small" :empty-text="$t('auto.TraceMonitor.18.21efd8')">
          <el-table-column prop="channel" :label="$t('auto.TraceMonitor.7.ebee05')" width="180" />
          <el-table-column prop="p50_ms" label="P50(ms)" width="140" />
          <el-table-column prop="p95_ms" label="P95(ms)" width="140" />
          <el-table-column prop="sample_size" :label="$t('auto.TraceMonitor.20.c1ce23')" width="120" />
        </el-table>
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.TraceMonitor.25.63b892')" name="lifecycle">
        <el-form :inline="true" @submit.prevent>
          <el-form-item :label="$t('auto.TraceMonitor.3.427bc5')">
            <el-input v-model.trim="lcForm.conversation_id" :placeholder="$t('auto.TraceMonitor.4.8b6554')" clearable style="width: 260px" />
          </el-form-item>
          <el-form-item label="Trace ID">
            <el-input v-model.trim="lcForm.trace_id" :placeholder="$t('auto.TraceMonitor.26.b0313a')" clearable style="width: 260px" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="loadingLc" @click="loadLifecycle">{{ $t('auto.TraceMonitor.62.bee912') }}</el-button>
          </el-form-item>
        </el-form>
        <el-table :data="lifecycleRows" border stripe size="small" :empty-text="$t('auto.TraceMonitor.18.21efd8')" class="mt-12">
          <el-table-column prop="trace_id" label="Trace" width="180" />
          <el-table-column prop="node" :label="$t('auto.TraceMonitor.19.3bf3c0')" width="180" />
          <el-table-column prop="channel" :label="$t('auto.TraceMonitor.7.ebee05')" width="120" />
          <el-table-column prop="direction" :label="$t('auto.TraceMonitor.27.a465db')" width="100" />
          <el-table-column prop="duration_ms" :label="$t('auto.TraceMonitor.28.7f62c0')" width="100" />
          <el-table-column :label="$t('auto.TraceMonitor.29.3fea7c')" width="100">
            <template #default="{ row }">
              <el-tag :type="row.abnormal ? 'danger' : 'success'" size="small">
                {{ row.abnormal ? '异常' : '正常' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="expected" :label="$t('auto.TraceMonitor.30.c633a3')" min-width="160" />
          <el-table-column prop="error" :label="$t('auto.TraceMonitor.31.c8d506')" min-width="160" show-overflow-tooltip />
        </el-table>
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.TraceMonitor.32.9b9378')" name="traces">
        <el-table :data="traceRows" border stripe size="small" :empty-text="$t('auto.TraceMonitor.18.21efd8')">
          <el-table-column prop="trace_id" label="Trace" width="180" />
          <el-table-column prop="conversation_id" :label="$t('auto.TraceMonitor.33.9a834e')" width="180" />
          <el-table-column prop="channel" :label="$t('auto.TraceMonitor.7.ebee05')" width="120" />
          <el-table-column prop="node_count" :label="$t('auto.TraceMonitor.34.93cb16')" width="90" />
          <el-table-column :label="$t('auto.TraceMonitor.35.8eaadf')" width="90">
            <template #default="{ row }">
              <el-tag :type="row.abnormal_count > 0 ? 'danger' : 'info'" size="small">{{ row.abnormal_count }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="first_at" :label="$t('auto.TraceMonitor.36.935c21')" width="180" />
          <el-table-column prop="last_at" :label="$t('auto.TraceMonitor.37.fea8c1')" width="180" />
          <el-table-column :label="$t('auto.TraceMonitor.38.2b6bc0')" width="100" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link @click="openTrace(row)">{{ $t('auto.TraceMonitor.67.607e7a') }}</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.TraceMonitor.13.c195df')" name="anomalies">
        <el-collapse v-if="anomaly" v-model="anomalyActive">
          <el-collapse-item :title="$t('auto.TraceMonitor.39.4a7da2')" name="sync_gap">
            <el-table :data="anomaly.sync_gap || []" border size="small" :empty-text="$t('auto.TraceMonitor.40.d81bb2')">
              <el-table-column prop="conversation_id" :label="$t('auto.TraceMonitor.33.9a834e')" width="200" />
              <el-table-column prop="channel" :label="$t('auto.TraceMonitor.7.ebee05')" width="120" />
              <el-table-column prop="message_count" :label="$t('auto.TraceMonitor.41.99e6f4')" width="100" />
            </el-table>
          </el-collapse-item>
          <el-collapse-item :title="$t('auto.TraceMonitor.42.d152be')" name="stuck_reachable">
            <el-table :data="anomaly.stuck_reachable || []" border size="small" :empty-text="$t('auto.TraceMonitor.40.d81bb2')">
              <el-table-column prop="conversation_id" :label="$t('auto.TraceMonitor.33.9a834e')" width="200" />
              <el-table-column prop="channel" :label="$t('auto.TraceMonitor.7.ebee05')" width="120" />
              <el-table-column prop="age_min" :label="$t('auto.TraceMonitor.43.0e6c96')" width="100" />
            </el-table>
          </el-collapse-item>
          <el-collapse-item :title="$t('auto.TraceMonitor.44.d45211')" name="stuck_unreachable">
            <el-table :data="anomaly.stuck_unreachable || []" border size="small" :empty-text="$t('auto.TraceMonitor.40.d81bb2')">
              <el-table-column prop="conversation_id" :label="$t('auto.TraceMonitor.33.9a834e')" width="200" />
              <el-table-column prop="channel" :label="$t('auto.TraceMonitor.7.ebee05')" width="120" />
              <el-table-column prop="age_min" :label="$t('auto.TraceMonitor.43.0e6c96')" width="100" />
            </el-table>
          </el-collapse-item>
          <el-collapse-item :title="$t('auto.TraceMonitor.45.9726d0')" name="unreachable">
            <el-table :data="anomaly.unreachable || []" border size="small" :empty-text="$t('auto.TraceMonitor.40.d81bb2')">
              <el-table-column prop="conversation_id" :label="$t('auto.TraceMonitor.33.9a834e')" width="200" />
              <el-table-column prop="channel" :label="$t('auto.TraceMonitor.7.ebee05')" width="120" />
            </el-table>
          </el-collapse-item>
          <el-collapse-item :title="$t('auto.TraceMonitor.46.209a8c')" name="node_abnormal">
            <el-table :data="anomaly.node_abnormal || []" border size="small" :empty-text="$t('auto.TraceMonitor.40.d81bb2')">
              <el-table-column prop="channel" :label="$t('auto.TraceMonitor.7.ebee05')" width="140" />
              <el-table-column prop="node" :label="$t('auto.TraceMonitor.19.3bf3c0')" width="180" />
              <el-table-column prop="abnormal_rate" :label="$t('auto.TraceMonitor.22.ab80c8')" width="120">
                <template #default="{ row }">{{ (row.abnormal_rate * 100).toFixed(1) }}%</template>
              </el-table-column>
            </el-table>
          </el-collapse-item>
        </el-collapse>
        <el-empty v-else :description="$t('auto.TraceMonitor.47.3a617d')" />
      </el-tab-pane>

      
      <el-tab-pane :label="$t('auto.TraceMonitor.48.03059f')" name="learn">
        <div class="learn-toolbar">
          <el-form :inline="true" @submit.prevent>
            <el-form-item :label="$t('auto.TraceMonitor.49.e54e95')">
              <el-input-number v-model="evalHours" :min="1" :max="168" :step="12" style="width: 140px" />
            </el-form-item>
            <el-form-item :label="$t('auto.TraceMonitor.50.cbda81')">
              <el-input-number v-model="evalLimit" :min="1" :max="200" :step="10" style="width: 130px" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="triggerLoading" @click="triggerEval">{{ $t('auto.TraceMonitor.68.34cc3b') }}</el-button>
              <el-button :loading="triggerLoading" @click="refreshLearn">{{ $t('auto.TraceMonitor.69.694fc5') }}</el-button>
            </el-form-item>
          </el-form>
          <el-alert
            v-if="triggerMsg"
            :type="triggerMsgType"
            :closable="false"
            :title="triggerMsg"
            class="mt-12"
            show-icon
          />
        </div>

        <el-divider content-position="left">{{ $t('auto.TraceMonitor.70.bd7701') }}</el-divider>
        <el-table :data="evalLogs" border stripe size="small" :empty-text="$t('auto.TraceMonitor.51.670f26')" class="mt-12">
          <el-table-column prop="created_at" :label="$t('auto.TraceMonitor.52.0ca2ff')" width="180" />
          <el-table-column prop="channel" :label="$t('auto.TraceMonitor.7.ebee05')" width="120" />
          <el-table-column prop="trace_id" label="Trace" width="180" show-overflow-tooltip />
          <el-table-column :label="$t('auto.TraceMonitor.53.108db1')" width="80">
            <template #default="{ row }">
              <el-tag :type="scoreType(row.score)" size="small">{{ row.score }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('auto.TraceMonitor.54.3a7170')" width="80">
            <template #default="{ row }">
              <el-tag :type="row.bad ? 'danger' : row.score >= 85 ? 'success' : 'warning'" size="small" effect="plain">
                {{ row.bad ? '差' : row.score >= 85 ? '优' : '中' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('auto.TraceMonitor.55.f29c54')" min-width="230">
            <template #default="{ row }">{{ dimText(row.dimensions) }}</template>
          </el-table-column>
          <el-table-column :label="$t('auto.TraceMonitor.56.41dfb0')" prop="reason" min-width="200" show-overflow-tooltip />
          <el-table-column :label="$t('auto.TraceMonitor.57.b20aad')" width="130">
            <template #default="{ row }">
              <span v-if="adjustedList(row.adjusted_chunks).length">{{ adjustedList(row.adjusted_chunks).length }} 个 chunk</span>
              <span v-else class="muted">{{ $t('auto.TraceMonitor.40.d81bb2') }}</span>
            </template>
          </el-table-column>
        </el-table>

        <el-divider content-position="left">{{ $t('auto.TraceMonitor.71.5e24ab') }}</el-divider>
        <el-table :data="weights" border stripe size="small" :empty-text="$t('auto.TraceMonitor.58.91336e')" class="mt-12">
          <el-table-column prop="id" label="Chunk ID" width="100" />
          <el-table-column prop="content" :label="$t('auto.TraceMonitor.59.2d711b')" min-width="280" show-overflow-tooltip />
          <el-table-column :label="$t('auto.TraceMonitor.60.4aac55')" width="120">
            <template #default="{ row }">
              <el-tag :type="row.weight < 1 ? 'danger' : row.weight > 1 ? 'success' : 'info'" size="small">
                {{ Number(row.weight).toFixed(2) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="hit_count" :label="$t('auto.TraceMonitor.61.e5852c')" width="100" />
        </el-table>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { MonitorApi as monitorApi } from '@/api/monitor'

const channels = ['douyin', 'xiaohongshu', 'wechat', 'telegram', 'whatsapp', 'instagram', 'facebook', 'line', 'kuaishou', 'weibo']

const { t } = useI18n()

const activeTab = ref('trace')
const form = reactive({ trace_id: '', conversation_id: '', msg_id: '', channel: '' })
const tree = ref(null)
const treeError = ref('')
const loadingTree = ref(false)
const expanded = reactive({})

const lcForm = reactive({ conversation_id: '', trace_id: '' })
const lifecycleRows = ref([])
const loadingLc = ref(false)

const health = ref(null)
const nodeRows = ref([])
const nodeWindow = ref('')
const latencyRows = ref([])
const traceRows = ref([])
const anomaly = ref(null)
const anomalyActive = ref(['sync_gap'])

const evalLogs = ref([]);
const weights = ref([])
const triggerLoading = ref(false)
const triggerMsg = ref('')
const triggerMsgType = ref('info')
const evalHours = ref(24)
const evalLimit = ref(20)
const learningLoaded = ref(false)

const NODE_LABEL = {
  ingest: t('auto.TraceMonitor.72.90375b'),
  ai_dispatch: 'AI 调度',
  outbound_enqueue: t('auto.TraceMonitor.73.f2fc7d'),
  inbox_sync: t('auto.TraceMonitor.74.72a37d'),
  downlink_fetch: t('auto.TraceMonitor.75.182a50'),
  delivered_ack: t('auto.TraceMonitor.76.aae560')
}

const nodeLabel = (n) => NODE_LABEL[n] || n || '—'

function fmt(v) {
  if (v === null || v === undefined || v === '') return '—'
  if (typeof v === 'string') return v
  try { return JSON.stringify(v, null, 2) } catch (e) { return String(v) }
}

function isAbnormal(sp) {
  return !!(sp.abnormal || sp.status === 'abnormal' || sp.error)
}

function statusType(sp) {
  if (sp.abnormal || sp.status === 'abnormal' || sp.error) return 'danger'
  if (sp.status === 'failed') return 'warning'
  return 'success'
}

function statusText(sp) {
  if (sp.abnormal || sp.status === 'abnormal' || sp.error) return t('auto.TraceMonitor.13.c195df')
  if (sp.status === 'failed') return t('auto.TraceMonitor.77.acd5cb')
  return t('auto.TraceMonitor.78.fd6e80')
}

function toggle(sp) {
  if (!sp || sp.id === undefined) return
  expanded[sp.id] = !expanded[sp.id]
}

const flattened = computed(() => {
  const t = tree.value
  if (!t || !t.spans) return []
  const out = []
  for (const sp of t.spans) {
    if (sp.span_kind !== 'lifecycle') continue
    out.push(sp)
    if (sp.node === 'ai_dispatch') {
      const turns = [...new Set(
        t.spans.filter(s => s.span_kind === 'agent_turn' && s.parent_node === 'ai_dispatch')
          .map(s => s.turn_index)
      )].sort((a, b) => a - b)
      for (const ti of turns) {
        const turn = t.spans.find(s => s.span_kind === 'agent_turn' && s.parent_node === 'ai_dispatch' && s.turn_index === ti)
        if (!turn) continue
        out.push({ _turn: true, ti, turn })
        const tools = t.spans
          .filter(s => s.span_kind === 'tool_call' && s.parent_node === 'ai_dispatch' && s.turn_index === ti)
          .sort((a, b) => (a.id || 0) - (b.id || 0))
        for (const tool of tools) out.push({ _tool: true, tool })
      }
    }
  }
  return out
});

const healthStats = computed(() => {
  const h = health.value
  if (!h) return []
  return [
    { label: t('auto.TraceMonitor.79.727844'), value: h.inbound_rate_per_min, unit: t('auto.TraceMonitor.80.a356de') },
    { label: t('auto.TraceMonitor.81.f1b6f8'), value: h.outbound_rate_per_min, unit: t('auto.TraceMonitor.80.a356de') },
    { label: t('auto.TraceMonitor.82.57a3b6'), value: h.pending_count, unit: t('auto.TraceMonitor.83.cc1bac'), cls: h.pending_count > 0 ? 'warn' : '' },
    { label: t('auto.TraceMonitor.84.4d8ddb'), value: h.oldest_pending_min, unit: t('auto.TraceMonitor.85.daf783'), cls: h.oldest_pending_min > 10 ? 'warn' : '' },
    { label: t('auto.TraceMonitor.86.f87f48'), value: h.delivered_count, unit: t('auto.TraceMonitor.83.cc1bac') },
    { label: t('auto.TraceMonitor.77.acd5cb'), value: h.failed_count, unit: t('auto.TraceMonitor.83.cc1bac'), cls: h.failed_count > 0 ? 'err' : '' },
    { label: t('auto.TraceMonitor.87.4e564c'), value: h.sync_gap_count, unit: t('auto.TraceMonitor.88.930882'), cls: h.sync_gap_count > 0 ? 'warn' : '' },
    { label: t('auto.TraceMonitor.89.1daba9'), value: h.stuck_reachable, unit: t('auto.TraceMonitor.88.930882'), cls: (h.stuck_reachable || 0) > 0 ? 'err' : '' },
    { label: t('auto.TraceMonitor.90.e7dc4b'), value: h.stuck_unreachable, unit: t('auto.TraceMonitor.88.930882'), cls: (h.stuck_unreachable || 0) > 0 ? 'err' : '' },
    { label: t('auto.TraceMonitor.91.a95f46'), value: h.abnormal_count, unit: t('auto.TraceMonitor.88.930882'), cls: h.abnormal_count > 0 ? 'warn' : '' }
  ]
})

function resetForm() {
  form.trace_id = ''
  form.conversation_id = ''
  form.msg_id = ''
  form.channel = ''
  tree.value = null
  treeError.value = ''
}

async function loadTree() {
  const params = {}
  if (form.trace_id) params.trace_id = form.trace_id
  if (form.conversation_id) params.conversation_id = form.conversation_id
  if (form.msg_id) params.msg_id = form.msg_id
  if (form.channel) params.channel = form.channel
  loadingTree.value = true
  treeError.value = ''
  try {
    const res = await monitorApi.traceTree(params)
    tree.value = res || null
    if (!tree.value || !tree.value.spans || !tree.value.spans.length) {
      treeError.value = t('auto.TraceMonitor.92.dba63f')
    } else {
      const firstAbn = tree.value.spans.find(s => isAbnormal(s));
      if (firstAbn) expanded[firstAbn.id] = true
    }
  } catch (e) {
    treeError.value = (e && e.message) || t('auto.TraceMonitor.93.0d66ed')
  } finally {
    loadingTree.value = false
  }
}

async function loadLifecycle() {
  const params = {}
  if (lcForm.conversation_id) params.conversation_id = lcForm.conversation_id
  if (lcForm.trace_id) params.trace_id = lcForm.trace_id
  if (!params.conversation_id && !params.trace_id) {
    ElMessage.warning(t('auto.TraceMonitor.94.deaf34'))
    return
  }
  loadingLc.value = true
  try {
    lifecycleRows.value = await monitorApi.lifecycle(params)
  } catch (e) {
    ElMessage.error((e && e.message) || t('auto.TraceMonitor.93.0d66ed'))
  } finally {
    loadingLc.value = false
  }
}

function openTrace(row) {
  form.conversation_id = row.conversation_id || ''
  form.trace_id = row.trace_id || ''
  form.msg_id = ''
  form.channel = ''
  activeTab.value = 'trace'
  loadTree()
}

async function loadHealth() {
  try { health.value = await monitorApi.health() } catch (e) { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
}
async function loadNodeHealth() {
  try {
    const res = await monitorApi.nodeHealth()
    nodeRows.value = res.nodes || []
    nodeWindow.value = res.window || ''
  } catch (e) { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
}
async function loadLatency() {
  try { latencyRows.value = await monitorApi.latency() } catch (e) { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
}
async function loadTraces() {
  try { traceRows.value = await monitorApi.traces({ limit: 50 }) } catch (e) { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
}
async function loadAnomalies() {
  try { anomaly.value = await monitorApi.anomalies() } catch (e) { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
}

function parseJSON(v, fallback) {
  if (v === null || v === undefined || v === '') return fallback
  if (typeof v === 'object') return v
  try { return JSON.parse(v) } catch (e) { return fallback }
}
function dimText(dims) {
  const d = parseJSON(dims, {})
  const order = [
    ['relevance', t('auto.TraceMonitor.95.326c88')],
    ['accuracy', t('auto.TraceMonitor.96.a8e10d')],
    ['usefulness', t('auto.TraceMonitor.97.93fa21')],
    ['safety', t('auto.TraceMonitor.98.6d266c')]
  ]
  return order.map(([k, label]) => `${label} ${d[k] != null ? d[k] : '—'}`).join(' · ')
}
function adjustedList(v) {
  const a = parseJSON(v, [])
  return Array.isArray(a) ? a : []
}
function scoreType(s) {
  if (s < 60) return 'danger'
  if (s >= 85) return 'success'
  return 'warning'
}
async function loadEvalLogs() {
  try { evalLogs.value = await monitorApi.evalLogs({ limit: 50 }) } catch (e) { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
}
async function loadWeights() {
  try { weights.value = await monitorApi.knowledgeWeights({ limit: 50 }) } catch (e) { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
}
async function refreshLearn() {
  await Promise.all([loadEvalLogs(), loadWeights()])
}
async function triggerEval() {
  triggerLoading.value = true
  triggerMsg.value = ''
  try {
    const res = await monitorApi.triggerEval({ hours: evalHours.value, limit: evalLimit.value })
    const processed = (res && res.processed) || 0
    triggerMsgType.value = 'success'
    triggerMsg.value = `已处理 ${processed} 条 trace（评分 + 调整知识库权重），详见下方记录`
    await refreshLearn()
  } catch (e) {
    triggerMsgType.value = 'error'
    triggerMsg.value = (e && e.message) || t('auto.TraceMonitor.99.b0b8b8')
  } finally {
    triggerLoading.value = false
  }
}

watch(activeTab, (tab) => {
  if (tab === 'health' && !health.value) loadHealth()
  if (tab === 'node' && !nodeRows.value.length) loadNodeHealth()
  if (tab === 'latency' && !latencyRows.value.length) loadLatency()
  if (tab === 'traces' && !traceRows.value.length) loadTraces()
  if (tab === 'anomalies' && !anomaly.value) loadAnomalies()
  if (tab === 'learn' && !learningLoaded.value) {
    learningLoaded.value = true
    refreshLearn()
  }
});
</script>

<style scoped>
.trace-monitor { padding: 8px; }
.mb-12 { margin-bottom: 12px; }
.mt-12 { margin-top: 12px; }
.trace-query { background: #f7f8fa; padding: 12px 12px 0; border-radius: 6px; }
.trace-summary { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.trace-count { margin-left: auto; color: #909399; font-size: 13px; }

.span-card { margin-bottom: 10px; border-radius: 8px; }
.span-card.abnormal { border-color: #f56c6c; }
.span-card.agent { border-left: 4px solid #409eff; }
.turn-card { margin-left: 28px; border-left: 4px solid #e6a23c; }
.tool-card { margin-left: 56px; border-left: 4px solid #67c23a; background: #fafcfa; }
.span-head { display: flex; align-items: center; gap: 10px; cursor: pointer; flex-wrap: wrap; }
.span-head .caret { width: 14px; color: #909399; }
.span-kind { font-size: 12px; padding: 1px 6px; border-radius: 4px; color: #fff; }
.kind-lifecycle { background: #409eff; }
.kind-turn { background: #e6a23c; }
.kind-tool { background: #67c23a; }
.span-title { font-weight: 600; }
.span-node { color: #909399; font-size: 12px; font-family: monospace; }
.span-dur { color: #606266; font-size: 13px; }
.span-time { color: #c0c4cc; font-size: 12px; }
.span-body { margin-top: 12px; }
.code-block {
  margin: 0;
  max-height: 280px;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-all;
  font-family: monospace;
  font-size: 12px;
  background: #f5f7fa;
  padding: 8px;
  border-radius: 4px;
}
.err-text { color: #f56c6c; }

.stat-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); gap: 12px; }
.stat-card { text-align: center; }
.stat-label { color: #909399; font-size: 13px; }
.stat-value { font-size: 26px; font-weight: 700; margin: 6px 0; }
.stat-value.warn { color: #e6a23c; }
.stat-value.err { color: #f56c6c; }
.stat-unit { color: #c0c4cc; font-size: 12px; }
.muted { color: #c0c4cc; }
.learn-toolbar { background: #f7f8fa; padding: 12px 12px 0; border-radius: 6px; }
</style>
