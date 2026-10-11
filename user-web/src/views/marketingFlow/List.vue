<template>
  <div class="marketing-flow-page">
    <el-card class="header-card">
      <div>
        <h2>{{ $t('营销流程') }}</h2>
        <p class="subtitle">{{ $t('可视化编排营销自动化流程') }}</p>
      </div>
      <el-button type="primary" @click="showCreateDialog">
        <el-icon><Plus /></el-icon>
        {{ $t('新建流程') }}
      </el-button>
    </el-card>

    <el-card>
      <el-table :data="flows" v-loading="loading" stripe>
        <el-table-column prop="name" ::label="$t('auto.List.f60c0dbc')" min-width="150" />
        <el-table-column prop="trigger" ::label="$t('auto.List.6a8a4f10')" min-width="150" />
        <el-table-column prop="stepCount" ::label="$t('auto.List.360cdc01')" width="100" />
        <el-table-column prop="runCount" ::label="$t('auto.List.857d65e7')" width="100" />
        <el-table-column prop="successRate" ::label="$t('auto.List.699146eb')" width="120">
          <template #default="{ row }">
            <el-progress :percentage="row.successRate" :status="row.successRate > 80 ? 'success' : 'warning'" />
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('auto.List.3fea7ca7')" width="100">
          <template #default="{ row }">
            <el-switch v-model="row.enabled" @change="toggleFlow(row)" />
          </template>
        </el-table-column>
        <el-table-column prop="updatedAt" :label="$t('auto.List.a001a226')" width="180" />
        <el-table-column :label="$t('auto.List.2b6bc0f2')" width="280" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="editFlow(row)">{{ $t('auto.List.95b351c8') }}</el-button>
            <el-button link type="primary" @click="runFlow(row)">{{ $t('auto.List.1a6aa24e') }}</el-button>
            <el-button link type="primary" @click="viewLogs(row)">{{ $t('auto.List.456d29ef') }}</el-button>
            <el-button link type="danger" @click="deleteFlow(row)">{{ $t('auto.List.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty :description="$t('auto.List.21efd88b')" />
        </template>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="900px" :close-on-click-modal="false">
      <el-form :model="form" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('auto.List.794d65af')" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.65939ccd')" prop="trigger">
          <el-select v-model="form.trigger" style="width: 100%">
            <el-option :label="$t('auto.List.e804a9a2')" value="user_register" />
            <el-option :label="$t('auto.List.43eb027a')" value="new_clue" />
            <el-option :label="$t('auto.List.16c7773f')" value="schedule" />
            <el-option :label="$t('auto.List.0cc990ba')" value="manual" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.0fb43600')">
          <div class="flow-builder">
            <div v-for="(step, idx) in form.steps" :key="idx" class="flow-step">
              <div class="step-header">
                <el-tag>步骤 {{ idx + 1 }}</el-tag>
                <el-button link type="danger" @click="removeStep(idx)" v-if="form.steps.length > 1">{{ $t('auto.List.2f4aaddd') }}</el-button>
              </div>
              <el-form-item :label="$t('auto.List.c500cfab')" :label-width="'80px'">
                <el-select v-model="step.action" style="width: 100%">
                  <el-option :label="$t('auto.List.f4f76ad4')" value="send_email" />
                  <el-option :label="$t('auto.List.2d728886')" value="send_sms" />
                  <el-option :label="$t('auto.List.71933ad7')" value="send_whatsapp" />
                  <el-option :label="$t('auto.List.736eaaae')" value="add_tag" />
                  <el-option :label="$t('auto.List.649332ee')" value="assign_sales" />
                  <el-option :label="$t('auto.List.87979227')" value="wait" />
                </el-select>
              </el-form-item>
              <el-form-item :label="$t('auto.List.aa203864')" :label-width="'80px'">
                <el-input-number v-model="step.delay" :min="0" :max="720" />
              </el-form-item>
              <el-form-item :label="$t('auto.List.3d0a2df9')" :label-width="'80px'">
                <el-input v-model="step.params" type="textarea" :rows="2" placeholder='{"template": "welcome"}' />
              </el-form-item>
            </div>
            <el-button @click="addStep" type="primary" link>{{ $t('auto.List.7a3f6a9e') }}</el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.List.625fb26b') }}</el-button>
        <el-button type="primary" @click="submitForm">{{ $t('auto.List.38cf16f2') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="logDialogVisible" :title="$t('auto.List.c84ddfe8')" width="900px">
      <div class="log-header">
        <span>流程：{{ currentFlowName }}</span>
        <span>共 {{ executionLogs.length }} 条日志</span>
      </div>
      <el-table :data="executionLogs" stripe>
        <el-table-column prop="started_at" :label="$t('auto.List.70b3635a')" width="180" />
        <el-table-column prop="status" :label="$t('auto.List.3fea7ca7')" width="100">
          <template #default="{ row }">
            <el-tag :type="logStatusType(row.status)">{{ logStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="current_node" :label="$t('auto.List.52b36576')" min-width="150" />
        <el-table-column :label="$t('auto.List.39f1374d')" width="120">
          <template #default="{ row }">
            {{ calcDuration(row.started_at, row.completed_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="error_message" :label="$t('auto.List.4604d502')" min-width="200">
          <template #default="{ row }">
            <span :style="{ color: row.error_message ? '#EF4444' : '' }">{{ row.error_message}}</span>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  getFlows,
  createFlow,
  updateFlow,
  deleteFlow as deleteFlowApi,
  runFlow as runFlowApi,
  toggleFlow as toggleFlowApi,
  getFlowLogs
} from '@/api/marketingFlow.js'
import { toList } from '@/utils/list.js'

const { t } = useI18n()

const loading = ref(false)
const flows = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref(t('auto.List.c53fe0f6'))
const logDialogVisible = ref(false)
const currentFlowName = ref('')
const executionLogs = ref([])
const formRef = ref()
const form = ref({
  id: 0,
  name: '',
  trigger: 'user_register',
  steps: [{ action: 'send_email', delay: 0, params: '' }]
})
const formRules = {
  name: [{ required: true, message: i18n.global.t(t('auto.List.e9a00996')), trigger: 'blur' }],
  trigger: [{ required: true, message: i18n.global.t(t('auto.List.1c8c6e2c')), trigger: 'change' }]
}

const refreshData = async () => {
  loading.value = true
  try {
    const res = await getFlows()
    flows.value = toList(res)
  } finally {
    loading.value = false
  }
}

const showCreateDialog = () => {
  form.value = { id: 0, name: '', trigger: 'user_register', steps: [{ action: 'send_email', delay: 0, params: '' }] }
  dialogTitle.value = t('auto.List.c53fe0f6')
  dialogVisible.value = true
}

const editFlow = (row) => {
  form.value = JSON.parse(JSON.stringify(row))
  dialogTitle.value = t('auto.List.96273602')
  dialogVisible.value = true
}

const addStep = () => {
  form.value.steps.push({ action: 'send_email', delay: 0, params: '' })
}

const removeStep = (idx) => {
  form.value.steps.splice(idx, 1)
}

const submitForm = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      if (form.value.id) {
        await updateFlow(form.value.id, form.value)
      } else {
        await createFlow(form.value)
      }
      ElMessage.success(i18n.global.t(t('auto.List.3b108349')))
      dialogVisible.value = false
      refreshData()
    } catch (error) {
      ElMessage.error(i18n.global.t(t('auto.List.6de920b4')))
    }
  })
}

const runFlow = async (row) => {
  await runFlowApi(row.id)
  ElMessage.success(i18n.global.t(t('auto.List.60c729d0')))
}

const toggleFlow = async (row) => {
  await toggleFlowApi(row.id, row.enabled)
  ElMessage.success(i18n.global.t(t('auto.List.7d6e0e1a')))
}

const viewLogs = async (row) => {
  try {
    const res = await getFlowLogs(row.id)
    executionLogs.value = res.list || []
    currentFlowName.value = row.name
    logDialogVisible.value = true
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.c93f334d')))
  }
}

const deleteFlow = async (row) => {
  try {
    await ElMessageBox.confirm(`确定删除流程 "${row.name}"？`, t('auto.List.e83a256e'), { type: 'warning' })
    await deleteFlowApi(row.id)
    ElMessage.success(i18n.global.t(t('auto.List.0007d170')))
    refreshData()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(i18n.global.t(t('auto.List.acf0664a')))
  }
}

const logStatusType = (status) => {
  const map = { running: 'warning', completed: 'success', failed: 'danger', cancelled: 'info' }
  return map[status]}
const logStatusText = (status) => {
  const map = { running: t('auto.List.d679aea3'), completed: t('auto.List.fad5222c'), failed: t('auto.List.acd5cb84'), cancelled: t('auto.List.2111ccbb') }
  return map[status] || status
}
const calcDuration = (start, end) => {
  if (!start || !end) return '-'
  const diff = new Date(end).getTime() - new Date(start).getTime()
  if (diff < 1000) return `${diff}ms`
  if (diff < 60000) return `${(diff / 1000).toFixed(1)}s`
  return `${Math.floor(diff / 60000)}m ${Math.floor((diff % 60000) / 1000)}s`
}

onMounted(() => refreshData())
</script>

<style scoped lang="scss">
.marketing-flow-page { padding: 20px; }
.header-card {
  margin-bottom: 20px;
  :deep(.el-card__body) { display: flex; justify-content: space-between; align-items: center; }
  h2 { margin: 0 0 8px 0; }
  .subtitle { color: #909399; margin: 0; }
}
.flow-builder {
  width: 100%;
  .flow-step {
    background: #f5f7fa;
    padding: 15px;
    border-radius: 4px;
    margin-bottom: 10px;
    .step-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 10px;
    }
  }
}
.log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  margin-bottom: 10px;
  color: #606266;
  font-size: 14px;
}
</style>
