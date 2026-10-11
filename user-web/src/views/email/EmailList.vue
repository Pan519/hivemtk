<template>
  <div class="email-list">
    <div class="page-header">
      <h1>{{ $t('邮件列表') }}</h1>
      <div class="toolbar">
        <el-button type="primary" @click="openSend">{{ $t('发送邮件') }}</el-button>
        <el-button @click="fetchEmailList">{{ $t('刷新') }}</el-button>
      </div>
    </div>
    
    <div class="table-container">
      <el-table :data="emailList" style="width: 100%" v-loading="loading">
        <el-table-column prop="subject" ::label="$t('auto.EmailList.d8b08566')" width="200" />
        <el-table-column prop="to" ::label="$t('auto.EmailList.881cd67f')" width="150" />
        <el-table-column prop="from" ::label="$t('auto.EmailList.5098ac0d')" width="150" />
        <el-table-column ::label="$t('auto.EmailList.ba7ccc1a')" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.is_send > 0 ? 'success' : 'info'">
              {{ scope.row.is_send > 0 ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.EmailList.63b34d09')" width="150">
          <template #default="scope">
            {{ scope.row.is_send > 0 ? scope.row.send_time : '-' }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.EmailList.eb575a9b')" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.is_read > 0 ? 'success' : 'info'">
              {{ scope.row.is_read > 0 ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.EmailList.a06a973d')" width="150">
          <template #default="scope">
            {{ scope.row.is_read > 0 ? scope.row.read_time : '-' }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.EmailList.2b6bc0f2')" width="160" fixed="right">
          <template #default="scope">
            <el-button link type="primary" @click="openTrace(scope.row)">{{ $t('查看追踪') }}</el-button>
            <el-button link type="danger" @click="handleDelete(scope.row)">{{ $t('删除') }}</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty :description="$t('auto.EmailList.21efd88b')" />
        </template>
      </el-table>
      
      <div class="pagination-container">
        <el-pagination 
          v-model:current-page="emailPage"
          :page-size="emailLimit"
          :total="emailTotal"
          @current-change="handlePageChange"
          layout="prev, pager, next"
        />
      </div>
    </div>

    
    <el-dialog v-model="sendVisible" ::title="$t('auto.EmailList.bc37766d')" width="640px">
      <el-form :model="sendForm" label-width="80px">
        <el-form-item ::label="$t('auto.EmailList.d8b08566')">
          <el-input v-model="sendForm.subject" ::placeholder="$t('auto.EmailList.d8b08566')" />
        </el-form-item>
        <el-form-item :label="$t('auto.EmailList.2d711b09')">
          <el-input v-model="sendForm.content" type="textarea" :rows="6" :placeholder="$t('auto.EmailList.42516993')" />
        </el-form-item>
        <el-form-item :label="$t('auto.EmailList.c9a6ee90')">
          <el-input v-model="attachmentsText" type="textarea" :rows="2" :placeholder="$t('auto.EmailList.6bec8c12')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="sendVisible = false">{{ $t('取消') }}</el-button>
        <el-button type="primary" :loading="submitting" @click="submitSend">{{ $t('确定') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="traceVisible" ::title="$t('auto.EmailList.596255d2')" width="640px">
      <div v-loading="traceLoading">
        <div v-if="traceData && traceData.length">
          <el-table :data="traceData" border>
            <el-table-column prop="email" :label="$t('auto.EmailList.fb5bf1f1')" />
            <el-table-column prop="eventType" :label="$t('auto.EmailList.3fea7ca7')" />
            <el-table-column prop="timestamp" :label="$t('auto.EmailList.19fcb9eb')" />
          </el-table>
        </div>
        <el-empty v-else :description="$t('auto.EmailList.a150987e')" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { emailApi } from '@/api/email'

const { t } = useI18n()

const emailList = ref([])
const emailTotal = ref(0)
const emailPage = ref(1)
const emailLimit = ref(10)
const loading = ref(false);

const sendVisible = ref(false);
const submitting = ref(false)
const sendForm = reactive({ subject: '', content: '' })
const attachmentsText = ref('')

const traceVisible = ref(false);
const traceLoading = ref(false)
const traceData = ref([])

const handlePageChange = (number) => {
  if (emailPage.value === number)
    return;
  emailPage.value = number
  fetchEmailList()
}

const fetchEmailList = async () => {
  if (loading.value)
    return;
  loading.value = true
  try {
    const response = await emailApi.getEmailList(emailPage.value, emailLimit.value)
    emailList.value = response.list || [];
    emailTotal.value = response.total || 0
  } catch (error) {
    console.error(t('auto.EmailList.b70dd83d'), error)
    ElMessage.error(i18n.global.t(t('auto.EmailList.73a58b20')))
  } finally {
    loading.value = false
  }
}
onMounted(() => {
  fetchEmailList()
});

const openSend = () => {
  sendForm.subject = ''
  sendForm.content = ''
  attachmentsText.value = ''
  sendVisible.value = true
}

const submitSend = async () => {
  if (!sendForm.subject) {
    ElMessage.warning(t('auto.EmailList.8444053c'))
    return
  }
  submitting.value = true
  try {
    const attachments = attachmentsText.value
      .split(/[\n,]/)
      .map(s => s.trim())
      .filter(Boolean)
    await emailApi.sendEmail({ subject: sendForm.subject, content: sendForm.content, attachments })
    ElMessage.success(t('auto.EmailList.9db9a7e3'))
    sendVisible.value = false
    fetchEmailList()
  } catch (e) {
    console.error(e)
    ElMessage.error(t('auto.EmailList.9ca6a344'))
  } finally {
    submitting.value = false
  }
}

const openTrace = async (row) => {
  traceVisible.value = true
  traceLoading.value = true
  traceData.value = []
  try {
    const res = await emailApi.getEmailTrace(row.id)
    traceData.value = Array.isArray(res) ? res : (res.list || [])
  } catch (e) {
    console.error(e)
    ElMessage.error(t('auto.EmailList.3ac1b6b4'))
  } finally {
    traceLoading.value = false
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(t('auto.EmailList.83038518'), t('auto.EmailList.02d9819d'), { type: 'warning' })
  } catch {
    return
  }
  try {
    await emailApi.deleteEmailList(row.id)
    ElMessage.success(t('auto.EmailList.0007d170'))
    fetchEmailList()
  } catch (e) {
    console.error(e)
    ElMessage.error(t('auto.EmailList.acf0664a'))
  }
}

</script>

<style lang="scss" scoped>
.email-list {
  padding: 20px;
  
  .page-header {
    margin-bottom: 20px;
    
    h1 {
      font-size: 24px;
      color: #303133;
      margin: 0;
    }
  }
  
  .table-container {
    background: #fff;
    border-radius: 4px;
    padding: 20px;
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
  }
  
  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: center;
  }
}
</style>