<template>
  <div class="script-template-page">
    <el-card class="header-card">
      <div>
        <h2>{{ $t('话术模板') }}</h2>
        <p class="subtitle">{{ $t('管理营销话术、销售话术模板') }}</p>
      </div>
      <el-button type="primary" @click="showCreateDialog">
        <el-icon><Plus /></el-icon>
        {{ $t('新增话术') }}
      </el-button>
    </el-card>

    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('话术列表') }}</span>
          <div>
            <el-select v-model="filterScenario" ::placeholder="$t('auto.List.1.5d1d73')" clearable style="width: 130px; margin-right: 10px">
              <el-option ::label="$t('auto.List.2.ff2c67')" value="opening" />
              <el-option ::label="$t('auto.List.3.0b8352')" value="follow" />
              <el-option ::label="$t('auto.List.4.b65c58')" value="objection" />
              <el-option ::label="$t('auto.List.5.655aed')" value="closing" />
            </el-select>
            <el-input v-model="searchKeyword" ::placeholder="$t('auto.List.6.59ecb4')" clearable style="width: 200px" />
          </div>
        </div>
      </template>
      <el-table :data="filteredScripts" v-loading="loading" stripe>
        <el-table-column prop="title" :label="$t('auto.List.7.32c65d')" min-width="180" />
        <el-table-column prop="scenario" :label="$t('auto.List.8.c93165')" width="100">
          <template #default="{ row }">
            <el-tag>{{ getScenarioText(row.scenario) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="content" :label="$t('auto.List.9.2d711b')" min-width="300" show-overflow-tooltip />
        <el-table-column prop="useCount" :label="$t('auto.List.10.194df9')" width="100" />
        <el-table-column prop="effectiveness" :label="$t('auto.List.11.65a647')" width="120">
          <template #default="{ row }">
            <el-rate v-model="row.effectiveness" disabled size="small" />
          </template>
        </el-table-column>
        <el-table-column prop="updatedAt" :label="$t('auto.List.12.a001a2')" width="180" />
        <el-table-column :label="$t('auto.List.13.2b6bc0')" width="250" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="useScript(row)">{{ $t('auto.List.28.ecff77') }}</el-button>
            <el-button link type="primary" @click="editScript(row)">{{ $t('auto.List.29.95b351') }}</el-button>
            <el-button link type="primary" @click="copyScript(row)">{{ $t('auto.List.30.79d3ab') }}</el-button>
            <el-button link type="danger" @click="deleteScript(row)">{{ $t('auto.List.31.2f4aad') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="700px">
      <el-form :model="form" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('auto.List.14.044363')" prop="title">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.15.7efcb0')" prop="scenario">
          <el-select v-model="form.scenario" style="width: 100%">
            <el-option :label="$t('auto.List.16.c0a428')" value="opening" />
            <el-option :label="$t('auto.List.17.994dcf')" value="follow" />
            <el-option :label="$t('auto.List.18.054579')" value="objection" />
            <el-option :label="$t('auto.List.19.87acf3')" value="closing" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.20.48dbef')" prop="category">
          <el-select v-model="form.category" :placeholder="$t('auto.List.21.9ef9e3')" filterable allow-create default-first-option style="width: 100%">
            <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.name" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.22.4395ee')" prop="content">
          <el-input v-model="form.content" type="textarea" :rows="6" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.23.ddc7d2')">
          <el-input v-model="form.variables" :placeholder="$t('auto.List.24.c0b802')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.25.14d342')">
          <el-input v-model="form.tags" :placeholder="$t('auto.List.26.c379e8')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.27.2432b5')">
          <el-input v-model="form.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.List.32.625fb2') }}</el-button>
        <el-button type="primary" @click="submitForm">{{ $t('auto.List.33.38cf16') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getScriptCategories, getScripts, createScript, updateScript, deleteScript as deleteScriptApi, useScript as useScriptApi } from '@/api/scriptTemplate.js'
import { toList } from '@/utils/list.js'

const { t } = useI18n()

const loading = ref(false)
const searchKeyword = ref('')
const filterScenario = ref('')
const scripts = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref(t('auto.List.34.fd722d'))
const formRef = ref()
const form = ref({
  id: 0,
  title: '',
  scenario: 'opening',
  category: '',
  content: '',
  variables: '',
  tags: '',
  remark: ''
})
const formRules = {
  title: [{ required: true, message: i18n.global.t(t('auto.List.35.96641a')), trigger: 'blur' }],
  scenario: [{ required: true, message: i18n.global.t(t('auto.List.36.75f77a')), trigger: 'change' }],
  category: [{ required: true, message: i18n.global.t(t('auto.List.21.9ef9e3')), trigger: 'change' }],
  content: [{ required: true, message: i18n.global.t(t('auto.List.37.a11cc7')), trigger: 'blur' }]
}

const categories = ref([]);
const loadCategories = async () => {
  try {
    const res = await getScriptCategories()
    categories.value = toList(res)
  } catch (e) {
    categories.value = []
  }
}

const filteredScripts = computed(() => {
  let result = scripts.value
  if (filterScenario.value) result = result.filter(s => s.scenario === filterScenario.value)
  if (searchKeyword.value) result = result.filter(s => s.title.includes(searchKeyword.value))
  return result
})

const getScenarioText = (s) => {
  const map = { opening: t('auto.List.16.c0a428'), follow: t('auto.List.17.994dcf'), objection: t('auto.List.18.054579'), closing: t('auto.List.19.87acf3') }
  return map[s] || s
}

const refreshData = async () => {
  loading.value = true
  try {
    const res = await getScripts()
    scripts.value = toList(res)
  } finally {
    loading.value = false
  }
}

const showCreateDialog = () => {
  form.value = { id: 0, title: '', scenario: 'opening', category: '', content: '', variables: '', tags: '', remark: '' }
  dialogTitle.value = t('auto.List.34.fd722d')
  dialogVisible.value = true
}

const editScript = (row) => {
  form.value = { ...row }
  dialogTitle.value = t('auto.List.38.28400f')
  dialogVisible.value = true
}

const submitForm = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      const payload = {
        title: form.value.title,
        category: form.value.category,
        scenario: form.value.scenario || '',
        content: form.value.content,
        variables: (form.value.variables || '')
          .split(',').map(s => String(s).trim()).filter(Boolean),
        tags: form.value.tags || '',
        is_public: false
      };
      if (form.value.id) {
        await updateScript(form.value.id, payload)
      } else {
        await createScript(payload)
      }
      ElMessage.success(i18n.global.t(t('auto.List.39.3b1083')))
      dialogVisible.value = false
      refreshData()
    } catch (error) {
      ElMessage.error(i18n.global.t(t('auto.List.40.6de920')))
    }
  })
}

const useScript = async (row) => {
  await useScriptApi(row.id)
  ElMessage.success(i18n.global.t(t('auto.List.41.634121')))
  navigator.clipboard.writeText(row.content)
}

const copyScript = (row) => {
  navigator.clipboard.writeText(row.content)
  ElMessage.success(i18n.global.t(t('auto.List.42.52e6ab')))
}

const deleteScript = async (row) => {
  try {
    await ElMessageBox.confirm(`确定删除话术 "${row.title}"？`, t('auto.List.43.e83a25'), { type: 'warning' })
    await deleteScriptApi(row.id)
    ElMessage.success(i18n.global.t(t('auto.List.44.0007d1')))
    refreshData()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(i18n.global.t(t('auto.List.45.acf066')))
  }
}

onMounted(() => { loadCategories(); refreshData() })
</script>

<style scoped lang="scss">
.script-template-page { padding: 20px; }
.header-card {
  margin-bottom: 20px;
  :deep(.el-card__body) { display: flex; justify-content: space-between; align-items: center; }
  h2 { margin: 0 0 8px 0; }
  .subtitle { color: #909399; margin: 0; }
}
.card-header { display: flex; justify-content: space-between; align-items: center; }
</style>
