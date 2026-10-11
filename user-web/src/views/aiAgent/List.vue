<template>
  <div class="ai-agent-list-page">
    
    <el-card class="header-card" shadow="never">
      <div class="header-content">
        <div>
          <h2>{{ $t('auto.List.44.fd92a0') }}</h2>
          <p class="subtitle">{{ $t('auto.List.45.74691f') }}</p>
        </div>
        <div class="header-actions">
          <el-button @click="loadList" :loading="loading">
            <el-icon><Refresh /></el-icon>
            {{ $t('刷新') }}
          </el-button>
          <el-button type="primary" @click="goCreate">
            <el-icon><Plus /></el-icon>
            {{ $t('新建智能体') }}
          </el-button>
        </div>
      </div>
    </el-card>

    
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filter" @submit.prevent>
        <el-form-item ::label="$t('auto.List.1.bc7a99')">
          <el-select v-model="filter.type" ::placeholder="$t('auto.List.2.dea2f9')" clearable style="width: 160px" @change="onSearch">
            <el-option ::label="$t('auto.List.3.c76c18')" value="sales" />
            <el-option ::label="$t('auto.List.4.323f35')" value="customer_service" />
            <el-option ::label="$t('auto.List.5.2d9fc1')" value="hybrid" />
          </el-select>
        </el-form-item>
        <el-form-item ::label="$t('auto.List.6.1b8578')">
          <el-select v-model="filter.status" ::placeholder="$t('auto.List.7.925332')" clearable style="width: 140px" @change="onSearch">
            <el-option ::label="$t('auto.List.8.0e4c2a')" :value="1" />
            <el-option ::label="$t('auto.List.9.27ea4f')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item ::label="$t('auto.List.10.648229')">
          <el-input
            v-model="filter.keyword"
            ::placeholder="$t('auto.List.11.e31366')"
            clearable
            style="width: 220px"
            @keyup.enter="onSearch"
            @clear="onSearch"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSearch">
            <el-icon><Search /></el-icon>
            {{ $t('搜索') }}
          </el-button>
          <el-button @click="resetFilter">{{ $t('重置') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    
    <el-card shadow="never">
      <el-table :data="list" v-loading="loading" stripe border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column ::label="$t('auto.List.12.ae73fc')" width="80" align="center">
          <template #default="{ row }">
            <el-avatar :size="36" :src="row.avatar" v-if="row.avatar">
              <el-icon><User /></el-icon>
            </el-avatar>
            <el-avatar :size="36" v-else>
              <el-icon><User /></el-icon>
            </el-avatar>
          </template>
        </el-table-column>
        <el-table-column prop="name" :label="$t('auto.List.13.d7ec2d')" min-width="140" show-overflow-tooltip />
        <el-table-column prop="agent_code" :label="$t('auto.List.14.cc6c35')" width="160" show-overflow-tooltip />
        <el-table-column :label="$t('auto.List.15.226b09')" width="120">
          <template #default="{ row }">
            <el-tag :type="getTypeTagType(row.agent_type)" size="small">{{ getTypeLabel(row.agent_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.16.3fea7c')" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="llm_model" :label="$t('auto.List.17.248d0f')" width="140" show-overflow-tooltip>
          <template #default="{ row }">{{ row.llm_model || '-' }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.18.295bb7')" width="170">
          <template #default="{ row }">
            <div class="lang-cell">
              <span>内部：{{ getLanguageLabel(row.internal_language || 'zh') }}</span>
              <span>目标：{{ row.target_language ? getLanguageLabel(row.target_language) : '跟随内部' }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.19.5d5dd3')" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <span class="persona-text">{{ truncateText(row.persona, 50) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.20.eca37c')" width="160">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.21.2b6bc0')" width="280" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="goEdit(row)">{{ $t('auto.List.46.95b351') }}</el-button>
            <el-button
              link
              :type="row.status === 1 ? 'warning' : 'success'"
              size="small"
              @click="onToggle(row)"
            >
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
            <el-button link type="primary" size="small" @click="openTestDialog(row)">{{ $t('auto.List.47.db06c7') }}</el-button>
            <el-button link type="info" size="small" @click="openBindingDialog(row)">{{ $t('auto.List.48.0dacc2') }}</el-button>
            <el-button link type="danger" size="small" @click="onDelete(row)">{{ $t('auto.List.49.2f4aad') }}</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty :description="$t('auto.List.22.fb73da')" />
        </template>
      </el-table>
    </el-card>

    
    <el-dialog v-model="testDialogVisible" :title="$t('auto.List.23.353e8d')" width="760px" top="6vh">
      <el-form :model="testForm" label-width="90px">
        <el-form-item :label="$t('auto.List.24.72fbb9')">
          <el-input :model-value="currentAgentLabel" disabled />
        </el-form-item>
        <el-form-item :label="$t('auto.List.25.1c2ecf')">
          <el-input v-model="testForm.customer_id" :placeholder="$t('auto.List.26.774823')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.27.b87b77')" required>
          <el-input
            v-model="testForm.message"
            type="textarea"
            :rows="3"
            :placeholder="$t('auto.List.28.cc633c')"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="testing" @click="runTest">
            <el-icon><VideoPlay /></el-icon>
            执行测试
          </el-button>
          <el-button @click="clearTestResult">{{ $t('auto.List.50.7e61ac') }}</el-button>
        </el-form-item>
      </el-form>

      
      <template v-if="testResult">
        <el-divider content-position="left">{{ $t('auto.List.51.c6c9da') }}</el-divider>
        <el-descriptions :column="2" border size="small">
          <el-descriptions-item :label="$t('auto.List.29.fab08d')" :span="2">
            <div class="reply-box">{{ testResult.reply || '-' }}</div>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.17.248d0f')">{{ testResult.llm_model || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.30.7f62c0')">{{ testResult.latency_ms || 0 }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.31.8f5f4e')">{{ testResult.cost_tokens || 0 }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.32.39e5fb')">
            <el-tag :type="testResult.transferred_to_human ? 'warning' : 'success'" size="small">
              {{ testResult.transferred_to_human ? '是' : '否' }}
            </el-tag>
            <span v-if="testResult.transfer_reason" class="transfer-reason">（{{ testResult.transfer_reason }}）</span>
          </el-descriptions-item>
        </el-descriptions>

        
        <div class="chain-title">{{ $t('auto.List.52.cda843') }}</div>
        <el-table :data="testResult.steps || []" stripe size="small" border>
          <el-table-column type="index" label="#" width="50" align="center" />
          <el-table-column prop="step" :label="$t('auto.List.33.52b365')" min-width="160" show-overflow-tooltip />
          <el-table-column :label="$t('auto.List.16.3fea7c')" width="90" align="center">
            <template #default="{ row }">
              <el-tag :type="getStepStatusType(row.status)" size="small">{{ row.status || '-' }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('auto.List.30.7f62c0')" width="100" align="center">
            <template #default="{ row }">{{ row.latency_ms || 0 }}</template>
          </el-table-column>
          <el-table-column prop="detail" :label="$t('auto.List.34.f26225')" min-width="200" show-overflow-tooltip />
          <el-table-column prop="error" :label="$t('auto.List.35.7030ff')" min-width="160" show-overflow-tooltip>
            <template #default="{ row }">
              <span v-if="row.error" class="error-text">{{ row.error }}</span>
              <span v-else>-</span>
            </template>
          </el-table-column>
        </el-table>
      </template>
      <el-empty v-else-if="!testing" :description="$t('auto.List.36.5910a5')" :image-size="80" />
    </el-dialog>

    
    <el-dialog v-model="bindingDialogVisible" :title="$t('auto.List.37.88cc83')" width="720px" top="6vh">
      <template v-if="currentAgent">
        <el-descriptions :column="2" border size="small">
          <el-descriptions-item :label="$t('auto.List.38.63ecdf')">{{ currentAgent.name }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.14.cc6c35')">{{ currentAgent.agent_code }}</el-descriptions-item>
        </el-descriptions>
      </template>

      <el-divider content-position="left">{{ $t('auto.List.53.f42c0d') }}</el-divider>
      <el-table :data="channelBindings" v-loading="bindingLoading" stripe size="small" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="channel_type" :label="$t('auto.List.39.474684')" width="120" />
        <el-table-column prop="account_id" :label="$t('auto.List.40.14d171')" min-width="140" show-overflow-tooltip />
        <el-table-column :label="$t('auto.List.41.7e5157')" width="90" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.is_primary" type="success" size="small">{{ $t('auto.List.54.7bbc73') }}</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.42.7854b5')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'" size="small">
              {{ row.enabled ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.20.eca37c')" width="160">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <template #empty>
          <el-empty :description="$t('auto.List.43.17be1d')" :image-size="60" />
        </template>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Plus, Search, User, VideoPlay } from '@element-plus/icons-vue'
import { listAgents, deleteAgent, toggleAgent, testAgent } from '@/api/aiAgent.js'
import { listBindingsByAgent } from '@/api/channelAgentBinding.js'
import { getLanguageLabel } from '@/constants/languages'

const router = useRouter()

const currentAgentLabel = computed(() => {
  const agent = currentAgent.value
  if (!agent) return ''
  const name = agent.name || ''
  const code = agent.agent_code || ''
  return code ? `${name}（${code}）` : name
});

const { t } = useI18n()

const loading = ref(false);
const list = ref([])
const filter = ref({ type: '', status: '', keyword: '' })

const testDialogVisible = ref(false);
const testing = ref(false)
const currentAgent = ref(null)
const testForm = ref({ customer_id: '', message: '' })
const testResult = ref(null)

const bindingDialogVisible = ref(false);
const bindingLoading = ref(false)
const channelBindings = ref([])

const typeMap = {
  sales: t('auto.List.55.facb6b'),
  customer_service: t('auto.List.56.afaab1'),
  hybrid: t('auto.List.57.bac247')
};

const getTypeLabel = (type) => typeMap[type] || type || '-'

const getTypeTagType = (type) => {
  if (type === 'sales') return 'success'
  if (type === 'customer_service') return 'warning'
  if (type === 'hybrid') return 'primary'
  return 'info'
}

const getStepStatusType = (status) => {
  if (status === 'ok') return 'success'
  if (status === 'fail') return 'danger'
  if (status === 'skip') return 'info'
  return 'info'
}

const truncateText = (text, len) => {
  if (!text) return '-'
  return text.length > len ? text.slice(0, len) + '...' : text
}

const formatTime = (t) => {
  if (!t) return '-'
  const d = new Date(t)
  if (isNaN(d.getTime())) return '-'
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const loadList = async () => {
  loading.value = true
  try {
    const params = {}
    if (filter.value.type) params.type = filter.value.type
    if (filter.value.status !== '' && filter.value.status !== null) params.status = filter.value.status
    if (filter.value.keyword) params.keyword = filter.value.keyword
    const res = await listAgents(params)
    list.value = res?.list || [];
  } catch (e) {
    ElMessage.error(t('auto.List.58.2d3424') + (e.message || t('auto.List.59.974e74')))
  } finally {
    loading.value = false
  }
};

const onSearch = () => {
  loadList()
}

const resetFilter = () => {
  filter.value = { type: '', status: '', keyword: '' }
  loadList()
}

const goCreate = () => {
  router.push('/aiAgent/create')
};

const goEdit = (row) => {
  router.push(`/aiAgent/edit/${row.id}`)
}

const onToggle = (row) => {
  const nextStatus = row.status === 1 ? 0 : 1
  const action = nextStatus === 1 ? t('auto.List.42.7854b5') : t('auto.List.60.710ad0')
  ElMessageBox.confirm(`确认${action}智能体「${row.name}」吗？`, t('auto.List.61.ce828a'), {
    type: 'warning',
    confirmButtonText: `确认${action}`,
    cancelButtonText: t('auto.List.62.625fb2')
  }).then(async () => {
    try {
      await toggleAgent(row.id, nextStatus)
      ElMessage.success(`${action}成功`)
      loadList()
    } catch (e) {
      ElMessage.error(`${action}失败：` + (e.message || t('auto.List.59.974e74')))
    }
  }).catch(() => {})
};

const onDelete = (row) => {
  ElMessageBox.confirm(`确认删除智能体「${row.name}」吗？删除后不可恢复。`, t('auto.List.63.50eaf9'), {
    type: 'warning',
    confirmButtonText: t('auto.List.64.631cd2'),
    cancelButtonText: t('auto.List.62.625fb2')
  }).then(async () => {
    try {
      await deleteAgent(row.id)
      ElMessage.success(i18n.global.t(t('auto.List.65.0007d1')))
      loadList()
    } catch (e) {
      ElMessage.error(t('auto.List.66.bc42bb') + (e.message || t('auto.List.59.974e74')))
    }
  }).catch(() => {})
};

const openTestDialog = (row) => {
  currentAgent.value = row
  testForm.value = { customer_id: '', message: '' }
  testResult.value = null
  testDialogVisible.value = true
};

const runTest = async () => {
  if (!testForm.value.message || !testForm.value.message.trim()) {
    ElMessage.warning(i18n.global.t(t('auto.List.67.fa967b')))
    return
  }
  if (!currentAgent.value) return
  testing.value = true
  testResult.value = null
  try {
    const data = {
      customer_id: testForm.value.customer_id || '',
      message: testForm.value.message
    }
    const res = await testAgent(currentAgent.value.id, data)
    testResult.value = res
    ElMessage.success(i18n.global.t(t('auto.List.68.2afc98')))
  } catch (e) {
    ElMessage.error(t('auto.List.69.4b3cb9') + (e.message || t('auto.List.59.974e74')))
  } finally {
    testing.value = false
  }
}

const clearTestResult = () => {
  testResult.value = null
  testForm.value.message = ''
}

const openBindingDialog = async (row) => {
  currentAgent.value = row
  bindingDialogVisible.value = true
  bindingLoading.value = true
  channelBindings.value = []
  try {
    const res = await listBindingsByAgent(row.id)
    channelBindings.value = res?.list || res || []
  } catch (e) {
    ElMessage.error(t('auto.List.70.8c8827') + (e.message || t('auto.List.59.974e74')))
  } finally {
    bindingLoading.value = false
  }
};

onMounted(() => {
  loadList()
});
</script>

<style scoped lang="scss">
.ai-agent-list-page { padding: 20px; }

.header-card {
  margin-bottom: 16px;
  :deep(.el-card__body) { padding: 16px 20px; }
  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    h2 { margin: 0 0 6px 0; font-size: 20px; }
    .subtitle { color: #909399; margin: 0; font-size: 13px; }
    .header-actions { display: flex; gap: 8px; }
  }
}

.filter-card { margin-bottom: 16px; }

.persona-text { color: #606266; }

.lang-cell {
  display: flex;
  flex-direction: column;
  font-size: 12px;
  color: #606266;
  line-height: 1.6;
}

.reply-box {
  white-space: pre-wrap;
  word-break: break-word;
  background: #f5f7fa;
  padding: 10px 12px;
  border-radius: 4px;
  line-height: 1.6;
  max-height: 200px;
  overflow-y: auto;
}

.transfer-reason { color: #F59E0B; margin-left: 6px; }

.error-text { color: #EF4444; }

.chain-title {
  margin: 16px 0 8px 0;
  font-weight: 600;
  color: #303133;
  font-size: 14px;
}
</style>
