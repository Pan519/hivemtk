<template>
  <div class="customer-event-page">
    <el-card class="header-card">
      <div>
        <h2>{{ $t('客户事件追踪') }}</h2>
        <p class="subtitle">{{ $t('追踪客户在系统中的所有行为事件') }}</p>
      </div>
      <el-button type="primary" @click="showCreateDialog">
        <el-icon><Plus /></el-icon>
        {{ $t('创建事件') }}
      </el-button>
    </el-card>

    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('总事件数') }}</div>
            <div class="stat-value">{{ stats.total }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('事件类型') }}</div>
            <div class="stat-value">{{ stats.typeCount }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('事件来源') }}</div>
            <div class="stat-value">{{ stats.sourceCount }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('事件总量') }}</div>
            <div class="stat-value">{{ stats.eventTypes }}</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('事件列表') }}</span>
          <div>
            <el-select v-model="filterType" ::placeholder="$t('auto.List.c7a0e919')" clearable style="width: 150px; margin-right: 10px">
              <el-option ::label="$t('auto.List.fe97625f')" value="view" />
              <el-option ::label="$t('auto.List.6764dc5e')" value="click" />
              <el-option ::label="$t('auto.List.45fdc386')" value="register" />
              <el-option ::label="$t('auto.List.703a2964')" value="purchase" />
              <el-option ::label="$t('auto.List.e78edce8')" value="share" />
            </el-select>
            <el-date-picker v-model="dateRange" type="daterange" :range-separator="$t('auto.List.981cbe31')" :start-placeholder="$t('auto.List.a3e3b883')" :end-placeholder="$t('auto.List.12f1d7ef')" value-format="YYYY-MM-DD" style="width: 240px" />
          </div>
        </div>
      </template>
      <el-table :data="filteredEvents" v-loading="loading" stripe>
        <el-table-column prop="eventType" :label="$t('auto.List.f974c846')" width="120">
          <template #default="{ row }">
            <el-tag :type="getEventTypeTag(row.eventType)">{{ row.eventType }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="userId" :label="$t('auto.List.30acd20d')" width="100" />
        <el-table-column prop="userName" :label="$t('auto.List.1fd02a90')" width="120" />
        <el-table-column prop="action" :label="$t('auto.List.a0496123')" min-width="150" />
        <el-table-column prop="target" :label="$t('auto.List.49e81f5e')" min-width="150" />
        <el-table-column prop="source" :label="$t('auto.List.26ca20b1')" width="100" />
        <el-table-column prop="ip" label="IP" width="130" />
        <el-table-column prop="createdAt" :label="$t('auto.List.19fcb9eb')" width="180" />
        <el-table-column :label="$t('auto.List.2b6bc0f2')" width="120" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="viewEvent(row)">{{ $t('auto.List.f26225bd') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.size"
        :total="pagination.total"
        layout="total, prev, pager, next, jumper"
        @current-change="loadEvents"
        style="margin-top: 15px; text-align: right"
      />
    </el-card>

    <el-dialog v-model="dialogVisible" :title="$t('auto.List.9c0ca814')" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item :label="$t('auto.List.52f15cfa')" required>
          <el-select v-model="form.customerId" filterable :placeholder="$t('auto.List.ec3fb565')" style="width: 100%">
            <el-option v-for="c in customerOptions" :key="c.id" :label="(c.name || c.id) + (c.phone ? ' (' + c.phone + ')' : '')" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.e84aba2a')">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.f974c846')">
          <el-select v-model="form.eventType" style="width: 100%">
            <el-option :label="$t('auto.List.9c5c5cdb')" value="view" />
            <el-option :label="$t('auto.List.4363c17e')" value="click" />
            <el-option :label="$t('auto.List.f1d4ff50')" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.65939ccd')">
          <el-input v-model="form.trigger" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.List.625fb26b') }}</el-button>
        <el-button type="primary" @click="submitForm">{{ $t('auto.List.d9ac9228') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="detailDialogVisible" :title="$t('auto.List.715f8295')" width="650px" :close-on-click-modal="false">
      <template v-if="currentEvent">
        <el-descriptions :column="2" border>
          <el-descriptions-item :label="$t('auto.List.f974c846')">
            <el-tag :type="getEventTypeTag(currentEvent.eventType)">{{ currentEvent.eventType }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.a0496123')">{{ currentEvent.action }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.30acd20d')">{{ currentEvent.userId }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.819767ad')">{{ currentEvent.userName }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.49e81f5e')">{{ currentEvent.target }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.26ca20b1')">{{ currentEvent.source }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.2dc9105c')">{{ currentEvent.ip }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.19fcb9eb')">{{ currentEvent.createdAt }}</el-descriptions-item>
        </el-descriptions>

        <el-divider content-position="left">{{ $t('auto.List.a9233d6d') }}</el-divider>
        <div class="event-data-json">
          <pre>{{ JSON.stringify(currentEvent, null, 2) }}</pre>
        </div>

        <template v-if="currentEvent.properties">
          <el-divider content-position="left">{{ $t('auto.List.33aeeb47') }}</el-divider>
          <el-descriptions :column="2" border>
            <el-descriptions-item
              v-for="(value, key) in currentEvent.properties"
              :key="key"
              :label="String(key)"
            >
              {{ value }}
            </el-descriptions-item>
          </el-descriptions>
        </template>
      </template>
      <template #footer>
        <el-button type="primary" @click="detailDialogVisible = false">{{ $t('auto.List.b15d9127') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getCustomerEventHistory, getEventStats, trackEvent, getGlobalEvents } from '@/api/customerEvent.js'
import { getCustomerList } from '@/api/customer360.js'

const { t } = useI18n()

const loading = ref(false)
const events = ref([])
const filterType = ref('')
const dateRange = ref([])
const dialogVisible = ref(false)
const form = ref({ name: '', eventType: 'view', trigger: '', userScope: 'all' })
const stats = ref({ total: 0, typeCount: 0, sourceCount: 0, eventTypes: 0 })
const pagination = ref({ page: 1, size: 20, total: 0 })
const detailDialogVisible = ref(false)
const currentEvent = ref(null)

const filteredEvents = computed(() => {
  let result = events.value
  if (filterType.value) result = result.filter(e => e.eventType === filterType.value)
  return result
})

const getEventTypeTag = (type) => {
  const map = { view: '', click: 'primary', register: 'success', purchase: 'warning', share: 'info' }
  return map[type]}

const mapEvent = (e) => {
  let data = {}
  try {
    data = typeof e.event_data === 'string' ? JSON.parse(e.event_data) : (e.event_data || {})
  } catch (_) {
    // JSON 解析失败时维持初始空对象
  }
  return {
    id: e.id,
    eventType: e.event_type || '',
    userId: e.customer_id || '',
    userName: data.user_name || data.name || e.customer_id || '',
    action: data.action || e.event_type || '',
    target: data.target || data.url || data.page || '',
    source: e.event_source || '',
    ip: data.ip || '',
    createdAt: (e.occurred_at || e.created_at || '').toString().replace('T', ' ').slice(0, 19),
    properties: data
  }
};

const loadEvents = async () => {
  loading.value = true
  try {
    const res = await getGlobalEvents({ page: pagination.value.page || 1, page_size: pagination.value.pageSize || 50 });
    const raw = Array.isArray(res) ? res : (res?.list || [])
    events.value = raw.map(mapEvent)
    pagination.value.total = res?.total ?? raw.length
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    const res = await getEventStats()
    const byType = res?.by_event_type || []
    const bySource = res?.by_event_source || []
    stats.value = {
      total: res?.total_events || 0,
      typeCount: byType.length,
      sourceCount: bySource.length,
      eventTypes: byType.reduce((s, x) => s + (x.count || 0), 0)
    }
  } catch (e) {
    console.error(t('auto.List.8522b582'), e)
  }
}

const customerOptions = ref([])
const loadCustomerOptions = async () => {
  try {
    const res = await getCustomerList({ pageSize: 200 })
    customerOptions.value = res?.list || []
  } catch (e) {
    console.error(t('auto.List.fbc9e953'), e)
  }
}

const showCreateDialog = () => {
  form.value = { customerId: '', name: '', eventType: 'view', trigger: '' }
  dialogVisible.value = true
  loadCustomerOptions()
}

const submitForm = async () => {
  if (!form.value.customerId) {
    ElMessage.warning(t('auto.List.6d00ea81'))
    return
  }
  await trackEvent({
    customer_id: form.value.customerId,
    event_type: form.value.eventType,
    event_source: 'manual',
    event_data: { name: form.value.name, trigger: form.value.trigger }
  })
  ElMessage.success(i18n.global.t(t('auto.List.39333ccb')))
  dialogVisible.value = false
  loadEvents()
}

const viewEvent = (row) => {
  currentEvent.value = row
  detailDialogVisible.value = true
}

onMounted(() => {
  loadEvents()
  loadStats()
})
</script>

<style scoped lang="scss">
.customer-event-page { padding: 20px; }
.header-card {
  margin-bottom: 20px;
  :deep(.el-card__body) { display: flex; justify-content: space-between; align-items: center; }
  h2 { margin: 0 0 8px 0; }
  .subtitle { color: #909399; margin: 0; }
}
.stats-row {
  margin-bottom: 20px;
  .stat-item {
    text-align: center;
    .stat-label { color: #909399; font-size: 14px; margin-bottom: 10px; }
    .stat-value { font-size: 28px; font-weight: bold; }
  }
}
.card-header { display: flex; justify-content: space-between; align-items: center; }
.event-data-json {
  background: #f5f7fa;
  border-radius: 4px;
  padding: 12px;
  max-height: 300px;
  overflow-y: auto;
}
.event-data-json pre {
  margin: 0;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
