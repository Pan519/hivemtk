<template>
  <div class="sms-drafts">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('短信草稿') }}</span>
          <el-button type="primary" @click="handleCreateDraft" :loading="loading">{{ $t('新建草稿') }}</el-button>
        </div>
      </template>
      
      <el-table :data="draftList" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" :label="$t('auto.Drafts.32c65d8d')" />
        <el-table-column prop="content" :label="$t('auto.Drafts.2d711b09')" show-overflow-tooltip />
        <el-table-column prop="createdAt" :label="$t('auto.Drafts.eca37cb0')" width="180" />
        <el-table-column prop="updatedAt" :label="$t('auto.Drafts.a001a226')" width="180" />
        <el-table-column :label="$t('auto.Drafts.2b6bc0f2')" width="200">
          <template #default="scope">
            <el-button size="small" @click="handleEdit(scope.row)">{{ $t('auto.Drafts.95b351c8') }}</el-button>
            <el-button size="small" type="success" @click="handleSend(scope.row)">{{ $t('auto.Drafts.1535fcfa') }}</el-button>
            <el-button size="small" type="danger" @click="handleDelete(scope.row)">{{ $t('auto.Drafts.2f4aaddd') }}</el-button>
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

    
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px" :before-close="handleDialogClose">
      <el-form :model="formData" label-width="80px" :rules="rules" ref="formRef">
        <el-form-item :label="$t('auto.Drafts.32c65d8d')" prop="title">
          <el-input v-model="formData.title" :placeholder="$t('auto.Drafts.87eb9514')" maxlength="100" show-word-limit />
        </el-form-item>
        <el-form-item :label="$t('auto.Drafts.2d711b09')" prop="content">
          <el-input
            v-model="formData.content"
            type="textarea"
            :placeholder="$t('auto.Drafts.56ad8041')"
            :rows="6"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="handleDialogClose">{{ $t('auto.Drafts.625fb26b') }}</el-button>
          <el-button type="primary" @click="handleSubmit">{{ $t('auto.Drafts.38cf16f2') }}</el-button>
        </span>
      </template>
    </el-dialog>

    
    <el-dialog v-model="sendDialogVisible" :title="$t('auto.Drafts.2d728886')" width="400px">
      <el-form :model="sendForm" label-width="80px">
        <el-form-item :label="$t('auto.Drafts.8098e2b4')">
          <el-input v-model="sendForm.phone" :placeholder="$t('auto.Drafts.6e4f4b33')" maxlength="11" show-word-limit />
        </el-form-item>
        <el-form-item :label="$t('auto.Drafts.d9336549')">
          <div class="content-preview">{{ currentDraft?.content }}</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="sendDialogVisible = false">{{ $t('auto.Drafts.625fb26b') }}</el-button>
          <el-button type="primary" @click="handleConfirmSend" :loading="sending">{{ $t('auto.Drafts.3ceeb0bb') }}</el-button>
        </span>
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

const { t } = useI18n()

const draftList = ref([]);
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)
const sending = ref(false)
const currentDraft = ref(null)

const dialogVisible = ref(false);
const sendDialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref(null)
const editingId = ref(null)

const formData = reactive({
  title: '',
  content: ''
});

const sendForm = reactive({
  phone: ''
});

const rules = {
  title: [
    { required: true, message: i18n.global.t(t('auto.Drafts.87eb9514')), trigger: 'blur' },
    { min: 1, max: 100, message: i18n.global.t(t('auto.Drafts.8ff90ce9')), trigger: 'blur' }
  ],
  content: [
    { required: true, message: i18n.global.t(t('auto.Drafts.56ad8041')), trigger: 'blur' },
    { min: 1, max: 500, message: i18n.global.t(t('auto.Drafts.75f753d2')), trigger: 'blur' }
  ]
};

const getDraftList = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      limit: pageSize.value
    }
    const response = await smsApi.getDraftList(params)
    draftList.value = response.list || [];
    total.value = response.total || 0
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.Drafts.6adbf4da')))
  } finally {
    loading.value = false
  }
};

const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  formData.title = ''
  formData.content = ''
  editingId.value = null
};

const handleCreateDraft = () => {
  resetForm()
  dialogTitle.value = t('auto.Drafts.b127f2bb')
  dialogVisible.value = true
};

const handleEdit = async (row) => {
  try {
    const response = await smsApi.getDraftDetail(row.id)
    const data = response;
    formData.title = data.title
    formData.content = data.content
    editingId.value = data.id
    dialogTitle.value = t('auto.Drafts.e53cd8aa')
    dialogVisible.value = true
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.Drafts.249f2141')))
  }
};

const handleDialogClose = () => {
  resetForm()
  dialogVisible.value = false
};

const handleSubmit = async () => {
  try {
    if (!formRef.value) return
    await formRef.value.validate()
    
    const draftData = {
      title: formData.title,
      content: formData.content
    }
    
    if (editingId.value) {
      await smsApi.updateDraft(editingId.value, draftData);
      ElMessage.success(i18n.global.t(t('auto.Drafts.b7773a7e')))
    } else {
      await smsApi.createDraft(draftData);
      ElMessage.success(i18n.global.t(t('auto.Drafts.83167530')))
    }
    
    dialogVisible.value = false
    getDraftList()
  } catch (error) {
    if (error !== false) {
      ElMessage.error(t('auto.Drafts.79585f92') + (error.message || t('auto.Drafts.974e7484')))
    }
  }
};

const handleSend = (row) => {
  currentDraft.value = row
  sendForm.phone = ''
  sendDialogVisible.value = true
};

const handleConfirmSend = async () => {
  if (!sendForm.phone || sendForm.phone.length !== 11) {
    ElMessage.warning(i18n.global.t(t('auto.Drafts.a32ab517')))
    return
  }
  
  sending.value = true
  try {
    await smsApi.sendDraft(currentDraft.value.id, { phone: sendForm.phone })
    ElMessage.success(i18n.global.t(t('auto.Drafts.b3136072')))
    sendDialogVisible.value = false
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.Drafts.8ea25782')))
  } finally {
    sending.value = false
  }
};

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要删除草稿「${row.title}」吗？`, t('auto.Drafts.631cd220'), {
      confirmButtonText: t('auto.Drafts.38cf16f2'),
      cancelButtonText: t('auto.Drafts.625fb26b'),
      type: 'warning'
    })
    await smsApi.deleteDraft(row.id)
    ElMessage.success(i18n.global.t(t('auto.Drafts.341f75ac')))
    getDraftList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(i18n.global.t(t('auto.Drafts.d2b74afe')))
    }
  }
};

const handleSizeChange = (val) => {
  pageSize.value = val
  getDraftList()
};

const handleCurrentChange = (val) => {
  currentPage.value = val
  getDraftList()
};

onMounted(() => {
  getDraftList()
})
</script>

<style scoped>
.sms-drafts {
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
  min-height: 80px;
  white-space: pre-wrap;
  word-break: break-word;
  color: #606266;
}
</style>
