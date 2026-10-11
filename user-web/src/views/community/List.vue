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
        <el-form-item ::label="$t('auto.List.ae2511fc')">
          <el-input v-model="searchForm.name" ::placeholder="$t('auto.List.31e59f42')" clearable />
        </el-form-item>
        <el-form-item ::label="$t('auto.List.1b857887')">
          <el-select v-model="searchForm.status" ::placeholder="$t('auto.List.33329123')" clearable>
            <el-option ::label="$t('auto.List.62f11ef3')" value="1" />
            <el-option ::label="$t('auto.List.27ea4fcb')" value="2" />
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
      <el-table-column prop="name" ::label="$t('auto.List.ae2511fc')" min-width="150" show-overflow-tooltip />
      <el-table-column prop="description" ::label="$t('auto.List.d1ef7e0b')" min-width="200" show-overflow-tooltip />
      <el-table-column prop="member_count" ::label="$t('auto.List.d522011b')" width="100" />
      <el-table-column prop="message_count" ::label="$t('auto.List.ee6d1b85')" width="100" />
      <el-table-column prop="status" ::label="$t('auto.List.1b857887')" width="100">
        <template #default="scope">
          <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
            {{ scope.row.status === 1 ? '正常' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" :label="$t('auto.List.eca37cb0')" width="180" />
      <el-table-column :label="$t('auto.List.2b6bc0f2')" width="340" fixed="right">
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
        <el-form-item :label="$t('auto.List.1014b33d')" prop="name">
          <el-input v-model="form.name" :placeholder="$t('auto.List.0c64160a')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.3bdd08ad')" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="4" :placeholder="$t('auto.List.97bb062a')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.3fea7ca7')" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">{{ $t('auto.List.fd6e80f1') }}</el-radio>
            <el-radio :label="2">{{ $t('auto.List.710ad08b') }}</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">{{ $t('auto.List.625fb26b') }}</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>

    
    <el-dialog
      v-model="membersDialogVisible"
      :title="$t('auto.List.366cd9f4')"
      width="800px"
      top="5vh"
    >
      <el-table :data="memberList" border style="width: 100%" v-loading="membersLoading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="nickname" :label="$t('auto.List.23eb0e60')" min-width="150" />
        <el-table-column prop="role" :label="$t('auto.List.464f3d4e')" width="120" />
        <el-table-column prop="join_time" :label="$t('auto.List.8653cd3d')" width="180" />
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
      :title="$t('auto.List.cd19640f')"
      width="800px"
      top="5vh"
    >
      <el-table :data="messageList" border style="width: 100%" v-loading="messagesLoading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="sender" :label="$t('auto.List.413dc937')" width="150" />
        <el-table-column prop="content" :label="$t('auto.List.2d711b09')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="type" :label="$t('auto.List.226b0912')" width="100" />
        <el-table-column prop="created_at" :label="$t('auto.List.19fcb9eb')" width="180" />
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
      :title="$t('auto.List.70352f02')"
      width="700px"
    >
      <div v-loading="statsLoading" class="stats-content">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-card class="stats-card">
              <div class="stats-value">{{ currentStats.member_count || 0 }}</div>
              <div class="stats-label">{{ $t('auto.List.df2b42fe') }}</div>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card class="stats-card">
              <div class="stats-value">{{ currentStats.message_count || 0 }}</div>
              <div class="stats-label">{{ $t('auto.List.22825d3e') }}</div>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card class="stats-card">
              <div class="stats-value">{{ currentStats.active_count || 0 }}</div>
              <div class="stats-label">{{ $t('auto.List.a4fc7501') }}</div>
            </el-card>
          </el-col>
        </el-row>
        <el-descriptions :column="2" border style="margin-top: 20px">
          <el-descriptions-item :label="$t('auto.List.1014b33d')">{{ currentGroup.name }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.31fde0c9')">{{ currentGroup.id }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.eca37cb0')">{{ currentGroup.created_at }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.3fea7ca7')">{{ currentGroup.status === 1 ? '正常' : '禁用' }}</el-descriptions-item>
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
    { required: true, message: i18n.global.t(t('auto.List.0c64160a')), trigger: 'blur' },
    { min: 2, max: 50, message: i18n.global.t(t('auto.List.946a3653')), trigger: 'blur' }
  ]
};

const dialogTitle = computed(() => {
  return isEdit.value ? t('auto.List.53dfd8b6') : t('auto.List.bd0362be')
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
    t('auto.List.02d9819d'),
    {
      confirmButtonText: t('auto.List.38cf16f2'),
      cancelButtonText: t('auto.List.625fb26b'),
      type: 'warning'
    }
  ).then(async () => {
    try {
      await communityApi.deleteGroup(row.id)
      ElMessage.success(i18n.global.t(t('auto.List.0007d170')))
      fetchGroupList()
    } catch (error) {
      ElMessage.error(i18n.global.t(t('auto.List.acf0664a')))
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
        ElMessage.success(isEdit.value ? t('auto.List.55aa6366') : t('auto.List.04a691b3'))
        dialogVisible.value = false
        fetchGroupList()
      } catch (error) {
        ElMessage.error(isEdit.value ? t('auto.List.930442e2') : t('auto.List.a889286a'))
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
    ElMessage.success(i18n.global.t(t('auto.List.b6d16a81')))
    fetchGroupList()
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.List.fddcd7c6')))
    console.error(error)
  }
};

const handleExport = async () => {
  try {
    await communityApi.exportData({})
    ElMessage.success(i18n.global.t(t('auto.List.105c8a51')))
  } catch (error) {
    ElMessage.error(i18n.global.t('auto.List.dd51ab50'))
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
