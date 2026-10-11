<template>
  <div class="message-hub-container">
    
    <div class="page-header">
      <h2>{{ $t('auto.List.50.1c057a') }}</h2>
      <div class="header-actions">
        <el-button type="primary" @click="openPushDialog">
          <el-icon><Plus /></el-icon>
          {{ $t('推送消息') }}
        </el-button>
        <el-button @click="loadStats">
          <el-icon><DataAnalysis /></el-icon>
          {{ $t('刷新统计') }}
        </el-button>
      </div>
    </div>

    
    <el-row :gutter="16" class="stats-row" v-loading="statsLoading">
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-value">{{ stats.total ?? 0 }}</div>
          <div class="stat-label">{{ $t('消息总数') }}</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card stat-inbound">
          <div class="stat-value">{{ stats.inbound ?? 0 }}</div>
          <div class="stat-label">{{ $t('接收消息') }}</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card stat-outbound">
          <div class="stat-value">{{ stats.outbound ?? 0 }}</div>
          <div class="stat-label">{{ $t('发送消息') }}</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card stat-unread">
          <div class="stat-value">{{ stats.unread ?? 0 }}</div>
          <div class="stat-label">{{ $t('未读消息') }}</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card stat-recent">
          <div class="stat-value">{{ stats.recent_24h ?? 0 }}</div>
          <div class="stat-label">{{ $t('auto.List.51.bc03e4') }}</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="hover" class="stat-card stat-platform">
          <div class="stat-value">{{ platformCount }}</div>
          <div class="stat-label">{{ $t('活跃平台数') }}</div>
        </el-card>
      </el-col>
    </el-row>

    
    <el-card v-if="stats.by_platform && Object.keys(stats.by_platform).length" class="platform-card" shadow="never">
      <template #header>
        <span>{{ $t('平台消息分布') }}</span>
      </template>
      <div class="platform-bars">
        <div v-for="(count, platform) in mergedByPlatform" :key="platform" class="platform-bar-item">
          <span class="platform-name">{{ platformLabel(platform) }}</span>
          <div class="bar-track">
            <div class="bar-fill" :style="{ width: barWidth(count) + '%' }"></div>
          </div>
          <span class="platform-count">{{ count }}</span>
        </div>
      </div>
    </el-card>

    
    <div class="search-form">
      <el-form :inline="true" :model="searchForm" class="search-form-content">
        <el-form-item :label="$t('auto.List.1.498128')">
          <el-select v-model="searchForm.platform" :placeholder="$t('auto.List.2.cca5c3')" clearable style="width: 140px">
            <el-option v-for="p in platformOptions" :key="p" :label="platformLabel(p)" :value="p" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.3.d050a0')">
          <el-input v-model="searchForm.account_id" :placeholder="$t('auto.List.3.d050a0')" clearable style="width: 160px" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.4.427bc5')">
          <el-input v-model="searchForm.conversation_id" :placeholder="$t('auto.List.4.427bc5')" clearable style="width: 160px" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.5.413dc9')">
          <el-input v-model="searchForm.sender_id" :placeholder="$t('auto.List.6.b21fd4')" clearable style="width: 140px" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.7.a465db')">
          <el-select v-model="searchForm.direction" :placeholder="$t('auto.List.8.a8b0c2')" clearable style="width: 120px">
            <el-option :label="$t('auto.List.9.3578f2')" value="inbound" />
            <el-option :label="$t('auto.List.10.36b4aa')" value="outbound" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.11.226b09')">
          <el-select v-model="searchForm.msg_type" :placeholder="$t('auto.List.12.c079e7')" clearable style="width: 120px">
            <el-option v-for="t in msgTypeOptions" :key="t" :label="t" :value="t" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.13.8ad629')">
          <el-select v-model="searchForm.is_read" :placeholder="$t('auto.List.8.a8b0c2')" clearable style="width: 100px">
            <el-option :label="$t('auto.List.14.894a10')" value="false" />
            <el-option :label="$t('auto.List.13.8ad629')" value="true" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.15.6778af')">
          <el-select v-model="searchForm.is_group" :placeholder="$t('auto.List.8.a8b0c2')" clearable style="width: 100px">
            <el-option :label="$t('auto.List.15.6778af')" value="true" />
            <el-option :label="$t('auto.List.16.f994ae')" value="false" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.17.cfb5f1')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('auto.List.18.7370e4')" clearable style="width: 160px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="resetSearch">
            <el-icon><RefreshRight /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    
    <el-table :data="messageList" border style="width: 100%" v-loading="loading"
              @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="44" />
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column :label="$t('auto.List.1.498128')" width="110">
        <template #default="{ row }">
          <el-tag :type="platformTagType(row.platform)" size="small">{{ platformLabel(row.platform) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="account_id" :label="$t('auto.List.19.7035c6')" width="120" show-overflow-tooltip />
      <el-table-column :label="$t('auto.List.7.a465db')" width="80">
        <template #default="{ row }">
          <el-tag :type="getDirectionTagType(row.direction)" size="small">
            {{ getDirectionLabel(row.direction) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="msg_type" :label="$t('auto.List.11.226b09')" width="90">
        <template #default="{ row }">
          <el-tag :type="getMsgTypeTagType(row.msg_type)" size="small" effect="plain">
            {{ getMsgTypeLabel(row.msg_type) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="sender_name" :label="$t('auto.List.20.16b213')" width="140" show-overflow-tooltip />
      <el-table-column prop="receiver_name" :label="$t('auto.List.21.8db216')" width="140" show-overflow-tooltip />
      <el-table-column prop="content" :label="$t('auto.List.22.2d711b')" min-width="240" show-overflow-tooltip />
      <el-table-column :label="$t('auto.List.23.9a834e')" width="120" show-overflow-tooltip>
        <template #default="{ row }">
          <span>{{ row.conversation_id || '-' }}</span>
          <el-tag v-if="row.is_group" type="info" size="small" style="margin-left: 4px">{{ $t('auto.List.52.532575') }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="AI" width="70">
        <template #default="{ row }">
          <el-tag v-if="row.is_ai_reply" type="danger" size="small">AI</el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('auto.List.24.3fea7c')" width="80">
        <template #default="{ row }">
          <el-tag :type="row.is_read ? 'info' : 'danger'" size="small">
            {{ row.is_read ? '已读' : '未读' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="sent_at" :label="$t('auto.List.25.63b34d')" width="170">
        <template #default="{ row }">{{ formatTime(row.sent_at) }}</template>
      </el-table-column>
      <el-table-column :label="$t('auto.List.26.2b6bc0')" width="170" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" size="small" link @click="handleViewDetail(row)">{{ $t('auto.List.53.f26225') }}</el-button>
          <el-button v-if="!row.is_read" type="warning" size="small" link @click="handleMarkRead(row)">{{ $t('auto.List.54.3105f7') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    
    <div class="batch-bar" v-if="selectedRows.length">
      <span>已选 {{ selectedRows.length }} 条</span>
      <el-button size="small" type="primary" @click="handleBatchMarkRead">{{ $t('auto.List.55.d30c40') }}</el-button>
    </div>

    
    <div class="pagination-container">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="pagination.total"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    
    <el-dialog v-model="pushDialogVisible" :title="$t('auto.List.27.38f580')" width="720px" top="5vh">
      <el-form :model="pushForm" :rules="pushRules" ref="pushFormRef" label-width="120px">
        <el-form-item :label="$t('auto.List.1.498128')" prop="platform">
          <el-select v-model="pushForm.platform" :placeholder="$t('auto.List.28.9134b4')" style="width: 100%">
            <el-option v-for="p in platformOptions" :key="p" :label="platformLabel(p)" :value="p" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.3.d050a0')" prop="account_id">
          <el-input v-model="pushForm.account_id" :placeholder="$t('auto.List.29.e4f16d')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.30.795e00')" prop="msg_id">
          <el-input v-model="pushForm.msg_id" :placeholder="$t('auto.List.31.1562af')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.7.a465db')" prop="direction">
          <el-radio-group v-model="pushForm.direction">
            <el-radio value="inbound">{{ $t('auto.List.9.3578f2') }}</el-radio>
            <el-radio value="outbound">{{ $t('auto.List.10.36b4aa') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="$t('auto.List.32.6d0071')" prop="msg_type">
          <el-select v-model="pushForm.msg_type" :placeholder="$t('auto.List.33.5e1872')" style="width: 100%">
            <el-option v-for="t in msgTypeOptions" :key="t" :label="t" :value="t" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.34.0fb440')">
          <el-input v-model="pushForm.sender_id" :placeholder="$t('auto.List.34.0fb440')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.35.37f17f')">
          <el-input v-model="pushForm.sender_name" :placeholder="$t('auto.List.36.0cd6ef')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.37.82147f')">
          <el-input v-model="pushForm.receiver_id" :placeholder="$t('auto.List.37.82147f')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.38.60471b')">
          <el-input v-model="pushForm.receiver_name" :placeholder="$t('auto.List.39.75534a')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.4.427bc5')">
          <el-input v-model="pushForm.conversation_id" :placeholder="$t('auto.List.40.0fa52e')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.22.2d711b')" prop="content">
          <el-input v-model="pushForm.content" type="textarea" :rows="4" :placeholder="$t('auto.List.41.b87b77')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.42.976d38')">
          <el-input v-model="pushForm.media_url" :placeholder="$t('auto.List.43.c1534b')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.15.6778af')">
          <el-switch v-model="pushForm.is_group" />
          <el-input v-if="pushForm.is_group" v-model="pushForm.group_id" :placeholder="$t('auto.List.44.9b7e09')"
                    style="margin-left: 12px; width: 60%" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.45.ae8297')">
          <el-switch v-model="pushForm.is_ai_reply" />
          <el-input v-if="pushForm.is_ai_reply" v-model="pushForm.ai_agent" :placeholder="$t('auto.List.46.5be750')"
                    style="margin-left: 12px; width: 60%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="pushDialogVisible = false">{{ $t('auto.List.56.625fb2') }}</el-button>
        <el-button type="primary" :loading="pushing" @click="submitPush">{{ $t('auto.List.57.6f36dc') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="detailDialogVisible" :title="$t('auto.List.47.1172fc')" width="820px" top="5vh">
      <div v-loading="detailLoading">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="ID">{{ currentMessage.id }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.30.795e00')">{{ currentMessage.msg_id }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.1.498128')">{{ platformLabel(currentMessage.platform) }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.19.7035c6')">{{ currentMessage.account_id }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.7.a465db')">{{ currentMessage.direction === 'inbound' ? '接收' : '发送' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.11.226b09')">{{ currentMessage.msg_type }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.20.16b213')">{{ currentMessage.sender_name || currentMessage.sender_id || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.21.8db216')">{{ currentMessage.receiver_name || currentMessage.receiver_id || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.4.427bc5')">{{ currentMessage.conversation_id || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.15.6778af')">
            <el-tag v-if="currentMessage.is_group" type="info" size="small">{{ $t('auto.List.58.0a60ac') }}</el-tag>
            <span v-else>{{ $t('auto.List.59.c9744f') }}</span>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.45.ae8297')">
            <el-tag v-if="currentMessage.is_ai_reply" type="danger" size="small">{{ currentMessage.ai_agent || 'AI' }}</el-tag>
            <span v-else>{{ $t('auto.List.59.c9744f') }}</span>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.13.8ad629')">
            <el-tag :type="currentMessage.is_read ? 'info' : 'danger'" size="small">
              {{ currentMessage.is_read ? '已读' : '未读' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.25.63b34d')">{{ formatTime(currentMessage.sent_at) }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.48.ef8fd0')">{{ formatTime(currentMessage.read_at) }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.49.5e2aaf')">{{ formatTime(currentMessage.created_at) }}</el-descriptions-item>
        </el-descriptions>
        <div class="detail-content-block">
          <div class="detail-content-label">{{ $t('auto.List.41.b87b77') }}</div>
          <div class="detail-content-body">{{ currentMessage.content || '(空)' }}</div>
        </div>
        <div v-if="currentMessage.media_url" class="detail-content-block">
          <div class="detail-content-label">{{ $t('auto.List.42.976d38') }}</div>
          <div class="detail-content-body">
            <a :href="currentMessage.media_url" target="_blank">{{ currentMessage.media_url }}</a>
          </div>
        </div>
        <div v-if="currentMessage.extra && Object.keys(currentMessage.extra).length" class="detail-content-block">
          <div class="detail-content-label">{{ $t('auto.List.60.874c5f') }}</div>
          <pre class="detail-content-json">{{ JSON.stringify(currentMessage.extra, null, 2) }}</pre>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, RefreshRight, DataAnalysis } from '@element-plus/icons-vue'
import { messageHubApi } from '@/api/messageHub'
import { getMyAgent } from '@/api/customerService'
import AgentSocket from '@/utils/agentSocket'
import { getChannelLabel, getChannelTagType, PLATFORM_GROUP_MEMBERS_REVERSE } from '@/constants/channel';
import { getDirectionLabel, getDirectionTagType } from '@/constants/direction';
import { getMsgTypeLabel, getMsgTypeTagType } from '@/constants/msgType'

const { t } = useI18n()

const loading = ref(false)
const statsLoading = ref(false)
const pushing = ref(false)
const detailLoading = ref(false)
const messageList = ref([])
const selectedRows = ref([])
const stats = ref({})
const platformOptions = ref([])
const msgTypeOptions = ref([])

const searchForm = reactive({
  platform: '',
  account_id: '',
  conversation_id: '',
  sender_id: '',
  direction: '',
  msg_type: '',
  is_read: '',
  is_group: '',
  keyword: ''
})

const pagination = reactive({ page: 1, pageSize: 20, total: 0 })

const mergedByPlatform = computed(() => {
  const src = stats.value.by_platform
  if (!src) return {}
  const out = {}
  for (const [platform, count] of Object.entries(src)) {
    const canonical = PLATFORM_GROUP_MEMBERS_REVERSE[platform] || platform
    out[canonical] = (out[canonical] || 0) + count
  }
  return out
});

const platformCount = computed(() => {
  return Object.keys(mergedByPlatform.value).length
})

const maxPlatformCount = computed(() => {
  const vals = Object.values(mergedByPlatform.value)
  return vals.length ? Math.max(...vals) : 1
})

const barWidth = (count) => {
  if (!maxPlatformCount.value) return 0
  return Math.round((count / maxPlatformCount.value) * 100)
}

const platformLabel = (p) => getChannelLabel(p);
const platformTagType = (p) => getChannelTagType(p)

const formatTime = (t) => {
  if (!t) return '-'
  try {
    const d = new Date(t)
    if (isNaN(d.getTime())) return String(t)
    return d.toLocaleString('zh-CN', { hour12: false })
  } catch (e) {
    return String(t)
  }
}

const pushDialogVisible = ref(false);
const pushFormRef = ref(null)
const pushForm = reactive({
  platform: 'wecom',
  account_id: '',
  msg_id: '',
  direction: 'inbound',
  msg_type: 'text',
  sender_id: '',
  sender_name: '',
  receiver_id: '',
  receiver_name: '',
  conversation_id: '',
  content: '',
  media_url: '',
  is_group: false,
  group_id: '',
  is_ai_reply: false,
  ai_agent: ''
})

const pushRules = {
  platform: [{ required: true, message: i18n.global.t(t('auto.List.61.72b562')), trigger: 'change' }],
  account_id: [{ required: true, message: i18n.global.t(t('auto.List.62.6e738f')), trigger: 'blur' }],
  direction: [{ required: true, message: i18n.global.t(t('auto.List.63.564589')), trigger: 'change' }],
  msg_type: [{ required: true, message: i18n.global.t(t('auto.List.64.23a62a')), trigger: 'change' }],
  content: [{ required: true, message: i18n.global.t(t('auto.List.65.fa967b')), trigger: 'blur' }]
}

const detailDialogVisible = ref(false);
const currentMessage = ref({})

let agentSocketInst = null;
let realtimeTimer = null
const scheduleRealtimeRefresh = () => {
  if (realtimeTimer) return
  realtimeTimer = setTimeout(() => {
    realtimeTimer = null
    fetchMessageList()
    loadStats()
  }, 3000)
}
const setupRealtime = async () => {
  if (agentSocketInst) return
  try {
    const me = await getMyAgent()
    const agentId = me?.agent_id || me?.data?.agent_id
    if (!agentId) return
    agentSocketInst = new AgentSocket(agentId, undefined, {
      onNewMessage: scheduleRealtimeRefresh,
      onNewSession: scheduleRealtimeRefresh,
      onSessionUpdate: scheduleRealtimeRefresh,
      onError: (e) => { console.warn('[messageHub ws]', e) }
    })
    agentSocketInst.connect()
  } catch (e) {
    console.warn('[messageHub] agent socket 初始化失败（不阻塞页面）:', e)
  }
}

onMounted(async () => {
  await loadPlatforms()
  await Promise.all([fetchMessageList(), loadStats()])
  setupRealtime()
})

onUnmounted(() => {
  if (realtimeTimer) { clearTimeout(realtimeTimer); realtimeTimer = null }
  if (agentSocketInst) { agentSocketInst.close(); agentSocketInst = null }
})

const loadPlatforms = async () => {
  try {
    const res = await messageHubApi.getPlatforms()
    platformOptions.value = res.platforms || []
    msgTypeOptions.value = res.msg_types || []
  } catch (e) {
    console.error(t('auto.List.66.1a4fc9'), e)
    ElMessage.error(t('auto.List.67.65a90c'))
  }
}

const loadStats = async () => {
  statsLoading.value = true
  try {
    const res = await messageHubApi.getStats({})
    stats.value = res || {}
  } catch (e) {
    console.error(t('auto.List.68.8522b5'), e)
    ElMessage.error(t('auto.List.69.7b7097'))
  } finally {
    statsLoading.value = false
  }
}

const buildParams = () => {
  const params = {
    page: pagination.page,
    page_size: pagination.pageSize
  }
  if (searchForm.platform) params.platform = searchForm.platform
  if (searchForm.account_id) params.account_id = searchForm.account_id
  if (searchForm.conversation_id) params.conversation_id = searchForm.conversation_id
  if (searchForm.sender_id) params.sender_id = searchForm.sender_id
  if (searchForm.direction) params.direction = searchForm.direction
  if (searchForm.msg_type) params.msg_type = searchForm.msg_type
  if (searchForm.is_read !== '') params.is_read = searchForm.is_read
  if (searchForm.is_group !== '') params.is_group = searchForm.is_group
  if (searchForm.keyword) params.keyword = searchForm.keyword
  return params
}

const fetchMessageList = async () => {
  loading.value = true
  try {
    const res = await messageHubApi.getMessages(buildParams())
    messageList.value = res.list || res.data || []
    pagination.total = res.total || 0
  } catch (e) {
    console.error(t('auto.List.70.f82be4'), e)
    ElMessage.error(t('auto.List.71.b854b3'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchMessageList()
}

const resetSearch = () => {
  Object.keys(searchForm).forEach(k => { searchForm[k] = '' })
  pagination.page = 1
  fetchMessageList()
}

const handleSizeChange = (v) => {
  pagination.pageSize = v
  pagination.page = 1
  fetchMessageList()
}

const handleCurrentChange = (v) => {
  pagination.page = v
  fetchMessageList()
}

const handleSelectionChange = (rows) => {
  selectedRows.value = rows
}

const handleViewDetail = async (row) => {
  detailDialogVisible.value = true
  detailLoading.value = true
  try {
    const res = await messageHubApi.getMessageById(row.id)
    currentMessage.value = res || row
  } catch (e) {
    currentMessage.value = row
  } finally {
    detailLoading.value = false
  }
}

const handleMarkRead = async (row) => {
  try {
    await messageHubApi.markRead([row.id])
    ElMessage.success(i18n.global.t(t('auto.List.72.2d1491')))
    row.is_read = true
    await loadStats()
  } catch (e) {
    ElMessage.error(i18n.global.t(t('auto.List.73.408627')))
  }
}

const handleBatchMarkRead = async () => {
  if (!selectedRows.value.length) return
  const ids = selectedRows.value.map(r => r.id)
  try {
    await messageHubApi.markRead(ids)
    ElMessage.success(`已标记 ${ids.length} 条消息为已读`)
    selectedRows.value.forEach(r => { r.is_read = true })
    await loadStats()
  } catch (e) {
    ElMessage.error(i18n.global.t(t('auto.List.74.29d152')))
  }
}

const openPushDialog = () => {
  Object.assign(pushForm, {
    platform: 'wecom',
    account_id: '',
    msg_id: '',
    direction: 'inbound',
    msg_type: 'text',
    sender_id: '',
    sender_name: '',
    receiver_id: '',
    receiver_name: '',
    conversation_id: '',
    content: '',
    media_url: '',
    is_group: false,
    group_id: '',
    is_ai_reply: false,
    ai_agent: ''
  })
  pushDialogVisible.value = true
}

const submitPush = async () => {
  if (!pushFormRef.value) return
  try {
    await pushFormRef.value.validate()
  } catch (e) {
    return
  }
  pushing.value = true
  try {
    const payload = { ...pushForm }
    if (!payload.msg_id) {
      payload.msg_id = `manual-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
    }
    const res = await messageHubApi.pushMessage(payload)
    if (res && res.duplicate) {
      ElMessage.warning(i18n.global.t(t('auto.List.75.fc3dfa')))
    } else {
      ElMessage.success(i18n.global.t(t('auto.List.76.3aefb5')))
    }
    pushDialogVisible.value = false
    pagination.page = 1
    await Promise.all([fetchMessageList(), loadStats()])
  } catch (e) {
    ElMessage.error(t('auto.List.77.448381') + (e?.message || t('auto.List.78.974e74')))
  } finally {
    pushing.value = false
  }
}
</script>

<style lang="scss" scoped>
.message-hub-container {
  padding: 20px;

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;

    h2 {
      margin: 0;
      font-size: 22px;
      color: #303133;
    }

    .header-actions {
      display: flex;
      gap: 8px;
    }
  }

  .stats-row {
    margin-bottom: 16px;

    .stat-card {
      text-align: center;
      padding: 12px 0;

      .stat-value {
        font-size: 26px;
        font-weight: 700;
        color: #303133;
        line-height: 1.2;
      }

      .stat-label {
        font-size: 13px;
        color: #909399;
        margin-top: 4px;
      }

      &.stat-inbound .stat-value { color: #10B981; }
      &.stat-outbound .stat-value { color: #F59E0B; }
      &.stat-unread .stat-value { color: #EF4444; }
      &.stat-recent .stat-value { color: #4F46E5; }
      &.stat-platform .stat-value { color: #909399; }
    }
  }

  .platform-card {
    margin-bottom: 16px;

    .platform-bars {
      display: flex;
      flex-direction: column;
      gap: 8px;
    }

    .platform-bar-item {
      display: flex;
      align-items: center;
      gap: 12px;

      .platform-name {
        width: 100px;
        font-size: 13px;
        color: #606266;
      }

      .bar-track {
        flex: 1;
        height: 14px;
        background: #f5f7fa;
        border-radius: 7px;
        overflow: hidden;
      }

      .bar-fill {
        height: 100%;
        background: linear-gradient(90deg, #4F46E5, #10B981);
        border-radius: 7px;
        transition: width 0.4s ease;
      }

      .platform-count {
        width: 60px;
        text-align: right;
        font-size: 13px;
        color: #303133;
        font-weight: 600;
      }
    }
  }

  .search-form {
    margin-bottom: 16px;
    padding: 14px;
    background-color: #f5f7fa;
    border-radius: 4px;

    .search-form-content {
      margin: 0;
    }
  }

  .batch-bar {
    margin-top: 12px;
    padding: 8px 12px;
    background: #ecf5ff;
    border-radius: 4px;
    display: flex;
    align-items: center;
    gap: 12px;
    color: #4F46E5;
    font-size: 14px;
  }

  .pagination-container {
    margin-top: 16px;
    display: flex;
    justify-content: center;
  }

  .detail-content-block {
    margin-top: 16px;

    .detail-content-label {
      font-size: 13px;
      color: #606266;
      margin-bottom: 8px;
      font-weight: 600;
    }

    .detail-content-body {
      padding: 12px;
      background: #f5f7fa;
      border-radius: 4px;
      line-height: 1.6;
      color: #303133;
      white-space: pre-wrap;
      word-break: break-all;
    }

    .detail-content-json {
      padding: 12px;
      background: #2d2d2d;
      color: #f8f8f2;
      border-radius: 4px;
      font-family: 'Monaco', 'Menlo', monospace;
      font-size: 12px;
      overflow-x: auto;
      margin: 0;
    }
  }
}
</style>
