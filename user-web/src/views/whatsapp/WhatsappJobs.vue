<template>
  <div class="page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.WhatsappJobs.20.e86df3') }}</span>
          <el-button type="primary" :icon="Plus" @click="showCreate = true">{{ $t('新建任务') }}</el-button>
        </div>
      </template>

      <el-form :inline="true" :model="search" class="search-form">
        <el-form-item :label="$t('auto.WhatsappJobs.1.3fea7c')">
          <el-select v-model="search.status" :placeholder="$t('auto.WhatsappJobs.2.a8b0c2')" clearable>
            <el-option :label="$t('auto.WhatsappJobs.3.6139a6')" value="pending" />
            <el-option :label="$t('auto.WhatsappJobs.4.46e386')" value="running" />
            <el-option :label="$t('auto.WhatsappJobs.5.fad522')" value="completed" />
            <el-option :label="$t('auto.WhatsappJobs.6.066584')" value="failed" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="load">{{ $t('auto.WhatsappJobs.21.e5f71f') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="jobs" style="width:100%" v-loading="loading" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="draft_title" :label="$t('auto.WhatsappJobs.7.22b433')" min-width="200" />
        <el-table-column prop="status" :label="$t('auto.WhatsappJobs.1.3fea7c')" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="total_count" :label="$t('auto.WhatsappJobs.8.9ed7d3')" width="80" align="center" />
        <el-table-column prop="sent_count" :label="$t('auto.WhatsappJobs.9.6f9369')" width="80" align="center" />
        <el-table-column prop="failed_count" :label="$t('auto.WhatsappJobs.10.acd5cb')" width="80" align="center" />
        <el-table-column prop="created_at" :label="$t('auto.WhatsappJobs.11.eca37c')" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.WhatsappJobs.12.2b6bc0')" width="150" fixed="right">
          <template #default="scope">
            <el-button size="small" type="primary" link @click="viewJob(scope.row)">{{ $t('auto.WhatsappJobs.22.f26225') }}</el-button>
            <el-button size="small" type="danger" link @click="confirmDelete(scope.row)">{{ $t('auto.WhatsappJobs.23.2f4aad') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="load"
          @current-change="load"
        />
      </div>
    </el-card>

    
    <el-dialog v-model="showCreate" :title="$t('auto.WhatsappJobs.13.f61371')" width="500px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item :label="$t('auto.WhatsappJobs.14.2b1806')">
          <el-select v-model="createForm.draft_id" :placeholder="$t('auto.WhatsappJobs.15.708c9d')" style="width: 100%">
            <el-option v-for="d in drafts" :key="d.id" :label="d.title" :value="d.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.WhatsappJobs.16.3e0aa8')">
          <el-input
            v-model="createForm.phone_numbers"
            type="textarea"
            :rows="4"
            :placeholder="$t('auto.WhatsappJobs.17.281f4a')"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreate = false">{{ $t('auto.WhatsappJobs.24.625fb2') }}</el-button>
        <el-button type="primary" :loading="creating" @click="submit">{{ $t('auto.WhatsappJobs.25.d9ac92') }}</el-button>
      </template>
    </el-dialog>

    
    <el-dialog v-model="showDetail" :title="$t('auto.WhatsappJobs.18.4a98bf')" width="500px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ detail.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.WhatsappJobs.7.22b433')">{{ detail.draft_title }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.WhatsappJobs.1.3fea7c')">
          <el-tag :type="getStatusType(detail.status)">{{ getStatusText(detail.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('auto.WhatsappJobs.8.9ed7d3')">{{ detail.total_count }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.WhatsappJobs.9.6f9369')">{{ detail.sent_count }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.WhatsappJobs.10.acd5cb')">{{ detail.failed_count }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.WhatsappJobs.11.eca37c')">{{ formatDate(detail.created_at) }}</el-descriptions-item>
        <el-descriptions-item :label="$t('auto.WhatsappJobs.19.7cbadb')">{{ formatDate(detail.finished_at) }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '@/api/whatsapp'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import { getOrderStatusLabel, getOrderStatusTagType } from '@/constants/orderStatus';

const { t } = useI18n()

const loading = ref(false)
const creating = ref(false)
const jobs = ref([])
const drafts = ref([])
const showCreate = ref(false)
const showDetail = ref(false)
const detail = ref({})

const search = reactive({ status: '' })
const pagination = reactive({ page: 1, page_size: 10, total: 0 })
const createForm = reactive({ draft_id: '', phone_numbers: '' })

const formatDate = (d) => {
  if (!d) return '-'
  const date = new Date(d)
  return isNaN(date.getTime()) ? '-' : date.toLocaleString('zh-CN', { hour12: false })
}

const getStatusType = (s) => {
  const types = { pending: 'info', running: 'warning', completed: 'success', failed: 'danger' }
  return types[s] || ''
}
const getStatusText = (s) => {
  const texts = { pending: t('auto.WhatsappJobs.3.6139a6'), running: t('auto.WhatsappJobs.4.46e386'), completed: t('auto.WhatsappJobs.5.fad522'), failed: t('auto.WhatsappJobs.6.066584') }
  return texts[s] || s
}

const loadDrafts = async () => {
  try {
    const res = await api.listDrafts()
    drafts.value = Array.isArray(res) ? res : (res?.list || [])
  } catch (e) {
    ElMessage.error(i18n.global.t(t('auto.WhatsappJobs.26.ef848e')))
    drafts.value = []
  }
}

const load = async () => {
  loading.value = true
  try {
    const res = await api.listJobs({
      page: pagination.page,
      page_size: pagination.page_size,
      status: search.status
    })
    if (Array.isArray(res)) {
      jobs.value = res
      pagination.total = res.length
    } else if (res && res.list) {
      jobs.value = res.list
      pagination.total = res.total || 0
    } else {
      jobs.value = []
      pagination.total = 0
    }
  } catch (e) {
    ElMessage.error(i18n.global.t(t('auto.WhatsappJobs.27.6f4db5')))
    jobs.value = []
  } finally {
    loading.value = false
  }
}

const submit = async () => {
  if (!createForm.draft_id) { ElMessage.warning(i18n.global.t(t('auto.WhatsappJobs.28.03b81b'))); return }
  if (!createForm.phone_numbers.trim()) { ElMessage.warning(i18n.global.t(t('auto.WhatsappJobs.29.8dd245'))); return }
  creating.value = true
  try {
    const phones = createForm.phone_numbers.split(/[\n,;]+/).map(s => s.trim()).filter(Boolean)
    await api.createJob({ draft_id: createForm.draft_id, phone_numbers: phones })
    ElMessage.success(i18n.global.t(t('auto.WhatsappJobs.30.fdbfcf')))
    showCreate.value = false
    createForm.draft_id = ''
    createForm.phone_numbers = ''
    load()
  } catch (err) {
    ElMessage.error(t('auto.WhatsappJobs.31.a2bddc') + (err.message || t('auto.WhatsappJobs.32.974e74')))
  } finally {
    creating.value = false
  }
}

const viewJob = async (row) => {
  try {
    const res = await api.getJob(row.id)
    detail.value = res || row
  } catch (e) {
    detail.value = row
  }
  showDetail.value = true
}

const confirmDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确认删除任务 #${row.id}?`, t('auto.WhatsappJobs.33.02d981'), { type: 'warning' })
    await api.deleteJob(row.id)
    ElMessage.success(i18n.global.t(t('auto.WhatsappJobs.34.0007d1')))
    load()
  } catch (err) {
    if (err !== 'cancel' && err?.message) {
      ElMessage.error(t('auto.WhatsappJobs.35.ad23f0') + err.message)
    }
  }
}

onMounted(() => {
  loadDrafts()
  load()
})
</script>

<style scoped>
.page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; font-weight: 600; }
.search-form { margin-bottom: 16px; }
.pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
</style>
