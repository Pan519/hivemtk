<template>
  <div class="faq-list-page">
    
    <el-card class="header-card" shadow="never">
      <div class="header-content">
        <div>
          <h2>{{ $t('auto.List.e571cd55') }}</h2>
          <p class="subtitle">{{ $t('auto.List.8360227b') }}</p>
        </div>
        <div class="header-actions">
          <el-button @click="loadList" :loading="loading">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
          <el-button @click="goTest">
            <el-icon><Aim /></el-icon>
            匹配测试
          </el-button>
          <el-button type="primary" @click="goCreate">
            <el-icon><Plus /></el-icon>
            新建 FAQ
          </el-button>
        </div>
      </div>
    </el-card>

    
    <el-row :gutter="20" class="stat-row" v-if="stats">
      <el-col :span="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-label">{{ $t('auto.List.9ed7d3ad') }}</div>
          <div class="stat-value">{{ stats.total || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-label">{{ $t('auto.List.7854b52a') }}</div>
          <div class="stat-value text-success">{{ stats.enabled || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-label">{{ $t('auto.List.ba5a6007') }}</div>
          <div class="stat-value text-primary">{{ enableRate }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-label">{{ $t('auto.List.cc59b3db') }}</div>
          <div class="stat-value text-info">{{ categoryCount }}</div>
        </el-card>
      </el-col>
    </el-row>

    
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filter" @submit.prevent>
        <el-form-item :label="$t('auto.List.9699a50e')">
          <el-input
            v-model="filter.keyword"
            :placeholder="$t('auto.List.6f505a66')"
            clearable
            style="width: 240px"
            @keyup.enter="onSearch"
            @clear="onSearch"
          />
        </el-form-item>
        <el-form-item :label="$t('auto.List.d0771a42')">
          <el-select v-model="filter.category" :placeholder="$t('auto.List.1a750305')" clearable style="width: 160px" @change="onSearch">
            <el-option v-for="c in categoryOptions" :key="c" :label="c" :value="c" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.764aca4e')">
          <el-input
            v-model="filter.intent"
            :placeholder="$t('auto.List.3814c5cd')"
            clearable
            style="width: 160px"
            @keyup.enter="onSearch"
            @clear="onSearch"
          />
        </el-form-item>
        <el-form-item :label="$t('auto.List.3fea7ca7')">
          <el-select v-model="filter.enabled" :placeholder="$t('auto.List.443483c9')" clearable style="width: 130px" @change="onSearch">
            <el-option :label="$t('auto.List.7854b52a')" :value="true" />
            <el-option :label="$t('auto.List.710ad08b')" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="resetFilter">{{ $t('auto.List.4b9c3271') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    
    <el-card shadow="never">
      <el-table :data="list" v-loading="loading" stripe border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column :label="$t('auto.List.5dc99f6e')" min-width="280" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="question-text">{{ row.question }}</div>
            <div v-if="row.keywords && row.keywords.length" class="keyword-list">
              <el-tag v-for="k in row.keywords.slice(0, 5)" :key="k" size="small" type="info" effect="plain" class="keyword-tag">
                {{ k }}
              </el-tag>
              <span v-if="row.keywords.length > 5" class="keyword-more">+{{ row.keywords.length - 5 }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.1a1dcce5')" min-width="220" show-overflow-tooltip>
          <template #default="{ row }">{{ truncate(row.answer, 80) }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.de92100c')" width="180" show-overflow-tooltip>
          <template #default="{ row }">
            <el-tag
              v-if="getKBName(row.kb_id)"
              size="small"
              type="success"
              effect="plain"
            >
              {{ getKBName(row.kb_id) }}
            </el-tag>
            <span v-else class="text-muted">{{ $t('auto.List.f64a3850') }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.d0771a42')" width="120" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.category" size="small" type="info">{{ row.category }}</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.764aca4e')" width="120" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.intent" size="small" type="warning">{{ row.intent }}</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.8f3599b6')" width="100" align="center">
          <template #default="{ row }">
            <el-progress
              :percentage="Math.round((row.confidence || 0) * 100)"
              :stroke-width="10"
              :show-text="true"
              :color="getConfColor(row.confidence)"
              style="width: 80px"
            />
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.e5852c46')" width="100" align="center">
          <template #default="{ row }">
            <el-tag size="small" :type="row.hit_count > 0 ? 'success' : 'info'">
              {{ row.hit_count || 0 }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.3fea7ca7')" width="90" align="center">
          <template #default="{ row }">
            <el-switch
              :model-value="row.enabled !== false"
              :loading="row._toggling"
              @change="(v) => onToggle(row, v)"
            />
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.a001a226')" width="160">
          <template #default="{ row }">{{ formatTime(row.updated_at) }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.2b6bc0f2')" width="180" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="goEdit(row)">{{ $t('auto.List.95b351c8') }}</el-button>
            <el-button link type="primary" size="small" @click="onTestOne(row)">{{ $t('auto.List.db06c78d') }}</el-button>
            <el-button link type="danger" size="small" @click="onDelete(row)">{{ $t('auto.List.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty :description="$t('auto.List.804d71b0')" />
        </template>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="filter.page"
          v-model:page-size="filter.page_size"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadList"
          @current-change="loadList"
        />
      </div>
    </el-card>

    
    <el-dialog v-model="testDialogVisible" :title="$t('auto.List.8dc65c8b')" width="640px" top="6vh">
      <el-form :model="testForm" label-width="80px">
        <el-form-item :label="$t('auto.List.dd60233d')">
          <el-input
            v-model="testForm.msg"
            type="textarea"
            :rows="2"
            :placeholder="$t('auto.List.4fb34744')"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="testLoading" @click="runTest">
            <el-icon><Aim /></el-icon>
            执行匹配
          </el-button>
        </el-form-item>
      </el-form>
      <template v-if="testResult.length">
        <el-divider content-position="left">{{ $t('auto.List.8cfdfd29') }}</el-divider>
        <el-table :data="testResult" stripe size="small" border>
          <el-table-column type="index" :label="$t('auto.List.a4dc00d8')" width="70" align="center" />
          <el-table-column :label="$t('auto.List.5dc99f6e')" min-width="200" show-overflow-tooltip>
            <template #default="{ row }">{{ row.entry?.question }}</template>
          </el-table-column>
          <el-table-column :label="$t('auto.List.8f3599b6')" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="getMatchTagType(row.score)" size="small">{{ (row.score * 100).toFixed(0) }}%</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('auto.List.194ffd85')" width="120" align="center">
            <template #default="{ row }">
              <el-tag size="small" type="info">{{ row.match_type || '-' }}</el-tag>
            </template>
          </el-table-column>
        </el-table>
      </template>
      <el-empty v-else-if="!testLoading && testTouched" :description="$t('auto.List.f10c9104')" :image-size="60" />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Plus, Search, Aim } from '@element-plus/icons-vue'
import { faqApi } from '@/api/faq'
import { listKBs } from '@/api/knowledgeBase'

const router = useRouter()

const { t } = useI18n()

const loading = ref(false);
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const stats = ref(null)
const filter = ref({
  keyword: '',
  category: '',
  intent: '',
  enabled: ''
})

const kbMap = ref({});

const loadKBMap = async () => {
  try {
    const res = await listKBs({ kb_type: 'faq', page: 1, page_size: 200 }).catch(() => null)
    const items = Array.isArray(res) ? res : res?.list || res?.items || []
    const map = {}
    items.forEach((it) => {
      map[it.id] = { name: it.name, kb_code: it.kb_code }
    })
    kbMap.value = map
  } catch (e) {
    kbMap.value = {};
  }
};

const getKBName = (kbId) => {
  if (kbId == null || kbId === '' || kbId === 0) return ''
  const info = kbMap.value[kbId]
  if (!info) return ''
  return info.kb_code ? `[${info.kb_code}] ${info.name}` : info.name
}

const categoryOptions = computed(() => {
  const set = new Set()
  list.value.forEach((it) => { if (it.category) set.add(it.category) })
  return Array.from(set)
});

const enableRate = computed(() => {
  if (!stats.value?.total) return '0%'
  return `${Math.round((stats.value.enabled / stats.value.total) * 100)}%`
})

const categoryCount = computed(() => {
  const set = new Set()
  list.value.forEach((it) => { if (it.category) set.add(it.category) })
  return set.size
})

const testDialogVisible = ref(false);
const testLoading = ref(false)
const testTouched = ref(false)
const testForm = ref({ msg: '' })
const testResult = ref([])

const truncate = (text, len) => {
  if (!text) return '-'
  return text.length > len ? text.slice(0, len) + '...' : text
};

const formatTime = (t) => {
  if (!t) return '-'
  const d = new Date(t)
  if (isNaN(d.getTime())) return '-'
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const getConfColor = (c) => {
  if (c >= 0.8) return '#67c23a'
  if (c >= 0.6) return '#e6a23c'
  return '#f56c6c'
}

const getMatchTagType = (s) => {
  if (s >= 0.8) return 'success'
  if (s >= 0.6) return 'warning'
  return 'info'
}

const loadList = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    if (filter.value.keyword) params.keyword = filter.value.keyword
    if (filter.value.category) params.category = filter.value.category
    if (filter.value.intent) params.intent = filter.value.intent
    if (filter.value.enabled !== '' && filter.value.enabled !== null) {
      params.enabled = filter.value.enabled
    }
    const res = await faqApi.list(params)
    list.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) {
    ElMessage.error(t('auto.List.db4e4c72') + (e.message || t('auto.List.974e7484')))
  } finally {
    loading.value = false
  }
};

const loadStats = async () => {
  try {
    const res = await faqApi.stats()
    stats.value = res || { total: 0, enabled: 0 }
  } catch (e) {
    console.warn(t('auto.List.f1051bb8'), e?.message);
  }
}

const onSearch = () => {
  page.value = 1
  loadList()
}

const resetFilter = () => {
  filter.value = {
    keyword: '',
    category: '',
    intent: '',
    enabled: ''
  }
  page.value = 1
  loadList()
}

const goCreate = () => {
  router.push('/faq/editor')
};

const goEdit = (row) => {
  router.push(`/faq/editor/${row.id}`)
}

const goTest = () => {
  testTouched.value = false
  testResult.value = []
  testForm.value = { msg: '' }
  testDialogVisible.value = true
}

const onToggle = async (row, v) => {
  row._toggling = true
  try {
    await faqApi.update(row.id, { ...stripToPayload(row), enabled: v })
    ElMessage.success(v ? t('auto.List.53ace430') : t('auto.List.1c1ed981'))
    await loadList()
    await loadStats()
  } catch (e) {
    ElMessage.error(t('auto.List.428b3b05') + (e.message || t('auto.List.974e7484')))
  } finally {
    row._toggling = false
  }
};

const stripToPayload = (row) => ({
  question: row.question,
  answer: row.answer,
  keywords: row.keywords || [],
  category: row.category || '',
  intent: row.intent || '',
  confidence: row.confidence || 0
})

const onDelete = (row) => {
  ElMessageBox.confirm(`确认删除 FAQ「${truncate(row.question, 30)}」吗？删除后不可恢复。`, t('auto.List.50eaf94d'), {
    type: 'warning',
    confirmButtonText: t('auto.List.631cd220'),
    cancelButtonText: t('auto.List.625fb26b')
  }).then(async () => {
    try {
      await faqApi.remove(row.id)
      ElMessage.success(t('auto.List.0007d170'))
      await loadList()
      await loadStats()
    } catch (e) {
      ElMessage.error(t('auto.List.bc42bb8f') + (e.message || t('auto.List.974e7484')))
    }
  }).catch(() => {})
};

const onTestOne = (row) => {
  testTouched.value = true
  testForm.value.msg = row.question
  testDialogVisible.value = true
  runTest()
};

const runTest = async () => {
  if (!testForm.value.msg || !testForm.value.msg.trim()) {
    ElMessage.warning(t('auto.List.2b216e37'))
    return
  }
  testLoading.value = true
  testTouched.value = true
  try {
    const res = await faqApi.match({ msg: testForm.value.msg, top_k: 5 })
    testResult.value = res || []
  } catch (e) {
    testResult.value = []
    ElMessage.error(t('auto.List.a4ccd593') + (e.message || t('auto.List.974e7484')))
  } finally {
    testLoading.value = false
  }
};

onMounted(() => {
  loadKBMap()
  loadList()
  loadStats()
})
</script>

<style scoped>
.faq-list-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.header-card :deep(.el-card__body) {
  padding: 18px 24px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
}

.header-content h2 {
  margin: 0 0 4px 0;
  font-size: 20px;
  font-weight: 600;
}

.subtitle {
  margin: 0;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.stat-row .stat-card :deep(.el-card__body) {
  text-align: center;
  padding: 16px;
}

.stat-label {
  color: var(--el-text-color-secondary);
  font-size: 13px;
  margin-bottom: 6px;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.text-success { color: var(--el-color-success); }
.text-primary { color: var(--el-color-primary); }
.text-info { color: var(--el-color-info); }
.text-muted { color: #c0c4cc; font-size: 12px; }

.filter-card :deep(.el-card__body) {
  padding: 16px 20px 0 20px;
}

.question-text {
  font-weight: 500;
  margin-bottom: 4px;
  color: var(--el-text-color-primary);
}

.keyword-list {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 4px;
}

.keyword-tag {
  font-size: 11px;
}

.keyword-more {
  color: var(--el-text-color-secondary);
  font-size: 11px;
  margin-left: 4px;
}

.pagination-wrap {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
