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
            <el-select v-model="filterScenario" ::placeholder="$t('auto.List.5d1d7351')" clearable style="width: 130px; margin-right: 10px">
              <el-option ::label="$t('auto.List.ff2c676a')" value="opening" />
              <el-option ::label="$t('auto.List.0b8352e6')" value="follow" />
              <el-option ::label="$t('auto.List.b65c587d')" value="objection" />
              <el-option ::label="$t('auto.List.655aed99')" value="closing" />
            </el-select>
            <el-input v-model="searchKeyword" ::placeholder="$t('auto.List.59ecb46e')" clearable style="width: 200px" />
          </div>
        </div>
      </template>
      <el-table :data="filteredScripts" v-loading="loading" stripe>
        <el-table-column prop="title" :label="$t('auto.List.32c65d8d')" min-width="180" />
        <el-table-column prop="scenario" :label="$t('auto.List.c931653c')" width="100">
          <template #default="{ row }">
            <el-tag>{{ getScenarioText(row.scenario) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="content" :label="$t('auto.List.2d711b09')" min-width="300" show-overflow-tooltip />
        <el-table-column prop="useCount" :label="$t('auto.List.194df98d')" width="100" />
        <el-table-column prop="effectiveness" :label="$t('auto.List.65a64798')" width="120">
          <template #default="{ row }">
            <el-rate v-model="row.effectiveness" disabled size="small" />
          </template>
        </el-table-column>
        <el-table-column prop="updatedAt" :label="$t('auto.List.a001a226')" width="180" />
        <el-table-column :label="$t('auto.List.2b6bc0f2')" width="250" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="useScript(row)">{{ $t('auto.List.ecff77a8') }}</el-button>
            <el-button link type="primary" @click="editScript(row)">{{ $t('auto.List.95b351c8') }}</el-button>
            <el-button link type="primary" @click="copyScript(row)">{{ $t('auto.List.79d3abe9') }}</el-button>
            <el-button link type="danger" @click="deleteScript(row)">{{ $t('auto.List.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="700px">
      <el-form :model="form" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('auto.List.04436328')" prop="title">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.7efcb0ce')" prop="scenario">
          <el-select v-model="form.scenario" style="width: 100%">
            <el-option :label="$t('auto.List.c0a428ee')" value="opening" />
            <el-option :label="$t('auto.List.994dcf28')" value="follow" />
            <el-option :label="$t('auto.List.054579de')" value="objection" />
            <el-option :label="$t('auto.List.87acf327')" value="closing" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.48dbefe8')" prop="category">
          <el-select v-model="form.category" :placeholder="$t('auto.List.9ef9e35f')" filterable allow-create default-first-option style="width: 100%">
            <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.name" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.4395eead')" prop="content">
          <el-input v-model="form.content" type="textarea" :rows="6" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.ddc7d28b')">
          <el-input v-model="form.variables" :placeholder="$t('auto.List.c0b8027e')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.14d34236')">
          <el-input v-model="form.tags" :placeholder="$t('auto.List.c379e80c')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.2432b575')">
          <el-input v-model="form.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.List.625fb26b') }}</el-button>
        <el-button type="primary" @click="submitForm">{{ $t('auto.List.38cf16f2') }}</el-button>
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
const dialogTitle = ref(t('auto.List.fd722d31'))
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
  title: [{ required: true, message: i18n.global.t(t('auto.List.96641a78')), trigger: 'blur' }],
  scenario: [{ required: true, message: i18n.global.t(t('auto.List.75f77abd')), trigger: 'change' }],
  category: [{ required: true, message: i18n.global.t(t('auto.List.9ef9e35f')), trigger: 'change' }],
  content: [{ required: true, message: i18n.global.t(t('auto.List.a11cc7a6')), trigger: 'blur' }]
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
  const map = { opening: t('auto.List.c0a428ee'), follow: t('auto.List.994dcf28'), objection: t('auto.List.054579de'), closing: t('auto.List.87acf327') }
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
  dialogTitle.value = t('auto.List.fd722d31')
  dialogVisible.value = true
}

const editScript = (row) => {
  form.value = { ...row }
  dialogTitle.value = t('auto.List.28400f19')
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
      ElMessage.success(i18n.global.t(t('auto.List.3b108349')))
      dialogVisible.value = false
      refreshData()
    } catch (error) {
      ElMessage.error(i18n.global.t(t('auto.List.6de920b4')))
    }
  })
}

const useScript = async (row) => {
  await useScriptApi(row.id)
  ElMessage.success(i18n.global.t(t('auto.List.63412144')))
  navigator.clipboard.writeText(row.content)
}

const copyScript = (row) => {
  navigator.clipboard.writeText(row.content)
  ElMessage.success(i18n.global.t(t('auto.List.52e6abbe')))
}

const deleteScript = async (row) => {
  try {
    await ElMessageBox.confirm(`确定删除话术 "${row.title}"？`, t('auto.List.e83a256e'), { type: 'warning' })
    await deleteScriptApi(row.id)
    ElMessage.success(i18n.global.t(t('auto.List.0007d170')))
    refreshData()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(i18n.global.t(t('auto.List.acf0664a')))
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
