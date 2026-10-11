<template>
  <div class="community-container">
    
    <div class="page-header">
      <h2>{{ $t('社群管理') }}</h2>
      <div class="action-buttons">
        <el-button type="success" @click="handleExport">
          <el-icon><Download /></el-icon>
          {{ $t('导出') }}
        </el-button>
        <el-button type="warning" @click="handleImport">
          <el-icon><Upload /></el-icon>
          {{ $t('导入') }}
        </el-button>
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          {{ $t('新增分组') }}
        </el-button>
      </div>
    </div>

    
    <div class="search-form">
      <el-form :inline="true" :model="searchForm" class="search-form-content">
        <el-form-item ::label="$t('auto.List.1.ae2511')">
          <el-input v-model="searchForm.name" ::placeholder="$t('auto.List.2.31e59f')" clearable />
        </el-form-item>
        <el-form-item ::label="$t('auto.List.3.1b8578')">
          <el-select v-model="searchForm.status" ::placeholder="$t('auto.List.4.333291')" clearable>
            <el-option ::label="$t('auto.List.5.62f11e')" value="1" />
            <el-option ::label="$t('auto.List.6.27ea4f')" value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            {{ $t('搜索') }}
          </el-button>
          <el-button @click="resetSearch">
            <el-icon><RefreshRight /></el-icon>
            {{ $t('重置') }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    
    <el-table :data="groupList" border style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" ::label="$t('auto.List.1.ae2511')" min-width="150" show-overflow-tooltip />
      <el-table-column prop="description" ::label="$t('auto.List.7.d1ef7e')" min-width="200" show-overflow-tooltip />
      <el-table-column prop="member_count" ::label="$t('auto.List.8.d52201')" width="100" />
      <el-table-column prop="message_count" ::label="$t('auto.List.9.ee6d1b')" width="100" />
      <el-table-column prop="status" ::label="$t('auto.List.3.1b8578')" width="100">
        <template #default="scope">
          <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
            {{ scope.row.status === 1 ? '正常' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" :label="$t('auto.List.10.eca37c')" width="180" />
      <el-table-column :label="$t('auto.List.11.2b6bc0')" width="340" fixed="right">
        <template #default="scope">
          <el-button type="primary" size="small" @click="handleEdit(scope.row)">
            <el-icon><Edit /></el-icon>
            编辑
          </el-button>
          <el-button type="info" size="small" @click="handleViewMembers(scope.row)">
            <el-icon><User /></el-icon>
            成员
          </el-button>
          <el-button type="success" size="small" @click="handleViewMessages(scope.row)">
            <el-icon><ChatDotSquare /></el-icon>
            消息
          </el-button>
          <el-button type="warning" size="small" @click="handleViewStats(scope.row)">
            <el-icon><DataAnalysis /></el-icon>
            统计
          </el-button>
          <el-button type="danger" size="small" @click="handleDelete(scope.row)">
            <el-icon><Delete /></el-icon>
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    
    <div class="pagination-container">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="pagination.total"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item :label="$t('auto.List.12.1014b3')" prop="name">
          <el-input v-model="form.name" :placeholder="$t('auto.List.13.0c6416')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.14.3bdd08')" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="4" :placeholder="$t('auto.List.15.97bb06')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.16.3fea7c')" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">{{ $t('auto.List.28.fd6e80') }}</el-radio>
            <el-radio :label="2">{{ $t('auto.List.29.710ad0') }}</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">{{ $t('auto.List.30.625fb2') }}</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>

    
    <el-dialog
      v-model="membersDialogVisible"
      :title="$t('auto.List.17.366cd9')"
      width="800px"
      top="5vh"
    >
      <el-table :data="memberList" border style="width: 100%" v-loading="membersLoading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="nickname" :label="$t('auto.List.18.23eb0e')" min-width="150" />
        <el-table-column prop="role" :label="$t('auto.List.19.464f3d')" width="120" />
        <el-table-column prop="join_time" :label="$t('auto.List.20.8653cd')" width="180" />
      </el-table>
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="membersPagination.page"
          v-model:page-size="membersPagination.pageSize"
          :page-sizes="[10, 20, 50]"
          layout="total, prev, pager, next"
          :total="membersPagination.total"
          @current-change="handleMembersPageChange"
        />
      </div>
    </el-dialog>

    
    <el-dialog
      v-model="messagesDialogVisible"
      :title="$t('auto.List.21.cd1964')"
      width="800px"
      top="5vh"
    >
      <el-table :data="messageList" border style="width: 100%" v-loading="messagesLoading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="sender" :label="$t('auto.List.22.413dc9')" width="150" />
        <el-table-column prop="content" :label="$t('auto.List.23.2d711b')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="type" :label="$t('auto.List.24.226b09')" width="100" />
        <el-table-column prop="created_at" :label="$t('auto.List.25.19fcb9')" width="180" />
      </el-table>
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="messagesPagination.page"
          v-model:page-size="messagesPagination.pageSize"
          :page-sizes="[10, 20, 50]"
          layout="total, prev, pager, next"
          :total="messagesPagination.total"
          @current-change="handleMessagesPageChange"
        />
      </div>
    </el-dialog>

    
    <el-dialog
      v-model="statsDialogVisible"
      :title="$t('auto.List.26.70352f')"
      width="700px"
    >
      <div v-loading="statsLoading" class="stats-content">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-card class="stats-card">
              <div class="stats-value">{{ currentStats.member_count || 0 }}</div>
              <div class="stats-label">{{ $t('auto.List.31.df2b42') }}</div>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card class="stats-card">
              <div class="stats-value">{{ currentStats.message_count || 0 }}</div>
              <div class="stats-label">{{ $t('auto.List.32.22825d') }}</div>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card class="stats-card">
              <div class="stats-value">{{ currentStats.active_count || 0 }}</div>
              <div class="stats-label">{{ $t('auto.List.33.a4fc75') }}</div>
            </el-card>
          </el-col>
        </el-row>
        <el-descriptions :column="2" border style="margin-top: 20px">
          <el-descriptions-item :label="$t('auto.List.12.1014b3')">{{ currentGroup.name }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.27.31fde0')">{{ currentGroup.id }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.10.eca37c')">{{ currentGroup.created_at }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.16.3fea7c')">{{ currentGroup.status === 1 ? '正常' : '禁用' }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus, Search, RefreshRight, Edit, Delete, User, ChatDotSquare,
  DataAnalysis, Download, Upload
} from '@element-plus/icons-vue'
import { communityApi } from '@/api/community'

const { t } = useI18n()

const loading = ref(false);
const submitting = ref(false)
const membersLoading = ref(false)
const messagesLoading = ref(false)
const statsLoading = ref(false)
const groupList = ref([])
const dialogVisible = ref(false)
const membersDialogVisible = ref(false)
const messagesDialogVisible = ref(false)
const statsDialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)
const currentGroup = ref({})
const memberList = ref([])
const messageList = ref([])
const currentStats = ref({})

const searchForm = reactive({
  name: '',
  status: ''
});

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
});

const membersPagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
});

const messagesPagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
});

const form = reactive({
  id: null,
  name: '',
  description: '',
  status: 1
});

const rules = {
  name: [
    { required: true, message: i18n.global.t(t('auto.List.13.0c6416')), trigger: 'blur' },
    { min: 2, max: 50, message: i18n.global.t(t('auto.List.34.946a36')), trigger: 'blur' }
  ]
};

const dialogTitle = computed(() => {
  return isEdit.value ? t('auto.List.35.53dfd8') : t('auto.List.36.bd0362')
});

onMounted(() => {
  fetchGroupList()
});

const fetchGroupList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      name: searchForm.name,
      status: searchForm.status ? parseInt(searchForm.status) : undefined
    }
    const res = await communityApi.getGroups(params)
    groupList.value = res.list || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
};

const handleSearch = () => {
  pagination.page = 1
  fetchGroupList()
}

const resetSearch = () => {
  searchForm.name = ''
  searchForm.status = ''
  pagination.page = 1
  fetchGroupList()
}

const handleSizeChange = (val) => {
  pagination.pageSize = val
  pagination.page = 1
  fetchGroupList()
}

const handleCurrentChange = (val) => {
  pagination.page = val
  fetchGroupList()
}

const handleAdd = () => {
  isEdit.value = false
  dialogVisible.value = true
  resetForm()
}

const handleEdit = (row) => {
  isEdit.value = true
  dialogVisible.value = true
  form.id = row.id
  form.name = row.name
  form.description = row.description
  form.status = row.status
}

