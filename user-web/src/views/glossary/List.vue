<template>
  <div class="glossary-list-page">
    
    <el-card class="header-card" shadow="never">
      <div class="header-content">
        <div>
          <h2>{{ $t('术语表管理') }}</h2>
          <p class="subtitle">{{ $t('auto.List.40.660669') }}</p>
        </div>
        <div class="header-actions">
          <el-button @click="loadList" :loading="loading">
            <el-icon><Refresh /></el-icon>
            {{ $t('刷新') }}
          </el-button>
          <el-button type="primary" @click="openCreateDialog">
            <el-icon><Plus /></el-icon>
            {{ $t('新增术语') }}
          </el-button>
          <el-button type="success" plain @click="openValidateDialog">
            <el-icon><Search /></el-icon>
            {{ $t('校验预览') }}
          </el-button>
        </div>
      </div>
    </el-card>

    
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filter" @submit.prevent>
        <el-form-item ::label="$t('auto.List.1.bcac61')">
          <el-select v-model="filter.category" ::placeholder="$t('auto.List.2.0d0d2d')" clearable style="width: 160px" @change="onSearch">
            <el-option v-for="opt in categoryOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item ::label="$t('auto.List.3.1b8578')">
          <el-select v-model="filter.status" ::placeholder="$t('auto.List.4.925332')" clearable style="width: 140px" @change="onSearch">
            <el-option ::label="$t('auto.List.5.0e4c2a')" :value="1" />
            <el-option ::label="$t('auto.List.6.27ea4f')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item ::label="$t('auto.List.7.648229')">
          <el-input
            v-model="filter.keyword"
            ::placeholder="$t('auto.List.8.aabc58')"
            clearable
            style="width: 240px"
            @keyup.enter="onSearch"
            @clear="onSearch"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSearch">
            <el-icon><Search /></el-icon>
            {{ $t('搜索') }}
          </el-button>
          <el-button @click="resetFilter">{{ $t('重置') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    
    <el-card shadow="never">
      <el-table :data="list" v-loading="loading" stripe border>
        <el-table-column prop="term_id" label="term_id" min-width="160" show-overflow-tooltip />
        <el-table-column :label="$t('auto.List.9.d0771a')" width="120">
          <template #default="{ row }">
            <el-tag :type="getCategoryTagType(row.category)" size="small">
              {{ getCategoryLabel(row.category) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.10.412e8c')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.preserve ? 'warning' : 'info'" size="small">
              {{ row.preserve ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.11.36577c')" min-width="280">
          <template #default="{ row }">
            <div class="trans-cell">
              <el-tag
                v-for="(text, lang) in normalizeTranslations(row.translations)"
                :key="lang"
                size="small"
                type="success"
                class="trans-tag"
              >
                <span class="trans-lang">{{ getLanguageLabel(lang) }}</span>
                <span class="trans-text">{{ truncateText(text, 20) }}</span>
              </el-tag>
              <span v-if="!hasTranslations(row.translations)" class="empty-text">-</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="pattern" :label="$t('auto.List.12.ec0fcb')" min-width="160" show-overflow-tooltip>
          <template #default="{ row }">{{ row.pattern || '-' }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.13.3fea7c')" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('auto.List.14.eca37c')" width="160">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.15.2b6bc0')" width="180" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="openEditDialog(row)">{{ $t('auto.List.41.95b351') }}</el-button>
            <el-button
              link
              :type="row.status === 1 ? 'warning' : 'success'"
              size="small"
              @click="onToggleStatus(row)"
            >
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
            <el-button link type="danger" size="small" @click="onDelete(row)">{{ $t('auto.List.42.2f4aad') }}</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty :description="$t('auto.List.16.1081a8')" />
        </template>
      </el-table>

      
      <div class="pagination-wrap" v-if="total > 0">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          background
          @size-change="loadList"
          @current-change="loadList"
        />
      </div>
    </el-card>

    
    <el-dialog
      v-model="formDialogVisible"
      ::title="$t('auto.List.17.bd0279')"
      width="720px"
      top="6vh"
      :close-on-click-modal="false"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="110px" v-loading="formLoading">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="term_id" prop="term_id">
              <el-input v-model="form.term_id" :placeholder="$t('auto.List.18.65a913')" :disabled="isEdit" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item ::label="$t('auto.List.1.bcac61')" prop="category">
              <el-select v-model="form.category" style="width: 100%" ::placeholder="$t('auto.List.19.a19bfd')">
                <el-option v-for="opt in categoryOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item ::label="$t('auto.List.20.703c69')">
              <el-switch v-model="form.preserve" />
              <div class="form-tip">{{ $t('auto.List.43.e6b978') }}</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item ::label="$t('auto.List.3.1b8578')">
              <el-select v-model="form.status" style="width: 100%">
                <el-option ::label="$t('auto.List.5.0e4c2a')" :value="1" />
                <el-option ::label="$t('auto.List.6.27ea4f')" :value="0" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item ::label="$t('auto.List.21.19c7ab')">
              <el-input v-model="form.pattern" :placeholder="$t('auto.List.22.a03752')" />
              <div class="form-tip">{{ $t('auto.List.44.02f561') }}</div>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item ::label="$t('auto.List.23.d25df6')">
              <div class="trans-editor">
                <div v-for="(item, idx) in translationList" :key="idx" class="trans-row">
                  <el-select
                    v-model="item.lang"
                    ::placeholder="$t('auto.List.24.eb5d3f')"
                    filterable
                    style="width: 180px"
                  >
                    <el-option
                      v-for="lang in languageOptions"
                      :key="lang.value"
                      :label="lang.label"
                      :value="lang.value"
                    />
                  </el-select>
                  <el-input v-model="item.text" ::placeholder="$t('auto.List.25.885c32')" style="flex: 1" />
                  <el-button link type="danger" size="small" @click="removeTranslation(idx)">
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </div>
                <el-button size="small" type="primary" plain @click="addTranslation">
                  <el-icon><Plus /></el-icon>
                  {{ $t('添加译法') }}
                </el-button>
              </div>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="formDialogVisible = false">{{ $t('取消') }}</el-button>
        <el-button type="primary" :loading="submitting" @click="submitForm">
          {{ $t('保存') }}
        </el-button>
      </template>
    </el-dialog>

    
    <el-dialog
      v-model="validateDialogVisible"
      ::title="$t('auto.List.26.fde388')"
      width="760px"
      top="6vh"
      :close-on-click-modal="false"
    >
      <el-form :model="validateForm" label-width="90px">
        <el-form-item ::label="$t('auto.List.27.07ae79')" required>
          <el-select v-model="validateForm.lang" ::placeholder="$t('auto.List.28.d8761b')" style="width: 220px">
            <el-option
              v-for="lang in languageOptions"
              :key="lang.value"
              :label="lang.label"
              :value="lang.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item ::label="$t('auto.List.29.288688')" required>
          <el-input
            v-model="validateForm.text"
            type="textarea"
            :rows="5"
            :placeholder="$t('auto.List.30.8604fa')"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="validating" @click="runValidate">
            <el-icon><Search /></el-icon>
            {{ $t('执行校验') }}
          </el-button>
          <el-button @click="clearValidateResult">{{ $t('清空') }}</el-button>
        </el-form-item>
      </el-form>

      <template v-if="validateResult">
        <el-divider content-position="left">{{ $t('校验结果') }}</el-divider>
        <el-descriptions :column="2" border size="small" class="validate-summary">
          <el-descriptions-item ::label="$t('auto.List.31.568a78')">
            {{ validateResult.hit_count ?? (validateResult.hits?.length || 0) }}
          </el-descriptions-item>
          <el-descriptions-item ::label="$t('auto.List.32.b5ca6b')">
            <el-tag :type="validateViolations.length ? 'danger' : 'success'" size="small">
              {{ validateViolations.length }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>

        <div class="result-section">
          <div class="result-title">{{ $t('命中的术语') }}</div>
          <el-table :data="validateHits" stripe size="small" border>
            <el-table-column prop="term_id" label="term_id" min-width="140" show-overflow-tooltip />
            <el-table-column prop="matched_text" :label="$t('auto.List.33.f8605e')" min-width="140" show-overflow-tooltip>
              <template #default="{ row }">{{ row.matched_text || row.match || '-' }}</template>
            </el-table-column>
            <el-table-column prop="translation" :label="$t('auto.List.34.167e59')" min-width="160" show-overflow-tooltip>
              <template #default="{ row }">{{ row.translation || row.translated || '-' }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.35.24fb86')" width="80" align="center">
              <template #default="{ row }">
                <el-tag v-if="row.preserve" type="warning" size="small">{{ $t('auto.List.45.0a60ac') }}</el-tag>
                <span v-else>-</span>
              </template>
            </el-table-column>
            <template #empty>
              <el-empty :description="$t('auto.List.36.25906e')" :image-size="50" />
            </template>
          </el-table>
        </div>

        <div class="result-section" v-if="validateViolations.length">
          <div class="result-title violation-title">{{ $t('违规项') }}</div>
          <el-table :data="validateViolations" stripe size="small" border>
            <el-table-column prop="term_id" label="term_id" min-width="140" show-overflow-tooltip />
            <el-table-column prop="issue" :label="$t('auto.List.37.254dd6')" min-width="200" show-overflow-tooltip>
              <template #default="{ row }">
                <span class="violation-text">{{ row.issue || row.reason || row.message || '-' }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="snippet" :label="$t('auto.List.38.9010ca')" min-width="200" show-overflow-tooltip>
              <template #default="{ row }">{{ row.snippet || row.context || '-' }}</template>
            </el-table-column>
          </el-table>
        </div>
      </template>
      <el-empty v-else-if="!validating" ::description="$t('auto.List.39.5bc9b1')" :image-size="80" />
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Plus, Search, Delete } from '@element-plus/icons-vue'
import {
  listGlossaries,
  createGlossary,
  updateGlossary,
  deleteGlossary,
  validateGlossary
} from '@/api/glossary.js'
import {
  SUPPORTED_LANGUAGES,
  getLanguageLabel
} from '@/constants/languages'

const categoryOptions = [
  { value: 'brand', label: t('auto.List.46.94c1f8') },
  { value: 'product', label: t('auto.List.47.882337') },
  { value: 'entity', label: t('auto.List.48.d041a5') },
  { value: 'technical', label: t('auto.List.49.85f2d0') },
  { value: 'custom', label: t('auto.List.50.f1d4ff') }
];
const categoryMap = categoryOptions.reduce((acc, o) => { acc[o.value] = o.label; return acc }, {})
const getCategoryLabel = (v) => categoryMap[v] || v || '-'
const getCategoryTagType = (v) => {
  if (v === 'brand') return 'danger'
  if (v === 'product') return 'warning'
  if (v === 'entity') return 'primary'
  if (v === 'technical') return 'success'
  return 'info'
}

const languageOptions = SUPPORTED_LANGUAGES;

const { t } = useI18n()

const loading = ref(false);
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const filter = ref({ category: '', status: '', keyword: '' })

const formDialogVisible = ref(false);
const formLoading = ref(false)
const submitting = ref(false)
const isEdit = ref(false)
const formRef = ref(null)
const editingTermId = ref('')
const form = ref({
  term_id: '',
  category: 'brand',
  preserve: false,
  pattern: '',
  status: 1,
  translations: {}
})
const translationList = ref([]);
const rules = {
  term_id: [{ required: true, message: t('auto.List.51.835e38'), trigger: 'blur' }],
  category: [{ required: true, message: t('auto.List.52.8bb820'), trigger: 'change' }]
}

const validateDialogVisible = ref(false);
const validating = ref(false)
const validateForm = ref({ lang: 'en', text: '' })
const validateResult = ref(null)

const validateHits = computed(() => {
  const r = validateResult.value
  if (!r) return []
  return Array.isArray(r.hits) ? r.hits : (Array.isArray(r.items) ? r.items : [])
})

const validateViolations = computed(() => {
  const r = validateResult.value
  if (!r) return []
  return Array.isArray(r.violations) ? r.violations : (Array.isArray(r.issues) ? r.issues : [])
})

const normalizeTranslations = (t) => {
  if (!t) return {}
  if (Array.isArray(t)) {
    const obj = {};
    t.forEach((item) => {
      if (item && item.lang) obj[item.lang] = item.text || item.translation || ''
    })
    return obj
  }
  return t
};

const hasTranslations = (t) => {
  const obj = normalizeTranslations(t)
  return Object.keys(obj).length > 0
}

const truncateText = (text, len) => {
  if (text === undefined || text === null || text === '') return '-'
  const s = String(text)
  return s.length > len ? s.slice(0, len) + '...' : s
}

const formatTime = (t) => {
  if (!t) return '-'
  const d = new Date(t)
  if (isNaN(d.getTime())) return '-'
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const loadList = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    if (filter.value.category) params.category = filter.value.category
    if (filter.value.status !== '' && filter.value.status !== null && filter.value.status !== undefined) {
      params.status = filter.value.status
    }
    if (filter.value.keyword) params.keyword = filter.value.keyword
    const res = await listGlossaries(params)
    if (Array.isArray(res)) {
      list.value = res
      total.value = res.length
    } else {
      list.value = res?.list || res?.items || []
      total.value = res?.total ?? list.value.length
    }
  } catch (e) {
    ElMessage.error(i18n.global.t(t('auto.List.53.150847')) + '：' + (e.message || t('auto.List.54.974e74')))
  } finally {
    loading.value = false
  }
};

const onSearch = () => {
  page.value = 1
  loadList()
}

const resetFilter = () => {
  filter.value = { category: '', status: '', keyword: '' }
  page.value = 1
  loadList()
}

const resetForm = () => {
  form.value = {
    term_id: '',
    category: 'brand',
    preserve: false,
    pattern: '',
    status: 1,
    translations: {}
  }
  translationList.value = []
  editingTermId.value = ''
};

const addTranslation = () => {
  translationList.value.push({ lang: '', text: '' })
}

const removeTranslation = (idx) => {
  translationList.value.splice(idx, 1)
}

const openCreateDialog = () => {
  isEdit.value = false
  resetForm()
  formDialogVisible.value = true
}

const openEditDialog = async (row) => {
  isEdit.value = true
  editingTermId.value = row.term_id
  resetForm()
  formDialogVisible.value = true
  formLoading.value = true
  try {
    const detail = row.term_id && !hasTranslations(row.translations)
      ? await getGlossarySafe(row.term_id)
      : row;
    form.value = {
      term_id: detail.term_id || row.term_id,
      category: detail.category || 'brand',
      preserve: !!detail.preserve,
      pattern: detail.pattern || '',
      status: detail.status === 0 ? 0 : 1,
      translations: {}
    }
    const obj = normalizeTranslations(detail.translations)
    translationList.value = Object.keys(obj).map((lang) => ({ lang, text: obj[lang] || '' }))
  } catch (e) {
    ElMessage.error(i18n.global.t(t('auto.List.55.298eeb')) + '：' + (e.message || t('auto.List.54.974e74')))
  } finally {
    formLoading.value = false
  }
}

const getGlossarySafe = async (termId) => {
  try {
    const { getGlossary } = await import('@/api/glossary.js')
    return await getGlossary(termId)
  } catch (e) {
    return {}
  }
}

const submitForm = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
  } catch (_) {
    return
  }
  const translations = {};
  let hasInvalid = false
  translationList.value.forEach((item) => {
    if (item.lang && item.text) {
      translations[item.lang] = item.text
    } else if (item.lang || item.text) {
      hasInvalid = true
    }
  })
  if (hasInvalid) {
    ElMessage.warning(i18n.global.t(t('auto.List.56.cd2761')))
  }

  const payload = {
    term_id: form.value.term_id,
    category: form.value.category,
    preserve: form.value.preserve,
    pattern: form.value.pattern || '',
    status: form.value.status,
    translations
  }

  submitting.value = true
  try {
    if (isEdit.value) {
      await updateGlossary(editingTermId.value, payload)
      ElMessage.success(i18n.global.t(t('auto.List.57.55aa63')))
    } else {
      await createGlossary(payload)
      ElMessage.success(i18n.global.t(t('auto.List.58.04a691')))
    }
    formDialogVisible.value = false
    loadList()
  } catch (e) {
    ElMessage.error((isEdit.value ? i18n.global.t(t('auto.List.59.930442')) : i18n.global.t(t('auto.List.60.a88928'))) + '：' + (e.message || t('auto.List.54.974e74')))
  } finally {
    submitting.value = false
  }
}

const onToggleStatus = (row) => {
  const next = row.status === 1 ? 0 : 1
  const action = next === 1 ? t('auto.List.61.7854b5') : t('auto.List.62.710ad0')
  ElMessageBox.confirm(
    `确认${action}术语「${row.term_id}」吗？`,
    i18n.global.t(t('auto.List.63.ce828a')),
    { type: 'warning', confirmButtonText: i18n.global.t(t('auto.List.64.e83a25')) + action, cancelButtonText: i18n.global.t(t('auto.List.65.625fb2')) }
  ).then(async () => {
    try {
      await updateGlossary(row.term_id, { status: next })
      ElMessage.success(`${action}成功`)
      loadList()
    } catch (e) {
      ElMessage.error(`${action}失败：` + (e.message || t('auto.List.54.974e74')))
    }
  }).catch(() => {})
};

const onDelete = (row) => {
  ElMessageBox.confirm(
    `确认删除术语「${row.term_id}」吗？删除后不可恢复。`,
    i18n.global.t(t('auto.List.66.50eaf9')),
    { type: 'warning', confirmButtonText: i18n.global.t(t('auto.List.67.631cd2')), cancelButtonText: i18n.global.t(t('auto.List.65.625fb2')) }
  ).then(async () => {
    try {
      await deleteGlossary(row.term_id)
      ElMessage.success(i18n.global.t(t('auto.List.68.0007d1')))
      loadList()
    } catch (e) {
      ElMessage.error(i18n.global.t(t('auto.List.69.acf066')) + '：' + (e.message || t('auto.List.54.974e74')))
    }
  }).catch(() => {})
};

const openValidateDialog = () => {
  validateForm.value = { lang: 'en', text: '' }
  validateResult.value = null
  validateDialogVisible.value = true
};

const runValidate = async () => {
  if (!validateForm.value.text || !validateForm.value.text.trim()) {
    ElMessage.warning(i18n.global.t(t('auto.List.70.412ecf')))
    return
  }
  if (!validateForm.value.lang) {
    ElMessage.warning(i18n.global.t(t('auto.List.71.c91ca0')))
    return
  }
  validating.value = true
  validateResult.value = null
  try {
    const res = await validateGlossary({
      text: validateForm.value.text,
      lang: validateForm.value.lang
    })
    validateResult.value = res || { hits: [], violations: [] }
    ElMessage.success(i18n.global.t(t('auto.List.72.0607a3')))
  } catch (e) {
    ElMessage.error(i18n.global.t(t('auto.List.73.136b33')) + '：' + (e.message || t('auto.List.54.974e74')))
  } finally {
    validating.value = false
  }
}

const clearValidateResult = () => {
  validateResult.value = null
  validateForm.value.text = ''
}

onMounted(() => {
  loadList()
});
</script>

<style scoped lang="scss">
.glossary-list-page { padding: 20px; }

.header-card {
  margin-bottom: 16px;
  :deep(.el-card__body) { padding: 16px 20px; }
  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    h2 { margin: 0 0 6px 0; font-size: 20px; }
    .subtitle { color: #909399; margin: 0; font-size: 13px; }
    .header-actions { display: flex; gap: 8px; }
  }
}

.filter-card { margin-bottom: 16px; }

.trans-cell {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  .trans-tag {
    display: inline-flex;
    align-items: center;
    max-width: 100%;
    .trans-lang { font-weight: 600; margin-right: 4px; }
    .trans-text { color: #303133; }
  }
  .empty-text { color: #909399; }
}

.pagination-wrap {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  line-height: 1.4;
  margin-top: 2px;
}

.trans-editor {
  width: 100%;
  .trans-row {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 8px;
  }
}

.validate-summary { margin-bottom: 16px; }

.result-section { margin-top: 16px; }
.result-title {
  font-weight: 600;
  margin-bottom: 8px;
  color: #303133;
}
.violation-title { color: #F56C6C; }
.violation-text { color: #F56C6C; }
</style>
