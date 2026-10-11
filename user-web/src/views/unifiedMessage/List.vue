<template>
  <div class="unified-message-container">
    
    <el-card class="header-card" shadow="never">
      <div class="page-header">
        <div class="header-text">
          <h2>{{ $t('统一消息') }}</h2>
          <p class="subtitle">{{ $t('auto.List.147e092d') }}</p>
        </div>
        <div class="header-actions">
          <el-button @click="fetchMessageList">
            <el-icon><RefreshRight /></el-icon>
            {{ $t('刷新') }}
          </el-button>
        </div>
      </div>
    </el-card>

    
    <el-tabs v-model="activeChannel" class="channel-tabs" @tab-change="handleChannelChange">
      <el-tab-pane v-for="ch in channelTabs" :key="ch.value" :name="ch.value">
        <template #label>
          <span class="tab-label">
            <el-icon v-if="ch.icon" class="tab-icon"><component :is="ch.icon" /></el-icon>
            <span>{{ ch.label }}</span>
            <el-badge
              v-if="channelCounts[ch.value] !== undefined"
              :value="channelCounts[ch.value]"
              :max="999"
              class="tab-badge"
            />
          </span>
        </template>
      </el-tab-pane>
    </el-tabs>

    
    <div class="search-form">
      <el-form :inline="true" :model="searchForm" class="search-form-content">
        <el-form-item ::label="$t('auto.List.8f8025f7')">
          <el-input
            v-model="searchForm.keyword"
            ::placeholder="$t('auto.List.e844f76c')"
            clearable
            style="width: 220px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item ::label="$t('auto.List.f85e3338')">
          <el-select v-model="searchForm.type" ::placeholder="$t('auto.List.e6855e6d')" clearable style="width: 140px">
            <el-option ::label="$t('auto.List.ba3f5410')" value="system" />
            <el-option ::label="$t('auto.List.9258c2d9')" value="user" />
            <el-option ::label="$t('auto.List.2e02d0a7')" value="notification" />
            <el-option ::label="$t('auto.List.abb905e9')" value="ai" />
          </el-select>
        </el-form-item>
        <el-form-item ::label="$t('auto.List.1b857887')">
          <el-select v-model="searchForm.status" ::placeholder="$t('auto.List.33329123')" clearable style="width: 130px">
            <el-option ::label="$t('auto.List.508ed9a8')" value="unread" />
            <el-option ::label="$t('auto.List.a0f8db91')" value="read" />
            <el-option ::label="$t('auto.List.bf5faf5c')" value="processing" />
          </el-select>
        </el-form-item>
        <el-form-item ::label="$t('auto.List.ed7faea1')">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            :range-separator="$t('auto.List.981cbe31')"
            :start-placeholder="$t('auto.List.a3e3b883')"
            :end-placeholder="$t('auto.List.12f1d7ef')"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 320px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            {{ $t('搜索') }}
          </el-button>
          <el-button @click="resetSearch">
            <el-icon><RefreshRight /></el-icon>
            {{ $t('重置') }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    
    <el-table
      :data="messageList"
      border
      style="width: 100%"
      v-loading="loading"
      ::empty-text="$t('auto.List.2d2d82c3')"
    >
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column
        prop="message_id"
        ::label="$t('auto.List.d15d9d68')"
        min-width="180"
        show-overflow-tooltip
      />
      <el-table-column ::label="$t('auto.List.349bd7bd')" width="110">
        <template #default="{ row }">
          <el-tag
            :type="getChannelTagType(row.channel || row.platform)"
            size="small"
            effect="plain"
          >
            {{ getChannelLabel(row.channel || row.platform) || '-' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column ::label="$t('auto.List.a0ee2507')" width="100">
        <template #default="{ row }">
          <el-tag :type="getTypeTagType(row.type || row.content_type)" size="small">
            {{ getTypeLabel(row.type || row.content_type) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column ::label="$t('auto.List.f2acfbc6')" min-width="280" show-overflow-tooltip>
        <template #default="{ row }">
          <div class="content-cell">
            <el-icon v-if="row.pinned" class="pin-icon"><Top /></el-icon>
            <span class="content-text">{{ row.content || row.text || '(无内容)' }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column
        prop="sender_name"
        ::label="$t('auto.List.be89f69b')"
        min-width="120"
        show-overflow-tooltip
      />
      <el-table-column
        prop="chat_id"
        ::label="$t('auto.List.351d26cb')"
        min-width="120"
        show-overflow-tooltip
      />
      <el-table-column ::label="$t('auto.List.1b857887')" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusTagType(row.status)" size="small">
            {{ getStatusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column ::label="$t('auto.List.ed7faea1')" min-width="170" show-overflow-tooltip>
        <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
      </el-table-column>
      <el-table-column ::label="$t('auto.List.4414e3c5')" width="240" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" size="small" link @click="handleViewDetail(row)">
            <el-icon><View /></el-icon>{{ $t('详情') }}
          </el-button>
          <el-button
            v-if="row.status === 'unread'"
            type="success"
            size="small"
            link
            @click="handleMarkRead(row)"
          >
            <el-icon><Check /></el-icon>{{ $t('已读') }}
          </el-button>
          <el-button type="warning" size="small" link @click="handleResend(row)">
            <el-icon><RefreshRight /></el-icon>{{ $t('重发') }}
          </el-button>
          <el-button type="info" size="small" link @click="handleCopyContent(row)">
            <el-icon><CopyDocument /></el-icon>{{ $t('复制') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    
    <div class="pagination-container">
      <el-pagination
        :current-page="pagination.page"
        :page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="pagination.total"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    
    <el-dialog v-model="detailDialogVisible" :title="$t('auto.List.1172fc0f')" width="800px" top="5vh">
      <div v-loading="detailLoading">
        <el-card class="detail-card">
          <template #header>
            <div class="detail-header">
              <span class="detail-title">{{ currentMessage.message_id }}</span>
              <el-tag
                :type="['pending', 'processing'].includes(currentMessage.status) ? 'warning' : 'success'"
              >
                {{ currentMessage.status }}
              </el-tag>
            </div>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item :label="$t('auto.List.b801b5e3')">{{ currentMessage.id }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.6d00710a')">{{ currentMessage.content_type }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.413dc937')">{{ currentMessage.sender_name }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.9a834e9a')">{{ currentMessage.chat_id }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.63b34d09')">{{ currentMessage.created_at }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.List.a001a226')">{{ currentMessage.updated_at }}</el-descriptions-item>
          </el-descriptions>
          <div class="message-content">
            <div class="content-label">{{ $t('auto.List.b87b7756') }}</div>
            <div class="content-body">{{ currentMessage.content }}</div>
          </div>
        </el-card>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Search, RefreshRight, View, Check, Top, CopyDocument } from '@element-plus/icons-vue'
import { unifiedMessageApi } from '@/api/unifiedMessage'
import { getMyAgent } from '@/api/customerService'
import AgentSocket from '@/utils/agentSocket'
import { CHANNEL_OPTIONS, getChannelLabel, getChannelTagType } from '@/constants/channel'

const { t } = useI18n()

const loading = ref(false);
const detailLoading = ref(false)
const messageList = ref([])
const currentMessage = ref({})
const detailDialogVisible = ref(false)

const channelTabs = [
  { value: 'all', label: t('auto.List.a8b0c204'), icon: 'Grid' },
  ...CHANNEL_OPTIONS.map((c) => ({ value: c.value, label: c.label, icon: c.icon }))
];
const activeChannel = ref('all')
const channelCounts = ref({ all: 0 })

const searchForm = reactive({
  keyword: '',
  type: '',
  status: ''
});
const dateRange = ref([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
});

const TYPE_LABEL = {
  system: t('auto.List.058ebdf2'),
  user: t('auto.List.b492d06f'),
  notification: t('auto.List.5660bcd2'),
  ai: 'AI 回复',
  text: t('auto.List.97d07614'),
  image: t('auto.List.20def794'),
  file: t('auto.List.2a0c4740'),
  event: t('auto.List.10b2761d')
};
const TYPE_TAG = {
  system: 'info',
  user: 'primary',
  notification: 'warning',
  ai: 'success',
  // 文本消息是最常见的内容类型（实测占绝大多数行），这里原先留空串，
  // 导致 el-tag 收到非法 type=""：Vue 每行打两条 Invalid prop 告警，
  // 且颜色回落到 default(primary) 蓝色。中性语义与同表的 system/file 一致，用 'info'。
  text: 'info',
  image: 'success',
  file: 'info',
  event: 'warning'
}
const getTypeLabel = (v) => TYPE_LABEL[v] || (v ? String(v) : '-')
// 未知 type 同样不能回落空串（空串是 el-tag 非法值），中性用 'info'，
// 与本文件 getStatusTagType 的兜底保持一致。
const getTypeTagType = (v) => TYPE_TAG[v] || 'info'

const STATUS_LABEL = {
  unread: t('auto.List.894a101e'),
  read: t('auto.List.8ad629e5'),
  processing: t('auto.List.5d459d55'),
  pending: t('auto.List.047109de'),
  sent: t('auto.List.93d15922'),
  failed: t('auto.List.acd5cb84'),
  received: t('auto.List.b4554f4c')
};
const STATUS_TAG = {
  unread: 'danger',
  read: 'info',
  processing: 'warning',
  pending: 'warning',
  sent: 'success',
  failed: 'danger',
  received: 'success'
}
const getStatusLabel = (v) => STATUS_LABEL[v] || (v ? String(v) : '-')
const getStatusTagType = (v) => STATUS_TAG[v] || 'info'

let agentSocketInst = null;
let realtimeTimer = null

const setupRealtime = async () => {
  if (agentSocketInst) return
  try {
    const me = await getMyAgent()
    const agentId = me?.agent_id || me?.data?.agent_id
    if (!agentId) return
    agentSocketInst = new AgentSocket(agentId, undefined, {
      onNewMessage: () => scheduleRealtimeRefresh(),
      onNewSession: () => scheduleRealtimeRefresh(),
      onSessionUpdate: () => scheduleRealtimeRefresh(),
      onError: (e) => { console.warn('[unifiedMessage ws]', e) }
    })
    agentSocketInst.connect()
  } catch (e) { console.warn("[request] 后台调用失败(已忽略):", e) }
};

const scheduleRealtimeRefresh = () => {
  if (realtimeTimer) return
  realtimeTimer = setTimeout(() => {
    realtimeTimer = null
    fetchMessageList()
  }, 3000)
};

onMounted(() => {
  fetchMessageList()
  setupRealtime()
})

onUnmounted(() => {
  if (realtimeTimer) { clearTimeout(realtimeTimer); realtimeTimer = null }
  if (agentSocketInst) { agentSocketInst.close(); agentSocketInst = null }
})

const handleChannelChange = (val) => {
  pagination.page = 1
  fetchMessageList()
};

const fetchMessageList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: searchForm.keyword,
      type: searchForm.type,
      status: searchForm.status
    }
    if (activeChannel.value && activeChannel.value !== 'all') {
      params.channel = activeChannel.value
    }
    if (Array.isArray(dateRange.value) && dateRange.value.length === 2) {
      params.start_time = dateRange.value[0]
      params.end_time = dateRange.value[1]
    }
    const res = await unifiedMessageApi.getMessages(params)
    messageList.value = res.list || []
    pagination.total = res.total || 0

    const counts = { all: messageList.value.length };
    messageList.value.forEach((m) => {
      const ch = m.channel || m.platform
      if (ch) counts[ch] = (counts[ch] || 0) + 1
    })
    channelCounts.value = counts
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
};

const handleSearch = () => {
  pagination.page = 1
  fetchMessageList()
}

const resetSearch = () => {
  searchForm.keyword = ''
  searchForm.type = ''
  searchForm.status = ''
  dateRange.value = []
  activeChannel.value = 'all'
  pagination.page = 1
  fetchMessageList()
}

const handleSizeChange = (val) => {
  pagination.pageSize = val
  pagination.page = 1
  fetchMessageList()
}

const handleCurrentChange = (val) => {
  pagination.page = val
  fetchMessageList()
}

const handleMarkRead = (row) => {
  if (!row) return
  row.status = 'read'
  ElMessage.success(i18n.global.t(t('auto.List.2d149186')))
};

const handleResend = (row) => {
  ElMessage.info(i18n.global.t(t('auto.List.f519eca8')) + (row.message_id || row.id || '-'))
};

const handleCopyContent = async (row) => {
  const text = row?.content || row?.text || ''
  if (!text) {
    ElMessage.warning(i18n.global.t(t('auto.List.75728639')))
    return
  }
  try {
    if (navigator?.clipboard?.writeText) {
      await navigator.clipboard.writeText(text)
    } else {
      const ta = document.createElement('textarea')
      ta.value = text
      document.body.appendChild(ta)
      ta.select()
      document.execCommand('copy')
      document.body.removeChild(ta)
    }
    ElMessage.success(i18n.global.t(t('auto.List.4fb42e6e')))
  } catch (e) {
    ElMessage.error(i18n.global.t(t('auto.List.b2cb4bb3')))
  }
};

const formatTime = (val) => {
  if (!val) return '-'
  try {
    const d = new Date(val)
    if (Number.isNaN(d.getTime())) return val
    const pad = (n) => String(n).padStart(2, '0')
    return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
  } catch (e) {
    return val
  }
};

const handleViewDetail = async (row) => {
  currentMessage.value = row
  detailDialogVisible.value = true
  await fetchMessageDetail(row.id)
};

const fetchMessageDetail = async (id) => {
  detailLoading.value = true
  try {
    const res = await unifiedMessageApi.getMessageById(id)
    currentMessage.value = res || {}
  } catch (error) {
    console.error(error)
  } finally {
    detailLoading.value = false
  }
}
</script>

<style lang="scss" scoped>
.unified-message-container {
  padding: 20px;

  .header-card {
    margin-bottom: 16px;
  }

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 16px;

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
    }
  }

  .channel-tabs {
    margin-bottom: 16px;
    background: #fff;
    padding: 4px 16px 0;
    border-radius: 6px;
    border: 1px solid #ebeef5;

    .tab-label {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      .tab-icon {
        font-size: 14px;
      }
      .tab-badge {
        margin-left: 4px;
      }
    }

    :deep(.el-tabs__header) {
      margin-bottom: 0;
    }
    :deep(.el-tabs__nav-wrap::after) {
      background-color: transparent;
    }
  }

  .search-form {
    margin-bottom: 16px;
    padding: 15px;
    background-color: #f5f7fa;
    border-radius: 4px;

    .search-form-content {
      margin: 0;
    }
  }

  .content-cell {
    display: flex;
    align-items: center;
    gap: 6px;

    .pin-icon {
      color: #e6a23c;
    }
    .content-text {
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }

  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: center;
  }

  .detail-card {
    .detail-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .detail-title {
        font-size: 18px;
        font-weight: bold;
      }
    }

    .message-content {
      margin-top: 20px;

      .content-label {
        font-size: 14px;
        color: #606266;
        margin-bottom: 10px;
        font-weight: bold;
      }

      .content-body {
        padding: 15px;
        background-color: #f5f7fa;
        border-radius: 4px;
        line-height: 1.6;
        color: #303133;
        white-space: pre-wrap;
      }
    }
  }
}
</style>
