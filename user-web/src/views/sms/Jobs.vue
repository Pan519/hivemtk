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
        <el-table-column prop="name" :label="$t('auto.Jobs.78caf711')" />
        <el-table-column prop="total" :label="$t('auto.Jobs.9ed7d3ad')" width="100" />
        <el-table-column prop="sent" :label="$t('auto.Jobs.93d15922')" width="100" />
        <el-table-column prop="failed" :label="$t('auto.Jobs.249f7b07')" width="100" />
        <el-table-column prop="status" :label="$t('auto.Jobs.3fea7ca7')" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="scheduleTime" :label="$t('auto.Jobs.f0c9d29d')" width="180">
          <template #default="scope">
            {{ scope.row.scheduleTime || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" :label="$t('auto.Jobs.eca37cb0')" width="180" />
        <el-table-column :label="$t('auto.Jobs.2b6bc0f2')" width="250">
          <template #default="scope">
            <el-button size="small" @click="handleView(scope.row)">{{ $t('auto.Jobs.607e7a4f') }}</el-button>
            <el-button size="small" type="warning" @click="handlePause(scope.row)" v-if="scope.row.status === 'running'">{{ $t('auto.Jobs.8d63ef38') }}</el-button>
            <el-button size="small" type="success" @click="handleResume(scope.row)" v-if="scope.row.status === 'paused'">{{ $t('auto.Jobs.27ca568b') }}</el-button>
            <el-button size="small" type="danger" @click="handleStop(scope.row)" v-if="scope.row.status === 'running' || scope.row.status === 'paused'">{{ $t('auto.Jobs.095e938e') }}</el-button>
            <el-button size="small" type="danger" @click="handleDelete(scope.row)" v-if="scope.row.status === 'completed' || scope.row.status === 'failed'">{{ $t('auto.Jobs.2f4aaddd') }}</el-button>
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

    
    <el-dialog v-model="dialogVisible" :title="$t('auto.Jobs.6ef95fe9')" width="600px" :before-close="handleDialogClose">
      <el-form :model="formData" label-width="100px" :rules="rules" ref="formRef">
        <el-form-item :label="$t('auto.Jobs.78caf711')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('auto.Jobs.64615c10')" maxlength="100" show-word-limit />
        </el-form-item>
        
        <el-form-item :label="$t('auto.Jobs.2b1806df')" prop="draftId">
          <el-select v-model="formData.draftId" :placeholder="$t('auto.Jobs.03b81bb8')">
            <el-option
              v-for="draft in draftOptions"
              :key="draft.id"
              :label="draft.title"
              :value="draft.id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item :label="$t('auto.Jobs.d9336549')" v-if="selectedDraft">
          <div class="content-preview">{{ selectedDraft.content }}</div>
        </el-form-item>
        
        <el-form-item :label="$t('auto.Jobs.8290e940')" prop="phoneList">
          <div class="phone-input-section">
            <el-input
              v-model="tempPhone"
              :placeholder="$t('auto.Jobs.6e4f4b33')"
              maxlength="11"
              @keyup.enter="addPhone"
            >
              <template #append>
                <el-button @click="addPhone">{{ $t('auto.Jobs.b58c7549') }}</el-button>
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
              :aria-label="$t('auto.Jobs.77f4314d')"
              @change="handleFileChange"
            />
            <!-- 链接态按钮要用独立的 link prop（见 components/AgentBindingDialog.vue:44 等既有写法），
                 type 的合法值是 default/primary/success/warning/info/danger/text/空串，不含 link；
                 写 type="link" 会触发 prop 校验告警并回落到 default，链接外观失效。 -->
            <el-button link size="small" @click="handleBatchUpload">{{ $t('auto.Jobs.c3202e4a') }}</el-button>
            <span class="tips-text">{{ $t('auto.Jobs.21fe0ce8') }}</span>
          </div>
        </el-form-item>
        
        <el-form-item :label="$t('auto.Jobs.6aa351f5')">
          <el-radio v-model="formData.sendType" label="immediate">{{ $t('auto.Jobs.1176c5b5') }}</el-radio>
          <el-radio v-model="formData.sendType" label="schedule">{{ $t('auto.Jobs.74e5c0d5') }}</el-radio>
        </el-form-item>
        
        <el-form-item :label="$t('auto.Jobs.4d9f36eb')" v-if="formData.sendType === 'schedule'" prop="scheduleTime">
          <el-date-picker
            v-model="formData.scheduleTime"
            type="datetime"
            :placeholder="$t('auto.Jobs.828c82fa')"
            value-format="YYYY-MM-DD HH:mm:ss"
            :disabled-date="disabledBeforeToday"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="handleDialogClose">{{ $t('auto.Jobs.625fb26b') }}</el-button>
          <el-button type="primary" @click="handleSubmit">{{ $t('auto.Jobs.38cf16f2') }}</el-button>
        </span>
      </template>
    </el-dialog>

    
    <el-dialog v-model="detailVisible" :title="$t('auto.Jobs.4a98bf0c')" width="700px">
      <el-descriptions :column="1" border>
        <el-descriptions-item :label="$t('auto.Jobs.78caf711')">{{ currentJob?.name }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.Jobs.9ed7d3ad')">{{ currentJob?.total }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.Jobs.93d15922')">{{ currentJob?.sent }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.Jobs.249f7b07')">{{ currentJob?.failed }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.Jobs.3fea7ca7')">
          <el-tag :type="getStatusType(currentJob?.status)">
            {{ getStatusText(currentJob?.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.Jobs.f0c9d29d')">{{ currentJob?.scheduleTime || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.Jobs.eca37cb0')">{{ currentJob?.createdAt }}</el-descriptions-item>
      </el-descriptions>
      
      <div class="job-detail-list" v-if="jobDetails.length > 0">
        <h3>{{ $t('auto.Jobs.4865e72a') }}</h3>
        <el-table :data="jobDetails" style="width: 100%">
          <el-table-column prop="phone" :label="$t('auto.Jobs.8098e2b4')" width="150" />
          <el-table-column prop="status" :label="$t('auto.Jobs.3fea7ca7')" width="100">
            <template #default="scope">
              <el-tag :type="getStatusType(scope.row.status)">
                {{ getStatusText(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="sendTime" :label="$t('auto.Jobs.63b34d09')" width="180" />
          <el-table-column prop="errorMsg" :label="$t('auto.Jobs.4604d502')" show-overflow-tooltip />
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
    { required: true, message: i18n.global.t(t('auto.Jobs.64615c10')), trigger: 'blur' },
    { min: 1, max: 100, message: i18n.global.t(t('auto.Jobs.01ecdfb1')), trigger: 'blur' }
  ],
  draftId: [
    { required: true, message: i18n.global.t(t('auto.Jobs.03b81bb8')), trigger: 'blur' }
  ],
  phoneList: [
    { required: true, message: i18n.global.t(t('auto.Jobs.f319668c')), trigger: 'blur' },
    { min: 1, message: i18n.global.t(t('auto.Jobs.f319668c')), trigger: 'blur' }
  ],
  scheduleTime: [
    { required: true, message: i18n.global.t(t('auto.Jobs.88824ad0')), trigger: 'change' }
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
    pending: t('auto.Jobs.6139a699'),
    running: t('auto.Jobs.46e38679'),
    paused: t('auto.Jobs.a2d930fd'),
    completed: t('auto.Jobs.fad5222c'),
    failed: t('auto.Jobs.1c83d797')
  }
  return statusMap[status] || t('auto.Jobs.1622dc9b')
};

const disabledBeforeToday = (time) => {
  return time.getTime() < Date.now() - 8.64e7;
};

const getDraftOptions = async () => {
  try {
    const response = await smsApi.getDraftList({ page: 1, limit: 100 })
    draftOptions.value = response.list || [];
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.Jobs.6adbf4da')))
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
    ElMessage.error(i18n.global.t(t('auto.Jobs.f71aceaf')))
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
    ElMessage.warning(i18n.global.t(t('auto.Jobs.6e4f4b33')))
    return
  }
  
  if (!/^1[3-9]\d{9}$/.test(phone)) {
    ElMessage.warning(i18n.global.t(t('auto.Jobs.a32ab517')))
    return
  }
  
  if (formData.phoneList.includes(phone)) {
    ElMessage.warning(i18n.global.t(t('auto.Jobs.d9426204')))
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
      ElMessageBox.alert(DOMPurify.sanitize(messages.join('<br>')), t('auto.Jobs.65804308'), {
        dangerouslyUseHTMLString: true
      })
    } else {
      ElMessage.warning(i18n.global.t(t('auto.Jobs.220948f2')))
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
    ElMessage.success(i18n.global.t(t('auto.Jobs.3799bdae')))
    
    dialogVisible.value = false
    getJobList()
  } catch (error) {
    if (error !== false) {
      ElMessage.error(t('auto.Jobs.a05f79f3') + (error.message || t('auto.Jobs.974e7484')))
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
    ElMessage.error(i18n.global.t(t('auto.Jobs.b1e72534')))
  }
};

const handlePause = async (row) => {
  try {
    await smsApi.pauseJob(row.id)
    ElMessage.success(i18n.global.t(t('auto.Jobs.60833eaa')))
    getJobList()
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.Jobs.fcc066a3')))
  }
};

const handleResume = async (row) => {
  try {
    await smsApi.resumeJob(row.id)
    ElMessage.success(i18n.global.t(t('auto.Jobs.8dbf03ae')))
    getJobList()
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.Jobs.d8cb5a1b')))
  }
};

const handleStop = async (row) => {
  try {
    await ElMessageBox.confirm(t('auto.Jobs.344da85c'), t('auto.Jobs.fb79c8d6'), {
      confirmButtonText: t('auto.Jobs.38cf16f2'),
      cancelButtonText: t('auto.Jobs.625fb26b'),
      type: 'warning'
    })
    await smsApi.stopJob(row.id)
    ElMessage.success(i18n.global.t(t('auto.Jobs.b20399d6')))
    getJobList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(i18n.global.t(t('auto.Jobs.ad89cba9')))
    }
  }
};

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(t('auto.Jobs.0c217b5b'), t('auto.Jobs.631cd220'), {
      confirmButtonText: t('auto.Jobs.38cf16f2'),
      cancelButtonText: t('auto.Jobs.625fb26b'),
      type: 'warning'
    })
    await smsApi.deleteJob(row.id)
    ElMessage.success(i18n.global.t(t('auto.Jobs.4aa7ab36')))
    getJobList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(i18n.global.t(t('auto.Jobs.0a70cd06')))
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