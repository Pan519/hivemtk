<template>
  <div class="ai-agent-list-page">
    
    <el-card class="header-card" shadow="never">
      <div class="header-content">
        <div>
          <h2>{{ $t('auto.List.fd92a00b') }}</h2>
          <p class="subtitle">{{ $t('auto.List.74691fe6') }}</p>
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
        <el-form-item ::label="$t('auto.List.bc7a997c')">
          <el-select v-model="filter.type" ::placeholder="$t('auto.List.dea2f931')" clearable style="width: 160px" @change="onSearch">
            <el-option ::label="$t('auto.List.c76c1861')" value="sales" />
            <el-option ::label="$t('auto.List.323f35da')" value="customer_service" />
            <el-option ::label="$t('auto.List.2d9fc100')" value="hybrid" />
          </el-select>
        </el-form-item>
        <el-form-item ::label="$t('auto.List.1b857887')">
          <el-select v-model="filter.status" ::placeholder="$t('auto.List.92533254')" clearable style="width: 140px" @change="onSearch">
            <el-option ::label="$t('auto.List.0e4c2a13')" :value="1" />
            <el-option ::label="$t('auto.List.27ea4fcb')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item ::label="$t('auto.List.64822934')">
          <el-input
            v-model="filter.keyword"
            ::placeholder="$t('auto.List.e3136615')"
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
        <el-table-column ::label="$t('auto.List.ae73fceb')" width="80" align="center">
          <template #default="{ row }">
            <el-avatar :size="36" :src="row.avatar" v-if="row.avatar">
              <el-icon><User /></el-icon>
            </el-avatar>
            <el-avatar :size="36" v-else>
              <el-icon><User /></el-icon>
            </el-avatar>
          </template>
        </el-table-column>
        <el-table-column prop="name" :label="$t('auto.List.d7ec2d3f')" min-width="140" show-overflow-tooltip />
        <el-table-column prop="agent_code" :label="$t('auto.List.cc6c35a3')" width="160" show-overflow-tooltip />
        <el-table-column :label="$t('auto.List.226b0912')" width="120">
          <template #default="{ row }">
            <el-tag :type="getTypeTagType(row.agent_type)" size="small">{{ getTypeLabel(row.agent_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.3fea7ca7')" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="llm_model" :label="$t('auto.List.248d0fb8')" width="140" show-overflow-tooltip>
          <template #default="{ row }">{{ row.llm_model || '-' }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.295bb704')" width="170">
          <template #default="{ row }">
            <div class="lang-cell">
              <span>内部：{{ getLanguageLabel(row.internal_language || 'zh') }}</span>
              <span>目标：{{ row.target_language ? getLanguageLabel(row.target_language) : '跟随内部' }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.5d5dd350')" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <span class="persona-text">{{ truncateText(row.persona, 50) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.eca37cb0')" width="160">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.2b6bc0f2')" width="280" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="goEdit(row)">{{ $t('auto.List.95b351c8') }}</el-button>
            <el-button
              link
              :type="row.status === 1 ? 'warning' : 'success'"
              size="small"
              @click="onToggle(row)"
            >
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
            <el-button link type="primary" size="small" @click="openTestDialog(row)">{{ $t('auto.List.db06c78d') }}</el-button>
            <el-button link type="info" size="small" @click="openBindingDialog(row)">{{ $t('auto.List.0dacc2d9') }}</el-button>
            <el-button link type="danger" size="small" @click="onDelete(row)">{{ $t('auto.List.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty :description="$t('auto.List.fb73da0f')" />
        </template>
      </el-table>
    </el-card>

    
    <el-dialog v-model="testDialogVisible" :title="$t('auto.List.353e8df4')" width="760px" top="6vh">
      <el-form :model="testForm" label-width="90px">
        <el-form-item :label="$t('auto.List.72fbb9b5')">
          <el-input :model-value="currentAgentLabel" disabled />
        </el-form-item>
        <el-form-item :label="$t('auto.List.1c2ecf70')">
          <el-input v-model="testForm.customer_id" :placeholder="$t('auto.List.77482360')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.b87b7756')" required>
          <el-input
            v-model="testForm.message"
            type="textarea"
            :rows="3"
            :placeholder="$t('auto.List.cc633c28')"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="testing" @click="runTest">
            <el-icon><VideoPlay /></el-icon>
            执行测试
          </el-button>
          <el-button @click="clearTestResult">{{ $t('auto.List.7e61acdd') }}</el-button>
        </el-form-item>
      </el-form>

      
      <template v-if="testResult">
        <el-divider content-position="left">{{ $t('auto.List.c6c9da2b') }}</el-divider>
        <el-descriptions :column="2" border size="small">
          <el-descriptions-item :label="$t('auto.List.fab08d35')" :span="2">
            <div class="reply-box">{{ testResult.reply || '-' }}</div>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.248d0fb8')">{{ testResult.llm_model || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.7f62c087')">{{ testResult.latency_ms || 0 }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.8f5f4e35')">{{ testResult.cost_tokens || 0 }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.39e5fbd5')">
            <el-tag :type="testResult.transferred_to_human ? 'warning' : 'success'" size="small">
              {{ testResult.transferred_to_human ? '是' : '否' }}
            </el-tag>
            <span v-if="testResult.transfer_reason" class="transfer-reason">（{{ testResult.transfer_reason }}）</span>
          </el-descriptions-item>
        </el-descriptions>

        
        <div class="chain-title">{{ $t('auto.List.cda843b5') }}</div>
        <el-table :data="testResult.steps || []" stripe size="small" border>
          <el-table-column type="index" label="#" width="50" align="center" />
          <el-table-column prop="step" :label="$t('auto.List.52b36576')" min-width="160" show-overflow-tooltip />
          <el-table-column :label="$t('auto.List.3fea7ca7')" width="90" align="center">
            <template #default="{ row }">
              <el-tag :type="getStepStatusType(row.status)" size="small">{{ row.status || '-' }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('auto.List.7f62c087')" width="100" align="center">
            <template #default="{ row }">{{ row.latency_ms || 0 }}</template>
          </el-table-column>
          <el-table-column prop="detail" :label="$t('auto.List.f26225bd')" min-width="200" show-overflow-tooltip />
          <el-table-column prop="error" :label="$t('auto.List.7030ff64')" min-width="160" show-overflow-tooltip>
            <template #default="{ row }">
              <span v-if="row.error" class="error-text">{{ row.error }}</span>
              <span v-else>-</span>
            </template>
          </el-table-column>
        </el-table>
      </template>
      <el-empty v-else-if="!testing" :description="$t('auto.List.5910a5ed')" :image-size="80" />
    </el-dialog>

    
    <el-dialog v-model="bindingDialogVisible" :title="$t('auto.List.88cc83f2')" width="720px" top="6vh">
      <template v-if="currentAgent">
        <el-descriptions :column="2" border size="small">
          <el-descriptions-item :label="$t('auto.List.63ecdffc')">{{ currentAgent.name }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.cc6c35a3')">{{ currentAgent.agent_code }}</el-descriptions-item>
        </el-descriptions>
      </template>

      <el-divider content-position="left">{{ $t('auto.List.f42c0df9') }}</el-divider>
      <el-table :data="channelBindings" v-loading="bindingLoading" stripe size="small" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="channel_type" :label="$t('auto.List.474684d1')" width="120" />
        <el-table-column prop="account_id" :label="$t('auto.List.14d1718c')" min-width="140" show-overflow-tooltip />
        <el-table-column :label="$t('auto.List.7e5157d5')" width="90" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.is_primary" type="success" size="small">{{ $t('auto.List.7bbc7364') }}</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.7854b52a')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'" size="small">
              {{ row.enabled ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.eca37cb0')" width="160">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <template #empty>
          <el-empty :description="$t('auto.List.17be1d7f')" :image-size="60" />
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
  sales: t('auto.List.facb6bff'),
  customer_service: t('auto.List.afaab1c7'),
  hybrid: t('auto.List.bac24769')
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
    ElMessage.error(t('auto.List.2d34242a') + (e.message || t('auto.List.974e7484')))
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
  const action = nextStatus === 1 ? t('auto.List.7854b52a') : t('auto.List.710ad08b')
  ElMessageBox.confirm(`确认${action}智能体「${row.name}」吗？`, t('auto.List.ce828a79'), {
    type: 'warning',
    confirmButtonText: `确认${action}`,
    cancelButtonText: t('auto.List.625fb26b')
  }).then(async () => {
    try {
      await toggleAgent(row.id, nextStatus)
      ElMessage.success(`${action}成功`)
      loadList()
    } catch (e) {
      ElMessage.error(`${action}失败：` + (e.message || t('auto.List.974e7484')))
    }
  }).catch(() => {})
};

const onDelete = (row) => {
  ElMessageBox.confirm(`确认删除智能体「${row.name}」吗？删除后不可恢复。`, t('auto.List.50eaf94d'), {
    type: 'warning',
    confirmButtonText: t('auto.List.631cd220'),
    cancelButtonText: t('auto.List.625fb26b')
  }).then(async () => {
    try {
      await deleteAgent(row.id)
      ElMessage.success(i18n.global.t(t('auto.List.0007d170')))
      loadList()
    } catch (e) {
      ElMessage.error(t('auto.List.bc42bb8f') + (e.message || t('auto.List.974e7484')))
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
    ElMessage.warning(i18n.global.t(t('auto.List.fa967bd0')))
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
    ElMessage.success(i18n.global.t(t('auto.List.2afc9835')))
  } catch (e) {
    ElMessage.error(t('auto.List.4b3cb9fe') + (e.message || t('auto.List.974e7484')))
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
    ElMessage.error(t('auto.List.8c8827ad') + (e.message || t('auto.List.974e7484')))
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
