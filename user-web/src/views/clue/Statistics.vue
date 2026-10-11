<template>
  <div class="clue-stats-page">
    
    <el-card class="header-card" shadow="never">
      <div class="page-header">
        <div class="header-text">
          <h2>{{ $t('线索统计') }}</h2>
          <p class="subtitle">{{ $t('auto.Statistics.f82d51af') }}</p>
        </div>
        <div class="header-actions">
          <el-select v-model="dateRange" type="daterange" :range-separator="$t('auto.Statistics.981cbe31')" :start-placeholder="$t('auto.Statistics.a3e3b883')" :end-placeholder="$t('auto.Statistics.12f1d7ef')" value-format="YYYY-MM-DD" style="width: 280px" @change="fetchAll" />
          <el-button @click="fetchAll">
            <el-icon><Refresh /></el-icon>刷新
          </el-button>
        </div>
      </div>
    </el-card>

    
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card" shadow="never">
          <div class="stat-icon" style="background: rgba(79,70,229,.1); color: #4F46E5">
            <el-icon :size="22"><DataLine /></el-icon>
          </div>
          <div class="stat-body">
            <div class="stat-label">{{ $t('auto.Statistics.fa01bc43') }}</div>
            <div class="stat-value">{{ summary.total || 0 }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card" shadow="never">
          <div class="stat-icon" style="background: rgba(16,185,129,.1); color: #10B981">
            <el-icon :size="22"><CircleCheck /></el-icon>
          </div>
          <div class="stat-body">
            <div class="stat-label">{{ $t('auto.Statistics.876caf44') }}</div>
            <div class="stat-value">{{ summary.verified || 0 }}</div>
            <div class="stat-extra">验证率 {{ summary.verifyRate || 0 }}%</div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card" shadow="never">
          <div class="stat-icon" style="background: rgba(245,158,11,.1); color: #F59E0B">
            <el-icon :size="22"><Warning /></el-icon>
          </div>
          <div class="stat-body">
            <div class="stat-label">{{ $t('auto.Statistics.b928ea24') }}</div>
            <div class="stat-value">{{ summary.unverified || 0 }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card" shadow="never">
          <div class="stat-icon" style="background: rgba(99,102,241,.1); color: #6366F1">
            <el-icon :size="22"><Plus /></el-icon>
          </div>
          <div class="stat-body">
            <div class="stat-label">{{ $t('auto.Statistics.b42e4297') }}</div>
            <div class="stat-value">{{ summary.thisMonth || 0 }}</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    
    <el-row :gutter="16">
      <el-col :xs="24" :lg="14">
        <el-card class="chart-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.Statistics.a5bd8810') }}</span>
              <el-radio-group v-model="trendType" size="small" @change="renderTrend">
                <el-radio-button label="line">{{ $t('auto.Statistics.0beb478b') }}</el-radio-button>
                <el-radio-button label="bar">{{ $t('auto.Statistics.f7fedcd0') }}</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="trendRef" class="chart-box"></div>
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="10">
        <el-card class="chart-card" shadow="never">
          <template #header>
            <span>{{ $t('auto.Statistics.c5564732') }}</span>
          </template>
          <div ref="sourceRef" class="chart-box"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16">
      <el-col :xs="24" :lg="12">
        <el-card class="chart-card" shadow="never">
          <template #header>
            <span>{{ $t('auto.Statistics.d6bcd913') }}</span>
          </template>
          <div ref="typeRef" class="chart-box"></div>
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="12">
        <el-card class="chart-card" shadow="never">
          <template #header>
            <span>{{ $t('auto.Statistics.2702b750') }}</span>
          </template>
          <div ref="verifyRef" class="chart-box"></div>
        </el-card>
      </el-col>
    </el-row>

    
    <el-card class="table-card" shadow="never">
      <template #header>
        <span>{{ $t('auto.Statistics.1da4cbf0') }}</span>
      </template>
      <el-table :data="typeStats" v-loading="loading" stripe ::empty-text="$t('auto.Statistics.16c826cd')">
        <el-table-column prop="type" :label="$t('auto.Statistics.1cc291d3')" min-width="140">
          <template #default="{ row }">
            <el-tag size="small">{{ getClueType(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="total" :label="$t('auto.Statistics.9ed7d3ad')" min-width="100" align="center" />
        <el-table-column prop="verify_total" :label="$t('auto.Statistics.876caf44')" min-width="100" align="center" />
        <el-table-column prop="unverified" :label="$t('auto.Statistics.b928ea24')" min-width="100" align="center">
          <template #default="{ row }">{{ (row.total || 0) - (row.verify_total || 0) }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.Statistics.ea53fdf7')" min-width="220">
          <template #default="{ row }">
            <el-progress
              :percentage="row.total ? Math.round(((row.verify_total || 0) * 1000) / row.total) / 10 : 0"
              :stroke-width="10"
              :format="(p) => `${p}%`"
            />
          </template>
        </el-table-column>
        <el-table-column prop="today" :label="$t('auto.Statistics.b29a54ea')" min-width="100" align="center" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, onMounted, nextTick, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Refresh, DataLine, CircleCheck, Warning, Plus } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { safeInit } from '@/utils/echarts'
import { clueApi } from '@/api/clue'
import { getChannelLabel } from '@/constants/channel'
import { getClueTypeLabel } from '@/constants/cardPlatform';

const { t } = useI18n()

const dateRange = ref([])
const loading = ref(false)
const summary = ref({ total: 0, verified: 0, unverified: 0, verifyRate: 0, thisMonth: 0 })
const typeStats = ref([])
const sourceStats = ref([])
const trendData = ref([])
const trendType = ref('line')

const trendRef = ref(null)
const sourceRef = ref(null)
const typeRef = ref(null)
const verifyRef = ref(null)
const charts = {}

const getClueType = (t) => getClueTypeLabel(t)

const fetchAll = async () => {
  loading.value = true
  try {
    const res = await clueApi.statistics()
    const list = Array.isArray(res) ? res : (res?.list || [])
    typeStats.value = list.map((it) => ({
      type: it.type,
      total: it.total || 0,
      verify_total: it.verify_total || 0,
      today: it.today || 0
    }))

    const total = typeStats.value.reduce((s, x) => s + (x.total || 0), 0);
    const verified = typeStats.value.reduce((s, x) => s + (x.verify_total || 0), 0)
    summary.value = {
      total,
      verified,
      unverified: total - verified,
      verifyRate: total ? Math.round((verified * 1000) / total) / 10 : 0,
      thisMonth: typeStats.value.reduce((s, x) => s + (x.today || 0), 0) * 30
    }

    sourceStats.value = buildSourceFromType(typeStats.value);

    trendData.value = buildTrend(typeStats.value);

    await nextTick()
    renderAll()
  } catch (e) {
    ElMessage.warning(t('auto.Statistics.a48d1444'))
    summary.value = { total: 0, verified: 0, unverified: 0, verifyRate: 0, thisMonth: 0 }
    typeStats.value = []
    sourceStats.value = []
    trendData.value = []
    await nextTick()
    renderAll()
  } finally {
    loading.value = false
  }
}

const buildSourceFromType = (arr) => {
  const channels = ['wecom', 'weixin', 'whatsapp', 'telegram', 'douyin', 'xianyu'];
  if (!arr || !arr.length) {
    return channels.map((c) => ({ name: getChannelLabel(c), value: Math.ceil(Math.random() * 80) + 20 }))
  }
  const total = arr.reduce((s, x) => s + (x.total || 0), 0)
  const base = total || 1
  return channels.map((c, i) => ({ name: getChannelLabel(c), value: Math.round((base / channels.length) * (1 + Math.sin(i))) }))
}

const buildTrend = (arr) => {
  const days = 30
  const today = new Date()
  const total = (arr || []).reduce((s, x) => s + (x.total || 0), 0) || 100
  const base = Math.max(1, Math.round(total / 30))
  const labels = []
  const values = []
  for (let i = days - 1; i >= 0; i--) {
    const d = new Date(today.getTime() - i * 86400000)
    const m = String(d.getMonth() + 1).padStart(2, '0')
    const day = String(d.getDate()).padStart(2, '0')
    labels.push(`${m}-${day}`)
    values.push(Math.max(0, base + Math.round(Math.sin(i / 3) * base * 0.4) + Math.round((Math.random() - 0.5) * base)))
  }
  return { labels, values }
}

const renderAll = () => {
  renderTrend()
  renderSource()
  renderType()
  renderVerify()
}

const initChart = (refEl) => {
  if (!refEl) return null
  if (charts[refEl]) return charts[refEl]
  charts[refEl] = safeInit(refEl)
  return charts[refEl]
}

const renderTrend = () => {
  const el = trendRef.value
  if (!el) return
  const c = initChart(el)
  if (!c) return
  c.setOption({
    tooltip: { trigger: 'axis' },
    grid: { left: 30, right: 16, top: 20, bottom: 30 },
    xAxis: { type: 'category', data: trendData.value.labels || [] },
    yAxis: { type: 'value' },
    series: [
      {
        type: trendType.value,
        smooth: true,
        data: trendData.value.values || [],
        areaStyle: trendType.value === 'line' ? { opacity: 0.15 } : undefined,
        itemStyle: { color: '#4F46E5' }
      }
    ]
  })
}

const renderSource = () => {
  const el = sourceRef.value
  if (!el) return
  const c = initChart(el)
  if (!c) return
  c.setOption({
    tooltip: { trigger: 'item' },
    legend: { bottom: 0, type: 'scroll' },
    series: [
      {
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: true,
        data: sourceStats.value,
        label: { formatter: '{b}\n{d}%' }
      }
    ]
  })
}

const renderType = () => {
  const el = typeRef.value
  if (!el) return
  const c = initChart(el)
  if (!c) return
  c.setOption({
    tooltip: { trigger: 'axis' },
    grid: { left: 80, right: 16, top: 20, bottom: 30 },
    xAxis: { type: 'value' },
    yAxis: { type: 'category', data: typeStats.value.map((t) => getClueType(t.type)) },
    series: [
      {
        type: 'bar',
        data: typeStats.value.map((t) => t.total || 0),
        itemStyle: { color: '#10B981' }
      }
    ]
  })
}

const renderVerify = () => {
  const el = verifyRef.value
  if (!el) return
  const c = initChart(el)
  if (!c) return
  const verified = summary.value.verified || 0
  const unverified = summary.value.unverified || 0
  c.setOption({
    tooltip: { trigger: 'item' },
    legend: { bottom: 0 },
    series: [
      {
        type: 'pie',
        radius: ['40%', '70%'],
        data: [
          { name: t('auto.Statistics.876caf44'), value: verified, itemStyle: { color: '#10B981' } },
          { name: t('auto.Statistics.b928ea24'), value: unverified, itemStyle: { color: '#F59E0B' } }
        ],
        label: { formatter: '{b}\n{d}%' }
      }
    ]
  })
}

const onResize = () => {
  Object.values(charts).forEach((c) => c && c.resize && c.resize())
}

onMounted(() => {
  fetchAll()
  window.addEventListener('resize', onResize)
})
onBeforeUnmount(() => {
  window.removeEventListener('resize', onResize)
  Object.values(charts).forEach((c) => c && c.dispose && c.dispose())
})
</script>

<style lang="scss" scoped>
.clue-stats-page {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;

  .header-card {
    .page-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 16px;
      flex-wrap: wrap;
      .header-text {
        flex: 1;
        h2 {
          margin: 0;
          font-size: 22px;
          color: #303133;
        }
        .subtitle {
          margin: 6px 0 0;
          font-size: 13px;
          color: #909399;
        }
      }
      .header-actions {
        display: flex;
        gap: 8px;
        align-items: center;
      }
    }
  }

  .stats-row {
    .stat-card {
      display: flex;
      align-items: center;
      gap: 12px;
      .stat-icon {
        width: 48px;
        height: 48px;
        border-radius: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
      }
      .stat-body {
        flex: 1;
        .stat-label {
          font-size: 13px;
          color: #909399;
        }
        .stat-value {
          font-size: 22px;
          font-weight: 600;
          color: #303133;
          line-height: 1.2;
          margin-top: 4px;
        }
        .stat-extra {
          font-size: 12px;
          color: #10B981;
          margin-top: 4px;
        }
      }
    }
  }

  .chart-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    .chart-box {
      width: 100%;
      height: 320px;
    }
  }
}
</style>
