<template>
  <div class="ops-overview">
    <PageHeader :title="$t('auto.Index.99d5464f')" :subtitle="$t('auto.Index.e1e8ef54')" />

    
    <el-row :gutter="20" class="status-row">
      <el-col :span="6">
        <el-card class="status-card status-ok">
          <div class="status-label">{{ $t('auto.Index.a838a981') }}</div>
          <div class="status-value">{{ systemStatus.text }}</div>
          <div class="status-meta">运行时长 {{ uptime }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="status-card">
          <div class="status-label">{{ $t('auto.Index.753e9e75') }}</div>
          <div class="status-value">{{ stats.onlineUsers }}</div>
          <div class="status-meta">总用户 {{ stats.totalUsers }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="status-card">
          <div class="status-label">{{ $t('auto.Index.3fa57183') }}</div>
          <div class="status-value">{{ stats.todayMessages }}</div>
          <div class="status-meta">环比 {{ stats.messagesDelta }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card :class="['status-card', alertsCount > 0 ? 'status-warn' : '']">
          <div class="status-label">{{ $t('auto.Index.d5aede72') }}</div>
          <div class="status-value">{{ alertsCount }}</div>
          <div class="status-meta">{{ $t('auto.Index.6e1e0000') }}</div>
        </el-card>
      </el-col>
    </el-row>

    
    <el-row :gutter="20" class="mt-20">
      <el-col :span="12">
        <el-card :header="$t('auto.Index.dd300fd9')">
          <el-table :data="moduleStatus" stripe>
            <el-table-column prop="name" :label="$t('auto.Index.fac54c34')" />
            <el-table-column prop="status" :label="$t('auto.Index.3fea7ca7')">
              <template #default="{ row }">
                <el-tag :type="row.status === 'ok' ? 'success' : row.status === 'warn' ? 'warning' : 'danger'">
                  {{ row.status === 'ok' ? '正常' : row.status === 'warn' ? '告警' : '异常' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="latency" :label="$t('auto.Index.93d0286b')" />
            <el-table-column prop="qps" label="QPS" />
            <el-table-column prop="errorRate" :label="$t('auto.Index.0c8524d7')" />
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card :header="$t('auto.Index.bbb49e4b')">
          <el-table :data="recentOps" stripe max-height="400">
            <el-table-column prop="time" :label="$t('auto.Index.19fcb9eb')" width="100" />
            <el-table-column prop="user" :label="$t('auto.Index.f9ac4b2a')" />
            <el-table-column prop="action" :label="$t('auto.Index.2b6bc0f2')" />
            <el-table-column prop="result" :label="$t('auto.Index.5ad7f5a8')">
              <template #default="{ row }">
                <el-tag :type="row.result === 'success' ? 'success' : 'danger'" size="small">
                  {{ row.result === 'success' ? '成功' : '失败' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import PageHeader from '@/components/PageHeader.vue'
import { http } from '@/utils/request'
import { getSystemStatus, getOpsStats, getModuleStatus, getRecentOperations } from '@/api/cockpit'

const { t } = useI18n()

const systemStatus = ref({ text: t('auto.Index.26b5bd49'), code: 'loading' })
const uptime = ref('0h 0m')
const stats = ref({ onlineUsers: 0, totalUsers: 0, todayMessages: 0, messagesDelta: '0%' })
const moduleStatus = ref([])
const recentOps = ref([])
const alertsCount = ref(0)

function formatUptime(sec) {
  if (!sec || sec < 0) return '0h 0m'
  const h = Math.floor(sec / 3600)
  const m = Math.floor((sec % 3600) / 60)
  return `${h}h ${m}m`
}

function mapModules(nodeHealth) {
  if (!Array.isArray(nodeHealth)) return []
  return nodeHealth.map(n => ({
    name: n.channel || n.node_id || n.platform || '-',
    status: n.status === 'healthy' || n.status === 'ok' ? 'ok'
      : n.status === 'degraded' || n.status === 'warn' ? 'warn' : 'error',
    latency: n.latency_ms ?? n.avg_latency_ms ?? '-',
    qps: n.qps ?? n.rate_per_min ?? '-',
    errorRate: n.error_rate ?? n.failure_rate ?? '-'
  }))
}

function mapRecentOps(anomalies) {
  if (!anomalies) return []
  const rows = []
  const groups = [
    { key: 'sync_gap', label: t('auto.Index.78b323b1') },
    { key: 'stuck_reachable', label: t('auto.Index.a2ea2848') },
    { key: 'stuck_unreachable', label: t('auto.Index.0fb8b251') },
    { key: 'unreachable', label: t('auto.Index.834f7cd8') },
    { key: 'node_abnormal', label: t('auto.Index.209a8c83') }
  ]
  for (const g of groups) {
    const items = anomalies[g.key] || []
    for (const it of items.slice(0, 10)) {
      rows.push({
        time: it.last_seen || it.detected_at || it.updated_at || '-',
        user: it.account_id || it.customer_id || '-',
        action: g.label,
        result: 'fail'
      })
    }
  }
  return rows.slice(0, 20)
}

async function loadOverview() {
  try {
    const [statusRes, opsStatsRes, modulesRes, opsRes] = await Promise.all([
      getSystemStatus(),
      getOpsStats(),
      getModuleStatus(),
      getRecentOperations(20)
    ])

    const s = statusRes.data || {}
    systemStatus.value = {
      text: s.status || (s.cpu_usage > 90 ? t('auto.Index.aa0eab9d') : t('auto.Index.fd6e80f1')),
      code: s.status || 'ok'
    }
    uptime.value = formatUptime(s.system_uptime ?? s.uptime)

    const o = opsStatsRes.data || {}
    stats.value = {
      onlineUsers: s.online_users ?? s.active_users ?? 0,
      totalUsers: s.total_users ?? 0,
      todayMessages: (o.inbound_rate_per_min * 60 + o.outbound_rate_per_min * 60) || 0,
      messagesDelta: '0%'
    }

    moduleStatus.value = mapModules(modulesRes.data || modulesRes.data?.nodes || [])
    recentOps.value = mapRecentOps(opsRes.data || {})

    try {
      const alertsRes = await http.get('/api/monitor/alerts/unread')
      alertsCount.value = alertsRes?.count || alertsRes?.unread_count || 0
    } catch {
      alertsCount.value = 0
    }
  } catch (err) {
    console.error('load overview failed:', err)
  }
}

let timer = null
onMounted(() => {
  loadOverview()
  timer = setInterval(loadOverview, 30000)
})
onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.ops-overview {
  padding: 20px;
}
.status-row { margin-top: 16px; }
.status-card {
  text-align: center;
  padding: 12px 0;
}
.status-card .status-label { font-size: 14px; color: #909399; }
.status-card .status-value { font-size: 32px; font-weight: 600; color: #303133; margin: 12px 0; }
.status-card .status-meta { font-size: 12px; color: #909399; }
.status-ok { border-left: 4px solid #67c23a; }
.status-warn { border-left: 4px solid #e6a23c; }
.mt-20 { margin-top: 20px; }
</style>
