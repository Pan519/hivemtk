<template>
  <div class="sms-jobs">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('短信任务') }}</span>
          <el-button type="primary" @click="handleCreateJob" :loading="loading">{{ $t('创建任务') }}</el-button>
        </div>
      </template>
      
      <el-table :data="jobList" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" :label="$t('auto.Jobs.1.78caf7')" />
        <el-table-column prop="total" :label="$t('auto.Jobs.2.9ed7d3')" width="100" />
        <el-table-column prop="sent" :label="$t('auto.Jobs.3.93d159')" width="100" />
        <el-table-column prop="failed" :label="$t('auto.Jobs.4.249f7b')" width="100" />
        <el-table-column prop="status" :label="$t('auto.Jobs.5.3fea7c')" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="scheduleTime" :label="$t('auto.Jobs.6.f0c9d2')" width="180">
          <template #default="scope">
            {{ scope.row.scheduleTime || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" :label="$t('auto.Jobs.7.eca37c')" width="180" />
        <el-table-column :label="$t('auto.Jobs.8.2b6bc0')" width="250">
          <template #default="scope">
            <el-button size="small" @click="handleView(scope.row)">{{ $t('auto.Jobs.24.607e7a') }}</el-button>
            <el-button size="small" type="warning" @click="handlePause(scope.row)" v-if="scope.row.status === 'running'">{{ $t('auto.Jobs.25.8d63ef') }}</el-button>
            <el-button size="small" type="success" @click="handleResume(scope.row)" v-if="scope.row.status === 'paused'">{{ $t('auto.Jobs.26.27ca56') }}</el-button>
            <el-button size="small" type="danger" @click="handleStop(scope.row)" v-if="scope.row.status === 'running' || scope.row.status === 'paused'">{{ $t('auto.Jobs.27.095e93') }}</el-button>
            <el-button size="small" type="danger" @click="handleDelete(scope.row)" v-if="scope.row.status === 'completed' || scope.row.status === 'failed'">{{ $t('auto.Jobs.28.2f4aad') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    
    <el-dialog v-model="dialogVisible" :title="$t('auto.Jobs.9.6ef95f')" width="600px" :before-close="handleDialogClose">
      <el-form :model="formData" label-width="100px" :rules="rules" ref="formRef">
        <el-form-item :label="$t('auto.Jobs.1.78caf7')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('auto.Jobs.10.64615c')" maxlength="100" show-word-limit />
        </el-form-item>
        
        <el-form-item :label="$t('auto.Jobs.11.2b1806')" prop="draftId">
          <el-select v-model="formData.draftId" :placeholder="$t('auto.Jobs.12.03b81b')">
            <el-option
              v-for="draft in draftOptions"
              :key="draft.id"
              :label="draft.title"
              :value="draft.id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item :label="$t('auto.Jobs.13.d93365')" v-if="selectedDraft">
          <div class="content-preview">{{ selectedDraft.content }}</div>
        </el-form-item>
        
        <el-form-item :label="$t('auto.Jobs.14.8290e9')" prop="phoneList">
          <div class="phone-input-section">
            <el-input
              v-model="tempPhone"
              :placeholder="$t('auto.Jobs.15.6e4f4b')"
              maxlength="11"
              @keyup.enter="addPhone"
            >
              <template #append>
                <el-button @click="addPhone">{{ $t('auto.Jobs.29.b58c75') }}</el-button>
              </template>
            </el-input>
          </div>
          
          <div class="phone-list">
            <el-tag
              v-for="(phone, index) in formData.phoneList"
              :key="index"
              closable
              @close="removePhone(index)"
            >
              {{ phone }}
            </el-tag>
          </div>
          
          <div class="phone-tips">
            <input
              ref="fileInputRef"
              type="file"
              accept=".csv,.txt"
              style="display: none"
              :aria-label="$t('auto.Jobs.16.77f431')"
              @change="handleFileChange"
            />
            <!-- 链接态按钮要用独立的 link prop（见 components/AgentBindingDialog.vue:44 等既有写法），
                 type 的合法值是 default/primary/success/warning/info/danger/text/空串，不含 link；
                 写 type="link" 会触发 prop 校验告警并回落到 default，链接外观失效。 -->
            <el-button link size="small" @click="handleBatchUpload">{{ $t('auto.Jobs.30.c3202e') }}</el-button>
            <span class="tips-text">{{ $t('auto.Jobs.31.21fe0c') }}</span>
          </div>
        </el-form-item>
        
        <el-form-item :label="$t('auto.Jobs.17.6aa351')">
          <el-radio v-model="formData.sendType" label="immediate">{{ $t('auto.Jobs.32.1176c5') }}</el-radio>
          <el-radio v-model="formData.sendType" label="schedule">{{ $t('auto.Jobs.33.74e5c0') }}</el-radio>
        </el-form-item>
        
        <el-form-item :label="$t('auto.Jobs.18.4d9f36')" v-if="formData.sendType === 'schedule'" prop="scheduleTime">
          <el-date-picker
            v-model="formData.scheduleTime"
            type="datetime"
            :placeholder="$t('auto.Jobs.19.828c82')"
            value-format="YYYY-MM-DD HH:mm:ss"
            :disabled-date="disabledBeforeToday"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="handleDialogClose">{{ $t('auto.Jobs.34.625fb2') }}</el-button>
          <el-button type="primary" @click="handleSubmit">{{ $t('auto.Jobs.35.38cf16') }}</el-button>
        </span>
      </template>
    </el-dialog>

    
    <el-dialog v-model="detailVisible" :title="$t('auto.Jobs.20.4a98bf')" width="700px">
      <el-descriptions :column="1" border>
        <el-descriptions-item :label="$t('auto.Jobs.1.78caf7')">{{ currentJob?.name }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.Jobs.2.9ed7d3')">{{ currentJob?.total }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.Jobs.3.93d159')">{{ currentJob?.sent }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.Jobs.4.249f7b')">{{ currentJob?.failed }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.Jobs.5.3fea7c')">
          <el-tag :type="getStatusType(currentJob?.status)">
            {{ getStatusText(currentJob?.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.Jobs.6.f0c9d2')">{{ currentJob?.scheduleTime || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.Jobs.7.eca37c')">{{ currentJob?.createdAt }}</el-descriptions-item>
      </el-descriptions>
      
      <div class="job-detail-list" v-if="jobDetails.length > 0">
        <h3>{{ $t('auto.Jobs.36.4865e7') }}</h3>
        <el-table :data="jobDetails" style="width: 100%">
          <el-table-column prop="phone" :label="$t('auto.Jobs.21.8098e2')" width="150" />
          <el-table-column prop="status" :label="$t('auto.Jobs.5.3fea7c')" width="100">
            <template #default="scope">
              <el-tag :type="getStatusType(scope.row.status)">
                {{ getStatusText(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="sendTime" :label="$t('auto.Jobs.22.63b34d')" width="180" />
          <el-table-column prop="errorMsg" :label="$t('auto.Jobs.23.4604d5')" show-overflow-tooltip />
        </el-table>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import smsApi from '@/api/sms'
import { getOrderStatusLabel, getOrderStatusTagType } from '@/constants/orderStatus';
import DOMPurify from 'dompurify'
import { toList } from '@/utils/list'

const { t } = useI18n()

const jobList = ref([]);
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)

const dialogVisible = ref(false);
const detailVisible = ref(false)
const formRef = ref(null)
const currentJob = ref(null)
const jobDetails = ref([])

const draftOptions = ref([]);
const selectedDraft = computed(() => {
  return draftOptions.value.find(draft => draft.id === formData.draftId)
})

const formData = reactive({
  name: '',
  draftId: '',
  phoneList: [],
  sendType: 'immediate',
  scheduleTime: null
});

const tempPhone = ref('');

const rules = {
  name: [
    { required: true, message: i18n.global.t(t('auto.Jobs.10.64615c')), trigger: 'blur' },
    { min: 1, max: 100, message: i18n.global.t(t('auto.Jobs.37.01ecdf')), trigger: 'blur' }
  ],
  draftId: [
    { required: true, message: i18n.global.t(t('auto.Jobs.12.03b81b')), trigger: 'blur' }
  ],
  phoneList: [
    { required: true, message: i18n.global.t(t('auto.Jobs.38.f31966')), trigger: 'blur' },
    { min: 1, message: i18n.global.t(t('auto.Jobs.38.f31966')), trigger: 'blur' }
  ],
  scheduleTime: [
    { required: true, message: i18n.global.t(t('auto.Jobs.39.88824a')), trigger: 'change' }
  ]
};

const getStatusType = (status) => {
  const statusMap = {
    pending: 'info',
    running: 'success',
    paused: 'warning',
    completed: 'success',
    failed: 'danger'
  }
  return statusMap[status] || 'info'
};

const getStatusText = (status) => {
  const statusMap = {
    pending: t('auto.Jobs.40.6139a6'),
    running: t('auto.Jobs.41.46e386'),
    paused: t('auto.Jobs.42.a2d930'),
    completed: t('auto.Jobs.43.fad522'),
    failed: t('auto.Jobs.44.1c83d7')
  }
  return statusMap[status] || t('auto.Jobs.45.1622dc')
};

const disabledBeforeToday = (time) => {
  return time.getTime() < Date.now() - 8.64e7;
};

const getDraftOptions = async () => {
  try {
    const response = await smsApi.getDraftList({ page: 1, limit: 100 })
    draftOptions.value = response.list || [];
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.Jobs.46.6adbf4')))
  }
};

const getJobList = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      limit: pageSize.value
    }
    const response = await smsApi.getJobList(params)
    jobList.value = response.list || [];
    total.value = response.total || 0
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.Jobs.47.f71ace')))
  } finally {
    loading.value = false
  }
};

const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  formData.name = ''
  formData.draftId = ''
  formData.phoneList = []
  formData.sendType = 'immediate'
  formData.scheduleTime = null
  tempPhone.value = ''
};

const addPhone = () => {
  const phone = tempPhone.value.trim()
  if (!phone) {
    ElMessage.warning(i18n.global.t(t('auto.Jobs.15.6e4f4b')))
    return
  }
  
  if (!/^1[3-9]\d{9}$/.test(phone)) {
    ElMessage.warning(i18n.global.t(t('auto.Jobs.48.a32ab5')))
    return
  }
  
  if (formData.phoneList.includes(phone)) {
    ElMessage.warning(i18n.global.t(t('auto.Jobs.49.d94262')))
    return
  }
  
  formData.phoneList.push(phone)
  tempPhone.value = ''
};

const removePhone = (index) => {
  formData.phoneList.splice(index, 1)
};

const fileInputRef = ref(null);

const handleBatchUpload = () => {
  fileInputRef.value?.click()
};

const handleFileChange = (event) => {
  const file = event.target.files?.[0]
  if (!file) return

  const reader = new FileReader()
  reader.onload = (e) => {
    const content = e.target?.result
    if (typeof content !== 'string') return

    const rawNumbers = content.split(/[\r\n,;]+/).map(s => s.trim()).filter(Boolean);

    const validPhones = []
    const invalidPhones = []
    const duplicatePhones = []
    const existingSet = new Set(formData.phoneList)

    for (const raw of rawNumbers) {
      if (!/^1[3-9]\d{9}$/.test(raw)) {
        invalidPhones.push(raw)
        continue
      }
      if (existingSet.has(raw) || validPhones.includes(raw)) {
        duplicatePhones.push(raw)
        continue
      }
      validPhones.push(raw)
      existingSet.add(raw)
    }

    if (validPhones.length > 0) {
      formData.phoneList.push(...validPhones)
    }

    const messages = []
    if (validPhones.length > 0) {
      messages.push(`成功导入 ${validPhones.length} 个手机号`)
    }
    if (duplicatePhones.length > 0) {
      messages.push(`已去除 ${duplicatePhones.length} 个重复号码`)
    }
    if (invalidPhones.length > 0) {
      const preview = invalidPhones.slice(0, 5).join(', ')
      const suffix = invalidPhones.length > 5 ? `... 等共 ${invalidPhones.length} 个` : ''
      messages.push(`发现 ${invalidPhones.length} 个无效号码: ${preview}${suffix}`)
    }

    if (messages.length > 0) {
      ElMessageBox.alert(DOMPurify.sanitize(messages.join('<br>')), t('auto.Jobs.50.658043'), {
        dangerouslyUseHTMLString: true
      })
    } else {
      ElMessage.warning(i18n.global.t(t('auto.Jobs.51.220948')))
    }
  }
  reader.readAsText(file)

  event.target.value = '';
};

const handleCreateJob = () => {
  resetForm()
  getDraftOptions()
  dialogVisible.value = true
};

const handleDialogClose = () => {
  resetForm()
  dialogVisible.value = false
};

const handleSubmit = async () => {
  try {
    if (!formRef.value) return
    await formRef.value.validate()
    
    const jobData = {
      name: formData.name,
      phoneList: formData.phoneList,
      content: selectedDraft.value.content
    }
    
    if (formData.sendType === 'schedule' && formData.scheduleTime) {
      jobData.scheduleTime = formData.scheduleTime
    }
    
    await smsApi.createJob(jobData)
    ElMessage.success(i18n.global.t(t('auto.Jobs.52.3799bd')))
    
    dialogVisible.value = false
    getJobList()
  } catch (error) {
    if (error !== false) {
      ElMessage.error(t('auto.Jobs.53.a05f79') + (error.message || t('auto.Jobs.54.974e74')))
    }
  }
};

const handleView = async (row) => {
  try {
    const response = await smsApi.getJobDetail(row.id)
    currentJob.value = response
    try {
      const recordsRes = await smsApi.getJobRecords(row.id, { page: 1, page_size: 50 })
      const recordsData = toList(recordsRes)
      jobDetails.value = Array.isArray(recordsData) ? recordsData : []
    } catch (e) {
      jobDetails.value = []
    }
    detailVisible.value = true
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.Jobs.55.b1e725')))
  }
};

const handlePause = async (row) => {
  try {
    await smsApi.pauseJob(row.id)
    ElMessage.success(i18n.global.t(t('auto.Jobs.56.60833e')))
    getJobList()
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.Jobs.57.fcc066')))
  }
};

const handleResume = async (row) => {
  try {
    await smsApi.resumeJob(row.id)
    ElMessage.success(i18n.global.t(t('auto.Jobs.58.8dbf03')))
    getJobList()
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.Jobs.59.d8cb5a')))
  }
};

const handleStop = async (row) => {
  try {
    await ElMessageBox.confirm(t('auto.Jobs.60.344da8'), t('auto.Jobs.61.fb79c8'), {
      confirmButtonText: t('auto.Jobs.35.38cf16'),
      cancelButtonText: t('auto.Jobs.34.625fb2'),
      type: 'warning'
    })
    await smsApi.stopJob(row.id)
    ElMessage.success(i18n.global.t(t('auto.Jobs.62.b20399')))
    getJobList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(i18n.global.t(t('auto.Jobs.63.ad89cb')))
    }
  }
};

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(t('auto.Jobs.64.0c217b'), t('auto.Jobs.65.631cd2'), {
      confirmButtonText: t('auto.Jobs.35.38cf16'),
      cancelButtonText: t('auto.Jobs.34.625fb2'),
      type: 'warning'
    })
    await smsApi.deleteJob(row.id)
    ElMessage.success(i18n.global.t(t('auto.Jobs.66.4aa7ab')))
    getJobList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(i18n.global.t(t('auto.Jobs.67.0a70cd')))
    }
  }
};

const handleSizeChange = (val) => {
  pageSize.value = val
  getJobList()
};

const handleCurrentChange = (val) => {
  currentPage.value = val
  getJobList()
};

onMounted(() => {
  getJobList()
})
</script>

<style scoped>
.sms-jobs {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.content-preview {
  background-color: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  min-height: 60px;
  white-space: pre-wrap;
  word-break: break-word;
  color: #606266;
}

.phone-input-section {
  margin-bottom: 10px;
}

.phone-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-bottom: 10px;
  max-height: 150px;
  overflow-y: auto;
  padding: 10px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.phone-tips {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #909399;
  font-size: 14px;
}

.tips-text {
  color: #909399;
}

.job-detail-list {
  margin-top: 20px;
}

.job-detail-list h3 {
  margin-bottom: 10px;
  font-size: 16px;
  color: #303133;
}
</style>