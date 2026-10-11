<template>
  <div class="page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.WhatsappDrafts.5fee9de7') }}</span>
          <el-button type="primary" :icon="Plus" @click="openCreate">{{ $t('新建草稿') }}</el-button>
        </div>
      </template>

      <el-form :inline="true" :model="search" class="search-form">
        <el-form-item :label="$t('auto.WhatsappDrafts.32c65d8d')">
          <el-input v-model="search.title" :placeholder="$t('auto.WhatsappDrafts.b1140e36')" clearable @keyup.enter="load" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="load">{{ $t('auto.WhatsappDrafts.e5f71fc3') }}</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="drafts" style="width:100%" v-loading="loading" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" :label="$t('auto.WhatsappDrafts.32c65d8d')" min-width="200" />
        <el-table-column prop="content" :label="$t('auto.WhatsappDrafts.2d711b09')" min-width="300" show-overflow-tooltip />
        <el-table-column prop="updated_at" :label="$t('auto.WhatsappDrafts.a001a226')" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.updated_at) }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.WhatsappDrafts.2b6bc0f2')" width="200" fixed="right">
          <template #default="scope">
            <el-button size="small" type="primary" link @click="openEdit(scope.row)">{{ $t('auto.WhatsappDrafts.95b351c8') }}</el-button>
            <el-button size="small" type="danger" link @click="confirmDelete(scope.row)">{{ $t('auto.WhatsappDrafts.2f4aaddd') }}</el-button>
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

    
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑草稿' : '新建草稿'" width="600px" @close="resetForm">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item :label="$t('auto.WhatsappDrafts.32c65d8d')" prop="title">
          <el-input v-model="form.title" :placeholder="$t('auto.WhatsappDrafts.87eb9514')" />
        </el-form-item>
        <el-form-item :label="$t('auto.WhatsappDrafts.2d711b09')" prop="content">
          <el-input type="textarea" v-model="form.content" :rows="6" :placeholder="$t('auto.WhatsappDrafts.90fbb862')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.WhatsappDrafts.625fb26b') }}</el-button>
        <el-button type="primary" :loading="saving" @click="submit">{{ $t('auto.WhatsappDrafts.be5fbbe3') }}</el-button>
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
  title: [{ required: true, message: i18n.global.t(t('auto.WhatsappDrafts.96641a78')), trigger: 'blur' }],
  content: [{ required: true, message: i18n.global.t(t('auto.WhatsappDrafts.a11cc7a6')), trigger: 'blur' }]
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
    ElMessage.error(i18n.global.t(t('auto.WhatsappDrafts.199b5384')))
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
      ElMessage.success(i18n.global.t(t('auto.WhatsappDrafts.55aa6366')))
    } else {
      await api.createDraft({ title: form.title, content: form.content })
      ElMessage.success(i18n.global.t(t('auto.WhatsappDrafts.04a691b3')))
    }
    dialogVisible.value = false
    load()
  } catch (err) {
    ElMessage.error((isEdit.value ? t('auto.WhatsappDrafts.32ac152b') : t('auto.WhatsappDrafts.d9ac9228')) + t('auto.WhatsappDrafts.8c470138') + (err.message || t('auto.WhatsappDrafts.974e7484')))
  } finally {
    saving.value = false
  }
}

const confirmDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确认删除草稿 "${row.title}"?`, t('auto.WhatsappDrafts.02d9819d'), { type: 'warning' })
    await api.deleteDraft(row.id)
    ElMessage.success(i18n.global.t(t('auto.WhatsappDrafts.0007d170')))
    load()
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error(t('auto.WhatsappDrafts.ad23f072') + (err?.message || t('auto.WhatsappDrafts.974e7484')))
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