const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除分组 "${row.name}" 吗？`,
    t('auto.List.37.02d981'),
    {
      confirmButtonText: t('auto.List.38.38cf16'),
      cancelButtonText: t('auto.List.30.625fb2'),
      type: 'warning'
    }
  ).then(async () => {
    try {
      await communityApi.deleteGroup(row.id)
      ElMessage.success(i18n.global.t(t('auto.List.39.0007d1')))
      fetchGroupList()
    } catch (error) {
      ElMessage.error(i18n.global.t(t('auto.List.40.acf066')))
      console.error(error)
    }
  }).catch((e) => {
    if (e !== 'cancel' && e !== 'close')
      throw e;
  })
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (isEdit.value) {
          await communityApi.updateGroup(form.id, form)
        } else {
          await communityApi.createGroup(form)
        }
        ElMessage.success(isEdit.value ? t('auto.List.41.55aa63') : t('auto.List.42.04a691'))
        dialogVisible.value = false
        fetchGroupList()
      } catch (error) {
        ElMessage.error(isEdit.value ? t('auto.List.43.930442') : t('auto.List.44.a88928'))
        console.error(error)
      } finally {
        submitting.value = false
      }
    }
  })
}

const handleDialogClose = () => {
  resetForm()
}

const resetForm = () => {
  form.id = null
  form.name = ''
  form.description = ''
  form.status = 1
  if (formRef.value) {
    formRef.value.resetFields()
  }
}

const handleViewMembers = (row) => {
  currentGroup.value = row
  membersDialogVisible.value = true
  membersPagination.page = 1
  fetchMembers()
};

const fetchMembers = async () => {
  membersLoading.value = true
  try {
    const res = await communityApi.getMembers({
      group_id: currentGroup.value.id,
      page: membersPagination.page,
      page_size: membersPagination.pageSize
    })
    memberList.value = res.list || []
    membersPagination.total = res.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    membersLoading.value = false
  }
}

const handleMembersPageChange = (val) => {
  membersPagination.page = val
  fetchMembers()
}

const handleViewMessages = (row) => {
  currentGroup.value = row
  messagesDialogVisible.value = true
  messagesPagination.page = 1
  fetchMessages()
};

const fetchMessages = async () => {
  messagesLoading.value = true
  try {
    const res = await communityApi.getMessages({
      group_id: currentGroup.value.id,
      page: messagesPagination.page,
      page_size: messagesPagination.pageSize
    })
    messageList.value = res.list || []
    messagesPagination.total = res.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    messagesLoading.value = false
  }
}

const handleMessagesPageChange = (val) => {
  messagesPagination.page = val
  fetchMessages()
}

const handleViewStats = async (row) => {
  currentGroup.value = row
  statsDialogVisible.value = true
  statsLoading.value = true
  try {
    const res = await communityApi.getStats({ group_id: row.id })
    currentStats.value = res || {}
  } catch (error) {
    console.error(error)
  } finally {
    statsLoading.value = false
  }
};

const handleImport = async () => {
  try {
    await communityApi.importData({})
    ElMessage.success(i18n.global.t(t('auto.List.45.b6d16a')))
    fetchGroupList()
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.46.fddcd7')))
    console.error(error)
  }
};

const handleExport = async () => {
  try {
    await communityApi.exportData({})
    ElMessage.success(i18n.global.t(t('auto.List.47.105c8a')))
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.48.dd51ab')))
    console.error(error)
  }
};
</script>

<style lang="scss" scoped>
.community-container {
  padding: 20px;

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h2 {
      margin: 0;
      font-size: 24px;
      color: #303133;
    }

    .action-buttons {
      display: flex;
      gap: 10px;
    }
  }

  .search-form {
    margin-bottom: 20px;
    padding: 15px;
    background-color: #f5f7fa;
    border-radius: 4px;

    .search-form-content {
      margin: 0;
    }
  }

  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: center;
  }

  .stats-content {
    .stats-card {
      text-align: center;
      padding: 10px 0;

      .stats-value {
        font-size: 24px;
        font-weight: bold;
        color: #4F46E5;
        margin-bottom: 5px;
      }

      .stats-label {
        font-size: 14px;
        color: #909399;
      }
    }
  }
}
</style>
