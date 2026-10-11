<template>
  <div class="lead-group-selection">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('从线索库选择群体') }}</span>
        </div>
      </template>

      
      <el-form :model="filterForm" inline label-width="80px" style="margin-bottom: 20px;">
        <el-form-item :label="$t('auto.LeadGroupSelection.1.3fea7c')">
          <el-select v-model="filterForm.status" :placeholder="$t('auto.LeadGroupSelection.2.7090b1')" clearable>
            <el-option
              v-for="opt in leadStatusOptions"
              :key="opt.value"
              :label="opt.label"
              :value="opt.value"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item :label="$t('auto.LeadGroupSelection.3.26ca20')">
          <el-input v-model="filterForm.source" :placeholder="$t('auto.LeadGroupSelection.4.6c9445')" clearable />
        </el-form-item>
        
        <el-form-item :label="$t('auto.LeadGroupSelection.5.571d04')">
          <el-slider v-model="filterForm.minScore" :min="0" :max="100" show-input />
        </el-form-item>
        
        <el-form-item :label="$t('auto.LeadGroupSelection.6.e5f71f')">
          <el-input v-model="filterForm.search" :placeholder="$t('auto.LeadGroupSelection.7.ef2b01')" clearable />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="applyFilters">{{ $t('auto.LeadGroupSelection.13.c2fe62') }}</el-button>
          <el-button @click="resetFilters">{{ $t('auto.LeadGroupSelection.14.4b9c32') }}</el-button>
        </el-form-item>
      </el-form>

      
      <el-table 
        :data="leads" 
        style="width: 100%; margin-top: 20px;" 
        v-loading="loading"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="name" :label="$t('auto.LeadGroupSelection.8.60d045')" width="120" />
        <el-table-column prop="phone" :label="$t('auto.LeadGroupSelection.9.8098e2')" width="150" />
        <el-table-column prop="email" :label="$t('auto.LeadGroupSelection.10.3bc5e6')" width="200" />
        <el-table-column prop="company" :label="$t('auto.LeadGroupSelection.11.41716f')" width="150" />
        <el-table-column prop="source" :label="$t('auto.LeadGroupSelection.3.26ca20')" width="100" />
        <el-table-column prop="score" :label="$t('auto.LeadGroupSelection.12.b3dc71')" width="100">
          <template #default="{ row }">
            <el-tag :type="getScoreType(row.score)">
              {{ row.score }}分
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('auto.LeadGroupSelection.1.3fea7c')" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
      
      <el-pagination
        class="pagination"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="pagination.currentPage"
        :page-sizes="[10, 20, 50, 100]"
        :page-size="pagination.pageSize"
        layout="total, sizes, prev, pager, next, jumper"
        :total="pagination.total"
        v-if="pagination.total > 0"
      />
    </el-card>

    
    <el-card style="margin-top: 20px;" v-if="selectedLeads.length > 0">
      <template #header>
        <div class="card-header">
          <span>已选择 {{ selectedLeads.length }} 个联系人</span>
          <el-button type="primary" @click="proceedToMessaging">{{ $t('auto.LeadGroupSelection.15.507ec9') }}</el-button>
        </div>
      </template>
      
      <el-table :data="selectedLeads" style="width: 100%;">
        <el-table-column prop="name" :label="$t('auto.LeadGroupSelection.8.60d045')" width="120" />
        <el-table-column prop="phone" :label="$t('auto.LeadGroupSelection.9.8098e2')" width="150" />
        <el-table-column prop="email" :label="$t('auto.LeadGroupSelection.10.3bc5e6')" width="200" />
        <el-table-column prop="company" :label="$t('auto.LeadGroupSelection.11.41716f')" width="150" />
        <el-table-column prop="score" :label="$t('auto.LeadGroupSelection.12.b3dc71')" width="100">
          <template #default="{ row }">
            <el-tag :type="getScoreType(row.score)">
              {{ row.score }}分
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as clueApi from '@/api/clue'
import { useRouter } from 'vue-router'
import { LEAD_STATUS_OPTIONS_SIMPLE, getLeadStatusLabel, getLeadStatusTagType } from '@/constants/leadStatus'

const router = useRouter()
const { t } = useI18n()

const loading = ref(false)
const leads = ref([])
const selectedLeads = ref([])
const leadStatusOptions = LEAD_STATUS_OPTIONS_SIMPLE

const pagination = reactive({
  currentPage: 1,
  pageSize: 20,
  total: 0
})

const filterForm = reactive({
  status: '',
  source: '',
  minScore: 0,
  search: ''
})

const fetchLeads = async () => {
  try {
    loading.value = true
    const response = await clueApi.clueApi.list(
      pagination.currentPage,
      pagination.pageSize
    );
    const clueData = response?.list || response || [];
    leads.value = clueData.map(clue => ({
      id: clue.ID,
      name: clue.Name || clue.name || t('auto.LeadGroupSelection.16.1622dc'),
      phone: clue.Account || clue.account || '',
      email: '',
      company: clue.Address || clue.address || clue.City || clue.city || '',
      source: 'clue',
      score: 80,
      status: 'new'
    }));
    pagination.total = response?.total || response?.length || 0
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.LeadGroupSelection.17.aa2f7e')))
  } finally {
    loading.value = false
  }
}

const applyFilters = () => {
  pagination.currentPage = 1
  fetchLeads()
}

const resetFilters = () => {
  filterForm.status = ''
  filterForm.source = ''
  filterForm.minScore = 0
  filterForm.search = ''
  pagination.currentPage = 1
  fetchLeads()
}

const handleSelectionChange = (selection) => {
  selectedLeads.value = selection
}

const proceedToMessaging = () => {
  localStorage.setItem('selectedLeads', JSON.stringify(selectedLeads.value.map(lead => lead.id)));
  router.push('/whatsapp/group-messaging');
}

const getScoreType = (score) => {
  if (score >= 80) return 'success'
  if (score >= 60) return 'warning'
  return 'danger'
}

const getStatusType = (status) => {
  switch (status) {
    case 'new': return 'info'
    case 'contacted': return 'warning'
    case 'converted': return 'success'
    case 'lost': return 'danger'
    default: return 'info'
  }
}

const getStatusText = (status) => {
  switch (status) {
    case 'new': return t('auto.LeadGroupSelection.18.43eb02')
    case 'contacted': return t('auto.LeadGroupSelection.19.6b00b8')
    case 'converted': return t('auto.LeadGroupSelection.20.f9a204')
    case 'lost': return t('auto.LeadGroupSelection.21.f22a4d')
    default: return t('auto.LeadGroupSelection.16.1622dc')
  }
}

const handleSizeChange = (val) => {
  pagination.pageSize = val
  fetchLeads()
}

const handleCurrentChange = (val) => {
  pagination.currentPage = val
  fetchLeads()
}

onMounted(() => {
  fetchLeads()
});
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
</style>
