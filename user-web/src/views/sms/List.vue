<template>
  <div class="sms-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('短信列表') }}</span>
          <el-button type="primary" @click="handleSendSms">{{ $t('发送短信') }}</el-button>
        </div>
      </template>
      
      <div class="search-form">
        <el-form :inline="true" :model="searchForm">
          <el-form-item :label="$t('auto.List.8098e2b4')">
            <el-input v-model="searchForm.phone" :placeholder="$t('auto.List.6e4f4b33')" clearable />
          </el-form-item>
          <el-form-item :label="$t('auto.List.3fea7ca7')">
            <el-select v-model="searchForm.status" :placeholder="$t('auto.List.e1c965ef')" clearable>
              <el-option :label="$t('auto.List.57a3b675')" value="pending" />
              <el-option :label="$t('auto.List.702513ff')" value="sending" />
              <el-option :label="$t('auto.List.93d15922')" value="sent" />
              <el-option :label="$t('auto.List.9ca6a344')" value="failed" />
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('auto.List.63b34d09')">
            <el-date-picker
              v-model="searchForm.dateRange"
              type="daterange"
              :range-separator="$t('auto.List.981cbe31')"
              :start-placeholder="$t('auto.List.b44c0f33')"
              :end-placeholder="$t('auto.List.1d468be9')"
              value-format="YYYY-MM-DD"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">{{ $t('auto.List.e5f71fc3') }}</el-button>
            <el-button @click="handleReset">{{ $t('auto.List.4b9c3271') }}</el-button>
          </el-form-item>
        </el-form>
      </div>
      
      <el-table :data="smsList" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="phone" :label="$t('auto.List.8098e2b4')" />
        <el-table-column prop="content" :label="$t('auto.List.4e963cc0')" show-overflow-tooltip />
        <el-table-column prop="status" :label="$t('auto.List.3fea7ca7')" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sendTime" :label="$t('auto.List.63b34d09')" width="180" />
        <el-table-column :label="$t('auto.List.2b6bc0f2')" width="150">
          <template #default="scope">
            <el-button size="small" @click="handleView(scope.row)">{{ $t('auto.List.607e7a4f') }}</el-button>
            <el-button size="small" type="danger" @click="handleResend(scope.row)" v-if="scope.row.status === 'failed'">{{ $t('auto.List.61051e84') }}</el-button>
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

    
    <el-dialog v-model="sendDialogVisible" :title="$t('auto.List.2d728886')" width="500px">
      <el-form :model="sendForm" label-width="80px">
        <el-form-item :label="$t('auto.List.8098e2b4')">
          <el-input v-model="sendForm.phone" :placeholder="$t('auto.List.6e4f4b33')" maxlength="11" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.2d711b09')">
          <el-input v-model="sendForm.content" type="textarea" :rows="5" :placeholder="$t('auto.List.56ad8041')" maxlength="500" show-word-limit />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="sendDialogVisible = false">{{ $t('auto.List.625fb26b') }}</el-button>
        <el-button type="primary" :loading="sending" @click="submitSendSms">{{ $t('auto.List.3ceeb0bb') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import smsApi from '@/api/sms'
import { getOrderStatusLabel, getOrderStatusTagType } from '@/constants/orderStatus';

const { t } = useI18n()

const searchForm = ref({
  phone: '',
  status: '',
  dateRange: []
});

const smsList = ref([]);
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)

const getStatusType = (status) => {
  const statusMap = {
    pending: 'info',
    sending: 'warning',
    sent: 'success',
    failed: 'danger'
  }
  return statusMap[status] || 'info'
};

const getStatusText = (status) => {
  const statusMap = {
    pending: t('auto.List.57a3b675'),
    sending: t('auto.List.702513ff'),
    sent: t('auto.List.93d15922'),
    failed: t('auto.List.9ca6a344')
  }
  return statusMap[status] || t('auto.List.1622dc9b')
};

const getSmsList = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      limit: pageSize.value,
      phone: searchForm.value.phone,
      status: searchForm.value.status,
      startDate: searchForm.value.dateRange?.[0] || '',
      endDate: searchForm.value.dateRange?.[1] || ''
    }
    const response = await smsApi.getSmsList(params)
    smsList.value = response.list || [];
    total.value = response.total || 0
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.82bcb581')))
  } finally {
    loading.value = false
  }
};

const handleSearch = () => {
  currentPage.value = 1
  getSmsList()
};

const handleReset = () => {
  searchForm.value = {
    phone: '',
    status: '',
    dateRange: []
  }
  handleSearch()
};

const handleView = async (row) => {
  try {
    const response = await smsApi.getSmsDetail(row.id)
    ElMessageBox.alert(JSON.stringify(response, null, 2), t('auto.List.65cd4bd3'))
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.b48708dd')))
  }
};

const handleResend = async (row) => {
  try {
    await smsApi.resendSms(row.id)
    ElMessage.success(i18n.global.t(t('auto.List.97f07ad5')))
    getSmsList()
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.61376101')))
  }
};

const handleSizeChange = (val) => {
  pageSize.value = val
  getSmsList()
};

const handleCurrentChange = (val) => {
  currentPage.value = val
  getSmsList()
};

onMounted(() => {
  getSmsList()
})

const sendDialogVisible = ref(false);
const sending = ref(false)
const sendForm = reactive({ phone: '', content: '' })

const handleSendSms = () => {
  sendForm.phone = ''
  sendForm.content = ''
  sendDialogVisible.value = true
}

const submitSendSms = async () => {
  if (!/^1[3-9]\d{9}$/.test(sendForm.phone)) {
    ElMessage.warning(t('auto.List.a32ab517'))
    return
  }
  if (!sendForm.content.trim()) {
    ElMessage.warning(t('auto.List.56ad8041'))
    return
  }
  sending.value = true
  try {
    await smsApi.sendSms({ phone: sendForm.phone, content: sendForm.content })
    ElMessage.success(t('auto.List.b3136072'))
    sendDialogVisible.value = false
    getSmsList()
  } catch (error) {
    ElMessage.error(t('auto.List.8ea25782'))
  } finally {
    sending.value = false
  }
}
</script>

<style scoped>
.sms-list {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>