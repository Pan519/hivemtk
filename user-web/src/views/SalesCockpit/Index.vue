<template>
  <div class="sales-cockpit">
    <PageHeader :title="$t('auto.Index.1.a97d88')" :subtitle="$t('auto.Index.2.a47a3b')" />

    
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="capability-card">
          <div class="capability-icon">🤖</div>
          <div class="capability-name">{{ $t('auto.Index.19.2778e1') }}</div>
          <div class="capability-stat">
            <span class="big-num">{{ cockpit.react.totalRuns }}</span>
            <span class="unit">{{ $t('auto.Index.20.bb24ac') }}</span>
          </div>
          <div class="capability-meta">{{ $t('auto.Index.21.ef320c') }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="capability-card">
          <div class="capability-icon">📋</div>
          <div class="capability-name">{{ $t('auto.Index.22.a191ae') }}</div>
          <div class="capability-stat">
            <span class="big-num">{{ cockpit.sop.executions }}</span>
            <span class="unit">{{ $t('auto.Index.23.46e386') }}</span>
          </div>
          <div class="capability-meta">{{ $t('auto.Index.24.be6e40') }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="capability-card">
          <div class="capability-icon">🔍</div>
          <div class="capability-name">{{ $t('auto.Index.25.423e78') }}</div>
          <div class="capability-stat">
            <span class="big-num">{{ cockpit.rag.queries }}</span>
            <span class="unit">{{ $t('auto.Index.26.68ea25') }}</span>
          </div>
          <div class="capability-meta">{{ $t('auto.Index.27.eae9c4') }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="capability-card">
          <div class="capability-icon">📡</div>
          <div class="capability-name">{{ $t('auto.Index.28.58dfae') }}</div>
          <div class="capability-stat">
            <span class="big-num">{{ cockpit.reach.sentToday }}</span>
            <span class="unit">{{ $t('auto.Index.29.ec5058') }}</span>
          </div>
          <div class="capability-meta">{{ $t('auto.Index.30.07dc3a') }}</div>
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
        <el-card :header="$t('auto.Index.3.bb42b1')">
          <div ref="intentChart" style="height: 300px;"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card :header="$t('auto.Index.4.306bd4')">
          <el-table :data="cockpit.llmRoutes" stripe>
            <el-table-column prop="scenario" :label="$t('auto.Index.5.c93165')" />
            <el-table-column prop="primary" :label="$t('auto.Index.6.ac51ee')" />
            <el-table-column prop="fallback" :label="$t('auto.Index.7.2cf8d0')" />
            <el-table-column prop="qps" label="QPS" />
            <el-table-column prop="errorRate" :label="$t('auto.Index.8.0c8524')" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    
    <el-row :gutter="20" class="mt-20">
      <el-col :span="12">
        <el-card :header="$t('auto.Index.9.47af0a')">
          <el-table :data="cockpit.channelHealth" stripe>
            <el-table-column prop="channel" :label="$t('auto.Index.10.ebee05')" />
            <el-table-column prop="qpsLimit" :label="$t('auto.Index.11.e90619')" />
            <el-table-column prop="qpsUsed" :label="$t('auto.Index.12.0cdb5e')" />
            <el-table-column :label="$t('auto.Index.13.3fea7c')">
              <template #default="{ row }">
                <el-tag :type="row.qpsUsed / row.qpsLimit > 0.8 ? 'warning' : 'success'">
                  {{ row.qpsUsed / row.qpsLimit > 0.8 ? '繁忙' : '正常' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="errorRate" :label="$t('auto.Index.8.0c8524')" />
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card :header="$t('auto.Index.14.4e1818')">
          <el-table :data="cockpit.topTools" stripe>
            <el-table-column type="index" label="#" width="50" />
            <el-table-column prop="name" :label="$t('auto.Index.15.20dce2')" />
            <el-table-column prop="category" :label="$t('auto.Index.16.d0771a')" />
            <el-table-column prop="calls" :label="$t('auto.Index.17.00d5bd')" />
            <el-table-column prop="avgLatency" :label="$t('auto.Index.18.a1df2e')" />
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
  closureText: '—', closureMeta: t('auto.Index.31.fb4ca1'),
  collectionText: '—', collectionMeta: t('auto.Index.31.fb4ca1'),
  overdueText: '—', overdueMeta: t('auto.Index.31.fb4ca1')
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
      closureText: '—', closureMeta: t('auto.Index.32.866b79'),
      collectionText: '—', collectionMeta: t('auto.Index.32.866b79'),
      overdueText: '—', overdueMeta: t('auto.Index.32.866b79')
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
    name: r.intent || t('auto.Index.33.ecf7eb'),
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
