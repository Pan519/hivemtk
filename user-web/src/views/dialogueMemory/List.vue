<template>
  <div class="dialogue-memory-page">
    
    <el-card class="header-card" shadow="never">
      <div class="header-content">
        <div>
          <h2>{{ $t('对话记忆') }}</h2>
          <p class="subtitle">{{ $t('auto.List.23.cb3995') }}</p>
        </div>
        <div class="header-actions">
          <el-button @click="refreshAll" :loading="loading">
            <el-icon><Refresh /></el-icon>
            {{ $t('刷新') }}
          </el-button>
        </div>
      </div>
    </el-card>

    
    <el-row :gutter="16" class="stats-row">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-value">{{ pagination.total }}</div>
          <div class="stat-label">{{ $t('客户记忆数') }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card stat-high">
          <div class="stat-value">{{ intentCount('high') }}</div>
          <div class="stat-label">{{ $t('高意向客户') }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card stat-medium">
          <div class="stat-value">{{ intentCount('medium') }}</div>
          <div class="stat-label">{{ $t('中意向客户') }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card stat-objection">
          <div class="stat-value">{{ totalObjections }}</div>
          <div class="stat-label">{{ $t('累计异议记录') }}</div>
        </el-card>
      </el-col>
    </el-row>

    
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('记忆统计列表') }}</span>
          <div class="header-controls">
            <el-input
              v-model="filterCustomerID"
              ::placeholder="$t('auto.List.1.ca3542')"
              clearable
              size="small"
              style="width: 200px"
              @keyup.enter="onSearch"
              @clear="onSearch"
            />
            <el-button size="small" type="primary" @click="onSearch">{{ $t('查询') }}</el-button>
          </div>
        </div>
      </template>
      <el-table :data="memoryList" v-loading="listLoading" stripe>
        <el-table-column :label="$t('auto.List.2.ff0b20')" min-width="160">
          <template #default="{ row }">
            <div class="customer-cell">
              <span class="customer-name">{{ row.customer_name || row.customer_id || '-' }}</span>
              <span class="customer-id">ID: {{ row.customer_id || '-' }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.3.ce092e')" width="100">
          <template #default="{ row }">
            <el-tag :type="getIntentType(row.intent_level || row.purchase_intent)" size="small">
              {{ getIntentLabel(row.intent_level || row.purchase_intent) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.4.342cea')" width="90" align="center">
          <template #default="{ row }">{{ row.message_count || 0 }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.5.19b1dc')" width="90" align="center">
          <template #default="{ row }">{{ (row.objections || []).length }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.6.507e6a')" width="90" align="center">
          <template #default="{ row }">
            <el-tag size="small" effect="plain">{{ keyFactsCount(row.key_facts) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.7.6a07e0')" width="90" align="center">
          <template #default="{ row }">{{ (row.sop_history || []).length }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.8.7f2846')" width="160">
          <template #default="{ row }">{{ formatTime(row.last_active_at || row.updated_at) }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.9.2b6bc0')" width="160" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="viewDetail(row)">{{ $t('auto.List.24.fddda3') }}</el-button>
            <el-button link type="success" size="small" @click="buildContext(row)">{{ $t('auto.List.25.d956fe') }}</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty :description="$t('auto.List.10.85149b')" />
        </template>
      </el-table>
      <div class="pagination-container" v-if="pagination.total > 0">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, prev, pager, next, jumper"
          :total="pagination.total"
          @current-change="loadList"
          @size-change="loadList"
        />
      </div>
    </el-card>

    
    <!-- detailLoading 已移除：v-loading 不能挂在 el-dialog 上（teleport 渲染，根节点非单个元素，
         Vue 每次渲染都报 "Runtime directive used on component with non-element root node"），
         且它从声明起就从未被赋 true，是死状态。弹窗内的短消息/长期记忆分区各自有
         shortLoading / longLoading 转圈，删掉不影响任何可见行为。 -->
    <el-dialog v-model="detailVisible" :title="$t('auto.List.11.f9b256')" width="920px" top="5vh">
      <template v-if="currentMemory">
        
        <el-descriptions :column="3" border size="small">
          <el-descriptions-item :label="$t('auto.List.12.3f61de')">{{ currentMemory.customer_id || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.13.79559f')">{{ currentMemory.customer_name || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.3.ce092e')">
            <el-tag :type="getIntentType(currentMemory.intent_level || currentMemory.purchase_intent)" size="small">
              {{ getIntentLabel(currentMemory.intent_level || currentMemory.purchase_intent) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.4.342cea')">{{ currentMemory.message_count || 0 }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.5.19b1dc')">{{ (currentMemory.objections || []).length }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.8.7f2846')">{{ formatTime(currentMemory.last_active_at || currentMemory.updated_at) }}</el-descriptions-item>
        </el-descriptions>

        <div class="context-bar">
          <el-button type="primary" size="small" :loading="contextLoading" @click="buildContext(currentMemory)">
            <el-icon><Connection /></el-icon> 构建完整上下文
          </el-button>
          <el-input
            v-if="contextText"
            v-model="contextText"
            type="textarea"
            :rows="4"
            readonly
            :placeholder="$t('auto.List.14.65e5d2')"
            class="context-text"
          />
        </div>

        <el-tabs v-model="activeTab" class="detail-tabs">
          
          <el-tab-pane :label="$t('auto.List.15.3932ef')" name="short">
            <div class="message-stream" v-loading="shortLoading">
              <template v-if="shortMessages.length">
                <div v-for="(msg, idx) in shortMessages" :key="idx" :class="['message-item', msg.role]">
                  <div class="msg-avatar">
                    <el-tag :type="getRoleType(msg.role)" size="small" effect="dark">
                      {{ getRoleLabel(msg.role) }}
                    </el-tag>
                  </div>
                  <div class="msg-body">
                    <div class="msg-content">{{ msg.content }}</div>
                    <div class="msg-time">{{ formatTime(msg.timestamp || msg.created_at) }}</div>
                  </div>
                </div>
              </template>
              <el-empty v-else :description="$t('auto.List.16.2f0723')" :image-size="80" />
            </div>
          </el-tab-pane>

          
          <el-tab-pane :label="$t('auto.List.17.de49f7')" name="long">
            <div v-loading="longLoading">
              <template v-if="longTermText">
                <el-alert :title="longTermText" type="info" :closable="false" show-icon />
              </template>
              <el-empty v-else :description="$t('auto.List.18.40bdff')" :image-size="80" />
            </div>
          </el-tab-pane>

          
          <el-tab-pane :label="$t('auto.List.6.507e6a')" name="facts">
            <template v-if="keyFactsList.length">
              <el-descriptions :column="2" border size="small">
                <el-descriptions-item v-for="item in keyFactsList" :key="item.key" :label="item.key">
                  {{ item.value }}
                </el-descriptions-item>
              </el-descriptions>
            </template>
            <el-empty v-else :description="$t('auto.List.19.2e89e2')" :image-size="80" />
          </el-tab-pane>

          
          <el-tab-pane :label="`异议记录 (${(currentMemory.objections || []).length})`" name="objections">
            <template v-if="(currentMemory.objections || []).length">
              <el-timeline>
                <el-timeline-item
                  v-for="(obj, idx) in currentMemory.objections"
                  :key="idx"
                  :timestamp="formatTime(obj.time || obj.created_at)"
                  placement="top"
                >
                  <div class="objection-content">{{ obj.objection || obj.content }}</div>
                  <div v-if="obj.response" class="objection-response">回复：{{ obj.response }}</div>
                </el-timeline-item>
              </el-timeline>
            </template>
            <el-empty v-else :description="$t('auto.List.20.c0974a')" :image-size="80" />
          </el-tab-pane>

          
          <el-tab-pane :label="$t('auto.List.3.ce092e')" name="purchase">
            <div class="intent-block">
              <div class="intent-current">
                <span class="label">{{ $t('auto.List.26.507a08') }}</span>
                <el-tag :type="getIntentType(currentMemory.intent_level || currentMemory.purchase_intent)" size="large">
                  {{ getIntentLabel(currentMemory.intent_level || currentMemory.purchase_intent) }}
                </el-tag>
              </div>
              <el-divider content-position="left">{{ $t('auto.List.27.98003c') }}</el-divider>
              <el-radio-group v-model="intentLevel" @change="updateIntent">
                <el-radio-button value="high">{{ $t('auto.List.28.2601e2') }}</el-radio-button>
                <el-radio-button value="medium">{{ $t('auto.List.29.af4713') }}</el-radio-button>
                <el-radio-button value="low">{{ $t('auto.List.30.b20283') }}</el-radio-button>
              </el-radio-group>
            </div>
          </el-tab-pane>

          
          <el-tab-pane :label="`意图轨迹 (${(currentMemory.intent_trail || []).length})`" name="trail">
            <template v-if="(currentMemory.intent_trail || []).length">
              <div class="intent-trail">
                <el-tag
                  v-for="(t, idx) in currentMemory.intent_trail"
                  :key="idx"
                  :type="getIntentTagType(t.intent || t)"
                  size="small"
                  effect="plain"
                  class="trail-tag"
                >
                  {{ idx + 1 }}. {{ getIntentName(t.intent || t) }}{{ t.confidence ? ` (${Math.round((t.confidence || 0) * 100)}%)` : '' }}
                </el-tag>
              </div>
            </template>
            <el-empty v-else :description="$t('auto.List.21.727aad')" :image-size="80" />
          </el-tab-pane>

          
          <el-tab-pane :label="`SOP 历史 (${(currentMemory.sop_history || []).length})`" name="sop">
            <template v-if="(currentMemory.sop_history || []).length">
              <el-timeline>
                <el-timeline-item v-for="(s, idx) in currentMemory.sop_history" :key="idx" placement="top">
                  <div class="sop-history-item">
                    <el-tag type="warning" size="small" effect="plain">SOP: {{ s.sop_id || s.sop_name || '-' }}</el-tag>
                    <span class="sop-node">节点: {{ s.node_id || '-' }}</span>
                    <span class="sop-action">动作: {{ s.action || '-' }}</span>
                    <span v-if="s.time || s.created_at" class="sop-time">{{ formatTime(s.time || s.created_at) }}</span>
                  </div>
                </el-timeline-item>
              </el-timeline>
            </template>
            <el-empty v-else :description="$t('auto.List.22.53f73b')" :image-size="80" />
          </el-tab-pane>
        </el-tabs>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Refresh, Connection } from '@element-plus/icons-vue'
import { memoryApi } from '@/api/dialogueMemory.js'

const { t } = useI18n()

const loading = ref(false)
const listLoading = ref(false)
const memoryList = ref([])
const filterCustomerID = ref('')
const pagination = ref({ page: 1, pageSize: 20, total: 0 })

const detailVisible = ref(false);
const currentMemory = ref(null)
const activeTab = ref('short')
const shortMessages = ref([])
const shortLoading = ref(false)
const longTermText = ref('')
const longLoading = ref(false)
const intentLevel = ref('')
const contextText = ref('')
const contextLoading = ref(false)

const keyFactsList = computed(() => {
  const facts = currentMemory.value?.key_facts
  if (!facts) return []
  if (Array.isArray(facts)) {
    return facts.map(f => ({ key: f.key || f.name, value: String(f.value ?? f) }))
  }
  if (typeof facts === 'object') {
    return Object.entries(facts).map(([key, value]) => ({ key, value: String(value) }))
  }
  return []
})

const totalObjections = computed(() => {
  return memoryList.value.reduce((sum, m) => sum + (m.objections || []).length, 0)
})

const intentCount = (level) => {
  return memoryList.value.filter(m => (m.intent_level || m.purchase_intent) === level).length
}

const keyFactsCount = (facts) => {
  if (!facts) return 0
  if (Array.isArray(facts)) return facts.length
  if (typeof facts === 'object') return Object.keys(facts).length
  return 0
}

const getIntentType = (level) => {
  if (level === 'high') return 'success'
  if (level === 'medium') return 'warning'
  if (level === 'low') return 'info'
  return 'info'
}

const getIntentLabel = (level) => {
  const map = { high: t('auto.List.28.2601e2'), medium: t('auto.List.29.af4713'), low: t('auto.List.30.b20283') }
  return map[level] || t('auto.List.31.0ebf4f')
}

const getIntentTagType = (type) => {
  if (!type) return 'info'
  if (type === 'purchase') return 'success'
  if (type === 'price_inquiry') return 'primary'
  if (type.startsWith('objection')) return 'warning'
  if (type === 'churn' || type === 'complaint') return 'danger'
  return 'info'
}

const getIntentName = (type) => {
  const map = {
    price_inquiry: t('auto.List.32.be446c'), purchase: t('auto.List.3.ce092e'),
    objection_price: t('auto.List.33.6347cc'), objection_need: t('auto.List.34.01b833'),
    objection_trust: t('auto.List.35.62d02b'), objection_competitor: t('auto.List.36.8180b6'),
    objection_timing: t('auto.List.37.4742a9'), ask_product: t('auto.List.38.516a44'),
    ask_service: t('auto.List.39.62588a'), after_sale: t('auto.List.40.ab76fc'),
    churn: t('auto.List.41.b0be74'), social: t('auto.List.42.7a0b74'), greeting: t('auto.List.43.3f7b30'),
    complaint: t('auto.List.44.e19d1d'), unknown: t('auto.List.45.1622dc')
  }
  return map[type] || type
}

const getRoleType = (role) => {
  if (role === 'user' || role === 'customer') return 'primary'
  if (role === 'assistant' || role === 'ai' || role === 'bot') return 'success'
  return 'info'
}

const getRoleLabel = (role) => {
  if (role === 'user' || role === 'customer') return t('auto.List.2.ff0b20')
  if (role === 'assistant' || role === 'ai' || role === 'bot') return 'AI'
  if (role === 'system') return t('auto.List.46.8a8b89')
  return role || t('auto.List.45.1622dc')
}

const formatTime = (t) => {
  if (!t) return '-'
  const d = new Date(t)
  if (isNaN(d.getTime())) return '-'
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const onSearch = () => {
  pagination.value.page = 1
  loadList()
}

const loadList = async () => {
  listLoading.value = true
  try {
    const params = {
      page: pagination.value.page,
      page_size: pagination.value.pageSize
    }
    if (filterCustomerID.value) params.customer_id = filterCustomerID.value
    const res = await memoryApi.list(params)
    memoryList.value = res?.list || [];
    pagination.value.total = res?.total || 0
  } catch (e) {
    ElMessage.error(t('auto.List.47.94ca59') + (e.message || t('auto.List.48.974e74')))
  } finally {
    listLoading.value = false
  }
}

const refreshAll = async () => {
  loading.value = true
  try {
    await loadList()
  } finally {
    loading.value = false
  }
}

const viewDetail = async (row) => {
  currentMemory.value = row
  detailVisible.value = true
  activeTab.value = 'short'
  intentLevel.value = row.intent_level || row.purchase_intent || 'low'
  contextText.value = ''
  longTermText.value = ''
  await Promise.all([loadShortMessages(row.customer_id), loadLongTerm(row.customer_id)])
}

const loadShortMessages = async (customerId) => {
  if (!customerId) {
    shortMessages.value = []
    return
  }
  shortLoading.value = true
  try {
    const res = await memoryApi.getShortTerm({ customer_id: customerId, limit: 50 })
    shortMessages.value = Array.isArray(res) ? res : (res?.list || (res?.messages || []))
  } catch (e) {
    ElMessage.error(t('auto.List.49.667534') + (e.message || t('auto.List.48.974e74')))
    shortMessages.value = []
  } finally {
    shortLoading.value = false
  }
}

const loadLongTerm = async (customerId) => {
  if (!customerId) {
    longTermText.value = ''
    return
  }
  longLoading.value = true
  try {
    const res = await memoryApi.getLongTerm({ customer_id: customerId, limit: 10 })
    if (typeof res === 'string') {
      longTermText.value = res
    } else if (res?.summary) {
      longTermText.value = res.summary
    } else if (Array.isArray(res) && res.length) {
      longTermText.value = res.map(s => s.summary || s.content || JSON.stringify(s)).join('\n')
    } else {
      longTermText.value = res?.long_term_summary || ''
    }
  } catch (e) {
    ElMessage.error(t('auto.List.50.b871d6') + (e.message || t('auto.List.48.974e74')))
    longTermText.value = ''
  } finally {
    longLoading.value = false
  }
}

const buildContext = async (row) => {
  if (!row?.customer_id) {
    ElMessage.warning(i18n.global.t(t('auto.List.51.330056')))
    return
  }
  contextLoading.value = true
  try {
    const res = await memoryApi.buildContext({ customer_id: row.customer_id })
    if (typeof res === 'string') {
      contextText.value = res
    } else {
      contextText.value = JSON.stringify(res, null, 2)
    }
    ElMessage.success(i18n.global.t(t('auto.List.52.cb7059')))
  } catch (e) {
    ElMessage.error(t('auto.List.53.0f2820') + (e.message || t('auto.List.48.974e74')))
  } finally {
    contextLoading.value = false
  }
}

const updateIntent = async (level) => {
  if (!currentMemory.value?.customer_id) {
    ElMessage.warning(i18n.global.t(t('auto.List.51.330056')))
    return
  }
  try {
    await memoryApi.updatePurchaseIntent({
      customer_id: currentMemory.value.customer_id,
      intent_level: level
    })
    currentMemory.value.intent_level = level
    currentMemory.value.purchase_intent = level
    ElMessage.success(i18n.global.t(t('auto.List.54.11da4e')))
    loadList()
  } catch (e) {
    ElMessage.error(t('auto.List.55.87ef41') + (e.message || t('auto.List.48.974e74')))
  }
}

onMounted(() => refreshAll())
</script>

<style scoped lang="scss">
.dialogue-memory-page { padding: 20px; }

.header-card {
  margin-bottom: 16px;
  :deep(.el-card__body) { padding: 16px 20px; }
  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    h2 { margin: 0 0 6px 0; font-size: 20px; }
    .subtitle { color: #909399; margin: 0; font-size: 13px; }
  }
}

.stats-row {
  margin-bottom: 16px;
  .stat-card {
    text-align: center;
    .stat-value { font-size: 28px; font-weight: bold; color: #303133; }
    .stat-label { color: #909399; font-size: 13px; margin-top: 6px; }
    &.stat-high .stat-value { color: #10B981; }
    &.stat-medium .stat-value { color: #F59E0B; }
    &.stat-objection .stat-value { color: #EF4444; }
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  .header-controls { display: flex; gap: 8px; }
}

.customer-cell {
  display: flex;
  flex-direction: column;
  .customer-name { font-weight: 500; }
  .customer-id { font-size: 12px; color: #909399; }
}

.pagination-container {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.context-bar {
  margin-top: 16px;
  .context-text { margin-top: 12px; }
}

.detail-tabs { margin-top: 16px; }

.message-stream {
  max-height: 420px;
  overflow-y: auto;
  padding: 8px;
  background: #f5f7fa;
  border-radius: 6px;
  .message-item {
    display: flex;
    margin-bottom: 16px;
    &.user, &.customer { flex-direction: row; }
    &.assistant, &.ai, &.bot { flex-direction: row-reverse; }
    .msg-avatar { flex-shrink: 0; }
    .msg-body {
      max-width: 70%;
      margin: 0 10px;
      .msg-content {
        padding: 8px 12px;
        border-radius: 8px;
        background: #fff;
        box-shadow: 0 1px 2px rgba(0,0,0,0.05);
        word-break: break-all;
      }
      .msg-time {
        font-size: 11px;
        color: #909399;
        margin-top: 4px;
      }
    }
    &.assistant .msg-body, &.ai .msg-body, &.bot .msg-body { text-align: right; }
  }
}

.intent-block {
  .intent-current {
    margin-bottom: 16px;
    .label { font-size: 14px; margin-right: 8px; }
  }
}

.intent-trail {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  .trail-tag { margin: 2px; }
}

.objection-content { font-weight: 500; }
.objection-response { color: #10B981; font-size: 13px; margin-top: 4px; }

.sop-history-item {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
  .sop-node, .sop-action { font-size: 13px; color: #606266; }
  .sop-time { font-size: 12px; color: #909399; }
}
</style>
