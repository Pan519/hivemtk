<template>
  <div class="page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.WhatsappDrafts.9.5fee9d') }}</span>
          <el-button type="primary" :icon="Plus" @click="openCreate">{{ $t('新建草稿') }}</el-button>
        </div>
      </template>

      <el-form :inline="true" :model="search" class="search-form">
        <el-form-item :label="$t('auto.WhatsappDrafts.1.32c65d')">
          <el-input v-model="search.title" :placeholder="$t('auto.WhatsappDrafts.2.b1140e')" clearable @keyup.enter="load" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="load">{{ $t('auto.WhatsappDrafts.10.e5f71f') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="drafts" style="width:100%" v-loading="loading" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" :label="$t('auto.WhatsappDrafts.1.32c65d')" min-width="200" />
        <el-table-column prop="content" :label="$t('auto.WhatsappDrafts.3.2d711b')" min-width="300" show-overflow-tooltip />
        <el-table-column prop="updated_at" :label="$t('auto.WhatsappDrafts.4.a001a2')" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.updated_at) }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.WhatsappDrafts.5.2b6bc0')" width="200" fixed="right">
          <template #default="scope">
            <el-button size="small" type="primary" link @click="openEdit(scope.row)">{{ $t('auto.WhatsappDrafts.11.95b351') }}</el-button>
            <el-button size="small" type="danger" link @click="confirmDelete(scope.row)">{{ $t('auto.WhatsappDrafts.12.2f4aad') }}</el-button>
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

    
    <el-dialog v-model="dialogVisible" ::title="$t('auto.WhatsappDrafts.6.050da0')" width="600px" @close="resetForm">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item :label="$t('auto.WhatsappDrafts.1.32c65d')" prop="title">
          <el-input v-model="form.title" :placeholder="$t('auto.WhatsappDrafts.7.87eb95')" />
        </el-form-item>
        <el-form-item :label="$t('auto.WhatsappDrafts.3.2d711b')" prop="content">
          <el-input type="textarea" v-model="form.content" :rows="6" :placeholder="$t('auto.WhatsappDrafts.8.90fbb8')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.WhatsappDrafts.13.625fb2') }}</el-button>
        <el-button type="primary" :loading="saving" @click="submit">{{ $t('auto.WhatsappDrafts.14.be5fbb') }}</el-button>
      </template>
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

const { t } = useI18n()

const loading = ref(false)
const saving = ref(false)
const drafts = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()

const form = reactive({
  id: null,
  title: '',
  content: ''
})

const rules = {
  title: [{ required: true, message: i18n.global.t(t('auto.WhatsappDrafts.15.96641a')), trigger: 'blur' }],
  content: [{ required: true, message: i18n.global.t(t('auto.WhatsappDrafts.16.a11cc7')), trigger: 'blur' }]
}

const search = reactive({ title: '' })
const pagination = reactive({ page: 1, page_size: 10, total: 0 })

const formatDate = (d) => {
  if (!d) return '-'
  const date = new Date(d)
  return isNaN(date.getTime()) ? '-' : date.toLocaleString('zh-CN', { hour12: false })
}

const load = async () => {
  loading.value = true
  try {
    const res = await api.listDrafts({
      page: pagination.page,
      page_size: pagination.page_size,
      title: search.title
    });
    if (Array.isArray(res)) {
      drafts.value = res
      pagination.total = res.length
    } else if (res && res.list) {
      drafts.value = res.list
      pagination.total = res.total || 0
    } else {
      drafts.value = []
      pagination.total = 0
    }
  } catch (err) {
    ElMessage.error(i18n.global.t(t('auto.WhatsappDrafts.17.199b53')))
    drafts.value = []
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  isEdit.value = false
  Object.assign(form, { id: null, title: '', content: '' })
  dialogVisible.value = true
}

const openEdit = (row) => {
  isEdit.value = true
  Object.assign(form, { id: row.id, title: row.title, content: row.content })
  dialogVisible.value = true
}

const resetForm = () => {
  formRef.value?.resetFields()
  Object.assign(form, { id: null, title: '', content: '' })
}

const submit = async () => {
  try {
    await formRef.value.validate()
  } catch {
    return
  }
  saving.value = true
  try {
    if (isEdit.value && form.id) {
      await api.updateDraft(form.id, { title: form.title, content: form.content })
      ElMessage.success(i18n.global.t(t('auto.WhatsappDrafts.18.55aa63')))
    } else {
      await api.createDraft({ title: form.title, content: form.content })
      ElMessage.success(i18n.global.t(t('auto.WhatsappDrafts.19.04a691')))
    }
    dialogVisible.value = false
    load()
  } catch (err) {
    ElMessage.error((isEdit.value ? t('auto.WhatsappDrafts.20.32ac15') : t('auto.WhatsappDrafts.21.d9ac92')) + t('auto.WhatsappDrafts.22.8c4701') + (err.message || t('auto.WhatsappDrafts.23.974e74')))
  } finally {
    saving.value = false
  }
}

const confirmDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确认删除草稿 "${row.title}"?`, t('auto.WhatsappDrafts.24.02d981'), { type: 'warning' })
    await api.deleteDraft(row.id)
    ElMessage.success(i18n.global.t(t('auto.WhatsappDrafts.25.0007d1')))
    load()
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error(t('auto.WhatsappDrafts.26.ad23f0') + (err?.message || t('auto.WhatsappDrafts.23.974e74')))
    }
  }
}

onMounted(load)
</script>

<style scoped>
.page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; font-weight: 600; }
.search-form { margin-bottom: 16px; }
.pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
</style>
