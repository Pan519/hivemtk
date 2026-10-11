<template>
  <div class="deliverability">
    <el-row :gutter="16" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-label">{{ $t('auto.Deliverability.9.5b6ba0') }}</div>
          <div class="stat-value">{{ stats.deliveryRate?.toFixed(1) }}%</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card success">
          <div class="stat-label">{{ $t('auto.Deliverability.10.eb9eb1') }}</div>
          <div class="stat-value">{{ stats.openRate?.toFixed(1) }}%</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card warning">
          <div class="stat-label">{{ $t('auto.Deliverability.11.f30630') }}</div>
          <div class="stat-value">{{ stats.softBounce }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card danger">
          <div class="stat-label">{{ $t('auto.Deliverability.12.3df270') }}</div>
          <div class="stat-value">{{ stats.hardBounce }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.Deliverability.13.6690ec') }}</span>
          <el-button size="small" @click="refresh">{{ $t('auto.Deliverability.14.694fc5') }}</el-button>
        </div>
      </template>
      <el-table :data="domainReputation" v-loading="loading">
        <el-table-column prop="domain" :label="$t('auto.Deliverability.1.190980')" />
        <el-table-column prop="reputation" :label="$t('auto.Deliverability.2.9ed41e')" width="120">
          <template #default="{ row }">
            <el-rate v-model="row.reputation" disabled :max="5" show-score />
          </template>
        </el-table-column>
        <el-table-column prop="sentLast24h" :label="$t('auto.Deliverability.3.ba6b9f')" width="120" />
        <el-table-column prop="delivered" :label="$t('auto.Deliverability.4.f87f48')" width="100" />
        <el-table-column prop="bounced" :label="$t('auto.Deliverability.5.3f04ec')" width="80" />
        <el-table-column prop="complained" :label="$t('auto.Deliverability.6.e19d1d')" width="80" />
        <el-table-column prop="blacklisted" :label="$t('auto.Deliverability.7.bcbfe2')" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.blacklisted" type="danger" size="small">{{ $t('auto.Deliverability.15.23e66f') }}</el-tag>
            <el-tag v-else type="success" size="small">{{ $t('auto.Deliverability.16.fd6e80') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Deliverability.8.2b6bc0')" width="120" fixed="right">
          <template #default="{ row }">
            <el-button link type="warning" @click="suspendDomain(row)" :disabled="row.blacklisted">
              暂停使用
            </el-button>
            <el-button link type="primary" @click="viewDetails(row)">{{ $t('auto.Deliverability.17.f26225') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card class="bounce-card">
      <template #header><span>{{ $t('auto.Deliverability.18.870308') }}</span></template>
      <div ref="bounceChartRef" class="chart" />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n';
import { ElMessage, ElMessageBox } from 'element-plus'
import * as echarts from 'echarts'
import { safeInit } from '@/utils/echarts'
import { http } from '@/utils/request'

const { t } = useI18n()

const stats = ref({ deliveryRate: 0, openRate: 0, softBounce: 0, hardBounce: 0 })
const domainReputation = ref([])
const loading = ref(false)
const bounceChartRef = ref()

async function load() {
  loading.value = true
  try {
    const [s, domains, bounces] = await Promise.all([
      http.get('/api/email/deliverability'),
      http.get('/api/email/domain-reputation'),
      http.get('/api/email/bounces/breakdown')
    ])
    stats.value = s || stats.value
    domainReputation.value = domains || []
    renderBounceChart(bounces || [])
  } finally {
    loading.value = false
  }
}

function refresh() { load() }

function renderBounceChart(bounces) {
  if (!bounceChartRef.value) return
  const chart = safeInit(bounceChartRef.value)
  chart.setOption({
    tooltip: { trigger: 'item' },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      data: bounces.map((b) => ({ name: b.isp, value: b.count }))
    }]
  })
}

async function suspendDomain(row) {
  try {
    await ElMessageBox.confirm(`确认暂停域名「${row.domain}」？`, t('auto.Deliverability.19.9678f3'), { type: 'warning' })
  } catch (e) {
    if (e === 'cancel' || e === 'close') return
    throw e
  }
  await http.post(`/api/email/domains/${row.id}/suspend`, {})
  ElMessage.success(t('auto.Deliverability.20.a2d930'))
  await load()
}

function viewDetails(row) {
  window.open(`/email/domain-detail?id=${row.id}`, '_blank')
}

onMounted(load)
</script>

<style scoped>
.deliverability { padding: 16px; }
.stats-row { margin-bottom: 16px; }
.stat-card { text-align: center; padding: 8px; }
.stat-label { color: #64748B; font-size: 12px; }
.stat-value { font-size: 32px; font-weight: 700; margin: 8px 0; }
.stat-card.success .stat-value { color: #10B981; }
.stat-card.warning .stat-value { color: #F59E0B; }
.stat-card.danger .stat-value { color: #EF4444; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.bounce-card { margin-top: 16px; }
.chart { width: 100%; height: 320px; }
</style>
