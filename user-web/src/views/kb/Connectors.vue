<template>
  <div class="app-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.Connectors.18.476ab4') }}</span>
          <el-button type="primary" size="small" @click="openCreate">{{ $t('auto.Connectors.19.080ac6') }}</el-button>
        </div>
      </template>

      <el-alert
        type="info"
        :closable="false"
        show-icon
        :title="$t('auto.Connectors.1.d2f5df')"
        style="margin-bottom: 12px"
      />

      <el-table v-loading="loading" :data="rows" border size="small">
        <el-table-column prop="name" :label="$t('auto.Connectors.2.d7ec2d')" min-width="120" />
        <el-table-column prop="type" :label="$t('auto.Connectors.3.226b09')" width="100" />
        <el-table-column :label="$t('auto.Connectors.4.1abcfd')" width="80">
          <template #default="{ row }">
            <el-tag v-if="hasCred(row)" type="success" size="small">{{ $t('auto.Connectors.20.3794b7') }}</el-tag>
            <el-tag v-else type="info" size="small">{{ $t('auto.Connectors.21.3f37e9') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Connectors.5.3fea7c')" width="90">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'" size="small">
              {{ row.enabled ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_status" :label="$t('auto.Connectors.6.34ac0a')" width="100" />
        <el-table-column prop="last_summary" :label="$t('auto.Connectors.7.3ae146')" min-width="180" show-overflow-tooltip />
        <el-table-column :label="$t('auto.Connectors.8.2b6bc0')" width="220" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" :loading="syncingId === row.id" @click="sync(row)">{{ $t('auto.Connectors.22.6a620e') }}</el-button>
            <el-button size="small" @click="openEdit(row)">{{ $t('auto.Connectors.23.95b351') }}</el-button>
            <el-button size="small" type="danger" @click="remove(row)">{{ $t('auto.Connectors.24.2f4aad') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" ::title="$t('auto.Connectors.9.5591a2')" width="560px">
      <el-form label-width="110px">
        <el-form-item :label="$t('auto.Connectors.3.226b09')">
          <el-select v-model="form.type" style="width: 100%" :disabled="!!editId">
            <el-option v-for="t in types" :key="t" :value="t" :label="typeLabel(t)" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.Connectors.2.d7ec2d')">
          <el-input v-model="form.name" maxlength="100" :placeholder="$t('auto.Connectors.10.bb46ae')" />
        </el-form-item>
        <el-form-item label="Client ID">
          <el-input v-model="form.client_id" :placeholder="$t('auto.Connectors.11.d8525c')" autocomplete="off" />
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="form.api_key" :placeholder="$t('auto.Connectors.12.2ac9ed')" show-password autocomplete="new-password" />
        </el-form-item>
        <el-form-item v-if="form.type === 'ima'" :label="$t('auto.Connectors.13.c72e13')">
          <el-input v-model="form.knowledge_base_id" :placeholder="$t('auto.Connectors.14.4297f6')" />
        </el-form-item>
        <el-form-item v-if="form.type === 'webhttp'" :label="$t('auto.Connectors.15.b42a03')">
          <el-input v-model="form.manifest_url" placeholder="https://.../manifest.json" />
        </el-form-item>
        <el-form-item :label="$t('auto.Connectors.16.25431f')">
          <el-input-number v-model="form.interval_minutes" :min="0" :max="1440" /> 分钟（0=仅手动）
        </el-form-item>
        <el-form-item :label="$t('auto.Connectors.17.7854b5')">
          <el-switch v-model="form.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.Connectors.25.625fb2') }}</el-button>
        <el-button type="primary" :loading="saving" @click="save">{{ $t('auto.Connectors.26.be5fbb') }}</el-button>
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

const typeLabel = (t) => ({ ima: t('auto.Connectors.27.695f0c'), webhttp: t('auto.Connectors.28.33617f') }[t] || t)
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
    ElMessage.success(t('auto.Connectors.29.f8dfed'))
    dialogVisible.value = false
    await load()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Connectors.30.6de920'))
  } finally {
    saving.value = false
  }
}

const sync = async (row) => {
  syncingId.value = row.id
  try {
    const res = await kbConnectorApi.sync(row.id)
    ElMessage.success(res?.message || t('auto.Connectors.31.02f667'))
    await load()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Connectors.32.d61036'))
  } finally {
    syncingId.value = null
  }
}

const remove = async (row) => {
  await ElMessageBox.confirm(`删除连接器「${row.name}」？已导入的文档会保留在知识库。`, t('auto.Connectors.33.50eaf9'), { type: 'warning' })
  await kbConnectorApi.remove(row.id)
  ElMessage.success(t('auto.Connectors.34.5cc232'))
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
