<template>
  <div class="workflow-orchestrator-page">
    <el-card class="header-card">
      <div class="header-left">
        <h2>{{ $t('工作流编排') }}</h2>
        <p class="subtitle">{{ $t('可视化编排业务工作流，支持触发器、条件分支和子流程') }}</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="showCreateDialog">
          <el-icon><Plus /></el-icon>
          {{ $t('新建工作流') }}
        </el-button>
      </div>
    </el-card>

    <el-card>
      <div class="filter-bar">
        <el-input
          v-model="filterKeyword"
          ::placeholder="$t('auto.List.c502df5d')"
          style="width: 300px"
          clearable
          :prefix-icon="Search"
          @clear="onFilterChange"
          @keyup.enter="onFilterChange"
        />
        <el-select v-model="filterStatus" ::placeholder="$t('auto.List.1b857887')" clearable style="width: 140px" @change="onFilterChange">
          <el-option :label="$t('auto.List.22b4334f')" value="draft" />
          <el-option :label="$t('auto.List.dca0c13b')" value="published" />
          <el-option :label="$t('auto.List.c3ba167c')" value="archived" />
        </el-select>
        <el-button type="primary" :icon="Search" @click="onFilterChange">{{ $t('查询') }}</el-button>
      </div>

      <el-table :data="filteredFlows" v-loading="loading" stripe>
        <el-table-column prop="workflow_id" ::label="$t('auto.List.9d5b99b1')" min-width="160" />
        <el-table-column prop="version" ::label="$t('auto.List.941bbdf1')" width="80" align="center">
          <template #default="{ row }">
            <el-tag size="small">{{ row.version }}v</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="name" ::label="$t('auto.List.f0213e4c')" min-width="180" />
        <el-table-column prop="status" ::label="$t('auto.List.1b857887')" width="110" align="center">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" effect="dark">{{ statusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" ::label="$t('auto.List.013ab345')" width="180">
          <template #default="{ row }">{{ formatTime(row.updated_at) }}</template>
        </el-table-column>
        <el-table-column ::label="$t('auto.List.4414e3c5')" width="380" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openEditor(row)">{{ $t('编辑') }}</el-button>
            <el-button
              v-if="row.status === 'draft'"
              link
              type="success"
              @click="handlePublish(row)"
              >{{ $t('发布') }}</el-button
            >
            <el-button
              v-if="row.status === 'published'"
              link
              type="warning"
              @click="handleArchive(row)"
              >{{ $t('归档') }}</el-button
            >
            <el-button link type="info" @click="handleExecute(row)">{{ $t('执行') }}</el-button>
            <el-button link type="danger" @click="handleDelete(row)">{{ $t('删除') }}</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty ::description="$t('auto.List.bea3d2b6')" />
        </template>
      </el-table>

      <div class="pagination-bar">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="pagination.total"
          @size-change="onSizeChange"
          @current-change="loadData"
        />
      </div>
    </el-card>

    
    <el-dialog v-model="dialogVisible" ::title="$t('auto.List.4e19187e')" width="600px" :close-on-click-modal="false">
      <el-form :model="form" :rules="formRules" ref="formRef" label-width="120px">
        <el-form-item :label="$t('auto.List.dfc57f14')" prop="workflow_id">
          <el-input v-model="form.workflow_id" ::placeholder="$t('auto.List.59c7234f')" />
        </el-form-item>
        <el-form-item ::label="$t('auto.List.f0213e4c')" prop="name">
          <el-input v-model="form.name" ::placeholder="$t('auto.List.53893125')" />
        </el-form-item>
        <el-form-item ::label="$t('auto.List.d1ef7e0b')">
          <el-input v-model="form.description" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('取消') }}</el-button>
        <el-button type="primary" @click="submitCreate">{{ $t('创建') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="executeDialogVisible" ::title="$t('auto.List.68429f7a')" width="560px">
      <el-form :model="executeForm" label-width="120px">
        <el-form-item ::label="$t('auto.List.9d5b99b1')">
          <el-input v-model="executeForm.workflow_id" disabled />
        </el-form-item>
        <el-form-item ::label="$t('auto.List.57668c5d')">
          <el-input
            v-model="executeForm.trigger_payload"
            type="textarea"
            :rows="6"
            :placeholder="$t('JSON 格式，例如 {&quot;order_id&quot;: 123}')"
          />
        </el-form-item>
        <el-form-item v-if="lastExecution" ::label="$t('auto.List.f5f380a1')">
          <el-tag :type="execStatusType(lastExecution.status)" effect="dark">
            {{ execStatusText(lastExecution.status) }}
          </el-tag>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="executeDialogVisible = false">{{ $t('取消') }}</el-button>
        <el-button type="primary" :loading="executing" @click="doExecute">{{ $t('执行') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import { workflowOrchestratorApi } from '@/api/workflowOrchestrator.js'

const router = useRouter()
const { t } = useI18n()

const loading = ref(false)
const filterKeyword = ref('')
const filterStatus = ref('')
const flows = ref([])
const lastExecution = ref(null)

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const filteredFlows = computed(() => {
  let list = flows.value;
  if (filterKeyword.value) {
    const kw = filterKeyword.value.toLowerCase()
    list = list.filter(f => (f.workflow_id || '').toLowerCase().includes(kw) || (f.name || '').toLowerCase().includes(kw))
  }
  return list
})

const dialogVisible = ref(false);
const formRef = ref()
const form = reactive({
  workflow_id: '',
  name: '',
  description: ''
})
const formRules = {
  workflow_id: [{ required: true, message: t('auto.List.4015e61f'), trigger: 'blur' }],
  name: [{ required: true, message: t('auto.List.06e2f88f'), trigger: 'blur' }]
}

const showCreateDialog = () => {
  form.workflow_id = ''
  form.name = ''
  form.description = ''
  dialogVisible.value = true
}

const submitCreate = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    const definition = { nodes: [], edges: [] }
    await workflowOrchestratorApi.createVersion({
      workflow_id: form.workflow_id,
      name: form.name,
      description: form.description,
      definition
    })
    ElMessage.success(t('auto.List.04a691b3'))
    dialogVisible.value = false
    loadData()
  } catch (e) {
    if (e !== false) ElMessage.error(t('auto.List.a2bddca3') + (e?.message || ''))
  }
}

const executeDialogVisible = ref(false);
const executing = ref(false)
const executeForm = reactive({
  workflow_id: '',
  trigger_payload: '{}'
})

const handleExecute = async (row) => {
  executeForm.workflow_id = row.workflow_id
  executeForm.trigger_payload = '{}'
  lastExecution.value = null
  executeDialogVisible.value = true
}

const doExecute = async () => {
  executing.value = true
  try {
    let payload = {}
    try {
      payload = JSON.parse(executeForm.trigger_payload || '{}')
    } catch {
      ElMessage.error(t('auto.List.dc790a92'))
      executing.value = false
      return
    }
    const result = await workflowOrchestratorApi.execute({
      workflow_id: executeForm.workflow_id,
      trigger_payload: payload
    })
    const data = result?.data || result
    lastExecution.value = data
    ElMessage.success(t('auto.List.ba530c37'))
    if (data?.id) {
      router.push({ name: 'WorkflowOrchestratorExecution', params: { id: data.id } })
    }
    executeDialogVisible.value = false
  } catch (e) {
    ElMessage.error(t('auto.List.23cc6892') + (e?.message || ''))
  } finally {
    executing.value = false
  }
}

const openEditor = (row) => {
  router.push({ name: 'WorkflowOrchestratorEditor', params: { workflow_id: row.workflow_id } })
};

const handlePublish = async (row) => {
  try {
    await ElMessageBox.confirm(`确认发布版本 ${row.version}v？`, t('auto.List.9861a6d5'), { type: 'warning' })
    await workflowOrchestratorApi.publishVersion(row.id)
    ElMessage.success(t('auto.List.2fff079b'))
    loadData()
  } catch (_) { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
}

const handleArchive = async (row) => {
  try {
    await ElMessageBox.confirm(`确认归档版本 ${row.version}v？`, t('auto.List.b3f4b1b1'), { type: 'warning' })
    await workflowOrchestratorApi.archiveVersion(row.id)
    ElMessage.success(t('auto.List.c3ba167c'))
    loadData()
  } catch (_) { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确认删除版本 ${row.version}v？此操作不可恢复`, t('auto.List.50eaf94d'), { type: 'error' })
    await workflowOrchestratorApi.deleteVersion(row.id)
    ElMessage.success(t('auto.List.0007d170'))
    loadData()
  } catch (_) { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
}

const loadData = async () => {
  loading.value = true
  try {
    const params = {
      workflow_id: '',
      status: filterStatus.value || '',
      page: pagination.page,
      page_size: pagination.pageSize
    };
    const result = await workflowOrchestratorApi.listVersions(params)
    const payload = result?.data || result || {}
    flows.value = payload.list || []
    pagination.total = payload.total || 0
  } catch (e) {
    flows.value = []
    pagination.total = 0
    ElMessage.warning(t('auto.List.4bbdceb5') + (e?.message || ''))
  } finally {
    loading.value = false
  }
};

const onFilterChange = () => {
  pagination.page = 1
  loadData()
};

const onSizeChange = () => {
  pagination.page = 1
  loadData()
};

// 两个 type 映射的兜底都不能用空串：el-tag 的 type 合法值只有 primary/success/info/
// warning/danger，传 '' 会触发 `Invalid prop: validation failed for prop "type"` 告警，
// 并回落到 default(primary) 显示成蓝色，与「未知状态」的中性语义不符。
const statusType = (s) => ({ draft: 'info', published: 'success', archived: 'warning' }[s] || 'info');
const statusText = (s) => ({ draft: t('auto.List.22b4334f'), published: t('auto.List.dca0c13b'), archived: t('auto.List.c3ba167c') }[s] || s)
const execStatusType = (s) => ({ running: 'primary', completed: 'success', failed: 'danger', terminated: 'warning' }[s] || 'info')
const execStatusText = (s) => ({ running: t('auto.List.d679aea3'), completed: t('auto.List.fad5222c'), failed: t('auto.List.acd5cb84'), terminated: t('auto.List.2554120a') }[s] || s)

const formatTime = (t) => {
  if (!t) return '-'
  const d = new Date(t)
  return d.toLocaleString('zh-CN', { hour12: false })
}

onMounted(loadData)
</script>

<style scoped>
.workflow-orchestrator-page {
  padding: 16px;
}
.header-card {
  margin-bottom: 16px;
}
.header-card .header-left h2 {
  margin: 0 0 4px 0;
}
.header-card .subtitle {
  margin: 0;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}
.header-card :deep(.el-card__body) {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.filter-bar {
  margin-bottom: 16px;
  display: flex;
  gap: 8px;
  align-items: center;
}
.pagination-bar {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
