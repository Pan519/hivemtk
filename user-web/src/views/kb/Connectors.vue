<template>
  <div class="app-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.Connectors.476ab442') }}</span>
          <el-button type="primary" size="small" @click="openCreate">{{ $t('auto.Connectors.080ac62c') }}</el-button>
        </div>
      </template>

      <el-alert
        type="info"
        :closable="false"
        show-icon
        :title="$t('auto.Connectors.d2f5df47')"
        style="margin-bottom: 12px"
      />

      <el-table v-loading="loading" :data="rows" border size="small">
        <el-table-column prop="name" :label="$t('auto.Connectors.d7ec2d3f')" min-width="120" />
        <el-table-column prop="type" :label="$t('auto.Connectors.226b0912')" width="100" />
        <el-table-column :label="$t('auto.Connectors.1abcfdd7')" width="80">
          <template #default="{ row }">
            <el-tag v-if="hasCred(row)" type="success" size="small">{{ $t('auto.Connectors.3794b783') }}</el-tag>
            <el-tag v-else type="info" size="small">{{ $t('auto.Connectors.3f37e954') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Connectors.3fea7ca7')" width="90">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'" size="small">
              {{ row.enabled ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_status" :label="$t('auto.Connectors.34ac0ac3')" width="100" />
        <el-table-column prop="last_summary" :label="$t('auto.Connectors.3ae14696')" min-width="180" show-overflow-tooltip />
        <el-table-column :label="$t('auto.Connectors.2b6bc0f2')" width="220" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" :loading="syncingId === row.id" @click="sync(row)">{{ $t('auto.Connectors.6a620e3c') }}</el-button>
            <el-button size="small" @click="openEdit(row)">{{ $t('auto.Connectors.95b351c8') }}</el-button>
            <el-button size="small" type="danger" @click="remove(row)">{{ $t('auto.Connectors.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" ::title="$t('auto.Connectors.5591a231')" width="560px">
      <el-form label-width="110px">
        <el-form-item :label="$t('auto.Connectors.226b0912')">
          <el-select v-model="form.type" style="width: 100%" :disabled="!!editId">
            <el-option v-for="t in types" :key="t" :value="t" :label="typeLabel(t)" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.Connectors.d7ec2d3f')">
          <el-input v-model="form.name" maxlength="100" :placeholder="$t('auto.Connectors.bb46ae41')" />
        </el-form-item>
        <el-form-item label="Client ID">
          <el-input v-model="form.client_id" :placeholder="$t('auto.Connectors.d8525cf3')" autocomplete="off" />
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="form.api_key" :placeholder="$t('auto.Connectors.2ac9ed5f')" show-password autocomplete="new-password" />
        </el-form-item>
        <el-form-item v-if="form.type === 'ima'" :label="$t('auto.Connectors.c72e13f9')">
          <el-input v-model="form.knowledge_base_id" :placeholder="$t('auto.Connectors.4297f6c8')" />
        </el-form-item>
        <el-form-item v-if="form.type === 'webhttp'" :label="$t('auto.Connectors.b42a0344')">
          <el-input v-model="form.manifest_url" placeholder="https://.../manifest.json" />
        </el-form-item>
        <el-form-item :label="$t('auto.Connectors.25431f2a')">
          <el-input-number v-model="form.interval_minutes" :min="0" :max="1440" /> 分钟（0=仅手动）
        </el-form-item>
        <el-form-item :label="$t('auto.Connectors.7854b52a')">
          <el-switch v-model="form.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.Connectors.625fb26b') }}</el-button>
        <el-button type="primary" :loading="saving" @click="save">{{ $t('auto.Connectors.be5fbbe3') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { kbConnectorApi } from '@/api/kbConnector'

const { t } = useI18n()

const rows = ref([])
const types = ref([])
const loading = ref(false)
const saving = ref(false)
const syncingId = ref(null)
const dialogVisible = ref(false)
const editId = ref(null)
const form = ref(defaultForm())

function defaultForm() {
  return { type: 'ima', name: '', client_id: '', api_key: '', knowledge_base_id: '', manifest_url: '', interval_minutes: 0, enabled: true }
}

const typeLabel = (t) => ({ ima: t('auto.Connectors.695f0c38'), webhttp: t('auto.Connectors.33617f28') }[t] || t)
const hasCred = (row) => {
  const m = row.config_masked || ''
  return m.includes('client_id') || m.includes('"api_key"')
}

const load = async () => {
  loading.value = true
  try {
    const res = await kbConnectorApi.list()
    const data = res?.data ?? res ?? {}
    rows.value = data.list ?? []
    types.value = data.types ?? ['ima', 'webhttp']
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  editId.value = null
  form.value = defaultForm()
  dialogVisible.value = true
}

const openEdit = (row) => {
  editId.value = row.id
  const masked = {}
  try { Object.assign(masked, JSON.parse(row.config_masked || '{}')) } catch { /* 忽略 */ }
  form.value = {
    type: row.type, name: row.name,
    client_id: masked.client_id || '', api_key: '', knowledge_base_id: masked.knowledge_base_id || '',
    manifest_url: masked.manifest_url || '', interval_minutes: row.interval_minutes || 0,
    enabled: !!row.enabled
  }
  dialogVisible.value = true
}

const buildPayload = () => {
  const p = { type: form.value.type, name: form.value.name, enabled: form.value.enabled, interval_minutes: form.value.interval_minutes }
  if (form.value.type === 'ima') {
    p.client_id = form.value.client_id
    if (form.value.api_key) p.api_key = form.value.api_key
    if (form.value.knowledge_base_id) p.knowledge_base_id = form.value.knowledge_base_id
  } else {
    if (form.value.manifest_url) p.manifest_url = form.value.manifest_url
    if (form.value.api_key) p.api_key = form.value.api_key
  }
  return p
}

const save = async () => {
  saving.value = true
  try {
    if (editId.value) await kbConnectorApi.update(editId.value, buildPayload())
    else await kbConnectorApi.create(buildPayload())
    ElMessage.success(t('auto.Connectors.f8dfedcd'))
    dialogVisible.value = false
    await load()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Connectors.6de920b4'))
  } finally {
    saving.value = false
  }
}

const sync = async (row) => {
  syncingId.value = row.id
  try {
    const res = await kbConnectorApi.sync(row.id)
    ElMessage.success(res?.message || t('auto.Connectors.02f667de'))
    await load()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Connectors.d61036e7'))
  } finally {
    syncingId.value = null
  }
}

const remove = async (row) => {
  await ElMessageBox.confirm(`删除连接器「${row.name}」？已导入的文档会保留在知识库。`, t('auto.Connectors.50eaf94d'), { type: 'warning' })
  await kbConnectorApi.remove(row.id)
  ElMessage.success(t('auto.Connectors.5cc23262'))
  await load()
}

onMounted(load)
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
