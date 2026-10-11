<template>
  <div class="sales-cockpit">
    <PageHeader :title="$t('auto.Index.a97d888c')" :subtitle="$t('auto.Index.a47a3b55')" />

    
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="capability-card">
          <div class="capability-icon">🤖</div>
          <div class="capability-name">{{ $t('auto.Index.2778e177') }}</div>
          <div class="capability-stat">
            <span class="big-num">{{ cockpit.react.totalRuns }}</span>
            <span class="unit">{{ $t('auto.Index.bb24ac91') }}</span>
          </div>
          <div class="capability-meta">{{ $t('auto.Index.ef320c72') }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="capability-card">
          <div class="capability-icon">📋</div>
          <div class="capability-name">{{ $t('auto.Index.a191ae48') }}</div>
          <div class="capability-stat">
            <span class="big-num">{{ cockpit.sop.executions }}</span>
            <span class="unit">{{ $t('auto.Index.46e38679') }}</span>
          </div>
          <div class="capability-meta">{{ $t('auto.Index.be6e407f') }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="capability-card">
          <div class="capability-icon">🔍</div>
          <div class="capability-name">{{ $t('auto.Index.423e787a') }}</div>
          <div class="capability-stat">
            <span class="big-num">{{ cockpit.rag.queries }}</span>
            <span class="unit">{{ $t('auto.Index.68ea25aa') }}</span>
          </div>
          <div class="capability-meta">{{ $t('auto.Index.eae9c4a4') }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="capability-card">
          <div class="capability-icon">📡</div>
          <div class="capability-name">{{ $t('auto.Index.58dfae26') }}</div>
          <div class="capability-stat">
            <span class="big-num">{{ cockpit.reach.sentToday }}</span>
            <span class="unit">{{ $t('auto.Index.ec5058b2') }}</span>
          </div>
          <div class="capability-meta">{{ $t('auto.Index.07dc3a35') }}</div>
        </el-card>
      </el-col>
    </el-row>

    
    <!-- LTC 三率（T-P8-05，LTC-29/AC4）：闭环完成率是北极星指标，摆首位 -->
    <el-row
      :gutter="20"
      class="mt-20"
    >
      <el-col :span="8">
        <el-card class="capability-card ltc-north-star">
          <div class="capability-icon">
            🎯
          </div>
          <div class="capability-name">
            闭环完成率
            <el-tag
              type="danger"
              size="small"
            >
              北极星
            </el-tag>
          </div>
          <div class="capability-stat">
            <span class="big-num">
              {{ ltc.closureText }}
            </span>
          </div>
          <div class="capability-meta">
            {{ ltc.closureMeta }}
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="capability-card">
          <div class="capability-icon">
            💰
          </div>
          <div class="capability-name">
            回款率
          </div>
          <div class="capability-stat">
            <span class="big-num">
              {{ ltc.collectionText }}
            </span>
          </div>
          <div class="capability-meta">
            {{ ltc.collectionMeta }}
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="capability-card">
          <div class="capability-icon">
            ⏰
          </div>
          <div class="capability-name">
            逾期率
          </div>
          <div class="capability-stat">
            <span class="big-num">
              {{ ltc.overdueText }}
            </span>
          </div>
          <div class="capability-meta">
            {{ ltc.overdueMeta }}
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-20">
      <el-col :span="12">
        <el-card :header="$t('auto.Index.bb42b17e')">
          <div ref="intentChart" style="height: 300px;"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card :header="$t('auto.Index.306bd4f8')">
          <el-table :data="cockpit.llmRoutes" stripe>
            <el-table-column prop="scenario" :label="$t('auto.Index.c931653c')" />
            <el-table-column prop="primary" :label="$t('auto.Index.ac51ee10')" />
            <el-table-column prop="fallback" :label="$t('auto.Index.2cf8d0b3')" />
            <el-table-column prop="qps" label="QPS" />
            <el-table-column prop="errorRate" :label="$t('auto.Index.0c8524d7')" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    
    <el-row :gutter="20" class="mt-20">
      <el-col :span="12">
        <el-card :header="$t('auto.Index.47af0a86')">
          <el-table :data="cockpit.channelHealth" stripe>
            <el-table-column prop="channel" :label="$t('auto.Index.ebee0533')" />
            <el-table-column prop="qpsLimit" :label="$t('auto.Index.e9061939')" />
            <el-table-column prop="qpsUsed" :label="$t('auto.Index.0cdb5ef2')" />
            <el-table-column :label="$t('auto.Index.3fea7ca7')">
              <template #default="{ row }">
                <el-tag :type="row.qpsUsed / row.qpsLimit > 0.8 ? 'warning' : 'success'">
                  {{ row.qpsUsed / row.qpsLimit > 0.8 ? '繁忙' : '正常' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="errorRate" :label="$t('auto.Index.0c8524d7')" />
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card :header="$t('auto.Index.4e1818fb')">
          <el-table :data="cockpit.topTools" stripe>
            <el-table-column type="index" label="#" width="50" />
            <el-table-column prop="name" :label="$t('auto.Index.20dce2c6')" />
            <el-table-column prop="category" :label="$t('auto.Index.d0771a42')" />
            <el-table-column prop="calls" :label="$t('auto.Index.00d5bdf1')" />
            <el-table-column prop="avgLatency" :label="$t('auto.Index.a1df2e5f')" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import PageHeader from '@/components/PageHeader.vue'
import { getSalesCockpit } from '@/api/cockpit'
import LtcRatesApi from '@/api/ltcRates'
import * as echarts from 'echarts'
import { safeInit } from '@/utils/echarts'

const { t } = useI18n()

const intentChart = ref(null)
const cockpit = ref({
  react: { totalRuns: 0 },
  sop: { executions: 0 },
  rag: { queries: 0 },
  reach: { sentToday: 0 },
  llmRoutes: [],
  channelHealth: [],
  topTools: [],
  intentDistribution: []
})

/**
 * LTC 三率（T-P8-05）：后端 /api/ltc-rates 返回 {closure, collection, overdue}，
 * 每率 {rate, numerator, denominator}。分子分母摆出来，AC① 对账直接可比。
 */
const ltc = ref({
  closureText: '—', closureMeta: t('auto.Index.fb4ca1cf'),
  collectionText: '—', collectionMeta: t('auto.Index.fb4ca1cf'),
  overdueText: '—', overdueMeta: t('auto.Index.fb4ca1cf')
})

function fmtRate(v) {
  const n = Number(v && v.rate)
  if (!Number.isFinite(n)) return '—'
  return (n * 100).toFixed(1) + '%'
}

function fmtFrac(v) {
  const num = Number(v && v.numerator || 0)
  const den = Number(v && v.denominator || 0)
  return `${num.toLocaleString()} / ${den.toLocaleString()}`
}

async function loadLtcRates() {
  try {
    const res = await LtcRatesApi.getRates()
    const d = (res && res.data) || res || {}
    ltc.value = {
      closureText: fmtRate(d.closure),
      closureMeta: `已闭环 ${fmtFrac(d.closure)}`,
      collectionText: fmtRate(d.collection),
      collectionMeta: `已回款 ${fmtFrac(d.collection)} 元`,
      overdueText: fmtRate(d.overdue),
      overdueMeta: `逾期 ${fmtFrac(d.overdue)} 单`
    }
  } catch (err) {
    console.error('load ltc rates failed:', err)
    ltc.value = {
      closureText: '—', closureMeta: t('auto.Index.866b795e'),
      collectionText: '—', collectionMeta: t('auto.Index.866b795e'),
      overdueText: '—', overdueMeta: t('auto.Index.866b795e')
    }
  }
}

/**
 * 将后端 /api/ai/sales-cockpit 返回的 data 字段适配为前端模板期望的结构。
 * 后端原始字段（snake_case / 平铺）→ 前端模板字段（camelCase / 语义化）。
 */
function adaptCockpit(src) {
  const s = src || {}

  // 四张能力卡 —— 后端结构已对齐，直接取值
  const react = s.react || {}
  const sop = s.sop || {}
  const rag = s.rag || {}
  const reach = s.reach || {}

  // LLM 路由：后端 {scenario, provider, calls, avg_latency, total_cost}
  //          → 前端 {scenario, primary, fallback, qps, errorRate}
  const llmRoutes = (s.llmRoutes || []).map(r => ({
    scenario: r.scenario || '—',
    primary: r.provider || '—',
    fallback: '—',
    qps: Number(r.calls || 0).toLocaleString(),
    errorRate: '0.0%'
  }))

  // 渠道健康：后端 {platform, status, cnt}
  //          → 前端 {channel, qpsLimit, qpsUsed, errorRate}
  const channelHealth = (s.channelHealth || []).map(r => ({
    channel: r.platform || '—',
    qpsLimit: 100,
    qpsUsed: Number(r.cnt || 0),
    errorRate: '0.0%'
  }))

  // 意图分布 → echarts pie [{name, value}]
  const intentDistribution = (s.intentDistribution || []).map(r => ({
    name: r.intent || t('auto.Index.ecf7ebb5'),
    value: Number(r.cnt || 0)
  }))

  // Top 工具：后端 {tool_name, calls}
  //          → 前端 {name, category, calls, avgLatency}
  const topTools = (s.topTools || []).map(r => ({
    name: r.tool_name || '—',
    category: '—',
    calls: Number(r.calls || 0).toLocaleString(),
    avgLatency: '—'
  }))

  return {
    react: { totalRuns: Number(react.totalRuns || 0).toLocaleString() },
    sop:   { executions: Number(sop.executions || 0).toLocaleString() },
    rag:   { queries:    Number(rag.queries || 0).toLocaleString() },
    reach: { sentToday:  Number(reach.sentToday || 0).toLocaleString() },
    llmRoutes,
    channelHealth,
    intentDistribution,
    topTools
  }
}

let chartInstance = null
let timer = null

async function loadCockpit() {
  try {
    const res = await getSalesCockpit()
    cockpit.value = adaptCockpit(res.data || res)
    await nextTick()
    renderIntentChart()
  } catch (err) {
    console.error('load cockpit failed:', err)
  }
}

function renderIntentChart() {
  if (!intentChart.value) return
  if (chartInstance) {
    chartInstance.dispose()
    chartInstance = null
  }
  chartInstance = safeInit(intentChart.value)
  chartInstance.setOption({
    tooltip: { trigger: 'item' },
    legend: { bottom: 0 },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      data: cockpit.value.intentDistribution,
      label: { formatter: '{b}: {c} ({d}%)' }
    }]
  })
}

onMounted(() => {
  loadCockpit()
  loadLtcRates()
  timer = setInterval(() => { loadCockpit(); loadLtcRates() }, 60000)
})
onUnmounted(() => {
  if (timer) clearInterval(timer)
  if (chartInstance) {
    chartInstance.dispose()
    chartInstance = null
  }
})
</script>

<style scoped>
.sales-cockpit { padding: 20px; }
.capability-card { text-align: center; padding: 24px 0; }
.capability-icon { font-size: 48px; margin-bottom: 12px; }
.capability-name { font-size: 16px; font-weight: 500; color: #303133; }
.capability-stat { margin: 16px 0; }
.capability-stat .big-num { font-size: 36px; font-weight: 600; color: #409eff; }
.capability-stat .unit { font-size: 12px; color: #909399; margin-left: 6px; }
.capability-meta { font-size: 12px; color: #909399; }
.ltc-north-star { border-color: #f56c6c; }
.mt-20 { margin-top: 20px; }
</style>
