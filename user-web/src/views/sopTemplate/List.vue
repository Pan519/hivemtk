<template>
  <div class="sop-template-list-page">
    
    <el-card class="header-card" shadow="never">
      <div class="header-content">
        <div>
          <h2>{{ $t('auto.List.b9ce6f82') }}</h2>
          <p class="subtitle">{{ $t('auto.List.68fe9153') }}</p>
        </div>
        <div class="header-actions">
          <el-button @click="loadList" :loading="loading">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
          <el-button type="primary" @click="goCreate">
            <el-icon><Plus /></el-icon>
            新建模板
          </el-button>
        </div>
      </div>
    </el-card>

    
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filter" @submit.prevent>
        <el-form-item :label="$t('auto.List.d7ec2d3f')">
          <el-input
            v-model="filter.keyword"
            :placeholder="$t('auto.List.91a34091')"
            clearable
            style="width: 200px"
            @keyup.enter="onSearch"
            @clear="onSearch"
          />
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
        <el-form-item :label="$t('auto.List.d316330b')">
          <el-select v-model="filter.stage" :placeholder="$t('auto.List.7f3b19c5')" clearable style="width: 160px" @change="onSearch">
            <el-option v-for="s in stageOptions" :key="s" :label="s" :value="s" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.3fea7ca7')">
          <el-select v-model="filter.enabled" :placeholder="$t('auto.List.a8b0c204')" clearable style="width: 130px" @change="onSearch">
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
        <el-table-column :label="$t('auto.List.a5d1c511')" min-width="180" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="tpl-name">{{ row.name }}</div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.764aca4e')" width="120" align="center">
          <template #default="{ row }">
            <el-tag size="small" type="warning">{{ row.intent || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.d316330b')" width="100" align="center">
          <template #default="{ row }">
            <el-tag size="small" type="info">{{ row.stage || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.de92100c')" width="180" show-overflow-tooltip>
          <template #default="{ row }">
            <el-tag
              v-if="getKBName(row.kb_id)"
              size="small"
              type="warning"
              effect="plain"
            >
              {{ getKBName(row.kb_id) }}
            </el-tag>
            <span v-else class="text-muted">{{ $t('auto.List.f64a3850') }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.03ae7940')" min-width="240" show-overflow-tooltip>
          <template #default="{ row }">
            <code class="tpl-code">{{ truncate(row.template, 90) }}</code>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.ee8ecb9e')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.priority >= 50 ? 'success' : row.priority >= 0 ? 'info' : 'warning'" size="small">
              {{ row.priority || 0 }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.8f3599b6')" width="100" align="center">
          <template #default="{ row }">
            <el-progress
              :percentage="Math.round((row.confidence || 0) * 100)"
              :stroke-width="10"
              :color="getConfColor(row.confidence)"
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
        <el-table-column :label="$t('auto.List.2b6bc0f2')" width="160" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="goEdit(row)">{{ $t('auto.List.95b351c8') }}</el-button>
            <el-button link type="danger" size="small" @click="onDelete(row)">{{ $t('auto.List.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty :description="$t('auto.List.ecc8222a')" />
        </template>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadList"
          @current-change="loadList"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Plus, Search } from '@element-plus/icons-vue'
import { sopTemplateApi } from '@/api/sopTemplate'
import { listKBs } from '@/api/knowledgeBase'

const router = useRouter()

const { t } = useI18n()

const loading = ref(false);
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const filter = ref({
  keyword: '',
  intent: '',
  stage: '',
  enabled: '',
  page: 1,
  page_size: 20
})

const stageOptions = ['initial', 'middle', 'late', 'objection', 'closing'];

const kbMap = ref({});
const loadKBMap = async () => {
  try {
    const res = await listKBs({ kb_type: 'sop', page: 1, page_size: 200 }).catch(() => null)
    const items = Array.isArray(res) ? res : res?.list || res?.items || []
    const map = {}
    items.forEach((it) => {
      map[it.id] = { name: it.name, kb_code: it.kb_code }
    })
    kbMap.value = map
  } catch {
    kbMap.value = {}
  }
}
const getKBName = (kbId) => {
  if (kbId == null || kbId === '' || kbId === 0) return ''
  const info = kbMap.value[kbId]
  if (!info) return ''
  return info.kb_code ? `[${info.kb_code}] ${info.name}` : info.name
}

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

const loadList = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    if (filter.value.keyword) params.keyword = filter.value.keyword
    if (filter.value.intent) params.intent = filter.value.intent
    if (filter.value.stage) params.stage = filter.value.stage
    if (filter.value.enabled !== '' && filter.value.enabled !== null) {
      params.enabled = filter.value.enabled
    }
    const res = await sopTemplateApi.list(params)
    list.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) {
    ElMessage.error(t('auto.List.3a1ee2d2') + (e.message || t('auto.List.974e7484')))
  } finally {
    loading.value = false
  }
};

const onSearch = () => {
  page.value = 1
  loadList()
}

const resetFilter = () => {
  filter.value = {
    keyword: '',
    intent: '',
    stage: '',
    enabled: '',
    page: 1,
    page_size: pageSize.value
  }
  page.value = 1
  loadList()
}

const goCreate = () => {
  router.push('/sop-template/editor')
};

const goEdit = (row) => {
  router.push(`/sop-template/editor/${row.id}`)
}

const onToggle = async (row, v) => {
  row._toggling = true
  try {
    await sopTemplateApi.update(row.id, {
      name: row.name,
      intent: row.intent,
      stage: row.stage,
      template: row.template,
      vars: row.vars || '',
      priority: row.priority || 0,
      confidence: row.confidence || 0,
      enabled: v
    })
    ElMessage.success(v ? t('auto.List.53ace430') : t('auto.List.1c1ed981'))
    await loadList()
  } catch (e) {
    ElMessage.error(t('auto.List.428b3b05') + (e.message || t('auto.List.974e7484')))
  } finally {
    row._toggling = false
  }
};

const onDelete = (row) => {
  ElMessageBox.confirm(`确认删除模板「${row.name}」吗？删除后不可恢复。`, t('auto.List.50eaf94d'), {
    type: 'warning',
    confirmButtonText: t('auto.List.631cd220'),
    cancelButtonText: t('auto.List.625fb26b')
  }).then(async () => {
    try {
      await sopTemplateApi.remove(row.id)
      ElMessage.success(t('auto.List.0007d170'))
      await loadList()
    } catch (e) {
      ElMessage.error(t('auto.List.bc42bb8f') + (e.message || t('auto.List.974e7484')))
    }
  }).catch(() => {})
};

onMounted(() => {
  loadKBMap()
  loadList()
})
</script>

<style scoped>
.sop-template-list-page {
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

.filter-card :deep(.el-card__body) {
  padding: 16px 20px 0 20px;
}

.tpl-name {
  font-weight: 500;
  color: var(--el-text-color-regular);
}
.tpl-code {
  font-family: 'SF Mono', Monaco, Consolas, monospace;
  font-size: 12px;
  color: #606266;
  background: var(--el-fill-color-light);
  padding: 2px 6px;
  border-radius: 3px;
}
.text-muted { color: #c0c4cc; font-size: 12px; }

.pagination-wrap {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
