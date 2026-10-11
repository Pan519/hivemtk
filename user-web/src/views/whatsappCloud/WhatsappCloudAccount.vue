<template>
  <div class="whatsapp-cloud">
    <el-card class="card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.WhatsappCloudAccount.19.baff3b') }}</span>
          <el-button type="primary" @click="openCreate">{{ $t('auto.WhatsappCloudAccount.20.34374b') }}</el-button>
        </div>
      </template>

      <el-table :data="tableData" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" :label="$t('auto.WhatsappCloudAccount.1.d7ec2d')" min-width="120" />
        <el-table-column prop="remark" :label="$t('auto.WhatsappCloudAccount.2.2432b5')" min-width="120" show-overflow-tooltip />
        <el-table-column prop="phone_number_id" label="PhoneNumberID" min-width="130" show-overflow-tooltip />
        <el-table-column prop="whatsapp_business_id" label="WABA ID" min-width="120" show-overflow-tooltip />
        <el-table-column :label="$t('auto.WhatsappCloudAccount.3.5f1e0e')" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="row.webhook_enabled ? 'success' : 'info'" size="small">
              {{ row.webhook_enabled ? '已启用' : '未启用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.WhatsappCloudAccount.4.bfdd28')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.ai_agent_enabled ? 'success' : 'info'" size="small">
              {{ row.ai_agent_enabled ? '已开启' : '未开启' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.WhatsappCloudAccount.5.2b6bc0')" width="320" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openTestSend(row)">{{ $t('auto.WhatsappCloudAccount.14.e1c652') }}</el-button>
            <el-button link type="primary" @click="openBind(row)">{{ $t('auto.WhatsappCloudAccount.21.c9e87d') }}</el-button>
            <el-button link type="primary" @click="copyCallbackUrl(row)">{{ $t('auto.WhatsappCloudAccount.22.e73c73') }}</el-button>
            <el-button link type="primary" @click="openEdit(row)">{{ $t('auto.WhatsappCloudAccount.23.95b351') }}</el-button>
            <el-button link type="danger" @click="handleDelete(row)">{{ $t('auto.WhatsappCloudAccount.24.2f4aad') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-if="total > pageSize"
        class="pager"
        background
        layout="total, prev, pager, next"
        :total="total"
        :page-size="pageSize"
        :current-page="page"
        @current-change="handlePage"
      />
    </el-card>

    
    <el-dialog v-model="dialogVisible" ::title="$t('auto.WhatsappCloudAccount.6.571f04')" width="640px">
      <el-form :model="form" label-width="150px">
        <el-form-item :label="$t('auto.WhatsappCloudAccount.1.d7ec2d')" prop="name">
          <el-input v-model="form.name" :placeholder="$t('auto.WhatsappCloudAccount.7.17327f')" />
        </el-form-item>
        <el-form-item :label="$t('auto.WhatsappCloudAccount.2.2432b5')" prop="remark">
          <el-input v-model="form.remark" :placeholder="$t('auto.WhatsappCloudAccount.8.67f6a8')" />
        </el-form-item>
        <el-form-item label="PhoneNumberID" prop="phone_number_id">
          <el-input v-model="form.phone_number_id" :placeholder="$t('auto.WhatsappCloudAccount.9.8f4179')" />
        </el-form-item>
        <el-form-item label="WABA ID" prop="whatsapp_business_id">
          <el-input v-model="form.whatsapp_business_id" placeholder="WhatsApp Business Account ID" />
        </el-form-item>
        <el-form-item label="Access Token" prop="access_token">
          <el-input v-model="form.access_token" type="password" show-password :placeholder="$t('auto.WhatsappCloudAccount.10.e466a8')" />
        </el-form-item>
        <el-form-item label="App Secret" prop="app_secret">
          <el-input v-model="form.app_secret" type="password" show-password :placeholder="$t('auto.WhatsappCloudAccount.11.bb6463')" />
        </el-form-item>
        <el-form-item label="Verify Token" prop="verify_token">
          <el-input v-model="form.verify_token" :placeholder="$t('auto.WhatsappCloudAccount.12.16428c')" />
          <div class="form-hint">{{ $t('auto.WhatsappCloudAccount.25.87a987') }}</div>
        </el-form-item>

        <el-divider content-position="left">{{ $t('auto.WhatsappCloudAccount.26.f96de9') }}</el-divider>
        <el-form-item :label="$t('auto.WhatsappCloudAccount.13.5c7f1a')" v-if="dialogMode === 'edit'">
          <el-input :model-value="callbackUrl(form.id)" readonly>
            <template #append>
              <el-button @click="copyText(callbackUrl(form.id))">{{ $t('auto.WhatsappCloudAccount.27.79d3ab') }}</el-button>
            </template>
          </el-input>
          <div class="form-hint">{{ $t('auto.WhatsappCloudAccount.28.938169') }}</div>
        </el-form-item>
        <el-form-item :label="$t('auto.WhatsappCloudAccount.3.5f1e0e')">
          <el-switch v-model="form.webhook_enabled" />
          <span class="form-hint-inline">{{ $t('auto.WhatsappCloudAccount.29.1aad39') }}</span>
        </el-form-item>
        <el-form-item :label="$t('auto.WhatsappCloudAccount.4.bfdd28')">
          <el-switch v-model="form.ai_agent_enabled" />
          <span class="form-hint-inline">{{ $t('auto.WhatsappCloudAccount.30.809880') }}</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.WhatsappCloudAccount.31.625fb2') }}</el-button>
        <el-button type="primary" :loading="saving" @click="submit">{{ $t('auto.WhatsappCloudAccount.32.be5fbb') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="testVisible" :title="$t('auto.WhatsappCloudAccount.14.e1c652')" width="520px">
      <el-form :model="testForm" label-width="90px">
        <el-form-item :label="$t('auto.WhatsappCloudAccount.15.dc2d48')">
          <el-input v-model="testForm.to" :placeholder="$t('auto.WhatsappCloudAccount.16.504107')" />
        </el-form-item>
        <el-form-item :label="$t('auto.WhatsappCloudAccount.17.b87b77')">
          <el-input v-model="testForm.content" type="textarea" :rows="3" :placeholder="$t('auto.WhatsappCloudAccount.18.35af65')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="testVisible = false">{{ $t('auto.WhatsappCloudAccount.31.625fb2') }}</el-button>
        <el-button type="primary" :loading="testing" @click="submitTestSend">{{ $t('auto.WhatsappCloudAccount.33.1535fc') }}</el-button>
      </template>
    </el-dialog>

    
    <agent-binding-dialog
      v-model:visible="bindVisible"
      channel-type="whatsapp"
      :account-id="String(bindRow?.id)"
      :account-name="bindRow?.name"
      :current-agent-id="bindRow?.ai_agent_id ? String(bindRow.ai_agent_id) : ''"
      @saved="onBindSaved"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import AgentBindingDialog from '@/components/AgentBindingDialog.vue'
import {
  listWhatsappCloud,
  getWhatsappCloud,
  createWhatsappCloud,
  updateWhatsappCloud,
  deleteWhatsappCloud,
  testSendWhatsappCloud
} from '@/api/whatsappCloud'

const { t } = useI18n()

const loading = ref(false)
const saving = ref(false)
const testing = ref(false)
const tableData = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const dialogVisible = ref(false)
const dialogMode = ref('create')
const form = reactive({
  id: null,
  name: '',
  remark: '',
  phone_number_id: '',
  whatsapp_business_id: '',
  access_token: '',
  app_secret: '',
  verify_token: '',
  webhook_enabled: false,
  ai_agent_enabled: false
})

const testVisible = ref(false)
const testForm = reactive({ id: null, to: '', content: '' })

const bindVisible = ref(false)
const bindRow = ref(null)

function callbackUrl(id) {
  if (!id) return ''
  return `${window.location.origin}/api/webhook/whatsapp/${id}`
}

async function loadData() {
  loading.value = true
  try {
    const res = await listWhatsappCloud({ page: page.value, page_size: pageSize.value })
    const arr = Array.isArray(res) ? res : (res.data || [])
    tableData.value = arr
    total.value = res && res.total != null ? res.total : arr.length
  } catch (e) {
    ElMessage.error(t('auto.WhatsappCloudAccount.34.c4fa9b') + (e.message || e))
  } finally {
    loading.value = false
  }
}

function resetForm() {
  Object.assign(form, {
    id: null,
    name: '',
    remark: '',
    phone_number_id: '',
    whatsapp_business_id: '',
    access_token: '',
    app_secret: '',
    verify_token: '',
    webhook_enabled: false,
    ai_agent_enabled: false
  })
}

function openCreate() {
  dialogMode.value = 'create'
  resetForm()
  dialogVisible.value = true
}

async function openEdit(row) {
  dialogMode.value = 'edit'
  resetForm()
  try {
    const detail = await getWhatsappCloud(row.id)
    Object.assign(form, {
      id: detail.id,
      name: detail.name || '',
      remark: detail.remark || '',
      phone_number_id: detail.phone_number_id || '',
      whatsapp_business_id: detail.whatsapp_business_id || '',
      access_token: '',
      app_secret: '',
      verify_token: detail.verify_token || '',
      webhook_enabled: !!detail.webhook_enabled,
      ai_agent_enabled: !!detail.ai_agent_enabled
    })
    dialogVisible.value = true
  } catch (e) {
    ElMessage.error(t('auto.WhatsappCloudAccount.35.db870b') + (e.message || e))
  }
}

async function submit() {
  const payload = {
    name: form.name,
    remark: form.remark || undefined,
    phone_number_id: form.phone_number_id || undefined,
    whatsapp_business_id: form.whatsapp_business_id || undefined,
    access_token: form.access_token || undefined,
    app_secret: form.app_secret || undefined,
    verify_token: form.verify_token || undefined,
    webhook_enabled: form.webhook_enabled,
    ai_agent_enabled: form.ai_agent_enabled
  }
  saving.value = true
  try {
    if (dialogMode.value === 'edit') {
      await updateWhatsappCloud(form.id, payload)
    } else {
      await createWhatsappCloud(payload)
    }
    ElMessage.success(t('auto.WhatsappCloudAccount.36.3b1083'))
    dialogVisible.value = false
    loadData()
  } catch (e) {
    ElMessage.error(t('auto.WhatsappCloudAccount.37.78d235') + (e.message || e))
  } finally {
    saving.value = false
  }
}

async function handleDelete(row) {
  try {
    await ElMessageBox.confirm(`确认删除账号「${row.name}」？`, t('auto.WhatsappCloudAccount.38.02d981'), { type: 'warning' })
  } catch {
    return
  }
  try {
    await deleteWhatsappCloud(row.id)
    ElMessage.success(t('auto.WhatsappCloudAccount.39.5cc232'))
    loadData()
  } catch (e) {
    ElMessage.error(t('auto.WhatsappCloudAccount.40.bc42bb') + (e.message || e))
  }
}

function openTestSend(row) {
  testForm.id = row.id
  testForm.to = ''
  testForm.content = ''
  testVisible.value = true
}

async function submitTestSend() {
  if (!testForm.to || !testForm.content) {
    ElMessage.warning(t('auto.WhatsappCloudAccount.41.6943cc'))
    return
  }
  testing.value = true
  try {
    await testSendWhatsappCloud(testForm.id, { to: testForm.to, content: testForm.content })
    ElMessage.success(t('auto.WhatsappCloudAccount.42.9db9a7'))
    testVisible.value = false
  } catch (e) {
    ElMessage.error(t('auto.WhatsappCloudAccount.43.9704ea') + (e.message || e))
  } finally {
    testing.value = false
  }
}

function openBind(row) {
  bindRow.value = row
  bindVisible.value = true
}

function onBindSaved() {
  ElMessage.success('AI 绑定已保存')
  loadData()
}

function copyText(text) {
  if (!text) return
  navigator.clipboard.writeText(text).then(
    () => ElMessage.success(t('auto.WhatsappCloudAccount.44.4fb42e')),
    () => ElMessage.error(t('auto.WhatsappCloudAccount.45.41da1c'))
  )
}

function copyCallbackUrl(row) {
  copyText(callbackUrl(row.id))
}

function handlePage(p) {
  page.value = p
  loadData()
}

onMounted(loadData)
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.pager {
  margin-top: 16px;
  justify-content: flex-end;
}
.form-hint {
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
  margin-top: 4px;
}
.form-hint-inline {
  font-size: 12px;
  color: #909399;
  margin-left: 10px;
}
</style>
