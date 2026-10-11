<template>
  <div class="whatsapp-cloud">
    <el-card class="card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.WhatsappCloudAccount.baff3bb1') }}</span>
          <el-button type="primary" @click="openCreate">{{ $t('auto.WhatsappCloudAccount.34374b71') }}</el-button>
        </div>
      </template>

      <el-table :data="tableData" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" :label="$t('auto.WhatsappCloudAccount.d7ec2d3f')" min-width="120" />
        <el-table-column prop="remark" :label="$t('auto.WhatsappCloudAccount.2432b575')" min-width="120" show-overflow-tooltip />
        <el-table-column prop="phone_number_id" label="PhoneNumberID" min-width="130" show-overflow-tooltip />
        <el-table-column prop="whatsapp_business_id" label="WABA ID" min-width="120" show-overflow-tooltip />
        <el-table-column :label="$t('auto.WhatsappCloudAccount.5f1e0e72')" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="row.webhook_enabled ? 'success' : 'info'" size="small">
              {{ row.webhook_enabled ? '已启用' : '未启用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.WhatsappCloudAccount.bfdd28c7')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.ai_agent_enabled ? 'success' : 'info'" size="small">
              {{ row.ai_agent_enabled ? '已开启' : '未开启' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.WhatsappCloudAccount.2b6bc0f2')" width="320" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openTestSend(row)">{{ $t('auto.WhatsappCloudAccount.e1c652bc') }}</el-button>
            <el-button link type="primary" @click="openBind(row)">{{ $t('auto.WhatsappCloudAccount.c9e87d56') }}</el-button>
            <el-button link type="primary" @click="copyCallbackUrl(row)">{{ $t('auto.WhatsappCloudAccount.e73c73c1') }}</el-button>
            <el-button link type="primary" @click="openEdit(row)">{{ $t('auto.WhatsappCloudAccount.95b351c8') }}</el-button>
            <el-button link type="danger" @click="handleDelete(row)">{{ $t('auto.WhatsappCloudAccount.2f4aaddd') }}</el-button>
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

    
    <el-dialog v-model="dialogVisible" :title="dialogMode === 'edit' ? '编辑 WhatsApp Cloud 账号' : '新增 WhatsApp Cloud 账号'" width="640px">
      <el-form :model="form" label-width="150px">
        <el-form-item :label="$t('auto.WhatsappCloudAccount.d7ec2d3f')" prop="name">
          <el-input v-model="form.name" :placeholder="$t('auto.WhatsappCloudAccount.17327f9c')" />
        </el-form-item>
        <el-form-item :label="$t('auto.WhatsappCloudAccount.2432b575')" prop="remark">
          <el-input v-model="form.remark" :placeholder="$t('auto.WhatsappCloudAccount.67f6a89e')" />
        </el-form-item>
        <el-form-item label="PhoneNumberID" prop="phone_number_id">
          <el-input v-model="form.phone_number_id" :placeholder="$t('auto.WhatsappCloudAccount.8f4179c9')" />
        </el-form-item>
        <el-form-item label="WABA ID" prop="whatsapp_business_id">
          <el-input v-model="form.whatsapp_business_id" placeholder="WhatsApp Business Account ID" />
        </el-form-item>
        <el-form-item label="Access Token" prop="access_token">
          <el-input v-model="form.access_token" type="password" show-password :placeholder="$t('auto.WhatsappCloudAccount.e466a877')" />
        </el-form-item>
        <el-form-item label="App Secret" prop="app_secret">
          <el-input v-model="form.app_secret" type="password" show-password :placeholder="$t('auto.WhatsappCloudAccount.bb6463bb')" />
        </el-form-item>
        <el-form-item label="Verify Token" prop="verify_token">
          <el-input v-model="form.verify_token" :placeholder="$t('auto.WhatsappCloudAccount.16428c45')" />
          <div class="form-hint">{{ $t('auto.WhatsappCloudAccount.87a98781') }}</div>
        </el-form-item>

        <el-divider content-position="left">{{ $t('auto.WhatsappCloudAccount.f96de929') }}</el-divider>
        <el-form-item :label="$t('auto.WhatsappCloudAccount.5c7f1a99')" v-if="dialogMode === 'edit'">
          <el-input :model-value="callbackUrl(form.id)" readonly>
            <template #append>
              <el-button @click="copyText(callbackUrl(form.id))">{{ $t('auto.WhatsappCloudAccount.79d3abe9') }}</el-button>
            </template>
          </el-input>
          <div class="form-hint">{{ $t('auto.WhatsappCloudAccount.938169cf') }}</div>
        </el-form-item>
        <el-form-item :label="$t('auto.WhatsappCloudAccount.5f1e0e72')">
          <el-switch v-model="form.webhook_enabled" />
          <span class="form-hint-inline">{{ $t('auto.WhatsappCloudAccount.1aad39ea') }}</span>
        </el-form-item>
        <el-form-item :label="$t('auto.WhatsappCloudAccount.bfdd28c7')">
          <el-switch v-model="form.ai_agent_enabled" />
          <span class="form-hint-inline">{{ $t('auto.WhatsappCloudAccount.809880b9') }}</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.WhatsappCloudAccount.625fb26b') }}</el-button>
        <el-button type="primary" :loading="saving" @click="submit">{{ $t('auto.WhatsappCloudAccount.be5fbbe3') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="testVisible" :title="$t('auto.WhatsappCloudAccount.e1c652bc')" width="520px">
      <el-form :model="testForm" label-width="90px">
        <el-form-item :label="$t('auto.WhatsappCloudAccount.dc2d48d5')">
          <el-input v-model="testForm.to" :placeholder="$t('auto.WhatsappCloudAccount.50410759')" />
        </el-form-item>
        <el-form-item :label="$t('auto.WhatsappCloudAccount.b87b7756')">
          <el-input v-model="testForm.content" type="textarea" :rows="3" :placeholder="$t('auto.WhatsappCloudAccount.35af65ac')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="testVisible = false">{{ $t('auto.WhatsappCloudAccount.625fb26b') }}</el-button>
        <el-button type="primary" :loading="testing" @click="submitTestSend">{{ $t('auto.WhatsappCloudAccount.1535fcfa') }}</el-button>
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
    ElMessage.error(t('auto.WhatsappCloudAccount.c4fa9b74') + (e.message || e))
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
    ElMessage.error(t('auto.WhatsappCloudAccount.db870bff') + (e.message || e))
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
    ElMessage.success(t('auto.WhatsappCloudAccount.3b108349'))
    dialogVisible.value = false
    loadData()
  } catch (e) {
    ElMessage.error(t('auto.WhatsappCloudAccount.78d235cc') + (e.message || e))
  } finally {
    saving.value = false
  }
}

async function handleDelete(row) {
  try {
    await ElMessageBox.confirm(`确认删除账号「${row.name}」？`, t('auto.WhatsappCloudAccount.02d9819d'), { type: 'warning' })
  } catch {
    return
  }
  try {
    await deleteWhatsappCloud(row.id)
    ElMessage.success(t('auto.WhatsappCloudAccount.5cc23262'))
    loadData()
  } catch (e) {
    ElMessage.error(t('auto.WhatsappCloudAccount.bc42bb8f') + (e.message || e))
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
    ElMessage.warning(t('auto.WhatsappCloudAccount.6943cc55'))
    return
  }
  testing.value = true
  try {
    await testSendWhatsappCloud(testForm.id, { to: testForm.to, content: testForm.content })
    ElMessage.success(t('auto.WhatsappCloudAccount.9db9a7e3'))
    testVisible.value = false
  } catch (e) {
    ElMessage.error(t('auto.WhatsappCloudAccount.9704ea75') + (e.message || e))
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
    () => ElMessage.success(t('auto.WhatsappCloudAccount.4fb42e6e')),
    () => ElMessage.error(t('auto.WhatsappCloudAccount.41da1ca4'))
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
