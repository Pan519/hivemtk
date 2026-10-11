<template>
  <div class="geo-page p-4">
    <!-- 概览卡片区 -->
    <el-row :gutter="12" class="mb-4">
      <el-col :xs="12" :sm="6">
        <el-card shadow="hover">
          <div class="stat-block">
            <div class="stat-label">AI 可见率（近 {{ days }} 天）</div>
            <div class="stat-value" :style="{ color: visibilityColor(summary.current_avg) }">
              {{ pct(summary.current_avg) }}
            </div>
            <div class="stat-sub" v-if="hasCompare">
              环比
              <span :class="summary.change >= 0 ? 'trend-up' : 'trend-down'">
                {{ summary.change >= 0 ? '↑' : '↓' }}{{ pct(Math.abs(summary.change)) }}
              </span>
            </div>
            <div class="stat-sub muted" v-else>{{ $t('auto.VisibilityBoard.15.7b04cb') }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card shadow="hover">
          <div class="stat-block">
            <div class="stat-label">{{ $t('auto.VisibilityBoard.16.95660d') }}</div>
            <div class="stat-value">{{ formatInt(summary.total_probes) }}</div>
            <div class="stat-sub">品牌命中 {{ formatInt(summary.total_brand_hits) }} 次</div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card shadow="hover">
          <div class="stat-block">
            <div class="stat-label">{{ $t('auto.VisibilityBoard.17.96247f') }}</div>
            <div class="stat-value">{{ engines.length }}</div>
            <div class="stat-sub">最弱引擎：{{ weakestEngine || '—' }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card shadow="hover">
          <div class="stat-block">
            <div class="stat-label">{{ $t('auto.VisibilityBoard.18.3ae768') }}</div>
            <div class="stat-value" :class="totalNegative > 0 ? 'text-danger' : ''">{{ formatInt(totalNegative) }}</div>
            <div class="stat-sub">{{ $t('auto.VisibilityBoard.19.e614c4') }}</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 筛选 + 操作 -->
    <el-card class="mb-4" shadow="never">
      <div class="flex items-center gap-3 flex-wrap">
        <el-radio-group v-model="days" @change="loadAll">
          <el-radio-button :value="7">{{ $t('auto.VisibilityBoard.20.d765ef') }}</el-radio-button>
          <el-radio-button :value="30">{{ $t('auto.VisibilityBoard.21.c059be') }}</el-radio-button>
          <el-radio-button :value="90">{{ $t('auto.VisibilityBoard.22.9be1fd') }}</el-radio-button>
        </el-radio-group>
        <el-button type="primary" :loading="probing" @click="runProbe" plain>
          <el-icon><VideoPlay /></el-icon><span>{{ $t('auto.VisibilityBoard.23.b12607') }}</span>
        </el-button>
        <el-button @click="loadAll" :loading="loading">{{ $t('auto.VisibilityBoard.24.694fc5') }}</el-button>
      </div>
    </el-card>

    <!-- 可见性趋势折线图 -->
    <el-card class="mb-4" shadow="never">
      <template #header>
        <span class="font-bold">{{ $t('auto.VisibilityBoard.25.36a180') }}</span>
      </template>
      <div v-if="!loading && !hasTrendData" class="py-12 text-center text-gray-400">
        暂无观测数据。点击上方「立即探测一轮」，或在「多模型验证」页运行验证后回来看趋势。
      </div>
      <div ref="trendChartRef" v-show="hasTrendData" style="height: 320px"></div>
    </el-card>

    <!-- 引擎对比 -->
    <el-card class="mb-4" shadow="never">
      <template #header>
        <span class="font-bold">{{ $t('auto.VisibilityBoard.26.12ed10') }}</span>
      </template>
      <el-empty v-if="!loading && engines.length === 0" :description="$t('auto.VisibilityBoard.1.e32543')" :image-size="60" />
      <el-table v-else :data="engines" v-loading="loading" size="default">
        <el-table-column prop="engine" :label="$t('auto.VisibilityBoard.2.6c1b19')" min-width="160">
          <template #default="{ row }">
            <span class="font-semibold">{{ row.engine }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.VisibilityBoard.3.b89857')" min-width="220">
          <template #default="{ row }">
            <div class="flex items-center gap-2">
              <el-progress
                :percentage="Math.round(row.visibility * 100)"
                :stroke-width="10"
                :color="visibilityColor(row.visibility)"
                style="flex: 1"
              />
              <span class="font-mono w-12 text-right">{{ pct(row.visibility) }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="probe_count" :label="$t('auto.VisibilityBoard.4.500003')" width="100" align="right" sortable />
        <el-table-column prop="brand_hits" :label="$t('auto.VisibilityBoard.5.648129')" width="100" align="right" sortable />
        <el-table-column :label="$t('auto.VisibilityBoard.6.a76c59')" width="100" align="right">
          <template #default="{ row }">{{ (row.avg_citations || 0).toFixed(2) }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.VisibilityBoard.7.f3c630')" width="90" align="right">
          <template #default="{ row }">
            <el-tag v-if="row.negative_count > 0" type="danger" size="small">{{ row.negative_count }}</el-tag>
            <span v-else class="text-gray-400">0</span>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 探针运行记录 -->
    <el-card shadow="never">
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-bold">探针运行记录（最近 {{ runsLimit }} 条）</span>
          <div class="flex items-center gap-2">
            <el-select v-model="runsEngine" :placeholder="$t('auto.VisibilityBoard.8.1d4c0f')" clearable size="small" style="width: 160px" @change="loadRuns">
              <el-option v-for="e in engines" :key="e.engine" :label="e.engine" :value="e.engine" />
            </el-select>
            <el-button size="small" @click="loadRuns" :loading="runsLoading">{{ $t('auto.VisibilityBoard.24.694fc5') }}</el-button>
          </div>
        </div>
      </template>
      <el-table :data="runs" v-loading="runsLoading" size="small">
        <el-table-column prop="created_at" :label="$t('auto.VisibilityBoard.9.19fcb9')" width="160">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column prop="engine" :label="$t('auto.VisibilityBoard.2.6c1b19')" width="150">
          <template #default="{ row }"><el-tag size="small" effect="plain">{{ row.engine }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="query" :label="$t('auto.VisibilityBoard.10.042e36')" min-width="180" show-overflow-tooltip />
        <el-table-column :label="$t('auto.VisibilityBoard.11.09307c')" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="row.brand_mentioned ? 'success' : 'info'" size="small">
              {{ row.brand_mentioned ? '提及' : '未提及' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.VisibilityBoard.12.2b1c16')" width="80" align="center">
          <template #default="{ row }">
            <el-tag size="small" :type="sentimentType(row.sentiment)">{{ sentimentLabel(row.sentiment) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="latency_ms" :label="$t('auto.VisibilityBoard.13.7f62c0')" width="100" align="right" />
        <el-table-column :label="$t('auto.VisibilityBoard.14.54e5c9')" min-width="240">
          <template #default="{ row }">
            <el-tooltip :content="snippet(row.response)" placement="top" :show-after="200">
              <span class="text-gray-500 text-xs">{{ snippet(row.response) }}</span>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { VideoPlay } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { safeInit } from '@/utils/echarts'
import { getEngineCompare, listProbeRuns, probeAllEngines } from '@/api/geoProbe.js'
import { toList } from '@/utils/list'

const { t } = useI18n()

const days = ref(30)
const loading = ref(false)
const summary = ref({})
const engines = ref([])
const daily = ref([])

const runs = ref([])
const runsLoading = ref(false)
const runsEngine = ref('')
const runsLimit = 50
const probing = ref(false)

const trendChartRef = ref(null)
let chartInstance = null

const hasTrendData = computed(() => (summary.value?.points?.length || 0) > 0)
const hasCompare = computed(() => (summary.value?.total_probes || 0) > 0 && (summary.value?.previous_avg || 0) > 0)
const totalNegative = computed(() => engines.value.reduce((s, e) => s + (e.negative_count || 0), 0))
const weakestEngine = computed(() => {
  if (!engines.value.length) return ''
  const worst = [...engines.value].sort((a, b) => a.visibility - b.visibility)[0]
  return `${worst.engine} ${pct(worst.visibility)}`
})

const pct = (v) => `${((Number(v) || 0) * 100).toFixed(1)}%`
const visibilityColor = (v) => {
  const n = (Number(v) || 0) * 100
  if (n >= 60) return '#22c55e'
  if (n >= 30) return '#f59e0b'
  return '#ef4444'
}
const formatInt = (n) => Number(n || 0).toLocaleString('zh-CN')
const formatTime = (t) => (t ? new Date(t).toLocaleString('zh-CN', { hour12: false }) : '-')
const sentimentType = (s) => ({ positive: 'success', negative: 'danger' }[s] || 'info')
const sentimentLabel = (s) => ({ positive: t('auto.VisibilityBoard.27.1b614b'), negative: t('auto.VisibilityBoard.7.f3c630') }[s] || t('auto.VisibilityBoard.28.62731a'))
const snippet = (text) => {
  const t = String(text || '').replace(/\s+/g, ' ').trim()
  return t.length > 60 ? t.slice(0, 60) + '…' : t || '—'
}

const renderTrendChart = () => {
  const points = summary.value?.points || []
  if (!trendChartRef.value || points.length === 0) return
  if (!chartInstance) {
    chartInstance = safeInit(trendChartRef.value)
    window.addEventListener('resize', handleResize)
  }
  const dates = points.map((p) => p.date)
  chartInstance.setOption({
    tooltip: { trigger: 'axis' },
    legend: { data: [t('auto.VisibilityBoard.3.b89857'), t('auto.VisibilityBoard.4.500003')] },
    grid: { left: 60, right: 60, top: 40, bottom: 30 },
    xAxis: { type: 'category', data: dates },
    yAxis: [
      { type: 'value', name: t('auto.VisibilityBoard.3.b89857'), min: 0, max: 1, axisLabel: { formatter: (v) => `${Math.round(v * 100)}%` } },
      { type: 'value', name: t('auto.VisibilityBoard.4.500003'), minInterval: 1 }
    ],
    series: [
      {
        name: t('auto.VisibilityBoard.3.b89857'),
        type: 'line',
        smooth: true,
        data: points.map((p) => p.visibility),
        itemStyle: { color: '#409eff' },
        areaStyle: { opacity: 0.12 }
      },
      {
        name: t('auto.VisibilityBoard.4.500003'),
        type: 'bar',
        yAxisIndex: 1,
        data: points.map((p) => p.probe_count),
        itemStyle: { color: '#d3dce6', opacity: 0.7 },
        barMaxWidth: 18
      }
    ]
  }, true)
}

const handleResize = () => chartInstance && chartInstance.resize()

const loadAll = async () => {
  loading.value = true
  try {
    const data = await getEngineCompare(days.value)
    summary.value = data?.summary || {}
    engines.value = data?.engines || []
    daily.value = data?.daily || []
    await nextTick()
    renderTrendChart()
    await loadRuns()
  } catch (e) {
    ElMessage.error(e?.message || t('auto.VisibilityBoard.29.00141a'))
  } finally {
    loading.value = false
  }
}

const loadRuns = async () => {
  runsLoading.value = true
  try {
    const res = await listProbeRuns(runsEngine.value || undefined, runsLimit)
    runs.value = toList(res)
  } catch {
    runs.value = []
  } finally {
    runsLoading.value = false
  }
}

const runProbe = async () => {
  probing.value = true
  try {
    await probeAllEngines([], t('auto.VisibilityBoard.30.7e7bd4'))
    ElMessage.success(t('auto.VisibilityBoard.31.eea9ca'))
    setTimeout(loadAll, 8000)
  } catch (e) {
    ElMessage.error(e?.message || t('auto.VisibilityBoard.32.51049e'))
  } finally {
    probing.value = false
  }
}

onMounted(loadAll)
onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  if (chartInstance) chartInstance.dispose()
})
</script>

<style lang="scss" scoped>
.stat-block {
  text-align: center;
  padding: 4px 0;
}
.stat-label {
  font-size: 12px;
  color: #909399;
}
.stat-value {
  font-size: 28px;
  font-weight: 700;
  margin: 4px 0;
  color: #303133;
}
.stat-value.text-danger {
  color: #ef4444;
}
.stat-sub {
  font-size: 12px;
  color: #606266;
}
.stat-sub.muted {
  color: #c0c4cc;
}
.trend-up {
  color: #22c55e;
  font-weight: 600;
}
.trend-down {
  color: #ef4444;
  font-weight: 600;
}
.font-mono {
  font-family: 'JetBrains Mono', Consolas, monospace;
}
</style>
