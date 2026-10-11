<template>
  <div class="user-segment-page">
    <el-card class="header-card">
      <div class="header-content">
        <h2>{{ $t('auto.List.29.79ed66') }}</h2>
        <p class="subtitle">{{ $t('auto.List.30.8a13c4') }}</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="showCreateDialog">
          <el-icon><Plus /></el-icon>
          {{ $t('创建分群') }}
        </el-button>
        <el-button @click="$router.push('/userSegment/builder')">
          <el-icon><SetUp /></el-icon>
          {{ $t('规则构建器') }}
        </el-button>
        <el-button @click="refreshData">
          <el-icon><Refresh /></el-icon>
          {{ $t('刷新') }}
        </el-button>
      </div>
    </el-card>

    
    <el-row :gutter="20" class="overview-row">
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card">
          <div class="stat-label">{{ $t('总用户数') }}</div>
          <div class="stat-value">{{ stats.totalUsers ?? stats.total ?? 0 }}</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card">
          <div class="stat-label">{{ $t('高价值客户') }}</div>
          <div class="stat-value" style="color: #10B981">{{ stats.highValue ?? stats.high_value ?? 0 }}</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card">
          <div class="stat-label">{{ $t('活跃客户') }}</div>
          <div class="stat-value" style="color: #4F46E5">{{ stats.active ?? stats.active_count ?? 0 }}</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card">
          <div class="stat-label">{{ $t('流失风险') }}</div>
          <div class="stat-value" style="color: #EF4444">{{ stats.churnRisk ?? stats.churn_risk ?? 0 }}</div>
        </el-card>
      </el-col>
    </el-row>

    
    <el-card class="matrix-card">
      <template #header>
        <span>{{ $t('auto.List.31.87e724') }}</span>
      </template>
      <div class="rfm-matrix">
        <table>
          <thead>
            <tr>
              <th></th>
              <th>{{ $t('auto.List.32.19ec90') }}</th>
              <th>{{ $t('auto.List.33.2ec604') }}</th>
              <th>{{ $t('auto.List.34.b758d3') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <th>{{ $t('auto.List.35.c67247') }}</th>
              <td class="cell high-value" @click="onCellClick('high_value')">
                <div class="cell-name">{{ $t('auto.List.36.66f84a') }}</div>
                <div class="cell-count">{{ matrixCounts.high_value || 0 }} 人</div>
              </td>
              <td class="cell good" @click="onCellClick('important_develop')">
                <div class="cell-name">{{ $t('auto.List.37.adda3b') }}</div>
                <div class="cell-count">{{ matrixCounts.important_develop || 0 }} 人</div>
              </td>
              <td class="cell warn" @click="onCellClick('important_keep')">
                <div class="cell-name">{{ $t('auto.List.38.1c57e2') }}</div>
                <div class="cell-count">{{ matrixCounts.important_keep || 0 }} 人</div>
              </td>
            </tr>
            <tr>
              <th>{{ $t('auto.List.39.f84a72') }}</th>
              <td class="cell good" @click="onCellClick('general_value')">
                <div class="cell-name">{{ $t('auto.List.40.6ecb76') }}</div>
                <div class="cell-count">{{ matrixCounts.general_value || 0 }} 人</div>
              </td>
              <td class="cell normal" @click="onCellClick('general_develop')">
                <div class="cell-name">{{ $t('auto.List.41.cef1bd') }}</div>
                <div class="cell-count">{{ matrixCounts.general_develop || 0 }} 人</div>
              </td>
              <td class="cell warn" @click="onCellClick('general_keep')">
                <div class="cell-name">{{ $t('auto.List.42.1d517c') }}</div>
                <div class="cell-count">{{ matrixCounts.general_keep || 0 }} 人</div>
              </td>
            </tr>
            <tr>
              <th>{{ $t('auto.List.43.b85c4d') }}</th>
              <td class="cell normal" @click="onCellClick('potential_value')">
                <div class="cell-name">{{ $t('auto.List.44.74ba3e') }}</div>
                <div class="cell-count">{{ matrixCounts.potential_value || 0 }} 人</div>
              </td>
              <td class="cell warn" @click="onCellClick('potential_develop')">
                <div class="cell-name">{{ $t('auto.List.45.2df631') }}</div>
                <div class="cell-count">{{ matrixCounts.potential_develop || 0 }} 人</div>
              </td>
              <td class="cell low-value" @click="onCellClick('churn')">
                <div class="cell-name">{{ $t('auto.List.46.ddd6e7') }}</div>
                <div class="cell-count">{{ matrixCounts.churn || 0 }} 人</div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </el-card>

    
    <el-card class="segments-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.List.47.e204cd') }}</span>
          <el-input
            v-model="searchKeyword"
            :placeholder="$t('auto.List.1.a1bbf5')"
            style="width: 200px"
            clearable
          />
        </div>
      </template>
      <el-table :data="filteredSegments" v-loading="loading" stripe>
        <el-table-column prop="name" :label="$t('auto.List.2.33429b')" min-width="150" />
        <el-table-column prop="type" :label="$t('auto.List.3.65f306')" width="120">
          <template #default="{ row }">
            <el-tag :type="getTypeTagType(row.type)">{{ getTypeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="userCount" :label="$t('auto.List.4.71cb56')" width="100" />
        <el-table-column prop="criteria" :label="$t('auto.List.5.012c7f')" min-width="250" />
        <el-table-column prop="createdAt" :label="$t('auto.List.6.eca37c')" width="180" />
        <el-table-column prop="status" :label="$t('auto.List.7.3fea7c')" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)">{{ getStatusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.8.2b6bc0')" width="250" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="viewUsers(row)">{{ $t('auto.List.48.9519f4') }}</el-button>
            <el-button link type="primary" @click="editSegment(row)">{{ $t('auto.List.49.95b351') }}</el-button>
            <el-button link type="danger" @click="deleteSegment(row)">{{ $t('auto.List.50.2f4aad') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="700px">
      <el-form :model="form" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('auto.List.2.33429b')" prop="name">
          <el-input v-model="form.name" :placeholder="$t('auto.List.9.4bb33c')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.3.65f306')" prop="type">
          <el-select v-model="form.type" :placeholder="$t('auto.List.10.099cb7')" style="width: 100%">
            <el-option :label="$t('auto.List.11.ef764b')" value="RFM" />
            <el-option :label="$t('auto.List.12.f1d4ff')" value="CUSTOM" />
            <el-option :label="$t('auto.List.13.578fa4')" value="BEHAVIOR" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.14.787b74')" v-if="form.type === 'RFM'">
          <el-input-number v-model="form.rDays" :min="1" :max="365" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.15.b25268')" v-if="form.type === 'RFM'">
          <el-input-number v-model="form.fCount" :min="1" :max="100" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.16.ad8ce7')" v-if="form.type === 'RFM'">
          <el-input-number v-model="form.mAmount" :min="0" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.5.012c7f')" prop="criteria">
          <el-input
            v-model="form.criteria"
            type="textarea"
            :rows="4"
            :placeholder="$t('auto.List.17.ac5bc6')"
          />
        </el-form-item>
        <el-form-item :label="$t('auto.List.18.2432b5')">
          <el-input v-model="form.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.List.51.625fb2') }}</el-button>
        <el-button type="primary" @click="submitForm">{{ $t('auto.List.52.38cf16') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog
      v-model="viewUsersDialogVisible"
      ::title="$t('auto.List.19.fb3844')"
      width="900px"
    >
      <el-table :data="segmentUsers" v-loading="viewUsersLoading" stripe max-height="500">
        <el-table-column prop="id" :label="$t('auto.List.20.30acd2')" width="80" />
        <el-table-column prop="name" :label="$t('auto.List.21.819767')" min-width="120" />
        <el-table-column prop="email" :label="$t('auto.List.22.3bc5e6')" min-width="180" />
        <el-table-column prop="phone" :label="$t('auto.List.23.8098e2')" width="140" />
        <el-table-column prop="rScore" :label="$t('auto.List.24.d96b97')" width="80" align="center" />
        <el-table-column prop="fScore" :label="$t('auto.List.25.72c524')" width="80" align="center" />
        <el-table-column prop="mScore" :label="$t('auto.List.26.4e30b7')" width="80" align="center" />
        <el-table-column prop="segmentLabel" :label="$t('auto.List.27.96636b')" width="120">
          <template #default="{ row: user }">
            <el-tag type="primary">{{ user.segmentLabel}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="lastOrderTime" :label="$t('auto.List.28.9bc1d7')" width="170" />
        <el-table-column :label="$t('auto.List.8.2b6bc0')" width="100" fixed="right">
          <template #default="{ row: user }">
            <el-button link type="primary" @click="viewUserDetail(user)">{{ $t('auto.List.53.f26225') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="viewUsersDialogVisible = false">{{ $t('auto.List.54.b15d91') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import {
  getUserSegments,
  createUserSegment,
  updateUserSegment,
  deleteUserSegment,
  getSegmentStats,
  getSegmentUsers
} from '@/api/userSegment.js'
import { toList } from '@/utils/list.js'
import {
  SEGMENT_TYPE,
  GROUP_STATUS,
  getStatusLabel as getStatusLabelFromArr,
  getStatusTagType as getStatusTagTypeFromArr,
} from '@/constants/status';
import { getEnabledLabel, getEnabledTagType } from '@/constants/enabled'

const router = useRouter()

const { t } = useI18n()

const loading = ref(false)
const searchKeyword = ref('')
const segments = ref([])
const stats = ref({ totalUsers: 0, highValue: 0, active: 0, churnRisk: 0 })
const matrixCounts = ref({})
const dialogVisible = ref(false)
const dialogTitle = ref(t('auto.List.55.c28915'))
const viewUsersDialogVisible = ref(false)
const viewUsersLoading = ref(false)
const currentSegment = ref(null)
const segmentUsers = ref([])
const formRef = ref()
const form = ref({
  id: 0,
  name: '',
  type: 'RFM',
  rDays: 30,
  fCount: 3,
  mAmount: 1000,
  criteria: '',
  remark: '',
  status: 'active'
})
const formRules = {
  name: [{ required: true, message: i18n.global.t(t('auto.List.9.4bb33c')), trigger: 'blur' }],
  type: [{ required: true, message: i18n.global.t(t('auto.List.10.099cb7')), trigger: 'change' }],
  criteria: [{ required: true, message: i18n.global.t(t('auto.List.56.e419bd')), trigger: 'blur' }]
}

const filteredSegments = computed(() => {
  if (!searchKeyword.value) return segments.value
  return segments.value.filter(s => s.name.includes(searchKeyword.value))
})

const LEGACY_SEGMENT_TYPE = { RFM: 'RFM', CUSTOM: t('auto.List.12.f1d4ff'), BEHAVIOR: t('auto.List.57.a04961') };
const getTypeLabel = (type) => LEGACY_SEGMENT_TYPE[type] || getStatusLabelFromArr(type, SEGMENT_TYPE)
const getTypeTagType = (type) => {
  const map = { RFM: 'success', CUSTOM: 'primary', BEHAVIOR: 'warning' }
  return map[type] || getStatusTagTypeFromArr(type, SEGMENT_TYPE) || 'info'
}

const getStatusLabel = (s) => getEnabledLabel(s);
const getStatusTagType = (s) => getEnabledTagType(s)

const refreshData = async () => {
  loading.value = true
  try {
    const [segRes, statsRes] = await Promise.all([
      getUserSegments(),
      getSegmentStats()
    ])
    segments.value = toList(segRes)
    const s = statsRes || {}
    stats.value = {
      totalUsers: s.totalUsers ?? s.total_users ?? s.total ?? 0,
      highValue: s.highValue ?? s.high_value ?? 0,
      active: s.active ?? s.active_count ?? 0,
      churnRisk: s.churnRisk ?? s.churn_risk ?? 0
    }
    matrixCounts.value = s.matrix || s.matrix_counts || s.rfmMatrix || {}
  } catch (error) {
    console.error(t('auto.List.58.65804d'), error)
  } finally {
    loading.value = false
  }
}

const onCellClick = (key) => {
  const map = {
    high_value: t('auto.List.36.66f84a'),
    important_develop: t('auto.List.37.adda3b'),
    important_keep: t('auto.List.38.1c57e2'),
    general_value: t('auto.List.40.6ecb76'),
    general_develop: t('auto.List.41.cef1bd'),
    general_keep: t('auto.List.42.1d517c'),
    potential_value: t('auto.List.44.74ba3e'),
    potential_develop: t('auto.List.45.2df631'),
    churn: t('auto.List.46.ddd6e7')
  }
  ElMessage.info(`查看分群：${map[key] || key}`)
  const fakeSegment = { id: key, name: map[key] };
  viewUsers(fakeSegment)
}

const showCreateDialog = () => {
  form.value = { id: 0, name: '', type: 'RFM', rDays: 30, fCount: 3, mAmount: 1000, criteria: '', remark: '', status: 'active' }
  dialogTitle.value = t('auto.List.55.c28915')
  dialogVisible.value = true
}

const editSegment = (row) => {
  form.value = { ...row }
  dialogTitle.value = t('auto.List.59.086612')
  dialogVisible.value = true
}

const submitForm = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      if (form.value.id) {
        await updateUserSegment(form.value.id, form.value)
        ElMessage.success(i18n.global.t(t('auto.List.60.55aa63')))
      } else {
        await createUserSegment(form.value)
        ElMessage.success(i18n.global.t(t('auto.List.61.04a691')))
      }
      dialogVisible.value = false
      refreshData()
    } catch (error) {
      ElMessage.error(i18n.global.t(t('auto.List.62.5fa802')))
    }
  })
}

const deleteSegment = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要删除分群 "${row.name}" 吗？`, t('auto.List.63.e83a25'), { type: 'warning' })
    await deleteUserSegment(row.id)
    ElMessage.success(i18n.global.t(t('auto.List.64.0007d1')))
    refreshData()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error(i18n.global.t(t('auto.List.65.acf066')))
  }
}

const viewUsers = async (row) => {
  currentSegment.value = { id: row.id, name: row.name }
  viewUsersDialogVisible.value = true
  viewUsersLoading.value = true
  try {
    const res = await getSegmentUsers(row.id)
    segmentUsers.value = toList(res)
  } catch {
    ElMessage.error(i18n.global.t(t('auto.List.66.7cd2c6')))
    segmentUsers.value = []
  } finally {
    viewUsersLoading.value = false
  }
}

const viewUserDetail = (user) => {
  router.push({ path: '/customer360', query: { id: user.id, name: user.name } });
}

onMounted(() => {
  refreshData()
})
</script>

<style scoped lang="scss">
.user-segment-page {
  padding: 20px;
}
.header-card {
  margin-bottom: 20px;
  :deep(.el-card__body) {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .header-content h2 {
    margin: 0 0 8px 0;
  }
  .subtitle {
    color: #909399;
    margin: 0;
  }
  .header-actions {
    display: flex;
    gap: 10px;
  }
}
.overview-row {
  margin-bottom: 20px;
  .stat-card {
    text-align: center;
    .stat-label {
      color: #909399;
      font-size: 14px;
      margin-bottom: 10px;
    }
    .stat-value {
      font-size: 32px;
      font-weight: bold;
    }
  }
}
.matrix-card {
  margin-bottom: 20px;
  .rfm-matrix {
    table {
      width: 100%;
      border-collapse: collapse;
      th, td {
        border: 1px solid #ebeef5;
        padding: 12px;
        text-align: center;
      }
      th {
        background: #f5f7fa;
      }
      .cell {
        &.high-value { background: #10B981; color: white; }
        &.good { background: #4F46E5; color: white; }
        &.warn { background: #F59E0B; color: white; }
        &.normal { background: #909399; color: white; }
        &.low-value { background: #EF4444; color: white; }
      }
    }
  }
}
.segments-card {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}
</style>
