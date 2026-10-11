<template>
  <div class="wecom-data">
    <el-card shadow="never">
      <div class="toolbar">
        <el-select
          v-model="selectedAccountId"
          :placeholder="$t('auto.Data.1.915fcb')"
          clearable
          filterable
          style="width: 240px"
          @change="onAccountChange"
        >
          <el-option
            v-for="a in accounts"
            :key="a.account.id"
            ::label="$t('auto.Data.2.a5e6e3')"
            :value="a.account.id"
          />
        </el-select>
        <el-button type="primary" :icon="Refresh" :loading="loadingAny" @click="refreshAll">{{ $t('auto.Data.26.498519') }}</el-button>
        <span class="tip">{{ $t('auto.Data.27.49507c') }}</span>
      </div>

      
      <div class="summary-row">
        <div class="summary-card">
          <div class="summary-value">{{ total.customers || 0 }}</div>
          <div class="summary-label">{{ $t('auto.Data.28.a4577f') }}</div>
        </div>
        <div class="summary-card">
          <div class="summary-value">{{ total.groups || 0 }}</div>
          <div class="summary-label">{{ $t('auto.Data.29.ce77df') }}</div>
        </div>
        <div class="summary-card">
          <div class="summary-value">{{ total.messages || 0 }}</div>
          <div class="summary-label">{{ $t('auto.Data.30.22825d') }}</div>
        </div>
      </div>

      
      <div class="chart-grid">
        <el-card shadow="hover" class="chart-card">
          <div ref="sourceChart" class="chart-box"></div>
        </el-card>
        <el-card shadow="hover" class="chart-card">
          <div ref="groupSizeChart" class="chart-box"></div>
        </el-card>
        <el-card shadow="hover" class="chart-card chart-card-wide">
          <div ref="msgTrendChart" class="chart-box chart-box-wide"></div>
        </el-card>
      </div>

      <el-tabs v-model="activeTab" @tab-change="onTabChange">
        
        <el-tab-pane :label="$t('auto.Data.3.ff0b20')" name="customers">
          <el-table :data="customers" v-loading="loading.customers" border stripe size="small">
            <el-table-column prop="avatar" :label="$t('auto.Data.4.4c50ee')" width="60">
              <template #default="{ row }">
                <el-avatar v-if="row.avatar" :src="row.avatar" size="small" />
                <span v-else>{{ row.name?.charAt(0) || '?' }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="name" :label="$t('auto.Data.5.d7ec2d')" min-width="120" />
            <el-table-column :label="$t('auto.Data.6.226b09')" width="90">
              <template #default="{ row }">{{ row.type === 1 ? '企业' : '个人' }}</template>
            </el-table-column>
            <el-table-column prop="corp_name" :label="$t('auto.Data.7.04c9e3')" min-width="120" />
            <el-table-column prop="position" :label="$t('auto.Data.8.b652de')" min-width="100" />
            <el-table-column prop="external_userid" :label="$t('auto.Data.9.941aeb')" min-width="160" show-overflow-tooltip />
            <el-table-column prop="remark" :label="$t('auto.Data.10.2432b5')" min-width="100" show-overflow-tooltip />
            <el-table-column :label="$t('auto.Data.11.2b6bc0')" width="100" fixed="right">
              <template #default="{ row }">
                <el-button size="small" type="primary" @click="openSendTo(row.external_userid)">{{ $t('auto.Data.31.94e705') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            class="pager"
            background
            layout="total, prev, pager, next"
            :total="total.customers"
            :current-page="page.customers"
            :page-size="pageSize"
            @current-change="(p) => { page.customers = p; loadCustomers() }"
          />
        </el-tab-pane>

        
        <el-tab-pane :label="$t('auto.Data.12.8ba76f')" name="groups">
          <el-table :data="groups" v-loading="loading.groups" border stripe size="small">
            <el-table-column prop="name" :label="$t('auto.Data.13.b65ab8')" min-width="160" />
            <el-table-column prop="owner" :label="$t('auto.Data.14.4ef262')" min-width="120" />
            <el-table-column prop="member_count" :label="$t('auto.Data.15.eb9d06')" width="90" />
            <el-table-column prop="chat_id" :label="$t('auto.Data.16.3e0e8a')" min-width="160" show-overflow-tooltip />
            <el-table-column :label="$t('auto.Data.17.eca37c')" width="160">
              <template #default="{ row }">{{ row.created_at ? formatTime(row.created_at) : '-' }}</template>
            </el-table-column>
          </el-table>
          <el-pagination
            class="pager"
            background
            layout="total, prev, pager, next"
            :total="total.groups"
            :current-page="page.groups"
            :page-size="pageSize"
            @current-change="(p) => { page.groups = p; loadGroups() }"
          />
        </el-tab-pane>

        
        <el-tab-pane :label="$t('auto.Data.18.14d342')" name="tags">
          <el-table :data="tags" v-loading="loading.tags" border stripe size="small">
            <el-table-column prop="tag_name" :label="$t('auto.Data.19.0718bf')" min-width="160" />
            <el-table-column prop="tag_id" :label="$t('auto.Data.20.e2c5cd')" min-width="140" />
          </el-table>
        </el-tab-pane>

        
        <el-tab-pane :label="$t('auto.Data.21.046826')" name="messages">
          <el-table :data="messages" v-loading="loading.messages" border stripe size="small">
            <el-table-column prop="to_user" :label="$t('auto.Data.22.8805bb')" min-width="160" show-overflow-tooltip />
            <el-table-column prop="msg_type" :label="$t('auto.Data.6.226b09')" width="90" />
            <el-table-column :label="$t('auto.Data.23.2d711b')" min-width="200" show-overflow-tooltip>
              <template #default="{ row }">
                {{ previewContent(row) }}
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.Data.24.3fea7c')" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
                  {{ row.status === 1 ? '成功' : '失败' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.Data.25.63b34d')" width="160">
              <template #default="{ row }">{{ row.created_at ? formatTime(row.created_at) : '-' }}</template>
            </el-table-column>
          </el-table>
          <el-pagination
            class="pager"
            background
            layout="total, prev, pager, next"
            :total="total.messages"
            :current-page="page.messages"
            :page-size="pageSize"
            @current-change="(p) => { page.messages = p; loadMessages() }"
          />
        </el-tab-pane>
      </el-tabs>
    </el-card>

    
    <WeComSendDialog v-model:visible="sendVisible" :account-id="selectedAccountId" :external-userid="sendExternalUserID" />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import { safeInit } from '@/utils/echarts'
import { wecomAccountApi } from '@/api/wecomAccount'
import WeComSendDialog from '@/components/WeComSendDialog.vue'
import { toList } from '@/utils/list'

const { t } = useI18n()

const activeTab = ref('customers')
const pageSize = 20
const chartPageSize = 100

const accounts = ref([])
const selectedAccountId = ref(null)

const customers = ref([])
const groups = ref([])
const tags = ref([])
const messages = ref([])

const total = reactive({ customers: 0, groups: 0, messages: 0 })
const page = reactive({ customers: 1, groups: 1, messages: 1 })

const loading = reactive({ customers: false, groups: false, tags: false, messages: false, charts: false })
const loadingAny = ref(false)

const sendVisible = ref(false)
const sendExternalUserID = ref('')

const sourceChart = ref(null);
const groupSizeChart = ref(null)
const msgTrendChart = ref(null)
let sourceChartInst = null
let groupSizeChartInst = null
let msgTrendChartInst = null

const formatTime = (v) => {
  if (!v) return '-'
  const d = new Date(v)
  if (isNaN(d.getTime())) return v
  const p = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}`
}

const dayKey = (v) => {
  const d = new Date(v)
  if (isNaN(d.getTime())) return null
  const p = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())}`
}

const previewContent = (row) => {
  if (row.msg_type === 'text') return row.content
  if (row.msg_type === 'image') return '[图片] ' + (row.media_id || '')
  if (row.msg_type === 'link') return '[链接] ' + (row.title || row.url || '')
  return row.content || '-'
}

const loadAccounts = async () => {
  try {
    const res = await wecomAccountApi.listAccounts()
    accounts.value = toList(res)
    if (accounts.value.length && !selectedAccountId.value) {
      selectedAccountId.value = accounts.value[0].account.id
    }
  } catch (e) {
    ElMessage.error(t('auto.Data.32.51fceb') + (e.message || e))
  }
}

const loadCustomers = async () => {
  loading.customers = true
  try {
    const res = await wecomAccountApi.getCustomers({ page: page.customers, page_size: pageSize })
    customers.value = toList(res.list)
    total.customers = res.total || 0
  } catch (e) {
    ElMessage.error(t('auto.Data.33.40bf9d') + (e.message || e))
  } finally {
    loading.customers = false
  }
}

const loadGroups = async () => {
  loading.groups = true
  try {
    const res = await wecomAccountApi.getGroups({ page: page.groups, page_size: pageSize })
    groups.value = toList(res.list)
    total.groups = res.total || 0
  } catch (e) {
    ElMessage.error(t('auto.Data.34.523fac') + (e.message || e))
  } finally {
    loading.groups = false
  }
}

const loadTags = async () => {
  loading.tags = true
  try {
    const res = await wecomAccountApi.getTags()
    tags.value = toList(res)
  } catch (e) {
    ElMessage.error(t('auto.Data.35.d1e871') + (e.message || e))
  } finally {
    loading.tags = false
  }
}

const loadMessages = async () => {
  loading.messages = true
  try {
    const res = await wecomAccountApi.getMessages({ page: page.messages, page_size: pageSize })
    messages.value = toList(res.list)
    total.messages = res.total || 0
  } catch (e) {
    ElMessage.error(t('auto.Data.36.68a249') + (e.message || e))
  } finally {
    loading.messages = false
  }
}

const loadChartData = async () => {
  loading.charts = true
  try {
    const [cRes, gRes, mRes] = await Promise.all([
      wecomAccountApi.getCustomers({ page: 1, page_size: chartPageSize }),
      wecomAccountApi.getGroups({ page: 1, page_size: chartPageSize }),
      wecomAccountApi.getMessages({ page: 1, page_size: chartPageSize })
    ])
    const allCustomers = toList(cRes.list)
    const allGroups = toList(gRes.list)
    const allMessages = toList(mRes.list)

    const sourceMap = {};
    for (const c of allCustomers) {
      const k = c.source || t('auto.Data.37.1622dc')
      sourceMap[k] = (sourceMap[k] || 0) + 1
    }
    const sourceData = Object.keys(sourceMap).map((k) => ({ name: k, value: sourceMap[k] }))

    const topGroups = [...allGroups]
      .sort((a, b) => (Number(b.member_count) || 0) - (Number(a.member_count) || 0))
      .slice(0, 10);
    const groupNames = topGroups.map((g) => g.name || g.chat_id || t('auto.Data.38.532575'))
    const groupSizes = topGroups.map((g) => Number(g.member_count) || 0)

    const days = [];
    const now = new Date()
    for (let i = 29; i >= 0; i--) {
      const d = new Date(now)
      d.setDate(now.getDate() - i)
      const p = (n) => String(n).padStart(2, '0')
      days.push(`${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())}`)
    }
    const trendMap = {}
    for (const d of days) trendMap[d] = 0
    for (const m of allMessages) {
      const k = dayKey(m.created_at)
      if (k && trendMap[k] !== undefined) trendMap[k] += 1
    }
    const trendData = days.map((d) => trendMap[d])

    nextTick(() => {
      renderSourceChart(sourceData)
      renderGroupSizeChart(groupNames, groupSizes)
      renderMsgTrendChart(days, trendData)
    })
  } catch (e) {
    ElMessage.error(t('auto.Data.39.ab4618') + (e.message || e))
  } finally {
    loading.charts = false
  }
};

const renderSourceChart = (data) => {
  if (sourceChartInst) sourceChartInst.dispose()
  if (!sourceChart.value) return
  sourceChartInst = safeInit(sourceChart.value)
  sourceChartInst.setOption({
    title: { text: t('auto.Data.40.e8f6bb'), left: 'center' },
    tooltip: { trigger: 'item' },
    legend: { bottom: 0, type: 'scroll' },
    series: [
      {
        name: t('auto.Data.41.d85de8'),
        type: 'pie',
        radius: '55%',
        center: ['50%', '45%'],
        data: data.length ? data : [{ name: t('auto.Data.42.21efd8'), value: 1 }],
        label: { formatter: '{b}: {c}' }
      }
    ]
  })
}

const renderGroupSizeChart = (names, sizes) => {
  if (groupSizeChartInst) groupSizeChartInst.dispose()
  if (!groupSizeChart.value) return
  groupSizeChartInst = safeInit(groupSizeChart.value)
  groupSizeChartInst.setOption({
    title: { text: t('auto.Data.43.f0ffef'), left: 'center' },
    tooltip: { trigger: 'axis' },
    grid: { left: 80, right: 20, top: 50, bottom: 30 },
    xAxis: { type: 'value' },
    yAxis: { type: 'category', data: names.slice().reverse() },
    series: [
      {
        name: t('auto.Data.15.eb9d06'),
        type: 'bar',
        data: sizes.slice().reverse(),
        itemStyle: { color: '#409EFF' }
      }
    ]
  })
}

const renderMsgTrendChart = (days, data) => {
  if (msgTrendChartInst) msgTrendChartInst.dispose()
  if (!msgTrendChart.value) return
  msgTrendChartInst = safeInit(msgTrendChart.value)
  msgTrendChartInst.setOption({
    title: { text: t('auto.Data.44.59a763'), left: 'center' },
    tooltip: { trigger: 'axis' },
    grid: { left: 50, right: 20, top: 50, bottom: 50 },
    xAxis: { type: 'category', data: days, axisLabel: { rotate: 45, interval: 3 } },
    yAxis: { type: 'value' },
    series: [
      {
        name: t('auto.Data.45.342cea'),
        type: 'line',
        smooth: true,
        data: data,
        areaStyle: { opacity: 0.15 },
        itemStyle: { color: '#67C23A' }
      }
    ]
  })
}

const resizeCharts = () => {
  if (sourceChartInst) sourceChartInst.resize()
  if (groupSizeChartInst) groupSizeChartInst.resize()
  if (msgTrendChartInst) msgTrendChartInst.resize()
}

const onTabChange = (tab) => {
  if (tab === 'customers' && !customers.value.length) loadCustomers()
  if (tab === 'groups' && !groups.value.length) loadGroups()
  if (tab === 'tags' && !tags.value.length) loadTags()
  if (tab === 'messages' && !messages.value.length) loadMessages()
}

const onAccountChange = () => {
  refreshAll();
}

const openSendTo = (externalUserid) => {
  if (!selectedAccountId.value) {
    ElMessage.warning(t('auto.Data.46.dcf5bd'))
    return
  }
  sendExternalUserID.value = externalUserid
  sendVisible.value = true
}

const refreshAll = () => {
  loadingAny.value = true
  Promise.all([
    loadCustomers(),
    loadGroups(),
    loadMessages(),
    loadChartData()
  ]).finally(() => { loadingAny.value = false })
}

onMounted(() => {
  loadAccounts()
  loadCustomers()
  loadGroups()
  loadMessages()
  loadChartData()
  window.addEventListener('resize', resizeCharts)
})

onUnmounted(() => {
  window.removeEventListener('resize', resizeCharts)
  if (sourceChartInst) sourceChartInst.dispose()
  if (groupSizeChartInst) groupSizeChartInst.dispose()
  if (msgTrendChartInst) msgTrendChartInst.dispose()
})
</script>

<style scoped>
.wecom-data {
  padding: 4px;
}
.toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}
.tip {
  color: #909399;
  font-size: 12px;
}
.summary-row {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}
.summary-card {
  flex: 1;
  min-width: 160px;
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
  color: #fff;
  border-radius: 10px;
  padding: 18px 20px;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.25);
}
.summary-value {
  font-size: 28px;
  font-weight: 700;
  line-height: 1.2;
}
.summary-label {
  font-size: 13px;
  opacity: 0.9;
  margin-top: 4px;
}
.chart-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 16px;
}
.chart-card-wide {
  grid-column: 1 / 3;
}
.chart-box {
  width: 100%;
  height: 300px;
}
.chart-box-wide {
  height: 320px;
}
.pager {
  margin-top: 12px;
  justify-content: flex-end;
}
</style>
