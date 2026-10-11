<template>
  <div class="lead-group-selection">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('从线索库选择群体') }}</span>
        </div>
      </template>

      
      <el-form :model="filterForm" inline label-width="80px" style="margin-bottom: 20px;">
        <el-form-item :label="$t('auto.LeadGroupSelection.3fea7ca7')">
          <el-select v-model="filterForm.status" :placeholder="$t('auto.LeadGroupSelection.7090b1b5')" clearable>
            <el-option
              v-for="opt in leadStatusOptions"
              :key="opt.value"
              :label="opt.label"
              :value="opt.value"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item :label="$t('auto.LeadGroupSelection.26ca20b1')">
          <el-input v-model="filterForm.source" :placeholder="$t('auto.LeadGroupSelection.6c9445a4')" clearable />
        </el-form-item>
        
        <el-form-item :label="$t('auto.LeadGroupSelection.571d04f0')">
          <el-slider v-model="filterForm.minScore" :min="0" :max="100" show-input />
        </el-form-item>
        
        <el-form-item :label="$t('auto.LeadGroupSelection.e5f71fc3')">
          <el-input v-model="filterForm.search" :placeholder="$t('auto.LeadGroupSelection.ef2b0109')" clearable />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="applyFilters">{{ $t('auto.LeadGroupSelection.c2fe6253') }}</el-button>
          <el-button @click="resetFilters">{{ $t('auto.LeadGroupSelection.4b9c3271') }}</el-button>
        </el-form-item>
      </el-form>

      
      <el-table 
        :data="leads" 
        style="width: 100%; margin-top: 20px;" 
        v-loading="loading"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="name" :label="$t('auto.LeadGroupSelection.60d0458a')" width="120" />
        <el-table-column prop="phone" :label="$t('auto.LeadGroupSelection.8098e2b4')" width="150" />
        <el-table-column prop="email" :label="$t('auto.LeadGroupSelection.3bc5e602')" width="200" />
        <el-table-column prop="company" :label="$t('auto.LeadGroupSelection.41716f87')" width="150" />
        <el-table-column prop="source" :label="$t('auto.LeadGroupSelection.26ca20b1')" width="100" />
        <el-table-column prop="score" :label="$t('auto.LeadGroupSelection.b3dc7140')" width="100">
          <template #default="{ row }">
            <el-tag :type="getScoreType(row.score)">
              {{ row.score }}分
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('auto.LeadGroupSelection.3fea7ca7')" width="100">
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
          <el-button type="primary" @click="proceedToMessaging">{{ $t('auto.LeadGroupSelection.507ec9f3') }}</el-button>
        </div>
      </template>
      
      <el-table :data="selectedLeads" style="width: 100%;">
        <el-table-column prop="name" :label="$t('auto.LeadGroupSelection.60d0458a')" width="120" />
        <el-table-column prop="phone" :label="$t('auto.LeadGroupSelection.8098e2b4')" width="150" />
        <el-table-column prop="email" :label="$t('auto.LeadGroupSelection.3bc5e602')" width="200" />
        <el-table-column prop="company" :label="$t('auto.LeadGroupSelection.41716f87')" width="150" />
        <el-table-column prop="score" :label="$t('auto.LeadGroupSelection.b3dc7140')" width="100">
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
      name: clue.Name || clue.name || t('auto.LeadGroupSelection.1622dc9b'),
      phone: clue.Account || clue.account || '',
      email: '',
      company: clue.Address || clue.address || clue.City || clue.city || '',
      source: 'clue',
      score: 80,
      status: 'new'
    }));
    pagination.total = response?.total || response?.length || 0
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.LeadGroupSelection.aa2f7e51')))
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
    case 'new': return t('auto.LeadGroupSelection.43eb027a')
    case 'contacted': return t('auto.LeadGroupSelection.6b00b8a1')
    case 'converted': return t('auto.LeadGroupSelection.f9a204fc')
    case 'lost': return t('auto.LeadGroupSelection.f22a4db3')
    default: return t('auto.LeadGroupSelection.1622dc9b')
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
