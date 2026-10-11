<template>
  <div class="agent-status-page">
    <el-card class="header-card">
      <div>
        <h2>{{ $t('坐席状态') }}</h2>
        <p class="subtitle">{{ $t('维护客服在线状态、当前会话数与接待上限') }}</p>
      </div>
      <el-button type="primary" @click="openCreate">
        <el-icon><Plus /></el-icon>
        {{ $t('新增坐席') }}
      </el-button>
    </el-card>

    <el-row :gutter="20" class="stat-row">
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('坐席总数') }}</div>
            <div class="stat-value">{{ agents.length }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('在线坐席') }}</div>
            <div class="stat-value" style="color: #10B981">{{ onlineCount }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('当前会话') }}</div>
            <div class="stat-value">{{ totalSessions }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="stat-label">{{ $t('离线坐席') }}</div>
            <div class="stat-value" style="color: #EF4444">{{ agents.length - onlineCount }}</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('坐席列表') }}</span>
          <el-input v-model="search" ::placeholder="$t('auto.AgentStatus.1c2dec3e')" clearable style="width: 240px" />
        </div>
      </template>
      <el-table :data="filtered" v-loading="loading" stripe>
        <el-table-column prop="agent_id" :label="$t('auto.AgentStatus.79f701de')" width="120" />
        <el-table-column prop="agent_name" :label="$t('auto.AgentStatus.60d0458a')" min-width="120" />
        <el-table-column prop="status" :label="$t('auto.AgentStatus.3fea7ca7')" width="100">
          <template #default="{ row }">
            <el-tag :type="getAgentStatusTagType(row.status)" size="small">
              {{ getAgentStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="active_sessions" :label="$t('auto.AgentStatus.00e43644')" width="100" align="center" />
        <el-table-column prop="max_sessions" :label="$t('auto.AgentStatus.4974bc83')" width="100" align="center" />
        <el-table-column prop="last_active_at" :label="$t('auto.AgentStatus.7f2846ef')" min-width="180">
          <template #default="{ row }">
            {{ formatTime(row.last_active_at || row.online_at) }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.AgentStatus.2b6bc0f2')" width="280" fixed="right">
          <template #default="{ row }">
            <el-button v-if="row.status !== 'online'" link type="success" @click="handleGoOnline(row)">{{ $t('auto.AgentStatus.879eb99c') }}</el-button>
            <el-button v-else link type="warning" @click="handleGoOffline(row)">{{ $t('auto.AgentStatus.4805dd77') }}</el-button>
            <el-button link type="primary" @click="openStatusDialog(row)">{{ $t('auto.AgentStatus.c2211da7') }}</el-button>
            <el-button link type="info" @click="loadSessions(row)">{{ $t('auto.AgentStatus.a85ed23c') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    
    <el-dialog v-model="createVisible" :title="$t('auto.AgentStatus.725ec468')" width="420px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item :label="$t('auto.AgentStatus.79f701de')" required>
          <el-input-number v-model="createForm.id" :min="1" :step="1" :placeholder="$t('auto.AgentStatus.21488ac3')" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('auto.AgentStatus.60d0458a')" required>
          <el-input v-model="createForm.name" />
        </el-form-item>
        <el-form-item :label="$t('auto.AgentStatus.4974bc83')">
          <el-input-number v-model="createForm.max_sessions" :min="1" :max="50" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">{{ $t('auto.AgentStatus.625fb26b') }}</el-button>
        <el-button type="primary" @click="submitCreate">{{ $t('auto.AgentStatus.be5fbbe3') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="statusVisible" :title="$t('auto.AgentStatus.37eee7bf')" width="380px">
      <el-form :model="statusForm" label-width="100px">
        <el-form-item :label="$t('auto.AgentStatus.3fea7ca7')">
          <el-select v-model="statusForm.status" style="width: 100%">
            <el-option :label="$t('auto.AgentStatus.68905cf3')" value="online" />
            <el-option :label="$t('auto.AgentStatus.0f187d38')" value="busy" />
            <el-option :label="$t('auto.AgentStatus.50d4a850')" value="offline" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.AgentStatus.4974bc83')">
          <el-input-number v-model="statusForm.max_sessions" :min="1" :max="50" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="statusVisible = false">{{ $t('auto.AgentStatus.625fb26b') }}</el-button>
        <el-button type="primary" @click="submitStatus">{{ $t('auto.AgentStatus.be5fbbe3') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="sessionsVisible" :title="`坐席 ${currentAgent.agent_name || currentAgent.agent_id || currentAgent.id} 的会话`" width="640px">
      <el-table :data="agentSessions" v-loading="sessionsLoading" stripe>
        <el-table-column prop="id" :label="$t('auto.AgentStatus.427bc5b3')" min-width="100" />
        <el-table-column prop="user_name" :label="$t('auto.AgentStatus.ff0b2077')" min-width="120" />
        <el-table-column :label="$t('auto.AgentStatus.ebee0533')" width="100">
          <template #default="{ row }">
            {{ getChannelLabel(row.platform) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('auto.AgentStatus.3fea7ca7')" width="100">
          <template #default="{ row }">
            <el-tag :type="getSessionStatusTagType(row.status)" size="small">
              {{ getSessionStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="started_at" :label="$t('auto.AgentStatus.592c5958')" min-width="160">
          <template #default="{ row }">
            {{ formatTime(row.started_at) }}
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  listAllAgents,
  createAgent,
  updateAgentStatus,
  goOnline,
  goOffline,
  getAgentSessions
} from '@/api/customerService.js'
import { getChannelLabel } from '@/constants/channel';
import { getEnabledLabel, getEnabledTagType } from '@/constants/enabled';
import { AGENT_STATUS, SESSION_STATUS, getStatusLabel, getStatusTagType } from '@/constants/status';

const getAgentStatusLabel = (s) => getStatusLabel(s, AGENT_STATUS);
const getAgentStatusTagType = (s) => getStatusTagType(s, AGENT_STATUS)
const getSessionStatusLabel = (s) => getStatusLabel(s, SESSION_STATUS);
const getSessionStatusTagType = (s) => getStatusTagType(s, SESSION_STATUS)

const { t } = useI18n()

const loading = ref(false)
const search = ref('')
const agents = ref([])

const createVisible = ref(false)
const createForm = reactive({ id: '', name: '', max_sessions: 5 })

const statusVisible = ref(false)
const statusForm = reactive({ id: null, status: 'online', max_sessions: 5 })

const sessionsVisible = ref(false)
const sessionsLoading = ref(false)
const agentSessions = ref([])
const currentAgent = ref({})

const onlineCount = computed(() => agents.value.filter((a) => a.status === 'online').length)
const totalSessions = computed(() => agents.value.reduce((s, a) => s + (a.active_sessions || 0), 0))

const filtered = computed(() => {
  if (!search.value) return agents.value
  const kw = search.value.toLowerCase()
  return agents.value.filter(
    (a) => String(a.id).includes(kw) || a.name?.toLowerCase().includes(kw)
  )
})

const formatTime = (val) => {
  if (!val) return '-'
  const d = new Date(val)
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleString('zh-CN', { hour12: false })
}

const loadAgents = async () => {
  loading.value = true
  try {
    const res = await listAllAgents()
    const data = res || []
    agents.value = Array.isArray(data) ? data : data.list || []
  } catch (e) {
    ElMessage.error(i18n.global.t(t('auto.AgentStatus.446f1a0e')))
    agents.value = []
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  Object.assign(createForm, { id: 1, name: '', max_sessions: 5 })
  createVisible.value = true
}

const submitCreate = async () => {
  if (!createForm.id || !createForm.name) {
    ElMessage.warning(i18n.global.t(t('auto.AgentStatus.d1e4bff7')))
    return
  }
  const agentId = Number(createForm.id)
  if (!Number.isInteger(agentId) || agentId <= 0) {
    ElMessage.warning(i18n.global.t(t('auto.AgentStatus.4b6ae211')))
    return
  }
  try {
    await createAgent({
      agent_id: agentId,
      agent_name: createForm.name,
      max_sessions: createForm.max_sessions
    })
    ElMessage.success(i18n.global.t(t('auto.AgentStatus.a463c226')))
    createVisible.value = false
    await loadAgents()
  } catch (e) {
    ElMessage.error(e?.message || t('auto.AgentStatus.a889286a'))
  }
}

const handleGoOnline = async (row) => {
  try {
    await goOnline(row.agent_id)
    ElMessage.success(i18n.global.t(t('auto.AgentStatus.8d4a7497')))
    await loadAgents()
  } catch (e) {
    ElMessage.error(i18n.global.t(t('auto.AgentStatus.f70c66cb')))
  }
}

const handleGoOffline = async (row) => {
  try {
    await goOffline(row.agent_id)
    ElMessage.success(i18n.global.t(t('auto.AgentStatus.0a666759')))
    await loadAgents()
  } catch (e) {
    ElMessage.error(i18n.global.t(t('auto.AgentStatus.485df0cc')))
  }
}

const openStatusDialog = (row) => {
  statusForm.id = row.agent_id
  statusForm.status = row.status || 'online'
  statusForm.max_sessions = row.max_sessions || 5
  statusVisible.value = true
}

const submitStatus = async () => {
  try {
    await updateAgentStatus(statusForm.id, {
      status: statusForm.status,
      max_sessions: statusForm.max_sessions
    })
    ElMessage.success(i18n.global.t(t('auto.AgentStatus.7d6e0e1a')))
    statusVisible.value = false
    await loadAgents()
  } catch (e) {
    ElMessage.error(i18n.global.t(t('auto.AgentStatus.930442e2')))
  }
}

const loadSessions = async (row) => {
  currentAgent.value = row
  agentSessions.value = []
  sessionsVisible.value = true
  sessionsLoading.value = true
  try {
    const res = await getAgentSessions(row.agent_id)
    const data = res || []
    agentSessions.value = Array.isArray(data) ? data : data.list || []
  } catch (e) {
    ElMessage.error(i18n.global.t(t('auto.AgentStatus.c59cad21')))
    agentSessions.value = []
  } finally {
    sessionsLoading.value = false
  }
}

onMounted(() => loadAgents())
</script>

<style scoped lang="scss">
.agent-status-page { padding: 20px; }
.header-card {
  margin-bottom: 20px;
  :deep(.el-card__body) { display: flex; justify-content: space-between; align-items: center; }
  h2 { margin: 0 0 8px 0; }
  .subtitle { color: #909399; margin: 0; }
}
.stat-row {
  margin-bottom: 20px;
  .stat-item {
    text-align: center;
    .stat-label { color: #909399; font-size: 14px; margin-bottom: 10px; }
    .stat-value { font-size: 28px; font-weight: bold; }
  }
}
.card-header { display: flex; justify-content: space-between; align-items: center; }
</style>
