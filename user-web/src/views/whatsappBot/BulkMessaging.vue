<template>
  <div class="bulk-messaging">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('批量消息发送') }}</span>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane :label="$t('auto.BulkMessaging.1.b5f159')" name="send">
          <el-form :model="sendForm" :rules="sendRules" ref="sendFormRef" label-width="120px">
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item :label="$t('auto.BulkMessaging.2.579cc5')" prop="templateId">
                  <el-select 
                    v-model="sendForm.templateId" 
                    :placeholder="$t('auto.BulkMessaging.3.60debf')"
                    style="width: 100%"
                    @change="onTemplateChange"
                  >
                    <el-option
                      v-for="template in templates"
                      :key="template.id"
                      :label="template.name"
                      :value="template.id"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('auto.BulkMessaging.4.f4ff7f')" prop="accountId">
                  <el-select 
                    v-model="sendForm.accountId" 
                    :placeholder="$t('auto.BulkMessaging.5.0add9b')"
                    style="width: 100%"
                  >
                    <el-option
                      v-for="account in accounts"
                      :key="account.id"
                      :label="account.name"
                      :value="account.id"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item :label="$t('auto.BulkMessaging.6.9d7363')" prop="leadIds">
              <el-transfer
                v-model="sendForm.leadIds"
                filterable
                :filter-placeholder="$t('auto.BulkMessaging.7.259bea')"
                :data="transferData"
                ::titles="$t('auto.BulkMessaging.8.eb48d6')"
                ::button-texts="$t('auto.BulkMessaging.9.25f27e')"
                style="width: 100%"
              />
            </el-form-item>

            <el-form-item :label="$t('auto.BulkMessaging.10.740556')">
              <el-card shadow="never" class="preview-card">
                <div v-html="renderedContent"></div>
              </el-card>
            </el-form-item>

            <el-form-item :label="$t('auto.BulkMessaging.11.f0c9d2')">
              <el-date-picker
                v-model="sendForm.scheduleAt"
                type="datetime"
                :placeholder="$t('auto.BulkMessaging.12.900309')"
                format="YYYY-MM-DD HH:mm:ss"
                value-format="YYYY-MM-DD HH:mm:ss"
              />
            </el-form-item>

            <el-form-item>
              <el-button 
                type="primary" 
                @click="handleSend" 
                :loading="sending"
                :disabled="!canSend"
              >
                {{ sendForm.scheduleAt ? '计划发送' : '立即发送' }}
              </el-button>
              <el-button @click="resetForm">{{ $t('auto.BulkMessaging.35.4b9c32') }}</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane :label="$t('auto.BulkMessaging.2.579cc5')" name="templates">
          <div style="margin-bottom: 20px;">
            <el-button type="primary" @click="showTemplateDialog = true">{{ $t('auto.BulkMessaging.36.156a04') }}</el-button>
          </div>

          <el-table :data="templates" style="width: 100%" v-loading="templateLoading">
            <el-table-column prop="name" :label="$t('auto.BulkMessaging.13.a5d1c5')" width="200" />
            <el-table-column prop="content" :label="$t('auto.BulkMessaging.14.03ae79')" show-overflow-tooltip />
            <el-table-column prop="category" :label="$t('auto.BulkMessaging.15.d0771a')" width="120" />
            <el-table-column prop="isActive" :label="$t('auto.BulkMessaging.16.3fea7c')" width="100">
              <template #default="{ row }">
                <el-tag :type="row.isActive ? 'success' : 'danger'">
                  {{ row.isActive ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="updatedAt" :label="$t('auto.BulkMessaging.17.a001a2')" width="180" />
            <el-table-column :label="$t('auto.BulkMessaging.18.2b6bc0')" width="200">
              <template #default="{ row }">
                <el-button size="small" @click="editTemplate(row)">{{ $t('auto.BulkMessaging.37.95b351') }}</el-button>
                <el-button size="small" type="danger" @click="deleteTemplate(row)">{{ $t('auto.BulkMessaging.38.2f4aad') }}</el-button>
              </template>
            </el-table-column>
          </el-table>

          <el-pagination
            class="pagination"
            @size-change="handleTemplateSizeChange"
            @current-change="handleTemplateCurrentChange"
            :current-page="templatePagination.currentPage"
            :page-sizes="[10, 20, 50, 100]"
            :page-size="templatePagination.pageSize"
            layout="total, sizes, prev, pager, next, jumper"
            :total="templatePagination.total"
            v-if="templatePagination.total > 0"
          />
        </el-tab-pane>

        <el-tab-pane :label="$t('auto.BulkMessaging.19.4865e7')" name="records">
          <el-table :data="sendRecords" style="width: 100%" v-loading="recordLoading">
            <el-table-column prop="templateName" :label="$t('auto.BulkMessaging.13.a5d1c5')" width="200" />
            <el-table-column prop="totalCount" :label="$t('auto.BulkMessaging.20.8887e3')" width="100" />
            <el-table-column prop="sentCount" :label="$t('auto.BulkMessaging.21.93d159')" width="100" />
            <el-table-column prop="failedCount" :label="$t('auto.BulkMessaging.22.acd5cb')" width="100">
              <template #default="{ row }">
                <el-tag :type="row.failedCount > 0 ? 'danger' : 'success'">
                  {{ row.failedCount }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" :label="$t('auto.BulkMessaging.16.3fea7c')" width="120">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)">
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="createdAt" :label="$t('auto.BulkMessaging.23.eca37c')" width="180" />
            <el-table-column :label="$t('auto.BulkMessaging.18.2b6bc0')" width="150">
              <template #default="{ row }">
                <el-button size="small" @click="viewRecordDetails(row)">{{ $t('auto.BulkMessaging.39.f26225') }}</el-button>
              </template>
            </el-table-column>
          </el-table>

          <el-pagination
            class="pagination"
            @size-change="handleRecordSizeChange"
            @current-change="handleRecordCurrentChange"
            :current-page="recordPagination.currentPage"
            :page-sizes="[10, 20, 50, 100]"
            :page-size="recordPagination.pageSize"
            layout="total, sizes, prev, pager, next, jumper"
            :total="recordPagination.total"
            v-if="recordPagination.total > 0"
          />
        </el-tab-pane>
      </el-tabs>
    </el-card>

    
    <el-dialog :title="templateDialogTitle" v-model="showTemplateDialog" width="60%">
      <el-form :model="templateForm" :rules="templateRules" ref="templateFormRef" label-width="100px">
        <el-form-item :label="$t('auto.BulkMessaging.13.a5d1c5')" prop="name">
          <el-input v-model="templateForm.name" :placeholder="$t('auto.BulkMessaging.24.8f21b9')" />
        </el-form-item>
        <el-form-item :label="$t('auto.BulkMessaging.25.b87b77')" prop="content">
          <el-input 
            v-model="templateForm.content" 
            type="textarea" 
            :rows="6"
            placeholder="请输入消息内容，可使用变量如：{{name}}、{{company}}等"
            maxlength="4096"
            show-word-limit
          />
          <div class="template-help">
            <p><strong>{{ $t('auto.BulkMessaging.40.11b591') }}</strong></p>
            <ul>
              <li><code>{{name}}</code>{{ $t('auto.BulkMessaging.41.f52891') }}</li>
              <li><code>{{phone}}</code>{{ $t('auto.BulkMessaging.42.24906c') }}</li>
              <li><code>{{email}}</code>{{ $t('auto.BulkMessaging.43.85160a') }}</li>
              <li><code>{{company}}</code>{{ $t('auto.BulkMessaging.44.cf423a') }}</li>
              <li><code>{{source}}</code>{{ $t('auto.BulkMessaging.45.2649a2') }}</li>
            </ul>
          </div>
        </el-form-item>
        <el-form-item :label="$t('auto.BulkMessaging.15.d0771a')" prop="category">
          <el-select v-model="templateForm.category" :placeholder="$t('auto.BulkMessaging.26.8bb820')">
            <el-option :label="$t('auto.BulkMessaging.27.102b37')" value="marketing" />
            <el-option :label="$t('auto.BulkMessaging.28.2039ae')" value="customer_service" />
            <el-option :label="$t('auto.BulkMessaging.29.d920c6')" value="notification" />
            <el-option :label="$t('auto.BulkMessaging.30.725965')" value="product" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.BulkMessaging.16.3fea7c')" prop="isActive">
          <el-switch v-model="templateForm.isActive" :active-text="$t('auto.BulkMessaging.31.7854b5')" :inactive-text="$t('auto.BulkMessaging.32.710ad0')" />
        </el-form-item>
        <el-form-item :label="$t('auto.BulkMessaging.33.3bdd08')">
          <el-input 
            v-model="templateForm.description" 
            type="textarea" 
            :rows="3"
            :placeholder="$t('auto.BulkMessaging.34.ace725')"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showTemplateDialog = false">{{ $t('auto.BulkMessaging.46.625fb2') }}</el-button>
          <el-button type="primary" @click="saveTemplate">{{ $t('auto.BulkMessaging.47.38cf16') }}</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import DOMPurify from 'dompurify';
import * as bulkMessagingApi from '@/api/bulkMessaging'
import * as clueApi from '@/api/clue'
import * as whatsappApi from '@/api/whatsapp'
import { toList } from '@/utils/list'

const { t } = useI18n()

const activeTab = ref('send')
const sending = ref(false)
const templateLoading = ref(false)
const recordLoading = ref(false)
const showTemplateDialog = ref(false)
const templateDialogTitle = ref('')
const sendFormRef = ref()
const templateFormRef = ref()

const sendForm = reactive({
  templateId: '',
  accountId: '',
  leadIds: [],
  scheduleAt: null
});

const sendRules = {
  templateId: [{ required: true, message: i18n.global.t(t('auto.BulkMessaging.3.60debf')), trigger: 'change' }],
  accountId: [{ required: true, message: i18n.global.t(t('auto.BulkMessaging.5.0add9b')), trigger: 'change' }],
  leadIds: [{ required: true, message: i18n.global.t(t('auto.BulkMessaging.48.cd4e48')), trigger: 'change' }]
}

const templateForm = reactive({
  id: '',
  name: '',
  content: '',
  category: 'marketing',
  isActive: true,
  description: ''
});

const templateRules = {
  name: [{ required: true, message: i18n.global.t(t('auto.BulkMessaging.24.8f21b9')), trigger: 'blur' }],
  content: [{ required: true, message: i18n.global.t(t('auto.BulkMessaging.49.fa967b')), trigger: 'blur' }]
}

const templates = ref([]);
const accounts = ref([])
const allLeads = ref([])
const sendRecords = ref([])
const renderedContent = ref('')

const templatePagination = reactive({
  currentPage: 1,
  pageSize: 20,
  total: 0
});

const recordPagination = reactive({
  currentPage: 1,
  pageSize: 20,
  total: 0
})

const transferData = computed(() => {
  return allLeads.value.map(lead => ({
    key: lead.id,
    label: `${lead.name} (${lead.phone})`,
    disabled: false
  }))
});

const canSend = computed(() => {
  return sendForm.templateId && sendForm.accountId && sendForm.leadIds.length > 0
})

onMounted(() => {
  loadTemplates()
  loadAccounts()
  loadLeads()
  loadSendRecords()
})

const loadTemplates = async () => {
  try {
    templateLoading.value = true
    const response = await bulkMessagingApi.getTemplates()
    templates.value = toList(response)
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.BulkMessaging.50.801541')))
  } finally {
    templateLoading.value = false
  }
}

const loadAccounts = async () => {
  try {
    const response = await whatsappApi.getAccounts({ limit: 100 })
    accounts.value = response || [];
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.BulkMessaging.51.9e2a3a')))
  }
}

const loadLeads = async () => {
  try {
    const response = await clueApi.clueApi.list(1, 1000);
    const clueData = response?.list || response || [];
    allLeads.value = clueData.map(clue => ({
      id: clue.ID,
      name: clue.Name || clue.name || t('auto.BulkMessaging.52.1622dc'),
      phone: clue.Account || clue.account || '',
      email: '',
      company: clue.Address || clue.address || clue.City || clue.city || '',
      source: 'clue',
      score: 80,
      status: 'new'
    }))
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.BulkMessaging.53.9360c6')))
  }
}

const loadSendRecords = async () => {
  try {
    recordLoading.value = true
    const response = await bulkMessagingApi.getSendRecords({
      page: recordPagination.currentPage,
      limit: recordPagination.pageSize
    })
    sendRecords.value = toList(response);
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.BulkMessaging.54.a2a7f3')))
  } finally {
    recordLoading.value = false
  }
}

const onTemplateChange = (templateId) => {
  const template = templates.value.find(t => t.id === templateId)
  if (template) {
    const rawHtml = template.content
      .replace('{{name}}', t('auto.BulkMessaging.55.615db5'))
      .replace('{{phone}}', '13800138000')
      .replace('{{email}}', 'zhangsan@example.com')
      .replace('{{company}}', t('auto.BulkMessaging.56.5c95e3'))
      .replace('{{source}}', t('auto.BulkMessaging.57.5ff495'));
    renderedContent.value = DOMPurify.sanitize(rawHtml, { USE_PROFILES: { html: true } })
  }
}

const handleSend = async () => {
  try {
    await sendFormRef.value.validate()
    
    sending.value = true
    
    const response = await bulkMessagingApi.sendBulkMessage(sendForm)
    
    ElMessage.success(response.message || t('auto.BulkMessaging.58.fe785b'))
    
    resetForm();
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.BulkMessaging.59.9ca6a3')))
  } finally {
    sending.value = false
  }
}

const resetForm = () => {
  sendForm.templateId = ''
  sendForm.accountId = ''
  sendForm.leadIds = []
  sendForm.scheduleAt = null
}

const editTemplate = (row) => {
  Object.assign(templateForm, row)
  templateDialogTitle.value = t('auto.BulkMessaging.60.c6aa35')
  showTemplateDialog.value = true
}

const deleteTemplate = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要删除模板 "${row.name}" 吗？`, t('auto.BulkMessaging.61.631cd2'), {
      confirmButtonText: t('auto.BulkMessaging.47.38cf16'),
      cancelButtonText: t('auto.BulkMessaging.46.625fb2'),
      type: 'warning'
    })
    
    await bulkMessagingApi.deleteTemplate(row.id)
    ElMessage.success(i18n.global.t(t('auto.BulkMessaging.62.15e7d5')))
    loadTemplates()
  } catch (error) { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
}

const saveTemplate = async () => {
  try {
    await templateFormRef.value.validate()
    
    if (templateForm.id) {
      await bulkMessagingApi.updateTemplate(templateForm.id, templateForm)
      ElMessage.success(i18n.global.t(t('auto.BulkMessaging.63.b846b3')))
    } else {
      await bulkMessagingApi.createTemplate(templateForm)
      ElMessage.success(i18n.global.t(t('auto.BulkMessaging.64.49c62b')))
    }
    
    showTemplateDialog.value = false
    resetTemplateForm()
    loadTemplates()
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.BulkMessaging.65.5fa802')))
  }
}

const resetTemplateForm = () => {
  templateForm.id = ''
  templateForm.name = ''
  templateForm.content = ''
  templateForm.category = 'marketing'
  templateForm.isActive = true
  templateForm.description = ''
}

const viewRecordDetails = (row) => {
  ElMessageBox.alert(DOMPurify.sanitize(`
    <div style="text-align: left;">
      <p><strong>模板名称:</strong> ${row.templateName}</p>
      <p><strong>总数量:</strong> ${row.totalCount}</p>
      <p><strong>已发送:</strong> ${row.sentCount}</p>
      <p><strong>失败:</strong> ${row.failedCount}</p>
      <p><strong>状态:</strong> ${getStatusText(row.status)}</p>
      <p><strong>创建时间:</strong> ${row.createdAt}</p>
    </div>
  `), t('auto.BulkMessaging.66.f18912'), {
    dangerouslyUseHTMLString: true,
    customClass: 'record-detail-box'
  });
}

const getStatusType = (status) => {
  switch (status) {
    case 'completed': return 'success'
    case 'timeout': return 'warning'
    case 'failed': return 'danger'
    case 'pending': return 'info'
    default: return 'info'
  }
}

const getStatusText = (status) => {
  switch (status) {
    case 'completed': return t('auto.BulkMessaging.67.fad522')
    case 'timeout': return t('auto.BulkMessaging.68.0e7b6c')
    case 'failed': return t('auto.BulkMessaging.22.acd5cb')
    case 'pending': return t('auto.BulkMessaging.69.57a3b6')
    default: return status
  }
}

const handleTemplateSizeChange = (val) => {
  templatePagination.pageSize = val
  loadTemplates()
}

const handleTemplateCurrentChange = (val) => {
  templatePagination.currentPage = val
  loadTemplates()
}

const handleRecordSizeChange = (val) => {
  recordPagination.pageSize = val
  loadSendRecords()
}

const handleRecordCurrentChange = (val) => {
  recordPagination.currentPage = val
  loadSendRecords()
}
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}

.preview-card {
  background-color: #f5f5f5;
  padding: 15px;
  border-radius: 4px;
}

.template-help {
  margin-top: 10px;
  padding: 10px;
  background-color: #f0f9ff;
  border-radius: 4px;
}

.template-help ul {
  margin: 5px 0;
  padding-left: 20px;
}

.template-help code {
  background-color: #eee;
  padding: 2px 4px;
  border-radius: 3px;
  font-family: monospace;
}

:deep(.record-detail-box) {
  text-align: left;
}
</style>
