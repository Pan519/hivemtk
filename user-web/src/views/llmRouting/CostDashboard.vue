<template>
  <div class="llm-cost">
    <el-row :gutter="16" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-label">{{ $t('auto.CostDashboard.3.a539a9') }}</div>
          <div class="stat-value">¥{{ Number(stats.total_cost || 0).toFixed(2) }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-label">{{ $t('auto.CostDashboard.4.b615f5') }}</div>
          <div class="stat-value">{{ formatNumber(stats.total_tokens) }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card warning">
          <div class="stat-label">{{ $t('auto.CostDashboard.5.0bedb9') }}</div>
          <div class="stat-value">{{ stats.total_failed || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card success">
          <div class="stat-label">{{ $t('auto.CostDashboard.6.11d2f8') }}</div>
          <div class="stat-value">{{ stats.total_calls ? Math.round(stats.total_success / stats.total_calls * 100) : 0 }}%</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16">
      <el-col :span="16">
        <el-card>
          <template #header><span>{{ $t('auto.CostDashboard.7.d5bd40') }}</span></template>
          <div ref="trendRef" class="chart" />
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header><span>{{ $t('auto.CostDashboard.8.ba70ac') }}</span></template>
          <div ref="distRef" class="chart" />
        </el-card>
      </el-col>
    </el-row>

    <el-card class="breakdown-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.CostDashboard.9.f0067d') }}</span>
          <el-radio-group v-model="breakdown" size="small">
            <el-radio-button label="scenario">{{ $t('auto.CostDashboard.10.e1ec3b') }}</el-radio-button>
            <el-radio-button label="model">{{ $t('auto.CostDashboard.11.046f78') }}</el-radio-button>
            <el-radio-button label="agent">{{ $t('auto.CostDashboard.12.a40967') }}</el-radio-button>
          </el-radio-group>
        </div>
      </template>
      <el-table :data="breakdownData" v-loading="loading">
        <el-table-column prop="name" :label="breakdownName" />
        <el-table-column prop="call_count" :label="$t('auto.CostDashboard.1.00d5bd')" width="120" />
        <el-table-column prop="total_tokens" label="Token" width="120">
          <template #default="{ row }">{{ formatNumber(row.total_tokens) }}</template>
        </el-table-column>
        <el-table-column prop="total_cost" :label="$t('auto.CostDashboard.2.1f9507')" width="120">
          <template #default="{ row }">{{ Number(row.total_cost || 0).toFixed(4) }}</template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n';
import * as echarts from 'echarts'
import { safeInit } from '@/utils/echarts'
import { http } from '@/utils/request'
import { formatCompactNumber } from '@/utils/format'

const { t } = useI18n()

const stats = ref({})
const breakdown = ref('scenario')
const breakdownData = ref([])
const loading = ref(false)
const trendRef = ref()
const distRef = ref()

const breakdownName = computed(() => ({
  scenario: t('auto.CostDashboard.13.c93165'),
  model: t('auto.CostDashboard.14.8000f1'),
  agent: t('auto.CostDashboard.15.72fbb9')
}[breakdown.value]))

const formatNumber = formatCompactNumber

async function load() {
  loading.value = true
  try {
    const [cs, u] = await Promise.all([
      http.get('/api/llm/cost-stats', { window: 'month' }),
      http.get('/api/llm/usage', { window: 'month', group_by: breakdown.value })
    ])
    // 后端 cost-stats 返回 by_model/by_provider 分组、usage 返回汇总+by_* —— 前端做映射
    stats.value = u || {}
    const rows = (breakdown.value === 'model' ? (cs?.by_model || (u.by_provider || [])) : (u['by_' + breakdown.value] || []))
    breakdownData.value = rows.map((r) => ({ ...r, name: r.provider || r.scenario || r.agent || '-' }))
    renderTrend(breakdownData.value)
    renderDist(breakdownData.value)
  } finally {
    loading.value = false
  }
}

function renderTrend(rows) {
  if (!trendRef.value) return
  const top = (rows || []).slice(0, 8)
  if (!top.length) return
  const chart = safeInit(trendRef.value)
  chart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: top.map((d) => d.name) },
    yAxis: [
      { type: 'value', name: '¥' },
      { type: 'value', name: t('auto.CostDashboard.1.00d5bd'), position: 'right' }
    ],
    series: [
      { name: t('auto.CostDashboard.16.7e237e'), type: 'bar', data: top.map((d) => Number(d.total_cost || 0)) },
      { name: t('auto.CostDashboard.1.00d5bd'), type: 'bar', yAxisIndex: 1, data: top.map((d) => Number(d.call_count || 0)) }
    ]
  })
}

function renderDist(data) {
  if (!distRef.value) return
  const chart = safeInit(distRef.value)
  chart.setOption({
    tooltip: { trigger: 'item' },
    series: [{
      type: 'pie',
      radius: '70%',
      data: (data || []).slice(0, 8).map((d) => ({ name: d.name || '-', value: Number(d.total_cost || 0) }))
    }]
  })
}

watch(breakdown, load)
onMounted(load)
</script>

<style scoped>
.llm-cost { padding: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { text-align: center; padding: 8px; }
.stat-label { color: #64748B; font-size: 12px; }
.stat-value { font-size: 28px; font-weight: 700; margin: 8px 0; }
.chart { width: 100%; height: 320px; }
.breakdown-card { margin-top: 16px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
</style>
